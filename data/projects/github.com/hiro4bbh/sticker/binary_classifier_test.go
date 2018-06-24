package sticker

import (
	"bytes"
	"log"
	"math/rand"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestClassifAllyToBinaryClass(t *testing.T) {
	goassert.New(t, []bool{false, false, true}).Equal(ClassifyAllToBinaryClass([]float32{-1.0, 0.0, 2.0}))
}

func TestBinaryClassifierPredict(t *testing.T) {
	X := FeatureVectors{
		FeatureVector{KeyValue32{0, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
		FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{8, 6.0}},
		FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{5, 7.0}, KeyValue32{8, 6.0}},
	}
	bc := &BinaryClassifier{
		Bias:   -1.0,
		Weight: SparsifyVector([]float32{-1.0, 1.0, 2.0, 0.0, 0.0, 0.0, 0.0, 0.0, 3.0, 0.0}),
	}
	goassert.New(t, []float32{-2.0, 0.0, 23.0, 23.0}).Equal(bc.PredictAll(X))
	goassert.New(t, []float32{-2.0, 0.0, 23.0, 23.0}, []uint32{1, 2, 3, 3}).Equal(bc.PredictAndCountAll(X))
}

func TestBinaryClassifierReportPerformance(t *testing.T) {
	X, Y1, Y2 := FeatureVectors{
		FeatureVector{KeyValue32{0, -2.0}}, FeatureVector{KeyValue32{0, -1.0}}, FeatureVector{KeyValue32{0, 0.0}},
		FeatureVector{KeyValue32{0, 1.0}}, FeatureVector{KeyValue32{0, 2.0}},
	}, []bool{false, false, false, true, true}, []bool{false, false, false, false, false}
	bc := &BinaryClassifier{
		Bias:   -1.0,
		Weight: SparsifyVector([]float32{1.0}),
	}
	goassert.New(t, uint(3), uint(1), uint(0), uint(1), []float32{-3.0, -2.0, -1.0, 0.0, 1.0}, []bool{false, false, false, false, true}).Equal(bc.ReportPerformance(X, Y1))
	goassert.New(t, uint(4), uint(0), uint(1), uint(0), []float32{-3.0, -2.0, -1.0, 0.0, 1.0}, []bool{false, false, false, false, true}).Equal(bc.ReportPerformance(X, Y2))
}

func TestBinaryClassifierTrainer_L1Logistic_PrimalSGD(t *testing.T) {
	C, epsilon, debug := float32(10.0), float32(0.01), (*log.Logger)(nil)
	// Case 1: fully-separable 1x2 points
	X1, Y1 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
	}, []bool{false, true}
	bsvc1 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X1, Y1, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y1).Equal(ClassifyAllToBinaryClass(bsvc1.PredictAll(X1)))
	// Case 2: fully-separable 2x2 points
	X2, Y2 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{false, false, true, true}
	bsvc2 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X2, Y2, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y2).Equal(ClassifyAllToBinaryClass(bsvc2.PredictAll(X2)))
	// Case 3: fully-separable 1 and 3 points
	X3, Y3 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{true, false, true, true}
	bsvc3 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X3, Y3, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y3).Equal(ClassifyAllToBinaryClass(bsvc3.PredictAll(X3)))
	// Case 4: fully-separable 3 and 1 points
	X4, Y4 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{false, false, false, true}
	bsvc4 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X4, Y4, C, epsilon, nil)).(*BinaryClassifier)
	goassert.New(t, Y4).Equal(ClassifyAllToBinaryClass(bsvc4.PredictAll(X4)))
	// Expect any debug log.
	var debugBuffer bytes.Buffer
	goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X1, Y1, C, epsilon, log.New(&debugBuffer, "", 0)))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func TestBinaryClassifierTrainer_L1SVC_PrimalSGD(t *testing.T) {
	C, epsilon, debug := float32(1.0), float32(0.01), (*log.Logger)(nil)
	// Case 1: fully-separable 1x2 points
	X1, Y1 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
	}, []bool{false, true}
	bsvc1 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X1, Y1, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y1).Equal(ClassifyAllToBinaryClass(bsvc1.PredictAll(X1)))
	// Case 2: fully-separable 2x2 points
	X2, Y2 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{false, false, true, true}
	bsvc2 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X2, Y2, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y2).Equal(ClassifyAllToBinaryClass(bsvc2.PredictAll(X2)))
	// Case 3: fully-separable 1 and 3 points
	X3, Y3 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{true, false, true, true}
	bsvc3 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X3, Y3, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(t, Y3).Equal(ClassifyAllToBinaryClass(bsvc3.PredictAll(X3)))
	// Case 4: fully-separable 3 and 1 points
	X4, Y4 := FeatureVectors{
		FeatureVector{KeyValue32{0, 0.0}, KeyValue32{1, 0.0}},
		FeatureVector{KeyValue32{0, -1.0}, KeyValue32{1, -1.0}},
		FeatureVector{KeyValue32{0, 1.0}, KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 2.0}},
	}, []bool{false, false, false, true}
	bsvc4 := goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X4, Y4, C, epsilon, nil)).(*BinaryClassifier)
	goassert.New(t, Y4).Equal(ClassifyAllToBinaryClass(bsvc4.PredictAll(X4)))
	// Expect any debug log.
	var debugBuffer bytes.Buffer
	goassert.New(t).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X1, Y1, C, epsilon, log.New(&debugBuffer, "", 0)))
	goassert.New(t, true).Equal(debugBuffer.String() != "")
}

