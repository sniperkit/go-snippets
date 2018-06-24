package plugin

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/hiro4bbh/sticker"
)

// LabelTreeParameters has parameters of label trees.
type LabelTreeParameters struct {
	// AssignerName is the used LeftRightAssigner name.
	AssignerName string
	// AssignInitializerName is the used LeftRightAssignInitializer name.
	AssignInitializerName string
	// ClassifierTrainerName is the used BinaryClassifierTrainer name.
	ClassifierTrainerName string
	// C is the inverse of the penalty parameter used by BinaryClassifierTrainer.
	C float32
	// Epsilon is the tolerance parameter used by BinaryClassifierTrainer.
	Epsilon float32
	// FeatureSubSamplerName is the used DatasetFeatureSubSampler name.
	FeatureSubSamplerName string
	// K is the maximum number of labels in the distribution in each terminal leaf.
	K uint
	// MaxEntriesInLeaf is the maximum number of entries in each terminal leaf.
	MaxEntriesInLeaf uint
	// SuppVecK is the maximum number of support vectors in summary of LabelTree.
	SuppVecK uint
}

// NewLabelTreeParameters returns a new LabelTreeParameters with default values.
func NewLabelTreeParameters() *LabelTreeParameters {
	return &LabelTreeParameters{
		AssignerName:          DefaultLeftRightAssignerName,
		AssignInitializerName: DefaultLeftRightAssignInitializerName,
		ClassifierTrainerName: "L1SVC_PrimalSGD",
		C:                     1.0,
		Epsilon:               0.01,
		FeatureSubSamplerName: DefaultDatasetFeatureSubSamplerName,
		K:                20,
		MaxEntriesInLeaf: 100,
		SuppVecK:         10,
	}
}

// LabelTree is the data structure for trees in LabelForest.
// LabelTree can have at most 2^64 - 1 leaves.
type LabelTree struct {
	// SplitterSet is the map from a leaf id to the splitter used in the leaf.
	// If the splitter of the leaf is not nil, the splitter used for deciding whether x goes to the left or the right.
	// If it is nil, the leaf is terminal.
	SplitterSet map[uint64]*sticker.BinaryClassifier
	// LabelFreqSet is the map from a leaf id to the label frequency table in the leaf.
	// The table is constructed from the training dataset.
	// In the terminal leaf, it is used for prediction.
	LabelFreqSet map[uint64]sticker.SparseVector
	// The following members are not required.
	//
	// SummarySet is the map from a leaf id to the summary for the non-terminal leaf.
	// The entries in this summary is considered to provide compact and useful information in best-effort, so this specification would be loose and rapidly changing.
	SummarySet map[uint64]map[string]interface{}
}

