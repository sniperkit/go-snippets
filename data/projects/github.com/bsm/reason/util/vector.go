package util

import (
	"math"

	"github.com/bsm/reason/internal/sparsedense"
)

// NewVectorFromSlice initalizes a vector using a slice of weights.
func NewVectorFromSlice(weights ...float64) *Vector {
	return &Vector{Dense: weights}
}

// Len returns the number of non-zero weights in the vector.
func (vv *Vector) Len() int {
	if vv.Dense != nil {
		n := 0
		for _, v := range vv.Dense {
			if v != 0 {
				n++
			}
		}
		return n
	}

	return len(vv.Sparse)
}

// Weight returns the total of all weights of the vector.
func (vv *Vector) Weight() float64 {
	sum := 0.0
	vv.ForEachValue(func(w float64) bool { sum += w; return true })
	return sum
}

// Min returns the position of the minimum weight in the vector (with value).
func (vv *Vector) Min() (pos int, weight float64) {
	pos = -1
	vv.ForEach(func(i int, w float64) bool {
		if i == 0 || w < weight {
			pos = i
			weight = w
		}
		return true
	})
	return pos, weight
}

// Max returns the position od the maximum weight in the vector (with value).
func (vv *Vector) Max() (pos int, weight float64) {
	pos = -1
	vv.ForEach(func(i int, w float64) bool {
		if w > weight {
			pos = i
			weight = w
		}
		return true
	})
	return
}

// Clone creates a copy of the vector
func (vv *Vector) Clone() *Vector {
	nv := &Vector{SparseCap: vv.SparseCap}
	if vv.Dense != nil {
		nv.Dense = make([]float64, len(vv.Dense))
		copy(nv.Dense, vv.Dense)
	} else {
		nv.Sparse = make(map[int64]float64, len(vv.Sparse))
		for i, w := range vv.Sparse {
			nv.Sparse[i] = w
		}
	}
	return nv
}

// Get returns a value at index
func (vv *Vector) Get(index int) float64 {
	if index < 0 {
		return 0.0
	}

	if vv.Dense != nil && index < len(vv.Dense) {
		return vv.Dense[index]
	} else if vv.Sparse != nil {
		return vv.Sparse[int64(index)]
	}
	return 0.0
}

// Set sets a weight at index
func (vv *Vector) Set(index int, weight float64) {
	if index < 0 {
		return
	}

	if vv.Dense != nil {
		if n := index + 1; n > cap(vv.Dense) {
			dense := make([]float64, n, 2*n)
			copy(dense, vv.Dense)
			vv.Dense = dense
		} else if n > len(vv.Dense) {
			vv.Dense = vv.Dense[:n]
		}
		vv.Dense[index] = weight
		return
	}

	if n := int64(index + 1); n > vv.SparseCap {
		vv.SparseCap = n
	}
	if vv.Sparse == nil {
		vv.Sparse = make(map[int64]float64, 1)
	}
	vv.Sparse[int64(index)] = weight
	vv.tryConvertToDense()
}

// Add increments a weight at index by delta.
func (vv *Vector) Add(index int, delta float64) {
	if index < 0 {
		return
	}

	if vv.Dense != nil {
		if index < len(vv.Dense) {
			vv.Dense[index] += delta
		} else {
			vv.Set(index, delta)
		}
		return
	}

	if n := int64(index + 1); n > vv.SparseCap {
		vv.SparseCap = n
	}
	if vv.Sparse == nil {
		vv.Sparse = make(map[int64]float64, 1)
	}
	vv.Sparse[int64(index)] += delta
	vv.tryConvertToDense()
}

// ForEach iterates over each index/value
func (vv *Vector) ForEach(iter func(int, float64) bool) {
	if vv.Dense != nil {
		for i, w := range vv.Dense {
			if w > 0 && !iter(i, w) {
				break
			}
		}
		return
	}

	for i, w := range vv.Sparse {
		if !iter(int(i), w) {
			break
		}
	}
}

// ForEachValue iterates over each value
func (vv *Vector) ForEachValue(iter func(float64) bool) {
	if vv.Dense != nil {
		for _, w := range vv.Dense {
			if w > 0 && !iter(w) {
				break
			}
		}
		return
	}

	for _, w := range vv.Sparse {
		if !iter(w) {
			break
		}
	}
}

// Clear removes all weights from the vector.
func (vv *Vector) Clear() {
	if vv.Dense != nil {
		vv.Dense = vv.Dense[:0]
	}
	if vv.Sparse != nil {
		vv.Sparse = make(map[int64]float64)
		vv.SparseCap = 0
	}
}

// Mean returns the mean value of the vector.
func (vv *Vector) Mean() float64 {
	_, _, mu := vv.csm()
	return mu
}

