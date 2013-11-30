// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

type indexRequest struct {
	PropertyKeys []string `json:"property_keys"`
}

// An Index improves the speed of looking up nodes in the database.
type Index struct {
	db           *Database
	Label        string
	PropertyKeys []string `json:"property_keys"`
}

// Drop removes the index.
func (idx *Index) Drop() error {
	uri := join(idx.db.Url, "schema/index", idx.Label, idx.PropertyKeys[0])
	ne := NeoError{}
	resp, err := idx.db.Session.Delete(uri, nil, &ne)
	if err != nil {
		return err
	}
	if resp.Status() == 404 {
		return NotFound
	}
	if resp.Status() != 204 {
		return ne
	}
	return nil
}

// CreateIndex starts a background job in the database that will create and
// populate the new index of a specified property on nodes of a given label.
func (db *Database) CreateIndex(label, property string) (*Index, error) {
	uri := join(db.Url, "schema/index", label)
	payload := indexRequest{[]string{property}}
	result := Index{db: db}
	ne := NeoError{}
	resp, err := db.Session.Post(uri, payload, &result, &ne)
	if err != nil {
		return nil, err
	}
	switch resp.Status() {
	case 200:
		return &result, nil // Success
	case 404:
		return nil, NotFound
	}
	return nil, ne
}

// Indexes lists indexes for a label.
func (db *Database) Indexes(label string) ([]*Index, error) {
	uri := join(db.Url, "schema/index", label)
	result := []*Index{}
	ne := NeoError{}
	resp, err := db.Session.Get(uri, nil, &result, &ne)
	if err != nil {
		return result, err
	}
	if resp.Status() == 404 {
		return result, NotFound
	}
	if resp.Status() != 200 {
		return result, ne
	}
	for _, idx := range result {
		idx.db = db
	}
	return result, nil
}
