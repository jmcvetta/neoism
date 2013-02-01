// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"github.com/jmcvetta/restclient"
	"net/http"
	"net/url"
)

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	url           *url.URL // Root URL for REST API
	client        *http.Client
	rc            *restclient.Client
	Nodes         *NodeManager
	Relationships *RelationshipManager
	//
	HrefNode      string
	HrefRefNode   string
	HrefNodeIndex string
	HrefRelIndex  string
	HrefExtInfo   string
	HrefRelTypes  string
	HrefBatch     string
	HrefCypher    string
	Version       string
}

// A serviceRoot describes services available on the Neo4j server
type serviceRoot struct {
	Extensions    interface{} `json:"extensions"`
	HrefNode      string      `json:"node"`
	HrefRefNode   string      `json:"reference_node"`
	HrefNodeIndex string      `json:"node_index"`
	HrefRelIndex  string      `json:"relationship_index"`
	HrefExtInfo   string      `json:"extensions_info"`
	HrefRelTypes  string      `json:"relationship_types"`
	HrefBatch     string      `json:"batch"`
	HrefCypher    string      `json:"cypher"`
	Version       string      `json:"neo4j_version"`
}

func Connect(uri string) (*Database, error) {
	var sr serviceRoot
	var e neoError
	db := &Database{
		client: new(http.Client),
		rc:     restclient.New(),
	}
	u, err := url.Parse(uri)
	if err != nil {
		return db, err
	}
	db.url = u
	db.Nodes = &NodeManager{
		db:      db,
		Indexes: &NodeIndexManager{},
	}
	db.Nodes.Indexes.db = db
	db.Relationships = &RelationshipManager{
		db:      db,
		Indexes: &RelationshipIndexManager{},
	}
	req := restclient.RestRequest{
		Url:    u.String(),
		Method: restclient.GET,
		Result: &sr,
		Error:  &e,
	}
	status, err := db.rc.Do(&req)
	if err != nil {
		logPretty(req)
		return db, err
	}
	switch {
	case status == 404:
		return db, InvalidDatabase
	case status != 200 || sr.Version == "":
		logPretty(req)
		return db, BadResponse
	}
	// Populate Database struct
	db.HrefNode = sr.HrefNode
	db.HrefRefNode = sr.HrefRefNode
	db.HrefNodeIndex = sr.HrefNodeIndex
	db.HrefRelIndex = sr.HrefRelIndex
	db.HrefExtInfo = sr.HrefExtInfo
	db.HrefRelTypes = sr.HrefRelTypes
	db.HrefBatch = sr.HrefBatch
	db.HrefCypher = sr.HrefCypher
	db.Version = sr.Version
	// Set HrefIndex so the generic indexManager knows what URL to use when
	// creating a NodeIndex.
	db.Nodes.Indexes.HrefIndex = sr.HrefNodeIndex
	db.Relationships.Indexes.HrefIndex = sr.HrefRelIndex
	// Success!
	return db, nil
}
