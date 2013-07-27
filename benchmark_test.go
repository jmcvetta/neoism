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

// 18.5.2. Create relationship
func BenchmarkNodeChain(b *testing.B) {
	b.StopTimer()
	db := connectBench(b)
	b.StartTimer()
	n0, _ := db.CreateNode(Props{})
	defer n0.Delete()
	lastNode := n0
	for i := 0; i < b.N; i++ {
		nextNode, _ := db.CreateNode(Props{})
		defer nextNode.Delete()
		r0, _ := lastNode.Relate("knows", nextNode.Id(), Props{})
		defer r0.Delete()
	}
	b.StopTimer()
}