// TrainLabelTree returns a trained LabelTree on the given dataset.
// The 16 MSBs of seed are used as the tree id which is reported in the debug log.
//
// This function returns an error if the height of the tree is greater than 64 or in training the tree.
func TrainLabelTree(ds *sticker.Dataset, params *LabelTreeParameters, seed int64, debug *log.Logger) (*LabelTree, error) {
	leftRightAssigner, ok := LeftRightAssigners[params.AssignerName]
	if !ok {
		return nil, fmt.Errorf("unknown LeftRightAssigner: %s", params.AssignerName)
	}
	leftRightAssignInitializer, ok := LeftRightAssignInitializers[params.AssignInitializerName]
	if !ok {
		return nil, fmt.Errorf("unknown LeftRightAssignInitializer: %s", params.AssignInitializerName)
	}
	binaryClassifierTrainer, ok := sticker.BinaryClassifierTrainers[params.ClassifierTrainerName]
	if !ok {
		return nil, fmt.Errorf("unknown BinaryClassiferTrainer: %s", params.ClassifierTrainerName)
	}
	featureSubSampler, ok := DatasetFeatureSubSamplers[params.FeatureSubSamplerName]
	if !ok {
		return nil, fmt.Errorf("unknown DatasetFeatureSubSampler: %s", params.FeatureSubSamplerName)
	}
	if debug != nil {
		debug.Printf("TrainLabelTree(seed>>48=%d): starting with sizeOfDataset=%d, params=%#v ...", seed>>48, ds.Size(), params)
	}
	rng := rand.New(rand.NewSource(seed))
	tree := &LabelTree{
		SplitterSet:  make(map[uint64]*sticker.BinaryClassifier),
		LabelFreqSet: make(map[uint64]sticker.SparseVector),
		SummarySet:   make(map[uint64]map[string]interface{}),
	}
	// Train the tree in the depth-first way.
	stackId := []uint64{0x1}
	stackDs := []*sticker.Dataset{ds}
	for len(stackId) > 0 {
		leafId, subds := stackId[len(stackId)-1], stackDs[len(stackDs)-1]
		stackId, stackDs = stackId[:len(stackId)-1], stackDs[:len(stackDs)-1]
		// Set the label frequency table of the leaf.
		labelFreq := make(sticker.SparseVector)
		for _, yi := range subds.Y {
			for _, label := range yi {
				labelFreq[label]++
			}
		}
		K := params.K
		if K > uint(len(labelFreq)) {
			K = uint(len(labelFreq))
		}
		labelRankTopK := sticker.RankTopK(labelFreq, K)
		labelFreqTopK := make(sticker.SparseVector)
		for _, label := range labelRankTopK {
			labelFreqTopK[label] = labelFreq[label]
		}
		tree.LabelFreqSet[leafId] = labelFreqTopK
		if subds.X == nil {
			continue
		}
		// Sub-sample features.
		subsubds, err := featureSubSampler(subds, seed+int64(leafId))
		if err != nil {
			return nil, fmt.Errorf("DatasetFeatureSubSampler(%s): %s", params.FeatureSubSamplerName, err)
		}
		// Optimize the left/right allocation.
		if debug != nil {
			debug.Printf("TrainLabelTree(seed>>48=%d,leafId=0b%b): optimizing the left/right allocation on sub-dataset (size=%d) ...", seed>>48, leafId, subsubds.Size())
		}
		delta := leftRightAssignInitializer(subsubds, params, rng, debug)
		if err := leftRightAssigner(subsubds, delta, nil); err != nil {
			return nil, fmt.Errorf("LeftRightAssigner(%s): %s", params.AssignerName, err)
		}
		// Optimize the hyper-plane for classifying all entries in left and right.
		nLeftRights := make(map[bool]int)
		for _, deltai := range delta {
			nLeftRights[deltai]++
		}
		if nLeftRights[false] == 0 || nLeftRights[true] == 0 {
			// There is no need to split the dataset.
			continue
		}
		if debug != nil {
			debug.Printf("TrainLabelTree(seed>>48=%d,leafId=0b%b): training the splitter: %d in left and %d in right ...", seed>>48, leafId, nLeftRights[false], nLeftRights[true])
		}
		splitter, err := binaryClassifierTrainer(subsubds.X, delta, params.C, params.Epsilon, nil)
		if err != nil {
			if debug != nil {
				debug.Printf("TrainLabelTree(seed>>48=%d,leafId=0b%b): BinaryClassifierTrainer(%s): %s", seed>>48, leafId, params.ClassifierTrainerName, err)
			}
		}
		// Divide the sub-dataset into the left/right one with the trained splitter.
		tn, fn, fp, tp, Z, predDelta := splitter.ReportPerformance(subds.X, delta)
		nPredLeftRights := make(map[bool]int)
		for _, predDeltai := range predDelta {
			nPredLeftRights[predDeltai]++
		}
		if nPredLeftRights[false] == 0 || nPredLeftRights[true] == 0 {
			// There is no need to split the dataset.
			continue
		}
		tree.SplitterSet[leafId] = splitter
		leftSubds, rightSubds := &sticker.Dataset{
			Y: make(sticker.LabelVectors, 0, nPredLeftRights[true]),
		}, &sticker.Dataset{
			Y: make(sticker.LabelVectors, 0, nPredLeftRights[false]),
		}
		leftLabelFreq, rightLabelFreq := make(sticker.SparseVector), make(sticker.SparseVector)
		for i, predDeltai := range predDelta {
			yi := subds.Y[i]
			var labelFreqi sticker.SparseVector
			if predDeltai {
				rightSubds.Y = append(rightSubds.Y, yi)
				labelFreqi = rightLabelFreq
			} else {
				leftSubds.Y = append(leftSubds.Y, yi)
				labelFreqi = leftLabelFreq
			}
			for _, label := range yi {
				labelFreqi[label]++
			}
		}
		leftLabelRankTopK, rightLabelRankTopK := sticker.RankTopK(leftLabelFreq, params.K), sticker.RankTopK(rightLabelFreq, params.K)
		leftLabelTopK, rightLabelTopK := sticker.InvertRanks(leftLabelRankTopK), sticker.InvertRanks(rightLabelRankTopK)
		// Create the training summary
		if debug != nil {
			debug.Printf("TrainLabelTree(seed>>48=%d,leafId=0b%b): creating the training summary (TN=%d,FN=%d,FP=%d,TP=%d) ...", seed>>48, leafId, tn, fn, fp, tp)
		}
		summary := make(map[string]interface{})
		tree.SummarySet[leafId] = summary
		// Summarize the training performance of the splitter.
		splitPerf := make(map[string]interface{})
		splitPerf["tn"], splitPerf["fn"], splitPerf["fp"], splitPerf["tp"] = int(tn), int(fn), int(fp), int(tp)
		splitPerfSumZPerLabel, splitPerfNentriesPerLabel := make(map[uint32]float32), make(map[uint32]int)
		for i, zi := range Z {
			for _, label := range subds.Y[i] {
				_, leftok := leftLabelTopK[label]
				_, rightok := rightLabelTopK[label]
				if leftok || rightok {
					splitPerfSumZPerLabel[label] += zi
					splitPerfNentriesPerLabel[label]++
				}
			}
		}
		splitPerf["sumZPerLabel"], splitPerf["nentriesPerLabel"] = splitPerfSumZPerLabel, splitPerfNentriesPerLabel
		summary["splitPerf"] = splitPerf
		// Summarize the top-SuppVecK support vectors of the splitter if the splitter is trained by the solver using dual problems.
		if splitter.Beta != nil {
			suppVecIdBetas := make(sticker.KeyValues32OrderedByValue, len(subds.X))
			for i := range subds.X {
				suppVecIdBetas[i] = sticker.KeyValue32{uint32(i), splitter.Beta[i]}
			}
			sort.Sort(sort.Reverse(suppVecIdBetas))
			K := params.SuppVecK
			if K > uint(len(suppVecIdBetas)) {
				K = uint(len(suppVecIdBetas))
			}
			suppVecs := make([]interface{}, K)
			for rank := uint(0); rank < K; rank++ {
				i := suppVecIdBetas[rank].Key
				xi, yi := subds.X[i], subds.Y[i]
				suppVec := make(map[string]interface{})
				suppVec["beta"] = suppVecIdBetas[rank].Value
				labels := make([]int, 0, len(yi))
				for _, label := range yi {
					labels = append(labels, int(label))
				}
				sort.Ints(labels)
				suppVec["featureVector"], suppVec["labels"], suppVec["delta"] = xi, labels, delta[i]
				suppVecs[rank] = suppVec
			}
			summary["suppVecs"] = suppVecs
			// Discard Beta of the splitter.
			splitter.Beta = nil
		}
		// Training the left and right.
		if leafId>>63 == 1 {
			return nil, fmt.Errorf("height of tree cannot be greater than 64")
		}
		// In order to force tree balances, the BOTH child leaf should have at least MaxEntriesInLeaf entries.
		if nPredLeftRights[false] >= int(params.MaxEntriesInLeaf) && nPredLeftRights[true] >= int(params.MaxEntriesInLeaf) {
			rightSubds.X = make(sticker.FeatureVectors, 0, nPredLeftRights[true])
			for i, predDeltai := range predDelta {
				if predDeltai {
					rightSubds.X = append(rightSubds.X, subds.X[i])
				}
			}
			leftSubds.X = make(sticker.FeatureVectors, 0, nPredLeftRights[false])
			for i, predDeltai := range predDelta {
				if !predDeltai {
					leftSubds.X = append(leftSubds.X, subds.X[i])
				}
			}
		}
		stackId = append(stackId, 2*leafId+1)
		stackDs = append(stackDs, rightSubds)
		stackId = append(stackId, 2*leafId+0)
		stackDs = append(stackDs, leftSubds)
	}
	if debug != nil {
		debug.Printf("TrainLabelTree(seed>>48=%d): finished training", seed>>48)
	}
	return tree, nil
}

