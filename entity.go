// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strings"
)

// An entity is an object - either a Node or a Relationship - in a Neo4j graph
// database.  An entity may optinally be assigned an arbitrary set of key:value
// properties.
type entity interface {
	// HrefProperty()   string
	// HrefProperties() string
	// HrefSelf()       string
	SetProperty(key string, value string) error
	GetProperty(key string) (string, error)
	DeleteProperty(key string) error
	Delete() error
	Properties() (Properties, error)
	SetProperties(p Properties) error
	DeleteProperties() error
}

type baseEntity struct {
	HrefProperty   string
	HrefProperties string
	HrefSelf       string
	db             *Database
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (e *baseEntity) do(rr *restclient.RequestResponse) (status int, err error) {
	return e.db.rc.Do(rr)
}

// Properties is a bag of key/value pairs that describe an baseEntity.
type Properties map[string]string

// EmptyProps is an empty Properties map.
var EmptyProps = Properties{}

// SetProperty sets the single property key to value.
func (e *baseEntity) SetProperty(key string, value string) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	rr := restclient.RequestResponse{
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
func (e *baseEntity) Property(key string) (string, error) {
	var val string
	if e.HrefProperties == "" {
		return val, FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &val,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logPretty(ne)
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
func (e *baseEntity) DeleteProperty(key string) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	parts := []string{e.HrefProperties, key}
	uri := strings.Join(parts, "/")
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logPretty(ne)
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logPretty(ne)
	return BadResponse
}

// Delete removes the object from the DB.
func (e *baseEntity) Delete() error {
	if e.HrefSelf == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    e.HrefSelf,
		Method: restclient.DELETE,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	switch {
	case err != nil:
		logPretty(ne)
		return err
	case status == 204:
		// Successful deletion!
		return nil
	case status == 409:
		return CannotDelete
	}
	logPretty(ne)
	return BadResponse
}

// Properties fetches all properties
func (e *baseEntity) Properties() (Properties, error) {
	props := make(map[string]string)
	if e.HrefProperties == "" {
		return props, FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: restclient.GET,
		Result: &props,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logPretty(ne)
		return props, err
	}
	// Status code 204 indicates no properties on this node
	if status == 204 {
		props = map[string]string{}
	}
	return props, nil
}

// SetProperties updates all properties, overwriting any existing properties.
func (e *baseEntity) SetProperties(p Properties) error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: restclient.PUT,
		Data:   &p,
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logPretty(ne)
		return err
	}
	if status == 204 {
		return nil // Success!
	}
	logPretty(ne)
	return BadResponse
}

// DeleteProperties deletes all properties.
func (e *baseEntity) DeleteProperties() error {
	if e.HrefProperties == "" {
		return FeatureUnavailable
	}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    e.HrefProperties,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := e.do(&rr)
	if err != nil {
		logPretty(ne)
		return err
	}
	switch status {
	case 204:
		return nil // Success!
	case 404:
		return NotFound
	}
	logPretty(ne)
	return BadResponse
}
