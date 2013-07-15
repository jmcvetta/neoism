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
	idx, err := db.CreateNodeIndex("words", "", "")
	if err != nil {
		b.Fatal(err)
	}
	defer idx.Delete()
	findOrCreate := func(word string) *Node {
		nodes, err := idx.Find("word", word)
		if err != nil {
			b.Fatal(err)
		}
		var n *Node
		for _, n = range nodes {
			continue
		}
		if n != nil {
			return n
		}
		n, err = db.CreateNode(Properties{"word": word})
		if err != nil {
			b.Fatal(err)
		}
		err = idx.Add(n, "word", word)
		if err != nil {
			b.Fatal(err)
		}
		return n
	}
	//
	// Read book file
	//
	file, err := os.Open("tale-of-two-cities.txt")
	if err != nil {
		b.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var s0, s1 string
	var n0, n1 *Node
	//
	// Start Benchmark
	//
	b.StartTimer()
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, word := range fields {
			s0 = s1
			s1 = strings.Trim(word, ",-;!?/")
			if s0 == "" {
				continue
			}
		    if n1 != nil {
				n0 = n1
			} else {
				n0 = findOrCreate(s0)
				defer n0.Delete()
			}
			if s1 == "" {
				continue
			}
			n1 = findOrCreate(s1)
			defer n1.Delete()
			r0, err := n0.Relate("precedes", n1.Id(), nil)
			if err != nil {
				b.Fatal(err)
			}
			defer r0.Delete()
			r1, err := n1.Relate("follows", n0.Id(), nil)
			if err != nil {
				b.Fatal(err)
			}
			defer r1.Delete()
			fmt.Println(s0, "-->", s1)
		}
	}
	b.StopTimer() // Stop timer before return, so deferred deletes are not timed
}

