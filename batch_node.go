// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/restclient"
)

func (b *Batch) CreateNode(p *Properties) *BatchNode {
	bn := new(BatchNode)
	bn.batch = b
	j := job{
		Url:    b.db.HrefNode,
		Method: restclient.POST,
		Body:   p,
	}
	bn.id = b.Add(&j)
	return bn
}

// A BatchNode represents a Node created as part of a Batch job, which has not
// yet been instantiated in the db.  The BatchNode supports a limited subset of
// Node methods, all of which can be carried out as part of a Batch operation.
type BatchNode struct {
	batch *Batch
	id    int // ID within this batch
	job   *job
}

// NodeIdentity returns a string notation for referring to this BatchNode in
// other batch operations.
func (bn *BatchNode) NodeIdentity() string {
	return fmt.Sprintf("{ %v }", bn.id)
}

// Node returns a Node object if the batch has been executed, or an error.
func (bn *BatchNode) Node() (*Node, error) {
	j := bn.job.resultJson
	if j == nil {
		return nil, BatchNotExecuted
	}
	nr := new(nodeResponse)
	err := json.Unmarshal([]byte(*j), nr)
	if err != nil {
		return nil, err
	}
	n := nr.Node(bn.batch.db)
	return n, nil
}

/*
// Relate creates a relationship of relType, with specified properties, 
// from this Node to the node identified by destId.
func (bn *BatchNode) Relate(relType string, dest NodeIdentifier, p Properties) *BatchRelationship {
	targetUri := join(bn.NodeIdentity(), "relationships")
	body := map[string]interface{}{
		"to":   dest.NodeIdentity(),
		"type": relType,
	}
	if p != nil {
		body["data"] = &p
	}
	j := job{
		Url:            targetUri,
		Method:         restclient.POST,
		Body:           body,
		resultTemplate: relationshipResponse{},
	}
	bn.batch.Add(&j)
	r := BatchRelationship{
		batch: bn.batch,
		job:   &j,
	}
	return &r
}
*/
