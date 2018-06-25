package crawler

type Options struct {
	SameHostOnly  bool
	MaxVisits     int
	FollowingProb float64
	Tolerance     float64
	Parallel      bool
	FileType     string
}

func NewOptions() *Options {
	// Use defaults except for Extender
	return &Options{
		true,
		0,
		0.85,
		0.0001,
		false,
		"bin",
	}
}
