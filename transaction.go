// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"errors"
	"github.com/jmcvetta/restclient"
)

type txRequest struct {
	Statements []*CypherQuery `json:"statements"`
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
	ne := NeoError{}
	payload := txRequest{Statements: qs}
	res := txResponse{}
	rr := restclient.RequestResponse{
		Url:    db.HrefTransaction,
		Method: "POST",
		Data:   payload,
		Result: &res,
		Error:  &ne,
	}
	status, err := db.Rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status != 201 {
		return nil, ne
	}
	t := Tx{
		db:         db,
		hrefCommit: res.Commit,
		Location:   rr.HttpResponse.Header.Get("location"),
		Errors:     res.Errors,
	}
	err = res.unmarshal(qs)
	if err != nil {
		return &t, err
	}
	if len(t.Errors) != 0 {
		return &t, TxQueryError
	}
	return &t, err
}

// A Tx is an in-progress database transaction.
type Tx struct {
	db         *Database
	hrefCommit string
	Location   string
	Errors     []TxError
}

func (t *Tx) Commit() error {
	ne := NeoError{}
	rr := restclient.RequestResponse{
		Url:    t.hrefCommit,
		Method: "POST",
		Error:  &ne,
	}
	status, err := t.db.Rc.Do(&rr)
	if err != nil {
		return err
	}
	if status != 200 {
		return ne
	}
	return nil // Success
}
