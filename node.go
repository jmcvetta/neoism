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

// do is a convenience wrapper around the embedded restclient's Do() method.
func (nm *NodeManager) do(rr *restclient.RestRequest) (status int, err error) {
	return nm.db.rc.Do(rr)
}

// CreateNode creates a Node in the database.
func (m *NodeManager) Create(p Properties) (*Node, error) {
	n := Node{}
	n.db = m.db
	res := new(nodeResponse)
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    m.db.info.Node,
		Method: restclient.POST,
		Data:   &p,
		Result: res,
		Error:  ne,
	}
	status, err := m.do(&rr)
	if err != nil || status != 201 {
		logPretty(ne)
		return &n, err
	}
	if res.HrefSelf == "" {
		logPretty(ne)
		return &n, BadResponse
	}
	n.populate(res)
	return &n, nil
}

// GetNode fetches a Node from the database
func (m *NodeManager) Get(id int) (*Node, error) {
	uri := join(m.db.info.Node, strconv.Itoa(id))
	return m.getNodeByUri(uri)
}

// getNodeByUri fetches a Node from the database based on its URI.
func (m *NodeManager) getNodeByUri(uri string) (*Node, error) {
	res := new(nodeResponse)
	ne := new(neoError)
	n := Node{}
	n.db = m.db
	rr := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: res,
		Error:  ne,
	}
	status, err := m.do(&rr)
	switch {
	case status == 404:
		return &n, NotFound
	case status != 200 || res.HrefSelf == "":
		logPretty(ne)
		return &n, BadResponse
	}
	if err != nil {
		logPretty(ne)
		return &n, err
	}
	n.populate(res)
	return &n, nil
}

type nodeResponse struct {
	HrefProperty   string `json:"property"`
	HrefProperties string `json:"properties"`
	HrefSelf       string `json:"self"`
	// HrefData       interface{} `json:"data"`
	// HrefExtensions interface{} `json:"extensions"`
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
}

// populate uses the values from a nodeResponse object to populate the fields on
// this Node.
func (n *Node) populate(r *nodeResponse) {
	n.HrefProperty = r.HrefProperty
	n.HrefProperties = r.HrefProperties
	n.HrefSelf = r.HrefSelf
	// n.HrefData = r.HrefData
	// n.HrefExtensions = r.HrefExtensions
	n.HrefOutgoingRels = r.HrefOutgoingRels
	n.HrefTraverse = r.HrefTraverse
	n.HrefAllTypedRels = r.HrefAllTypedRels
	n.HrefOutgoing = r.HrefOutgoing
	n.HrefIncomingRels = r.HrefIncomingRels
	n.HrefCreateRel = r.HrefCreateRel
	n.HrefPagedTraverse = r.HrefPagedTraverse
	n.HrefAllRels = r.HrefAllRels
	n.HrefIncomingTypedRels = r.HrefIncomingTypedRels
}

// A node in a Neo4j database
type Node struct {
	baseEntity
	HrefOutgoingRels      string
	HrefTraverse          string
	HrefAllTypedRels      string
	HrefOutgoing          string
	HrefIncomingRels      string
	HrefCreateRel         string
	HrefPagedTraverse     string
	HrefAllRels           string
	HrefIncomingTypedRels string
}

// Id gets the ID number of this Node.
func (n *Node) Id() int {
	l := len(n.db.info.Node)
	s := n.HrefSelf[l:]
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
	resArray := []relationshipResponse{}
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &resArray,
		Error:  &ne,
	}
	status, err := n.do(&rr)
	if err != nil {
		return m, err
	}
	for _, res := range resArray {
		rel := Relationship{}
		rel.db = n.db
		rel.populate(&res)
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
	return n.getRelationships(n.HrefAllRels, types...)
}

// Incoming gets all incoming Relationships for this Node.
func (n *Node) Incoming(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.HrefIncomingRels, types...)
}

// Outgoing gets all outgoing Relationships for this Node.
func (n *Node) Outgoing(types ...string) (map[int]Relationship, error) {
	return n.getRelationships(n.HrefOutgoingRels, types...)
}

// Relate creates a relationship of relType, with specified properties, 
// from this Node to the node identified by destId.
func (n *Node) Relate(relType string, destId int, p Properties) (*Relationship, error) {
	rel := Relationship{}
	rel.db = n.db
	res := new(relationshipResponse)
	ne := new(neoError)
	srcUri := join(n.HrefSelf, "relationships")
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
		Result: &res,
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
	rel.populate(res)
	return &rel, nil
}
