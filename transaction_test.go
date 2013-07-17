// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"testing"
)

/*
['CREATE (n:Person {name: "james t kirk"}) RETURN n',
 'CREATE (m:Person {name: "dr mccoy"}) RETURN m',
 'CREATE (q:Person {name: "spock"})',
 'START n=node(*) RETURN n',
 'START n=node(*) MATCH n-[r]->m RETURN n,r,m']

*/

func TestTxBegin(t *testing.T) {
	db := connectTest(t)
	type name struct {
		Name string `json:"name"`
	}
	stmts := []*CypherStatement{
		&CypherStatement{
			Statement:  "CREATE (n:Person {props}) RETURN n",
			Parameters: map[string]interface{}{"props": map[string]string{"name": "James T Kirk"}},
			Data:       [][]string{},
		},
		&CypherStatement{
			Statement: "CREATE (m:Person {name: \"dr mccoy\"}) RETURN m",
		},
	}
	_, err := db.BeginTx(stmts)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range stmts {
		logPretty(s)
	}
}
