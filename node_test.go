// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

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
	db := connectTest(t)
	// Create
	n0, err := db.CreateNode(nil)
	if err != nil {
		t.Error(err)
	}
	defer n0.Delete()
	// Confirm creation
	_, err = db.Node(n0.Id())
	if err != nil {
		t.Error(err)
	}
}

// 18.4.2. Create Node with properties
func TestCreateNodeWithProperties(t *testing.T) {
	db := connectTest(t)
	// Create
	props0 := Props{}
	props0["foo"] = "bar"
	props0["spam"] = "eggs"
	n0, err := db.CreateNode(props0)
	if err != nil {
		t.Error(err)
	}
	defer n0.Delete()
	// Confirm creation
	_, err = db.Node(n0.Id())
	if err != nil {
		t.Error(err)
	}
	// Confirm properties
	props1, _ := n0.Properties()
	assert.Equalf(t, props0, props1, "Node properties not as expected")
}

// 18.4.3. Get node
func TestGetNode(t *testing.T) {
	db := connectTest(t)
	// Create
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
	// Get Node
	n1, err := db.Node(n0.Id())
	if err != nil {
		t.Error(err)
	}
	// Confirm nodes are the same
	assert.Equalf(t, n0.Id(), n1.Id(), "Nodes do not have same ID")
}

// 18.4.4. Get non-existent node
func TestGetNonexistentNode(t *testing.T) {
	db := connectTest(t)
	// Create a node
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
	// Try to get non-existent node with next Id
	implausible := n0.Id() + 1000
	_, err := db.Node(implausible)
	assert.Equal(t, err, NotFound)
}

// 18.4.5. Delete node
func TestDeleteNode(t *testing.T) {
	db := connectTest(t)
	// Create then delete a node
	n0, _ := db.CreateNode(Props{})
	id := n0.Id()
	err := n0.Delete()
	if err != nil {
		t.Error(err)
	}
	// Check that node is no longer in db
	_, err = db.Node(id)
	assert.Equal(t, err, NotFound)
}

// 18.4.6. Nodes with relationships can not be deleted;
func TestDeleteNodeWithRelationships(t *testing.T) {
	db := connectTest(t)
	// Create
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
	n1, _ := db.CreateNode(Props{})
	defer n1.Delete()
	r0, _ := n0.Relate("knows", n1.Id(), Props{})
	defer r0.Delete()
	// Attempt to delete node without deleting relationship
	err := n0.Delete()
	assert.Equalf(t, err, CannotDelete, "Should not be possible to delete node with relationship.")
}

// 18.7.1. Set property on node
func TestSetPropertyOnNode(t *testing.T) {
	db := connectTest(t)
	// Create
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
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
}

// 18.7.1. Set property on node
func TestSetBadPropertyOnNode(t *testing.T) {
	db := connectTest(t)
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
	key := ""
	value := rndStr(t)
	err := n0.SetProperty(key, value)
	if _, ok := err.(NeoError); !ok {
		t.Fatal(err)
	}
}

// 18.7.2. Update node properties
func TestUpdatePropertyOnNode(t *testing.T) {
	db := connectTest(t)
	// Create
	props0 := Props{rndStr(t): rndStr(t)}
	props1 := Props{rndStr(t): rndStr(t)}
	n0, _ := db.CreateNode(props0)
	defer n0.Delete()
	// Update
	err := n0.SetProperties(props1)
	if err != nil {
		t.Error(err)
	}
	// Confirm
	checkProps, _ := n0.Properties()
	assert.Equalf(t, props1, checkProps, "Did not recover expected properties after updating with SetProperties().")
}

// 18.7.3. Get properties for node
func TestGetPropertiesForNode(t *testing.T) {
	db := connectTest(t)
	// Create
	props := Props{rndStr(t): rndStr(t)}
	n0, _ := db.CreateNode(props)
	defer n0.Delete()
	// Get properties & confirm
	checkProps, err := n0.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, props, checkProps, "Did not return expected properties.")
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
	db := connectTest(t)
	// Create
	props := Props{
		rndStr(t): rndStr(t),
		rndStr(t): rndStr(t),
	}
	n0, _ := db.CreateNode(props)
	defer n0.Delete()
	// Delete properties
	err := n0.DeleteProperties()
	if err != nil {
		t.Error(err)
	}
	// Confirm deletion
	checkProps, _ := n0.Properties()
	assert.Equalf(t, Props{}, checkProps, "Properties should be empty after call to DeleteProperties()")
	n0.Delete()
	err = n0.DeleteProperties()
	assert.Equal(t, NotFound, err)
}

// 18.7.7. Delete a named property from a node
func TestDeleteNamedPropertyFromNode(t *testing.T) {
	db := connectTest(t)
	// Create
	props0 := Props{"foo": "bar"}
	props1 := Props{"foo": "bar", "spam": "eggs"}
	n0, _ := db.CreateNode(props1)
	defer n0.Delete()
	// Delete
	err := n0.DeleteProperty("spam")
	if err != nil {
		t.Error(err)
	}
	// Confirm
	checkProps, _ := n0.Properties()
	assert.Equalf(t, props0, checkProps, "Failed to remove named property with DeleteProperty().")
	//
	// Delete non-existent property
	//
	err = n0.DeleteProperty("eggs")
	assert.NotEqual(t, nil, err)
	//
	// Delete and check 404
	//
	n0.Delete()
	err = n0.DeleteProperty("spam")
	assert.Equal(t, NotFound, err)
}

func TestNodeProperty(t *testing.T) {
	db := connectTest(t)
	props := Props{"foo": "bar"}
	n0, _ := db.CreateNode(props)
	defer n0.Delete()
	value, err := n0.Property("foo")
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, value, "bar", "Incorrect value when getting single property.")
	//
	// Check Not Found
	//
	n0.Delete()
	_, err = n0.Property("foo")
	assert.Equal(t, NotFound, err)
}
