package magicsquare_test

import (
	"fmt"
	"testing"

	"reflect"

	"github.com/sahilm/go-concurrency-patterns/magicsquare"
)

func TestPermute(t *testing.T) {
	cases := []struct {
		in   []int
		want [][]int
	}{
		{
			[]int{},
			[][]int{},
		},
		{
			[]int{1},
			[][]int{{1}},
		},
		{
			[]int{1, 2},
			[][]int{
				{1, 2},
				{2, 1},
			},
		},
		{
			[]int{1, 2, 3},
			[][]int{
				{1, 2, 3},
				{1, 3, 2},
				{2, 1, 3},
				{2, 3, 1},
				{3, 1, 2},
				{3, 2, 1},
			},
		},
		{
			[]int{1, 2, 3, 4},
			[][]int{
				{1, 2, 3, 4},
				{1, 2, 4, 3},
				{1, 3, 2, 4},
				{1, 3, 4, 2},
				{1, 4, 2, 3},
				{1, 4, 3, 2},
				{2, 1, 3, 4},
				{2, 1, 4, 3},
				{2, 3, 1, 4},
				{2, 3, 4, 1},
				{2, 4, 1, 3},
				{2, 4, 3, 1},
				{3, 1, 2, 4},
				{3, 1, 4, 2},
				{3, 2, 1, 4},
				{3, 2, 4, 1},
				{3, 4, 1, 2},
				{3, 4, 2, 1},
				{4, 1, 2, 3},
				{4, 1, 3, 2},
				{4, 2, 1, 3},
				{4, 2, 3, 1},
				{4, 3, 1, 2},
				{4, 3, 2, 1}},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("when %v", c.in), func(t *testing.T) {
			got := magicsquare.Permute(c.in...)
			eq := reflect.DeepEqual(got, c.want)
			if !eq {
				t.Errorf("got: %v, want: %v", got, c.want)
			}
		})
	}
}
