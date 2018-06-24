package plugin

import (
	"log"
	"math/rand"
	"sort"

	"github.com/hiro4bbh/sticker"
)

// LeftRightAssigner_greedyBottomRanks assigns left or right label as moving each entry which has that bottom-ranked labels from left to right.
// This is registered to sticker.LeftRightAssigners.
//
// This function returns no error currently.
func LeftRightAssigner_greedyBottomRanks(ds *sticker.Dataset, delta []bool, debug *log.Logger) error {
	for i := range delta {
		delta[i] = false
	}
	labelFreq := make(map[uint32]float32)
	for _, yi := range ds.Y {
		for _, label := range yi {
			labelFreq[label]++
		}
	}
	labelRankTopL := sticker.RankTopK(labelFreq, uint(len(labelFreq)))
	n, nlefts, nrights := ds.Size(), ds.Size(), 0
	for k := len(labelRankTopL) - 1; nlefts > n/2; k-- {
		if debug != nil {
			debug.Printf("LeftRightAssigner(greedyBottomRanks): optimizing the allocation on %d in left and %d in right ...", nlefts, nrights)
		}
		label := labelRankTopL[k]
		for i, yi := range ds.Y {
			if delta[i] {
				continue
			}
			for _, l := range yi {
				if l == label {
					delta[i] = true
					nlefts--
					nrights++
					break
				}
			}
		}
	}
	return nil
}

// LeftRightAssigner_nDCG assigns left or right on each label as maximizing the sum of left and right utilities with nDCGs.
// This is registered to LeftRightAssigners.
//
// This function return no error currently.
//
// NOTICE: In calculating nDCG, this function uses the base of logarithm is 2 because of precision.
func LeftRightAssigner_nDCG(ds *sticker.Dataset, delta []bool, debug *log.Logger) error {
	delta0, delta1 := make([]bool, len(delta)), make([]bool, len(delta))
	copy(delta0, delta)
	objval0 := float32(0.0)
	for {
		if debug != nil {
			nLeftRights := make(map[bool]int)
			for _, deltai := range delta0 {
				nLeftRights[deltai]++
			}
			debug.Printf("LeftRightAssigner_nDCG: optimizing the allocation on %d in left and %d in right (objval0=%g) ...", nLeftRights[false], nLeftRights[true], objval0)
		}
		// Construct the left/right label distributions r^-/r^+.
		leftLabelFreq, rightLabelFreq := make(map[uint32]float32), make(map[uint32]float32)
		for i, yi := range ds.Y {
			labelFreq := leftLabelFreq
			if delta0[i] {
				labelFreq = rightLabelFreq
			}
			Z := 1.0 / sticker.IdealDCG(uint(len(yi)))
			for _, label := range yi {
				labelFreq[label] += Z
			}
		}
		leftLabelRanks, rightLabelRanks := sticker.RankTopK(leftLabelFreq, uint(len(leftLabelFreq))), sticker.RankTopK(rightLabelFreq, uint(len(rightLabelFreq)))
		leftLabelInvRanks, rightLabelInvRanks := sticker.InvertRanks(leftLabelRanks), sticker.InvertRanks(rightLabelRanks)
		objval := float32(0.0)
		for i, yi := range ds.Y {
			// v_i^- = L_{nDCG@K}(y_i, r^-), v_i^+ = L_{nDCG@K}(y_i, r^+)
			vn, vp := float32(0.0), float32(0.0)
			for _, label := range yi {
				if rank, ok := leftLabelInvRanks[label]; ok {
					vn += 1.0 / sticker.LogBinary32(1.0+float32(rank))
				}
				if rank, ok := rightLabelInvRanks[label]; ok {
					vp += 1.0 / sticker.LogBinary32(1.0+float32(rank))
				}
			}
			if vn < vp {
				delta1[i] = true
				objval += vp
			} else {
				delta1[i] = false
				objval += vn
			}
		}
		if objval0 >= objval {
			break
		}
		delta0, delta1 = delta1, delta0
		objval0 = objval
	}
	copy(delta, delta0)
	return nil
}

// LeftRightAssigner_none assigns left or right on each label with the given initialized delta, so the label assignment won't change from the given initialized delta.
// This is registered to LeftRightAssigners.
//
// This function returns no error.
func LeftRightAssigner_none(ds *sticker.Dataset, delta []bool, debug *log.Logger) error {
	return nil
}

// LeftRightAssigner is the type of the left/right assigners.
// An assigner modifies delta to store the result of the assignment.
// delta is also used as the initial value.
type LeftRightAssigner func(ds *sticker.Dataset, delta []bool, debug *log.Logger) error

// LeftRightAssigners is the map from the assigner name to the left/right assigner.
var LeftRightAssigners = map[string]LeftRightAssigner{
	"greedyBottomRanks": LeftRightAssigner_greedyBottomRanks,
	"nDCG":              LeftRightAssigner_nDCG,
	"none":              LeftRightAssigner_none,
}

