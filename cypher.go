// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
	"github.com/jmcvetta/restclient"
)

type cypherRequest struct {
	Query  string            `json:"query"`
	Params map[string]string `json:"params"`
}

// A CypherResult is returned when a cypher query is executed.
type CypherResult struct {
	Columns []string   `json:"columns"`
	Data    [][]string `json:"data"`
}

func (db *Database) CypherQuery(query string, params map[string]string) (*CypherResult, error) {
	result := new(CypherResult)
	ne := new(neoError)
	data := cypherRequest{
		Query:  query,
		Params: params,
	}
	req := restclient.RestRequest{
		Url:    db.info.Cypher,
		Method: restclient.POST,
		Data:   data,
		Result: result,
		Error:  ne,
	}
	status, err := db.rc.Do(&req)
	if err != nil {
		return result, err
	}
	if status != 200 {
		logPretty(req)
		return result, BadResponse
	}
	return result, nil
}
