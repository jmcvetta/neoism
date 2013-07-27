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

// unmarshall populates a slice of CypherQuery object with result data returned
// from the server.
func (tr *txResponse) unmarshall(qs []*CypherQuery) error {
	if len(tr.Results) != len(qs) {
		return errors.New("Result count does not match query count")
	}
	for i, s := range qs {
		s.cr = tr.Results[i]
		if s.Result != nil {
			err := s.Unmarshall(s.Result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// BeginTx opens a new transaction, executing zero or more cypher queries
// inside the transaction.
func (db *Database) BeginTx(qs []*CypherQuery) (*Tx, error) {
	ne := new(neoError)
	payload := txRequest{Statements: qs}
	res := txResponse{}
	rr := restclient.RequestResponse{
		Url:            db.HrefTransaction,
		Method:         "POST",
		Data:           payload,
		Result:         &res,
		Error:          &ne,
		ExpectedStatus: 201,
	}
	_, err := db.rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	t := Tx{
		Location: rr.HttpResponse.Header.Get("location"),
		Commit:   res.Commit,
		Errors:   res.Errors,
	}
	err = res.unmarshall(qs)
	return &t, err
}

// A Tx is an in-progress database transaction.
type Tx struct {
	Location string
	Commit   string
	Errors   []TxError
}
