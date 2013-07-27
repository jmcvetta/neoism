neo4j - Neo4j client for Go
===========================

![Neo4j + Gopher Logo](https://raw.github.com/jmcvetta/neo4j/master/neo4j_gopher.png)

Package `neo4j` is a [Go](http://golang.org) client library providing access to
the [Neo4j](http://www.neo4j.org) graph database via its REST API.


# Requirements

[Go 1.1](http://golang.org/doc/go1.1) or later is required.  Earlier versions
of Go cannot be used, as they can't unmarshall JSON into an embedded struct.


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
* Legacy Indexing (create/edit/delete/add node/remove node/find/query)
* Cypher (query with and without parameters) - still under active development,
  API should not be considered stable.

## In Progress:

* Transactional endpoint (Neo4j 2.0)
* Batch Cypher Queries - transactional endpoint does not provide an adequate
  substitute, as it has no means to reference result of previous statements.

## To Do:

* Node labels (Neo4j 2.0)
* Schema index (Neo4j 2.0)
* ~~Unique Indexes~~ - probably will not expand support for legacy indexing.
* ~~Automatic Indexes~~ - "
* Traversals - May never be supported due to security concerns.  From the
  manual:  "The Traversal REST Endpoint executes arbitrary Groovy code under
  the hood as part of the evaluators definitions. In hosted and open
  environments, this can constitute a security risk."
* Built-In Graph Algorithms
* Gremlin


# License

This is Free Software, released under the terms of the [GPL
v3](http://www.gnu.org/copyleft/gpl.html).

