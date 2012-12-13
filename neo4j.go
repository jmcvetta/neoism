// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	// "github.com/kr/pretty"
	"github.com/jmcvetta/restclient"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var (
	InvalidDatabase    = errors.New("Invalid database.  Check URI.")
	BadResponse        = errors.New("Bad response from Neo4j server.")
	NotFound           = errors.New("Cannot find in database.")
	FeatureUnavailable = errors.New("Feature unavailable")
	CannotDelete       = errors.New("The node cannot be deleted. Check that the node is orphaned before deletion.")
)

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	url    *url.URL // Root URL for REST API
	client *http.Client
	rc     *restclient.Client
	info   *serviceRootInfo
}

// A neoError is populated by api calls when there is an error.
type neoError struct {
	Mesage     string   `json:"message"`
	Exception  string   `json:"exception"`
	StackTrace []string `json:"stacktrace"`
}

// A serviceRootInfo describes services available on the Neo4j server
type serviceRootInfo struct {
	neoError
	Extensions interface{} `json:"extensions"`
	Node       string      `json:"node"`
	RefNode    string      `json:"reference_node"`
	NodeIndex  string      `json:"node_index"`
	RelIndex   string      `json:"relationship_index"`
	ExtInfo    string      `json:"extensions_info"`
	RelTypes   string      `json:"relationship_types"`
	Batch      string      `json:"batch"`
	Cypher     string      `json:"cypher"`
	Version    string      `json:"neo4j_version"`
}

type restCall struct {
	Url     string      // Absolute URL to call
	Method  string      // HTTP method to use 
	Content interface{} // Data to JSON-encode and include with call
	Result  interface{} // JSON-encoded data in respose will be unmarshalled into Result
}

func (db *Database) rest(r *restCall) (status int, err error) {
	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return
	}
	if r.Content != nil {
		// log.Println(pretty.Sprintf("Content: %# v", r.Content))
		var b []byte
		b, err = json.Marshal(r.Content)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(r.Method, r.Url, buf)
		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Accept", "application/json")
	// log.Println(pretty.Sprintf("Request: %# v", req))
	resp, err := db.client.Do(req)
	if err != nil {
		return
	}
	status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	// Ignore unmarshall errors - worst case is, r.Result will be nil
	json.Unmarshal(data, &r.Result)
	if status < 200 || status >= 300 {
		res := &r.Result
		// log.Println(*res)
		info, ok := (*res).(neoError)
		if ok {
			log.Println("Got error response code:", status)
			log.Println(info.Mesage)
			log.Println(info.Exception)
			log.Println(info.StackTrace)
		}
	}
	return
}

// Joins URL fragments
func join(fragments ...string) string {
	parts := []string{}
	for _, v := range fragments {
		v = strings.Trim(v, "/")
		parts = append(parts, v)
	}
	return strings.Join(parts, "/")
}
func Connect(uri string) (*Database, error) {
	var info serviceRootInfo
	db := &Database{
		client: new(http.Client),
		rc:     restclient.New(),
		info:   &info,
	}
	u, err := url.Parse(uri)
	if err != nil {
		return db, err
	}
	db.url = u
	r := restclient.RestRequest{
		Url:    u.String(),
		Method: restclient.GET,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := db.rc.Do(&r)
	if err != nil {
		log.Println(info.Mesage)
		log.Println(info.Exception)
		log.Println(info.StackTrace)
		return db, err
	}
	switch {
	case status == 200 && db.info.Version != "":
		return db, nil // Success!
	case status == 404:
		return db, InvalidDatabase
	}
	log.Println(info.Mesage)
	log.Println(info.Exception)
	log.Println(info.StackTrace)
	return db, BadResponse
}

// CreateNode creates a Node in the database.
func (db *Database) CreateNode(p Properties) (*Node, error) {
	var info nrInfo
	n := Node{nrBase{
		db:   db,
		info: &info,
	}}
	c := restCall{
		Url:     db.info.Node,
		Method:  "POST",
		Content: &p,
		Result:  &info,
	}
	code, err := db.rest(&c)
	if err != nil || code != 201 {
		return &n, err
	}
	if info.Self == "" {
		return &n, BadResponse
	}
	return &n, nil
}

// GetNode fetches a Node from the database
func (db *Database) GetNode(id int) (*Node, error) {
	uri := join(db.info.Node, strconv.Itoa(id))
	return db.getNodeByUri(uri)
}

// GetNode fetches a Node from the database based on its uri
func (db *Database) getNodeByUri(uri string) (*Node, error) {
	var info nrInfo
	n := Node{nrBase{
		db:   db,
		info: &info,
	}}
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &info,
	}
	code, err := db.rest(&c)
	switch {
	case code == 404:
		return &n, NotFound
	case code != 200 || info.Self == "":
		return &n, BadResponse
	}
	if err != nil {
		return &n, err
	}
	// n.Info = &info
	return &n, nil
}