// DecodeLabelTreeWithGobDecoder decodes LabelTree using decoder.
//
// This function returns an error in decoding.
func DecodeLabelTreeWithGobDecoder(tree *LabelTree, decoder *gob.Decoder) error {
	tree.SplitterSet = make(map[uint64]*sticker.BinaryClassifier)
	tree.LabelFreqSet = make(map[uint64]sticker.SparseVector)
	tree.SummarySet = make(map[uint64]map[string]interface{})
	stackId := []uint64{0x1}
	for len(stackId) > 0 {
		leafId := stackId[len(stackId)-1]
		stackId = stackId[:len(stackId)-1]
		var has [5]bool
		if err := decoder.Decode(&has); err != nil {
			return fmt.Errorf("DecodeLabelTree: leafId=0b%b: header: %s", leafId, err)
		}
		if has[0] {
			var splitter sticker.BinaryClassifier
			if err := decoder.Decode(&splitter); err != nil {
				return fmt.Errorf("DecodeLabelTree: leafId=0b%b: Splitter: %s", leafId, err)
			}
			tree.SplitterSet[leafId] = &splitter
		}
		if has[1] {
			var labelFreq sticker.SparseVector
			if err := decoder.Decode(&labelFreq); err != nil {
				return fmt.Errorf("DecodeLabelTree: leafId=0b%b: LabelFreq: %s", leafId, err)
			}
			tree.LabelFreqSet[leafId] = labelFreq
		}
		if has[2] {
			var summary map[string]interface{}
			if err := decoder.Decode(&summary); err != nil {
				return fmt.Errorf("DecodeLabelTree: leafId=0b%b: Summary: %s", leafId, err)
			}
			tree.SummarySet[leafId] = summary
		}
		if leafId>>63 == 1 {
			return fmt.Errorf("height of tree cannot be greater than 64")
		}
		if has[3] {
			stackId = append(stackId, 2*leafId+0)
		}
		if has[4] {
			stackId = append(stackId, 2*leafId+1)
		}
	}
	return nil
}

