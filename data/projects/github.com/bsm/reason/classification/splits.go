package classification

import (
	"math"

	"github.com/bsm/reason/internal/splits"
	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for classifications.
type SplitCriterion interface {
	// Range returns the range of the split merit
	Range(pre *util.Vector) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.Vector, post *util.VectorDistribution) float64
}

// DefaultSplitCriterion returns the default split criterion:
//   InformationGain{MinBranchFraction: 0.1}
func DefaultSplitCriterion() SplitCriterion {
	return InformationGain{MinBranchFraction: 0.1}
}

// --------------------------------------------------------------------

// GiniImpurity determines split merit using Gini impurity
type GiniImpurity struct{}

// Range implements SplitCriterion
func (GiniImpurity) Range(_ *util.Vector) float64 { return 1.0 }

// Merit implements SplitCriterion
func (GiniImpurity) Merit(pre *util.Vector, post *util.VectorDistribution) float64 {
	if pre == nil || post == nil {
		return 0.0
	}

	total := 0.0
	post.ForEach(func(_ int, vv *util.Vector) bool {
		total += vv.Weight()
		return true
	})
	if total == 0 {
		return 0.0
	}

	merit := 0.0
	post.ForEach(func(_ int, vv *util.Vector) bool {
		sum := vv.Weight()
		merit += sum / total * calcGiniSplit(vv, sum)
		return true
	})
	return splits.NormMerit(merit)
}

func calcGiniSplit(vv *util.Vector, sum float64) float64 {
	res := 1.0
	vv.ForEachValue(func(v float64) bool {
		sub := v / sum
		res -= sub * sub
		return true
	})
	return res
}

// InformationGain determines split merit through information gain
type InformationGain struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFraction float64
}

// Range implements SplitCriterion
func (InformationGain) Range(pre *util.Vector) float64 {
	if pre != nil {
		if sz := pre.Len(); sz > 2 {
			return math.Log2(float64(sz))
		}
	}
	return math.Log2(2.0)
}

// Merit implements SplitCriterion
func (c InformationGain) Merit(pre *util.Vector, post *util.VectorDistribution) float64 {
	if pre == nil || post == nil {
		return 0.0
	}

	total := 0.0
	count := 0
	post.ForEach(func(_ int, vv *util.Vector) bool {
		total += vv.Weight()
		count++
		return true
	})
	if count < 2 || total == 0 {
		return 0.0
	}

	if min := c.MinBranchFraction; min > 0 {
		count = 0
		post.ForEach(func(_ int, vv *util.Vector) bool {
			if vv.Weight()/total > min {
				if count++; count > 1 {
					return false
				}
			}
			return true
		})
	}
	if count < 2 {
		return 0.0
	}

	e1, e2 := pre.Entropy(), 0.0
	post.ForEach(func(_ int, vv *util.Vector) bool {
		e2 += vv.Weight() * vv.Entropy()
		return true
	})
	return splits.NormMerit(e1 - e2/total)
}

// --------------------------------------------------------------------

// GainRatio normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

// Merit implements SplitCriterion
func (c GainRatio) Merit(pre *util.Vector, post *util.VectorDistribution) float64 {
	merit := c.SplitCriterion.Merit(pre, post)
	penalty := new(splits.GainRatioPenality)
	post.ForEach(func(_ int, vv *util.Vector) bool {
		penalty.Weight += vv.Weight()
		return true
	})
	post.ForEach(func(_ int, vv *util.Vector) bool {
		penalty.Update(vv.Weight())
		return true
	})
	return splits.NormMerit(merit / penalty.Value())
}
