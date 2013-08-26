// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"github.com/jmcvetta/restclient"
	"sort"
	"strconv"
	"strings"
)

// Relationship fetches a Relationship from by id.
func (db *Database) Relationship(id int) (*Relationship, error) {
	rel := Relationship{}
	rel.Db = db
	uri := join(db.Url, "relationship", strconv.Itoa(id))
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &rel,
		Error:  &ne,
	}
	status, err := db.Rc.Do(&rr)
	if err != nil {
		return &rel, err
	}
	switch status {
	default:
		logPretty(ne)
		err = ne
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
	ne := NeoError{}
	c := restclient.RequestResponse{
		Url:    db.HrefRelTypes,
		Method: "GET",
		Result: &reltypes,
		Error:  &ne,
	}
	status, err := db.Rc.Do(&c)
	if err != nil {
		return reltypes, err
	}
	if status == 200 {
		sort.Sort(sort.StringSlice(reltypes))
		return reltypes, nil // Success!
	}
	logPretty(ne)
	return reltypes, ne
}

// A Relationship is a directional connection between two Nodes, with an
// optional set of arbitrary properties.
type Relationship struct {
	entity
	Type       string      `json:"type"`
	HrefStart  string      `json:"start"`
	HrefEnd    string      `json:"end"`
	Data       interface{} `json:"data"`
	Extensions interface{} `json:"extensions"`
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
	return r.Db.getNodeByUri(r.HrefStart)
}

// End gets the ending Node of this Relationship.
func (r *Relationship) End() (*Node, error) {
	return r.Db.getNodeByUri(r.HrefEnd)
}

// A Rels is a collection of relationships.
type Rels []*Relationship

// Map presents a set of relationships as a map associating relationship IDs to
// relationship objects.
func (r *Rels) Map() map[int]*Relationship {
	m := make(map[int]*Relationship, len(*r))
	for _, rel := range *r {
		m[rel.Id()] = rel
	}
	return m
}
