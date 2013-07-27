// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"encoding/json"
	"errors"
	"github.com/jmcvetta/restclient"
)

// A CypherQuery is a statement in the Cypher query language, with optional
// parameters and result.  If Result value is supplied, result data will be
// unmarshalled into it when the query is executed. Result must be a pointer
// to a slice of slices of structs - e.g. &[][]someStruct{}.
type CypherQuery struct {
	Statement  string                 `json:"statement"`
	Parameters map[string]interface{} `json:"parameters"`
	Result     interface{}            `json:"-"`
	columns    []string
	data       [][]*json.RawMessage
}

// Columns returns the names, in order, of the columns returned for this query.
// Empty if query has not been executed.
func (cq *CypherQuery) Columns() []string {
	return cq.columns
}

// Unmarshall decodes result data into v, which must be a pointer to a slice of
// slices of structs - e.g. &[][]someStruct{}.  Struct fields are matched up
// with fields returned by the cypher query using the `json:"fieldName"` tag.
func (cq *CypherQuery) Unmarshall(v interface{}) error {
	// We do a round-trip thru the JSON marshaller.  A fairly simple way to
	// do type-safe unmarshalling, but perhaps not the most efficient solution.
	rs := make([]map[string]*json.RawMessage, len(cq.data))
	for rowNum, row := range cq.data {
		m := map[string]*json.RawMessage{}
		for colNum, col := range row {
			name := cq.columns[colNum]
			m[name] = col
		}
		rs[rowNum] = m
	}
	b, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		logPretty(err)
		return err
	}
	return json.Unmarshal(b, v)
}

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
	Commit  string
	Results []struct {
		Columns []string
		Data    [][]*json.RawMessage
	}
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
		r := tr.Results[i]
		s.columns = r.Columns
		s.data = r.Data
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
