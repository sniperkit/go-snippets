package plugin

import (
	"log"
	"math/rand"

	"github.com/hiro4bbh/sticker"
)

// BinaryClassifierTrainer_L1SVC_DualCD trains a L1-Support Vector Classifier with Dual Coordinate Descent.
// This is registered to sticker.BinaryClassifierTrainers.
//
// This function returns no error currently.
//
// Reference: C. Hsieh, K. Chang, C. Lin, S. S. Keerthi, and S. Sundararajan. "A Dual Coordinate Descent Method for Large-Scale Linear SVM." Proceedings of the 25th international conference on Machine learning, ACM, 2008.
func BinaryClassifierTrainer_L1SVC_DualCD(X sticker.FeatureVectors, Y []bool, C, epsilon float32, debug *log.Logger) (*sticker.BinaryClassifier, error) {
	rng := rand.New(rand.NewSource(0))
	n, d := len(X), X.Dim()
	b, w := float32(0.0), make([]float32, d)
	Qdiag, beta := make([]float32, n), make([]float32, n)
	pi := make([]int, n)
	for i, xi := range X {
		q := float32(1.0 * 1.0)
		for _, xipair := range xi {
			q += xipair.Value * xipair.Value
		}
		Qdiag[i] = q
		pi[i] = i
	}
	nactives := n
	maxG, minG := sticker.Inf32(+1), sticker.Inf32(-1)
	for t := 0; t < 1000; t++ {
		maxPG, minPG := sticker.Inf32(-1), sticker.Inf32(+1)
		// Shuffle the active entries.
		for i_ := 0; i_ < nactives-1; i_++ {
			j_ := i_ + rng.Intn(n-i_)
			pi[i_], pi[j_] = pi[j_], pi[i_]
		}
		for i_ := 0; i_ < nactives; i_++ {
			i := pi[i_]
			xi, yi, betai := X[i], Y[i], beta[i]
			// G: the gradient of the unconstrained case.
			//   G = t(e_i) Q \beta - 1 = y_i t(w) x_i - 1
			G := b * 1.0
			for _, xipair := range xi {
				G += w[xipair.Key] * xipair.Value
			}
			if !yi {
				G *= -1.0
			}
			G -= 1.0
			// Shrink the active entries if possible.
			if (betai == 0.0 && maxG < G) || (betai == C && G < minG) {
				nactives--
				pi[i_], pi[nactives] = pi[nactives], pi[i_]
				i_--
				continue
			}
			// PG: the projected gradient: G or 0.0 if the optimal value is out of the contraint region.
			//   PG = min{0, G}  if \beta_i = 0
			//        G          if 0 < \beta_i < C
			//        max{0, G}  if \beta_i = C
			PG := G
			if betai == 0.0 {
				if PG > 0.0 {
					PG = 0.0
				}
			} else if betai == C {
				if PG < 0.0 {
					PG = 0.0
				}
			}
			if maxPG < PG {
				maxPG = PG
			}
			if minPG > PG {
				minPG = PG
			}
			if sticker.Abs32(PG) <= 1.0e-06 {
				continue
			}
			// Calculate the next value in the constraint region.
			// d = G/(t(e_i) Q e_i)
			d := G / Qdiag[i]
			newbetai := betai - d
			if newbetai < 0.0 {
				newbetai = 0.0
			} else if newbetai > C {
				newbetai = C
			}
			beta[i] = newbetai
			// Update w = \sum_{i=1}^n \beta_i y_i x_i.
			delta := beta[i] - betai
			if !yi {
				delta *= -1.0
			}
			b += delta * 1.0
			for _, xipair := range xi {
				w[xipair.Key] += delta * xipair.Value
			}
		}
		status := ""
		// Check if the projected gradient is enough small.
		if maxPG-minPG <= epsilon {
			if nactives == n {
				// Checked the PGs of all entries.
				break
			}
			// Check all entries at the next loop.
			status = ": check all entries at the next loop"
			nactives = n
			maxG, minG = sticker.Inf32(+1), sticker.Inf32(-1)
		} else {
			if maxPG > 0 {
				maxG = maxPG
			} else {
				maxG = sticker.Inf32(+1)
			}
			if minPG < 0 {
				minG = minPG
			} else {
				minG = sticker.Inf32(-1)
			}
		}
		if debug != nil {
			debug.Printf("BinaryClassifierTrainer(L1SVC_DualCD): t=%d: maxPG-minPG=%g, nactives=%d%s", t, maxPG-minPG, nactives, status)
		}
	}
	return &sticker.BinaryClassifier{
		Bias:   b,
		Weight: sticker.SparsifyVector(w),
		Beta:   beta,
	}, nil
}

