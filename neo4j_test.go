// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/bmizerany/assert"
	"log"
	"testing"
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
	assert.Equal(t, err, NodeNotFound)
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
	assert.Equal(t, err, NodeNotFound)
}

func TestCreateRel(t *testing.T) {
	db := connect(t)
	props := Properties{}
	node0, _ := db.CreateNode(props)
	node1, _ := db.CreateNode(props)
	log.Println("node0", node0)
	log.Println("node0.Info", node0.Info)
	log.Println("node1", node1)
	log.Println("node1.Info", node1.Info)
	rel, err := node0.Relate("knows", node1.Id(), props)
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
}

func TestGetRelationship(t *testing.T) {
	// TODO: Write get relationship test
}
