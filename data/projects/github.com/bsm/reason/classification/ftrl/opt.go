package ftrl

import (
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/bsm/reason/classification/ftrl/internal"
	"github.com/bsm/reason/core"
)

// Optimizer represents an FTRL optimiser. Regressions
// can only predict values between 0 and 1. For correct results,
// please ensure that all your target values are within that range.
type Optimizer struct {
	opt        *internal.Optimizer
	target     *core.Feature
	predictors []string
	offsets    []int
	config     Config
	mu         sync.RWMutex
}

// Load loads an Optimizer from a reader.
func Load(r io.Reader, config *Config) (*Optimizer, error) {
	opt := new(internal.Optimizer)
	if _, err := opt.ReadFrom(r); err != nil {
		return nil, err
	}

	predictors, offsets, _ := parseFeatures(opt.Model.Features, opt.Target)
	return newOptimizer(opt, predictors, offsets, config)
}

// New inits a new Optimizer using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Optimizer, error) {
	predictors, offsets, size := parseFeatures(model.Features, target)
	opt := internal.NewOptimizer(model, target, size)
	return newOptimizer(opt, predictors, offsets, config)
}

func newOptimizer(opt *internal.Optimizer, predictors []string, offsets []int, c *Config) (*Optimizer, error) {
	feat := opt.Model.Feature(opt.Target)
	if feat == nil {
		return nil, fmt.Errorf("ftrl: unknown feature %q", opt.Target)
	}
	for _, feat := range opt.Model.Features {
		if feat.Strategy != core.Feature_VOCABULARY {
			return nil, fmt.Errorf("ftrl: feature's %q strategy %q is not supported", feat.Name, feat.Strategy.String())
		}
	}

	var config Config
	if c != nil {
		config = *c
	}
	config.Norm()

	return &Optimizer{
		opt:        opt,
		target:     feat,
		predictors: predictors,
		offsets:    offsets,
		config:     config,
	}, nil
}

// Predict performs prediction
func (o *Optimizer) Predict(x core.Example) float64 {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.predict(x, nil)
}

// Trains trains the optimizer with an example and a weight.
func (o *Optimizer) Train(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	var y float64 // target
	switch o.target.Kind {
	case core.Feature_CATEGORICAL:
		if v := o.target.Category(x); v < 0 {
			return
		} else if v > 0 {
			y = 1.0
		}
	case core.Feature_NUMERICAL:
		if v := o.target.Number(x); math.IsNaN(v) {
			return
		} else {
			y = v
		}
	default:
		return
	}

	t := make(map[int]float64, len(o.predictors))

	o.mu.Lock()
	defer o.mu.Unlock()

	delta := o.predict(x, t) - y
	for i, name := range o.predictors {
		feat := o.opt.Model.Features[name]
		bucket, val := featureBV(feat, x, o.offsets[i])
		if bucket < 0 {
			continue
		}

		// calculate gradient
		g := delta * val
		G := g * g

		// calculate sigma
		s := (math.Sqrt(o.opt.Sums[bucket]+G) - math.Sqrt(o.opt.Sums[bucket])) / o.config.Alpha

		// update
		o.opt.Weights[bucket] += g - s*t[bucket]
		o.opt.Sums[bucket] += G
	}
}

// WriteTo implements io.WriterTo
func (o *Optimizer) WriteTo(w io.Writer) (int64, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.opt.WriteTo(w)
}

func (o *Optimizer) predict(x core.Example, t map[int]float64) float64 {
	var wTx float64

	for i, name := range o.predictors {
		feat := o.opt.Model.Features[name]
		bucket, val := featureBV(feat, x, o.offsets[i])
		if bucket < 0 {
			continue
		}

		sign := 1.0
		if o.opt.Weights[bucket] < 0 {
			sign = -1.0
		}
		fabs := o.opt.Weights[bucket] * sign

		if fabs <= o.config.L1 {
			if t != nil {
				t[bucket] = 0
			}
		} else {
			step := o.config.L2 + (o.config.Beta+math.Sqrt(o.opt.Sums[bucket]))/o.config.Alpha
			factor := sign * (o.config.L1 - fabs) / step
			if t != nil {
				t[bucket] = factor
			}
			wTx += factor * val
		}
	}
	return 1 / (1 + math.Exp(-math.Max(math.Min(wTx, 35), -35)))
}
