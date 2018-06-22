package algo

import (
	"log"
	"testing"
)

type Matrix [][]int

func rotate(m Matrix, n int) {
	for k := 0; k < n/2; k++ {
		for j := k; j < n-1-k; j++ {
			tmp := m[k][j]
			m[k][j] = m[j][n-k-1]
			m[j][n-k-1] = m[n-k-1][n-j-1]
			m[n-k-1][n-j-1] = m[n-j-1][k]
			m[n-j-1][k] = tmp
		}
	}
}

func rotate2(m [][]int, n int) {
	// transpose
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			temp := m[i][j]
			m[i][j] = m[j][i]
			m[j][i] = temp
		}
	}

	// rolling over
	for i := 0; i < n/2; i++ {
		for j := 0; j < n; j++ {
			tmp := m[i][j]
			m[i][j] = m[n-1-i][j]
			m[n-1-i][j] = tmp
		}
	}
}

func _initMatrix() {
	dim := 3
	m := make([][]int, dim)
	for i := range m {
		m[i] = make([]int, dim)
	}
	log.Printf("%v\n", m)
}

func TestRotate(t *testing.T) {
	m := Matrix{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	rotate(m, 3)
	log.Printf("rotate ==> %v", m)

	m = [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	rotate2(m, 3)
	log.Printf("transpose ==> %v", m)
}

func TestDiv(t *testing.T) {
	log.Printf("%v", 3/2)
}
