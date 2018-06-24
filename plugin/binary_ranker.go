package plugin

import (
	"log"
	"math/rand"

	"github.com/hiro4bbh/sticker"
)

// BinaryRankerTrainer_L1SVC_PrimalSGD trains a L1-Support Vector Ranker with primal stochastic gradient descent.
// This is registered to BinaryRankerTrainers.
//
// This is the optimized implementation based on sticker.BinaryClassifierTrainer_L1SVC_PrimalSGD.
//
// This function returns no error currently.
func BinaryRankerTrainer_L1SVC_PrimalSGD(X sticker.FeatureVectors, pairIndices [][2]int, pairMargins []float32, pairCs []float32, epsilon float32, debug *log.Logger) (*sticker.BinaryClassifier, error) {
	rng := rand.New(rand.NewSource(0))
	n, d := len(pairIndices), X.Dim()
	w := make([]float32, d)
	Q, A := make([]float32, n), make([]float32, n)
	// pi holds the permutation indices on all entries.
	pi := make([]int, n)
	for i := range pairIndices {
		ip, in := pairIndices[i][0], pairIndices[i][1]
		var xpi, xni sticker.FeatureVector
		if ip >= 0 {
			xpi = X[ip]
		}
		if in >= 0 {
			xni = X[in]
		}
		qi, sumAbsXi := float32(0.0), float32(0.0)
		jp, jn := 0, 0
		for jp < len(xpi) && jn < len(xni) {
			delta := float32(0.0)
			if xpi[jp].Key < xni[jn].Key {
				delta = xpi[jp].Value - 0
				jp++
			} else if xpi[jp].Key == xni[jn].Key {
				delta = xpi[jp].Value - xni[jn].Value
				jp++
				jn++
			} else {
				delta = 0 - xni[jn].Value
				jn++
			}
			qi += delta * delta
			sumAbsXi += sticker.Abs32(delta)
		}
		for ; jp < len(xpi); jp++ {
			delta := xpi[jp].Value - 0
			qi += delta * delta
			sumAbsXi += sticker.Abs32(delta)
		}
		for ; jn < len(xni); jn++ {
			delta := 0 - xni[jn].Value
			qi += delta * delta
			sumAbsXi += sticker.Abs32(delta)
		}
		Q[i], A[i] = qi, sumAbsXi
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
		for i_ := 0; i_ < n; i_++ {
			i := pi[i_]
			ip, in := pairIndices[i][0], pairIndices[i][1]
			var xpi, xni sticker.FeatureVector
			if ip >= 0 {
				xpi = X[ip]
			}
			if in >= 0 {
				xni = X[in]
			}
			qi, sumAbsXi := Q[i], A[i]
			// linear predictor: z_i = b + t(w)(x_i^{(+)} - x_i^{(-)})
			zi := float32(0.0)
			for _, xpipair := range xpi {
				zi += w[xpipair.Key] * xpipair.Value
			}
			for _, xnipair := range xni {
				zi -= w[xnipair.Key] * xnipair.Value
			}
			// loss: l_i = \max\{0, (1 - \rho_i) - (+1)z_i\}
			lossi := (1.0 - pairMargins[i]) - zi
			if lossi > 0.0 {
				// Step size: s_i = (+1) \min\{(eta0/t) C_i, l_i/(t(x_i^{(+)} - x_i^{(-)})(x_i^{(+)} - x_i^{(-)}))\}
				si := lossi / qi
				lambdai := eta0 / float32(t) * pairCs[i]
				if si > lambdai {
					si = lambdai
				}
				// Update the bias and weights: w_{t+1} = w_{t} + s_ix_i
				for _, xpipair := range xpi {
					w[xpipair.Key] += si * xpipair.Value
				}
				for _, xnipair := range xni {
					w[xnipair.Key] -= si * xnipair.Value
				}
				gL1 := sticker.Abs32(si * sumAbsXi)
				if maxGL1 < gL1 {
					maxGL1 = gL1
				}
			}
			t++
		}
		if debug != nil {
			debug.Printf("BinaryClassifierTrainer(L1SVC_PrimalSGD): epoch=%d: max||g||_1=%g", epoch, maxGL1)
		}
		// The termination condition is that the maximum of Linf-norm is less than or equal to epsilon.
		if maxGL1 <= epsilon {
			break
		}
	}
	return &sticker.BinaryClassifier{
		Bias:   float32(0.0),
		Weight: sticker.SparsifyVector(w),
	}, nil
}

// BinaryRankerTrainer is the type of binary ranker trainers.
// A trainer returns a new BinaryClassifier on positive/negative pair indices pairIndices on X with the specified pair margins.
// A negative values of pairIndices means the zero-vector.
// C is the penalty parameter slice for reweighting each entry.
// epsilon is the tolerance parameter for checking the convergence.
// debug is used for debug logs.
type BinaryRankerTrainer func(X sticker.FeatureVectors, pairIndices [][2]int, pairMargins []float32, C []float32, epsilon float32, debug *log.Logger) (*sticker.BinaryClassifier, error)

// BinaryRankerTrainers is the map from the binary classifier trainer name to the corresponding binary classifier trainer.
var BinaryRankerTrainers = map[string]BinaryRankerTrainer{
	"L1SVC_PrimalSGD": BinaryRankerTrainer_L1SVC_PrimalSGD,
}
