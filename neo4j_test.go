// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// +build !goci


package neo4j

import (
	"log"
	"testing"
	"github.com/bmizerany/assert"
)

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
	props := map[string]string{}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, props, p)
}

func TestCreateNodeProps(t *testing.T) {
	db := connect(t)
	props := map[string]string{"foo": "bar"}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, props, p)
}

func TestGetNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node0, _ := db.CreateNode(props)
	id := node0.Id()
	log.Println("id:", id)
	node1, err := db.GetNode(id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, node0, node1)
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
