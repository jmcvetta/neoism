// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	// "time"
	"encoding/json"
)

type CypherStatement struct {
	Statement  string                 `json:"statement"`
	Parameters map[string]interface{} `json:"parameters"`
	// Columns and Data are populated with the result from the server.  Data
	// is a struct into which the query result will be unmarshalled.
	Columns []string
	Data    interface{}
}

type txRequest struct {
	Statements []*CypherStatement `json:"statements"`
}

type txResponse struct {
	Commit  string `json:"commit"`
	Results []struct {
		Columns []string        `json:"columns"`
		Data    json.RawMessage `json:"data"`
	} `json:"results"`
	Transaction struct {
		Expires string `json:"expires"` // server returns unparseable timestamp
	} `json:"transaction"`
	Errors []struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (db *Database) BeginTx(stmts []*CypherStatement) (*Transaction, error) {
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
	logPretty(len(res.Results))
	logPretty(len(stmts))
	/*
		if len(res.Results) != len(stmts) {
			return nil, errors.New("WTF?")
		}
	*/
	return &tx, nil
}

type Transaction struct {
	Location string
	Commit   string
}
