// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

package neo4j

import (
	"github.com/bmizerany/assert"
	"log"
	"sort"
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
	n0, _ := db.Nodes.Create(empty)
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
	n0, _ := db.Nodes.Create(empty)
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
	n0, _ := db.Nodes.Create(empty)
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
	n0, _ := db.Nodes.Create(empty)
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
	idxName := rndStr(t)
	key0 := rndStr(t)
	key1 := rndStr(t)
	value0 := rndStr(t)
	value1 := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(idxName)
	n0, _ := db.Nodes.Create(empty)
	n1, _ := db.Nodes.Create(empty)
	n2, _ := db.Nodes.Create(empty)
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
	// This query should have returned a slice containing just two nodes, n1 and n0.
	assert.Equal(t, len(nodes), 2)
	nodeIds := []int{}
	for _, aNode := range nodes {
		nodeIds = append(nodeIds, aNode.Id())
	}
	assert.Tf(t, sort.SearchInts(nodeIds, n0.Id()) < len(nodeIds),
		"Find() failed to return node with id "+strconv.Itoa(n0.Id()))
	assert.Tf(t, sort.SearchInts(nodeIds, n1.Id()) < len(nodeIds),
		"Find() failed to return node with id "+strconv.Itoa(n1.Id()))
	//
	// TODO: Test that n0 and n1 are members of nodes
	//
	// Cleanup
	//
	n0.Delete()
	n1.Delete()
	n2.Delete()
	idx0.Delete()
}
