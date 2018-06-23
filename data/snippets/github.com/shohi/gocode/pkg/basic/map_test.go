package basic

import (
	"log"
	"strconv"
	"testing"
)

type state int

const (
	open state = iota
	halfopen
	closed
)

func TestMapTraverse(t *testing.T) {
	m := make(map[string]bool)
	m["hello"] = true
	m["world"] = false

	for v := range m {
		log.Println(v)
	}
}

func TestMapForCapacity(t *testing.T) {
	m := make(map[string]int, 2)
	log.Println(len(m))

	for v := range m {
		log.Println(v)
	}
}

func TestMapNoInitialize(t *testing.T) {
	var m map[string]int

	log.Println(m)

	// map must be intialized before use
	m = make(map[string]int)
	for k := 0; k < 10; k++ {
		m[strconv.Itoa(k)] = k
	}

	for k, v := range m {
		log.Println(k, v)
	}
}

func TestEnum(t *testing.T) {
	log.Printf("%v, %T\n", open, open)
	log.Printf("%v, %T\n", halfopen, halfopen)
	log.Println(open == 0)
}

func TestMapForNonexistKey(t *testing.T) {
	m := make(map[string][]byte)

	log.Println(m["hello"] == nil)
}

func TestMapIteration(t *testing.T) {
	m := map[string]int{
		"hello": 0,
		"world": 1,
	}

	// 1. iterate over key
	for k := range m {
		log.Println(k)
	}

	// 2. iterate over entry(key/value pair)
	for k, v := range m {
		log.Println(k, v)
	}
}

func TestMapForMutable(t *testing.T) {
	initSlice := make([]string, 0, 8)
	initSlice = append(initSlice, "world")
	m := make(map[string][]string)
	m["hello"] = initSlice
	strSlice := m["hello"]
	strSlice = append(strSlice, "new world")

	log.Println(m["hello"])
}

func TestMapWithIntKey(t *testing.T) {
	m := make(map[int]string)
	m[1] = "hello"

	for k, v := range m {
		log.Printf("key: %d, value: %s", k, v)
	}
}
