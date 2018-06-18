package bubble_test

import (
	"testing"
	"github.com/ygto/sort-algorithms/bubble"
	"math/rand"
	"fmt"
)

func TestBestCase(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := bubble.Sort(arr)
	testResult := testArray(result)
	if testResult != nil {
		t.Error(testResult)
	}
}
func TestWorstCase(t *testing.T) {
	arr := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	result := bubble.Sort(arr)
	testResult := testArray(result)
	if testResult != nil {
		t.Error(testResult)
	}
}
func TestRandomCase(t *testing.T) {
	var arr []int
	arraySize := 10
	for i := 0; i < arraySize; i++ {
		arr = append(arr, (rand.Int() % 100))
	}
	result := bubble.Sort(arr)
	testResult := testArray(result)
	if testResult != nil {
		t.Error(testResult)
	}
}
func TestRandomBiggerNumberCase(t *testing.T) {
	var arr []int
	arraySize := 10000
	for i := 0; i < arraySize; i++ {
		arr = append(arr, (rand.Int() % 100))
	}
	result := bubble.Sort(arr)
	testResult := testArray(result)
	if testResult != nil {
		t.Error(testResult)
	}
}
func testArray(arr []int, ) error {
	for i := 0; i < len(arr) - 1; i++ {
		for j := i; j < len(arr) - 1; j++ {
			if arr[i] > arr[j] {
				return fmt.Errorf("index[%d] :%d > index[%d] : %d", i, arr[i], j, arr[j])
			}
		}
	}
	return nil
}