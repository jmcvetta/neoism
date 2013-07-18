neo4j
=====

![Neo4j + Gopher Logo](https://raw.github.com/jmcvetta/neo4j/master/neo4j_gopher.png)
Neo4j client for [Go](http://golang.org).


# Requirements

Package `neo4j` requires [Go 1.1](http://golang.org/doc/go1.1).  Earlier
versions of Go cannot be used, as they can't unmarshall JSON into an embedded
struct.


# Documentation

See [Go Walker](http://gowalker.org/github.com/jmcvetta/neo4j) or
[GoDoc](http://godoc.org/github.com/jmcvetta/neo4j) for automatic
documentation.



# Status

[![Build Status](https://travis-ci.org/jmcvetta/neo4j.png?branch=master)](https://travis-ci.org/jmcvetta/neo4j)
[![Build Status](https://drone.io/github.com/jmcvetta/neo4j/status.png)](https://drone.io/github.com/jmcvetta/neo4j/latest)
[![Coverage Status](https://coveralls.io/repos/jmcvetta/neo4j/badge.png?branch=master)](https://coveralls.io/r/jmcvetta/neo4j)

This driver is a work in progress.  It is not yet complete, but may now be
suitable for use by others.  The code has an extensive set of integration
tests, but very little real-world testing.  YMMV; use in production at your own
risk.

## Production Note

If you decide to use `neo4j` in a production system, please let me know.  All
API changes will be made via Pull Request, so it's highly recommended you Watch
the repo Issues.  The API is **not** promised to be stable at this time.


## Completed:

* Node (create/edit/relate/delete/properties)
* Relationship (create/edit/delete/properties)
* Index (create/edit/delete/add node/remove node/find/query)
* Cypher (query with and without parameters) - still under active development,
  API should not be considered stable.

## To Do:

* Unique Indexes
* Automatic Indexes - Not sure how much there is to do here, but these are a
  seperate section in the REST API manual, that I have not yet read.
* Traversals - May never be supported due to security concerns.  From the
  manual:  "The Traversal REST Endpoint executes arbitrary Groovy code under
  the hood as part of the evaluators definitions. In hosted and open
  environments, this can constitute a security risk."
* Built-In Graph Algorithms
* Batch Operations
* Gremlin


# License

This is Free Software, released under the terms of the [GPL
v3](http://www.gnu.org/copyleft/gpl.html).

