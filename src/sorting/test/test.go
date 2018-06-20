package test

import (
	"strings"
	"runtime"
	"reflect"
	"fmt"
	"math/rand"
	"testing"
)

const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_N = "\x1b[0m"

func Run(testCount int, sort func([]int), t *testing.T) {
	lastList := new([]int)
	lastGenerated := new([]int)
	name := strings.Split(runtime.FuncForPC(reflect.ValueOf(sort).Pointer()).Name(), ".")
	for i := 0; i < testCount; i++ {
		length := rand.Intn(50) + 10
		list := generateList(length)
		lastList = &list
		unSorted := make([]int, len(list))
		lastGenerated = &unSorted
		copy(unSorted, list)
		sort(list)
		if success := isSorted(list); !success || len(list) != len(unSorted) || !sameCount(list, unSorted) {
			t.Error(CLR_R,
				"Sorting with function", "\""+name[len(name)-1]+"\"", "failed!\n",
				"Generated:", unSorted,
				"\n    Result:", list, CLR_N,
			)
			return
		}

	}
	fmt.Printf("%s", CLR_G)
	fmt.Printf("Successfully tested function \"%s\" with %d random lists without errors %s\n", name[len(name)-1], testCount, CLR_N)
	fmt.Printf("\tExample: \n")
	fmt.Println("\t", *lastGenerated)
	fmt.Println("\t", *lastList)
}

func generateList(length int) []int {
	list := make([]int, length)
	for i := range list {
		list[i] = rand.Intn(3*length) - length/2
	}
	return list
}

func isSorted(list []int) bool {
	for i := 0; i < len(list)-1; i++ {
		if list[i] > list[i+1] {
			return false
		}
	}
	return true
}

func sameCount(list1, list2 []int) bool {
	for i := range list1 {
		if count(list1[i], list1) != count(list1[i], list2) {
			return false
		}
	}
	return true
}

func count(value int, list []int) int {
	count := 0
	for i := range list {
		if list[i] == value {
			count++
		}
	}
	return count
}
