// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"errors"
	"github.com/kr/pretty"
	"strings"
)

var (
	InvalidDatabase    = errors.New("Invalid database.  Check URI.")
	BadResponse        = errors.New("Bad response from Neo4j server.")
	NotFound           = errors.New("Cannot find in database.")
	FeatureUnavailable = errors.New("Feature unavailable")
	CannotDelete       = errors.New("The node cannot be deleted. Check that the node is orphaned before deletion.")
)

// A neoError is populated by api calls when there is an error.
type neoError struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
}

// Joins URL fragments
func join(fragments ...string) string {
	parts := []string{}
	for _, v := range fragments {
		v = strings.Trim(v, "/")
		parts = append(parts, v)
	}
	return strings.Join(parts, "/")
}

func logError(e *neoError) {
	pretty.Printf("%# v", e)
}
