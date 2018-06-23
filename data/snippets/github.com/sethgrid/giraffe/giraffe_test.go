package giraffe

import "testing"

func TestRootNodeCreate(t *testing.T) {
	g, err := NewGraph("testGraph")
	if err != nil {
		t.Fatalf("unable to create graph, error: %v", err)
	}

	if g.NodeCount() != 1 {
		t.Error("expected a root node to be created")
	}
	if g.LastNodeID() != 0 {
		t.Errorf("got id %d, want id %d for root node", g.LastNodeID(), 0)
	}
}

func TestAddNewNodesWithDestinationsAndSources(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()
	n4 := g.InsertNode()

	g.Root().AddRelationship(n1)
	g.Root().AddRelationship(n2)
	n2.AddRelationship(n3)
	n4.AddRelationship(n3)

	destinations := extractIDs(g.Root().ListDestinations())

	if len(destinations) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(destinations), 2)
	}

	if want := []uint64{n1.ID, n2.ID}; !ContainsAll(destinations, want) {
		t.Errorf("actual destinations %v not expected %v", destinations, want)
	}

	sources := extractIDs(n3.ListSources())

	if len(sources) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(sources), 2)
	}

	if want := []uint64{n2.ID, n4.ID}; !ContainsAll(sources, want) {
		t.Errorf("actual sources %v not expected %v", sources, want)
	}

}

func TestFindRoots(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()
	n4 := g.InsertNode()

	g.Root().AddRelationship(n1)
	g.Root().AddRelationship(n2)
	n2.AddRelationship(n3)

	roots := g.FindRoots()
	expectedRoots := []uint64{g.Root().ID, n4.ID}

	if len(roots) != len(expectedRoots) {
		t.Errorf("got %d roots, wat %d roots", len(roots), len(expectedRoots))
	}

	if !ContainsAll(roots, expectedRoots) {
		t.Errorf("actual roots %v not expected %v", roots, expectedRoots)
	}
}

func TestAddNewDataNodesAndFindNode(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1, _ := g.InsertDataNode("key1", []byte("value1"))
	n2, _ := g.InsertDataNode("key2", []byte("value2"))
	n3, _ := g.InsertDataNode("key3", []byte("value3"))

	g.Root().AddRelationship(n1)
	g.Root().AddRelationship(n2)
	n2.AddRelationship(n3)

	destinations := extractIDs(g.Root().ListDestinations())

	if len(destinations) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(destinations), 2)
	}

	if want := []uint64{n1.ID, n2.ID}; !ContainsAll(destinations, want) {
		t.Errorf("actual destinations %v not expected %v", destinations, want)
	}

	foundNode, ok := g.FindNodeByKey("key3")
	if !ok {
		t.Fatal("unable to find node")
	}

	if string(foundNode.Value) != "value3" {
		t.Errorf("got `%s`, want `value3` data", foundNode.Value)
	}
}

func TestDuplicateKeyError(t *testing.T) {
	g, _ := NewConstraintGraph("testGraph", false, true)

	_, err := g.InsertDataNode("key", []byte("data"))
	if err != nil {
		t.Errorf("expected no error, got `%s`", err)
	}

	_, err = g.InsertDataNode("key", []byte("oops, same key"))
	if err == nil {
		t.Error("no error when there should have been")
	}

	if err.Error() != ErrKeyExists {
		t.Errorf("unexpected error message. got `%s`, want `%s`", err.Error(), ErrKeyExists)
	}
}

func TestCircularRelationship(t *testing.T) {
	g, _ := NewConstraintGraph("testGraph", false, false)
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()

	var err error

	err = g.Root().AddRelationship(n1)
	if err != nil {
		t.Error("should not error adding relationship")
	}
	err = n1.AddRelationship(n2)
	if err != nil {
		t.Error("should not error adding relationship")
	}
	err = n2.AddRelationship(n3)
	if err != nil {
		t.Error("should not error adding relationship")
	}
	err = n3.AddRelationship(g.Root())
	if err == nil {
		t.Error("should error adding circular relationship")
	}
}

