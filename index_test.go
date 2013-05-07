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
	db := connectTest(t)
	name := rndStr(t)
	template := join(db.HrefNodeIndex, name, "{key}/{value}")
	//
	// Create new index
	//
	// idx0, err := db.Nodes.Indexes.Create(name)
	idx0, err := db.CreateNodeIndex(name, "", "")
	if err != nil {
		t.Error(err)
	}
	defer idx0.Delete()
	assert.Equal(t, idx0.Name, name)
	assert.Equal(t, idx0.HrefTemplate, template)
	//
	// Get the index we just created
	//
	idx1, err := db.NodeIndex(name)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.Name, idx1.Name)
}

// 18.9.2. Create node index with configuration
func TestNodeIndexCreateWithConf(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	indexType := "fulltext"
	provider := "lucene"
	template := join(db.HrefNodeIndex, name, "{key}/{value}")
	//
	// Create new index
	//
	idx0, err := db.CreateNodeIndex(name, indexType, provider)
	if err != nil {
		t.Error(err)
	}
	defer idx0.Delete()
	assert.Equal(t, idx0.IndexType, indexType)
	assert.Equal(t, idx0.Provider, provider)
	assert.Equal(t, idx0.HrefTemplate, template)
	assert.Equal(t, idx0.Name, name)
	//
	// Get the index we just created
	//
	idx1, err := db.NodeIndex(name)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, idx0.Name, idx1.Name)
}

// 18.9.3. Delete node index
func TestDeleteNodeIndex(t *testing.T) {
	db := connectTest(t)
	// Include a space in the name to ensure correct URL escaping.
	name := rndStr(t) + " " + rndStr(t)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	err := idx0.Delete()
	if err != nil {
		t.Error(err)
	}
	_, err = db.NodeIndex(name)
	assert.Equal(t, err, NotFound)
}

// 18.9.4. List node indexes
func TestListNodeIndexes(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(name)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	defer idx0.Delete()
	indexes, err := db.NodeIndexes()
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
}

// 18.9.5. Add node to index
func TestAddNodeToIndex(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(name)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	err := idx0.Add(n0, key, value)
	if err != nil {
		t.Error(err)
	}
}

// 18.9.6. Remove all entries with a given node from an index
func TestRemoveNodeFromIndex(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(name)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, "", "")
	if err != nil {
		t.Error(err)
	}
}

// 18.9.7. Remove all entries with a given node and key from an indexj
func TestRemoveNodeAndKeyFromIndex(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(name)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, key, "")
	if err != nil {
		t.Error(err)
	}
}

// 18.9.8. Remove all entries with a given node, key and value from an index
func TestRemoveNodeKeyAndValueFromIndex(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(name)
	idx0, _ := db.CreateNodeIndex(name, "", "")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	idx0.Add(n0, key, value)
	err := idx0.Remove(n0, key, "")
	if err != nil {
		t.Error(err)
	}
}

// 18.9.9. Find node by exact match
func TestFindNodeByExactMatch(t *testing.T) {
	db := connectTest(t)
	// Create
	idxName := rndStr(t)
	key0 := rndStr(t)
	key1 := rndStr(t)
	value0 := rndStr(t)
	value1 := rndStr(t)
	// idx0, _ := db.Nodes.Indexes.Create(idxName)
	idx0, _ := db.CreateNodeIndex(idxName, "", "")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	n1, _ := db.Nodes.Create(EmptyProps)
	defer n1.Delete()
	n2, _ := db.Nodes.Create(EmptyProps)
	defer n2.Delete()
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
}

// 18.9.10. Find node by query
func TestFindNodeByQuery(t *testing.T) {
	db := connectTest(t)
	// Create
	// idx0, _ := db.Nodes.Indexes.Create("test index")
	idx0, _ := db.CreateNodeIndex("test index", "", "")
	defer idx0.Delete()
	key0 := rndStr(t)
	key1 := rndStr(t)
	value0 := rndStr(t)
	value1 := rndStr(t)
	n0, _ := db.Nodes.Create(EmptyProps)
	defer n0.Delete()
	idx0.Add(n0, key0, value0)
	idx0.Add(n0, key1, value1)
	n1, _ := db.Nodes.Create(EmptyProps)
	defer n1.Delete()
	idx0.Add(n1, key0, value0)
	n2, _ := db.Nodes.Create(EmptyProps)
	defer n2.Delete()
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
}
