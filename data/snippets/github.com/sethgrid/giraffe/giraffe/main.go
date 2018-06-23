package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/sethgrid/giraffe"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	//randomGraph()
	curriculumGraph()
}

func curriculumGraph() {
	g, _ := giraffe.NewConstraintGraph("math curriculum", true, true)

	root := g.Root() // same as g.Nodes[0]

	// this will be the entry point into the curriculum
	// Keys will be lesson titles, Values will be a foreign key into a relational DB (or something)
	root.Key = "Intro"
	root.Value = []byte("lesson_id 1")

	// this will be the culmination of this curriculum
	QF, _ := g.InsertDataNode("Quadratic Formula", []byte("lesson_id 55"))

	// All the different lessons that lead up to being able to understand the Quadratic Formula
	FQP, _ := g.InsertDataNode("Factoring Quadratic Polynomials", []byte("lesson_id 52"))
	FAE, _ := g.InsertDataNode("Factoring Algebraic Expressions", []byte("lesson_id 45"))
	FN, _ := g.InsertDataNode("Factoring Numerics", []byte("lesson_id 12"))
	GQ, _ := g.InsertDataNode("Graphing Quadratics", []byte("lesson_id 32"))
	GL, _ := g.InsertDataNode("Graphing Linear Equations", []byte("lesson_id 19"))
	SLE, _ := g.InsertDataNode("Solving Linear Equations", []byte("lesson_id 17"))
	ZPP, _ := g.InsertDataNode("Zero Product Property", []byte("lesson_id 27"))
	POLY, _ := g.InsertDataNode("Polynomials", []byte("lesson_id 9"))
	SUB, _ := g.InsertDataNode("Substitution", []byte("lesson_id 5"))

	// mapping relationships
	FN.AddRelationship(FAE)
	FAE.AddRelationship(FQP)
	GL.AddRelationship(GQ)
	GQ.AddRelationship(FQP)
	FQP.AddRelationship(QF)
	ZPP.AddRelationship(FQP)
	SUB.AddRelationship(SLE)
	SLE.AddRelationship(POLY)
	SLE.AddRelationship(ZPP)
	SLE.AddRelationship(GL)
	root.AddRelationship(FN)
	root.AddRelationship(POLY)
	root.AddRelationship(SUB)

	// save
	data, _ := g.GobEncode()
	ioutil.WriteFile("someFile", data, 0766)

	// load
	restoredData, _ := ioutil.ReadFile("someFile")
	restoredGraph := &giraffe.Graph{}
	restoredGraph.GobDecode(restoredData)

	// visualize
	fmt.Println(restoredGraph.ToVisJS(false, true, false))
}

func randomGraph() {
	// make a new graph
	g, _ := giraffe.NewGraph("example")

	// we automatically have access to the root node
	root := g.Root() // same as g.Nodes[0]

	// let's make a bunch of new nodes
	nodeCount := 50
	nodes := make([]*giraffe.Node, nodeCount)
	nodes[0] = root
	for i := 1; i < nodeCount; i++ {
		nodes[i], _ = g.InsertDataNode(fmt.Sprintf("node %d", i), []byte{})
	}

	// let's make random relationships between the nodes
	for i := 0; i < nodeCount*1; i++ {
		nodeA := int(rand.Intn(nodeCount - 1))
		nodeB := int(rand.Intn(nodeCount - 1))
		if nodeA == nodeB {
			nodeA = 0
		}

		nodes[nodeA].AddRelationship(nodes[nodeB])
	}

	// and we can now visualize the relationships
	fmt.Println(g.ToVisJS(true, false, false))

	// while we started with one root node, we can have many roots
	log.Println("root nodes: ", g.FindRoots())
}
