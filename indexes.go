// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
)

type NodeIndexManager struct {
	db *Database
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (nim *NodeIndexManager) do(rr *restclient.RestRequest) (status int, err error) {
	return nim.db.rc.Do(rr)
}

type NodeIndex struct {
	db           *Database
	Name         string
	HrefTemplate string
	HrefProvider string
	HrefType     string
}

func (ni *NodeIndex) populate(res *nodeIndexResponse) {
	ni.HrefTemplate = res.HrefTemplate
	ni.HrefProvider = res.HrefProvider
	ni.HrefType = res.HrefType
}

type nodeIndexResponse struct {
	HrefTemplate string `json:"template"`
	HrefProvider string `json:"provider"`
	HrefType     string `json:"type"`
}

// CreateIndex creates a new Index, with the name supplied, in the db.
func (nim *NodeIndexManager) Create(name string) (*NodeIndex, error) {
	type s struct {
		Name string `json:"name"`
	}
	data := s{Name: name}
	res := new(nodeIndexResponse)
	ne := new(neoError)
	idx := new(NodeIndex)
	idx.db = nim.db
	idx.Name = name
	rr := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &res,
		Error:  &ne,
	}
	status, err := nim.do(&rr)
	if err != nil {
		logError(ne)
		return idx, err
	}
	if status != 201 {
		logError(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	return idx, nil
}

func (nim *NodeIndexManager) CreateWithConf(name, indexType, provider string) (*NodeIndex, error) {
	idx := new(NodeIndex)
	idx.db = nim.db
	type conf struct {
		Type     string `json:"type"`
		Provider string `json:"provider"`
	}
	type s struct {
		Name   string `json:"name"`
		Config conf   `json:"config"`
	}
	data := s{
		Name: name,
		Config: conf{
			Type:     indexType,
			Provider: provider,
		},
	}
	res := new(nodeIndexResponse)
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: res,
		Error:  ne,
	}
	status, err := nim.do(&rr)
	if err != nil {
		logError(ne)
		return idx, err
	}
	if status != 201 {
		logError(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	return idx, nil
}
