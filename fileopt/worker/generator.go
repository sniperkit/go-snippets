package worker

import (
	"errors"
)

var (
	ErrProducteJob error = errors.New("Producer: generator producte job error")
)

type Generator interface {
	Generate() <-chan Jobber
}
