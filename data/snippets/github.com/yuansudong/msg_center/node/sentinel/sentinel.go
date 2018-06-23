package sentinel

type (
	// Sentinel is used as description
	Sentinel struct {
		l *leader
	}
)

// CreateSentinel is used as
func CreateSentinel(url string) *Sentinel {
	pSentinel := &Sentinel{}
	pSentinel.l = createLeader(url)
	return pSentinel
}

// Run is used as run leader mode and listen system signal
func (s *Sentinel) Run() {
	s.l.Run()
}
