package sticker

import (
	"log"
	"math/rand"
)

// ClassifyToBinaryClass returns true indicating positive label if the z value is positive, otherwise false indicating negative label.
func ClassifyToBinaryClass(z float32) bool {
	return z > float32(0.0)
}

// ClassifyAllToBinaryClass returns the bool slice whose entry indicates positive label if the corresponding z value is positive, otherwise negative label.
func ClassifyAllToBinaryClass(Z []float32) []bool {
	y := make([]bool, len(Z))
	for i, zi := range Z {
		y[i] = ClassifyToBinaryClass(zi)
	}
	return y
}

// BinaryClassifier is the data structure having information about binary classifiers.
//
// BinaryClassifier classifies the given entry x as positive if dot(Weight, x) + Bias > 0, otherwise as negative.
type BinaryClassifier struct {
	// Bias is the bias parameter.
	Bias float32
	// Weight is the weight parameter.
	Weight SparseVector
	// The following members are not required.
	//
	// Beta is used by some solvers (using dual problems) as the optimization target.
	// Weight can be expressed as the sum of y_ix_i weighted with the corresponding elements of Beta.
	Beta []float32
}

// BinaryClassifierTrainer_L1Logistic_PrimalSGD returns an trained BinaryClassifier with FTRL-Proximal (McMahan+ 2013) method for L1-penalized logistic regression.
// This can be used for estimating the probability which the given data point belongs to the positive class, and this algorithm would produce the smaller model.
//
// This function returns no error currently.
//
// References:
//
// (McMahan+ 2013) H. B. McMahan, et al. "Ad Click Prediction: a View from the Trenches." Proceedings of the 19th ACM SIGKDD International Conference on Knowledge Discovery and Data Mining, 2013.
func BinaryClassifierTrainer_L1Logistic_PrimalSGD(X FeatureVectors, Y []bool, C, epsilon float32, debug *log.Logger) (*BinaryClassifier, error) {
	rng := rand.New(rand.NewSource(0))
	// lambda is the penalty parameter.
	lambda := 1.0 / C
	// n is the number of data points, and d is the dimension of a feature vector.
	n, d := len(X), X.Dim()
	// $\alpha$ and $\beta$ is the hyper parameter for the learning rate.
	alpha, beta := float32(1.0), float32(1.0)
	// bias and weight is the classifier parameters.
	bias, weight := float32(0.0), make([]float32, d)
	// m, gSqSum is the auxiliary vectors and the sum of the squared sum of the gradients for the classifier parameters.
	// The first d elements are for weight, and the (d+1)-th element is for bias.
	m, gSqSum := make([]float32, d+1), make([]float32, d+1)
	// perm is the data point index slice for providing the random order at each round.
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	// t is the iteration number.
	t := 1
	// lossPenalty0 is the previous (loss+penalty).
	lossPenalty0 := Inf32(+1.0)
	// Repeat at most 100 epochs, because they are enough epochs for (loss+penalty) to converge.
	for epoch := 0; epoch < 100; epoch++ {
		// Permutate the data points.
		for i := 0; i < n-1; i++ {
			j := i + rng.Intn(n-i)
			perm[i], perm[j] = perm[j], perm[i]
		}
		// loss is the current loss.
		loss := float32(0.0)
		for _, i := range perm {
			// $\bm{x}_i$ is the current data point.
			xi := X[i]
			// $z_i = \bm{w}^\top\bm{x}_i$ is the linear predictor for $\bm{x}_i$.
			zi := bias * 1.0
			for _, xipair := range xi {
				if mj := m[xipair.Key]; Abs32(mj) > lambda {
					signMj := float32(-1.0)
					if mj > 0.0 {
						signMj = +1.0
					}
					wj := (alpha / (beta + Sqrt32(gSqSum[xipair.Key]))) * (mj - signMj*lambda)
					weight[xipair.Key] = wj
					zi += wj * xipair.Value
				} else {
					weight[xipair.Key] = 0.0
				}
			}
			// $p_i = 1/(1 + exp(-z_i))$ is the predicted probability for $\bm{x}_i$.
			pi := 1.0 / (1.0 + Exp32(-zi))
			// $y_i$ is the received correct class.
			yi := float32(0.0)
			if Y[i] {
				yi = +1.0
			}
			// Calculate the $l_i = -y_i\log(p_i) - (1 - y_i)\log(1 - p_i)$.
			if Y[i] {
				x := float32(0.0)
				if x < -zi {
					x = -zi
				}
				loss += x + Log32(Exp32(0-x)+Exp32(-zi-x))
			} else {
				x := float32(0.0)
				if x < zi {
					x = zi
				}
				loss += x + Log32(Exp32(0-x)+Exp32(zi-x))
			}
			// gBias is the gradient for the bias.
			gBias := -(yi - pi) * 1.0
			// Update the bias and the squared sum of the gradients for the bias.
			bias -= (alpha / (beta + Sqrt32(gSqSum[d]))) * gBias * 1.0
			gSqSum[d] += gBias * gBias
			for _, xipair := range xi {
				// gj is the gradient for the weight parameter.
				gj := -(yi - pi) * xipair.Value
				// Update the auxiliary vector and the squared sum of the gradient for the weight vector.
				gSqSumj := gSqSum[xipair.Key] + gj*gj
				sigmaj := (Sqrt32(gSqSumj) - Sqrt32(gSqSum[xipair.Key])) / alpha
				m[xipair.Key] += sigmaj*weight[xipair.Key] - gj
				gSqSum[xipair.Key] = gSqSumj
			}
			t++
		}
		// Calculate the penalty term.
		penalty := float32(0.0)
		for _, wfeature := range weight {
			penalty += Abs32(wfeature)
		}
		penalty *= lambda
		if debug != nil {
			debug.Printf("BinaryClassifierTrainer(L1Logistic_PrimalSGD): epoch=%d: lambda=%g, penalty=%g, loss=%g, penalty+loss=%g", epoch, lambda, penalty, loss, penalty+loss)
		}
		// Terminate if the relative difference between the previous (loss+penalty) and (loss+penalty) is below epsilon.
		lossPenalty := loss + penalty
		if !IsInf32(lossPenalty0, +1.0) && (lossPenalty0-lossPenalty)/lossPenalty0 < epsilon {
			break
		}
		lossPenalty0 = lossPenalty
	}
	return &BinaryClassifier{
		Bias:   bias,
		Weight: SparsifyVector(weight),
	}, nil
}

