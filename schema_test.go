// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	db := connectTest(t)
	defer cleanup(t, db)
	defer cleanupIndexes(t, db)
	label := rndStr(t)
	prop0 := rndStr(t)
	idx, err := db.CreateIndex(label, prop0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, label, idx.Label)
	assert.Equal(t, prop0, idx.PropertyKeys[0])
	_, err = db.CreateIndex("", "")
	assert.Equal(t, NotAllowed, err)
}

func TestIndexes(t *testing.T) {
	db := connectTest(t)
	defer cleanup(t, db)
	defer cleanupIndexes(t, db)
	label0 := rndStr(t)
	label1 := rndStr(t)
	prop0 := rndStr(t)
	prop1 := rndStr(t)
	indexes0, err := db.Indexes(label0)
	assert.Equal(t, err, nil)
	assert.Equal(t, 0, len(indexes0))
	_, err = db.CreateIndex(label0, prop0)
	assert.Equal(t, err, nil)
	_, err = db.CreateIndex(label1, prop0)
	assert.Equal(t, err, nil)
	_, err = db.CreateIndex(label1, prop1)
	assert.Equal(t, err, nil)
	indexes1, err := db.Indexes(label1)
	assert.Equal(t, err, nil)
	assert.Equal(t, 2, len(indexes1))
	indexes2, err := db.Indexes("")
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(indexes2))
}

func TestDropIndex(t *testing.T) {
	db := connectTest(t)
	defer cleanup(t, db)
	defer cleanupIndexes(t, db)
	label := rndStr(t)
	prop0 := rndStr(t)
	idx, _ := db.CreateIndex(label, prop0)
	indexes, _ := db.Indexes(label)
	assert.Equal(t, 1, len(indexes))
	err := idx.Drop()
	if err != nil {
		t.Fatal(err)
	}
	indexes, _ = db.Indexes(label)
	assert.Equal(t, 0, len(indexes))
	err = idx.Drop()
	assert.Equal(t, NotFound, err)
}

func cleanupIndexes(t *testing.T, db *Database) {
	labels, err := db.Labels()
	if err != nil {
		t.Fatal(err)
	}
	indexes := []*Index{}
	for _, l := range labels {
		idxs, err := db.Indexes(l)
		if err != nil {
			t.Fatal(err)
		}
		indexes = append(indexes, idxs...)
	}
	qs := make([]*CypherQuery, len(indexes))
	for i, idx := range indexes {
		// Cypher doesn't support properties in DROP statements
		l := idx.Label
		p := idx.PropertyKeys[0]
		stmt := fmt.Sprintf("DROP INDEX ON :%s(%s)", l, p)
		cq := CypherQuery{
			Statement: stmt,
		}
		qs[i] = &cq
	}
	// db.Rc.Log = true
	err = db.CypherBatch(qs)
	if err != nil {
		t.Fatal(err)
	}
}
