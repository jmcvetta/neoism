// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strings"
)

// Properties is a bag of key/value pairs that can describe Nodes
// and Relationships.
// type Properties map[string]string

// nrBase is the base type for Nodes and Relationships.
type Entity struct {
	db             *Database
	HrefProperty   string
	HrefProperties string
	HrefSelf       string
}

func (e *Entity) do(rr *restclient.RestRequest) (status int, err error) {
	return e.db.rc.Do(rr)
}

// A nrInfo is returned from the Neo4j server on successful operations 
// involving a Node or a Relationship.
type entityInfo struct {
	neoError
	//
	// Always filled on success
	//
	HrefProperty   string      `json:"property"`
	HrefProperties string      `json:"properties"`
	HrefSelf       string      `json:"self"`
	HrefData       interface{} `json:"data"`
	HrefExtensions interface{} `json:"extensions"`
	//
	// Filled only for Node operations
	//
	HrefOutgoingRels      string `json:"outgoing_relationships"`
	HrefTraverse          string `json:"traverse"`
	HrefAllTypedRels      string `json:"all_typed_relationships"`
	HrefOutgoing          string `json:"outgoing_typed_relationships"`
	HrefIncomingRels      string `json:"incoming_relationships"`
	HrefCreateRel         string `json:"create_relationship"`
	HrefPagedTraverse     string `json:"paged_traverse"`
	HrefAllRels           string `json:"all_relationships"`
	HrefIncomingTypedRels string `json:"incoming_typed_relationships"`
	//
	// Filled only for Relationship operations
	//
	HrefStart string `json:"start"`
	HrefType  string `json:"type"`
	HrefEnd   string `json:"end"`
}

// SetProperty sets the single property key to value.
func (e *Entity) SetProperty(key string, value string) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	rr := restclient.RestRequest{
		Url:    uri,
		Method: restclient.PUT,
		Data:   &value,
		Error:  new(neoError),
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil // Success!
	}
	return BadResponse
}

// GetProperty fetches the value of property key.
func (e *Entity) GetProperty(key string) (string, error) {
	var val string
	if e.HrefProperties == "" {
		return val, FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    uri,
		Method: "GET",
		Result: &val,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logError(ne)
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
func (e *Entity) DeleteProperty(key string) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logError(ne)
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logError(ne)
	return BadResponse
}

// Delete removes the object from the DB.
func (e *Entity) Delete() error {
	if e.HrefSelf == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    e.HrefSelf,
		Method: restclient.DELETE,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	switch {
	case err != nil:
		logError(ne)
		return err
	case status == 204:
		// Successful deletion!
		return nil
	case status == 409:
		logError(ne)
		return CannotDelete
	}
	logError(ne)
	return BadResponse
}

// Properties fetches all properties
func (e *Entity) Properties() (Properties, error) {
	props := make(map[string]string)
	if e.HrefProperties == "" {
		return props, FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    e.HrefProperties,
		Method: restclient.GET,
		Result: &props,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logError(ne)
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if status == 204 {
		props = map[string]string{}
	}
	return props, nil
}

// SetProperties updates all properties, overwriting any existing properties.
func (e *Entity) SetProperties(p Properties) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    e.HrefProperties,
		Method: restclient.PUT,
		Data:   &p,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logError(ne)
		return err
	}
	if status == 204 {
		return nil // Success!
	}
	logError(ne)
	return BadResponse
}

// DeleteProperties deletes all properties.
func (e *Entity) DeleteProperties() error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    e.HrefProperties,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logError(ne)
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logError(ne)
	return BadResponse
}
