// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strings"
)

/*

NOTE:  The API for working with Nodes and Relationships is basically identical. 
So we base both types on nrBase.  But is this a good idea?  Is the similarity
between Node and Relationship by design, or just coincidental?

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
