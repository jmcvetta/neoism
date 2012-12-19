// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"errors"
	"github.com/jmcvetta/restclient"
	"log"
	"net/http"
	"net/url"
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
	url           *url.URL // Root URL for REST API
	client        *http.Client
	rc            *restclient.Client
	info          *serviceRootInfo
	Nodes         *NodeManager
	Relationships *RelationshipManager
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
	db.Nodes = &NodeManager{
		db: db,
	}
	db.Relationships = &RelationshipManager{
		db: db,
	}
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

/*
func (db *Database) Nodes() *NodeManager {
	return &NodeManager{
		db: db,
	}
}

func (db *Database) Relationships() *RelationshipManager {
	return &RelationshipManager{
		db: db,
	}
}
*/

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
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.PUT,
		Data:   &value,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return err
	}
	if status == 204 {
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: "GET",
		Result: &val,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return val, err
	}
	switch status {
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return err
	}
	switch status {
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	switch {
	case err != nil:
		return err
	case status == 204:
		// Successful deletion!
		return nil
	case status == 409:
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &props,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if status == 204 {
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.PUT,
		Data:   &p,
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return err
	}
	if status == 204 {
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
	c := restclient.RestRequest{
		Url:    uri,
		Method: "DELETE",
		Error:  new(neoError),
	}
	status, err := e.db.rc.Do(&c)
	if err != nil {
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	return BadResponse
}
