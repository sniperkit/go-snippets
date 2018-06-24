package plugin

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/hiro4bbh/sticker"
)

// SelectItemsAMAP returns the sum of weights and the selected item ID slice as many as possible under the constraints that the sum of weights is at most W.
//
// This function returns an error if some weights are negative.
func SelectItemsAMAP(weights []float32, W float32) (float32, []int, error) {
	if W < 0.0 {
		return 0.0, nil, fmt.Errorf("W must be non-negative")
	}
	values, backptrs := make([]map[float32]float32, len(weights)), make([]map[float32]int, len(weights))
	for i, wi := range weights {
		if wi < 0.0 {
			return 0.0, nil, fmt.Errorf("weights must have non-negatives")
		}
		values[i], backptrs[i] = make(map[float32]float32), make(map[float32]int)
	}
	var getValue func(i int, weight float32) (float32, int)
	getValue = func(i int, weight float32) (float32, int) {
		if i < 0 {
			return 0, -1
		}
		if value, ok := values[i][weight]; ok {
			return value, backptrs[i][weight]
		}
		valueA, backptrA := getValue(i-1, weight)
		if wi := weights[i]; wi > weight {
			values[i][weight], backptrs[i][weight] = valueA, backptrA
		} else {
			valueB, _ := getValue(i-1, weight-wi)
			value, backptr := valueB+wi, i
			if value < valueA {
				value, backptr = valueA, backptrA
			}
			values[i][weight], backptrs[i][weight] = value, backptr
		}
		return values[i][weight], backptrs[i][weight]
	}
	value, backptr := getValue(len(weights)-1, W)
	ids := []int{}
	for backptr >= 0 {
		ids = append(ids, backptr)
		W -= weights[backptr]
		_, backptr = getValue(backptr-1, W)
	}
	sort.Ints(ids)
	return value, ids, nil
}

// BipartitionWeightedGraph bipartitions the given weighted graph with 1-spectral clustering method, and returns the bool slice which is true if it is in right.
// If the size of left or right is less than minLeftRight, this function won't try to bipartition.
//
// This function returns an error if A is not legal adjacency matrix.
func BipartitionWeightedGraph(n uint64, minLeftRight uint64, A []float32) ([]bool, error) {
	rng := rand.New(rand.NewSource(0))
	if n == 0 {
		return []bool{}, nil
	}
	if uint64(len(A)) != n*n {
		return nil, fmt.Errorf("A is not legal adjacency matrix")
	}
	// Check the connected components.
	compIds, ncomps := make([]uint64, n), uint64(0)
	stack := []uint64{}
	for root := uint64(0); root < n; root++ {
		if compIds[root] == 0 {
			stack = append(stack, root)
		}
		for len(stack) > 0 {
			i := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if compIds[i] == 0 {
				ncomps++
				compIds[i] = ncomps
			}
			for j := uint64(0); j < n; j++ {
				if w := A[i*n+j]; w > 0.0 {
					if compIds[j] == 0 {
						compIds[j] = compIds[i]
						stack = append(stack, j)
					}
				} else if w < 0.0 {
					return nil, fmt.Errorf("A is not legal adjacency matrix")
				}
			}
		}
	}
	if ncomps >= 2 {
		compIdOrder := make([]uint64, ncomps)
		for id := uint64(0); id < ncomps; id++ {
			compIdOrder[id] = id + 1
		}
		compSizes := make([]float32, ncomps)
		sumWeight := float32(0.0)
		for i := uint64(0); i < n; i++ {
			w := float32(0.0)
			for j := uint64(0); j < n; j++ {
				w += A[i*n+j]
			}
			compSizes[compIds[i]-1] += w
			sumWeight += w
		}
		value, items, err := SelectItemsAMAP(compSizes, sumWeight/2.0)
		if err != nil {
			return nil, fmt.Errorf("SelectItemsAMAP: %s", err)
		}
		if value >= float32(minLeftRight) && sumWeight-value >= float32(minLeftRight) {
			leftCompIdSet := make(map[uint64]bool)
			for _, item := range items {
				leftCompIdSet[compIdOrder[item]] = true
			}
			delta := make([]bool, n)
			for i, compId := range compIds {
				delta[i] = leftCompIdSet[compId]
			}
			return delta, nil
		}
	}
	// Try 1-spectral clustering on the connected components.
	// Calculate the negative graph Laplacian -L := -(D - A), D := diag(A1_n).
	// -L becomes the negative semi-definite matrix.
	negL := make([]float32, len(A))
	for i := uint64(0); i < n; i++ {
		di := float32(0.0)
		for j := uint64(0); j < n; j++ {
			negL[i*n+j] = A[i*n+j]
			di += A[i*n+j]
		}
		negL[i*n+i] = A[i*n+i] - di
	}
	// Initialize the random L1-unit vector.
	p, p0 := make([]float32, n), make([]float32, n)
	pL1 := float32(0.0)
	for i := uint64(0); i < n; i++ {
		p[i] = rng.Float32() - 0.5
		pL1 += sticker.Abs32(p[i])
	}
	for i := uint64(0); i < n; i++ {
		p[i] /= pL1
	}
	// Calculate the smallest eigen-value of -L.
	t, T := 0, 1000
	absLambdaMin := float32(0.0)
	for {
		if t++; t >= T {
			break
		}
		p, p0 = p0, p
		pL1, deltaL1 := float32(0.0), float32(0.0)
		for i := uint64(0); i < n; i++ {
			pi := float32(0.0)
			for j := uint64(0); j < n; j++ {
				pi += negL[i*n+j] * p0[j]
			}
			p[i] = pi
			pL1 += sticker.Abs32(p[i])
		}
		for i := uint64(0); i < n; i++ {
			p[i] /= pL1
			deltaL1 += sticker.Abs32(p[i] - p0[i])
		}
		absLambdaMin = pL1
		if deltaL1 <= 1.0e-05 {
			break
		}
	}
	// Re-initialize the random L1-unit vector.
	pL1 = 0.0
	for i := uint64(0); i < n; i++ {
		p[i] = rng.Float32() - 0.5
		pL1 += sticker.Abs32(p[i])
	}
	for i := uint64(0); i < n; i++ {
		p[i] /= pL1
	}
	// Calculate the second-largest eigen-value/vector of -L which are the largest eigen-value/vector of -L - lambdaMin respectively.
	t, T = 0, 1000
	for {
		if t++; t >= T {
			break
		}
		p, p0 = p0, p
		pL1, deltaL1 := float32(0.0), float32(0.0)
		s := float32(0.0)
		for i := uint64(0); i < n; i++ {
			pi := float32(0.0)
			for j := uint64(0); j < n; j++ {
				pi += negL[i*n+j] * p0[j]
			}
			p[i] = pi + absLambdaMin*p0[i]
			s += p[i]
		}
		s /= float32(n)
		for i := uint64(0); i < n; i++ {
			p[i] -= s
			pL1 += sticker.Abs32(p[i])
		}
		for i := uint64(0); i < n; i++ {
			p[i] /= pL1
			deltaL1 += sticker.Abs32(p[i] - p0[i])
		}
		if deltaL1 <= 1.0e-05 {
			break
		}
	}
	vertexPs := make(sticker.KeyValues32OrderedByValue, n)
	for i := uint64(0); i < n; i++ {
		vertexPs[i] = sticker.KeyValue32{uint32(i), p[i]}
	}
	sort.Sort(vertexPs)
	delta := make([]bool, n)
	for i := 0; i < len(vertexPs)/2; i++ {
		delta[vertexPs[i].Key] = true
	}
	return delta, nil
}
