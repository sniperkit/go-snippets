package pagerank

import (
	"testing"
)

func TestPageRank(t *testing.T) {
	adj_matrix := [][]float64{
		[]float64{0, 1, 0, 0, 0},
		[]float64{0, 0, 1, 0, 1},
		[]float64{0, 0, 0, 1, 0},
		[]float64{1, 1, 0, 0, 0},
		[]float64{0, 0, 0, 1, 0},
	}
	pg := NewPageRank(0.85)
	pg.GetVector(adj_matrix)
}
