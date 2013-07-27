// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/bmizerany/assert"
	"testing"
)

type resStruct0 struct {
	N struct {
		Name string
	}
}

type resStruct1 struct {
	M map[string]string
}

type resStruct2 struct {
	A   string `json:"a.name"`
	Rel string `json:"type(r)"`
	B   struct {
		Name string
	} `json:"b"`
}

func TestTxBegin(t *testing.T) {
	db := connectTest(t)
	type name struct {
		Name string `json:"name"`
	}
	res0 := []resStruct0{}
	res1 := []resStruct1{}
	res2 := []resStruct2{}
	q0 := CypherQuery{
		Statement:  "CREATE (n:Person {props}) RETURN n",
		Parameters: map[string]interface{}{"props": map[string]string{"name": "James T Kirk"}},
		Result:     &res0,
	}
	q1 := CypherQuery{
		Statement: "CREATE (m:Person {name: \"Dr McCoy\"}) RETURN m",
		Result:    &res1,
	}
	q2 := CypherQuery{
		Statement: `
				MATCH a:Person, b:Person
				WHERE a.name = "James T Kirk" AND b.name = "Dr McCoy"
				CREATE a-[r:Commands]->b
				RETURN a.name, type(r), b
			`,
		Parameters: map[string]interface{}{
			"n_name": "James T Kirk",
			"m_name": "dr mccoy",
		},
		Result: &res2,
	}

	assert.Equal(t, *new([]string), q1.Columns())
	stmts := []*CypherQuery{&q0, &q1, &q2}
	_, err := db.BeginTx(stmts)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res0))
	assert.Equal(t, "James T Kirk", res0[0].N.Name)
	assert.Equal(t, 1, len(res1))
	assert.Equal(t, "Dr McCoy", res1[0].M["name"])
	assert.Equal(t, 1, len(res2))
	assert.Equal(t, "James T Kirk", res2[0].A)
	assert.Equal(t, "Commands", res2[0].Rel)
	assert.Equal(t, "Dr McCoy", res2[0].B.Name)
}
