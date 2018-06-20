// Package graylog adds to github.com/go-mixins/logrus.ContextLogger the
// ability to log entries to Graylog
//
package graylog

import (
	"github.com/go-mixins/log/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// NewAsyncGraylogHook copied here to reduce import overhead
var NewAsyncGraylogHook = graylog.NewAsyncGraylogHook

// Some useful constants from Logrus to reduce import clutter
const (
	DebugLevel = logrus.DebugLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

// New creates new ContextLogger with Logrus side-logging to specified URIs
func New(graylogURIs ...string) *logrus.ContextLogger {
	res := logrus.New()
	for _, uri := range graylogURIs {
		if uri == "" {
			continue
		}
		res.Hooks.Add(graylog.NewAsyncGraylogHook(uri, nil))
	}
	return res
}
