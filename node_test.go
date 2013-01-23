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

// 18.4.5. Delete node
func TestDeleteNode(t *testing.T) {
	// Create then delete a node
	n0, _ := db.Nodes.Create(emptyProps)
	id := n0.Id()
	n0.Delete()
	// Check that node is no longer in db
	_, err := db.Nodes.Get(id)
	assert.Equal(t, err, NotFound)
}

// 18.4.6. Nodes with relationships can not be deleted;
func TestDeleteNodeWithRelationships(t *testing.T) {
	// Create 
	n0, _ := db.Nodes.Create(emptyProps)
	n1, _ := db.Nodes.Create(emptyProps)
	r0, _ := n0.Relate("knows", n1.Id(), emptyProps)
	// Attempt to delete node without deleting relationship
	err := n0.Delete()
	assert.Equalf(t, err, CannotDelete, "Should not be possible to delete node with relationship.")
	// Cleanup
	r0.Delete()
	n0.Delete()
	n1.Delete()
}

// 18.7.1. Set property on node
func TestSetPropertyOnNode(t *testing.T) {
	// Create
	n0, _ := db.Nodes.Create(emptyProps)
	key := rndStr(t)
	value := rndStr(t)
	err := n0.SetProperty(key, value)
	if err != nil {
		t.Error(err)
	}
	// Confirm
	props, _ := n0.Properties()
	checkVal, present := props[key]
	assert.Tf(t, present, "Expected property key not found")
	assert.Tf(t, checkVal == value, "Expected property value not found")
	// Cleanup
	n0.Delete()
}

// 18.7.2. Update node properties
func TestUpdatePropertyOnNode(t *testing.T) {
	// Create
	props0 := Properties{rndStr(t): rndStr(t)}
	props1 := Properties{rndStr(t): rndStr(t)}
	n0, _ := db.Nodes.Create(props0)
	// Update
	err := n0.SetProperties(props1)
	if err != nil {
		t.Error(err)
	}
	// Confirm
	checkProps, _ := n0.Properties()
	assert.Equalf(t, props1, checkProps, "Did not recover expected properties after updating with SetProperties().")
	// Cleanup
	n0.Delete()
}

// 18.7.3. Get properties for node
func TestGetPropertiesForNode(t *testing.T) {
	// Create
	props := Properties{rndStr(t): rndStr(t)}
	n0, _ := db.Nodes.Create(props)
	// Get properties & confirm
	checkProps, err := n0.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, props, checkProps, "Did not return expected properties.")
	// Cleanup
	n0.Delete()
}

//
// 18.7.4. Property values can not be null
//
// This section cannot be tested.  Properties - which is a map[string]string -
// cannot be instantiated with a nil value.  If you try, the code will not compile.
//

//
// 18.7.5. Property values can not be nested
//
// This section cannot be tested.  Properties is defined as map[string]string -
// only strings may be used as values.  If you try to create a nested
// Properties, the code will not compile.
//

// 18.7.6. Delete all properties from node
func TestDeleteAllPropertiesFromNode(t *testing.T) {
	// Create
	props := Properties{
		rndStr(t): rndStr(t),
		rndStr(t): rndStr(t),
	}
	n0, _ := db.Nodes.Create(props)
	// Delete properties
	err := n0.DeleteProperties()
	if err != nil {
		t.Error(err)
	}
	// Confirm deletion
	checkProps, _ := n0.Properties()
	assert.Equalf(t, emptyProps, checkProps, "Properties should be empty after call to DeleteProperties()")
	// Cleanup
	n0.Delete()
}

// 18.7.7. Delete a named property from a node
func TestDeleteNamedPropertyFromNode(t *testing.T) {
	// Create
	props0 := Properties{"foo": "bar"}
	props1 := Properties{"foo": "bar", "spam": "eggs"}
	n0, _ := db.Nodes.Create(props1)
	// Delete
	err := n0.DeleteProperty("spam")
	if err != nil {
		t.Error(err)
	}
	// Confirm
	checkProps, _ := n0.Properties()
	assert.Equalf(t, props0, checkProps, "Failed to remove named property with DeleteProperty().")
	// Cleanup
	n0.Delete()
}
