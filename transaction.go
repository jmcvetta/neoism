// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"errors"
)

// A Tx is an in-progress database transaction.
type Tx struct {
	db         *Database
	hrefCommit string
	Location   string
	Errors     []TxError
	Expires    string // Cannot unmarshall into time.Time :(
}

// A TxQueryError is returned when there is an error with one of the Cypher
// queries inside a transaction, but not with the transaction itself.
var TxQueryError = errors.New("Error with a query inside a transaction.")

// A TxError is an error with one of the statements submitted in a transaction,
// but not with the transaction itself.
type TxError struct {
	Code    int
	Status  string
	Message string
}

type txRequest struct {
	Statements []*CypherQuery `json:"statements"`
}

type txResponse struct {
	Commit      string
	Results     []cypherResult
	Transaction struct {
		Expires string
	}
	Errors []TxError
}

// unmarshal populates a slice of CypherQuery object with result data returned
// from the server.
func (tr *txResponse) unmarshal(qs []*CypherQuery) error {
	for i, res := range tr.Results {
		q := qs[i]
		q.cr = res
		if q.Result != nil {
			err := q.Unmarshal(q.Result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Begin opens a new transaction, executing zero or more cypher queries
// inside the transaction.
func (db *Database) Begin(qs []*CypherQuery) (*Tx, error) {
	payload := txRequest{Statements: qs}
	result := txResponse{}
	resp, err := db.Session.Post(db.HrefTransaction, payload, &result, nil)
	if err != nil {
		return nil, err
	}
	if resp.Status() != 201 {
		ne := NeoError{}
		resp.Unmarshal(&ne)
		return nil, ne
	}
	t := Tx{
		db:         db,
		hrefCommit: result.Commit,
		Location:   resp.HttpResponse().Header.Get("location"),
		Errors:     result.Errors,
		Expires:    result.Transaction.Expires,
	}
	err = result.unmarshal(qs)
	if err != nil {
		return &t, err
	}
	if len(t.Errors) != 0 {
		return &t, TxQueryError
	}
	return &t, err
}

// Commit commits an open transaction.
func (t *Tx) Commit() error {
	if len(t.Errors) > 0 {
		return TxQueryError
	}
	resp, err := t.db.Session.Post(t.hrefCommit, nil, nil, nil)
	if err != nil {
		return err
	}
	if resp.Status() != 200 {
		ne := NeoError{}
		resp.Unmarshal(&ne)
		return ne
	}
	return nil // Success
}

// Query executes statements in an open transaction.
func (t *Tx) Query(qs []*CypherQuery) error {
	payload := txRequest{Statements: qs}
	result := txResponse{}
	resp, err := t.db.Session.Post(t.Location, payload, &result, nil)
	if err != nil {
		return err
	}
	if resp.Status() == 404 {
		return NotFound
	}
	if resp.Status() != 200 {
		ne := NeoError{}
		resp.Unmarshal(&ne)
		return &ne
	}
	t.Expires = result.Transaction.Expires
	t.Errors = append(t.Errors, result.Errors...)
	err = result.unmarshal(qs)
	if err != nil {
		return err
	}
	if len(t.Errors) != 0 {
		return TxQueryError
	}
	return nil
}

// Rollback rolls back an open transaction.
func (t *Tx) Rollback() error {
	resp, err := t.db.Session.Delete(t.Location, nil)
	if err != nil {
		return err
	}
	if resp.Status() != 200 {
		ne := NeoError{}
		resp.Unmarshal(&ne)
		return ne
	}
	return nil // Success
}