// DefaultLeftRightAssignerName is the default LeftRightAssigner name.
const DefaultLeftRightAssignerName = "nDCG"

func _LeftRightAssignInitializer_topLabels(isGraph bool, ds *sticker.Dataset, params *LabelTreeParameters, rng *rand.Rand, debug *log.Logger) []bool {
	labelFreq := make(map[uint32]float32)
	for _, yi := range ds.Y {
		for _, label := range yi {
			labelFreq[label]++
		}
	}
	K := 16 * params.K
	if K > uint(len(labelFreq)) {
		K = uint(len(labelFreq))
	}
	labelRankTopK := sticker.RankTopK(labelFreq, K)
	labelInvRankTopK := sticker.InvertRanks(labelRankTopK)
	labelFreqTopK := make(map[uint32]float32)
	for _, label := range labelRankTopK {
		labelFreqTopK[label] = labelFreq[label]
	}
	A := make([]float32, K*K)
	if isGraph {
		for _, yi := range ds.Y {
			for _, label1 := range yi {
				rank1p, ok := labelInvRankTopK[label1]
				if !ok {
					continue
				}
				for _, label2 := range yi {
					if rank2p, ok := labelInvRankTopK[label2]; ok && rank1p != rank2p {
						A[uint(rank1p-1)*K+uint(rank2p-1)]++
					}
				}
			}
		}
	} else {
		for _, yi := range ds.Y {
			labelRanks := make(sticker.KeyValues32OrderedByValue, 0, len(yi))
			for _, label := range yi {
				if rank, ok := labelInvRankTopK[label]; ok {
					labelRanks = append(labelRanks, sticker.KeyValue32{label, float32(rank)})
				}
			}
			sort.Sort(labelRanks)
			for j := 0; j < len(labelRanks)-1; j++ {
				A[uint(labelRanks[j].Value-1)*K+uint(labelRanks[j+1].Value-1)]++
				A[uint(labelRanks[j+1].Value-1)*K+uint(labelRanks[j].Value-1)]++
			}
		}
	}
	deltaLabel, _ := BipartitionWeightedGraph(uint64(K), uint64(params.MaxEntriesInLeaf), A)
	delta := make([]bool, ds.Size())
	for i, yi := range ds.Y {
		dn, dp := 0, 0
		for _, label := range yi {
			if rank, ok := labelInvRankTopK[label]; ok {
				if deltaLabel[rank-1] {
					dp++
				} else {
					dn++
				}
			}
		}
		if dp >= dn {
			if dn+dp == 0 {
				delta[i] = rng.Float32() >= 0.5
			} else {
				delta[i] = true
			}
		} else if dp < dn {
			delta[i] = false
		}
	}
	return delta
}

// LeftRightAssignInitializer_topLabelGraph returns the delta slice initialized with the cutting of the top-label graph.
// This is registered to sticker.LeftRightAssignInitializers.
func LeftRightAssignInitializer_topLabelGraph(ds *sticker.Dataset, params *LabelTreeParameters, rng *rand.Rand, debug *log.Logger) []bool {
	return _LeftRightAssignInitializer_topLabels(true, ds, params, rng, debug)
}

// LeftRightAssignInitializer_topLabelGraph returns the delta slice initialized with the cutting of the top-labels tree.
// This is registered to sticker.LeftRightAssignInitializers.
func LeftRightAssignInitializer_topLabelTree(ds *sticker.Dataset, params *LabelTreeParameters, rng *rand.Rand, debug *log.Logger) []bool {
	return _LeftRightAssignInitializer_topLabels(false, ds, params, rng, debug)
}

// LeftRightAssignInitializer_uniform returns the delta slice initialized with the samples from uniform probability distribution.
// This is registered to LeftRightAssignInitializers.
func LeftRightAssignInitializer_uniform(ds *sticker.Dataset, params *LabelTreeParameters, rng *rand.Rand, debug *log.Logger) []bool {
	delta := make([]bool, ds.Size())
	for i := range delta {
		delta[i] = rng.Float32() >= 0.5
	}
	return delta
}

// LeftRightAssignInitializer is the type of the left/right assignment initializers.
// An initializer returns the initialized left/right assignment slice.
type LeftRightAssignInitializer func(ds *sticker.Dataset, params *LabelTreeParameters, rng *rand.Rand, debug *log.Logger) []bool

// LeftRightAssignInitializers is the map from the initializer name to the corresponding left/right assignment initializer.
var LeftRightAssignInitializers = map[string]LeftRightAssignInitializer{
	"topLabelGraph": LeftRightAssignInitializer_topLabelGraph,
	"topLabelTree":  LeftRightAssignInitializer_topLabelTree,
	"uniform":       LeftRightAssignInitializer_uniform,
}

// DefaultLeftRightAssignInitializerName is the default LeftRightAssignInitializer name.
const DefaultLeftRightAssignInitializerName = "uniform"
