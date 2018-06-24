package plugin

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/hiro4bbh/go-assert"
	"github.com/hiro4bbh/sticker"
)

func TestLabelTree(t *testing.T) {
	tree := &LabelTree{
		SplitterSet: map[uint64]*sticker.BinaryClassifier{
			0x1: {
				Bias:   0.0,
				Weight: sticker.SparsifyVector([]float32{1.0}),
			},
		},
		LabelFreqSet: map[uint64]sticker.SparseVector{
			0x1: {0: 1, 1: 2},
			0x2: {0: 1},
			0x3: {1: 2},
		},
	}
	goassert.New(t, false).Equal(tree.IsValidLeaf(0x0))
	goassert.New(t, true).Equal(tree.IsValidLeaf(0x1))
	goassert.New(t, true).Equal(tree.IsValidLeaf(0x2))
	goassert.New(t, true).Equal(tree.IsValidLeaf(0x3))
	goassert.New(t, false).Equal(tree.IsTerminalLeaf(0x1))
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x2))
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x3))
}

func TestLabelTreeClassify_Predict(t *testing.T) {
	tree := &LabelTree{
		SplitterSet: map[uint64]*sticker.BinaryClassifier{
			0x1: {
				Bias:   0.0,
				Weight: sticker.SparsifyVector([]float32{1.0}),
			},
		},
		LabelFreqSet: map[uint64]sticker.SparseVector{
			0x1: {0: 1, 1: 2},
			0x2: {0: 1},
			0x3: {1: 2},
		},
	}
	X := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}}, sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}, sticker.FeatureVector{sticker.KeyValue32{0, 1.0}},
	}
	leafIdSlice := tree.ClassifyAll(X)
	goassert.New(t, []uint64{0x2, 0x2, 0x3}).Equal(leafIdSlice)
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{}, sticker.LabelVector{}, sticker.LabelVector{}}).Equal(tree.PredictAll(leafIdSlice, 0))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0}, sticker.LabelVector{0}, sticker.LabelVector{1}}).Equal(tree.PredictAll(leafIdSlice, 1))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0, ^uint32(0)}, sticker.LabelVector{0, ^uint32(0)}, sticker.LabelVector{1, ^uint32(0)}}).Equal(tree.PredictAll(leafIdSlice, 2))
	leafIdSliceWithWeight, weightSlice := tree.ClassifyAllWithWeight(X)
	goassert.New(t, leafIdSlice).Equal(leafIdSliceWithWeight)
	goassert.New(t, []float32{1, 1, 1}).Equal(weightSlice)
}

