// Package regression contains tools for solving regression problems
package regression

import "github.com/bsm/reason/util"

// Predictions is a slice of predictions.
type Predictions []Prediction

// Best returns the most accurate prediction.
func (pp Predictions) Best() *Prediction {
	if n := len(pp); n != 0 {
		return &pp[n-1]
	}
	return nil
}

// Prediction is a standard prediction of a regression.
type Prediction struct{ util.StreamStats }