func TestSearch(t *testing.T) {
	g, _ := newTestGraph()

	if !g.Root().DepthFirstSearch(g.Nodes[11]) {
		t.Error("unable to find path root -> n11")
	}

	if g.Root().DepthFirstSearch(g.Nodes[12]) {
		t.Error("should not find a path root -> n12")
	}

	if !g.Root().BreadthFirstSearch(g.Nodes[9]) {
		t.Error("unable to find path root -> n9")
	}

	if g.Root().BreadthFirstSearch(g.Nodes[12]) {
		t.Error("should not find a path root -> n12")
	}
}

func TestRemoveRelationship(t *testing.T) {
	g, _ := newTestGraph()

	// cause nodes 9 and 10 (and 11) to be cut off (because 2 does not touch 6)
	err := g.Nodes[2].RemoveRelationship(g.Nodes[6])
	if err != nil {
		t.Fatalf("RemoveRelationship should not error, got `%v`", err)
	}

	for i := uint64(9); i <= 11; i++ {
		if g.Nodes[2].DepthFirstSearch(g.Nodes[i]) {
			t.Errorf("found node %d but should not have", i)
		}
	}
}

func TestDeleteNode(t *testing.T) {
	g, _ := newTestGraph()

	// cause nodes 9 and 10 (and 11) to be cut off (because 6 is gone)
	err := g.DeleteNode(g.Nodes[6])
	if err != nil {
		t.Fatalf("RemoveRelationship should not error, got `%v`", err)
	}

	for i := uint64(9); i <= 11; i++ {
		if g.Nodes[2].DepthFirstSearch(g.Nodes[i]) {
			t.Errorf("found node %d but should not have", i)
		}
	}

	if _, ok := g.Nodes[6]; ok {
		t.Error("node 6 found but should be deleted")
	}
}

func TestDeleteNodeByID(t *testing.T) {
	g, _ := newTestGraph()

	// cause nodes 9 and 10 (and 11) to be cut off (because 6 is gone)
	err := g.DeleteNodeByID(6)
	if err != nil {
		t.Fatalf("RemoveRelationship should not error, got `%v`", err)
	}

	for i := uint64(9); i <= 11; i++ {
		if g.Nodes[2].DepthFirstSearch(g.Nodes[i]) {
			t.Errorf("found node %d but should not have", i)
		}
	}

	if _, ok := g.Nodes[6]; ok {
		t.Error("node 6 found but should be deleted")
	}
}

func newTestGraph() (*Graph, error) {
	g, err := NewGraph("testGraph")
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()
	n4 := g.InsertNode()
	n5 := g.InsertNode()
	n6 := g.InsertNode()
	n7 := g.InsertNode()
	n8 := g.InsertNode()
	n9 := g.InsertNode()
	n10 := g.InsertNode()
	n11 := g.InsertNode()
	n12 := g.insertNode() // stranded
	_ = n12               // do nothing with it

	g.Root().AddRelationship(n1)
	g.Root().AddRelationship(n2)
	g.Root().AddRelationship(n3)
	n1.AddRelationship(n4)
	n2.AddRelationship(n5)
	n2.AddRelationship(n6)
	n3.AddRelationship(n7)
	n5.AddRelationship(n8)
	n6.AddRelationship(n9)
	n6.AddRelationship(n10)
	n10.AddRelationship(n11)

	/*
		above makes the following tree
		             0        12 (stranded)
		           / | \
		          1  2  3
		         /  / \  \
		        4  5   6  7
		          /   / \
		         8   9   10
		                  \
		                   11
	*/
	return g, err
}

func ContainsAll(superset, subset []uint64) bool {
	for _, sub := range subset {
		found := false
		for _, super := range superset {
			if sub == super {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