// Variance calculates the variance of the vector weights.
func (vv *Vector) Variance() float64 {
	if n, _, mu := vv.csm(); n > 1 {
		var ss, cs float64
		vv.ForEachValue(func(w float64) bool {
			dt := w - mu
			ss += dt * dt
			cs += dt
			return true
		})
		return (ss - cs*cs/float64(n)) / float64(n-1)
	}
	return math.NaN()
}

// Entropy calculates the entropy of the vector weights.
func (vv *Vector) Entropy() float64 {
	ent := 0.0
	sum := 0.0
	vv.ForEachValue(func(w float64) bool {
		ent -= w * math.Log2(w)
		sum += w
		return true
	})
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

// StdDev calculates the standard deviation
func (vv *Vector) StdDev() float64 { return math.Sqrt(vv.Variance()) }

// Normalize normalizes all values to the 0..1 range
func (vv *Vector) Normalize() {
	sum := vv.Weight()
	if sum == 0 {
		return
	}

	if vv.Dense != nil {
		for i, w := range vv.Dense {
			vv.Dense[i] = w / sum
		}
	} else if vv.Sparse != nil {
		for i, w := range vv.Sparse {
			vv.Sparse[i] = w / sum
		}
	}
}

func (vv *Vector) csm() (n int, sum, mu float64) {
	vv.ForEachValue(func(w float64) bool {
		n++
		sum += w
		return true
	})
	if n > 0 {
		mu = sum / float64(n)
	}
	return
}

func (vv *Vector) tryConvertToDense() {
	if !sparsedense.BetterOffDense(len(vv.Sparse), int(vv.SparseCap)) {
		return
	}

	dense := make([]float64, int(vv.SparseCap))
	vv.ForEach(func(i int, w float64) bool {
		dense[i] = w
		return true
	})
	vv.Dense = dense
	vv.Sparse = nil
	vv.SparseCap = 0
}

// --------------------------------------------------------------------

// Get returns the vector at index
func (x *VectorDistribution) Get(index int) *Vector {
	if index < 0 {
		return nil
	}

	if x.Dense != nil && index < len(x.Dense) {
		return x.Dense[index].Vector
	} else if x.Sparse != nil {
		return x.Sparse[int64(index)]
	}
	return nil
}

// Add increments a weight for index at value by delta.
func (x *VectorDistribution) Add(index, value int, delta float64) {
	if index < 0 {
		return
	}

	if x.Dense != nil {
		x.fetchDense(index).Add(value, delta)
		return
	}

	if x.Sparse == nil {
		x.Sparse = make(map[int64]*Vector)
	}
	x.fetchSparse(int64(index)).Add(value, delta)
}

// Len returns the number of elements in the distribution.
func (x *VectorDistribution) Len() int {
	if x.Dense != nil {
		n := 0
		x.ForEach(func(_ int, _ *Vector) bool { n++; return true })
		return n
	}

	return len(x.Sparse)
}

// ForEach iterates over each stats eleemnt in the distribution.
func (x *VectorDistribution) ForEach(iter func(index int, vv *Vector) bool) {
	if x.Dense != nil {
		for i, s := range x.Dense {
			if s.Vector != nil && !iter(i, s.Vector) {
				break
			}
		}
	} else if x.Sparse != nil {
		for i, s := range x.Sparse {
			if !iter(int(i), s) {
				break
			}
		}
	}
}

func (x *VectorDistribution) fetchDense(index int) *Vector {
	if n := index + 1; n > len(x.Dense) {
		dense := make([]VectorDistribution_Dense, n, n*2)
		copy(dense, x.Dense)
		x.Dense = dense
	} else if n > cap(x.Dense) {
		x.Dense = x.Dense[:n]
	}

	s := x.Dense[index].Vector
	if s == nil {
		s = new(Vector)
		x.Dense[index].Vector = s
	}
	return s
}

func (x *VectorDistribution) fetchSparse(index int64) *Vector {
	if n := index + 1; n > x.SparseCap {
		x.SparseCap = n
	}

	s, ok := x.Sparse[index]
	if !ok {
		if sparsedense.BetterOffDense(len(x.Sparse), int(x.SparseCap)) {
			x.convertToDense()
			return x.fetchDense(int(index))
		}

		s = new(Vector)
		x.Sparse[index] = s
	}
	return s
}

func (x *VectorDistribution) convertToDense() {
	dense := make([]VectorDistribution_Dense, int(x.SparseCap))
	x.ForEach(func(i int, vv *Vector) bool {
		dense[i].Vector = vv
		return true
	})
	x.Dense = dense
	x.Sparse = nil
	x.SparseCap = 0
}
