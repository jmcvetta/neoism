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
	// Confirm relationship was created with specified properties.
	props1, _ := r0.Properties()
	assert.Equalf(t, props0, props1, "Properties queried from relationship do not match properties it was created with.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.4. Delete relationship
func TestDeleteRelationship(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, err := n0.Relate("knows", n1.Id(), emptyProps)
	if err != nil {
		t.Error(err)
	}
	// Delete and confirm
	r0.Delete()
	_, err = db.Relationships.Get(r0.Id())
	assert.Equalf(t, err, NotFound, "Should not be able to Get() a deleted relationship.")
	// Cleanup
	n0.Delete()
	n1.Delete()
}

// 18.5.5. Get all properties on a relationship
func TestGetAllPropertiesOnRelationship(t *testing.T) {
	// Create
	props0 := Properties{"foo": "bar", "spam": "eggs"}
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), props0)
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

// 18.5.6. Set all properties on a relationship
func TestSetAllPropertiesOnRelationship(t *testing.T) {
	props0 := Properties{"foo": "bar"}
	props1 := Properties{"spam": "eggs"}
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), props0)
	// Set all properties
	r0.SetProperties(props1)
	// Confirm
	checkProps, _ := r0.Properties()
	assert.Equalf(t, checkProps, props1, "Failed to set all properties on relationship")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.7. Get single property on a relationship
func TestGetSinglePropertyOnRelationship(t *testing.T) {
	// Create
	props := Properties{"foo": "bar"}
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), props)
	// Get property
	value, err := r0.Property("foo")
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, value, "bar", "Incorrect value when getting single property.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.8. Set single property on a relationship
func TestSetSinglePropertyOnRelationship(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Set property
	r0.SetProperty("foo", "bar")
	// Confirm
	expected := Properties{"foo": "bar"}
	props, _ := r0.Properties()
	assert.Equalf(t, props, expected, "Failed to set single property on relationship.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.9. Get all relationships
func TestGetAllRelationships(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	r1, _ := n1.Relate("knows", n0.Id(), emptyProps)
	r2, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Check relationships
	rels, err := n0.Relationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 3, "Wrong number of relationships")
	for _, r := range []*Relationship{r0, r1, r2} {
		_, present := rels[r.Id()]
		assert.Tf(t, present, "Missing expected relationship")
	}
	// Cleanup
	r0.Delete()
	r1.Delete()
	r2.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.10. Get incoming relationships
func TestGetIncomingRelationships(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	r1, _ := n1.Relate("knows", n0.Id(), emptyProps)
	r2, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Check relationships
	rels, err := n0.Incoming()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 1, "Wrong number of relationships")
	_, present := rels[r1.Id()]
	assert.Tf(t, present, "Missing expected relationship")
	// Cleanup
	r0.Delete()
	r1.Delete()
	r2.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.11. Get outgoing relationships
func TestGetOutgoingRelationships(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	r1, _ := n1.Relate("knows", n0.Id(), emptyProps)
	r2, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Check relationships
	rels, err := n0.Outgoing()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 2, "Wrong number of relationships")
	for _, r := range []*Relationship{r0, r2} {
		_, present := rels[r.Id()]
		assert.Tf(t, present, "Missing expected relationship")
	}
	// Cleanup
	r0.Delete()
	r1.Delete()
	r2.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.12. Get typed relationships
func TestGetTypedRelationships(t *testing.T) {
	// Create
	relType0 := rndStr(t)
	relType1 := rndStr(t)
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate(relType0, n1.Id(), emptyProps)
	r1, _ := n0.Relate(relType1, n1.Id(), emptyProps)
	// Check one type of relationship
	rels, err := n0.Relationships(relType0)
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 1, "Wrong number of relationships")
	_, present := rels[r0.Id()]
	assert.Tf(t, present, "Missing expected relationship")
	// Check two types of relationship together
	rels, err = n0.Relationships(relType0, relType1)
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 2, "Wrong number of relationships")
	for _, r := range []*Relationship{r0, r1} {
		_, present := rels[r.Id()]
		assert.Tf(t, present, "Missing expected relationship")
	}
	// Cleanup
	r0.Delete()
	r1.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.5.13. Get relationships on a node without relationships
func TestGetRelationshipsOnNodeWithoutRelationships(t *testing.T) {
	n0, _ := db.Nodes.Create(emptyProps)
	rels, err := n0.Relationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, len(rels), 0, "Node with no relationships should return empty slice of relationships")
}
