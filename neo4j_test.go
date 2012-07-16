// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/bmizerany/assert"
	"log"
	"testing"
)

// Buckets of properties for convenient testing
var (
	jCash     = Properties{"cash": "johnny"}
	hWilliams = Properties{"williams": "hank"}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func connect(t *testing.T) *Database {
	//
	// Connect
	//
	db, err := NewDatabase("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCreateNode(t *testing.T) {
	db := connect(t)
	props := Properties{}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, props, p)
}

func TestCreateNodeProps(t *testing.T) {
	db := connect(t)
	props := Properties{"foo": "bar"}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, props, p)
}

func TestGetNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node0, _ := db.CreateNode(props)
	id := node0.Id()
	node1, err := db.GetNode(id)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, node0, node1)
}

func TestGetNonexistNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node0, _ := db.CreateNode(props)
	id := node0.Id()
	id = id + 50000 // Node with this id should (probably??) not yet exist
	_, err := db.GetNode(id)
	assert.Equal(t, err, NotFound)
}

func TestDeleteNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node, _ := db.CreateNode(props)
	id := node.Id()
	err := node.Delete()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.GetNode(id)
	assert.Equal(t, err, NotFound)
}

func TestCreateRel(t *testing.T) {
	db := connect(t)
	props := Properties{}
	relProps := Properties{"this one goes to": "11"}
	node0, _ := db.CreateNode(props)
	node1, _ := db.CreateNode(props)
	rel, err := node0.Relate("knows", node1.Id(), relProps)
	if err != nil {
		t.Error(err)
		return
	}
	start, err := rel.Start()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, node0.Id(), start.Id())
	end, err := rel.End()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, node1, end)
	newRelProps, err := rel.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, relProps, newRelProps)
}

func createRelationship(t *testing.T, p Properties) *Relationship {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	rel, err := node0.Relate("knows", node1.Id(), p)
	if err != nil {
		t.Error(err)
	}
	return rel
}

func TestGetRelationship(t *testing.T) {
	db := connect(t)
	jCash := Properties{"cash": "johnny"}
	rel0 := createRelationship(t, jCash)
	rel1, err := db.GetRelationship(rel0.Id())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, rel0, rel1)
}

func TestSetRelProps(t *testing.T) {
	rel := createRelationship(t, jCash)
	props, err := rel.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, jCash, props)
	err = rel.SetProperties(hWilliams)
	if err != nil {
		t.Error(err)
	}
	props, _ = rel.Properties()
	assert.Equal(t, hWilliams, props)
}

func TestGetRelProperty(t *testing.T) {
	rel := createRelationship(t, jCash)
	val0, err := rel.GetProperty("cash")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val0, jCash["cash"])
	_, err = rel.GetProperty("foobar")
	assert.Equal(t, NotFound, err)
}

func TestSetRelProperty(t *testing.T) {
	rel := createRelationship(t, jCash)
	err := rel.SetProperty("cash", "money")
	if err != nil {
		t.Error(err)
	}
	val, err := rel.GetProperty("cash")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val, "money")
	err = rel.SetProperty("spam", "eggs")
	if err != nil {
		t.Error(err)
	}
	val, err = rel.GetProperty("spam")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val, "eggs")
}

func TestGetAllRels(t *testing.T) {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	node2, _ := db.CreateNode(empty)
	rel0, _ := node0.Relate("knows", node1.Id(), jCash)
	rel1, _ := node0.Relate("knows", node2.Id(), hWilliams)
	rs, err := node0.AllRelationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, *rel0, rs[rel0.Id()])
	assert.Equal(t, *rel1, rs[rel1.Id()])
}
