// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"sort"
	"strconv"
	"strings"
)

type RelationshipManager struct {
	db *Database
}

// GetRelationship fetches a Relationship from the DB by id.
func (m *RelationshipManager) Get(id int) (*Relationship, error) {
	var info nrInfo
	rel := Relationship{nrBase{
		db:   m.db,
		info: &info,
	}}
	uri := join(m.db.url.String(), "relationship", strconv.Itoa(id))
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := m.db.rc.Do(&c)
	switch status {
	default:
		err = BadResponse
	case 200:
		err = nil // Success!
	case 404:
		err = NotFound
	}
	return &rel, err
}

// Types lists all existing relationship types
func (m *RelationshipManager) Types() ([]string, error) {
	reltypes := []string{}
	if m.db.info.RelTypes == "" {
		return reltypes, FeatureUnavailable
	}
	c := restclient.RestRequest{
		Url:    m.db.info.RelTypes,
		Method: restclient.GET,
		Result: &reltypes,
		Error:  new(neoError),
	}
	status, err := m.db.rc.Do(&c)
	if err != nil {
		return reltypes, err
	}
	if status == 200 {
		return reltypes, nil // Success!
	}
	sort.Sort(sort.StringSlice(reltypes))
	return reltypes, BadResponse
}

// A relationship in a Neo4j database
type Relationship struct {
	nrBase
}

// Id gets the ID number of this Relationship
func (r *Relationship) Id() int {
	parts := strings.Split(r.info.Self, "/")
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
	return r.db.Nodes.getNodeByUri(r.info.Start)
}

// End gets the ending Node of this Relationship.
func (r *Relationship) End() (*Node, error) {
	return r.db.Nodes.getNodeByUri(r.info.End)
}

// Type gets the type of this relationship
func (r *Relationship) Type() string {
	return r.info.Type
}
