// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	// "log"
)

type Batch struct {
	*Database
	Queue []*restclient.RestRequest
}

// Add queues a RestRequest for later execution as part of this batch.
func (b *Batch) Add(r *restclient.RestRequest) {
	b.Queue = append(b.Queue, r)
}
