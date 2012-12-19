// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
)

type IndexConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
}

// An indexInfo is returned from the Neo4j server on operations involving an Index.
type indexInfo struct {
	neoError
	Template string `json:"template"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
}

type Index struct {
	db   *Database
	info *indexInfo
}

// CreateIndex creates a new Index, with the name supplied, in the db.
func (db *Database) CreateIndex(name string) (*Index, error) {
	conf := IndexConfig{
		Name: name,
	}
	return db.CreateIndexFromConf(conf)
}

// CreateIndexFromConf creates a new Index based on an IndexConfig object
func (db *Database) CreateIndexFromConf(conf IndexConfig) (*Index, error) {
	var info indexInfo
	i := Index{
		db:   db,
		info: &info,
	}
	c := restclient.RestRequest{
		Url:    db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &conf,
		Result: &info,
		Error:  new(neoError),
	}
	status, err := db.rc.Do(&c)
	if err != nil {
		return &i, err
	}
	if status != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", status)
		log.Printf("    Result:", info)
		return &i, BadResponse
	}
	return &i, nil
}
