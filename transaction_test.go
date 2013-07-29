// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"strconv"
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
	_, err := db.Begin(stmts)
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

func TestTxCommit(t *testing.T) {
	db := connectTest(t)
	name := rndStr(t)
	qs := []*CypherQuery{
		&CypherQuery{
			Statement: `
				CREATE (n:Person {name: {name}})
				RETURN n
			`,
			Parameters: Props{"name": name},
		},
	}
	tx, err := db.Begin(qs)
	if err != nil {
		t.Fatal(err)
	}
	//
	// Confirm it does not exist before commit
	//
	res0 := []struct {
		N string `json:"n.name"`
	}{}
	q0 := CypherQuery{
		Statement: `
			MATCH n:Person
			WHERE n.name = {name}
			RETURN n.name
		`,
		Parameters: Props{"name": name},
		Result:     &res0,
	}
	err = db.Cypher(&q0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(res0))
	//
	// Commit and confirm creation
	//
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}
	err = db.Cypher(&q0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res0))
	//
	// Cleanup
	//
	q1 := CypherQuery{
		Statement: `
			MATCH n:Person
			WHERE n.name = {name}
			DELETE n
		`,
		Parameters: Props{"name": name},
	}
	err = db.Cypher(&q1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTxBadResultObj(t *testing.T) {
	db := connectTest(t)
	//
	// Struct not slice of structs
	//
	res0 := struct{ N string }{}
	qs := []*CypherQuery{
		&CypherQuery{
			Statement: `CREATE (n:Person) RETURN n`,
			Result:    &res0,
		},
	}
	_, err := db.Begin(qs)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Fatal(err)
	}
}

func TestTxBadQuery(t *testing.T) {
	db := connectTest(t)
	qs := []*CypherQuery{
		&CypherQuery{
			Statement: `CREATE (n:Person) RETURN n`,
		},
		&CypherQuery{
			Statement: `CREATE (n:Person) RETURN n`,
		},
		&CypherQuery{
			Statement: `foobar`,
		},
		&CypherQuery{
			Statement: `CREATE (n:Person) RETURN n`,
		},
	}
	tx, err := db.Begin(qs)
	assert.Equal(t, TxQueryError, err)
	numErr := len(tx.Errors)
	assert.T(t, numErr == 1, "Expected one tx error, got "+strconv.Itoa(numErr))
}
