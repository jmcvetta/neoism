// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/bmizerany/assert"
	"testing"
)

// 18.3.1. Send queries with parameters
func TestCypherSendQueryWithParameters(t *testing.T) {
	// Create
	idx0, _ := db.Nodes.Indexes.Create("name_index")
	n0, _ := db.Nodes.Create(Properties{"name": "I"})
	idx0.Add(n0, "name", "I")
	n1, _ := db.Nodes.Create(Properties{"name": "you"})
	// idx0.Add(n1, "name", "you")
	r0, _ := n0.Relate("know", n1.Id(), nil)
	r1, _ := n0.Relate("love", n1.Id(), nil)
	// Deferred Cleanup
	defer idx0.Delete()
	defer r0.Delete()
	defer r1.Delete()
	defer n0.Delete()
	defer n1.Delete()
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
