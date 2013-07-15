// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"testing"
	"bufio"
	"os"
	"fmt"
	"strings"
)

func BenchmarkTale(b *testing.B) {
	b.StopTimer()
	db, err := Connect("http://localhost:7474/db/data")
	if err != nil {
		b.Fatal(err)
	}
	file, err := os.Open("tale-of-two-cities.txt")
	if err != nil {
		b.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var prev, cur, next string
	// db.rc.Log = true
	//
	// Start Benchmark
	//
	b.StartTimer()
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, word := range fields {
			prev = cur
			cur = next
			next = word
			fmt.Println(prev, cur, next)
			if cur == "" {
				continue
			}
			curNode, err := db.CreateNode(Properties{"word": cur})
			if err != nil {
				b.Fatal(err)
			}
			defer curNode.Delete()
			if prev != "" {
				prevNode, err := db.CreateNode(Properties{"word": prev})
				if err != nil {
					b.Fatal(err)
				}
				defer prevNode.Delete()
				r, err := curNode.Relate("follows", prevNode.Id(), nil)
				if err != nil {
					b.Fatal(err)
				}
				defer r.Delete()
				fmt.Println(prev, "--follows-->", cur)
			}
			if next != "" {
				nextNode, err := db.CreateNode(Properties{"word": next})
				if err != nil {
					b.Fatal(err)
				}
				defer nextNode.Delete()
				r, err := curNode.Relate("precedes", nextNode.Id(), nil)
				if err != nil {
					b.Fatal(err)
				}
				defer r.Delete()
				fmt.Println(cur, "--precedes-->", next)
			}
		}
	}
	b.StopTimer() // Stop timer before return, so deferred deletes are not timed
}

