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
	db            *Database
	Name          string
	HrefTemplate  string
	Provider      string
	IndexType     string
	CaseSensitive bool
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
		logPretty(ne)
		return idx, err
	}
	if status != 201 {
		logPretty(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	return idx, nil
}

func (nim *NodeIndexManager) CreateWithConf(name, indexType, provider string) (*NodeIndex, error) {
	idx := new(NodeIndex)
	idx.db = nim.db
	idx.Name = name
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
		logPretty(ne)
		return idx, err
	}
	if status != 201 {
		logPretty(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	return idx, nil
}

func (nim *NodeIndexManager) All() ([]*NodeIndex, error) {
	res := map[string]nodeIndexResponse{}
	nis := []*NodeIndex{}
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    nim.db.info.NodeIndex,
		Method: restclient.GET,
		Result: &res,
		Error:  ne,
	}
	status, err := nim.do(&req)
	if err != nil {
		logPretty(ne)
		return nis, err
	}
	if status != 200 {
		logPretty(ne)
		return nis, BadResponse
	}
	for name, r := range res {
		n := NodeIndex{}
		n.db = nim.db
		n.Name = name
		n.populate(&r)
		nis = append(nis, &n)
	}
	return nis, nil
}

func (nim *NodeIndexManager) Get(name string) (*NodeIndex, error) {
	ni := new(NodeIndex)
	ni.Name = name
	ne := new(neoError)
	baseUri := nim.db.info.NodeIndex
	if baseUri == "" {
		return ni, FeatureUnavailable
	}
	uri := join(baseUri, name)
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Error:  ne,
	}
	status, err := nim.do(&req)
	if err != nil {
		logPretty(ne)
		return ni, err
	}
	switch status {
	// Success!
	case 200:
		return ni, nil
	case 400:
		return ni, NotFound
	}
	logPretty(ne)
	return ni, BadResponse
}

type nodeIndexResponse struct {
	HrefTemplate string `json:"template"`
	Provider     string `json:"provider"`      // Not always populated by server
	IndexType    string `json:"type"`          // Not always populated by server
	LowerCase    string `json:"to_lower_case"` // Not always populated by server
}

func (ni *NodeIndex) populate(res *nodeIndexResponse) {
	ni.HrefTemplate = res.HrefTemplate
	ni.Provider = res.Provider
	ni.IndexType = res.IndexType
	if res.LowerCase == "true" {
		ni.CaseSensitive = false
	} else {
		ni.CaseSensitive = true
	}
}
