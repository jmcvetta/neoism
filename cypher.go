// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neo4j

import (
	"github.com/jmcvetta/restclient"
)

type cypherRequest struct {
	Query string `json:"query"`
}

type cypherRequestParams struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}

// A CypherResult is returned when a cypher query is executed.
type CypherResult struct {
	Columns *[]string   `json:"columns"`
	Data    interface{} `json:"data"`
}

// Cypher executes a db query written in the Cypher language.  Data returned
// from the db is used to populate `result`, which should be a pointer to a
// slice of structs.  TODO:  Or a pointer to a two-dimensional array of structs?
func (db *Database) Cypher(query string, params map[string]interface{}, result interface{}) (columns []string, err error) {
	columns = []string{}
	cr := CypherResult{
		Columns: &columns,
		Data:    result,
	}
	ne := new(neoError)
	var data interface{}
	if params != nil {
		data = cypherRequestParams{
			Query:  query,
			Params: params,
		}
	} else {
		data = cypherRequest{
			Query: query,
		}
	}
	req := restclient.RequestResponse{
		Url:    db.HrefCypher,
		Method: "POST",
		Data:   data,
		Result: &cr,
		Error:  ne,
	}
	status, err := db.rc.Do(&req)
	if err != nil {
		return columns, err
	}
	if status != 200 {
		logPretty(req)
		return columns, BadResponse
	}
	return columns, nil
}