func createBenchmarkDatasetForBinaryClassifier() (FeatureVectors, []bool) {
	rng := rand.New(rand.NewSource(0))
	n, d := 1000, 25
	X, Y := make(FeatureVectors, 2*n), make([]bool, 2*n)
	X[0] = FeatureVector{KeyValue32{0, 0.2}, KeyValue32{1, 0.0}}
	X[n+0] = FeatureVector{KeyValue32{0, 0.8}, KeyValue32{1, 0.0}}
	for i := 1; i < n; i++ {
		X[i] = FeatureVector{KeyValue32{0, rng.Float32() - 2.0}, KeyValue32{1, rng.Float32() - 0.5}}
		X[n+i] = FeatureVector{KeyValue32{0, rng.Float32() + 2.0}, KeyValue32{1, rng.Float32() + 0.5}}
	}
	for i := 0; i < n; i++ {
		for j := 2; j < d; j++ {
			X[i] = append(X[i], KeyValue32{uint32(j), -0.5})
			X[n+i] = append(X[n+i], KeyValue32{uint32(j), +0.5})
		}
		Y[n+i] = true
	}
	return X, Y
}

func BenchmarkBinaryClassifierTrainer_L1Logistic_PrimalSGD(b *testing.B) {
	C, epsilon, debug := float32(1.0), float32(0.01), (*log.Logger)(nil)
	X, Y := createBenchmarkDatasetForBinaryClassifier()
	// Check the integrity
	bsvc := goassert.New(b).SucceedNew(BinaryClassifierTrainer_L1Logistic_PrimalSGD(X, Y, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(b, Y).Equal(ClassifyAllToBinaryClass(bsvc.PredictAll(X)))
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		BinaryClassifierTrainer_L1SVC_PrimalSGD(X, Y, C, epsilon, nil)
	}
}

func BenchmarkBinaryClassifierTrainer_L1SVC_PrimalSGD(b *testing.B) {
	C, epsilon, debug := float32(1.0), float32(0.01), (*log.Logger)(nil)
	X, Y := createBenchmarkDatasetForBinaryClassifier()
	// Check the integrity
	bsvc := goassert.New(b).SucceedNew(BinaryClassifierTrainer_L1SVC_PrimalSGD(X, Y, C, epsilon, debug)).(*BinaryClassifier)
	goassert.New(b, Y).Equal(ClassifyAllToBinaryClass(bsvc.PredictAll(X)))
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		BinaryClassifierTrainer_L1SVC_PrimalSGD(X, Y, C, epsilon, nil)
	}
}
