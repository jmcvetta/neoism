// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// +build !goci


package neo4j

import (
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	//
	// Connect
	//
	neo, err := NewDatabase("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(neo)
	//
	// Create
	//
	props := map[string]string{"foo": "bar"}
	node, err := neo.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(node)
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
