// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"log"
	"testing"
)

func connectBench(b *testing.B) *Database {
	log.SetFlags(log.Ltime | log.Lshortfile)
	db, err := Connect("http://localhost:7474/db/data")
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func BenchmarkNodeChain(b *testing.B) {
	b.StopTimer()
	db := connectBench(b)
	b.StartTimer()
	n0, _ := db.CreateNode(Props{"name": 0})
	defer n0.Delete()
	lastNode := n0
	for i := 1; i < b.N; i++ {
		nextNode, _ := db.CreateNode(Props{"name": i})
		defer nextNode.Delete()
		r0, _ := lastNode.Relate("knows", nextNode.Id(), Props{"name": i})
		defer r0.Delete()
	}
	b.StopTimer()
}

func BenchmarkNodeChainBatch(b *testing.B) {
	b.StopTimer()
	db := connectBench(b)
	b.StartTimer()
	qs := []*CypherQuery{}
	nodes := []int{}
	rels := []int{}
	cq := CypherQuery{
		Statement:  `CREATE (n:Person {name: {i}}) RETURN n`,
		Parameters: Props{"i": 0},
		Result:     &[]Node{},
	}
	qs = append(qs, &cq)
	nodes = append(nodes, 0)
	for i := 1; i < b.N; i++ {
		cq0 := CypherQuery{
			Statement:  `CREATE (n:Person {name: {i}}) RETURN n`,
			Parameters: Props{"i": i},
			Result:     &[]Node{},
		}
		qs = append(qs, &cq0)
		nodes = append(nodes, i)
		cq1 := CypherQuery{
			Statement:  `MATCH a:Person, b:Person WHERE a.name = {i} AND b.name = {k} CREATE a-[r:Knows {name: {i}}]->b RETURN id(r)`,
			Parameters: Props{"i": i, "k": i - 1},
			Result:     &[]Relationship{},
		}
		qs = append(qs, &cq1)
		rels = append(rels, i)
	}
	err := db.CypherBatch(qs)
	if err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
	//
	// Cleanup
	//
	qs = []*CypherQuery{}
	for _, r := range rels {
		cq := CypherQuery{
			Statement:  `MATCH ()-[r:Knows]->() WHERE r.name = {i} DELETE r`,
			Parameters: Props{"i": r},
		}
		qs = append(qs, &cq)
	}
	for _, n := range nodes {
		cq := CypherQuery{
			Statement:  `MATCH n:Person WHERE n.name = {i} DELETE n`,
			Parameters: Props{"i": n},
		}
		qs = append(qs, &cq)
	}
	err = db.CypherBatch(qs)
	if err != nil {
		b.Fatal(err)
	}
}