// BinaryClassifierTrainer_L1SVC_PrimalSGD trains a L1-Support Vector Classifier with primal stochastic gradient descent.
// This is registered to BinaryClassifierTrainers.
//
// The used update procedure is the one used by Online Passive-Aggressive Algorithm (Crammer+ 2006) with the dynamic penalty parameter depending on the round number t.
// This update is proven to be safe, that is, this leads to sane results even when the learning rate is large (Karampatziakis+ 2011, SubSection 4.2).
// Thus, although we fix the eta0 as 1.0 and the learning rate as eta0 / t, this algorithm is enough fast and accurate.
//
// This function returns no error currently.
//
// Reference:
//
// (Crammer+ 2006) K.Crammer, O. Dekel, J. Keshet, S. Shalev-Shwarts, and Y. Singer. "Online Passive-Aggressive Algorithms." Journal of Machine Learning Research, vol. 7, pp. 551-585, 2006.
//
// (Karampatziakis+ 2011) N. Karampatziakis, and J. Langford, "Online Importance Weight Aware Updates." Association for Uncertainty in Artificial Intelligence, 2011.
func BinaryClassifierTrainer_L1SVC_PrimalSGD(X FeatureVectors, Y []bool, C, epsilon float32, debug *log.Logger) (*BinaryClassifier, error) {
	rng := rand.New(rand.NewSource(0))
	n, d := len(X), X.Dim()
	b, w := float32(0.0), make([]float32, d)
	// Qdiag holds the squared L2-norm of each entry.
	Qdiag := make([]float32, n)
	// pi holds the permutation indices on all entries.
	pi := make([]int, n)
	for i, xi := range X {
		q := float32(1.0 * 1.0)
		for _, xipair := range xi {
			q += xipair.Value * xipair.Value
		}
		Qdiag[i] = q
		pi[i] = i
	}
	// eta0 is the ratio of the learning rate.
	eta0 := float32(1.0)
	// t is the number of update iterations.
	t := 1
	for epoch := 0; epoch < 1000; epoch++ {
		// Shuffle all entries.
		for i_ := 0; i_ < n-1; i_++ {
			j_ := i_ + rng.Intn(n-i_)
			pi[i_], pi[j_] = pi[j_], pi[i_]
		}
		maxGL1 := float32(0.0)
		for _, i := range pi {
			xi, yi := X[i], float32(-1.0)
			if Y[i] {
				yi = float32(+1.0)
			}
			// linear predictor: z_i = b + t(w)x_i
			zi := b * 1.0
			for _, xipair := range xi {
				zi += w[xipair.Key] * xipair.Value
			}
			// loss: l_i = C\max\{0, 1 - y_iz_i\}
			lossi := C * (1.0 - yi*zi)
			if lossi > 0.0 {
				// Step size: s_i = y_i \min\{eta0/t, l_i/(t(x_i)x_i)\}
				// Here, the weights are normalized for the case of size 1 sample.
				si := lossi / Qdiag[i]
				lambdai := eta0 / float32(t)
				if si > lambdai {
					si = lambdai
				}
				si *= yi
				// Update the bias and weights: w_{t+1} = w_{t} + s_ix_i
				gib := si * 1.0
				b += gib
				gL1 := Abs32(gib)
				for _, xipair := range xi {
					gij := si * xipair.Value
					w[xipair.Key] += gij
					gL1 += Abs32(gij)
				}
				if maxGL1 < gL1 {
					maxGL1 = gL1
				}
			}
			t++
		}
		if debug != nil {
			debug.Printf("BinaryClassifierTrainer(L1SVC_PrimalSGD): epoch=%d: max||g||_1=%g", epoch, maxGL1)
		}
		// The termination condition is that the maximum of L1-norm is less than or equal to epsilon.
		if maxGL1 <= epsilon {
			break
		}
	}
	return &BinaryClassifier{
		Bias:   b,
		Weight: SparsifyVector(w),
	}, nil
}

