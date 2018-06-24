package hoeffding

import (
	"github.com/bsm/reason/classification"
	common "github.com/bsm/reason/common/hoeffding"
)

// Config configures behaviour
type Config struct {
	common.Config

	// The split criterion to use for evaluating splits
	// Default: classification.DefaultSplitCriterion()
	SplitCriterion classification.SplitCriterion
}

// Norm inits and normalizes the config
func (c *Config) Norm() {
	c.Config.Norm()

	if c.SplitCriterion == nil {
		c.SplitCriterion = classification.DefaultSplitCriterion()
	}
}
