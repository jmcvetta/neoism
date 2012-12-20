// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strconv"
	"strings"
)

type NodeManager struct {
	db      *Database
	Indexes *NodeIndexManager
}

// CreateNode creates a Node in the database.
func (m *NodeManager) Create(p Properties) (*Node, error) {
	var info nrInfo
	n := Node{nrBase: nrBase{
		db:   m.db,
		info: &info,
	}}
	c := restclient.RestRequest{
		Url:    m.db.info.Node,
		Method: restclient.POST,
		Data:   &p,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := m.db.rc.Do(&c)
	if err != nil || status != 201 {
		return &n, err
	}
	if info.Self == "" {
		return &n, BadResponse
	}
	return &n, nil
}

// GetNode fetches a Node from the database
func (m *NodeManager) Get(id int) (*Node, error) {
	uri := join(m.db.info.Node, strconv.Itoa(id))
	return m.getNodeByUri(uri)
}

// GetNode fetches a Node from the database based on its uri
func (m *NodeManager) getNodeByUri(uri string) (*Node, error) {
	var info nrInfo
	n := Node{nrBase{
		db:   m.db,
		info: &info,
	}}
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := m.db.rc.Do(&c)
	switch {
	case status == 404:
		return &n, NotFound
	case status != 200 || info.Self == "":
		return &n, BadResponse
	}
	if err != nil {
		return &n, err
	}
	// n.Info = &info
	return &n, nil
}

// A node in a Neo4j database
type Node struct {
	nrBase
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.db.info.Node)
	s := n.info.Self[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return id
}

// getRelationships makes an api call to the supplied uri and returns a map 
// keying relationship IDs to Relationship objects.
func (n *Node) getRelationships(uri string, types ...string) (map[int]Relationship, error) {
	m := map[int]Relationship{}
	if uri == "" {
		return m, FeatureUnavailable
	}
	if types != nil {
		fragment := strings.Join(types, "&")
		parts := []string{uri, fragment}
		uri = strings.Join(parts, "/")
	}
	s := []nrInfo{}
	c := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &s,
		Error:  new(neoError),
	}
	status, err := n.db.rc.Do(&c)
	if err != nil {
		return m, err
	}
	for _, info := range s {
		rel := Relationship{nrBase{
			db:   n.db,
			info: &info,
		}}
		m[rel.Id()] = rel
	}
	if status == 200 {
		return m, nil // Success!
	}
	return m, BadResponse
}

// Relationships gets all Relationships for this Node, optionally filtered by 
// type, returning them as a map keyed on Relationship ID.
func (n *Node) Relationships(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.AllRels, types...)
}

// Incoming gets all incoming Relationships for this Node.
func (n *Node) Incoming(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.IncomingRels, types...)
}

// Outgoing gets all outgoing Relationships for this Node.
func (n *Node) Outgoing(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.info.OutgoingRels, types...)
}

// Relate creates a relationship of relType, with specified properties, 
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Properties) (*Relationship, error) {
	var info nrInfo
	rel := Relationship{nrBase{
		db:   n.db,
		info: &info,
	}}
	srcUri := join(n.info.Self, "relationships")
	destUri := join(n.db.info.Node, strconv.Itoa(destId))
	content := map[string]interface{}{
		"to":   destUri,
		"type": relType,
	}
	if p != nil {
		content["data"] = &p
	}
	c := restclient.RestRequest{
		Url:    srcUri,
		Method: restclient.POST,
		Data:   content,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := n.db.rc.Do(&c)
	if err != nil {
		return &rel, err
	}
	if status != 201 {
		return &rel, BadResponse
	}
	return &rel, nil
}
