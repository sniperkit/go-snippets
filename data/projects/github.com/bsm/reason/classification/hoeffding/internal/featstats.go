package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/hoeffding"
	"github.com/bsm/reason/util"
)

// FetchCategorical fetches categorical stats.
func (s *FeatureStats) FetchCategorical() *FeatureStats_Categorical {
	stats := s.GetCategorical()
	if stats == nil {
		stats = new(FeatureStats_Categorical)
		s.Kind = &FeatureStats_Categorical_{Categorical: stats}
	}
	return stats
}

// FetchNumerical fetches numerical stats.
func (s *FeatureStats) FetchNumerical() *FeatureStats_Numerical {
	stats := s.GetNumerical()
	if stats == nil {
		stats = new(FeatureStats_Numerical)
		s.Kind = &FeatureStats_Numerical_{Numerical: stats}
	}
	return stats
}

// --------------------------------------------------------------------

// PostSplit calculates a post-split distribution from previous observations.
func (s *FeatureStats_Categorical) PostSplit() *util.VectorDistribution {
	return &s.VectorDistribution
}

// Add adds an observation
func (s *FeatureStats_Categorical) Add(featCat, targetCat core.Category, weight float64) {
	s.VectorDistribution.Add(int(featCat), int(targetCat), weight)
}

// --------------------------------------------------------------------

// Add adds an observation
func (s *FeatureStats_Numerical) Add(featVal float64, targetCat core.Category, weight float64) {
	targetPos := int(targetCat)
	if v := s.Min.Get(targetPos); v == 0 || featVal < v {
		s.Min.Set(targetPos, featVal)
	}
	if v := s.Max.Get(targetPos); v == 0 || featVal > v {
		s.Max.Set(targetPos, featVal)
	}
	s.Stats.Add(targetPos, featVal, weight)
}

// PivotPoints determines the optimum split points for the range of values.
func (s *FeatureStats_Numerical) PivotPoints() []float64 {
	var tmin, tmax float64
	s.Min.ForEach(func(i int, min float64) bool {
		if tmin == 0 || min < tmin {
			tmin = min
		}
		if max := s.Max.Get(i); tmax == 0 || max > tmax {
			tmax = max
		}
		return true
	})
	return hoeffding.PivotPoints(tmin, tmax)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_Numerical) PostSplit(pivot float64) *util.VectorDistribution {
	res := new(util.VectorDistribution)
	s.Stats.ForEach(func(i int, x *util.StreamStats) bool {
		if min := s.Min.Get(i); min > 0 && pivot < min {
			res.Add(1, i, x.Weight)
		} else if max := s.Max.Get(i); max > 0 && pivot >= max {
			res.Add(0, i, x.Weight)
		} else {
			lt, eq, gt := x.Estimate(pivot)
			res.Add(0, i, lt+eq)
			res.Add(1, i, gt)
		}
		return true
	})
	return res
}
