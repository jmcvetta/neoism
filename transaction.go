// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
	"time"
)

type CypherStatement struct {
	Statement string                 `json:"statement"`
	Params    map[string]interface{} `json:"params"`
	// Columns and Data are populated with the result from the server.  Data
	// is a struct into which the query result will be unmarshalled.
	Columns []string
	Data    interface{}
}

type txRequest struct {
	Statements []*CypherStatement `json:"statements"`
}

type txInfo struct {
	Expires time.Time `json:"expires"`
}

type txResult struct {
	Columns *[]string   `json:"columns"`
	Data    interface{} `json:"data"`
}

type txResponse struct {
	Commit  string     `json:"commit"`
	Results []txResult `json:"results"`
}

func (db *Database) BeginTx(stmts []*CypherStatement) (*Transaction, error) {
	ne := new(neoError)
	payload := txRequest{Statements: stmts}
	txres := make([]txResult, len(stmts))
	for i, s := range stmts {
		txres[i] = txResult{
			Columns: &s.Columns,
			Data:    &s.Data,
		}
	}
	res := txResponse{Results: txres}
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
	tx := Transaction{
		Location: rr.HttpResponse.Header.Get("location"),
		Commit:   res.Commit,
	}
	return &tx, nil
}

type Transaction struct {
	Location string
	Commit   string
}
