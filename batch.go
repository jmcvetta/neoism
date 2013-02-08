// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"fmt"
	"github.com/jmcvetta/restclient"
	"sync"
)

type job struct {
	BatchId        int                    `json:"id"`     // Identifies this job within its Batch
	Method         restclient.Method      `json:"method"` // HTTP Method to use for this job
	Url            string                 `json:"to"`     // Target URL
	Body           map[string]interface{} `json:"body"`   // Request body
	resultId       int                    // Identifies DB object created by executing this job
	resultTemplate interface{}
}

type Batch struct {
	db       *Database
	lock     *sync.Mutex  // lock protects queue
	queue    []*job       // Orderd queue of jobs
	jobs     map[int]*job // Map associating a job with its batchId value
	executed bool         // Has his batch been executed?
}

func (b *Batch) Add(j *job) int {
	b.lock.Lock()
	defer b.lock.Unlock()
	nextId := len(b.queue)
	j.BatchId = nextId
	b.queue[nextId] = j
	return nextId
}

// Execute sends all jobs in the queue to the DB as a single operation.
func (b *Batch) Execute() (result map[int]*entity, err error) {
	type respItem struct {
		BatchId  int         `json:"id"`
		Location string      `json:"location"`
		Body     interface{} `json:"body"`
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
	}
	return result, nil
}

// A BatchNode represents a Node created as part of a Batch job, which has not
// yet been instantiated in the db.  The BatchNode supports a limited subset of
// Node methods, all of which can be carried out as part of a Batch operation.
type BatchNode struct {
	batch *Batch
	id    int // ID within this batch
}

// NodeIdentity returns a string notation for referring to this BatchNode in
// other batch operations.
func (bn *BatchNode) NodeIdentity() string {
	return fmt.Sprintf("{ %v }", bn.id)
}

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

// Relationship returns the concrete Relationship object created by execution of
// this BatchRelationship, or an error if the Batch has not yet been executed.
func (br *BatchRelationship) Relationship() (*Relationship, error) {
	if br.batch.executed {
		return br.batch.db.Relationships.Get(br.job.resultId)
	}
	return nil, BatchNotExecuted
}

// A BatchRelationship represents a Relationship created as part of a Batch job,
// which has not yet been instantiated in the db.  The BatchRelationship
// supports a limited subset of Relationship methods, all of which can be
// carried out as part of a Batch operation.
type BatchRelationship struct {
	batch *Batch
	job   *job
}
