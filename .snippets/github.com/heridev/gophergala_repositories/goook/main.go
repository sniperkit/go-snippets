package main

import (
	"fmt"
	"log"
)

func main() {
	ddc, err := LoadDDC("data/classes_clean")
	if err != nil {
		log.Fatal(err)
	}

	tree := NewTreeBayes()

	// FIXME: do some real learning here

	for id := range ddc.Classes {
		tree.Learn(id, []string{fmt.Sprintf("%d", id)})
	}

	pt := tree.ClassifyTree([]string{"100", "120", "123"})
	pt.Display(1, "", "    ", ddc)
}
