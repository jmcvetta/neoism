// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"testing"
)

func TestTxBegin(t *testing.T) {
	db := connectTest(t)
	type name struct {
		Name string `json:"name"`
	}
	stmts := []*CypherStatement{
		&CypherStatement{
			Statement:  "CREATE (n:Person {props}) RETURN n",
			Parameters: map[string]interface{}{"props": map[string]string{"name": "James T Kirk"}},
			Data: &[][]struct {
				Name string
			}{},
		},
		&CypherStatement{
			Statement: "CREATE (m:Person {name: \"dr mccoy\"}) RETURN m",
		},
		&CypherStatement{
			Statement: `
				MATCH a:Person, b:Person
				WHERE a.name = "James T Kirk" AND b.name = "dr mccoy"
				CREATE a-[r:Commands]->b
				RETURN a, type(r) AS rel_type, b
			`,
			Parameters: map[string]interface{}{
				"n_name": "James T Kirk",
				"m_name": "dr mccoy",
			},
		},
	}
	_, err := db.BeginTx(stmts)
	if err != nil {
		t.Fatal(err)
	}
}
