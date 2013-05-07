// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
)

func (db *Database) CreateRelIndex(name, idxType, provider string) (*RelationshipIndex, error) {
	idx, err := db.createIndex(db.HrefRelIndex, name, idxType, provider)
	if err != nil {
		return nil, err
	}
	return &RelationshipIndex{*idx}, nil
}

func (db *Database) RelationshipIndexes() ([]*RelationshipIndex, error) {
	indexes, err := db.indexes(db.HrefRelIndex)
	if err != nil {
		return nil, err
	}
	ris := make([]*RelationshipIndex, len(indexes))
	for i, idx := range indexes {
		ris[i] = &RelationshipIndex{*idx}
	}
	return ris, nil
}

func (db *Database) RelationshipIndex(name string) (*RelationshipIndex, error) {
	idx, err := db.index(db.HrefRelIndex, name)
	if err != nil {
		return nil, err
	}
	return &RelationshipIndex{*idx}, nil
}

// A RelationshipIndex is an index for searching Relationships.
type RelationshipIndex struct {
	index
}

// Remove deletes all entries with a given node, key and value from the index.
// If value or both key and value are the blank string, they are ignored.
func (rix *RelationshipIndex) Remove(r *Relationship, key, value string) error {
	return rix.remove(r, key, value)
}

