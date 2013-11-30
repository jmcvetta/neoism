// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

/*******************************************************************************

The Neo4j Manual section numbers quoted in the test suite refer to the manual
for milestone release 1.8.  http://docs.neo4j.org/chunked/1.8/

*******************************************************************************/

/*

To run these tests, a Neo4j 1.8 instance must be running on localhost on the
default port.

If the test suite complete successfully, all new nodes & relationships created
for testing will have been deleted in cleanup.  However if the test suite fails
to complete due to a panic, the cleanup code may not get called.  Therefore it
is not recommended to run this test suite on a db containing valuable data - run
it on a throwaway testing db instead! It's possible we could reduce this risk by
using defer() for cleanup.

*/

package neoism

import (
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/randutil"
	"log"
	"testing"
)

func connectTest(t *testing.T) *Database {
	log.SetFlags(log.Ltime | log.Lshortfile)
	db, err := Connect("http://localhost:7474/db/data")
	// db.Session.Log = true
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func cleanup(t *testing.T, db *Database) {
	qs := []*CypherQuery{
		&CypherQuery{
			Statement: `START r=rel(*) DELETE r`,
		},
		&CypherQuery{
			Statement: `START n=node(*) DELETE n`,
		},
	}
	err := db.CypherBatch(qs)
	if err != nil {
		t.Fatal(err)
	}
}

func rndStr(t *testing.T) string {
	// Neo4j doesn't like object names beginning with numerals.
	name, err := randutil.String(12, randutil.Alphabet)
	if err != nil {
		t.Fatal(err)
	}
	return name
}

func TestConnect(t *testing.T) {
	db := connectTest(t)
	assert.Equal(t, "http://localhost:7474/db/data", db.Url)
}

func TestConnectInvalidUrl(t *testing.T) {
	//
	//  Missing protocol scheme - url.Parse should fail
	//
	_, err := Connect("://foobar.com")
	if err == nil {
		t.Fatal("Expected error due to missing protocol scheme")
	}
	//
	// Unsupported protocol scheme - Session.Get should fail
	//
	_, err = Connect("foo://bar.com")
	if err == nil {
		t.Fatal("Expected error due to unsupported protocol scheme")
	}
	//
	// Not Found
	//
	_, err = Connect("http://localhost:7474/db/datadatadata")
	assert.Equal(t, InvalidDatabase, err)
	//
	// 200 Success and HTML returned
	//
	_, err = Connect("http://localhost:7474")
	assert.Equal(t, InvalidDatabase, err)
}
