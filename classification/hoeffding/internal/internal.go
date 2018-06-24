package internal

import "github.com/bsm/reason/util"

// SplitCandidate is a candidate for a split decision
type SplitCandidate struct {
	Feature string  // the feature name
	Merit   float64 // the split merit
	Range   float64 // the split range
	Pivot   float64 // the split pivot, for binary splits

	// Pre-split stats
	PreSplit *util.Vector
	// Post-split stats
	PostSplit *util.VectorDistribution
}

// SplitCandidates are a sortable collection of split candidates
type SplitCandidates []SplitCandidate

func (p SplitCandidates) Len() int           { return len(p) }
func (p SplitCandidates) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SplitCandidates) Less(i, j int) bool { return p[i].Merit < p[j].Merit }
