package util

import (
	"math"

	"github.com/bsm/reason/internal/sparsedense"
	"gonum.org/v1/gonum/stat/distuv"
)

// Add adds a new value with a weight
func (s *StreamStats) Add(value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return
	}

	wv := weight * value
	s.Weight += weight
	s.Sum += wv
	s.SumSquares += wv * value
}

// IsZero returns true if there are no values in the series
func (s *StreamStats) IsZero() bool { return s.Weight <= 0 }

// Mean returns a mean average
func (s *StreamStats) Mean() float64 { return s.Sum / s.Weight }

// Variance is the sample variance of the series
func (s *StreamStats) Variance() float64 {
	if s.Weight <= 1 {
		return 0.0
	}
	x := (s.Sum * s.Sum) / s.Weight
	return (s.SumSquares - x) / (s.Weight - 1)
}

// StdDev is the sample standard deviation of the series
func (s *StreamStats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Prob calculates the gaussian probability density of a value
func (s *StreamStats) Prob(value float64) float64 {
	if sig := s.StdDev(); !math.IsNaN(sig) {
		return distuv.Normal{Mu: s.Mean(), Sigma: sig}.Prob(value)
	}
	return math.NaN()
}

// Estimate estimates weight boundaries for a given value
func (s *StreamStats) Estimate(value float64) (lessThan float64, equalTo float64, greaterThan float64) {
	equalTo = s.Prob(value) * s.Weight

	if sig := s.StdDev(); !math.IsNaN(sig) {
		lessThan = distuv.Normal{
			Mu:    s.Mean(),
			Sigma: sig,
		}.CDF(value)*s.Weight - equalTo
	} else {
		lessThan = math.NaN()
	}

	if greaterThan = s.Weight - equalTo - lessThan; greaterThan < 0 {
		greaterThan = 0
	}

	return
}

// --------------------------------------------------------------------

// Get returns the series at index
func (x *StreamStatsDistribution) Get(index int) *StreamStats {
	if index < 0 {
		return nil
	}

	if x.Dense != nil && index < len(x.Dense) {
		return x.Dense[index].StreamStats
	} else if x.Sparse != nil {
		return x.Sparse[int64(index)]
	}
	return nil
}

// Add adds and observation of a weighted value at index.
func (x *StreamStatsDistribution) Add(index int, value, weight float64) {
	if index < 0 {
		return
	}

	if x.Dense != nil {
		x.fetchDense(index).Add(value, weight)
		return
	}

	if x.Sparse == nil {
		x.Sparse = make(map[int64]*StreamStats)
	}
	x.fetchSparse(int64(index)).Add(value, weight)
}

// Len returns the number of elements in the distribution.
func (x *StreamStatsDistribution) Len() int {
	if x.Dense != nil {
		n := 0
		x.ForEach(func(_ int, _ *StreamStats) bool { n++; return true })
		return n
	}

	return len(x.Sparse)
}

// ForEach iterates over each stats eleemnt in the distribution.
func (x *StreamStatsDistribution) ForEach(iter func(int, *StreamStats) bool) {
	if x.Dense != nil {
		for i, s := range x.Dense {
			if s.StreamStats != nil && !iter(i, s.StreamStats) {
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

func (x *StreamStatsDistribution) fetchDense(index int) *StreamStats {
	if n := index + 1; n > len(x.Dense) {
		dense := make([]StreamStatsDistribution_Dense, n, n*2)
		copy(dense, x.Dense)
		x.Dense = dense
	} else if n > cap(x.Dense) {
		x.Dense = x.Dense[:n]
	}

	s := x.Dense[index].StreamStats
	if s == nil {
		s = new(StreamStats)
		x.Dense[index].StreamStats = s
	}
	return s
}

func (x *StreamStatsDistribution) fetchSparse(index int64) *StreamStats {
	if n := index + 1; n > x.SparseCap {
		x.SparseCap = n
	}

	s, ok := x.Sparse[index]
	if !ok {
		if sparsedense.BetterOffDense(len(x.Sparse), int(x.SparseCap)) {
			x.convertToDense()
			return x.fetchDense(int(index))
		}

		s = new(StreamStats)
		x.Sparse[index] = s
	}
	return s
}

func (x *StreamStatsDistribution) convertToDense() {
	dense := make([]StreamStatsDistribution_Dense, int(x.SparseCap))
	x.ForEach(func(i int, s *StreamStats) bool {
		dense[i].StreamStats = s
		return true
	})
	x.Dense = dense
	x.Sparse = nil
	x.SparseCap = 0
}
