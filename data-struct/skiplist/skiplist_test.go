package skiplist

import (
	"log"
	"math/rand"
	"testing"
)

func init() {
	// modify random. Same seed generates same results.
	random = rand.New(rand.NewSource(1))
}

func TestSkipList(t *testing.T) {
	sl := NewSkipList()

	for i := 1; i < 19; i++ {
		sl.Insert(i, i)
	}
	sl.Display()
	if i, ok := sl.Search(4); ok {
		log.Printf("search 4 index: %d\n", i)
	}
	if ok := sl.Delete(4); ok {
		log.Println("delete 4 success")
	}
	sl.Display()
}

func TestSkipListInsert(t *testing.T) {
	sl := NewSkipList()
	for i := 1; i < 19; i++ {
		sl.Insert(i, i)
	}
	sl.Display()
	if value, ok := sl.Search(17); ok && value != 17 {
		t.Errorf("The value of Key 17 is not 17 but %d", value)
	}
}

func TestSkipListDelete(t *testing.T) {
	sl := NewSkipList()
	for i := 1; i < 19; i++ {
		sl.Insert(i, i)
	}
	sl.Display()
	if ok := sl.Delete(17); !ok {
		t.Errorf("Key 17 is not found")
	}
	sl.Display()
	if _, ok := sl.Search(17); ok {
		t.Errorf("The key 17 should not exists")
	}
}
