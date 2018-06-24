package eval

import (
	"math"
	"sync"
)

// LogLoss, or logarithmic loss, can be used when
// the raw output of the classifier is a numeric probability
// instead of a class label.
type LogLoss struct {
	Epsilon float64 // A small increment to add to avoid taking a log of zero. Default: 1e-15.

	sum, weight float64
	mu          sync.RWMutex
}

// NewLogLoss returns a new evaluator.
func NewLogLoss() *LogLoss {
	return &LogLoss{Epsilon: 1e-15}
}

// Record records the prediction (probability of the actually observed value).
// Assuming the predictions were:
//   [dog: 0.2, cat: 0.5, fish: 0.3]
//   [dog: 0.8, cat: 0.1, fish: 0.1]
//   [dog: 0.6, cat: 0.1, fish: 0.4]
// And the actual observations were:
//   * cat
//   * dog
//   * fish
// Then the recorded values should be:
//   e.Record(0.5, 1)
//   e.Record(0.8, 1)
//   e.Record(0.4, 1)
func (e *LogLoss) Record(probability float64, weight float64) {
	e.mu.Lock()
	e.weight += weight
	e.sum += weight * math.Log(probability+e.Epsilon)
	e.mu.Unlock()
}

//  Value calulates the logarithmic loss.
func (e *LogLoss) Value() float64 {
	var v float64
	e.mu.RLock()
	if e.weight > 0 {
		v = -e.sum / e.weight
	}
	e.mu.RUnlock()
	return v
}