func TestTrainLabelTree(t *testing.T) {
	// Test usual cases:
	//   This components should be tested already in unit tests, so this test is for integration test.
	//   This test uses only L1SVC_DualCD which is not default classifier anymore, because it uses support vectors used in creating a summary.
	// Test training:
	//   Summaries are not checked, because it is used in rapid-inspection.
	n := 100
	ds := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 4*n),
		Y: make(sticker.LabelVectors, 4*n),
	}
	delta1, delta2 := make([]bool, 4*n), make([]bool, 2*n)
	for i := 0; i < n; i++ {
		ds.X[4*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}}
		ds.X[4*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, 1.0}}
		ds.X[4*i+2] = sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}}
		ds.X[4*i+3] = sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, -1.0}}
		ds.Y[4*i+0] = sticker.LabelVector{0, 1, 3}
		ds.Y[4*i+1] = sticker.LabelVector{0, 1, 4}
		ds.Y[4*i+2] = sticker.LabelVector{0, 2, 5}
		ds.Y[4*i+3] = sticker.LabelVector{0, 2, 6}
		delta1[4*i+0], delta1[4*i+1], delta1[4*i+2], delta1[4*i+3] = false, false, true, true
		delta2[2*i+0], delta2[2*i+1] = false, true
	}
	// Set parameters as forcing the following:
	//   * cutting too many required support vectors
	//   * the minimum value ok K such that trees can grow fully (hitting the condition that it is unable to try to split in terminal leaves)
	params := NewLabelTreeParameters()
	params.SuppVecK = uint(8 * n)
	params.MaxEntriesInLeaf = uint(2 * n)
	tree := goassert.New(t).SucceedNew(TrainLabelTree(ds, params, 0, nil)).(*LabelTree)
	// Expected LabelTree structure:
	//   #b1: (0.0, 1.0)x + 0.0 <> 0
	//        #b10: (1.0, 0.0)x + 0.0 <> 0
	//        #b11: (1.0, 0.0)x + 0.0 <> 0
	goassert.New(t, false).Equal(tree.IsTerminalLeaf(0x1))
	lv1Reversed := assertBinaryClassAssignmentEqual(t, delta1, sticker.ClassifyAllToBinaryClass(tree.SplitterSet[0x1].PredictAll(ds.X)))
	leftDs, rightDs := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n),
		Y: make(sticker.LabelVectors, 2*n),
	}, &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n),
		Y: make(sticker.LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		leftDs.X[2*i+0], leftDs.X[2*i+1] = ds.X[4*i+0], ds.X[4*i+1]
		rightDs.X[2*i+0], rightDs.X[2*i+1] = ds.X[4*i+2], ds.X[4*i+3]
		leftDs.Y[2*i+0], leftDs.Y[2*i+1] = ds.Y[4*i+0], ds.Y[4*i+1]
		rightDs.Y[2*i+0], rightDs.Y[2*i+1] = ds.Y[4*i+2], ds.Y[4*i+3]
	}
	if lv1Reversed {
		leftDs, rightDs = rightDs, leftDs
	}
	goassert.New(t, false).Equal(tree.IsTerminalLeaf(0x2))
	lv20Reversed := assertBinaryClassAssignmentEqual(t, delta2, sticker.ClassifyAllToBinaryClass(tree.SplitterSet[0x2].PredictAll(leftDs.X)))
	goassert.New(t, false).Equal(tree.IsTerminalLeaf(0x3))
	lv21Reversed := assertBinaryClassAssignmentEqual(t, delta2, sticker.ClassifyAllToBinaryClass(tree.SplitterSet[0x3].PredictAll(rightDs.X)))
	leftLeftLabelFreq, leftRightLabelFreq := sticker.SparseVector{0: float32(n), 1: float32(n), 3: float32(n)}, sticker.SparseVector{0: float32(n), 1: float32(n), 4: float32(n)}
	rightLeftLabelFreq, rightRightLabelFreq := sticker.SparseVector{0: float32(n), 2: float32(n), 5: float32(n)}, sticker.SparseVector{0: float32(n), 2: float32(n), 6: float32(n)}
	if lv1Reversed {
		leftLeftLabelFreq, leftRightLabelFreq, rightLeftLabelFreq, rightRightLabelFreq = rightLeftLabelFreq, rightRightLabelFreq, leftLeftLabelFreq, leftRightLabelFreq
	}
	if lv20Reversed {
		leftLeftLabelFreq, leftRightLabelFreq = leftRightLabelFreq, leftLeftLabelFreq
	}
	if lv21Reversed {
		rightLeftLabelFreq, rightRightLabelFreq = rightRightLabelFreq, rightLeftLabelFreq
	}
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x4))
	goassert.New(t, leftLeftLabelFreq).Equal(tree.LabelFreqSet[0x4])
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x5))
	goassert.New(t, leftRightLabelFreq).Equal(tree.LabelFreqSet[0x5])
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x6))
	goassert.New(t, rightLeftLabelFreq).Equal(tree.LabelFreqSet[0x6])
	goassert.New(t, true).Equal(tree.IsTerminalLeaf(0x7))
	goassert.New(t, rightRightLabelFreq).Equal(tree.LabelFreqSet[0x7])
	// Test encoder/decoder
	var buf bytes.Buffer
	goassert.New(t, "LabelTree should be encoded with EncodeLabelTree").ExpectError(gob.NewEncoder(&buf).Encode(tree))
	buf.Reset()
	goassert.New(t).SucceedWithoutError(EncodeLabelTree(tree, &buf))
	var decodedTree LabelTree
	// gob.Decoder.Decode won't call LabelLeaf.GobDecode, because the encoder did not encode LabelLeaf.
	goassert.New(t).SucceedWithoutError(DecodeLabelTree(&decodedTree, &buf))
	goassert.New(t, tree).Equal(&decodedTree)
	// Test debug log
	var debugBuf bytes.Buffer
	goassert.New(t).SucceedNew(TrainLabelTree(ds, NewLabelTreeParameters(), 0, log.New(&debugBuf, "", 0)))
	goassert.New(t, true).Equal(debugBuf.String() != "")

	// Test the trivial case: datasets whose elements have with the same labels.
	dsConst := &sticker.Dataset{
		X: make(sticker.FeatureVectors, n),
		Y: make(sticker.LabelVectors, n),
	}
	for i := 0; i < n; i++ {
		dsConst.X[i] = sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}
		dsConst.Y[i] = sticker.LabelVector{0}
	}
	treeConst := goassert.New(t).SucceedNew(TrainLabelTree(dsConst, params, 0, nil)).(*LabelTree)
	goassert.New(t, true).Equal(treeConst.IsTerminalLeaf(0x1))
	goassert.New(t, sticker.SparseVector{0: float32(n)}).Equal(treeConst.LabelFreqSet[0x1])

	// Test the impossible split case: datasets whose elements have separable labels but those cannot be splitted by any hyper-plane.
	dsSingleton := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n),
		Y: make(sticker.LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		dsSingleton.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}
		dsSingleton.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}
		dsSingleton.Y[2*i+0] = sticker.LabelVector{0}
		dsSingleton.Y[2*i+1] = sticker.LabelVector{1}
	}
	treeSingleton := goassert.New(t).SucceedNew(TrainLabelTree(dsSingleton, params, 0, nil)).(*LabelTree)
	goassert.New(t, true).Equal(treeSingleton.IsTerminalLeaf(0x1))
	goassert.New(t, sticker.SparseVector{0: float32(n), 1: float32(n)}).Equal(treeSingleton.LabelFreqSet[0x1])
}

