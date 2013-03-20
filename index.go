// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"net/url"
	"strconv"
)

type indexManager struct {
	HrefIndex string
	db        *Database
}

type NodeIndexManager struct {
	indexManager
}

type RelationshipIndexManager struct {
	indexManager
}

// do is a convenience wrapper around the embedded restclient's Do() method.
func (im *indexManager) do(rr *restclient.RequestResponse) (status int, err error) {
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
	rr := restclient.RequestResponse{
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
	rr := restclient.RequestResponse{
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
	req := restclient.RequestResponse{
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
	idx := new(index)
	resp := new(indexResponse)
	idx.Name = name
	baseUri := im.HrefIndex
	if baseUri == "" {
		return idx, FeatureUnavailable
	}
	rawurl := join(baseUri, name)
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return idx, err
	}
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    u.String(),
		Method: restclient.GET,
		Error:  ne,
	}
	status, err := im.do(&req)
	if err != nil {
		logPretty(req)
		return idx, err
	}
	switch status {
	// Success!
	case 200:
		idx.populate(resp)
		return idx, nil
	case 404:
		return idx, NotFound
	}
	logPretty(ne)
	return idx, BadResponse
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

type indexResponse struct {
	HrefTemplate string `json:"template"`
	Provider     string `json:"provider"`      // Not always populated by server
	IndexType    string `json:"type"`          // Not always populated by server
	LowerCase    string `json:"to_lower_case"` // Not always populated by server
}

// A NodeIndex is an index for searching Nodes.
type NodeIndex struct {
	index
}

// A RelationshipIndex is an index for searching Relationships.
type RelationshipIndex struct {
	index
}

// uri returns the URI for this Index.
func (idx *index) uri() (string, error) {
	s := join(idx.HrefIndex, idx.Name)
	u, err := url.ParseRequestURI(s)
	return u.String(), err
}

// Delete removes a index from the database.
func (idx *index) Delete() error {
	uri, err := idx.uri()
	if err != nil {
		return err
	}
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  ne,
	}
	status, err := idx.db.rc.Do(&req)
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
func (idx *index) Add(n *Node, key, value string) error {
	uri, err := idx.uri()
	if err != nil {
		return err
	}
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
	req := restclient.RequestResponse{
		Url:    uri,
		Method: restclient.POST,
		Data:   data,
		Error:  ne,
	}
	status, err := idx.db.rc.Do(&req)
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
func (idx *index) Remove(n *Node, key, value string) error {
	uri, err := idx.uri()
	if err != nil {
		return err
	}
	// Since join() ignores fragments that are empty strings, joining an empty
	// value with a non-empty key produces a valid URL.  But joining a non-empty
	// value with an empty key would produce an invalid URL wherein they value is
	// conflated with the key.
	if key != "" {
		uri = join(uri, key, value)
	}
	uri = join(uri, strconv.Itoa(n.Id()))
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    uri,
		Method: restclient.DELETE,
		Error:  ne,
	}
	status, err := idx.db.rc.Do(&req)
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

// A NodeMap associates Node objects with their integer IDs.
type NodeMap map[int]*Node

// Find locates Nodes in the index by exact key/value match.
func (idx *index) Find(key, value string) (NodeMap, error) {
	nm := make(NodeMap)
	rawurl, err := idx.uri()
	if err != nil {
		return nm, err
	}
	rawurl = join(rawurl, key, value)
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nm, err
	}
	ne := new(neoError)
	resp := []nodeResponse{}
	req := restclient.RequestResponse{
		Url:    u.String(),
		Method: restclient.GET,
		Result: &resp,
		Error:  ne,
	}
	status, err := idx.db.rc.Do(&req)
	if err != nil {
		logPretty(ne)
		return nm, err
	}
	if status != 200 {
		logPretty(req)
		return nm, BadResponse
	}
	for _, r := range resp {
		n := Node{}
		n.db = idx.db
		n.populate(&r)
		nm[n.Id()] = &n
	}
	return nm, nil
}

// Query locatess Nodes by query, in the query language appropriate for a given Index.
func (idx *index) Query(query string) (NodeMap, error) {
	nm := make(NodeMap)
	rawurl, err := idx.uri()
	if err != nil {
		return nm, err
	}
	v := make(url.Values)
	v.Add("query", query)
	rawurl += "?" + v.Encode()
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nm, err
	}
	result := []nodeResponse{}
	req := restclient.RequestResponse{
		Url:    u.String(),
		Method: restclient.GET,
		Result: &result,
	}
	status, err := idx.db.rc.Do(&req)
	if err != nil {
		return nm, err
	}
	if status != 200 {
		logPretty(req)
		return nm, BadResponse
	}
	for _, r := range result {
		n := Node{}
		n.db = idx.db
		n.populate(&r)
		nm[n.Id()] = &n
	}
	return nm, nil
}
