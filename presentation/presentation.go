// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package main

import (
	"github.com/jmcvetta/neoism"
	"log"
)

func connect() *neoism.Database {
	db, err := neoism.Connect("localhost:7474")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func create() {
}
