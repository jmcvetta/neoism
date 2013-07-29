// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
)

type indexRequest struct {
	PropertyKeys []string `json:"property_keys"`
}

// An Index improves the speed of looking up nodes in the database.
type Index struct {
	db *Database
	Label        string
	PropertyKeys []string `json:"property-keys"`
}

// Drop removes the index.
func (idx *Index) Drop() error {
	url := join(idx.db.Url, "schema/index", idx.Label, idx.PropertyKeys[0])
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    url,
		Method: "DELETE",
		Error:  &ne,
	}
	status, err := idx.db.Rc.Do(&rr)
	if err != nil {
		return err
	}
	if status == 404 {
		return NotFound
	}
	if status != 204 {
		return ne
	}
	return nil
}

// CreateIndex starts a background job in the database that will create and
// populate the new index of a specified property on nodes of a given label.
func (db *Database) CreateIndex(label, property string) (*Index, error) {
	url := join(db.Url, "schema/index", label)
	payload := indexRequest{[]string{property}}
	ne := NeoError{}
	res := Index{db: db}
	rr := restclient.RequestResponse{
		Url:    url,
		Method: "POST",
		Data:   payload,
		Result: &res,
		Error:  &ne,
	}
	status, err := db.Rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, NotFound
	}
	if status != 200 {
		return nil, ne
	}
	return &res, nil
}

// Indexes lists indexes for a label.
func (db *Database) Indexes(label string) ([]*Index, error) {
	url := join(db.Url, "schema/index", label)
	ne := NeoError{}
	res := []*Index{}
	rr := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &res,
		Error:  &ne,
	}
	status, err := db.Rc.Do(&rr)
	if err != nil {
		return res, err
	}
	if status == 404 {
		return res, NotFound
	}
	if status != 200 {
		return res, ne
	}
	for _, idx := range res {
		idx.db = db
	}
	return res, nil
}
