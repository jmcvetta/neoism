// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

package neo4j

import (
	"github.com/bmizerany/assert"
	"testing"
)

// 18.5.1. Get Relationship by ID
func TestGetRelationshipById(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Get relationship
	r1, err := db.Relationships.Get(r0.Id())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, r0.Id(), r1.Id())
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.2. Create relationship
func TestCreateRelationship(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, err := n0.Relate("knows", n1.Id(), emptyProps)
	if err != nil {
		t.Error(err)
	}
	// Confirm relationship exists on both nodes
	rels, _ := n0.Outgoing("knows")
	_, present := rels[r0.Id()]
	assert.Tf(t, present, "Outgoing relationship not present on origin node.")
	rels, _ = n1.Incoming("knows")
	_, present = rels[r0.Id()]
	assert.Tf(t, present, "Incoming relationship not present on destination node.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.3. Create a relationship with properties
func TestCreateRelationshipWithProperties(t *testing.T) {
	// Create
	props0 := Properties{"foo": "bar", "spam": "eggs"}
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, err := n0.Relate("knows", n1.Id(), props0)
	if err != nil {
		t.Error(err)
	}
	// Confirm relationship was created with specified properties.  No need to
	// check success of creation itself, as that is handled by TestCreateRelationship().
	props1, err := r0.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, props0, props1, "Properties queried from relationship do not match properties it was created with.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}
