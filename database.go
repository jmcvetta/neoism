// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
	"net/http"
	"net/url"
)

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	url           *url.URL // Root URL for REST API
	client        *http.Client
	rc            *restclient.Client
	info          *serviceRoot
	Nodes         *NodeManager
	Relationships *RelationshipManager
}

// A serviceRoot describes services available on the Neo4j server
type serviceRoot struct {
	Extensions interface{} `json:"extensions"`
	Node       string      `json:"node"`
	RefNode    string      `json:"reference_node"`
	NodeIndex  string      `json:"node_index"`
	RelIndex   string      `json:"relationship_index"`
	ExtInfo    string      `json:"extensions_info"`
	RelTypes   string      `json:"relationship_types"`
	Batch      string      `json:"batch"`
	Cypher     string      `json:"cypher"`
	Version    string      `json:"neo4j_version"`
}

func Connect(uri string) (*Database, error) {
	var sr serviceRoot
	var e neoError
	db := &Database{
		client: new(http.Client),
		rc:     restclient.New(),
		info:   &sr,
	}
	u, err := url.Parse(uri)
	if err != nil {
		return db, err
	}
	db.url = u
	db.Nodes = &NodeManager{
		db: db,
		Indexes: &NodeIndexManager{
			db: db,
		},
	}
	db.Relationships = &RelationshipManager{
		db: db,
	}
	r := restclient.RestRequest{
		Url:    u.String(),
		Method: restclient.GET,
		Result: &sr,
		Error:  &e,
	}
	status, err := db.rc.Do(&r)
	if err != nil {
		log.Println(e.Message)
		log.Println(e.Exception)
		log.Println(e.Stacktrace)
		return db, err
	}
	switch {
	case status == 200 && db.info.Version != "":
		return db, nil // Success!
	case status == 404:
		return db, InvalidDatabase
	}
	log.Println(e.Message)
	log.Println(e.Exception)
	log.Println(e.Stacktrace)
	return db, BadResponse
}
