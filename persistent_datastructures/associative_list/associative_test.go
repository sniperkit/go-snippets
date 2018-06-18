package associative_list

import (
	"fmt"
	"testing"
)

func TestAssociative1(t *testing.T) {
	m1 := empty{}
	if m1.Value("xx") != nil {
		t.Fatalf("m1 got \"xxx\" is not equal nil")
	}
	m2 := m1.Set("a", 1)
	if m2.Value("a").(int) != 1 {
		t.Fatalf("m2 got \"a\" is not equal 1")
	}
	m3 := m2.Set("b", 2)
	if m3.Value("b").(int) != 2 {
		t.Fatalf("m3 got \"b\" is not equal 2")
	}
	m4 := m3.Set("a", 1)
	if m4.Value("a").(int) != 1 {
		t.Fatalf("m4 got \"a\" is not equal 1")
	}

	keys := []string{"a", "b", "a"}
	values := []int{1, 2, 1}
	index := 0
	validateFunc := func(k, v interface{}) {
		fmt.Printf("index: %d k: %v v: %v\n", index, k, v)
		if key, ok := k.(string); !ok || key != keys[index] {
			t.Fatalf("associative_list key: %s not equal right one %s", key, keys[index])
		}
		if value, ok := v.(int); !ok || value != values[index] {
			t.Fatalf("associative_list key: %s value: %d not equal right one %d", keys[index], value, values[index])
		}
		index++
	}
	m4.Iterate(validateFunc)

	fmt.Println("----------I'm a pretty line----------")

	m := Map{m4}
	keys = []string{"a", "b"}
	values = []int{1, 2}
	index = 0
	validateFunc = func(k, v interface{}) {
		fmt.Printf("index: %d k: %v v: %v\n", index, k, v)
		if key, ok := k.(string); !ok || key != keys[index] {
			t.Fatalf("associative_list key: %s not equal right one %s", key, keys[index])
		}
		if value, ok := v.(int); !ok || value != values[index] {
			t.Fatalf("associative_list key: %s value: %d not equal right one %d", keys[index], value, values[index])
		}
		index++
	}
	m.Iterate(validateFunc)
}
