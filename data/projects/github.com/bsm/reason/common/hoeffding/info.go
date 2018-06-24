package hoeffding

import (
	"fmt"
	"strings"
)

// TreeInfo contains tree information/stats
type TreeInfo struct {
	NumNodes    int // the total number of nodes
	NumLearning int // the number of learning leaves
	NumDisabled int // the number of disable leaves
	MaxDepth    int // the maximum depth
}

// SplitCandidateInfo contains information about
// a candidate of a split attempt.
type SplitCandidateInfo struct {
	Feature string  // the feature name
	Merit   float64 // the merit
}

// SplitAttemptInfo instances may be emitted as part of the the tree training.
// They contain information about attempted splits.
type SplitAttemptInfo struct {
	// The weight at the time of the evaluation
	Weight float64
	// Indicator of a successful split
	Success bool
	// The posssible merit gain of this split attempt.
	MeritGain float64
	// The hoeffding bound of this split attempt.
	HoeffdingBound float64
	// Split candidates
	Candidates []SplitCandidateInfo
}

// String returns a one-liner summary.
func (t *SplitAttemptInfo) String() string {
	candidates := make([]string, 0, 4)
	for _, ct := range t.Candidates {
		if ct.Merit == 0 {
			continue
		}

		s := fmt.Sprintf("%s=%.2f", ct.Feature, ct.Merit)
		if candidates = append(candidates, s); len(candidates) == 4 {
			break
		}
	}
	return fmt.Sprintf("Weight: %.1f, Success: %v, MeritGain: %.2f, HBound: %.2f, Candidates: [%s]",
		t.Weight, t.Success, t.MeritGain, t.HoeffdingBound, strings.Join(candidates, " "))
}
