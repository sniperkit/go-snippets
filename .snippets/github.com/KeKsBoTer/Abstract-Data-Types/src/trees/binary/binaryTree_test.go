package binary

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
	for k,v := range testMap{
		tree.Insert(v,k)
	}
	fmt.Println(tree)
}
