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
func (db *Database) CreateNode(p Properties) (*Node, error) {
	n := Node{}
	n.Db = db
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
	n.Db = db
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
	case status != 200 || n.hrefSelf == "":
		logPretty(ne)
		return &n, BadResponse
	}
	if err != nil {
		logPretty(ne)
		return &n, err
	}
	return &n, nil
}

// A node in a Neo4j database
type Node struct {
	Db             *Database
	hrefProperty   string `json:"property"`
	hrefProperties string `json:"properties"`
	hrefSelf       string `json:"self"`
	// hrefData       interface{} `json:"data"`
	// hrefExtensions interface{} `json:"extensions"`
	//
	hrefOutgoingRels      string `json:"outgoing_relationships"`
	hrefTraverse          string `json:"traverse"`
	hrefAllTypedRels      string `json:"all_typed_relationships"`
	hrefOutgoing          string `json:"outgoing_typed_relationships"`
	hrefIncomingRels      string `json:"incoming_relationships"`
	hrefCreateRel         string `json:"create_relationship"`
	hrefPagedTraverse     string `json:"paged_traverse"`
	hrefAllRels           string `json:"all_relationships"`
	hrefIncomingTypedRels string `json:"incoming_typed_relationships"`
}

// Do executes a REST request.
func (n *Node) Do(rr *restclient.RequestResponse) (status int, err error) {
	return n.Db.rc.Do(rr)
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.Db.HrefNode)
	s := n.hrefSelf[l:]
	s = strings.Trim(s, "/")
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return id
}

func (n *Node) HrefProperty() string {
	return n.hrefProperty
}

func (n *Node) HrefSelf() string {
	return n.hrefSelf
}

func (n *Node) HrefProperties() string {
	return n.hrefProperties
}

func (n *Node) Delete() error {
	return delete(n)
}

func (n *Node) Properties() (Properties, error) {
	return properties(n)
}

func (n *Node) Property(key string) (string, error) {
	return property(n, key)
}

func (n *Node) SetProperty(key, value string) error {
	return setProperty(n, key, value)
}

func (n *Node) SetProperties(p Properties) error {
	return setProperties(n, p)
}

func (n *Node) DeleteProperties() error {
	return deleteProperties(n)
}

func (n *Node) DeleteProperty(key string) error {
	return deleteProperty(n, key)
}

// getRelationships makes an api call to the supplied uri and returns a map
// keying relationship IDs to Relationship objects.
func (n *Node) getRelationships(uri string, types ...string) (map[int]Relationship, error) {
	m := map[int]Relationship{}
	if types != nil {
		fragment := strings.Join(types, "&")
		parts := []string{uri, fragment}
		uri = strings.Join(parts, "/")
	}
	rels := []Relationship{}
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:    uri,
		Method: "GET",
		Result: &rels,
		Error:  &ne,
	}
	status, err := n.Do(&rr)
	if err != nil {
		return m, err
	}
	for _, r := range rels {
		r.Db = n.Db
		m[r.Id()] = r
	}
	if status == 200 {
		return m, nil // Success!
	}
	return m, BadResponse
}

// Relationships gets all Relationships for this Node, optionally filtered by
// type, returning them as a map keyed on Relationship ID.
func (n *Node) Relationships(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.hrefAllRels, types...)
}

// Incoming gets all incoming Relationships for this Node.
func (n *Node) Incoming(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.hrefIncomingRels, types...)
}

// Outgoing gets all outgoing Relationships for this Node.
func (n *Node) Outgoing(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.hrefOutgoingRels, types...)
}

// Relate creates a relationship of relType, with specified properties,
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Properties) (*Relationship, error) {
	rel := Relationship{}
	rel.Db = n.Db
	ne := new(neoError)
	srcUri := join(n.hrefSelf, "relationships")
	destUri := join(n.Db.HrefNode, strconv.Itoa(destId))
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
	status, err := n.Db.rc.Do(&c)
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