// DecodeLabelTree decodes LabelTree from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelTreeWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelTree(tree *LabelTree, r io.Reader) error {
	return DecodeLabelTreeWithGobDecoder(tree, gob.NewDecoder(r))
}

// EncodeLabelTreeWithGobEncoder decodes LabelTree using encoder.
//
// This function returns an error in decoding.
func EncodeLabelTreeWithGobEncoder(tree *LabelTree, encoder *gob.Encoder) error {
	stackId := []uint64{0x1}
	for len(stackId) > 0 {
		leafId := stackId[len(stackId)-1]
		stackId = stackId[:len(stackId)-1]
		leftLeafId, rightLeafId := 2*leafId+0, 2*leafId+1
		hasLeft, hasRight := tree.LabelFreqSet[leftLeafId] != nil, tree.LabelFreqSet[rightLeafId] != nil
		has := [5]bool{tree.SplitterSet[leafId] != nil, tree.LabelFreqSet[leafId] != nil, tree.SummarySet[leafId] != nil, hasLeft, hasRight}
		if err := encoder.Encode(has); err != nil {
			return fmt.Errorf("EncodeLabelTree: leafId=0b%b: header: %s", leafId, err)
		}
		if has[0] {
			if err := encoder.Encode(tree.SplitterSet[leafId]); err != nil {
				return fmt.Errorf("EncodeLabelTree: leafId=0b%b: Splitter: %s", leafId, err)
			}
		}
		if has[1] {
			if err := encoder.Encode(tree.LabelFreqSet[leafId]); err != nil {
				return fmt.Errorf("EncodeLabelTree: leafId=0b%b: LabelFreq: %s", leafId, err)
			}
		}
		if has[2] {
			if err := encoder.Encode(tree.SummarySet[leafId]); err != nil {
				return fmt.Errorf("EncodeLabelTree: leafId=0b%b: Summary: %s", leafId, err)
			}
		}
		if hasLeft {
			stackId = append(stackId, leftLeafId)
		}
		if hasRight {
			stackId = append(stackId, rightLeafId)
		}
	}
	return nil
}

