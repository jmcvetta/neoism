// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"errors"
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
