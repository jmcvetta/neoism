// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

package neo4j

import (
	"github.com/bmizerany/assert"
	// "github.com/jmcvetta/randutil"
	// "log"
	// "sort"
	"testing"
)

// 18.4.1. Create Node
func TestCreateNode(t *testing.T) {
	// Create
	n0, err := db.Nodes.Create(emptyProps)
	if err != nil {
		t.Error(err)
	}
	// Confirm creation
	_, err = db.Nodes.Get(n0.Id())
	if err != nil {
		t.Error(err)
	}
	// Cleanup
	n0.Delete()
}

// 18.4.2. Create Node with properties
func TestCreateNodeWithProperties(t *testing.T) {
	// Create
	props0 := Properties{}
	props0["foo"] = "bar"
	props0["spam"] = "eggs"
	n0, err := db.Nodes.Create(props0)
	if err != nil {
		t.Error(err)
	}
	// Confirm creation
	_, err = db.Nodes.Get(n0.Id())
	if err != nil {
		t.Error(err)
	}
	// Confirm properties
	props1, _ := n0.Properties()
	assert.Equalf(t, props0, props1, "Node properties not as expected")
	// Cleanup
	n0.Delete()
}

// 18.4.3. Get node
func TestGetNode(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	// Get Node
	n1, err := db.Nodes.Get(n0.Id())
	if err != nil {
		t.Error(err)
	}
	// Confirm nodes are the same
	assert.Equalf(t, n0.Id(), n1.Id(), "Nodes do not have same ID")
	// Cleanup
	n0.Delete()
}

// 18.4.4. Get non-existent node
func TestGetNonexistentNode(t *testing.T) {
	// Create a node
	n0, _ := db.Nodes.Create(emptyProps)
	// Try to get non-existent node with next Id
	implausible := n0.Id() + 1000
	_, err := db.Nodes.Get(implausible)
	assert.Equal(t, err, NotFound)
	// Cleanup
	n0.Delete()
}
