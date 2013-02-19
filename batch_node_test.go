// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package neo4j

import (
"testing"
)

func TestBatchCreateNode(t *testing.T) {
	b := batchTest(t)
	batchNodes := make([]*BatchNode, 3)
	for i := 0; i < 3; i++ {
		bn := b.CreateNode(nil)
		batchNodes = append(batchNodes, bn)
	}
	result, err := b.Execute()
	if err != nil {
		t.Error(err)
	}
	for _, r := range result {
		// FIXME: r is an entity, which is not terribly useful to us.  
		bn, ok := r.(BatchNode)
		if !ok {
			t.Errorf("Type assertion failed")
			logPretty(r)
		}
		logPretty(bn)
		n, err := bn.Node()
		if err != nil {
			t.Error(err)
		}
		n.Delete()
	}
}