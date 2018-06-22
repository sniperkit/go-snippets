package giraffe

import (
	"strconv"
	"testing"
)

func BenchmarkGraphNewGraph(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = newTestGraph()
	}
}

func BenchmarkGraphInsertNode(b *testing.B) {
	g, _ := newTestGraph()
	for n := 0; n < b.N; n++ {
		g.InsertNode()
	}
}

func BenchmarkGraphDeleteNode(b *testing.B) {
	g, _ := NewGraph("benchGraph")
	for n := 0; n < b.N; n++ {
		g.InsertNode()
	}
	for n := 0; n < b.N; n++ {
		g.deleteNodeByID(uint64(n))
	}
}

func BenchmarkGraphInsertDataNode(b *testing.B) {
	g, _ := newTestGraph()
	for n := 0; n < b.N; n++ {
		g.InsertDataNode("some key", []byte("some data"))
	}
}

func BenchmarkDFS(b *testing.B) {
	g, _ := NewGraph("benchGraph")
	for n := 1; n < b.N; n++ {
		g.InsertNode()
		g.Nodes[uint64(n-1)].AddRelationship(g.Nodes[uint64(n)])
		g.Root().DepthFirstSearch(g.Nodes[uint64(n)])
	}
}

func BenchmarkBFS(b *testing.B) {
	g, _ := NewGraph("benchGraph")
	for n := 1; n < b.N; n++ {
		g.InsertNode()
		g.Nodes[uint64(n-1)].AddRelationship(g.Nodes[uint64(n)])
		g.Root().BreadthFirstSearch(g.Nodes[uint64(n)])
	}
}

func BenchmarkFindByKey(b *testing.B) {
	g, _ := NewGraph("benchGraph")
	for n := 1; n < b.N; n++ {
		key := strconv.Itoa(n)
		g.InsertDataNode(key, []byte("data"))
		g.Nodes[uint64(n-1)].AddRelationship(g.Nodes[uint64(n)])
		g.FindNodeByKey(key)
	}
}
