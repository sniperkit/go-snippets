package hoeffding

// Config configures behaviour
type Config struct {
	// The number of training instances a leaf node should observe
	// between split attempts.
	// Default: 200
	GracePeriod int

	// The number of training instances the tree should observe
	// between pruning attempts. To disable, set to <0.
	// Default: 100,000
	PrunePeriod int

	// The maximum number of active/learning leaf nodes. To prevent a
	// tree from growing too large and to avoid unnecessary computation
	// trees can be pruned to only focus on training nodes with the highest
	// merit and the best chance of a split.
	// Default: 1,000,000
	MaxLearningNodes int

	// The allowable error in a split decision - values closer
	// to zero will take longer to decide.
	// Default: 0.0000001
	SplitConfidence float64

	// Threshold below which a split will be forced to break ties
	// Default: 0.05
	TieThreshold float64

	// By enabling this option, tracing notification events will be
	// emitted via the Traces channel after each training cycle. This
	// is for debug purposes only. When enabled, you must consume
	// the Traces channel to avoid locked threads.
	// Default: false
	EnableTracing bool
}

// Norm inits and normalizes the config
func (c *Config) Norm() {
	if c.GracePeriod <= 0 {
		c.GracePeriod = 200
	}
	if c.PrunePeriod == 0 {
		c.PrunePeriod = 100000
	}
	if c.MaxLearningNodes <= 0 {
		c.MaxLearningNodes = 1000000
	}
	if c.SplitConfidence <= 0 {
		c.SplitConfidence = 1e-7
	}
	if c.TieThreshold <= 0 {
		c.TieThreshold = 0.05
	}
}
