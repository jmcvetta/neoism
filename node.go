// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strconv"
	"strings"
)

// CreateNode creates a Node in the database.
func (db *Database) CreateNode(p Props) (*Node, error) {
	n := Node{}
	n.db = db
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:            db.HrefNode,
		Method:         "POST",
		Data:           &p,
		Result:         &n,
		Error:          ne,
		ExpectedStatus: 201,
	}
	status, err := db.rc.Do(&rr)
	if err != nil {
		logPretty(status)
		logPretty(ne)
		return &n, err
	}
	return &n, nil
}

// Node fetches a Node from the database
func (db *Database) Node(id int) (*Node, error) {
	uri := join(db.HrefNode, strconv.Itoa(id))
	return db.getNodeByUri(uri)
}

// getNodeByUri fetches a Node from the database based on its URI.
func (db *Database) getNodeByUri(uri string) (*Node, error) {
	ne := new(neoError)
	n := Node{}
	n.db = db
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &n,
		Error:  ne,
	}
	status, err := db.rc.Do(&rr)
	switch {
	case status == 404:
		return &n, NotFound
	case status != 200 || n.HrefSelf == "":
		logPretty(ne)
		return &n, BadResponse
	}
	if err != nil {
		logPretty(ne)
		return &n, err
	}
	return &n, nil
}

// A Node is a node, with optional properties, in a graph.
type Node struct {
	entity
	// HrefSelf              string      `json:"self"`
	// HrefProperty          string      `json:"property"`
	// HrefProperties        string      `json:"properties"`
	HrefOutgoingRels      string                 `json:"outgoing_relationships"`
	HrefTraverse          string                 `json:"traverse"`
	HrefAllTypedRels      string                 `json:"all_typed_relationships"`
	HrefOutgoing          string                 `json:"outgoing_typed_relationships"`
	HrefIncomingRels      string                 `json:"incoming_relationships"`
	HrefCreateRel         string                 `json:"create_relationship"`
	HrefPagedTraverse     string                 `json:"paged_traverse"`
	HrefAllRels           string                 `json:"all_relationships"`
	HrefIncomingTypedRels string                 `json:"incoming_typed_relationships"`
	Data                  map[string]interface{} `json:"data"`
	Extensions            map[string]interface{} `json:"extensions"`
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.db.HrefNode)
	s := n.HrefSelf[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return id
}

// getRels makes an api call to the supplied uri and returns a map
// keying relationship IDs to Rel objects.
func (n *Node) getRels(uri string, types ...string) (Rels, error) {
	if types != nil {
		fragment := strings.Join(types, "&")
		parts := []string{uri, fragment}
		uri = strings.Join(parts, "/")
	}
	rels := Rels{}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &rels,
		Error:  &ne,
	}
	status, err := n.db.rc.Do(&rr)
	if err != nil {
		return rels, err
	}
	if status == 200 {
		return rels, nil // Success!
	}
	return rels, BadResponse
}

// Rels gets all Rels for this Node, optionally filtered by
// type, returning them as a map keyed on Rel ID.
func (n *Node) Relationships(types ...string) (Rels, error) {
	return n.getRels(n.HrefAllRels, types...)
}

// Incoming gets all incoming Rels for this Node.
func (n *Node) Incoming(types ...string) (Rels, error) {
	return n.getRels(n.HrefIncomingRels, types...)
}

// Outgoing gets all outgoing Rels for this Node.
func (n *Node) Outgoing(types ...string) (Rels, error) {
	return n.getRels(n.HrefOutgoingRels, types...)
}

// Relate creates a relationship of relType, with specified properties,
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Props) (*Relationship, error) {
	rel := Relationship{}
	rel.db = n.db
	ne := new(neoError)
	srcUri := join(n.HrefSelf, "relationships")
	destUri := join(n.db.HrefNode, strconv.Itoa(destId))
	content := map[string]interface{}{
		"to":   destUri,
		"type": relType,
	}
	if p != nil {
		content["data"] = &p
	}
	c := restclient.RequestResponse{
		Url:    srcUri,
		Method: "POST",
		Data:   content,
		Result: &rel,
		Error:  &ne,
	}
	status, err := n.db.rc.Do(&c)
	if err != nil {
		logPretty(ne)
		return &rel, err
	}
	if status != 201 {
		logPretty(ne)
		return &rel, BadResponse
	}
	return &rel, nil
}
