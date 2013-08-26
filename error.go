// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"errors"
)

// One of these errors is returned if we receive an error or unexpected response
// from the server.
var (
	InvalidDatabase = errors.New("Invalid database.  Check URI.")
	NotFound        = errors.New("Cannot find in database.")
	CannotDelete    = errors.New("The node cannot be deleted. Check that the node is orphaned before deletion.")
)

// A NeoError is populated by api calls when there is an error.
type NeoError struct {
	Message    string      `json:"message"`
	Exception  string      `json:"exception"`
	Stacktrace []string    `json:"stacktrace"`
	Cause      interface{} `json:"cause"` // New in Neo4j 2.0
}

// Error returns the error message supplied by the server.
func (ne NeoError) Error() string {
	return ne.Message
}
