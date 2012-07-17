// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/bmizerany/assert"
	"log"
	"testing"
)

// Buckets of properties for convenient testing
var (
	empty = Properties{}
	kirk  = Properties{"name": "kirk"}
	spock = Properties{"name": "spock"}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func connect(t *testing.T) *Database {
	//
	// Connect
	//
	db, err := NewDatabase("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCreateNode(t *testing.T) {
	db := connect(t)
	props := Properties{}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, props, p)
}

func TestCreateNodeProps(t *testing.T) {
	db := connect(t)
	props := Properties{"foo": "bar"}
	node, err := db.CreateNode(props)
	if err != nil {
		t.Fatal(err)
	}
	p, err := node.Properties()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, props, p)
}

func TestGetNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node0, _ := db.CreateNode(props)
	id := node0.Id()
	node1, err := db.GetNode(id)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, node0, node1)
}

func TestGetNonexistNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node0, _ := db.CreateNode(props)
	id := node0.Id()
	id = id + 50000 // Node with this id should (probably??) not yet exist
	_, err := db.GetNode(id)
	assert.Equal(t, err, NotFound)
}

func TestDeleteNode(t *testing.T) {
	db := connect(t)
	props := map[string]string{}
	node, _ := db.CreateNode(props)
	id := node.Id()
	err := node.Delete()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.GetNode(id)
	assert.Equal(t, err, NotFound)
}

func TestCreateRel(t *testing.T) {
	db := connect(t)
	props := Properties{}
	relProps := Properties{"this one goes to": "11"}
	node0, _ := db.CreateNode(props)
	node1, _ := db.CreateNode(props)
	rel, err := node0.Relate("knows", node1.Id(), relProps)
	if err != nil {
		t.Error(err)
		return
	}
	start, err := rel.Start()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, node0.Id(), start.Id())
	end, err := rel.End()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, node1, end)
	newRelProps, err := rel.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, relProps, newRelProps)
}

func createRelationship(t *testing.T, p Properties) *Relationship {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	rel, err := node0.Relate("knows", node1.Id(), p)
	if err != nil {
		t.Error(err)
	}
	return rel
}

func TestGetRelationship(t *testing.T) {
	db := connect(t)
	rel0 := createRelationship(t, empty)
	rel1, err := db.GetRelationship(rel0.Id())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, rel0, rel1)
}

func TestSetRelProps(t *testing.T) {
	rel := createRelationship(t, kirk)
	props, err := rel.Properties()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, kirk, props)
	err = rel.SetProperties(spock)
	if err != nil {
		t.Error(err)
	}
	props, _ = rel.Properties()
	assert.Equal(t, spock, props)
}

func TestGetRelProperty(t *testing.T) {
	rel := createRelationship(t, kirk)
	val0, err := rel.GetProperty("name")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val0, kirk["name"])
	_, err = rel.GetProperty("foobar")
	assert.Equal(t, NotFound, err)
}

func TestSetRelProperty(t *testing.T) {
	rel := createRelationship(t, kirk)
	err := rel.SetProperty("name", "mccoy")
	if err != nil {
		t.Error(err)
	}
	val, err := rel.GetProperty("name")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val, "mccoy")
	err = rel.SetProperty("spam", "eggs")
	if err != nil {
		t.Error(err)
	}
	val, err = rel.GetProperty("spam")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val, "eggs")
}

func TestGetAllRels(t *testing.T) {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	node2, _ := db.CreateNode(empty)
	node3, _ := db.CreateNode(empty)
	r0, _ := node0.Relate("knows", node1.Id(), kirk)
	r1, _ := node0.Relate("knows", node2.Id(), spock)
	rs, err := node0.AllRelationships()
	if err != nil {
		t.Error(err)
	}
	rels := []*Relationship{r0, r1}
	for _, v := range rels {
		_, ok := rs[v.Id()]
		if !ok {
			t.Errorf("Relationship ID %v not found in AllRelationships()", v.Id())
		}
	}
	// node3 has no relationships
	rs, err = node3.AllRelationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(rs))

}

func TestGetOutRels(t *testing.T) {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	node2, _ := db.CreateNode(empty)
	r0, _ := node0.Relate("knows", node1.Id(), kirk)
	r1, _ := node0.Relate("knows", node2.Id(), spock)
	rs, err := node0.OutgoingRelationships()
	if err != nil {
		t.Error(err)
	}
	rels := []*Relationship{r0, r1}
	for _, v := range rels {
		_, ok := rs[v.Id()]
		if !ok {
			t.Errorf("Relationship ID %v not found in OutgoingRelationships()", v.Id())
		}
	}
	// node1 has no outgoing relationships
	rs, err = node1.OutgoingRelationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(rs))
}

func TestGetInRels(t *testing.T) {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	r0, _ := node0.Relate("knows", node1.Id(), empty)
	// node0 has no incoming relationships
	rs, err := node0.IncomingRelationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(rs))
	// node1 has 1 incoming relationship, from node0
	rs, err = node1.IncomingRelationships()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(rs))
	_, ok := rs[r0.Id()]
	if !ok {
		t.Errorf("Relationship ID %v not found in OutgoingRelationships()", r0.Id())
	}
}

func TestTypedRels(t *testing.T) {
	db := connect(t)
	empty := Properties{}
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	node2, _ := db.CreateNode(empty)
	r0, _ := node0.Relate("knows", node1.Id(), kirk)
	node0.Relate("likes", node2.Id(), spock) // No need to capture the rel object, it won't be used
	// One "knows" relationship
	rs, err := node0.AllRelationships("knows")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(rs))
	_, ok := rs[r0.Id()]
	if !ok {
		t.Errorf("Relationship ID %v not found in OutgoingRelationships()", r0.Id())
	}
	// Two "knows" or "likes" relationships
	rs, err = node0.AllRelationships("knows", "likes")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(rs))
	// Zero "employs" relationships
	rs, err = node0.AllRelationships("employs")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(rs))
}

func TestRelTypes(t *testing.T) {
	// This test assumes only those types of relationship created by this test 
	// suite exist in the database.  If the test suite is run on a non-empty 
	// db, there is a good chance this test will fail because of that.
	knownTypes := []string{"knows", "likes"}
	db := connect(t)
	rts, err := db.RelationshipTypes()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, knownTypes, rts)
}
