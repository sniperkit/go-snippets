package ftrl

// Config configures behaviour
type Config struct {
	// Learn rate alpha parameter.
	// Default: 0.1
	Alpha float64
	// Learn rate beta parameter.
	// Default: 1.0
	Beta float64
	// Regularization strength #1.
	// Default: 1.0
	L1 float64
	// Regularization strength #2.
	// Default: 0.1
	L2 float64
}

// Norm inits and normalizes the config
func (c *Config) Norm() {
	if c.Alpha <= 0 {
		c.Alpha = 0.1
	}
	if c.Beta <= 0 {
		c.Beta = 1.0
	}
	if c.L1 <= 0 {
		c.L1 = 1.0
	}
	if c.L2 <= 0 {
		c.L2 = 0.1
	}
}