// GetRelationship fetches a Relationship from the DB by id.
func (db *Database) GetRelationship(id int) (*Relationship, error) {
	var info nrInfo
	rel := Relationship{nrBase{
		db:   db,
		info: &info,
	}}
	uri := join(db.url.String(), "relationship", strconv.Itoa(id))
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &info,
	}
	code, err := db.rest(&c)
	switch code {
	default:
		err = BadResponse
	case 200:
		err = nil // Success!
	case 404:
		err = NotFound
	}
	return &rel, err
}

// RelationshipTypes gets all existing relationship types from the DB
func (db *Database) RelationshipTypes() ([]string, error) {
	reltypes := []string{}
	if db.info.RelTypes == "" {
		return reltypes, FeatureUnavailable
	}
	c := restCall{
		Url:    db.info.RelTypes,
		Method: "GET",
		Result: &reltypes,
	}
	code, err := db.rest(&c)
	if err != nil {
		return reltypes, err
	}
	if code == 200 {
		return reltypes, nil // Success!
	}
	sort.Sort(sort.StringSlice(reltypes))
	return reltypes, BadResponse
}

// CreateIndex creates a new Index, with the name supplied, in the db.
func (db *Database) CreateIndex(name string) (*Index, error) {
	conf := IndexConfig{
		Name: name,
	}
	return db.CreateIndexFromConf(conf)
}

// CreateIndexFromConf creates a new Index based on an IndexConfig object
func (db *Database) CreateIndexFromConf(conf IndexConfig) (*Index, error) {
	var info indexInfo
	i := Index{
		db:   db,
		info: &info,
	}
	c := restCall{
		Url:     db.info.NodeIndex,
		Method:  "POST",
		Content: &conf,
		Result:  &info,
	}
	code, err := db.rest(&c)
	if err != nil {
		return &i, err
	}
	if code != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", code)
		log.Printf("    Result:", info)
		return &i, BadResponse
	}
	return &i, nil
}

// Properties is a bag of key/value pairs that can describe Nodes
// and Relationships.
type Properties map[string]string

// nrBase is the base type for Nodes and Relationships.
type nrBase struct {
	db   *Database
	info *nrInfo
}

// A nrInfo is returned from the Neo4j server on successful operations 
// involving a Node or a Relationship.
type nrInfo struct {
	neoError
	//
	// Always filled on success
	//
	Property   string      `json:"property"`
	Properties string      `json:"properties"`
	Self       string      `json:"self"`
	Data       interface{} `json:"data"`
	Extensions interface{} `json:"extensions"`
	//
	// Filled only for Node operations
	//
	OutgoingRels      string `json:"outgoing_relationships"`
	Traverse          string `json:"traverse"`
	AllTypedRels      string `json:"all_typed_relationships"`
	Outgoing          string `json:"outgoing_typed_relationships"`
	IncomingRels      string `json:"incoming_relationships"`
	CreateRel         string `json:"create_relationship"`
	PagedTraverse     string `json:"paged_traverse"`
	AllRels           string `json:"all_relationships"`
	IncomingTypedRels string `json:"incoming_typed_relationships"`
	//
	// Filled only for Relationship operations
	//
	Start string `json:"start"`
	Type  string `json:"type"`
	End   string `json:"end"`
}

// SetProperty sets the single property key to value.
func (e *nrBase) SetProperty(key string, value string) error {
	uri := e.info.Properties
	if uri == "" {
		return FeatureUnavailable
	}
	parts := []string{uri, key}
	uri = strings.Join(parts, "/")
	c := restCall{
		Url:     uri,
		Method:  "PUT",
		Content: &value,
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return err
	}
	if code == 204 {
		return nil // Success!
	}
	return BadResponse
}

// GetProperty fetches the value of property key.
func (e *nrBase) GetProperty(key string) (string, error) {
	var val string
	uri := e.info.Properties
	if uri == "" {
		return val, FeatureUnavailable
	}
	parts := []string{uri, key}
	uri = strings.Join(parts, "/")
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &val,
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return val, err
	}
	switch code {
	case 200:
		return val, nil
	case 404:
		return val, NotFound
	}
	return val, BadResponse
}

