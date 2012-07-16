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
	Info   *serviceRootInfo
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
func NewDatabase(uri string) (db *Database, err error) {
	var info serviceRootInfo
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	db = &Database{
		url:    u,
		client: new(http.Client),
		Info:   &info,
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
	if db.Info.Version == "" {
		err = BadResponse
		return
	}
	return
}

//
// Nodes
//

// A node in a Neo4j database
type Node struct {
	Info *nodeInfo
	Db   *Database
}

type Properties map[string]string

// A nodeInfo is returned from the Neo4j server on successful operations 
// involving a Node.
type nodeInfo struct {
	OutgoingRels      string      `json:"outgoing_relationships"`
	Data              interface{} `json:"data"`
	Traverse          string      `json:"traverse"`
	AllTypedRels      string      `json:"all_typed_relationships"`
	Property          string      `json:"property"`
	Self              string      `json:"self"`
	Outgoing          string      `json:"outgoing_typed_relationships"`
	Properties        string      `json:"properties"`
	IncomingRels      string      `json:"incoming_relationships"`
	Extensions        interface{} `json:"extensions"`
	CreateRel         string      `json:"create_relationship"`
	PagedTraverse     string      `json:"paged_traverse"`
	AllRels           string      `json:"all_relationships"`
	IncomingTypedRels string      `json:"incoming_typed_relationships"`
}

// CreateNode creates a Node in the database.
func (db *Database) CreateNode(p Properties) (*Node, error) {
	n := Node{
		Db: db,
	}
	var info nodeInfo
	c := restCall{
		Url:     db.Info.Node,
		Method:  "POST",
		Content: &p,
		Result:  &info,
	}
	code, err := db.rest(&c)
	if err != nil || code != 201 {
		return &n, err
	}
	n.Info = &info
	if n.Info.Self == "" {
		return &n, BadResponse
	}
	return &n, nil
}

// GetNode fetches a Node from the database
func (db *Database) GetNode(id int) (*Node, error) {
	uri := join(db.Info.Node, strconv.Itoa(id))
	return db.getNodeByUri(uri)
}

// GetNode fetches a Node from the database based on its uri
func (db *Database) getNodeByUri(uri string) (*Node, error) {
	n := Node{
		Db: db,
	}
	var info nodeInfo
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
	n.Info = &info
	return &n, nil
}

// Delete deletes a Node from the database
func (n *Node) Delete() error {
	c := restCall{
		Url:    n.Info.Self,
		Method: "DELETE",
	}
	code, err := n.Db.rest(&c)
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

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.Db.Info.Node)
	s := n.Info.Self[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		// Are both n.Info and n.Node valid?
		panic(err)
	}
	return id
}

// Properties gets the Node's properties map from the DB.
func (n *Node) Properties() (Properties, error) {
	props := make(map[string]string)
	if n.Info.Properties == "" {
		return props, FeatureUnavailable
	}
	c := restCall{
		Url:    n.Info.Properties,
		Method: "GET",
		Result: &props,
	}
	code, err := n.Db.rest(&c)
	if err != nil {
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if code == 204 {
		props = map[string]string{}
	}
	return props, nil
}

//
// Relationships
//

// A relInfo is returned from the Neo4j server on successful operations 
// involving a Relationship.
type relInfo struct {
	Start      string            `json:"start"`
	Data       map[string]string `json:"data"`
	Self       string            `json:"self"`
	Property   string            `json:"property"`
	Properties string            `json:"properties"`
	Type       string            `json:"type"`
	Extensions map[string]string `json:"extensions"`
	End        string            `json:"end"`
}

// A relationship in a Neo4j database
type Relationship struct {
	Info *relInfo
	Db   *Database
}

// Relate creates a relationship of relType, with specified properties, 
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Properties) (*Relationship, error) {
	var info relInfo
	rel := Relationship{
		Db:   n.Db,
		Info: &info,
	}
	srcUri := join(n.Info.Self, "relationships")
	destUri := join(n.Db.Info.Node, strconv.Itoa(destId))
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
	code, err := n.Db.rest(&c)
	if err != nil {
		return &rel, err
	}
	if code != 201 {
		return &rel, BadResponse
	}
	return &rel, nil
}

// Start gets the starting Node of this Relationship.
func (r *Relationship) Start() (*Node, error) {
	// log.Println("INFO", r.Info)
	return r.Db.getNodeByUri(r.Info.Start)
}

// End gets the ending Node of this Relationship.
func (r *Relationship) End() (*Node, error) {
	return r.Db.getNodeByUri(r.Info.End)
}

// Type gets the type of this relationship
func (r *Relationship) Type() string {
	return r.Info.Type
}

// Properties gets the Relationship's properties map from the DB.
func (r *Relationship) Properties() (Properties, error) {
	props := Properties{}
	if r.Info.Properties == "" {
		return props, FeatureUnavailable
	}
	c := restCall{
		Url:    r.Info.Properties,
		Method: "GET",
		Result: &props,
	}
	// code, err := r.Db.rest(&c)
	_, err := r.Db.rest(&c)
	if err != nil {
		return props, err
	}
	/*
		// Status code 204 indicates no properties on this Relationship
		if code == 204 {
			props = map[string]string{}
		}
	*/
	return props, nil
}

// GetRelationship fetches a Relationship from the DB by id.
func (db *Database) GetRelationship(id int) (*Relationship, error) {
	var info relInfo
	rel := Relationship{
		Db:   db,
		Info: &info,
	}
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

// Id gets the ID number of this Relationship
func (r *Relationship) Id() int {
	parts := strings.Split(r.Info.Self, "/")
	s := parts[len(parts)-1]
	id, err := strconv.Atoi(s)
	if err != nil {
		// Are both r.Info and r.Node valid?
		panic(err)
	}
	return id
}

// Delete deletes a Relationship from the database
func (r *Relationship) Delete() error {
	c := restCall{
		Url:    r.Info.Self,
		Method: "DELETE",
	}
	code, err := r.Db.rest(&c)
	switch {
	case err != nil:
		return err
	case code == 204:
		// Successful deletion!
		return nil
		/*
			case code == 409:
				return CannotDelete
		*/
	}
	return BadResponse
}

// SetProperties sets all properties on a Relationship, overwriting any
// existing properties.
func (r *Relationship) SetProperties(p Properties) error {
	c := restCall{
		Url:     r.Info.Properties,
		Method:  "PUT",
		Content: &p,
	}
	code, err := r.Db.rest(&c)
	if err != nil {
		return err
	}
	if code == 204 {
		return nil // Success!
	}
	return BadResponse
}

// GetProperty retrieves the value for the named property
func (r *Relationship) GetProperty(key string) (string, error) {
	var val string
	parts := []string{r.Info.Properties, key}
	u := strings.Join(parts, "/")
	c := restCall{
		Url:    u,
		Method: "GET",
		Result: &val,
	}
	code, err := r.Db.rest(&c)
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

// SetProperty sets the value for the named property
func (r *Relationship) SetProperty(key, value string) error {
	parts := []string{r.Info.Properties, key}
	u := strings.Join(parts, "/")
	c := restCall{
		Url:     u,
		Method:  "PUT",
		Content: &value,
	}
	code, err := r.Db.rest(&c)
	if err != nil {
		return err
	}
	if code == 204 {
		return nil // Success!
	}
	return BadResponse
}
