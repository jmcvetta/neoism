// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

package neo4j

import (
	"github.com/bmizerany/assert"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// 18.9.1. Create node index
func TestCreateNodeIndex(t *testing.T) {
	db := connect(t)
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
}

// 18.9.2. Create node index with configuration
func TestNodeIndexCreateWithConf(t *testing.T) {
	db := connect(t)
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
}

// 18.9.4. List node indexes
func TestListNodeIndexes(t *testing.T) {
	db := connect(t)
	name := rndStr(t)
	db.Nodes.Indexes.Create(name)
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
}

// 18.9.3. Delete node index
func TestDeleteNodeIndex(t *testing.T) {
	db := connect(t)
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

// 18.9.5. Add node to index
func TestAddNodeToIndex(t *testing.T) {
	db := connect(t)
	name := rndStr(t)
	key := rndStr(t)
	value := rndStr(t)
	idx0, _ := db.Nodes.Indexes.Create(name)
	n0, _ := db.Nodes.Create(empty)
	err := idx0.Add(n0, key, value)
	if err != nil {
		t.Error(err)
	}
	idx0.Delete()
	n0.Delete()
}
