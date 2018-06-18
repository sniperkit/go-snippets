package logging

import (
	apex "github.com/francoishill/log"
)

type LogDebugTracer interface {
	StopDebug(errPtr *error)
}

type localLogDebugTracer struct {
	e *apex.Entry
}

func (l *localLogDebugTracer) StopDebug(errPtr *error) { l.e.StopLevel(apex.DebugLevel, errPtr) }