// EncodeLabelTree encodes LabelTree to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelTreeWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelTree(tree *LabelTree, w io.Writer) error {
	return EncodeLabelTreeWithGobEncoder(tree, gob.NewEncoder(w))
}

// Classify returns the leaf id which x falls.
func (tree *LabelTree) Classify(x sticker.FeatureVector) uint64 {
	leafId := uint64(0x1)
	for {
		splitter := tree.SplitterSet[leafId]
		if splitter == nil {
			break
		}
		if sticker.ClassifyToBinaryClass(splitter.Predict(x)) {
			leafId = 2*leafId + 1
		} else {
			leafId = 2*leafId + 0
		}
	}
	return leafId
}

// ClassifyWithWeight returns the leaf ID and the weight which x falls.
// Weight will not affect any prediction result on single trees, it affects on ensembled trees.
func (tree *LabelTree) ClassifyWithWeight(x sticker.FeatureVector) (uint64, float32) {
	leafId, minWeight := uint64(0x1), sticker.Inf32(+1.0)
	for {
		splitter := tree.SplitterSet[leafId]
		if splitter == nil {
			break
		}
		z, count := splitter.PredictAndCount(x)
		if sticker.ClassifyToBinaryClass(z) {
			leafId = 2*leafId + 1
		} else {
			leafId = 2*leafId + 0
		}
		if minWeight > float32(count) {
			minWeight = float32(count)
		}
	}
	return leafId, minWeight
}

// ClassifyAll returns the leaf ID slice and the weight slice which each entry of X falls.
func (tree *LabelTree) ClassifyAll(X sticker.FeatureVectors) []uint64 {
	leafIds := make([]uint64, len(X))
	for i, x := range X {
		leafIds[i] = tree.Classify(x)
	}
	return leafIds
}

// ClassifyAllWithWeight returns the leaf ID slice and the weight slice which each entry of X falls.
func (tree *LabelTree) ClassifyAllWithWeight(X sticker.FeatureVectors) ([]uint64, []float32) {
	leafIds, weights := make([]uint64, len(X)), make([]float32, len(X))
	for i, x := range X {
		leafIds[i], weights[i] = tree.ClassifyWithWeight(x)
	}
	return leafIds, weights
}

// GobEncode returns the error always such that users should encode large LabelTree objects with EncodeLabelTree.
func (tree *LabelTree) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelTree should be encoded with EncodeLabelTree")
}

// IsValidLeaf returns true if the leaf id is valid, otherwise false.
func (tree *LabelTree) IsValidLeaf(leafId uint64) bool {
	return tree.LabelFreqSet[leafId] != nil
}

// IsTerminalLeaf returns true if the leaf is terminal, otherwise false.
func (tree *LabelTree) IsTerminalLeaf(leafId uint64) bool {
	return tree.SplitterSet[leafId] == nil
}

