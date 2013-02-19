// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/jmcvetta/restclient"
)

// CreateNode creates a new Node as part of a Batch
func (b *Batch) CreateNode(p *Properties) *BatchNode {
	bn := new(BatchNode)
	bn.batch = b
	op := operation{
		Url:    b.db.HrefNode,
		Method: restclient.POST,
		Body:   p,
	}
	bn.id = b.add(&op)
	return bn
}

// DeleteNode deletes the specified Node as part of a Batch
func (b *Batch) DeleteNode(id int) {
	u := join(b.db.HrefNode, strconv.Itoa(id))
	op := operation{
		Url:    u,
		Method: restclient.DELETE,
	}
	b.add(&op)
}

// A BatchNode represents a Node created as part of a Batch job, which has not
// yet been instantiated in the db.  The BatchNode supports a limited subset of
// Node methods, all of which can be carried out as part of a Batch operation.
type BatchNode struct {
	batch *Batch
	id    int // ID within this batch
	op   *operation
}

// NodeIdentity returns a string notation for referring to this BatchNode in
// other batch operations.
func (bn *BatchNode) NodeIdentity() string {
	return fmt.Sprintf("{ %v }", bn.id)
}

// Node returns a Node object if the batch has been executed, or an error.
func (bn *BatchNode) Node() (*Node, error) {
	op := bn.op.resultJson
	if op == nil {
		return nil, BatchNotExecuted
	}
	nr := new(nodeResponse)
	err := json.Unmarshal([]byte(*op), nr)
	if err != nil {
		return nil, err
	}
	n := nr.Node(bn.batch.db)
	return n, nil
}

