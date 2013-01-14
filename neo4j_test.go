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
	"sort"
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

// Tests API described in Neo4j Manual section 19.5. Relationship types
func TestRelationshipTypes(t *testing.T) {
	//
	// 19.5.1. Get relationship types
	//
	reltypes, err := db.Relationships.Types()
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"likes", "knows"}
	sort.Sort(sort.StringSlice(expected))
	assert.Equal(t, expected, reltypes)
}

// Tests API described in Neo4j Manual section 19.6. Node properties
func TestNodeProperties(t *testing.T) {
	//
	// 19.6.1. Set property on node
	//
	node0, _ := db.Nodes.Create(emptyProps)
	err := node0.SetProperty("name", "mccoy")
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.6.2. Update node properties
	//
	err = node0.SetProperties(spock)
	if err != nil {
		t.Fatal(err)
	}
	//
	// 19.6.3. Get properties for node
	//
	props, err := node0.Properties()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, spock, props)
	//
	// 19.6.4. Property values can not be null
	//
	// 19.6.5. Property values can not be nested
	//
	// These sections cannot be tested, because this library only accepts valid 
	// strings (the nil string, "", is still a valid string) as argument when 
	// setting properties.  It is not possible to write code that constructs an 
	// invalid request of this sort and still compiles.
	//
	// 19.6.6. Delete all properties from node
	//
	err = node0.DeleteProperties()
	if err != nil {
		t.Fatal(err)
	}
	props, _ = node0.Properties()
	assert.Equal(t, emptyProps, props)
	//
	// 19.6.7. Delete a named property from a node
	//
	node0.SetProperties(spock)
	node0.SetProperty("foo", "bar")
	node0.DeleteProperty("foo")
	if err != nil {
		t.Fatal(err)
	}
	props, _ = node0.Properties()
	assert.Equal(t, spock, props)
}
