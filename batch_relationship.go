// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
/*
	"fmt"
	"github.com/jmcvetta/restclient"
	"sync"
	"encoding/json"
*/
)

// A BatchRelationship represents a Relationship created as part of a Batch job,
// which has not yet been instantiated in the db.  The BatchRelationship
// supports a limited subset of Relationship methods, all of which can be
// carried out as part of a Batch operation.
type BatchRelationship struct {
	batch *Batch
	op   *operation
}

/*
// Relationship returns the concrete Relationship object created by execution of
// this BatchRelationship, or an error if the Batch has not yet been executed.
func (br *BatchRelationship) Relationship() (*Relationship, error) {
	if br.batch.executed {
		return br.batch.db.Relationships.Get(br.job.resultId)
	}
	return nil, BatchNotExecuted
}
*/
