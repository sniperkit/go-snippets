package logging

import (
	"fmt"
	"os"
	"strings"

	apex "github.com/francoishill/log"

	"github.com/go-zero-boilerplate/extended-apex-logger/utils/debug"
)

//NewApexLogger creates a new logger
func NewApexLogger(level apex.Level, handler apex.Handler, apexEntry *apex.Entry, exitOnEmergency bool) Logger {
	//TODO: This set level and handler is global, this should not happen inside `NewApexLogger`
	apex.SetLevel(level)
	apex.SetHandler(handler)

	return &apexLogger{
		level:           level,
		handler:         handler,
		apexEntry:       apexEntry,
		errStackTraces:  true, //TODO: Is this fine by default?
		exitOnEmergency: exitOnEmergency,
	}
}

type apexLogger struct {
	level           apex.Level
	handler         apex.Handler
	apexEntry       *apex.Entry
	errStackTraces  bool
	exitOnEmergency bool
}

func (l *apexLogger) Emergency(s string) {
	if l.exitOnEmergency {
		l.apexEntry.Fatalf(s)
		//Sure this is not required since apex Entry.Fatal already calls `os.Exit(1)`. But this is for "safety" if that ever changes in apex
		os.Exit(1)
	} else {
		l.apexEntry.Errorf(s)
	}
}
func (l *apexLogger) Alert(s string) {
	l.apexEntry.Errorf(s)
}
func (l *apexLogger) Critical(s string) {
	l.apexEntry.Errorf(s)
}
func (l *apexLogger) Error(s string) {
	l.apexEntry.Errorf(s)
}
func (l *apexLogger) Warn(s string) {
	l.apexEntry.Warnf(s)
}
func (l *apexLogger) Notice(s string) {
	l.apexEntry.Warnf(s)
}
func (l *apexLogger) Info(s string) {
	l.apexEntry.Infof(s)
}
func (l *apexLogger) Debug(s string) {
	l.apexEntry.Debugf(s)
}
func (l *apexLogger) Trace(s string) LogTracer {
	return l.apexEntry.Trace(fmt.Sprintf(s))
}
func (l *apexLogger) TraceDebug(s string) LogDebugTracer {
	return &localLogDebugTracer{l.apexEntry.TraceLevel(apex.DebugLevel, fmt.Sprintf(s))}
}

func (l *apexLogger) WithError(err error) Logger {
	newEntry := l.apexEntry.WithError(err)
	if l.errStackTraces {
		stack := strings.Replace(debug.GetFullStackTrace_Normal(false), "\n", "\\n", -1)
		newEntry = newEntry.WithField("stack", stack)
	}
	return NewApexLogger(l.level, l.handler, newEntry, l.exitOnEmergency)
}

func (l *apexLogger) WithField(key string, value interface{}) Logger {
	return NewApexLogger(l.level, l.handler, l.apexEntry.WithField(key, value), l.exitOnEmergency)
}

func (l *apexLogger) WithFields(fields map[string]interface{}) Logger {
	return NewApexLogger(l.level, l.handler, l.apexEntry.WithFields(apex.Fields(fields)), l.exitOnEmergency)
}

func (l *apexLogger) DeferredRecoverStack(debugMessage string) {
	if r := recover(); r != nil {
		logger := l.WithField("recovery", r).WithField("debug", debugMessage)
		stack := strings.Replace(debug.GetFullStackTrace_Normal(false), "\n", "\\n", -1)
		logger = logger.WithField("stack", stack)
		logger.Alert("Unhandled panic recovered")
	}
}
