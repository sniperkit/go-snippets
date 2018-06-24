package hoeffding

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"sort"
	"sync"

	"github.com/bsm/reason/classification"
	"github.com/bsm/reason/classification/hoeffding/internal"
	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
)

// Tree is an implementation of a Hoeffding tree.
type Tree struct {
	tree   *internal.Tree
	target *core.Feature

	config Config
	cycles int

	tn []*internal.Node
	mu sync.RWMutex
}

// Load loads a new tree from a reader.
func Load(r io.Reader, config *Config) (*Tree, error) {
	tt := new(internal.Tree)
	if _, err := tt.ReadFrom(r); err != nil {
		return nil, err
	}
	return newTree(tt, config)
}

// New inits a new Tree using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Tree, error) {
	return newTree(internal.NewTree(model, target), config)
}

func newTree(t *internal.Tree, c *Config) (*Tree, error) {
	var config Config
	if c != nil {
		config = *c
	}
	config.Norm()

	target := t.Model.Feature(t.Target)
	if target == nil {
		return nil, fmt.Errorf("hoeffding: unknown feature %q", t.Target)
	} else if !target.Kind.IsCategorical() {
		return nil, fmt.Errorf("hoeffding: feature %q is not categorical", t.Target)
	}

	return &Tree{
		tree:   t,
		target: target,
		config: config,
	}, nil
}

// Info returns information about the tree
func (t *Tree) Info() *common.TreeInfo {
	info := new(common.TreeInfo)

	t.mu.RLock()
	t.tree.Accumulate(t.tree.Root, 1, info)
	t.mu.RUnlock()

	return info
}

// Prune manually prunes the tree to limit it to maxLearningNodes.
func (t *Tree) Prune(maxLearningNodes int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(maxLearningNodes)
}

// Predict traverses the tree for the given example x and appends a prediction
// for every branch to dst, returning it in the end.
// The predictions will therefore increase in accuracy with the most accurate
// one being the last element of the returned slice.
func (t *Tree) Predict(dst classification.Predictions, x core.Example) classification.Predictions {
	t.mu.RLock()
	defer t.mu.RUnlock()

	t.tree.Traverse(x, t.tree.Root, nil, -1, func(node *internal.Node) {
		dst = append(dst, classification.Prediction{Vector: *node.Stats})
	})
	return dst
}

// Train passes an example x with a weight (usually 1.0) to the tree for training.
func (t *Tree) Train(x core.Example, weight float64) *common.SplitAttemptInfo {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, nodeRef, parent, parentIndex := t.tree.Traverse(x, t.tree.Root, nil, -1, nil)
	if node == nil && parentIndex > -1 {
		if split := parent.GetSplit(); split != nil {
			ref := t.tree.Add(nil)
			node = t.tree.Get(ref)
			split.Children.SetRef(parentIndex, ref)
		}
	}
	if node == nil {
		return nil
	}

	if leaf := node.GetLeaf(); leaf != nil {
		// Observe an example
		leaf.Observe(t.tree.Model, t.target, x, weight, node)

		// Pre-prune, if enabled
		if t.config.PrunePeriod > 0 {
			if t.cycles++; t.config.MaxLearningNodes > 0 && t.cycles%t.config.PrunePeriod == 0 {
				t.prune(t.config.MaxLearningNodes)
			}
		}

		// Check if a split should be attempted
		nodeWeight := node.Weight()
		if leaf.IsDisabled || int(nodeWeight-leaf.WeightAtLastEval) < t.config.GracePeriod {
			return nil
		}

		// Store new weight
		leaf.WeightAtLastEval = nodeWeight

		// Check if we have sufficient stats to perform the split
		if !node.IsSufficient() {
			return nil
		}

		// Try to split
		info := t.attemptSplit(leaf, node, nodeRef, nodeWeight)
		if info.Success {
			if parent == nil {
				t.tree.Root = nodeRef
			} else if split := parent.GetSplit(); split != nil {
				split.Children.SetRef(parentIndex, nodeRef)
			}
		}
		return info
	}

	return nil
}

