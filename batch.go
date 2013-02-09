// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"encoding/json"
	"github.com/jmcvetta/restclient"
	"sync"
)

type job struct {
	BatchId    int               `json:"id"`     // Identifies this job within its Batch
	Method     restclient.Method `json:"method"` // HTTP Method to use for this job
	Url        string            `json:"to"`     // Target URL
	Body       interface{}       `json:"body"`   // Request body
	resultJson *json.RawMessage  // JSON describing result of this job.  Should this be a pointer?
	// resultId       int                    // Identifies DB object created by executing this job
	// resultTemplate interface{}
	// resultEntity interface{}
}

type Batch struct {
	db       *Database
	lock     *sync.Mutex  // lock protects queue
	queue    []*job       // Orderd queue of jobs
	jobs     map[int]*job // Map associating a job with its batchId value
	executed bool         // Has his batch been executed?
}

// Add puts a job in the queue.
func (b *Batch) Add(j *job) int {
	b.lock.Lock()
	defer b.lock.Unlock()
	nextId := len(b.queue)
	j.BatchId = nextId
	b.jobs[nextId] = j
	b.queue = append(b.queue, j)
	return nextId
}

// Execute sends all jobs in the queue to the DB as a single operation.
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
		job := b.jobs[item.BatchId]
		job.resultJson = item.Body
	}
	return result, nil
}
