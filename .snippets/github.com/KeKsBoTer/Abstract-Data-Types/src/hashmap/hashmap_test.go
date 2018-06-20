package hashmap

import (
	"testing"
	"fmt"
)

var testStrings = []string{"123", "321", "dev", "test", "golang", "hash", "list", "map", "::::"}

func TestHashedMap_Insert(t *testing.T) {
	hMap := NewHashMap(3)
	for i, e := range testStrings {
		hMap.Insert(e, i*10)
	}
	fmt.Println(hMap)
}

func TestHashedMap_Delete(t *testing.T) {
	hMap := NewHashMap(3)
	for i, e := range testStrings {
		hMap.Insert(e, i*10)
	}
	fmt.Println(hMap)
	for _, e := range testStrings {
		hMap.Delete(e)
	}
	if size := hMap.Size(); size > 0 {
		t.Error(
			"For", hMap,
			"expected lenght", 0,
			"got", size,
		)
	}
}
