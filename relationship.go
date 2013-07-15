// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"sort"
	"strconv"
	"strings"
)

// Relationship fetches a Relationship from by id.
func (db *Database) Relationship(id int) (*Relationship, error) {
	rel := Relationship{}
	rel.db = db
	uri := join(db.url.String(), "relationship", strconv.Itoa(id))
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &rel,
		Error:  &ne,
	}
	status, err := db.rc.Do(&rr)
	if err != nil {
		logPretty(ne)
		return &rel, err
	}
	switch status {
	default:
		logPretty(ne)
		err = BadResponse
	case 200:
		err = nil // Success!
	case 404:
		err = NotFound
	}
	return &rel, err
}

// Types lists all existing relationship types
func (db *Database) RelTypes() ([]string, error) {
	reltypes := []string{}
	ne := new(neoError)
	c := restclient.RequestResponse{
		Url:    db.HrefRelTypes,
		Method: "GET",
		Result: &reltypes,
		Error:  &ne,
	}
	status, err := db.rc.Do(&c)
	if err != nil {
		logPretty(ne)
		return reltypes, err
	}
	if status == 200 {
		sort.Sort(sort.StringSlice(reltypes))
		return reltypes, nil // Success!
	}
	logPretty(ne)
	return reltypes, BadResponse
}

type Rels []*Relationship

// A Relationship is a directional connection between two Nodes, with an
// optional set of arbitrary properties.
type Relationship struct {
	entity
	HrefData       interface{} `json:"data"`
	HrefExtensions interface{} `json:"extensions"`
	HrefStart      string      `json:"start"`
	HrefType       string      `json:"type"`
	HrefEnd        string      `json:"end"`
}

func (r *Relationship) hrefSelf() string {
	return r.HrefSelf
}

// Id gets the ID number of this Relationship
func (r *Relationship) Id() int {
	parts := strings.Split(r.HrefSelf, "/")
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
	return r.db.getNodeByUri(r.HrefStart)
}

// End gets the ending Node of this Relationship.
func (r *Relationship) End() (*Node, error) {
	return r.db.getNodeByUri(r.HrefEnd)
}

// Type gets the type of this relationship
func (r *Relationship) Type() string {
	return r.HrefType
}
