// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
	"net/url"
	"strconv"
)

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	rc              *restclient.Client
	Url             string      `json:"-"` // Root URL for REST API
	HrefNode        string      `json:"node"`
	HrefRefNode     string      `json:"reference_node"`
	HrefNodeIndex   string      `json:"node_index"`
	HrefRelIndex    string      `json:"relationship_index"`
	HrefExtInfo     string      `json:"extensions_info"`
	HrefRelTypes    string      `json:"relationship_types"`
	HrefBatch       string      `json:"batch"`
	HrefCypher      string      `json:"cypher"`
	HrefTransaction string      `json:"transaction"`
	Version         string      `json:"neo4j_version"`
	Extensions      interface{} `json:"extensions"`
}

// Connect establishes a connection to the Neo4j server.
func Connect(uri string) (*Database, error) {
	var e NeoError
	db := &Database{
		rc: restclient.New(),
	}
	_, err := url.Parse(uri) // Sanity check
	if err != nil {
		return nil, err
	}
	db.Url = uri
	req := restclient.RequestResponse{
		Url:    db.Url,
		Method: "GET",
		Result: &db,
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
	case status != 200 || db.Version == "":
		log.Println("Status " + strconv.Itoa(status) + " trying to cconnect to " + uri)
		logPretty(req)
		return db, e
	}
	return db, nil
}

// A Props is a set of key/value properties.
type Props map[string]interface{}
