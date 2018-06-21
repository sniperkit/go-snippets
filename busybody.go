package busybody

import (
	"os"

	"github.com/zerklabs/auburn/log"
)

var (
	hostname string
)

const (
	HealthyState int = iota
	SuspiciousState
	FaultyState
)

type Introduction struct {
	Key       string
	Id        string
	Uri       string
	connected bool
	state     int
}

func init() {
	h, err := os.Hostname()
	if err != nil {
		log.Errorf("error fetching hostname: %s", err)
		os.Exit(1)
	}

	hostname = h
}
