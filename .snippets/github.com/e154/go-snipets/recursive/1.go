package main

import (
	"fmt"
	"encoding/json"
)

type Item struct {
	Name			string		`json:"name"`
	Parent			string		`json:"parent"`
	Data			string		`json:"data"`
}

type Tree struct {
	Name			string		`json:"name"`
	Nodes			[]*Tree		`json:"nodes"`
}

func renderTreeRecursive(i []*Item, t *Tree, c string) {

	for _, item := range i {
		if item.Parent == c {
			tree := new(Tree)
			tree.Name = item.Name
			tree.Nodes = make([]*Tree, 0)	// fix - nodes: null
			t.Nodes = append(t.Nodes, tree)
			renderTreeRecursive(i, tree, item.Name)
		}
	}

	return
}

func main(){

	items := []*Item{
		&Item{"itema1", "", "some data"},
		&Item{"itema2", "itema1", ""},
		&Item{"itema3", "itema2", ""},
		&Item{"itema4", "itema2", ""},
		&Item{"itema5", "itema2", ""},
	}

	tree := &Tree{}
	renderTreeRecursive(items, tree, "")
	b, _ := json.MarshalIndent(tree, " ", "  ")


	fmt.Println(string(b))
}