// BinaryClassifierTrainer is the type of binary classifier trainers.
// A trainer returns a new BinaryClassifier on X and Y.
// C is the inverse of the penalty parameter.
// epsilon is the tolerance parameter for checking the convergence.
// debug is used for debug logs.
type BinaryClassifierTrainer func(X FeatureVectors, Y []bool, C, epsilon float32, debug *log.Logger) (*BinaryClassifier, error)

// BinaryClassifierTrainers is the map from the binary classifier trainer name to the corresponding binary classifier trainer.
var BinaryClassifierTrainers = map[string]BinaryClassifierTrainer{
	"L1Logistic_PrimalSGD": BinaryClassifierTrainer_L1Logistic_PrimalSGD,
	"L1SVC_PrimalSGD":      BinaryClassifierTrainer_L1SVC_PrimalSGD,
}

// Predict returns the predicted value dot(Weight, x) + Bias.
func (bc *BinaryClassifier) Predict(x FeatureVector) float32 {
	z := bc.Bias
	for _, xipair := range x {
		z += bc.Weight[xipair.Key] * xipair.Value
	}
	return z
}

// PredictAndCount returns the predicted value dot(Weight, x) + Bias and the splitter count (the number of times the splitter hits).
func (bc *BinaryClassifier) PredictAndCount(x FeatureVector) (float32, uint32) {
	z, c := bc.Bias, uint32(0)
	for _, xipair := range x {
		if w, ok := bc.Weight[xipair.Key]; ok {
			z += w * xipair.Value
			c++
		}
	}
	return z, c
}

// PredictAll returns the predicted values dot(Weight, x) + Bias for each feature vector.
func (bc *BinaryClassifier) PredictAll(X FeatureVectors) []float32 {
	Z := make([]float32, len(X))
	for i, x := range X {
		Z[i] = bc.Predict(x)
	}
	return Z
}

// PredictAndCountAll returns the predicted values dot(Weight, x) + Bias and the splitter count for each feature vector.
func (bc *BinaryClassifier) PredictAndCountAll(X FeatureVectors) ([]float32, []uint32) {
	Z, counts := make([]float32, len(X)), make([]uint32, len(X))
	for i, x := range X {
		Z[i], counts[i] = bc.PredictAndCount(x)
	}
	return Z, counts
}

// ReportPerformance returns the true-negative/false-negative/false-positive/true-positive and the predicted values on X.
func (bc *BinaryClassifier) ReportPerformance(X FeatureVectors, Y []bool) (tn, fn, fp, tp uint, predVals []float32, Yhat []bool) {
	predVals = bc.PredictAll(X)
	Yhat = ClassifyAllToBinaryClass(predVals)
	for i, yihat := range Yhat {
		if yihat {
			if Y[i] {
				tp++
			} else {
				fp++
			}
		} else {
			if Y[i] {
				fn++
			} else {
				tn++
			}
		}
	}
	return
}
