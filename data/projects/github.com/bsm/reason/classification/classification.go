package classification

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// Predictions is a slice of predictions.
type Predictions []Prediction

// Best returns the most accurate prediction.
func (pp Predictions) Best() *Prediction {
	if n := len(pp); n != 0 {
		return &pp[n-1]
	}
	return nil
}

// Prediction is a standard prediction of a classification.
type Prediction struct {
	util.Vector
}

// Top returns the most probable category and its probability.
func (p *Prediction) Top() (core.Category, float64) {
	n, w := p.TopW()
	if w > 0 {
		return n, w / p.Weight()
	}
	return n, 0.0
}

// TopW returns the most probable category with its observed weight.
func (p *Prediction) TopW() (core.Category, float64) {
	n, w := p.Max()
	return core.Category(n), w
}

// W returns the weight of the given category.
func (p *Prediction) W(cat core.Category) float64 {
	return p.Get(int(cat))
}

// P returns the probability of the given category.
func (p *Prediction) P(cat core.Category) float64 {
	if w := p.W(cat); w > 0 {
		return w / p.Weight()
	}
	return 0.0
}
