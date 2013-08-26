neoism - Neo4j client for Go
===========================

![Neo4j + Gopher Logo](https://raw.github.com/jmcvetta/neoism/master/neoism.png)

Package `neoism` is a [Go](http://golang.org) client library providing access to
the [Neo4j](http://www.neo4j.org) graph database via its REST API.


# Requirements

[Go 1.1](http://golang.org/doc/go1.1) or later is required.

Neo4j Milestone 2.0.0-M03 or greater is required to run the full test suite.


# Documentation

See [Go Walker](http://gowalker.org/github.com/jmcvetta/neoism) or
[GoDoc](http://godoc.org/github.com/jmcvetta/neoism) for automatic
documentation.


# Status

[![Build Status](https://travis-ci.org/jmcvetta/neoism.png?branch=master)](https://travis-ci.org/jmcvetta/neoism)
[![Build Status](https://drone.io/github.com/jmcvetta/neoism/status.png)](https://drone.io/github.com/jmcvetta/neoism/latest)
[![Coverage Status](https://coveralls.io/repos/jmcvetta/neoism/badge.png?branch=master)](https://coveralls.io/r/jmcvetta/neoism)

This driver is fairly complete, and may now be suitable for general use.  The
code has an extensive set of integration tests, but little real-world testing.
YMMV; use in production at your own risk.

## Production Note

If you decide to use `neoism` in a production system, please let me know.  All
API changes will be made via Pull Request, so it's highly recommended you Watch
the repo Issues.  The API is fairly stable, but there are additions and small
changes from time to time.


## Completed:

* Node (create/edit/relate/delete/properties)
* Relationship (create/edit/delete/properties)
* Legacy Indexing (create/edit/delete/add node/remove node/find/query)
* Cypher queries
* Batched Cypher queries
* Transactional endpoint (Neo4j 2.0)
* Node labels (Neo4j 2.0)
* Schema index (Neo4j 2.0)


## To Do:

* ~~Unique Indexes~~ - probably will not expand support for legacy indexing.
* ~~Automatic Indexes~~ - "
* Traversals - May never be supported due to security concerns.  From the
  manual:  "The Traversal REST Endpoint executes arbitrary Groovy code under
  the hood as part of the evaluators definitions. In hosted and open
  environments, this can constitute a security risk."
* Built-In Graph Algorithms
* Gremlin


# Support

Paid support, development, related professional services, and proprietary
licensing terms for this package are available from [from the
author](mailto:jason.mcvetta@gmail.com).


# License

This is Free Software, released under the terms of the [GPL
v3](http://www.gnu.org/copyleft/gpl.html).
