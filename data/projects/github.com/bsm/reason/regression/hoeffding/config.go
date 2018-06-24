package hoeffding

import (
	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/regression"
)

// Config configures behaviour
type Config struct {
	common.Config

	// The split criterion to use for evaluating splits
	// Default: classification.DefaultSplitCriterion()
	SplitCriterion regression.SplitCriterion
}

// Norm inits and normalizes the config
func (c *Config) Norm() {
	c.Config.Norm()

	if c.SplitCriterion == nil {
		c.SplitCriterion = regression.DefaultSplitCriterion()
	}
}