func TestLabelForestClassify_Predict(t *testing.T) {
	forest := &LabelForest{
		Trees: []*LabelTree{
			{
				SplitterSet: map[uint64]*sticker.BinaryClassifier{
					0x1: {
						Bias:   -1.0,
						Weight: sticker.SparsifyVector([]float32{1.0}),
					},
				},
				LabelFreqSet: map[uint64]sticker.SparseVector{
					0x1: {0: 1, 1: 2, 9: 2},
					0x2: {0: 1, 9: 1},
					0x3: {1: 2, 9: 1},
				},
			},
			{
				SplitterSet: map[uint64]*sticker.BinaryClassifier{
					0x1: {
						Bias:   0.0,
						Weight: sticker.SparsifyVector([]float32{1.0}),
					},
				},
				LabelFreqSet: map[uint64]sticker.SparseVector{
					0x1: {0: 1, 2: 2, 9: 2},
					0x2: {0: 1, 9: 1},
					0x3: {2: 2, 9: 1},
				},
			},
			{
				SplitterSet: map[uint64]*sticker.BinaryClassifier{
					0x1: {
						Bias:   1.0,
						Weight: sticker.SparsifyVector([]float32{1.0, 1.0}),
					},
				},
				LabelFreqSet: map[uint64]sticker.SparseVector{
					0x1: {0: 1, 3: 2, 9: 2},
					0x2: {0: 1, 9: 1},
					0x3: {3: 2, 9: 1},
				},
			},
		},
	}
	X := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}}, sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}, sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
	}
	leafIdsSlice := forest.ClassifyAll(X)
	goassert.New(t, [][]uint64{{0x2, 0x2, 0x2}, {0x2, 0x2, 0x3}, {0x2, 0x3, 0x3}}).Equal(leafIdsSlice)
	leafIdsSliceWithWeight, weightsSlice := forest.ClassifyAllWithWeight(X)
	goassert.New(t, leafIdsSlice).Equal(leafIdsSliceWithWeight)
	goassert.New(t, [][]float32{{1, 1, 1}, {1, 1, 1}, {1, 1, 2}}).Equal(weightsSlice)
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{}, sticker.LabelVector{}, sticker.LabelVector{}}).Equal(forest.PredictAll(leafIdsSlice, 0))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0}, sticker.LabelVector{9}, sticker.LabelVector{9}}).Equal(forest.PredictAll(leafIdsSlice, 1))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0, 9}, sticker.LabelVector{9, 0}, sticker.LabelVector{9, 2}}).Equal(forest.PredictAll(leafIdsSlice, 2))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0, 9, ^uint32(0)}, sticker.LabelVector{9, 0, 3}, sticker.LabelVector{9, 2, 3}}).Equal(forest.PredictAll(leafIdsSlice, 3))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0, 9, ^uint32(0), ^uint32(0)}, sticker.LabelVector{9, 0, 3, ^uint32(0)}, sticker.LabelVector{9, 2, 3, 0}}).Equal(forest.PredictAll(leafIdsSlice, 4))
	goassert.New(t, sticker.LabelVectors{sticker.LabelVector{0, 9, ^uint32(0), ^uint32(0)}, sticker.LabelVector{9, 0, 3, ^uint32(0)}, sticker.LabelVector{9, 3, 2, 0}}).Equal(forest.PredictAllWithWeight(leafIdsSlice, weightsSlice, 4))
}

