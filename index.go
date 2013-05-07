// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
	"net/url"
	"strconv"
)

func (db *Database) createIndex(href, name, idxType, provider string) (*index, error) {
	idx := new(index)
	idx.db = db
	idx.Name = name
	type c struct {
		Type     string `json:"type,omitempty"`
		Provider string `json:"provider,omitempty"`
	}
	type p struct {
		Name   string `json:"name"`
		Config c      `json:"config,omitempty"`
	}
	payload := p{
		Name: name,
	}
	if idxType != "" || provider != "" {
		config := c{
			Type:     idxType,
			Provider: provider,
		}
		payload.Config = config
	}
	res := new(indexResponse)
	ne := new(neoError)
	rr := restclient.RequestResponse{
		Url:            href,
		Method:         "POST",
		Data:           &payload,
		Result:         res,
		Error:          ne,
		ExpectedStatus: 201,
	}
	status, err := db.rc.Do(&rr)
	if err != nil {
		log.Println(status)
		logPretty(ne)
		return idx, err
	}
	idx.populate(res)
	idx.HrefIndex = href
	return idx, nil
}

func (db *Database) indexes(href string) ([]*index, error) {
	res := map[string]indexResponse{}
	nis := []*index{}
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    href,
		Method: "GET",
		Result: &res,
		Error:  ne,
	}
	status, err := db.rc.Do(&req)
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
		n.db = db
		n.Name = name
		n.populate(&r)
		nis = append(nis, &n)
	}
	return nis, nil
}

func (db *Database) index(href, name string) (*index, error) {
	idx := new(index)
	resp := new(indexResponse)
	idx.Name = name
	baseUri := href
	rawurl := join(baseUri, name)
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return idx, err
	}
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    u.String(),
		Method: "GET",
		Error:  ne,
	}
	status, err := db.rc.Do(&req)
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

func (idx *index) populate(res *indexResponse) {
	idx.HrefTemplate = res.HrefTemplate
	idx.Provider = res.Provider
	idx.IndexType = res.IndexType
	if res.LowerCase == "true" {
		idx.CaseSensitive = false
	} else {
		idx.CaseSensitive = true
	}
}

type indexResponse struct {
	HrefTemplate string `json:"template"`
	Provider     string `json:"provider"`      // Not always populated by server
	IndexType    string `json:"type"`          // Not always populated by server
	LowerCase    string `json:"to_lower_case"` // Not always populated by server
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
		Method: "DELETE",
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
func (idx *index) add(e entity, key, value string) error {
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
		Uri:   e.HrefSelf(),
		Key:   key,
		Value: value,
	}
	req := restclient.RequestResponse{
		Url:    uri,
		Method: "POST",
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

func (idx *index) remove(e entity, key, value string) error {
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
	uri = join(uri, strconv.Itoa(e.Id()))
	ne := new(neoError)
	req := restclient.RequestResponse{
		Url:    uri,
		Method: "DELETE",
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

