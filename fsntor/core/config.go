package core

import "errors"

var ErrMalformedConfiguration = errors.New("malformated fsntor configuration")

type patterns []string
type executor struct {
	Command string `json:"command"`
	WaitFor string `json:"waitFor"`
	Timeout string `json:"timeout"`
}
type action []executor
type actions map[string]action
type Task struct {
	Patterns patterns
	Actions  actions
}

type Configurer interface {
	Parse([]byte) ([]Task, error)
}