// WriteTo implements io.WriterTo
func (t *Tree) WriteTo(w io.Writer) (int64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tree.WriteTo(w)
}

// WriteText writes text-based tree output to a writer
func (t *Tree) WriteText(w io.Writer) (int64, error) {
	buf := bufio.NewWriter(w)

	t.mu.RLock()
	defer t.mu.RUnlock()

	nn, err := t.tree.WriteText(w, t.tree.Root, "", "ROOT")
	if err != nil {
		return nn, err
	}
	return nn, buf.Flush()
}

// WriteDOT writes a graph in dot notation to a writer
func (t *Tree) WriteDOT(w io.Writer) (int64, error) {
	buf := bufio.NewWriter(w)
	nw := int64(0)

	n, err := buf.WriteString(`digraph ht {
  edge [fontsize=10];
  node [fontsize=10,shape=box];

`)
	nw += int64(n)
	if err != nil {
		return nw, err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	nn, err := t.tree.WriteDOT(buf, t.tree.Root, "N", "")
	nw += nn
	if err != nil {
		return nw, err
	}

	n, err = buf.WriteString("}\n")
	nw += int64(n)
	if err != nil {
		return nw, err
	}
	return nw, buf.Flush()
}

func (t *Tree) prune(maxLearningNodes int) {
	if maxLearningNodes < 0 {
		return
	}

	t.tn = t.tree.FilterLeaves(t.tn[:0])
	if len(t.tn) <= maxLearningNodes {
		return
	}

	// Sort leaves by weight (highest first)
	sort.Slice(t.tn, func(i, j int) bool {
		return t.tn[i].Weight() >= t.tn[j].Weight()
	})

	// Update node status
	for i, node := range t.tn {
		if leaf := node.GetLeaf(); leaf != nil {
			if i < maxLearningNodes {
				leaf.Enable()
			} else {
				leaf.Disable()
			}
		}
	}
}

func (t *Tree) attemptSplit(leaf *internal.LeafNode, node *internal.Node, nodeRef int64, weight float64) *common.SplitAttemptInfo {
	// Init split info
	info := &common.SplitAttemptInfo{Weight: weight}

	// Init candidates, including a null result
	candidates := make(internal.SplitCandidates, 1, len(leaf.FeatureStats)+1)

	// Calculate a split candiate from each of the leaf stats
	for name := range leaf.FeatureStats {
		if c := leaf.EvaluateSplit(name, t.config.SplitCriterion, node); c != nil {
			candidates = append(candidates, *c)
		}
	}

	// Sort candidates by merit, select first
	sort.Stable(sort.Reverse(candidates))
	best := candidates[0]

	// Calculate the gain between merits of the best and the second-best split
	meritGain := best.Merit
	if len(candidates) > 1 {
		meritGain -= candidates[1].Merit
	}

	// Update info
	info.MeritGain = meritGain
	info.Candidates = make([]common.SplitCandidateInfo, 0, len(candidates))
	for _, c := range candidates {
		info.Candidates = append(info.Candidates, common.SplitCandidateInfo{
			Feature: c.Feature,
			Merit:   c.Merit,
		})
	}

	// Give up if there is no merit gain
	if meritGain <= 0 {
		return info
	}

	// Calculate confidence interval + hoeffding bound
	interval := math.Log(1.0 / t.config.SplitConfidence)
	bound := math.Sqrt(best.Range * best.Range * interval * 0.5 / weight)
	info.HoeffdingBound = bound

	// Determine split
	if meritGain > bound || bound < t.config.TieThreshold {
		info.Success = true
		t.tree.Split(nodeRef, best.Feature, best.PreSplit, best.PostSplit, best.Pivot)
	}
	return info
}
