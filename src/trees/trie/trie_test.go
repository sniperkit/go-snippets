package trie

import (
	"testing"
	"fmt"
)

var testMap = map[string]int{
	"simon":   50,
	"daniel":  51,
	"luci":    40,
	"caro":    41,
	"tom":     30,
	"lucifer": 70,
	"pi":      20,
	"fritz":   53,
}

func TestTree_Insert(t *testing.T) {
	tree := NewTree()
	for k, v := range testMap {
		tree.Insert(k, v)
	}

	//insert all values again
	for k, v := range testMap {
		tree.Insert(k, v)
	}

	//check if all items are in the trie
	if treeLen, mapLen := tree.Size(), len(testMap); treeLen != mapLen {
		t.Error(
			"Checked length for ", tree,
			"expected", mapLen,
			"got", treeLen,
		)
		return
	}

	//check if all values where inserted correct
	for k := range testMap {
		if !tree.Contains(k) {
			t.Error(
				"Checked tree", tree,
				"for key", k,
				"received", false,
			)
			return
		}
	}

	fmt.Println("Success", tree)
}

func TestTree_Delete(t *testing.T) {
	tree := NewTree()
	for k, v := range testMap {
		tree.Insert(k, v)
	}
	for k := range testMap {
		tree.Delete(k)
	}
	if treeLen := tree.Size(); treeLen != 0{
		t.Error(
			"Deleted all keys in tree", tree,
			"expected length", 0,
			"received", treeLen,
		)
		return
	}
}
