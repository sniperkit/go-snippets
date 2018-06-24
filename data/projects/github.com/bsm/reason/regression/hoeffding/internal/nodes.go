package internal

import (
	"math"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression"
	"github.com/bsm/reason/util"
)

// NewNode inits a node
func NewNode(kind isNode_Kind, stats *util.StreamStats) *Node {
	if stats == nil {
		stats = new(util.StreamStats)
	}
	return &Node{Kind: kind, Stats: stats}
}

// IsSufficient returns true when a node has sufficient stats.
func (n *Node) IsSufficient() bool {
	stdev := n.Stats.StdDev()
	return stdev != 0.0 && !math.IsNaN(stdev)
}

// Weight returns the weight observed on the node.
func (n *Node) Weight() float64 {
	return n.Stats.Weight
}

// --------------------------------------------------------------------

func (n *SplitNode) childCat(feature *core.Feature, x core.Example) core.Category {
	switch feature.Kind {
	case core.Feature_CATEGORICAL:
		return feature.Category(x)
	case core.Feature_NUMERICAL:
		if feature.Number(x) < n.Pivot {
			return 0
		}
		return 1
	default:
		return core.NoCategory
	}
}

// --------------------------------------------------------------------

// Enable enables the node.
func (n *LeafNode) Enable() {
	if n.IsDisabled {
		n.IsDisabled = false
		n.FeatureStats = make(map[string]*FeatureStats)
	}
}

// Disable disables the node.
func (n *LeafNode) Disable() {
	n.IsDisabled = true
	n.FeatureStats = nil
}

// EvaluateSplit evaluates a split for a fiven feature.
// Returns nil if a split is not possible.
func (n *LeafNode) EvaluateSplit(feature string, crit regression.SplitCriterion, self *Node) *SplitCandidate {
	if n.IsDisabled {
		return nil
	}

	if n.FeatureStats == nil {
		return nil
	}

	stats, ok := n.FeatureStats[feature]
	if !ok {
		return nil
	}

	switch kind := stats.Kind.(type) {
	case *FeatureStats_Numerical_:
		var c *SplitCandidate
		s := kind.Numerical
		r := crit.Range(self.Stats)

		for _, pivot := range s.PivotPoints() {
			post := s.PostSplit(pivot)
			merit := crit.Merit(self.Stats, post)
			if c == nil || merit > c.Merit {
				c = &SplitCandidate{
					Feature:   feature,
					Merit:     merit,
					Range:     r,
					Pivot:     pivot,
					PreSplit:  self.Stats,
					PostSplit: post,
				}
			}
		}
		return c
	case *FeatureStats_Categorical_:
		if s := kind.Categorical; s.Len() > 1 {
			post := s.PostSplit()
			return &SplitCandidate{
				Feature:   feature,
				Merit:     crit.Merit(self.Stats, post),
				Range:     crit.Range(self.Stats),
				PreSplit:  self.Stats,
				PostSplit: post,
			}
		}
	}
	return nil
}

// Observe observes an example and updates internal stats.
func (n *LeafNode) Observe(m *core.Model, target *core.Feature, x core.Example, weight float64, self *Node) {
	// Get the target value, skip this example on "no value"
	targetVal := target.Number(x)
	if !core.IsNum(targetVal) {
		return
	}

	// Get example weight and update node stats
	self.Stats.Add(targetVal, weight)

	// Skip the remaining steps if this node is disabled
	if n.IsDisabled {
		return
	}

	// Ensure we have stats
	if n.FeatureStats == nil {
		n.FeatureStats = make(map[string]*FeatureStats)
	}

	// Update each predictor feature's stats with a target-value, predictor-value
	// and weight tuple
	for name, feat := range m.Features {
		if name == target.Name {
			continue // skip target, we are only interested in predictors
		}

		stats := n.FeatureStats[feat.Name]
		if stats == nil {
			stats = new(FeatureStats)
			n.FeatureStats[feat.Name] = stats
		}

		switch feat.Kind {
		case core.Feature_CATEGORICAL:
			if cat := feat.Category(x); core.IsCat(cat) {
				stats.FetchCategorical().Add(cat, targetVal, weight)
			}
		case core.Feature_NUMERICAL:
			if num := feat.Number(x); core.IsNum(num) {
				stats.FetchNumerical().Add(num, targetVal, weight)
			}
		}
	}
}
