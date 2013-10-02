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

func transaction(db *neoism.Database) {
	res := []struct {
		M   string `json:"m.name"` // `json` tag matches column name in query
		Rel string `json:"type(r)"`
		N   string `json:"n.name"`
	}{}
	qs := []*neoism.CypherQuery{
		&neoism.CypherQuery{
			Statement:  "CREATE (n {name: {crewman}, shirt: {shirt}}) RETURN n",
			Parameters: neoism.Props{"crewman": "Scottie", "shirt": "red"},
		},
		&neoism.CypherQuery{
			Statement: `START n=node(*), m=node(*)
				WHERE n.name = {name}, m.shirt = {color}
				CREATE (m)-[r:outranks]->(n)
				RETURN m.name, type(r), n.name`,
			Parameters: neoism.Props{"name": "Scottie", "color": "blue"},
			Result:     &res,
		},
	}
	// db.Session.Log = true
	tx, _ := db.Begin(qs)
	fmt.Println(tx)
	tx.Commit()
}

func cypherBatch(db *neoism.Database) {

}

func main() {
	db := connect()
	defer cleanup(db)
	create(db)
	transaction(db)
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
