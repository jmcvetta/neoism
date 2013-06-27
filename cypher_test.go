// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/bmizerany/assert"
	"strconv"
	"testing"
)

// 18.3.1. Send queries with parameters
func TestCypherParameters(t *testing.T) {
	var db *Database
	db = connectTest(t)
	// Create
	nameIdx, _ := db.CreateNodeIndex("name_index", "", "")
	defer nameIdx.Delete()
	floatIdx, _ := db.CreateNodeIndex("float_index", "", "")
	defer floatIdx.Delete()
	numIdx, _ := db.CreateNodeIndex("num_index", "", "")
	defer numIdx.Delete()
	n0, _ := db.CreateNode(Properties{"name": "I"})
	defer n0.Delete()
	nameIdx.Add(n0, "name", "I")
	n1, _ := db.CreateNode(Properties{"name": "you"})
	defer n1.Delete()
	n2, _ := db.CreateNode(Properties{"name": "num", "num": 42})
	defer n2.Delete()
	numIdx.Add(n2, "num", 42)
	n3, _ := db.CreateNode(Properties{"name": "float", "float": 3.14})
	defer n3.Delete()
	floatIdx.Add(n3, "float", 3.14)
	r0, _ := n0.Relate("knows", n1.Id(), nil)
	defer r0.Delete()
	r1, _ := n0.Relate("loves", n1.Id(), nil)
	defer r1.Delete()
	r2, _ := n0.Relate("understands", n2.Id(), nil)
	defer r2.Delete()
	//
	// Query with string parameters and integer results
	//
	query := `
		START n = node:name_index(name={startName})
		MATCH path = (n)-[r]->(m)
		WHERE m.name? = {name}
		RETURN id(n), id(r), id(m)
		`
	params := map[string]interface{}{
		"startName": "I",
		"name":      "you",
	}
	result0 := [][]int{}
	columns, err := db.Cypher(query, params, &result0)
	if err != nil {
		t.Error(err)
	}
	// Check result
	expCol := []string{"id(n)", "id(r)", "id(m)"}
	expDat0 := [][]int{
		[]int{n0.Id(), r0.Id(), n1.Id()},
		[]int{n0.Id(), r1.Id(), n1.Id()},
	}
	assert.Equal(t, expCol, columns)
	assert.Equal(t, expDat0, result0)
	//
	// Query with integer parameter and string results
	//
	query = `
		START n = node:num_index(num={num})
		RETURN n.name
		`
	params = map[string]interface{}{
		"num": 42,
	}
	result1 := [][]string{}
	columns, err = db.Cypher(query, params, &result1)
	if err != nil {
		t.Error(err)
	}
	expCol = []string{"n.name"}
	expDat1 := [][]string{[]string{"num"}}
	assert.Equal(t, expCol, columns)
	assert.Equal(t, expDat1, result1)
	//
	// Query with float parameter
	//
	query = `
		START n = node:float_index(float={float})
		RETURN n.name
		`
	params = map[string]interface{}{
		"float": 3.14,
	}
	result2 := [][]string{}
	columns, err = db.Cypher(query, params, &result2)
	if err != nil {
		t.Error(err)
	}
	expCol = []string{"n.name"}
	expDat2 := [][]string{[]string{"float"}}
	assert.Equal(t, expCol, columns)
	assert.Equal(t, expDat2, result2)
	//
	// Query with array parameter
	//
	query = `
		START n=node(*)
		WHERE id(n) IN {arr}
		RETURN n.name
		ORDER BY id(n)
		`
	params = map[string]interface{}{
		"arr": []int{n0.Id(), n1.Id()},
	}
	result3 := [][]string{}
	columns, err = db.Cypher(query, params, &result3)
	if err != nil {
		t.Error(err)
	}
	expCol = []string{"n.name"}
	expDat3 := [][]string{[]string{"I"}, []string{"you"}}
	assert.Equal(t, expCol, columns)
	assert.Equal(t, expDat3, result3)
}

// 18.3.2. Send a Query
func TestCypher(t *testing.T) {
	db := connectTest(t)
	// Create
	idx0, _ := db.CreateNodeIndex("name_index", "", "")
	defer idx0.Delete()
	n0, _ := db.CreateNode(Properties{"name": "I"})
	defer n0.Delete()
	idx0.Add(n0, "name", "I")
	n1, _ := db.CreateNode(Properties{"name": "you", "age": "69"})
	defer n1.Delete()
	r0, _ := n0.Relate("know", n1.Id(), nil)
	defer r0.Delete()
	// Query
	query := "start x = node(" + strconv.Itoa(n0.Id()) + ") match x -[r]-> n return type(r), n.name?, n.age?"
	// query := "START x = node:name_index(name=I) MATCH path = (x-[r]-friend) WHERE friend.name = you RETURN TYPE(r)"
	result := [][]string{}
	columns, err := db.Cypher(query, nil, &result)
	if err != nil {
		t.Error(err)
	}
	// Check result
	//
	// Our test only passes if Neo4j returns columns in the expected order - is
	// there any guarantee about order?
	expCol := []string{"type(r)", "n.name?", "n.age?"}
	expDat := [][]string{[]string{"know", "you", "69"}}
	assert.Equal(t, expCol, columns)
	assert.Equal(t, expDat, result)
}

func TestCypherBadQuery(t *testing.T) {
	db := connectTest(t)
	// Create
	idx0, _ := db.CreateNodeIndex("name_index", "", "")
	defer idx0.Delete()
	n0, _ := db.CreateNode(Properties{"name": "I"})
	defer n0.Delete()
	idx0.Add(n0, "name", "I")
	n1, _ := db.CreateNode(Properties{"name": "you", "age": "69"})
	defer n1.Delete()
	r0, _ := n0.Relate("know", n1.Id(), nil)
	defer r0.Delete()
	// Query
	query := "foobar("
	result := new(interface{})
	_, err := db.Cypher(query, nil, result)
	if err != BadResponse {
		t.Error(err)
	}
}
