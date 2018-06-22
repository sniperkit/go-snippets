package giraffe

import "testing"

func TestEncodeDecode(t *testing.T) {
	g, _ := newTestGraph()

	data, err := g.GobEncode()
	if err != nil {
		t.Fatalf("unable to encode graph - err `%v`", err)
	}

	decodedGraph := &Graph{}
	err = decodedGraph.GobDecode(data)
	if err != nil {
		t.Fatalf("unable to decode graph - err `%v`", err)
	}

	var targetNode *Node
	var ok bool
	if targetNode, ok = decodedGraph.Nodes[10]; !ok {
		t.Fatal("desired target node does not exist")
	}

	if !decodedGraph.Root().DepthFirstSearch(targetNode) {
		t.Fatalf("unable to find node %d in graph starting at root", targetNode.ID)
	}

	// verify internal structure is restored
	if decodedGraph.circularRelationship != true {
		t.Error("newTestGraph should allow circular relationships, decodedGraph too")
	}

	if len(g.Nodes[10].sources) != 1 {
		t.Errorf("node source restoration failed. Len of node 10 sources: %d", len(g.Nodes[10].sources))
	}

	if g.Nodes[10].sources[0].ID != uint64(6) {
		t.Errorf("node source for restoration failed. node 10 source node not 6, got %d", g.Nodes[10].sources[0].ID)
	}
}
