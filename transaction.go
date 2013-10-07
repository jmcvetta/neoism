// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"encoding/json"
	"errors"
	"reflect"
	"os"
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
	Commit  string
	Results []struct {
		Columns []string
		// Data    [][]*json.RawMessage
		Data []struct{
	        Row []*json.RawMessage
		}
	}
	Transaction struct {
		Expires string
	}
	Errors []TxError
}

// unmarshal populates a slice of CypherQuery object with result data returned
// from the server.
func (tr *txResponse) unmarshal(qs []*CypherQuery) error {
	if len(tr.Results) != len(qs) {
		return errors.New("Result count does not match query count")
	}
	for resultIdx, res := range tr.Results {
		if len(res.Data) != 1 {
			return errors.New("Confused by server response - please file an Issue with as much detail as you can provide!")
		}
		cq := qs[resultIdx]
		resType := reflect.TypeOf(cq.Result).Elem().Elem()
		// resSliceType := reflect.SliceOf(resType)
		r := reflect.New(resType).Elem()
		// s := reflect.New(resSliceType).Elem()
		//
		// Parse struct tags
		//
		fieldNameToNum := make(map[string]int)
		for n := 0; n < r.NumField(); n++ {
			tag := resType.Field(n).Tag.Get("neoism")
			if tag == "" {
				continue
			}
			fieldNameToNum[tag] = n

		}
		//
		// Sanity check
		//
		columnToField := make(map[int]int)
		for colNum, colName := range res.Columns {
			fieldNum, ok := fieldNameToNum[colName]
			if !ok {
				logPretty("Oh bugger, no field tagged " + colName)
			}
			columnToField[colNum] = fieldNum
		}
		logPretty(columnToField)
		rowMaps := make([]map[string]*json.RawMessage, len(res.Data))
		for rowNum, row := range res.Data {
			m := map[string]*json.RawMessage{}
			for colNum, jsonMsg := range row.Row {
				// name := res.Columns[colNum]
				// m[name] = col
				iface := reflect.New(resType).Interface()
				json.Unmarshal(jsonMsg, iface)
			}
			rowMaps[rowNum] = m
			// iface := reflect.New(resType).Interface()
		}
	}
	os.Exit(0)
	/*
	rs := make([]map[string]*json.RawMessage, len(cq.cr.Data))
	for rowNum, row := range cq.cr.Data {
		m := map[string]*json.RawMessage{}
		for colNum, col := range row {
			name := cq.cr.Columns[colNum]
			m[name] = col
		}
		rs[rowNum] = m
	}
	b, err := json.Marshal(rs)
	if err != nil {
		logPretty(err)
		return err
	}
	return json.Unmarshal(b, v)
	 */

	/*
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
	*/
	return nil
}

// Begin opens a new transaction, executing zero or more cypher queries
// inside the transaction.
func (db *Database) Begin(qs []*CypherQuery) (*Tx, error) {
	payload := txRequest{Statements: qs}
	result := txResponse{}
	ne := NeoError{}
	resp, err := db.Session.Post(db.HrefTransaction, payload, &result, &ne)
	if err != nil {
		return nil, err
	}
	if resp.Status() != 201 {
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
	ne := NeoError{}
	resp, err := t.db.Session.Post(t.hrefCommit, nil, nil, &ne)
	if err != nil {
		return err
	}
	if resp.Status() != 200 {
		return ne
	}
	return nil // Success
}

// Query executes statements in an open transaction.
func (t *Tx) Query(qs []*CypherQuery) error {
	payload := txRequest{Statements: qs}
	result := txResponse{}
	ne := NeoError{}
	resp, err := t.db.Session.Post(t.Location, payload, &result, &ne)
	if err != nil {
		return err
	}
	if resp.Status() == 404 {
		return NotFound
	}
	if resp.Status() != 200 {
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
	ne := NeoError{}
	resp, err := t.db.Session.Delete(t.Location, nil, &ne)
	if err != nil {
		return err
	}
	if resp.Status() != 200 {
		return ne
	}
	return nil // Success
}
