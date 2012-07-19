// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"bytes"
	"encoding/json"
	"errors"
	// "github.com/kr/pretty"
	"io/ioutil"
	// "log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	BadResponse        = errors.New("Bad response from Neo4j server.")
	NotFound           = errors.New("Cannot find in database.")
	FeatureUnavailable = errors.New("Feature unavailable")
	CannotDelete       = errors.New("The node cannot be deleted. Check that the node is orphaned before deletion.")
)

// An errorResponse is returned from the Neo4j server on errors.
type errorResponse struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
}

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	url    *url.URL // Root URL for REST API
	client *http.Client
	info   *serviceRootInfo
}

// A serviceRootInfo is returned from the Neo4j server on successful operations 
// involving a Node.
type serviceRootInfo struct {
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
	// log.Println(pretty.Sprintf("Response: %# v", resp))
	// Only try to unmarshal if status is 200 OK or 201 CREATED
	if status >= 200 && status <= 201 {
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(data, &r.Result)
		if err != nil {
			return
		}
		// log.Println(pretty.Sprintf("Result: %# v", r.Result))
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
func Connect(uri string) (db *Database, err error) {
	var info serviceRootInfo
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	db = &Database{
		url:    u,
		client: new(http.Client),
		info:   &info,
	}
	c := restCall{
		Url:    u.String(),
		Method: "GET",
		Result: &info,
	}
	code, err := db.rest(&c)
	if err != nil || code != 200 {
		return
	}
	if db.info.Version == "" {
		err = BadResponse
		return
	}
	return
}

// CreateNode creates a Node in the database.
func (db *Database) CreateNode(p Properties) (*Node, error) {
	var info neoInfo
	n := Node{neoEntity{
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
	var info neoInfo
	n := Node{neoEntity{
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
	var info neoInfo
	rel := Relationship{neoEntity{
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
	if code != 200 {
		err = BadResponse
	}
	return &rel, err
}

// RelationshipTypes gets all existing relationship types from the DB
func (db *Database) RelationshipTypes() ([]string, error) {
	ts := []string{}
	if db.info.RelTypes == "" {
		return ts, FeatureUnavailable
	}
	c := restCall{
		Url:    db.info.RelTypes,
		Method: "GET",
		Result: &ts,
	}
	code, err := db.rest(&c)
	if err != nil {
		return ts, err
	}
	if code == 200 {
		return ts, nil // Success!
	}
	return ts, BadResponse
}

// Properties is a bag of key/value pairs that can describe Nodes
// and Relationships.
type Properties map[string]string

// neoEntity is the base type for Nodes and Relationships
type neoEntity struct {
	info *neoInfo
	db   *Database
}

// A neoInfo is returned from the Neo4j server on successful operations 
// involving a Node or a Relationship
type neoInfo struct {
	//
	// Always filled
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

////////////////////////////////////////////////////////////////////////////////
//
// These operations can be performed on both Nodes and Relationships using 
// the same procedure with different URLs supplied in the neoInfo argument.
//
////////////////////////////////////////////////////////////////////////////////

// SetProperty sets the single property key to value.
func (e *neoEntity) SetProperty(key string, value string) error {
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
func (e *neoEntity) GetProperty(key string) (string, error) {
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
func (e *neoEntity) DeleteProperty(key string) error {
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
func (e *neoEntity) Delete() error {
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
func (e *neoEntity) Properties() (Properties, error) {
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
func (e *neoEntity) SetProperties(p Properties) error {
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
func (e *neoEntity) DeleteProperties() error {
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

////////////////////////////////////////////////////////////////////////////////
//
// Node
//
////////////////////////////////////////////////////////////////////////////////

// A node in a Neo4j database
type Node struct {
	neoEntity
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.db.info.Node)
	s := n.info.Self[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		// Are both n.Info and n.Node valid?
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
	s := []neoInfo{}
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
		rel := Relationship{neoEntity{
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
	var info neoInfo
	rel := Relationship{neoEntity{
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

// A relationship in a Neo4j database
type Relationship struct {
	neoEntity
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
