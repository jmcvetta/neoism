// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/bmizerany/assert"
	"strconv"
	"testing"
)

// 18.3.1. Send queries with parameters
func TestCypherSendQueryWithParameters(t *testing.T) {
	db := connectTest(t)
	// Create
	idx0, _ := db.Nodes.Indexes.Create("name_index")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(Properties{"name": "I"})
	defer n0.Delete()
	idx0.Add(n0, "name", "I")
	n1, _ := db.Nodes.Create(Properties{"name": "you"})
	defer n1.Delete()
	r0, _ := n0.Relate("know", n1.Id(), nil)
	defer r0.Delete()
	r1, _ := n0.Relate("love", n1.Id(), nil)
	defer r1.Delete()
	// Query
	query := "START x = node:name_index(name={startName}) MATCH path = (x-[r]-friend) WHERE friend.name = {name} RETURN TYPE(r)"
	params := map[string]string{
		"startName": "I",
		"name":      "you",
	}
	result, err := db.Cypher(query, params)
	if err != nil {
		t.Error(err)
	}
	// Check result
	expCol := []string{"TYPE(r)"}
	// Our test only passes if Neo4j returns "know" and "love" in this order.  Is
	// there any guarantee about order?  Can we modify the query to ensure order? 
	// Or is there a convenient way to sort result.Data here before checking it?
	expDat := [][]string{[]string{"know"}, []string{"love"}}
	assert.Equal(t, expCol, result.Columns)
	assert.Equal(t, expDat, result.Data)
}

// 18.3.2. Send a Query
func TestCypherSendQuery(t *testing.T) {
	db := connectTest(t)
	// Create
	idx0, _ := db.Nodes.Indexes.Create("name_index")
	defer idx0.Delete()
	n0, _ := db.Nodes.Create(Properties{"name": "I"})
	defer n0.Delete()
	idx0.Add(n0, "name", "I")
	n1, _ := db.Nodes.Create(Properties{"name": "you", "age": "69"})
	defer n1.Delete()
	r0, _ := n0.Relate("know", n1.Id(), nil)
	defer r0.Delete()
	// Query
	query := "start x = node(" + strconv.Itoa(n0.Id()) + ") match x -[r]-> n return type(r), n.name?, n.age?"
	// query := "START x = node:name_index(name=I) MATCH path = (x-[r]-friend) WHERE friend.name = you RETURN TYPE(r)"
	result, err := db.Cypher(query, nil)
	if err != nil {
		t.Error(err)
	}
	// Check result
	//
	// Our test only passes if Neo4j returns columns in the expected order - is 
	// there any guarantee about order?
	expCol := []string{"type(r)", "n.name?", "n.age?"}
	expDat := [][]string{[]string{"know", "you", "69"}}
	assert.Equal(t, expCol, result.Columns)
	assert.Equal(t, expDat, result.Data)
}
