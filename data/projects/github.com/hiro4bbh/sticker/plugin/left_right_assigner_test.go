package plugin

import (
	"bytes"
	"log"
	"math/rand"
	"testing"

	"github.com/hiro4bbh/go-assert"
	"github.com/hiro4bbh/sticker"
)

func assertBinaryClassAssignmentEqual(t *testing.T, expected, got []bool) bool {
	t.Helper()
	if expected[0] == got[0] {
		goassert.New(t, expected).Equal(got)
		return false
	}
	expectedNeg := make([]bool, len(expected))
	for i, b := range expected {
		expectedNeg[i] = !b
	}
	goassert.New(t, expectedNeg).Equal(got)
	return true
}

func TestLeftRightAssigner_greedyBottomRanks(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	// Test fully-separable case: mutable labels
	n1 := 25
	ds1, z1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n1),
		Y: make(sticker.LabelVectors, 2*n1),
	}, make([]bool, 2*n1)
	for i := 0; i < n1; i++ {
		ds1.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.Y[2*i+0] = sticker.LabelVector{0, 2}
		ds1.Y[2*i+1] = sticker.LabelVector{1, 3}
		z1[2*i+1] = true
	}
	delta1 := make([]bool, len(ds1.X))
	for i := n1; i < len(delta1); i++ {
		delta1[i] = true
	}
	LeftRightAssigner_greedyBottomRanks(ds1, delta1, nil)
	assertBinaryClassAssignmentEqual(t, z1, delta1)
	// Test fully-separable case: non-mutable labels
	n2 := 25
	ds2, z2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n2),
		Y: make(sticker.LabelVectors, 2*n2),
	}, make([]bool, 2*n2)
	for i := 0; i < n2; i++ {
		ds2.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.Y[2*i+0] = sticker.LabelVector{0, 1}
		ds2.Y[2*i+1] = sticker.LabelVector{0, 2}
		z2[2*i+1] = true
	}
	delta2 := make([]bool, len(ds2.X))
	for i := n2; i < len(delta2); i++ {
		delta2[i] = true
	}
	LeftRightAssigner_greedyBottomRanks(ds2, delta2, nil)
	assertBinaryClassAssignmentEqual(t, z2, delta2)
	// Test fully-separable case: non-mutable labels and sub-level labels
	n3 := 25
	ds3, z3 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 4*n3),
		Y: make(sticker.LabelVectors, 4*n3),
	}, make([]bool, 4*n3)
	for i := 0; i < n3; i++ {
		ds3.X[4*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+2] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+3] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.Y[4*i+0] = sticker.LabelVector{0, 1, 3}
		ds3.Y[4*i+1] = sticker.LabelVector{0, 1, 4}
		ds3.Y[4*i+2] = sticker.LabelVector{0, 2, 5}
		ds3.Y[4*i+3] = sticker.LabelVector{0, 2, 6}
		z3[4*i+2] = true
		z3[4*i+3] = true
	}
	delta3 := make([]bool, len(ds3.X))
	for i := 2 * n3; i < len(delta3); i++ {
		delta3[i] = true
	}
	LeftRightAssigner_greedyBottomRanks(ds3, delta3, nil)
	assertBinaryClassAssignmentEqual(t, z3, delta3)
	ds3_1, z3_1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_1.X[2*i+0] = ds3.X[4*i+0]
		ds3_1.Y[2*i+0] = ds3.Y[4*i+0]
		ds3_1.X[2*i+1] = ds3.X[4*i+1]
		ds3_1.Y[2*i+1] = ds3.Y[4*i+1]
		z3_1[2*i+1] = true
	}
	delta3_1 := make([]bool, len(ds3_1.X))
	for i := n3; i < len(delta3_1); i++ {
		delta3_1[i] = true
	}
	LeftRightAssigner_greedyBottomRanks(ds3_1, delta3_1, nil)
	assertBinaryClassAssignmentEqual(t, z3_1, delta3_1)
	ds3_2, z3_2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_2.X[2*i+0] = ds3.X[4*i+2]
		ds3_2.Y[2*i+0] = ds3.Y[4*i+2]
		ds3_2.X[2*i+1] = ds3.X[4*i+3]
		ds3_2.Y[2*i+1] = ds3.Y[4*i+3]
		z3_2[2*i+1] = true
	}
	delta3_2 := make([]bool, len(ds3_2.X))
	for i := n3; i < len(delta3_2); i++ {
		delta3_2[i] = true
	}
	LeftRightAssigner_greedyBottomRanks(ds3_2, delta3_2, nil)
	assertBinaryClassAssignmentEqual(t, z3_2, delta3_2)
	// Check debug logs
	var debugBuffer bytes.Buffer
	LeftRightAssigner_greedyBottomRanks(ds2, delta2, log.New(&debugBuffer, "", 0))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func TestLeftRightAssigner_nDCG(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	// Test fully-separable case: mutable labels
	n1 := 25
	ds1, z1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n1),
		Y: make(sticker.LabelVectors, 2*n1),
	}, make([]bool, 2*n1)
	for i := 0; i < n1; i++ {
		ds1.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.Y[2*i+0] = sticker.LabelVector{0, 2}
		ds1.Y[2*i+1] = sticker.LabelVector{1, 3}
		z1[2*i+1] = true
	}
	delta1 := make([]bool, len(ds1.X))
	for i := n1; i < len(delta1); i++ {
		delta1[i] = true
	}
	LeftRightAssigner_nDCG(ds1, delta1, nil)
	assertBinaryClassAssignmentEqual(t, z1, delta1)
	// Test fully-separable case: non-mutable labels
	n2 := 25
	ds2, z2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n2),
		Y: make(sticker.LabelVectors, 2*n2),
	}, make([]bool, 2*n2)
	for i := 0; i < n2; i++ {
		ds2.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.Y[2*i+0] = sticker.LabelVector{0, 1}
		ds2.Y[2*i+1] = sticker.LabelVector{0, 2}
		z2[2*i+1] = true
	}
	delta2 := make([]bool, len(ds2.X))
	for i := n2; i < len(delta2); i++ {
		delta2[i] = true
	}
	LeftRightAssigner_nDCG(ds2, delta2, nil)
	assertBinaryClassAssignmentEqual(t, z2, delta2)
	// Test fully-separable case: non-mutable labels and sub-level labels
	n3 := 25
	ds3, z3 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 4*n3),
		Y: make(sticker.LabelVectors, 4*n3),
	}, make([]bool, 4*n3)
	for i := 0; i < n3; i++ {
		ds3.X[4*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+2] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+3] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.Y[4*i+0] = sticker.LabelVector{0, 1, 3}
		ds3.Y[4*i+1] = sticker.LabelVector{0, 1, 4}
		ds3.Y[4*i+2] = sticker.LabelVector{0, 2, 5}
		ds3.Y[4*i+3] = sticker.LabelVector{0, 2, 6}
		z3[4*i+2] = true
		z3[4*i+3] = true
	}
	delta3 := make([]bool, len(ds3.X))
	for i := 2 * n3; i < len(delta3); i++ {
		delta3[i] = true
	}
	LeftRightAssigner_nDCG(ds3, delta3, nil)
	assertBinaryClassAssignmentEqual(t, z3, delta3)
	ds3_1, z3_1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_1.X[2*i+0] = ds3.X[4*i+0]
		ds3_1.Y[2*i+0] = ds3.Y[4*i+0]
		ds3_1.X[2*i+1] = ds3.X[4*i+1]
		ds3_1.Y[2*i+1] = ds3.Y[4*i+1]
		z3_1[2*i+1] = true
	}
	delta3_1 := make([]bool, len(ds3_1.X))
	for i := n3; i < len(delta3_1); i++ {
		delta3_1[i] = true
	}
	LeftRightAssigner_nDCG(ds3_1, delta3_1, nil)
	assertBinaryClassAssignmentEqual(t, z3_1, delta3_1)
	ds3_2, z3_2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_2.X[2*i+0] = ds3.X[4*i+2]
		ds3_2.Y[2*i+0] = ds3.Y[4*i+2]
		ds3_2.X[2*i+1] = ds3.X[4*i+3]
		ds3_2.Y[2*i+1] = ds3.Y[4*i+3]
		z3_2[2*i+1] = true
	}
	delta3_2 := make([]bool, len(ds3_2.X))
	for i := n3; i < len(delta3_2); i++ {
		delta3_2[i] = true
	}
	LeftRightAssigner_nDCG(ds3_2, delta3_2, nil)
	assertBinaryClassAssignmentEqual(t, z3_2, delta3_2)
	// Check debug logs
	var debugBuffer bytes.Buffer
	LeftRightAssigner_nDCG(ds2, delta2, log.New(&debugBuffer, "", 0))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func TestLeftRightAssigner_none(t *testing.T) {
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{sticker.KeyValue32{0, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{1, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{0, 2.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{1}, sticker.LabelVector{0},
		},
	}
	delta := []bool{false, true, false}
	goassert.New(t).SucceedWithoutError(LeftRightAssigner_none(ds, delta, nil))
	goassert.New(t, []bool{false, true, false}).Equal(delta)
}

func testLeftRightAssignInitializer(t *testing.T, name string) {
	leftRightAssignInitializer := LeftRightAssignInitializers[name]
	rng := rand.New(rand.NewSource(0))
	// Test fully-separable case: mutable labels
	n1 := 25
	ds1, z1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n1),
		Y: make(sticker.LabelVectors, 2*n1),
	}, make([]bool, 2*n1)
	for i := 0; i < n1; i++ {
		ds1.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds1.Y[2*i+0] = sticker.LabelVector{0, 2}
		ds1.Y[2*i+1] = sticker.LabelVector{1, 3}
		z1[2*i+1] = true
	}
	assertBinaryClassAssignmentEqual(t, z1, leftRightAssignInitializer(ds1, NewLabelTreeParameters(), rand.New(rand.NewSource(0)), nil))
	// Test fully-separable case: non-mutable labels
	n2 := 25
	ds2, z2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n2),
		Y: make(sticker.LabelVectors, 2*n2),
	}, make([]bool, 2*n2)
	for i := 0; i < n2; i++ {
		ds2.X[2*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.X[2*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds2.Y[2*i+0] = sticker.LabelVector{0, 1}
		ds2.Y[2*i+1] = sticker.LabelVector{0, 2}
		z2[2*i+1] = true
	}
	assertBinaryClassAssignmentEqual(t, z2, leftRightAssignInitializer(ds2, NewLabelTreeParameters(), rand.New(rand.NewSource(0)), nil))
	// Test fully-separable case: non-mutable labels and sub-level labels
	n3 := 25
	ds3, z3 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 4*n3),
		Y: make(sticker.LabelVectors, 4*n3),
	}, make([]bool, 4*n3)
	for i := 0; i < n3; i++ {
		ds3.X[4*i+0] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+1] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+2] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.X[4*i+3] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 1.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		ds3.Y[4*i+0] = sticker.LabelVector{0, 1, 3}
		ds3.Y[4*i+1] = sticker.LabelVector{0, 1, 4}
		ds3.Y[4*i+2] = sticker.LabelVector{0, 2, 5}
		ds3.Y[4*i+3] = sticker.LabelVector{0, 2, 6}
		z3[4*i+2] = true
		z3[4*i+3] = true
	}
	assertBinaryClassAssignmentEqual(t, z3, leftRightAssignInitializer(ds3, NewLabelTreeParameters(), rand.New(rand.NewSource(0)), nil))
	ds3_1, z3_1 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_1.X[2*i+0] = ds3.X[4*i+0]
		ds3_1.Y[2*i+0] = ds3.Y[4*i+0]
		ds3_1.X[2*i+1] = ds3.X[4*i+1]
		ds3_1.Y[2*i+1] = ds3.Y[4*i+1]
		z3_1[2*i+1] = true
	}
	assertBinaryClassAssignmentEqual(t, z3_1, leftRightAssignInitializer(ds3_1, NewLabelTreeParameters(), rand.New(rand.NewSource(0)), nil))
	ds3_2, z3_2 := &sticker.Dataset{
		X: make(sticker.FeatureVectors, 2*n3),
		Y: make(sticker.LabelVectors, 2*n3),
	}, make([]bool, 2*n3)
	for i := 0; i < n3; i++ {
		ds3_2.X[2*i+0] = ds3.X[4*i+2]
		ds3_2.Y[2*i+0] = ds3.Y[4*i+2]
		ds3_2.X[2*i+1] = ds3.X[4*i+3]
		ds3_2.Y[2*i+1] = ds3.Y[4*i+3]
		z3_2[2*i+1] = true
	}
	assertBinaryClassAssignmentEqual(t, z3_2, leftRightAssignInitializer(ds3_2, NewLabelTreeParameters(), rand.New(rand.NewSource(0)), nil))
}

func TestLeftRightAssignInitializer_topLabelGraph(t *testing.T) {
	testLeftRightAssignInitializer(t, "topLabelGraph")
}

func TestLeftRightAssignInitializer_topLabelTree(t *testing.T) {
	testLeftRightAssignInitializer(t, "topLabelTree")
}
