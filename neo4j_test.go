// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

package neo4j

import (
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/randutil"
	"log"
	"testing"
)

// Database connection used by all tests
var db *Database

// Buckets of properties for convenient testing
var (
	emptyProps = Properties{}
	kirk       = Properties{"name": "kirk"}
	spock      = Properties{"name": "spock"}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	var err error
	db, err = Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Panic(err)
	}
}

func rndStr(t *testing.T) string {
	name, err := randutil.AlphaString(12)
	if err != nil {
		t.Fatal(err)
	}
	return name
}