// Predict returns the top-K labels for the given result of Classify.
func (tree *LabelTree) Predict(leafId uint64, K uint) sticker.LabelVector {
	labelDist := make(map[uint32]float32)
	labelFreq := tree.LabelFreqSet[leafId]
	Z := float32(0.0)
	for _, freq := range labelFreq {
		Z += freq
	}
	for label, freq := range labelFreq {
		labelDist[label] += freq / Z
	}
	return sticker.RankTopK(labelDist, K)
}

// PredictAll returns the top-K labels for the given result of ClassifyAll.
func (tree *LabelTree) PredictAll(leafIdSlice []uint64, K uint) sticker.LabelVectors {
	YK := make(sticker.LabelVectors, len(leafIdSlice))
	for i, leafId := range leafIdSlice {
		YK[i] = tree.Predict(leafId, K)
	}
	return YK
}

// LabelForest is variously-modified FastXML (Prabhu+ 2014).
//
// References:
//
// (Prabhu+ 2014) Y. Prabhu, and M. Varma. "FastXML: A Fast, Accurate and Stable Tree-Classifier for Extreme Multi-Label Learning." Proceedings of the 20th ACM SIGKDD International Conference on Knowledge Discovery and Data Mining, pp. 263--272, 2014.
type LabelForest struct {
	// TreeParams is the parameters for training each LabelTree.
	TreeParams *LabelTreeParameters
	// Trees is the slice of trained trees.
	Trees []*LabelTree
	// The following members are not required.
	//
	// SummaryS is the sub-sampling summary.
	// This summary is considered to provide compact and useful information in best-effort, so this specification would be loose and rapidly changing.
	Summary map[string]interface{}
}

// TrainLabelForest returns a trained LabelForest on ds with multiple go-routines.
// The number of go-routines is runtime.GOMAXPROCS.
//
// This function returns the last error in training each tree in the forest by multiple go-routines.
func TrainLabelForest(ds *sticker.Dataset, ntrees uint, subSampler DatasetEntrySubSampler, params *LabelTreeParameters, debug *log.Logger) (*LabelForest, error) {
	forest := &LabelForest{
		TreeParams: params,
		Trees:      make([]*LabelTree, ntrees),
		Summary:    make(map[string]interface{}),
	}
	summaryDataCounts, summaryFeatureCounts, summaryLabelCounts := make(map[int]int), make(map[uint32]int), make(map[uint32]int)
	nworkers := runtime.GOMAXPROCS(0)
	if debug != nil {
		debug.Printf("training %d tree(s) with %d workers ...", ntrees, nworkers)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, nworkers)
	var mutexSummaries sync.Mutex
	var mutexLasterr sync.RWMutex
	var lasterr error
	latestTreeId := ^uint64(0)
	for t := uint(0); t < ntrees; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			treeId := atomic.AddUint64(&latestTreeId, 1)
			mutexLasterr.RLock()
			if lasterr != nil {
				mutexLasterr.RUnlock()
				debug.Printf("skipped to train #%d tree because of previous errors", treeId)
				return
			}
			mutexLasterr.RUnlock()
			if debug != nil {
				debug.Printf("training #%d tree ...", treeId)
			}
			indices := subSampler.SubSample(ds, uint(treeId))
			// Writing to the summaries starts here (protected by mutexSummaries).
			mutexSummaries.Lock()
			for _, i := range indices {
				summaryDataCounts[i]++
			}
			subds := ds.SubSet(indices)
			for _, xi := range subds.X {
				for _, xipair := range xi {
					summaryFeatureCounts[xipair.Key]++
				}
			}
			for _, yi := range subds.Y {
				for _, label := range yi {
					summaryLabelCounts[label]++
				}
			}
			// Writing to the summaries ends here (protected by mutexSummaries).
			mutexSummaries.Unlock()
			var err error
			forest.Trees[treeId], err = TrainLabelTree(subds, params, int64(treeId<<48), debug)
			if err != nil {
				err := fmt.Errorf("training #%d tree: %s", treeId, err)
				if debug != nil {
					debug.Printf("%s", err)
				}
				mutexLasterr.Lock()
				lasterr = err
				mutexLasterr.Unlock()
			}
			if debug != nil {
				debug.Printf("trained #%d tree", treeId)
			}
		}()
	}
	wg.Wait()
	if lasterr != nil {
		return nil, lasterr
	}
	if debug != nil {
		debug.Printf("creating the sub-sampling summary ...")
	}
	summaryDataHist, summaryFeatureHist, summaryLabelHist := make(map[int]int), make(map[int]int), make(map[int]int)
	for _, count := range summaryDataCounts {
		summaryDataHist[count]++
	}
	summaryDataHist[0] = ds.Size() - len(summaryDataCounts)
	for _, count := range summaryFeatureCounts {
		summaryFeatureHist[count]++
	}
	summaryFeatureHist[0] = ds.X.Dim() - len(summaryFeatureCounts)
	for _, count := range summaryLabelCounts {
		summaryLabelHist[count]++
	}
	summaryLabelHist[0] = ds.Y.Dim() - len(summaryLabelCounts)
	forest.Summary["dataHist"], forest.Summary["featureHist"], forest.Summary["labelHist"] = summaryDataHist, summaryFeatureHist, summaryLabelHist
	return forest, nil
}

