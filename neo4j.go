// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	// "strings"
	"errors"
	"log"
)

var (
	BadResponse = errors.New("Bad response from Neo4j server.")
)

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

// A node in a Neo4j database
type Node struct {
	Id int64
}

// A nodeInfo is returned from the Neo4j server on successful operations 
// involving a Node.
type nodeInfo struct {
	OutgoingRels      string            `json:"outgoing_relationships"`
	Data              map[string]string `json:"data"`
	Traverse          string            `json:"traverse"`
	AllTypedRels      string            `json:"all_typed_relationships"`
	Property          string            `json:"property"`
	Self              string            `json:"self"`
	Outgoing          string            `json:"outgoing_typed_relationships"`
	Properties        string            `json:"properties"`
	IncomingRels      string            `json:"incoming_relationships"`
	Extensions        map[string]string `json:"extensions"`
	CreateRel         string            `json:"create_relationship"`
	PagedTraverse     string            `json:"paged_traverse"`
	AllRels           string            `json:"all_relationships"`
	IncomingTypedRels string            `json:"incoming_typed_relationships"`
}

// An errorResponse is returned from the Neo4j server on errors.
type errorResponse struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
}

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

type restCall struct {
	Url    string      // Absolute URL to call
	Method string      // HTTP method to use 
	Data   interface{} // Data to JSON-encode and include with call
	Result interface{} // JSON-encoded data in respose will be unmarshalled into Result
}

func (db *Database) rest(r *restCall) error {
	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return err
	}
	if r.Data != nil {
		var b []byte
		b, err = json.Marshal(r.Data)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(r.Method, r.Url, buf)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
	} 
	req.Header.Add("Accept", "application/json")
	log.Println("Request:", req)
	resp, err := db.client.Do(req)
	if err != nil {
		return err
	}
	log.Println("Response:", resp)
	data, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &r.Result)
	if err != nil {
		return err
	}
	log.Println("Result:", r.Result)
	log.Printf("type(Result): %T", r.Result)
	return nil
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
	err = db.rest(&c)
	if err != nil {
		return
	}
	if info.Version == "" {
		err = BadResponse
		return
	}
	db.Info = &info
	return
}

//
// Nodes
//

/*
func (db *Database) CreateNode(props map[string]string) (n *Node, err error) {
	fragment := "/node"
	fragment = strings.Trim(fragment, "/")
	path := db.url.Path
	path = strings.TrimRight(path, "/")
	path = strings.Join([]string{path, fragment}, "/")
	u := db.url
	u.Path = path
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return
	}
	if props != nil {
		var b []byte
		req.Header.Add("Content-Type", "application/json")
		b, err = json.Marshal(props)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest("POST", u.String(), buf)
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("POST", u.String(), nil)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := db.client.Do(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &n)
	return
}
*/