func TestTrainLabelForest(t *testing.T) {
	// Test usual cases:
	//   This components should be tested already in unit tests, so this test is for integration test.
	// Test training
	n := 100
	ds := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n),
		Y: make(sticker.LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		ds.X[n*0+i], ds.X[n*1+i] = sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}, sticker.FeatureVector{sticker.KeyValue32{0, 0.0}}
		ds.Y[n*0+i], ds.Y[n*1+i] = sticker.LabelVector{0}, sticker.LabelVector{1}
	}
	params := NewLabelTreeParameters()
	subSampler := NewDeterministicDatasetEntrySubSampler(uint(n))
	forest := goassert.New(t).SucceedNew(TrainLabelForest(ds, 2, subSampler, params, nil)).(*LabelForest)
	goassert.New(t, params).Equal(forest.TreeParams)
	goassert.New(t, 2).Equal(len(forest.Trees))
	goassert.New(t, true).Equal(forest.Trees[0].IsTerminalLeaf(0x1))
	goassert.New(t, sticker.SparseVector{0: float32(n)}).Equal(forest.Trees[0].LabelFreqSet[0x1])
	goassert.New(t, true).Equal(forest.Trees[1].IsTerminalLeaf(0x1))
	goassert.New(t, sticker.SparseVector{1: float32(n)}).Equal(forest.Trees[1].LabelFreqSet[0x1])
	// Test encoder/decoder.
	var buf bytes.Buffer
	goassert.New(t, "LabelForest should be encoded with EncodeLabelForest").ExpectError(gob.NewEncoder(&buf).Encode(forest))
	buf.Reset()
	goassert.New(t).SucceedWithoutError(EncodeLabelForest(forest, &buf))
	var decodedForest LabelForest
	// gob.Decoder.Decode won't call LabelForest.GobDecode, because the encoder did not encode LabelForest.
	goassert.New(t).SucceedWithoutError(DecodeLabelForest(&decodedForest, &buf))
	goassert.New(t, forest).Equal(&decodedForest)
	// Test debug log
	var debugBuf bytes.Buffer
	goassert.New(t).SucceedNew(TrainLabelForest(ds, 2, subSampler, params, log.New(&debugBuf, "", 0)))
	goassert.New(t, true).Equal(debugBuf.String() != "")
}
