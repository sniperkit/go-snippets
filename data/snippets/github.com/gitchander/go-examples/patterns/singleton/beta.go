package main

import "sync"

type beta struct{}

func (beta) String() string {
	return "beta"
}

var (
	betaPtr  *beta
	betaOnce sync.Once
)

func Beta() *beta {
	betaOnce.Do(func() {
		betaPtr = &beta{}
	})
	return betaPtr
}