// BinaryClassifierTrainer_L2SVC_PrimalCD trains a L2-Support Vector Classifier with Primal Coordinate Descent.
// This is registered to sticker.BinaryClassifierTrainers.
//
// This is not recommended, because even if the dataset is normalized as saving its sparsity, it is too slow to converge due to its piece-wise quadratic form.
// It is difficult to control the scaling such that the magnitude of the first derivative equals to its corresponding newton step.
// Otherwise, the optimization would be slow even when the first derivative is not enough small.
// Furthermore, even if the optimization stops early, its performance is much worse than L1SVC_DualCD.
//
// This function returns no error currently.
//
// Reference: K. Chang, C. Hsieh, and C. Lin. "Coordinate Descent Method for Large-Scale L2-loss Linear Support Vector Machines." Journal of Machine Learning Research, vol. 9, pp. 1369-1398, 2008.
func BinaryClassifierTrainer_L2SVC_PrimalCD(X sticker.FeatureVectors, Y []bool, C, epsilon float32, debug *log.Logger) (*sticker.BinaryClassifier, error) {
	// sigma is the hyper-parameter for evaluating the decrease of the objective function value.
	beta, sigma := float32(0.5), float32(0.01)
	// b is the bias parameter, and w is the weight vector parameter.
	b, w := float32(0.0), make(sticker.SparseVector)
	// featureMapSet is the set of the map for each feature to the indices of the dataset.
	featureMapSet, valueMapSet := make(map[uint32][]int), make(map[uint32][]float32)
	for i, xi := range X {
		for _, xipair := range xi {
			featureMapSet[xipair.Key] = append(featureMapSet[xipair.Key], i)
			valueMapSet[xipair.Key] = append(valueMapSet[xipair.Key], xipair.Value)
		}
	}
	// Initialize weights with zero.
	for feature := range featureMapSet {
		w[feature] = 0.0
	}
	// HfeatureSet is the set of the H-value (like hessian) used in calculating the upper bound of the step size for each feature.
	HfeatureSet := make(map[uint32]float32)
	for feature, featureMap := range featureMapSet {
		valueMap := valueMapSet[feature]
		Hfeature := float32(0.0)
		for _, i := range featureMap {
			xifeature := valueMap[i]
			Hfeature += C * xifeature * xifeature
		}
		HfeatureSet[feature] = 1 + 2*Hfeature
	}
	featureMapSet[^uint32(0)] = []int{}
	Hintercept := float32(0.0)
	for i := range X {
		if Y[i] {
			Hintercept += C * 1.0 * 1.0
		} else {
			Hintercept += C * 1.0 * 1.0
		}
	}
	HfeatureSet[^uint32(0)] = 1 + 2*Hintercept
	// Z has the predictors.
	Z := make([]float32, len(X))
	for iter := 0; ; iter++ {
		// Remember the number of entries having the non-zero margin at updating the intercept.
		maxAbsDelta, maxAbsD1, nhasMargin := float32(0.0), float32(0.0), 0
		// The iteration order of golang's map is already random.
		for feature, featureMap := range featureMapSet {
			valueMap := valueMapSet[feature]
			// Penalty Term: 1/2 \|w\|_2^2
			// Loss term: \sum_{i=1}^n C \max\{0, 1 - Y[i] t(w) X[i]\}
			// First-order derivative: w[feature] + \sum_{i=1}^n C * 2 \max\{0, 1 - Y[i] t(w) X[i]\} * (-Y[i] X[i][feature]).
			// (Generalized) Second-order derivative: 1 + 2\sum_{i=1}^n C * 2 \delta[1 - Y[i] t(w) X[i] > 0] * X[i][feature]^2.
			// Here, the generalized derivative of max\{0, x\} is defined as \delta[x > 0].
			loss, d1, d2 := float32(0.0), float32(0.0), float32(0.0)
			if feature == ^uint32(0) {
				for i := range X {
					yi := float32(-1.0)
					if Y[i] {
						yi = +1.0
					}
					if margin := 1 - yi*Z[i]; margin > 0 {
						nhasMargin++
						loss += C * margin * margin
						d1 += C * yi * 1.0 * margin
						d2 += C * 1.0 * 1.0
					}
				}
			} else {
				for _, i := range featureMap {
					xifeature := valueMap[i]
					yi := float32(-1.0)
					if Y[i] {
						yi = +1.0
					}
					if margin := 1 - yi*Z[i]; margin > 0 {
						loss += C * margin * margin
						d1 += C * yi * xifeature * margin
						d2 += C * xifeature * xifeature
					}
				}
			}
			d1, d2 = w[feature]-2*d1, 1+2*d2
			r := d1 / d2
			if sticker.Abs32(r) <= 1.0e-05 {
				continue
			}
			// Search the step size linearly.
			// Theorem 1 of Chang+ (2008) ensures that any move no more than this upper bound always decreases the value of the objective function.
			lambda, lambdaUB := float32(1.0), d2/(HfeatureSet[feature]/2+sigma)
			delta := float32(0.0)
			for {
				delta = -lambda * r
				if lambda <= lambdaUB {
					break
				}
				newloss := float32(0.0)
				if feature == ^uint32(0) {
					for i := range X {
						yi := float32(-1.0)
						if Y[i] {
							yi = +1.0
						}
						if newmargin := 1 - yi*(Z[i]+delta*1.0); newmargin > 0 {
							newloss += C * newmargin * newmargin
						}
					}
				} else {
					for _, i := range featureMap {
						yi := float32(-1.0)
						if Y[i] {
							yi = +1.0
						}
						if newmargin := 1 - yi*(Z[i]+delta*valueMap[i]); newmargin > 0 {
							newloss += C * newmargin * newmargin
						}
					}
				}
				// Check the sufficient condition for the convergence.
				diff := (delta*delta/2 + delta*w[feature] + newloss) - loss
				if diff <= -sigma*delta*delta {
					break
				}
				lambda *= beta
			}
			// Update the weight for the feature.
			w[feature] += delta
			// Update each predictor
			if feature == ^uint32(0) {
				for i := range X {
					Z[i] += delta * 1.0
				}
			} else {
				for _, i := range featureMap {
					Z[i] += delta * valueMap[i]
				}
			}
			absDelta := sticker.Abs32(delta)
			if maxAbsDelta < absDelta {
				maxAbsDelta = absDelta
			}
			absD1 := sticker.Abs32(d1)
			if maxAbsD1 < absD1 {
				maxAbsD1 = absD1
			}
		}
		if debug != nil {
			debug.Printf("TrainerBinaryClassifier(L2SVC_PrimalCD): iter=%d: max|delta|=%.5g, max|d1|=%.5g, nhasMargin=%d, w[^uint64(0)]=%g", iter, maxAbsDelta, maxAbsD1, nhasMargin, w[^uint32(0)])
		}
		// The terminal condition is that any absolute value of first-derivative is no more than epsilon.
		if maxAbsD1 <= epsilon {
			break
		}
	}
	// Prune the intercept and weights enough small in float32.
	b = w[^uint32(0)]
	delete(w, ^uint32(0))
	return &sticker.BinaryClassifier{
		Bias:   b,
		Weight: w,
	}, nil
}

func init() {
	sticker.BinaryClassifierTrainers["L1SVC_DualCD"] = BinaryClassifierTrainer_L1SVC_DualCD
	sticker.BinaryClassifierTrainers["L2SVC_PrimalCD"] = BinaryClassifierTrainer_L2SVC_PrimalCD
}
