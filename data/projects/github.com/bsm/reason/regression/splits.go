package regression

import (
	"math"

	"github.com/bsm/reason/internal/splits"
	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for regressions
type SplitCriterion interface {
	// Range returns the range of the split merit.
	Range(pre *util.StreamStats) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.StreamStats, post *util.StreamStatsDistribution) float64
}

// DefaultSplitCriterion returns the default split criterion:
//   VarReductionSplitCriterion{MinWeight: 4.0}
func DefaultSplitCriterion() SplitCriterion {
	return VarianceReduction{MinWeight: 4.0}
}

// VarianceReduction performs splits using variance-reduction
type VarianceReduction struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 4.0
	MinWeight float64
}

// Range implements SplitCriterion
func (VarianceReduction) Range(_ *util.StreamStats) float64 { return 1.0 }

// Merit implements SplitCriterion
func (c VarianceReduction) Merit(pre *util.StreamStats, post *util.StreamStatsDistribution) float64 {
	if pre == nil {
		return 0.0
	}

	n := 0
	sumW := 0.0
	post.ForEach(func(_ int, s *util.StreamStats) bool {
		if w := s.Weight; w >= c.MinWeight {
			sumW += w
			n++
		}
		return true
	})
	if n < 2 || sumW == 0 {
		return 0.0
	}

	preV := pre.Variance()
	if math.IsNaN(preV) {
		return 0.0
	}

	postV := 0.0
	post.ForEach(func(_ int, s *util.StreamStats) bool {
		if w, v := s.Weight, s.Variance(); w >= c.MinWeight && !math.IsNaN(v) {
			postV += w * v / sumW
		}
		return true
	})
	return splits.NormMerit(preV - postV)
}

// --------------------------------------------------------------------

// GainRatio wraps a split criterion and normalises the merits
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

func (c GainRatio) Merit(pre *util.StreamStats, post *util.StreamStatsDistribution) float64 {
	merit := c.SplitCriterion.Merit(pre, post)
	penalty := new(splits.GainRatioPenality)
	post.ForEach(func(_ int, s *util.StreamStats) bool {
		penalty.Weight += s.Weight
		return true
	})
	post.ForEach(func(_ int, s *util.StreamStats) bool {
		penalty.Update(s.Weight)
		return true
	})
	return splits.NormMerit(merit / penalty.Value())
}
