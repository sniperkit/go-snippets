package pagerank

import (
	"runtime"
	"sync"
	"testing"
)

// Satisfies the pagerank.Node interface
type nodeExample struct {
	Base
	name string
	id   uint64
}

type adjacency struct {
	Left  string
	Right string
}

var wikiStructure []adjacency = []adjacency{
	{"b", "c"},
	{"c", "b"},
	{"d", "a"},
	{"d", "b"},
	{"e", "d"},
	{"e", "b"},
	{"e", "f"},
	{"f", "b"},
	{"f", "e"},
	{"v", "e"},
	{"v", "b"},
	{"w", "e"},
	{"w", "b"},
	{"x", "e"},
	{"x", "b"},
	{"y", "e"},
	{"z", "e"},
}

func TestPagerank(t *testing.T) {
	g := buildNewGraph(wikiStructure)

	// Compute
	g.Calculate(0.15, 500)

	for _, nodeMasked := range *g.nodes {
		node := nodeMasked.(*nodeExample)
		pr, err := g.Pagerank(node, true)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%s (%d): %d traversals, PR %.4f", node.name, node.id, node.Traversals(), pr)
	}
	t.Logf("%+v", g)
}

func TestPagerankParallel(t *testing.T) {
	graphs := make([]*Graph, runtime.NumCPU(), runtime.NumCPU())
	for i := range graphs {
		graphs[i] = buildNewGraph(wikiStructure)
	}

	// Compute
	wg := sync.WaitGroup{}
	for _, g := range graphs {
		wg.Add(1)
		go func(g *Graph) {
			g.Calculate(0.15, 100000)
			wg.Done()
		}(g)
	}
	wg.Wait()

	graph := graphs[0]

	for _, nodeMasked := range *graph.nodes {
		node := nodeMasked.(*nodeExample)
		pr, err := graph.Pagerank(node, true)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%s (%d): %d traversals, PR %.4f", node.name, node.id, node.Traversals(), pr)
	}
	t.Logf("%+v", graph)
}

func buildNewGraph(fromSource []adjacency) *Graph {
	// Map our nodes
	nodes := make([]Node, 0)
	nodeMap := make(map[string]int)
	edges := make(map[Node][]Node)

	letters := make([]string, 2, 2)
	var i uint64
	for _, link := range fromSource {
		letters[0] = link.Left
		letters[1] = link.Right
		for _, letter := range letters {
			if _, exists := nodeMap[letter]; exists {
				continue
			}
			nodes = append(nodes, &nodeExample{id: i, name: letter})
			nodeMap[letter] = len(nodes) - 1
			i++
		}
		//append(nodes[nodeMap[link.Left]].)
		//g.AddEdge(nodes[nodeMap[link.Left]], nodes[nodeMap[link.Right]])
		edges[nodes[nodeMap[letters[0]]]] = append(edges[nodes[nodeMap[letters[0]]]], nodes[nodeMap[letters[1]]])
	}

	getEdges := func(n Node) []Node {
		return edges[n]
	}

	g := NewGraph(31337, getEdges, &nodes)

	return g
}
