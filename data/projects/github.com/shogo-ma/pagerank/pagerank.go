package pagerank

import (
	"math"
)

type PageRank struct {
	Alpha float64
}

func NewPageRank(alpha float64) *PageRank {
	pagerank := new(PageRank)
	pagerank.Alpha = alpha
	return pagerank
}

func (pg *PageRank) GetVector(adj_matrix [][]float64) []float64 {
	transition_matris := pg.nomalize(adj_matrix)
	return pg.PowerMethod(transition_matris)
}

func (pg *PageRank) PowerMethod(transition [][]float64) []float64 {
	pagerank := pg.initVector(len(transition))
	var loss float64
	for loss > 0.0001 {
		vector := make([]float64, len(pagerank))
		for i := 0; i < len(pagerank); i++ {
			sum := 0.0
			for j := 0; j < len(transition); j++ {
				sum += pg.Alpha*(pagerank[j]*transition[j][i]) + (1.0-pg.Alpha)/float64(len(transition))
			}
			vector[i] = sum
		}
		loss = computeLoss(pagerank, vector)
		pagerank = vector
	}
	return pagerank
}

func computeLoss(vector, new_vector []float64) float64 {
	var loss float64
	for i := 0; i < len(vector); i++ {
		loss += math.Abs(vector[i] - new_vector[i])
	}
	return loss
}

func (pg *PageRank) initVector(column int) []float64 {
	vector := make([]float64, column)
	for i := 0; i < column; i++ {
		vector[i] = 1.0 / float64(column)
	}
	return vector
}

func (pg *PageRank) nomalize(matrix [][]float64) [][]float64 {
	p := make([][]float64, len(matrix))
	for i := 0; i < len(matrix); i++ {
		p[i] = make([]float64, len(matrix))
	}
	for i := range matrix {
		sum := 0.0
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
		for j := 0; j < len(matrix[i]); j++ {
			p[i][j] = matrix[i][j] / sum
		}
	}
	return p
}