// DecodeLabelForestWithGobDecoder decodes LabelForest using decoder.
//
// This function returns an error in decoding.
func DecodeLabelForestWithGobDecoder(forest *LabelForest, decoder *gob.Decoder) error {
	forest.TreeParams = &LabelTreeParameters{}
	if err := decoder.Decode(&forest.TreeParams); err != nil {
		return fmt.Errorf("DecodeLabelForest: TreeParams: %s", err)
	}
	var ntrees int
	if err := decoder.Decode(&ntrees); err != nil {
		return fmt.Errorf("DecodeLabelForest: ntrees: %s", err)
	}
	forest.Trees = make([]*LabelTree, ntrees)
	for treeId := 0; treeId < ntrees; treeId++ {
		forest.Trees[treeId] = &LabelTree{}
		if err := DecodeLabelTreeWithGobDecoder(forest.Trees[treeId], decoder); err != nil {
			return fmt.Errorf("DecodeLabelForest: #%d tree: %s", treeId, err)
		}
	}
	if err := decoder.Decode(&forest.Summary); err != nil {
		return fmt.Errorf("DecodeLabelForest: Summary: %s", err)
	}
	return nil
}

// DecodeLabelForest decodes LabelForest from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelForestWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelForest(forest *LabelForest, r io.Reader) error {
	return DecodeLabelForestWithGobDecoder(forest, gob.NewDecoder(r))
}

// EncodeLabelForestWithGobEncoder decodes LabelForest using encoder.
//
// This function returns an error in decoding.
func EncodeLabelForestWithGobEncoder(forest *LabelForest, encoder *gob.Encoder) error {
	if err := encoder.Encode(forest.TreeParams); err != nil {
		return fmt.Errorf("EncodeLabelForest: TreeParams: %s", err)
	}
	if err := encoder.Encode(len(forest.Trees)); err != nil {
		return fmt.Errorf("EncodeLabelForest: ntrees: %s", err)
	}
	for treeId, tree := range forest.Trees {
		if err := EncodeLabelTreeWithGobEncoder(tree, encoder); err != nil {
			return fmt.Errorf("EncodeLabelForest: #%d tree: %s", treeId, err)
		}
	}
	if err := encoder.Encode(forest.Summary); err != nil {
		return fmt.Errorf("EncodeLabelForest: Summary: %s", err)
	}
	return nil
}

// EncodeLabelForest encodes LabelForest to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelForestWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelForest(forest *LabelForest, w io.Writer) error {
	return EncodeLabelForestWithGobEncoder(forest, gob.NewEncoder(w))
}

