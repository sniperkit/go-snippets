// Package pagerank computes PageRank on up to MaxUInt32 (4,294,967,295) items via graph traversal.
//
// The demonstration that this method is effective in calculating PageRank
// comes via http://arxiv.org/pdf/1006.2880.pdf
//
// The key insight is that simple graph iteration and counting can approximate PageRank,
// and can also -- if the graph is stored -- be used to inexpensively recalculate PageRank
// when edges are added to or removed from the graph.
//
// This package does not ultimately store the entire graph traversal history, and therefore
// it does not assist with recalculating PageRank when edges are removed from the graph.
// However, it is designed to make the addition of edges easy and inexpensive, since it stores
// traversal counts from prior runs. The choice to store traversal counts (rather than the full
// traversal history) is the reason that edge removal is not facilitated.
package pagerank

import (
	"fmt"
	"math/rand"
)

type Node interface {
	IsStarter() bool // Should we start iterations from this node?
	Traverse()
	Traversals() uint64
}

type Graph struct {
	// Provided by the user in the constructor
	nodes        *[]Node
	edgesFor     func(Node) []Node
	randomSource *rand.Rand

	// Caches and other unexported values
	traversals      uint64
	jumpProbability float32
	calculated      bool
	summed          bool
}

func NewGraph(seed int64, edgesFor func(Node) []Node, nodes *[]Node) *Graph {
	return &Graph{
		randomSource: rand.New(rand.NewSource(seed)),
		edgesFor:     edgesFor,
		nodes:        nodes,
	}
}

func (g *Graph) Pagerank(node Node, normalized bool) (float32, error) {
	if !g.calculated {
		return 0, fmt.Errorf("Pagerank graph has not yet been calculated")
	}

	if !g.summed {
		for _, node := range *g.nodes {
			g.traversals += node.Traversals()
		}
		g.summed = true
	}

	if normalized {
		return float32(float64(node.Traversals()) / float64(g.traversals)), nil
	}

	return float32(float64(len(*g.nodes)) * float64(node.Traversals()) / float64(g.traversals)), nil
}

func (g *Graph) traverseFrom(node Node) {
	node.Traverse()

	// Terminate the traversal with probability 1/g.jumpProbability
	if x := g.randomSource.Float32(); x < g.jumpProbability {
		return
	}

	outlinks := g.edgesFor(node)

	// Terminate the traversal if the node has no outgoing links
	if len(outlinks) < 1 {
		return
	}

	// Continue the traversal from a randomly chosen outgoing link
	g.traverseFrom(outlinks[g.randomSource.Intn(len(outlinks))])
}

// Calculate runs the Pagerank computation on your graph in-place.
func (g *Graph) Calculate(JumpProbability float32, RoundsPerNode int) {
	g.jumpProbability = JumpProbability

	for _, node := range *g.nodes {
		// Only start at starter nodes
		if !node.IsStarter() {
			continue
		}

		for round := 0; round < RoundsPerNode; round++ {
			g.traverseFrom(node)
		}
	}

	g.calculated = true
}
