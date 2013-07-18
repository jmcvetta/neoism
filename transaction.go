// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	// "time"
	"encoding/json"
	"errors"
)

type CypherQuery struct {
	Statement  string                 `json:"statement"`
	Parameters map[string]interface{} `json:"parameters"`
	columns []string
	data []json.RawMessage
}

func (cq *CypherQuery) Columns() []string {
	return cq.columns
}

type txRequest struct {
	Statements []*CypherQuery `json:"statements"`
}

type txResponse struct {
	Commit  string
	Results []struct {
		Columns []string
		Data    []json.RawMessage
	}
	Transaction struct {
		Expires string
	}
	Errors []struct {
		Code    int
		Status  string
		Message string
	}
}

func (db *Database) BeginTx(stmts []*CypherQuery) (*Transaction, error) {
	ne := new(neoError)
	payload := txRequest{Statements: stmts}
	res := txResponse{}
	rr := restclient.RequestResponse{
		Url:            db.HrefTransaction,
		Method:         "POST",
		Data:           payload,
		Result:         &res,
		Error:          &ne,
		ExpectedStatus: 201,
	}
	db.rc.Log = true
	_, err := db.rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	logPretty(res)
	tx := Transaction{
		Location: rr.HttpResponse.Header.Get("location"),
		Commit:   res.Commit,
	}
	if len(res.Results) != len(stmts) {
		return nil, errors.New("WTF?")
	}
	for i, s := range stmts {
		r := res.Results[i]
		logPretty(r)
		s.columns = r.Columns
		s.data = r.Data
	}
	logPretty(stmts)
	return &tx, nil
}

type Transaction struct {
	Location string
	Commit   string
}
