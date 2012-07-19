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
	db, err := Connect("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// Tests API described in Neo4j Manual section 19.3. Nodes
func TestNode(t *testing.T) {
	db := connect(t)
	//
	// 19.3.1. Create Node
	//
	node0, err := db.CreateNode(empty)
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.3.2. Create Node with properties
	//
	node1, err := db.CreateNode(kirk)
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.3.3. Get node
	//
	check, err := db.GetNode(node0.Id())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, node0.Id(), check.Id())
	// Make sure we can also get a node created w/ properties
	_, err = db.GetNode(node1.Id())
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.3.4. Get non-existent node
	//
	badId := node1.Id() + 1000000 // Probably does not exist yet
	_, err = db.GetNode(badId)
	assert.Equal(t, NotFound, err)
	//
	// 19.3.5. Delete node
	//
	n0Id := node0.Id()
	err = node0.Delete()
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.GetNode(n0Id) // Make sure it's really gone
	assert.Equal(t, NotFound, err)
	//
	// 19.3.6. Nodes with relationships can not be deleted
	//
	node2, err := db.CreateNode(empty)
	if err != nil {
		t.Fatal(err)
	}
	_, err = node1.Relate("knows", node2.Id(), empty)
	if err != nil {
		t.Fatal(err)
	}
	err = node1.Delete()
	assert.Equal(t, CannotDelete, err)
}

// Tests API described in Neo4j Manual section 19.4. Relationships
func TestRelationships(t *testing.T) {
	//
	// 19.4.2. Create relationship
	//
	// This section must precede 19.4.1. in order to have an object in the DB for us to Get
	db := connect(t)
	node0, _ := db.CreateNode(empty)
	node1, _ := db.CreateNode(empty)
	rel0, err := node0.Relate("knows", node1.Id(), empty)
	if err != nil {
		t.Fatal(err)
	}
	start, err := rel0.Start()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, node0.Id(), start.Id())
	end, err := rel0.End()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, node1, end)
	//
	// 19.4.1. Get Relationship by ID
	//
	clone, err := db.GetRelationship(rel0.Id())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rel0, clone)
	//
	// 19.4.3. Create a relationship with properties
	//
	rel1, err := node0.Relate("knows", node1.Id(), kirk)
	if err != nil {
		t.Fatal(err)
	}
	props, _ := rel1.Properties()
	assert.Equal(t, kirk, props)
	//
	// 19.4.4. Delete relationship
	//
	r0Id := rel0.Id()
	err = rel0.Delete()
	if err != nil {
		t.Fatal(err)
	}
	// Make sure it's gone:
	_, err = db.GetRelationship(r0Id)
	log.Println(err)
	assert.Equal(t, NotFound, err)
	//
	// 19.4.6. Set all properties on a relationship
	//
	rel2, err := node0.Relate("knows", node1.Id(), empty)
	err = rel2.SetProperties(kirk)
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.4.5. Get all properties on a relationship
	//
	props, err = rel2.Properties()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, kirk, props)
	//
	// 19.4.7. Get single property on a relationship
	//
	s, err := rel1.GetProperty("name")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "kirk", s)
	//
	// 19.4.8. Set single property on a relationship
	//
	rel3, err := node0.Relate("knows", node1.Id(), empty)
	err = rel3.SetProperty("name", "kirk")
	if err != nil {
		t.Fatal(err)
	}
	s, _ = rel3.GetProperty("name")
	assert.Equal(t, "kirk", s)
	//
	// 19.4.9. Get all relationships
	//
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

func TestRelSetProps(t *testing.T) {
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

func TestRelGetProperty(t *testing.T) {
	rel := createRelationship(t, kirk)
	val0, err := rel.GetProperty("name")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val0, kirk["name"])
	_, err = rel.GetProperty("foobar")
	assert.Equal(t, NotFound, err)
}

func TestRelSetProperty(t *testing.T) {
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
	rs, err := node0.Relationships()
	if err != nil {
		t.Error(err)
	}
	rels := []*Relationship{r0, r1}
	for _, v := range rels {
		_, ok := rs[v.Id()]
		if !ok {
			t.Errorf("Relationship ID %v not found in Relationships()", v.Id())
		}
	}
	// node3 has no relationships
	rs, err = node3.Relationships()
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
	rs, err := node0.Outgoing()
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
	rs, err = node1.Outgoing()
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
	rs, err := node0.Incoming()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(rs))
	// node1 has 1 incoming relationship, from node0
	rs, err = node1.Incoming()
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
	rs, err := node0.Relationships("knows")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(rs))
	_, ok := rs[r0.Id()]
	if !ok {
		t.Errorf("Relationship ID %v not found in OutgoingRelationships()", r0.Id())
	}
	// Two "knows" or "likes" relationships
	rs, err = node0.Relationships("knows", "likes")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(rs))
	// Zero "employs" relationships
	rs, err = node0.Relationships("employs")
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

func TestNodeSetGetProperty(t *testing.T) {
	db := connect(t)
	node0, _ := db.CreateNode(empty)
	node0.SetProperty("spam", "eggs")
	s, err := node0.GetProperty("spam")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "eggs", s)
}