// Classify returns the leaf id slice for the given feature vector.
func (forest *LabelForest) Classify(x sticker.FeatureVector) []uint64 {
	leafIds := make([]uint64, len(forest.Trees))
	for treeId, tree := range forest.Trees {
		leafIds[treeId] = tree.Classify(x)
	}
	return leafIds
}

// ClassifyWithWeight returns the leaf id slice and the weight slice for the given feature vector.
func (forest *LabelForest) ClassifyWithWeight(x sticker.FeatureVector) ([]uint64, []float32) {
	leafIds, weights := make([]uint64, len(forest.Trees)), make([]float32, len(forest.Trees))
	for treeId, tree := range forest.Trees {
		leafIds[treeId], weights[treeId] = tree.ClassifyWithWeight(x)
	}
	return leafIds, weights
}

// ClassifyAll returns the slice of the leaf id slices for each feature vector.
func (forest *LabelForest) ClassifyAll(X sticker.FeatureVectors) [][]uint64 {
	leafIdsSlice := make([][]uint64, len(X))
	for i, x := range X {
		leafIdsSlice[i] = forest.Classify(x)
	}
	return leafIdsSlice
}

// ClassifyAllWithWeight returns the slice of the leaf id slices and the weight slices for each feature vector.
func (forest *LabelForest) ClassifyAllWithWeight(X sticker.FeatureVectors) ([][]uint64, [][]float32) {
	leafIdsSlice, weightsSlice := make([][]uint64, len(X)), make([][]float32, len(X))
	for i, x := range X {
		leafIdsSlice[i], weightsSlice[i] = forest.ClassifyWithWeight(x)
	}
	return leafIdsSlice, weightsSlice
}

// GobEncode returns the error always, because users should encode large LabelForest objects with EncodeLabelForest.
func (forest *LabelForest) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelForest should be encoded with EncodeLabelForest")
}

// Predict returns the top-K labels for the given result of Classify.
func (forest *LabelForest) Predict(leafIds []uint64, K uint) sticker.LabelVector {
	labelDist := make(map[uint32]float32)
	for treeId, tree := range forest.Trees {
		labelFreq := tree.LabelFreqSet[leafIds[treeId]]
		Z := float32(0.0)
		for _, freq := range labelFreq {
			Z += freq
		}
		for label, freq := range labelFreq {
			labelDist[label] += freq / Z
		}
	}
	return sticker.RankTopK(labelDist, K)
}

// PredictWithWeight returns the top-K labels for the given result of ClassifyWithWeight.
func (forest *LabelForest) PredictWithWeight(leafIds []uint64, weights []float32, K uint) sticker.LabelVector {
	labelDist := make(map[uint32]float32)
	for treeId, tree := range forest.Trees {
		labelFreq := tree.LabelFreqSet[leafIds[treeId]]
		Z := float32(0.0)
		for _, freq := range labelFreq {
			Z += freq
		}
		for label, freq := range labelFreq {
			labelDist[label] += freq / Z * weights[treeId]
		}
	}
	return sticker.RankTopK(labelDist, K)
}

// PredictAll returns the top-K labels for the given result of ClassifyAll.
func (forest *LabelForest) PredictAll(leafIdsSlice [][]uint64, K uint) sticker.LabelVectors {
	YK := make(sticker.LabelVectors, len(leafIdsSlice))
	for i, leafIds := range leafIdsSlice {
		YK[i] = forest.Predict(leafIds, K)
	}
	return YK
}

// PredictAllWithWeight returns the top-K labels for the given result of ClassifyAllWithWeight.
func (forest *LabelForest) PredictAllWithWeight(leafIdsSlice [][]uint64, weightsSlice [][]float32, K uint) sticker.LabelVectors {
	YK := make(sticker.LabelVectors, len(leafIdsSlice))
	for i, leafIds := range leafIdsSlice {
		YK[i] = forest.PredictWithWeight(leafIds, weightsSlice[i], K)
	}
	return YK
}
