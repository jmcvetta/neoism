// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strconv"
	"strings"
)

type NodeIndexManager struct {
	db *Database
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (nim *NodeIndexManager) do(rr *restclient.RestRequest) (status int, err error) {
	return nim.db.rc.Do(rr)
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
	name = encodeSpaces(name)
	baseUri := nim.db.info.NodeIndex
	if baseUri == "" {
		return ni, FeatureUnavailable
	}
	uri := join(baseUri, name)
	ne := new(neoError)
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
	case 404:
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

type NodeIndex struct {
	db            *Database
	Name          string
	HrefTemplate  string
	Provider      string
	IndexType     string
	CaseSensitive bool
}

// encodeSpaces encodes spaces in a string as %20.
func encodeSpaces(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

// uri returns the URI for this Index.
func (ni *NodeIndex) uri() string {
	name := encodeSpaces(ni.Name)
	return join(ni.db.info.NodeIndex, name)
}

// Delete removes a NodeIndex from the database.
func (ni *NodeIndex) Delete() error {
	uri := ni.uri()
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  ne,
	}
	status, err := ni.db.rc.Do(&req)
	if err != nil {
		logPretty(ne)
		return err
	}
	if status == 204 {
		// Success!
		return nil
	}
	logPretty(ne)
	return BadResponse
}

// Add associates a Node with the given key/value pair in the given index.
func (ni *NodeIndex) Add(n *Node, key, value string) error {
	uri := ni.uri()
	ne := new(neoError)
	type s struct {
		Uri   string `json:"uri"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	data := s{
		Uri:   n.HrefSelf,
		Key:   key,
		Value: value,
	}
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.POST,
		Data:   data,
		Error:  ne,
	}
	status, err := ni.db.rc.Do(&req)
	if err != nil {
		logPretty(ne)
		return err
	}
	if status == 201 {
		// Success!
		return nil
	}
	logPretty(ne)
	return BadResponse
}

// Remove removes all entries with a given node, key and value from an index. 
// If value or both key and value may be the blank string, they are ignored.
func (ni *NodeIndex) Remove(n *Node, key, value string) error {
	// If key is an empty string, it will be ignored by join().  However it is only
	// valid to specify a value if key is non-empty.
	uri := join(ni.uri(), strconv.Itoa(n.Id()), key)
	if key != "" {
		uri = join(uri, value)
	}
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  ne,
	}
	status, err := ni.db.rc.Do(&req)
	if err != nil {
		logPretty(ne)
		return err
	}
	if status == 204 {
		// Success!
		return nil
	}
	logPretty(req)
	return BadResponse
}
