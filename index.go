// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"strconv"
	"strings"
)

type indexManager struct {
	HrefIndex string
	db        *Database
}

type NodeIndexManager struct {
	indexManager
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (im *indexManager) do(rr *restclient.RestRequest) (status int, err error) {
	return im.db.rc.Do(rr)
}

// CreateIndex creates a new Index with the supplied name.
func (im *indexManager) Create(name string) (*index, error) {
	type s struct {
		Name string `json:"name"`
	}
	data := s{Name: name}
	res := new(indexResponse)
	ne := new(neoError)
	idx := new(index)
	idx.db = im.db
	idx.Name = name
	rr := restclient.RestRequest{
		Url:    im.HrefIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: &res,
		Error:  &ne,
	}
	status, err := im.do(&rr)
	if err != nil {
		logPretty(ne)
		return idx, err
	}
	if status != 201 {
		logPretty(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	idx.HrefIndex = im.HrefIndex
	return idx, nil
}

// CreateIndexWithConf creates a new Index with the supplied name and configuration.
func (im *indexManager) CreateWithConf(name, indexType, provider string) (*index, error) {
	idx := new(index)
	idx.db = im.db
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
	res := new(indexResponse)
	ne := new(neoError)
	rr := restclient.RestRequest{
		Url:    im.HrefIndex,
		Method: restclient.POST,
		Data:   &data,
		Result: res,
		Error:  ne,
	}
	status, err := im.do(&rr)
	if err != nil {
		logPretty(ne)
		return idx, err
	}
	if status != 201 {
		logPretty(ne)
		return idx, BadResponse
	}
	idx.populate(res)
	idx.HrefIndex = im.HrefIndex
	return idx, nil
}

func (im *indexManager) All() ([]*index, error) {
	res := map[string]indexResponse{}
	nis := []*index{}
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    im.HrefIndex,
		Method: restclient.GET,
		Result: &res,
		Error:  ne,
	}
	status, err := im.do(&req)
	if err != nil {
		logPretty(ne)
		return nis, err
	}
	if status != 200 {
		logPretty(ne)
		return nis, BadResponse
	}
	for name, r := range res {
		n := index{}
		n.db = im.db
		n.Name = name
		n.populate(&r)
		nis = append(nis, &n)
	}
	return nis, nil
}

func (im *indexManager) Get(name string) (*index, error) {
	ni := new(index)
	ni.Name = name
	name = encodeSpaces(name)
	baseUri := im.HrefIndex
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
	status, err := im.do(&req)
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

type indexResponse struct {
	HrefTemplate string `json:"template"`
	Provider     string `json:"provider"`      // Not always populated by server
	IndexType    string `json:"type"`          // Not always populated by server
	LowerCase    string `json:"to_lower_case"` // Not always populated by server
}

func (ni *index) populate(res *indexResponse) {
	ni.HrefTemplate = res.HrefTemplate
	ni.Provider = res.Provider
	ni.IndexType = res.IndexType
	if res.LowerCase == "true" {
		ni.CaseSensitive = false
	} else {
		ni.CaseSensitive = true
	}
}

type index struct {
	db            *Database
	Name          string
	HrefTemplate  string
	Provider      string
	IndexType     string
	CaseSensitive bool
	HrefIndex     string
}

type NodeIndex struct {
	index
}

// encodeSpaces encodes spaces in a string as %20.
func encodeSpaces(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

// uri returns the URI for this Index.
func (ni *index) uri() string {
	name := encodeSpaces(ni.Name)
	return join(ni.HrefIndex, name)
}

// Delete removes a index from the database.
func (ni *index) Delete() error {
	uri := ni.uri()
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  ne,
	}
	status, err := ni.db.rc.Do(&req)
	if err != nil {
		logPretty(req)
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
func (ni *index) Add(n *Node, key, value string) error {
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
// If value or both key and value are the blank string, they are ignored.
func (ni *index) Remove(n *Node, key, value string) error {
	uri := ni.uri()
	// Since join() ignores fragments that are empty strings, joining an empty
	// value with a non-empty key produces a valid URL.  But joining a non-empty
	// value with an empty key would produce an invalid URL wherein they value is
	// conflated with the key.
	if key != "" {
		uri = join(uri, key, value)
	}
	uri = join(uri, strconv.Itoa(n.Id()))
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

// Find locates a node in the index by exact key/value match.
func (ni *index) Find(key, value string) ([]*Node, error) {
	key = encodeSpaces(key)
	value = encodeSpaces(value)
	uri := join(ni.uri(), key, value)
	ne := new(neoError)
	nodes := []*Node{}
	resp := []nodeResponse{}
	req := restclient.RestRequest{
		Url:    uri,
		Method: restclient.GET,
		Result: &resp,
		Error:  ne,
	}
	status, err := ni.db.rc.Do(&req)
	if err != nil {
		logPretty(ne)
		return nodes, err
	}
	if status != 200 {
		logPretty(req)
		return nodes, BadResponse
	}
	for _, r := range resp {
		n := Node{}
		n.db = ni.db
		n.populate(&r)
		nodes = append(nodes, &n)
	}
	return nodes, nil
}
