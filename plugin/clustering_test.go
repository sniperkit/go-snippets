package plugin

import (
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestSelectItemsAMAP(t *testing.T) {
	goassert.New(t, "W must be non-negative").ExpectError(SelectItemsAMAP([]float32{}, -5))
	goassert.New(t, float32(0.0), []int{}).EqualWithoutError(SelectItemsAMAP([]float32{}, 0))
	goassert.New(t, "weights must have non-negatives").ExpectError(SelectItemsAMAP([]float32{1, 2, 5, -3}, 7))
	goassert.New(t, float32(7.0), []int{1, 2}).EqualWithoutError(SelectItemsAMAP([]float32{1, 2, 5, 3}, 7))
	goassert.New(t, float32(7.0), []int{1, 3}).EqualWithoutError(SelectItemsAMAP([]float32{1, 2, 3, 5}, 7))
	goassert.New(t, float32(11.0), []int{0, 1, 2, 3}).EqualWithoutError(SelectItemsAMAP([]float32{1, 2, 3, 5}, 12))
	goassert.New(t, float32(0.0), []int{}).EqualWithoutError(SelectItemsAMAP([]float32{3, 5, 7, 9}, 2))
}

func TestBipartitionWeightedGraph(t *testing.T) {
	goassert.New(t, []bool{}).EqualWithoutError(BipartitionWeightedGraph(0, 2, []float32{}))
	goassert.New(t, "A is not legal adjacency matrix").ExpectError(BipartitionWeightedGraph(4, 2, []float32{}))
	goassert.New(t, "A is not legal adjacency matrix").ExpectError(BipartitionWeightedGraph(2, 0, []float32{-1.0, -1.0, -1.0, -1.0}))
	// {0, 3}, {1, 2}
	assertBinaryClassAssignmentEqual(t, []bool{true, false, false, true}, goassert.New(t).SucceedNew(BipartitionWeightedGraph(4, 2, []float32{
		0.0, 0.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 0.0,
	})).([]bool))
	// {0, 3}, {1, 2} (connected)
	assertBinaryClassAssignmentEqual(t, []bool{false, true, true, false}, goassert.New(t).SucceedNew(BipartitionWeightedGraph(4, 2, []float32{
		0.0, 0.0, 1.0, 9.0,
		0.0, 0.0, 5.0, 0.0,
		1.0, 5.0, 0.0, 0.0,
		9.0, 0.0, 0.0, 0.0,
	})).([]bool))
}