// DeleteProperty deletes property key
func (e *nrBase) DeleteProperty(key string) error {
	uri := e.info.Properties
	if uri == "" {
		return FeatureUnavailable
	}
	parts := []string{uri, key}
	uri = strings.Join(parts, "/")
	c := restCall{
		Url:    uri,
		Method: "DELETE",
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return err
	}
	switch code {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	return BadResponse
}

// Delete removes the object from the DB.
func (e *nrBase) Delete() error {
	uri := e.info.Self
	if uri == "" {
		return FeatureUnavailable
	}
	c := restCall{
		Url:    uri,
		Method: "DELETE",
	}
	code, err := e.db.rest(&c)
	switch {
	case err != nil:
		return err
	case code == 204:
		// Successful deletion!
		return nil
	case code == 409:
		return CannotDelete
	}
	return BadResponse
}

// Properties fetches all properties
func (e *nrBase) Properties() (Properties, error) {
	props := make(map[string]string)
	uri := e.info.Properties
	if uri == "" {
		return props, FeatureUnavailable
	}
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &props,
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if code == 204 {
		props = map[string]string{}
	}
	return props, nil
}

// SetProperties updates all properties, overwriting any existing properties.
func (e *nrBase) SetProperties(p Properties) error {
	uri := e.info.Properties
	if uri == "" {
		return FeatureUnavailable
	}
	c := restCall{
		Url:     uri,
		Method:  "PUT",
		Content: &p,
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return err
	}
	if code == 204 {
		return nil // Success!
	}
	return BadResponse
}

// DeleteProperties deletes all properties.
func (e *nrBase) DeleteProperties() error {
	uri := e.info.Properties
	if uri == "" {
		return FeatureUnavailable
	}
	c := restCall{
		Url:    uri,
		Method: "DELETE",
	}
	code, err := e.db.rest(&c)
	if err != nil {
		return err
	}
	switch code {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	return BadResponse
}

/*******************************************************************************
 *
 * Node
 *
 ******************************************************************************/

// A node in a Neo4j database
type Node struct {
	nrBase
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.db.info.Node)
	s := n.info.Self[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return id
}

// getRelationships makes an api call to the supplied uri and returns a map 
// keying relationship IDs to Relationship objects.
func (n *Node) getRelationships(uri string, types ...string) (map[int]Relationship, error) {
	m := map[int]Relationship{}
	if uri == "" {
		return m, FeatureUnavailable
	}
	if types != nil {
		fragment := strings.Join(types, "&")
		parts := []string{uri, fragment}
		uri = strings.Join(parts, "/")
	}
	s := []nrInfo{}
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &s,
	}
	code, err := n.db.rest(&c)
	if err != nil {
		return m, err
	}
	for _, info := range s {
		rel := Relationship{nrBase{
			db:   n.db,
			info: &info,
		}}
		m[rel.Id()] = rel
	}
	if code == 200 {
		return m, nil // Success!
	}
	return m, BadResponse
}

// Relationships gets all Relationships for this Node, optionally filtered by 
// type, returning them as a map keyed on Relationship ID.
func (n *Node) Relationships(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.AllRels, types...)
}

// Incoming gets all incoming Relationships for this Node.
func (n *Node) Incoming(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.IncomingRels, types...)
}

// Outgoing gets all outgoing Relationships for this Node.
func (n *Node) Outgoing(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.OutgoingRels, types...)
}

// Relate creates a relationship of relType, with specified properties, 
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Properties) (*Relationship, error) {
	var info nrInfo
	rel := Relationship{nrBase{
		db:   n.db,
		info: &info,
	}}
	srcUri := join(n.info.Self, "relationships")
	destUri := join(n.db.info.Node, strconv.Itoa(destId))
	content := map[string]interface{}{
		"to":   destUri,
		"type": relType,
	}
	if p != nil {
		content["data"] = &p
	}
	c := restCall{
		Url:     srcUri,
		Method:  "POST",
		Content: content,
		Result:  &info,
	}
	code, err := n.db.rest(&c)
	if err != nil {
		return &rel, err
	}
	if code != 201 {
		return &rel, BadResponse
	}
	return &rel, nil
}

/*******************************************************************************
 *
 * Relationship
 *
 ******************************************************************************/

// A relationship in a Neo4j database
type Relationship struct {
	nrBase
}

// Id gets the ID number of this Relationship
func (r *Relationship) Id() int {
	parts := strings.Split(r.info.Self, "/")
	s := parts[len(parts)-1]
	id, err := strconv.Atoi(s)
	if err != nil {
		// Are both r.Info and r.Node valid?
		panic(err)
	}
	return id
}

// Start gets the starting Node of this Relationship.
func (r *Relationship) Start() (*Node, error) {
	// log.Println("INFO", r.Info)
	return r.db.getNodeByUri(r.info.Start)
}

// End gets the ending Node of this Relationship.
func (r *Relationship) End() (*Node, error) {
	return r.db.getNodeByUri(r.info.End)
}

// Type gets the type of this relationship
func (r *Relationship) Type() string {
	return r.info.Type
}

/*******************************************************************************
 *
 * Index
 *
 ******************************************************************************/

type IndexConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
}

// An indexInfo is returned from the Neo4j server on operations involving an Index.
type indexInfo struct {
	neoError
	Template string `json:"template"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
}

type Index struct {
	db   *Database
	info *indexInfo
}
