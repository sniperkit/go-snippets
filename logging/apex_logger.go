package logging

import (
	"fmt"
	"os"
	"strings"

	apex "github.com/francoishill/log"

	"github.com/go-zero-boilerplate/facades/debugging"
)

//NewApexLogger creates a new logger
func NewApexLogger(level apex.Level, handler apex.Handler, apexEntry *apex.Entry, exitOnEmergency bool) Logger {
	//TODO: This set level and handler is global, this should not happen inside `NewApexLogger`
	apex.SetLevel(level)
	apex.SetHandler(handler)

	return &apexLogger{
		Entry:           apexEntry,
		level:           level,
		handler:         handler,
		errStackTraces:  true, //TODO: Is this fine by default?
		exitOnEmergency: exitOnEmergency,
	}
}

type apexLogger struct {
	*apex.Entry
	level           apex.Level
	handler         apex.Handler
	errStackTraces  bool
	exitOnEmergency bool
}

func (l *apexLogger) Emergency(s string) {
	l.Entry.Emergency(s)
	if l.exitOnEmergency {
		os.Exit(1)
	}
}
func (l *apexLogger) Trace(s string) LogTracer {
	return l.Entry.Trace(fmt.Sprintf(s))
}
func (l *apexLogger) TraceDebug(s string) LogDebugTracer {
	return &localLogDebugTracer{l.Entry.TraceLevel(apex.DebugLevel, fmt.Sprintf(s))}
}

func (l *apexLogger) WithError(err error) Logger {
	newEntry := l.Entry.WithError(err)
	if l.errStackTraces {
		stack := strings.Replace(debugging.GetNormalStack(false), "\n", "\\n", -1)
		newEntry = newEntry.WithField("stack", stack)
	}
	return NewApexLogger(l.level, l.handler, newEntry, l.exitOnEmergency)
}

func (l *apexLogger) WithField(key string, value interface{}) Logger {
	return NewApexLogger(l.level, l.handler, l.Entry.WithField(key, value), l.exitOnEmergency)
}

func (l *apexLogger) WithFields(fields map[string]interface{}) Logger {
	return NewApexLogger(l.level, l.handler, l.Entry.WithFields(apex.Fields(fields)), l.exitOnEmergency)
}

func (l *apexLogger) DeferredRecoverStack(debugMessage string) {
	if r := recover(); r != nil {
		logger := l.WithField("recovery", r).WithField("debug", debugMessage)
		stack := strings.Replace(debugging.GetNormalStack(false), "\n", "\\n", -1)
		logger = logger.WithField("stack", stack)
		logger.Alert("Unhandled panic recovered")
	}
}
