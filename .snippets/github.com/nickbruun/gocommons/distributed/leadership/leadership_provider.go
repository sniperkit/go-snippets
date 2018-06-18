package leadership

// Leadership provider.
type LeadershipProvider interface {
	// Get a candidate.
	GetCandidate(data []byte, leadershipHandler LeadershipHandler) Candidate
}
