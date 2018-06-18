package logging

import (
	"os"

	apex "github.com/francoishill/log"
	"github.com/francoishill/log/handlers/multi"

	"github.com/go-zero-boilerplate/facades/logging/text_handler"
)

var logger Logger

//Get returns the singleton logger
func Get() Logger {
	return logger
}

//loggerRFC5424Simple is the same as LoggerRFC5424 but without params and still follows the RFC 5424 format: https://tools.ietf.org/html/rfc5424
//copied from github.com/go-zero-boilerplate/loggers
type loggerRFC5424Simple interface {
	Emergency(s string)
	Alert(s string)
	Critical(s string)
	Error(s string)
	Warn(s string)
	Notice(s string)
	Info(s string)
	Debug(s string)
}

//Logger is a the interface
type Logger interface {
	loggerRFC5424Simple

	Trace(s string) LogTracer
	TraceDebug(s string) LogDebugTracer

	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	DeferredRecoverStack(debugMessage string)
}

//InitLogger initializes the singleton logger instance
func InitLogger(level apex.Level, additionalHandlers []apex.Handler) {
	loggerFields := apex.Fields{}
	apexEntry := apex.WithFields(loggerFields)

	logHandlers := []apex.Handler{}

	logHandlers = append(logHandlers, text_handler.New(os.Stdout, os.Stderr, text_handler.DefaultTimeStampFormat, text_handler.DefaultMessageWidth))
	logHandlers = append(logHandlers, additionalHandlers...)

	exitOnEmergency := true
	logger = NewApexLogger(level, multi.New(logHandlers...), apexEntry, exitOnEmergency)
}
