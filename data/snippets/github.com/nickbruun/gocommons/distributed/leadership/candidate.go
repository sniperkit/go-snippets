package leadership

// Leadership handler.
//
// Invoked when a candidate becomes a leader. When the leadership ends, an
// empty struct is sent over end. Returning from the handler function results
// in termination of the leadership.
type LeadershipHandler func(end <-chan struct{})

// Candidate.
//
// Offers and handles leadership.
type Candidate interface {
	// Stop offering candidacy.
	Stop()
}
