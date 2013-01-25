// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"fmt"
	"github.com/bmizerany/assert"
	"log"
	"strconv"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// 18.9.1. Create node index
func TestCreateNodeIndex(t *testing.T) {
	name := rndStr(t)
	template := join(db.info.NodeIndex, name, "{key}/{value}")
	//
	// Create new index
	//
	idx0, err := db.Nodes.Indexes.Create(name)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.Name, name)
	assert.Equal(t, idx0.HrefTemplate, template)
	//
	// Get the index we just created
	//
	idx1, err := db.Nodes.Indexes.Get(name)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.Name, idx1.Name)
	//
	// Cleanup
	//
	idx0.Delete()
}

// 18.9.2. Create node index with configuration
func TestNodeIndexCreateWithConf(t *testing.T) {
	name := rndStr(t)
	indexType := "fulltext"
	provider := "lucene"
	template := join(db.info.NodeIndex, name, "{key}/{value}")
	//
	// Create new index
	//
	idx0, err := db.Nodes.Indexes.CreateWithConf(name, indexType, provider)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.IndexType, indexType)
	assert.Equal(t, idx0.Provider, provider)
	assert.Equal(t, idx0.HrefTemplate, template)
	assert.Equal(t, idx0.Name, name)
	//
	// Get the index we just created
	//
	idx1, err := db.Nodes.Indexes.Get(name)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.Name, idx1.Name)
	//
	// Cleanup
	//
	idx0.Delete()
}

// 18.9.3. Delete node index
func TestDeleteNodeIndex(t *testing.T) {
	// Include a space in the name to ensure correct URL escaping.
	name := rndStr(t) + " " + rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	err := idx0.Delete()
	if err != nil {
		t.Error(err)
	}
	_, err = db.Nodes.Indexes.Get(name)
	assert.Equal(t, err, NotFound)
}

// 18.9.4. List node indexes
func TestListNodeIndexes(t *testing.T) {
	name := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	indexes, err := db.Nodes.Indexes.All()
	if err != nil {
		t.Error(err)
	}
	valid := false
	for _, i := range indexes {
		if i.Name == name {
			valid = true
		}
	}
	assert.T(t, valid, "Newly created Index not found in listing of all Indexes.")
	//
	// Cleanup
	//
	idx0.Delete()
}

// 18.9.5. Add node to index
func TestAddNodeToIndex(t *testing.T) {
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	n0, _ := db.Nodes.Create(EmptyProps)
	err := idx0.Add(n0, key, value)
	if err != nil {
		t.Error(err)
	}
	//
	// Cleanup
	//
	n0.Delete()
	idx0.Delete()
}

// 18.9.6. Remove all entries with a given node from an index
func TestRemoveNodeFromIndex(t *testing.T) {
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	n0, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, "", "")
	if err != nil {
		t.Error(err)
	}
	// 
	// Cleanup
	//
	idx0.Delete()
	n0.Delete()
}

// 18.9.7. Remove all entries with a given node and key from an indexj
func TestRemoveNodeAndKeyFromIndex(t *testing.T) {
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	n0, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, key, "")
	if err != nil {
		t.Error(err)
	}
	// 
	// Cleanup
	//
	idx0.Delete()
	n0.Delete()
}

// 18.9.8. Remove all entries with a given node, key and value from an index
func TestRemoveNodeKeyAndValueFromIndex(t *testing.T) {
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	n0, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, key, "")
	if err != nil {
		t.Error(err)
	}
	//
	// Cleanup
	//
	n0.Delete()
	idx0.Delete()
}

// 18.9.9. Find node by exact match
func TestFindNodeByExactMatch(t *testing.T) {
	// Create
	idxName := rndStr(t)
	key0 := rndStr(t)
	key1 := rndStr(t)
	value0 := rndStr(t)
	value1 := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(idxName)
	n0, _ := db.Nodes.Create(EmptyProps)
	n1, _ := db.Nodes.Create(EmptyProps)
	n2, _ := db.Nodes.Create(EmptyProps)
	// These two will be located by Find() below
	idx0.Add(n0, key0, value0)
	idx0.Add(n1, key0, value0)
	// These two will NOT be located by Find() below
	idx0.Add(n2, key1, value0)
	idx0.Add(n2, key0, value1)
	//
	nodes, err := idx0.Find(key0, value0)
	if err != nil {
		t.Error(err)
	}
	// This query should have returned a map containing just two nodes, n1 and n0.
	assert.Equal(t, len(nodes), 2)
	_, present := nodes[n0.Id()]
	assert.Tf(t, present, "Find() failed to return node with id "+strconv.Itoa(n0.Id()))
	_, present = nodes[n1.Id()]
	assert.Tf(t, present, "Find() failed to return node with id "+strconv.Itoa(n1.Id()))
	// Cleanup
	n0.Delete()
	n1.Delete()
	n2.Delete()
	idx0.Delete()
}

// 18.9.10. Find node by query
func TestFindNodeByQuery(t *testing.T) {
	// Create
	idx0, _ := db.Nodes.Indexes.Create("test index")
	key0 := rndStr(t)
	key1 := rndStr(t)
	value0 := rndStr(t)
	value1 := rndStr(t)
	n0, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n0, key0, value0)
	idx0.Add(n0, key1, value1)
	n1, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n1, key0, value0)
	n2, _ := db.Nodes.Create(EmptyProps)
	idx0.Add(n2, rndStr(t), rndStr(t))
	// Retrieve
	luceneQuery0 := fmt.Sprintf("%v:%v AND %v:%v", key0, value0, key1, value1) // Retrieve n0 only
	luceneQuery1 := fmt.Sprintf("%v:%v", key0, value0)                         // Retrieve n0 and n1
	nodes0, err := idx0.Query(luceneQuery0)
	if err != nil {
		t.Error(err)
	}
	nodes1, err := idx0.Query(luceneQuery1)
	if err != nil {
		t.Error(err)
	}
	// Confirm
	assert.Equalf(t, len(nodes0), 1, "Query should have returned only one Node.")
	_, present := nodes0[n0.Id()]
	assert.Tf(t, present, "Query() failed to return node with id "+strconv.Itoa(n0.Id()))
	assert.Equalf(t, len(nodes1), 2, "Query should have returned exactly 2 Nodes.")
	_, present = nodes1[n0.Id()]
	assert.Tf(t, present, "Query() failed to return node with id "+strconv.Itoa(n0.Id()))
	_, present = nodes1[n1.Id()]
	assert.Tf(t, present, "Query() failed to return node with id "+strconv.Itoa(n1.Id()))
	// Cleanup
	idx0.Delete()
	n0.Delete()
	n1.Delete()
	n2.Delete()
}
