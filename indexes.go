// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
)

// An index can contain either nodes or relationships.
type Index struct {
	Template string `json:"template"`
	Provider string `json:"provider"`
	Type     string `json:"type"`
}

type IndexConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
}

type NodeIndexManager struct {
	db *Database
}

// CreateIndex creates a new Index, with the name supplied, in the db.
func (nim *NodeIndexManager) Create(name string) (*Index, error) {
	type req struct {
		Name string `json:"name"`
	}
	data := req{Name: name}
	var i Index
	var e neoError
	c := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &i,
		Error:  &e,
	}
	status, err := nim.db.rc.Do(&c)
	if err != nil {
		return &i, err
	}
	if status != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", status)
		log.Printf("    Error:", e)
		return &i, BadResponse
	}
	return &i, nil
}

func (nim *NodeIndexManager) CreateWithConf(name, indexType, provider string) (*Index, error) {
	type conf struct {
		Type     string `json:"type"`
		Provider string `json:"provider"`
	}
	type req struct {
		Name   string `json:"name"`
		Config conf   `json:"config"`
	}
	data := req{
		Name: name,
		Config: conf{
			Type:     indexType,
			Provider: provider,
		},
	}
	var i Index
	var e neoError
	c := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &i,
		Error:  &e,
	}
	status, err := nim.db.rc.Do(&c)
	if err != nil {
		return &i, err
	}
	if status != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", status)
		log.Printf("    Error:", e)
		return &i, BadResponse
	}
	return &i, nil
}
