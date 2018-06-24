package regression

import (
	"math"
	"sync"

	"github.com/bsm/reason/core"
)

// Evaluator is a basic regression evaluator
type Evaluator struct {
	weight float64 // total weight observed
	sum    float64 // sum of all values

	resSum  float64 // residual sum
	resSum2 float64 // residual sum of squares
	totSum2 float64 // total sum of squares

	mu sync.RWMutex
}

// NewEvaluator inits a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Record records a prediction
func (e *Evaluator) Record(predicted, actual, weight float64) {
	if !core.IsNum(predicted) || !core.IsNum(actual) {
		return
	}

	residual := actual - predicted

	e.mu.Lock()
	if e.weight != 0 {
		delta := actual - e.sum/e.weight
		e.totSum2 += delta * delta * weight
	}

	e.resSum += math.Abs(residual) * weight
	e.resSum2 += residual * residual * weight

	e.sum += actual * weight
	e.weight += weight
	e.mu.Unlock()
}

// Total returns the total weight observed
func (e *Evaluator) Total() float64 {
	e.mu.RLock()
	weight := e.weight
	e.mu.RUnlock()
	return weight
}

// Mean returns the mean value observed
func (e *Evaluator) Mean() float64 {
	e.mu.RLock()
	weight := e.weight
	sum := e.sum
	e.mu.RUnlock()

	if weight > 0 {
		return sum / weight
	}
	return 0.0
}

// MAE returns the mean absolute error
func (e *Evaluator) MAE() float64 {
	e.mu.RLock()
	weight := e.weight
	resSum := e.resSum
	e.mu.RUnlock()

	if weight > 0 {
		return resSum / weight
	}
	return 0.0
}

// MSE returns the mean square error
func (e *Evaluator) MSE() float64 {
	e.mu.RLock()
	weight := e.weight
	resSum2 := e.resSum2
	e.mu.RUnlock()

	if weight > 0 {
		return resSum2 / weight
	}
	return 0.0
}

// RMSE returns the root mean square error
func (e *Evaluator) RMSE() float64 {
	return math.Sqrt(e.MSE())
}

// R2 returns the R² coefficient of determination
func (e *Evaluator) R2() float64 {
	e.mu.RLock()
	resSum2 := e.resSum2
	totSum2 := e.totSum2
	e.mu.RUnlock()

	if totSum2 > 0 {
		return 1 - resSum2/totSum2
	}
	return 0.0
}
