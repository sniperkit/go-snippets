package logging

import (
	"github.com/go-zero-boilerplate/loggers"
)

//Logger is a the interface
type Logger interface {
	loggers.LoggerRFC5424Simple
	Trace(s string) LogTracer
	TraceDebug(s string) LogDebugTracer

	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	DeferredRecoverStack(debugMessage string)
}
