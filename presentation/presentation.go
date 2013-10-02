// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"log"
)

func connect() *neoism.Database {
	db, err := neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func create(db *neoism.Database) {
	kirk, err := db.CreateNode(neoism.Props{"name": "Kirk", "shirt": "yellow"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(kirk.Properties()) // Output: map[shirt:yellow name:Kirk] <nil>
	// Ignoring subsequent errors for brevity
	spock, _ := db.CreateNode(neoism.Props{"name": "Spock", "shirt": "blue"})
	mccoy, _ := db.CreateNode(neoism.Props{"name": "McCoy", "shirt": "blue"})
	r, _ := kirk.Relate("outranks", spock.Id(), nil) // No properties on this relationship
	start, _ := r.End()
	fmt.Println(start.Properties()) // Output: map[shirt:blue name:Spock] <nil>
	kirk.Relate("outranks", mccoy.Id(), nil)
	spock.Relate("outranks", mccoy.Id(), nil)
}

func cypher0(db *neoism.Database) *neoism.Node {
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}
	cq0 := neoism.CypherQuery{
		Statement: "CREATE (n {name: {crewman}, shirt: {shirt}}) RETURN n",
		// Use parameters instead of constructing a query string
		Parameters: neoism.Props{"crewman": "Scottie", "shirt": "red"},
		Result:     &res,
	}
	err := db.Cypher(&cq0)
	if err != nil {
		log.Fatal(err)
	}
	scottie := res[0].N // Only one row of data returned
	scottie.Db = db     // Must manually set Db with objects returned from Cypher query
	fmt.Println(scottie.Properties())
	// Output: map[shirt:red name:Scottie] <nil>
	return &scottie
}

func cypher1(db *neoism.Database, scottie *neoism.Node) {
	res := []struct {
		M   string `json:"m.name"` // `json` tag matches column name in query
		Rel string `json:"type(r)"`
		N   string `json:"n.name"`
	}{}
	cq1 := neoism.CypherQuery{
		// Use backticks for long statements - Cypher is whitespace indifferent
		Statement: `
				START n=node({id}), m=node(*)
				WHERE m.shirt = {color}
				CREATE (m)-[r:outranks]->(n)
				RETURN m.name, type(r), n.name
			`,
		Parameters: neoism.Props{"id": scottie.Id(), "color": "blue"},
		Result:     &res,
	}
	db.Cypher(&cq1)
	fmt.Println(res)
	// Output: [{Spock outranks Scottie} {McCoy outranks Scottie}]
}

func cypherBatch(db *neoism.Database) {

}

func main() {
	db := connect()
	defer cleanup(db)
	create(db)
	scottie := cypher0(db)
	cypher1(db, scottie)
}

func cleanup(db *neoism.Database) {
	qs := []*neoism.CypherQuery{
		&neoism.CypherQuery{
			Statement: `START r=rel(*) DELETE r`,
		},
		&neoism.CypherQuery{
			Statement: `START n=node(*) DELETE n`,
		},
	}
	err := db.CypherBatch(qs)
	if err != nil {
		log.Fatal(err)
	}
}
