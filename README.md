neo4j
=====

Neo4j client for [Go](http://golang.org).


# Documentation

See [GoDoc](http://godoc.org/github.com/jmcvetta/neo4j) for automatic
documentation.


# Status

This driver is a work in progress.  It is not yet complete, but may now be
suitable for use by others.  The code has an extensive set of integration
tests, but very little real-world testing.  YMMV; use in production at your own
risk.

## Completed:

* Node (create/edit/relate/delete/properties)
* Relationship (create/edit/delete/properties)
* Index (create/edit/delete/add node/remove node/find/query)

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
* Cypher
* Gremlin


# License

This is Free Software, released under the terms of the [GPL
v3](http://www.gnu.org/copyleft/gpl.html).  Resist intellectual serfdom - the
ownership of ideas is akin to slavery.  

