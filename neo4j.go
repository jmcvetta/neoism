// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	BadResponse = errors.New("Bad response from Neo4j server.")
	NodeNotFound = errors.New("Cannot find node in database.")
	FeatureUnavailable = errors.New("Feature unavailable")
	CannotDelete = errors.New("The node cannot be deleted. Check that the node is orphaned before deletion.")
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
	Url    string      // Absolute URL to call
	Method string      // HTTP method to use 
	Data   interface{} // Data to JSON-encode and include with call
	Result interface{} // JSON-encoded data in respose will be unmarshalled into Result
}

func (db *Database) rest(r *restCall) (status int, err error) {
	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return 
	}
	if r.Data != nil {
		var b []byte
		b, err = json.Marshal(r.Data)
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
	log.Println(pretty.Sprintf("Request: %# v", req))
	resp, err := db.client.Do(req)
	if err != nil {
		return 
	}
	status = resp.StatusCode
	log.Println(pretty.Sprintf("Response: %# v", resp))
	// Only try to unmarshal if status is 200 OK or 201 CREATED
	if status >= 200 && status <= 201 {
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(data, &r.Result)
		if err != nil {
			return 
		}
		log.Println(pretty.Sprintf("Result: %# v", r.Result))
		}
	return
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
func (db *Database) CreateNode(props map[string]string) (*Node, error) {
	n := Node{
		Db: db,
	}
	var info nodeInfo
	c := restCall{
		Url:    db.Info.Node,
		Method: "POST",
		Data:   &props,
		Result: &info,
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
	n := Node{
		Db: db,
	}
	var info nodeInfo
	parts := []string{db.Info.Node, strconv.Itoa(id)}
	uri := strings.Join(parts, "/")
	c := restCall{
		Url:    uri,
		Method: "GET",
		Result: &info,
	}
	code, err := db.rest(&c)
	switch {
	case code == 404:
		return &n, NodeNotFound
	case code != 200:
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
func (n *Node) Properties() (map[string]string, error) {
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
}
