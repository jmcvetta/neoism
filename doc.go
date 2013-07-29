package neo4j

/*

Package neo4j is a client library providing access to the Neo4j graph database
via its REST API.


Example Usage:

	package main

	import (
		"fmt"
		"github.com/jmcvetta/neo4j"
	)

	func main() {
		// No error handling in this example - bad, bad, bad!
		//
		// Connect to the Neo4j server
		//
		db, _ := neo4j.Connect("http://localhost:7474/db/data")
		kirk := "Captain Kirk"
		mccoy := "Dr McCoy"
		//
		// Create a node
		//
		n0, _ := db.CreateNode(neo4j.Props{"name": kirk})
		defer n0.Delete()  // Deferred clean up
		//
		// Create a node with a Cypher query
		//
		res0 := []struct {
			N neo4j.Node // Column "n" gets automagically unmarshalled into field N
		}{}
		cq0 := neo4j.CypherQuery{
			Statement: "CREATE (n:Person {name: {name}}) RETURN n",
			// Use parameters instead of constructing a query string
			Parameters: neo4j.Props{"name": mccoy},
			Result:     &res0,
		}
		db.Cypher(&cq0)
		n1 := res0[0].N
		n1.Db = db // Must manually set Db with objects returned from Cypher query
		//
		// Create a relationship
		//
		n1.Relate("reports to", n0.Id(), neo4j.Props{}) // Empty Props{} is okay
		//
		// Issue a query
		//
		res1 := []struct {
			A   string `json:"a.name"` // `json` tag matches column name in query
			Rel string `json:"type(r)"`
			B   string `json:"b.name"`
		}{}
		cq1 := neo4j.CypherQuery{
			// Use backticks for long statements - Cypher is whitespace indifferent
			Statement: `
				MATCH (a:Person)-[r]->(b)
				WHERE a.name = {name}
				RETURN a.name, type(r), b.name
			`,
			Parameters: neo4j.Props{"name": mccoy},
			Result:     &res1,
		}
		db.Cypher(&cq1)
		r := res1[0]
		fmt.Println(r.A, r.Rel, r.B)
		//
		// Clean up using a transaction
		//
		qs := []*neo4j.CypherQuery{
			&neo4j.CypherQuery{
				Statement: `
					MATCH (n:Person)-[r]->()
					WHERE n.name = {name}
					DELETE r
				`,
				Parameters: neo4j.Props{"name": mccoy},
			},
			&neo4j.CypherQuery{
				Statement: `
					MATCH n:Person
					WHERE n.name = {name}
					DELETE n
				`,
				Parameters: neo4j.Props{"name": mccoy},
			},
		}
		tx, _ := db.Begin(qs)
		tx.Commit()
	}
*/
