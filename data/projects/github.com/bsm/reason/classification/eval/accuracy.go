package eval

import (
	"sync"

	"github.com/bsm/reason/core"
)

// Accuracy is a basic classification evaluator. It measures how often the
// classifier makes the correct prediction. It is the ratio between the
// weight of correct predictions and the total weight of predictions.
type Accuracy struct {
	weight, correct float64
	mu              sync.RWMutex
}

// NewAccuracy inits a new evaluator.
func NewAccuracy() *Accuracy {
	return &Accuracy{}
}

// Record records a prediction.
func (e *Accuracy) Record(predicted, actual core.Category, weight float64) {
	if !core.IsCat(predicted) || !core.IsCat(actual) {
		return
	}

	e.mu.Lock()
	e.weight += weight
	if predicted == actual {
		e.correct += weight
	}
	e.mu.Unlock()
}

// Total returns the total weight observed.
func (e *Accuracy) Total() float64 {
	e.mu.RLock()
	weight := e.weight
	e.mu.RUnlock()
	return weight
}

// Correct returns the weight of correct observations.
func (e *Accuracy) Correct() float64 {
	e.mu.RLock()
	correct := e.correct
	e.mu.RUnlock()
	return correct
}

// Accuracy returns the rate of correct predictions.
func (e *Accuracy) Accuracy() float64 {
	t := e.Total()
	if t == 0 {
		return 0
	}
	return e.Correct() / t
}
