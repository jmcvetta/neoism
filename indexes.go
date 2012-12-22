// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
)

type IndexManager interface {
	Create(name string) *Index
	CreateWithConf(name, indexType, provider string) (*Index, error)
}

type Index interface {
	Name() string                          // Common Name of this index
	Template() string                      // Template for making REST calls to this Index
	Add(e Entity, key, value string) error // Add an entity to this index under key:value
}

type NodeIndex struct {
	template  string
	provider  string
	indexType string
}

func (ni *NodeIndex) Name() string {
	return "foobar!"
}

func (ni *NodeIndex) Template() string {
	return ni.template
}
func (ni *NodeIndex) Add(e Entity, key, value string) error {
	return nil
}

type nodeIndexManager struct {
	db *Database
}

// indexCreateResp is suitable for unmarshalling the JSON response from an index create operation.
type indexCreateResp struct {
	Template string `json:"template"`
	Provider string `json:"provider"`
	Type     string `json:"type"`
}

// CreateIndex creates a new Index, with the name supplied, in the db.
func (nim *nodeIndexManager) Create(name string) (Index, error) {
	var ni NodeIndex
	type req struct {
		Name string `json:"name"`
	}
	data := req{Name: name}
	var r indexCreateResp
	var e neoError
	c := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &r,
		Error:  &e,
	}
	status, err := nim.db.rc.Do(&c)
	if err != nil {
		return &ni, err
	}
	if status != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", status)
		log.Printf("    Error:", e)
		return &ni, BadResponse
	}
	ni.template = r.Template
	return &ni, nil
}

func (nim *nodeIndexManager) CreateWithConf(name, indexType, provider string) (Index, error) {
	var ni NodeIndex
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
	var r indexCreateResp
	var e neoError
	c := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &r,
		Error:  &e,
	}
	status, err := nim.db.rc.Do(&c)
	if err != nil {
		return &ni, err
	}
	if status != 201 {
		log.Printf("Unexpected response from server:")
		log.Printf("    Response code:", status)
		log.Printf("    Error:", e)
		return &ni, BadResponse
	}
	ni.template = r.Template
	ni.provider = r.Provider
	ni.indexType = r.Type
	return &ni, nil
}
