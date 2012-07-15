// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"testing"
	"log"
)

func TestCreateNode(t *testing.T) {
	neo, err := NewDatabase("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(neo)
}


func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}