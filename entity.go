// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strings"
)

// An entity is an object - either a Node or a Relationship - in a Neo4j graph
// database.  An entity may optinally be assigned an arbitrary set of key:value
// properties.
type entity struct {
	db             *Database
	HrefSelf       string `json:"self"`
	HrefProperty   string `json:"property"`
	HrefProperties string `json:"properties"`
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (e *entity) do(rr *restclient.RequestResponse) (status int, err error) {
	return e.db.rc.Do(rr)
}

// SetProperty sets the single property key to value.
func (e *entity) SetProperty(key string, value string) error {
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "PUT",
		Data:   &value,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	if status != 204 {
		return ne
	}
	return nil // Success!
}

// GetProperty fetches the value of property key.
func (e *entity) Property(key string) (string, error) {
	var val string
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &val,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return val, err
	}
	switch status {
	case 200:
	case 404:
		return val, NotFound
	default:
		return val, ne
	}
	return val, nil // Success!
}

// DeleteProperty deletes property key
func (e *entity) DeleteProperty(key string) error {
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logPretty(ne)
	return ne
}

// Delete removes the object from the DB.
func (e *entity) Delete() error {
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    e.HrefSelf,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	switch status {
	case 204:
	case 404:
		return NotFound
	case 409:
		return CannotDelete
	default:
		logPretty(status)
		logPretty(ne)
		return ne
	}
	return nil
}

// Properties fetches all properties
func (e *entity) Properties() (Props, error) {
	props := Props{}
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: "GET",
		Result: &props,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if status == 204 {
		props = Props{}
	}
	return props, nil
}

// SetProperties updates all properties, overwriting any existing properties.
func (e *entity) SetProperties(p Props) error {
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: "PUT",
		Data:   &p,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil // Success!
	}
	logPretty(ne)
	return ne
}

// DeleteProperties deletes all properties.
func (e *entity) DeleteProperties() error {
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logPretty(ne)
	return ne
}
