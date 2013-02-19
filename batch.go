// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"encoding/json"
	"github.com/jmcvetta/restclient"
	"sync"
)

// NewBatch returns a new, empty Batch
func (db *Database) NewBatch() *Batch {
	b := Batch{
		db: db,
		lock: new(sync.Mutex),
		queue: []*operation{},
		ops: make(map[int]*operation),
		executed: false,
	}
	return &b
}

type Batch struct {
	db       *Database
	lock     *sync.Mutex  // lock protects queue
	queue    []*operation       // Orderd queue of jobs
	ops     map[int]*operation // Map associating a job with its batchId value
	executed bool         // Has his batch been executed?
}

type operation struct {
	BatchId    int               `json:"id"`     // Identifies this operation within its Batch
	Method     restclient.Method `json:"method"` // HTTP Method to use for this operation
	Url        string            `json:"to"`     // Target URL
	Body       interface{}       `json:"body"`   // Request body
	resultJson *json.RawMessage  // JSON describing result of this operation.  Should this be a pointer?
	// resultId       int                    // Identifies DB object created by executing this operation
	// resultTemplate interface{}
	// resultEntity interface{}
}

// add puts an operations in the queue.
func (b *Batch) add(op *operation) int {
	b.lock.Lock()
	defer b.lock.Unlock()
	nextId := len(b.queue)
	op.BatchId = nextId
	b.ops[nextId] = op
	b.queue = append(b.queue, op)
	return nextId
}

// Execute sends all operations in the queue to the DB as a single operation.
func (b *Batch) Execute() (result map[int]*entity, err error) {
	type respItem struct {
		BatchId  int              `json:"id"`
		Location string           `json:"location"`
		Body     *json.RawMessage `json:"body"`
	}
	resp := make([]respItem, len(b.queue))
	ne := new(neoError)
	req := restclient.RestRequest{
		Url:    b.db.HrefBatch,
		Method: restclient.POST,
		Data:   b.queue,
		Error:  ne,
		Result: &resp,
	}
	status, err := b.db.Do(&req)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		logPretty(req)
		return nil, BadResponse
	}
	result = make(map[int]*entity, len(resp))
	for _, item := range resp {
		op := b.ops[item.BatchId]
		op.resultJson = item.Body
	}
	return result, nil
}
