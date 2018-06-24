package plugin

import (
	"bytes"
	"log"
	"math/rand"
	"testing"

	"github.com/hiro4bbh/go-assert"
	"github.com/hiro4bbh/sticker"
)

func TestBinaryClassifierTrainer_L1SVC_DualCD(t *testing.T) {
	C, epsilon, debug := float32(3.0), float32(0.1), (*log.Logger)(nil)
	// Case 1: fully-separable 1x2 points
	X1, Y1 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
	}, []bool{false, true}
	bsvc1 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X1, Y1, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y1).Equal(sticker.ClassifyAllToBinaryClass(bsvc1.PredictAll(X1)))
	// Case 2: fully-separable 2x2 points
	X2, Y2 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{false, false, true, true}
	bsvc2 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X2, Y2, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y2).Equal(sticker.ClassifyAllToBinaryClass(bsvc2.PredictAll(X2)))
	// Case 3: fully-separable 1 and 3 points
	X3, Y3 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{true, false, true, true}
	bsvc3 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X3, Y3, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y3).Equal(sticker.ClassifyAllToBinaryClass(bsvc3.PredictAll(X3)))
	// Case 4: fully-separable 3 and 1 points
	X4, Y4 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{false, false, false, true}
	bsvc4 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X4, Y4, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y4).Equal(sticker.ClassifyAllToBinaryClass(bsvc4.PredictAll(X4)))
	// Expect any debug log.
	var debugBuffer bytes.Buffer
	goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X1, Y1, C, epsilon, log.New(&debugBuffer, "", 0)))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func TestBinaryClassifierTrainer_L2SVC_PrimalCD(t *testing.T) {
	C, epsilon, debug := float32(3.0), float32(0.1), (*log.Logger)(nil)
	// Case 1: fully-separable 1x2 points
	X1, Y1 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
	}, []bool{false, true}
	bsvc1 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X1, Y1, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y1).Equal(sticker.ClassifyAllToBinaryClass(bsvc1.PredictAll(X1)))
	// Case 2: fully-separable 2x2 points
	X2, Y2 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{false, false, true, true}
	bsvc2 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X2, Y2, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y2).Equal(sticker.ClassifyAllToBinaryClass(bsvc2.PredictAll(X2)))
	// Case 3: fully-separable 1 and 3 points
	X3, Y3 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{true, false, true, true}
	bsvc3 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X3, Y3, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y3).Equal(sticker.ClassifyAllToBinaryClass(bsvc3.PredictAll(X3)))
	// Case 4: fully-separable 3 and 1 points
	X4, Y4 := sticker.FeatureVectors{
		sticker.FeatureVector{sticker.KeyValue32{0, 0.0}, sticker.KeyValue32{1, 0.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, -1.0}, sticker.KeyValue32{1, -1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 1.0}, sticker.KeyValue32{1, 1.0}},
		sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 2.0}},
	}, []bool{false, false, false, true}
	bsvc4 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X4, Y4, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(t, Y4).Equal(sticker.ClassifyAllToBinaryClass(bsvc4.PredictAll(X4)))
	// Expect any debug log.
	var debugBuffer bytes.Buffer
	goassert.New(t).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X1, Y1, C, epsilon, log.New(&debugBuffer, "", 0)))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func createBenchmarkDatasetForBinaryClassifier() (sticker.FeatureVectors, []bool) {
	rng := rand.New(rand.NewSource(0))
	n, d := 1000, 25
	X, Y := make(sticker.FeatureVectors, 2*n), make([]bool, 2*n)
	X[0] = sticker.FeatureVector{sticker.KeyValue32{0, 0.2}, sticker.KeyValue32{1, 0.0}}
	X[n+0] = sticker.FeatureVector{sticker.KeyValue32{0, 0.8}, sticker.KeyValue32{1, 0.0}}
	for i := 1; i < n; i++ {
		X[i] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() - 2.0}, sticker.KeyValue32{1, rng.Float32() - 0.5}}
		X[n+i] = sticker.FeatureVector{sticker.KeyValue32{0, rng.Float32() + 2.0}, sticker.KeyValue32{1, rng.Float32() + 0.5}}
	}
	for i := 0; i < n; i++ {
		for j := 2; j < d; j++ {
			X[i] = append(X[i], sticker.KeyValue32{uint32(j), -0.5})
			X[n+i] = append(X[n+i], sticker.KeyValue32{uint32(j), +0.5})
		}
		Y[n+i] = true
	}
	return X, Y
}

func BenchmarkBinaryClassifierTrainer_L1SVC_DualCD(b *testing.B) {
	C, epsilon, debug := float32(3.0), float32(0.1), (*log.Logger)(nil)
	X, Y := createBenchmarkDatasetForBinaryClassifier()
	// Check the integrity
	bsvc := goassert.New(b).SucceedNew(BinaryClassifierTrainer_L1SVC_DualCD(X, Y, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(b, Y).Equal(sticker.ClassifyAllToBinaryClass(bsvc.PredictAll(X)))
	// Check the number of support vectors.
	nSVs := 0
	for _, v := range bsvc.Beta {
		if sticker.Abs32(v) > 1.0e-05 {
			nSVs++
		}
	}
	goassert.New(b, 2).Equal(nSVs)
	if b.Failed() {
		b.FailNow()
	}
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		BinaryClassifierTrainer_L1SVC_DualCD(X, Y, C, epsilon, nil)
	}
}

func BenchmarkBinaryClassifierTrainer_L2SVC_PrimalCD(b *testing.B) {
	C, epsilon, debug := float32(3.0), float32(0.1), (*log.Logger)(nil)
	X, Y := createBenchmarkDatasetForBinaryClassifier()
	// Check the integrity
	bsvc := goassert.New(b).SucceedNew(BinaryClassifierTrainer_L2SVC_PrimalCD(X, Y, C, epsilon, debug)).(*sticker.BinaryClassifier)
	goassert.New(b, Y).Equal(sticker.ClassifyAllToBinaryClass(bsvc.PredictAll(X)))
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		BinaryClassifierTrainer_L2SVC_PrimalCD(X, Y, C, epsilon, nil)
	}
}
