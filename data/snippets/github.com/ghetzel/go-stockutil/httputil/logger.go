// Utilities for extracting and formatting data encountered in HTTP requests
package httputil

import (
	"strings"

	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger(`httputil`)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Notice
	Warning
	Error
	Fatal
)

type WritableLogger struct {
	level  LogLevel
	prefix string
}

func NewWritableLogger(level LogLevel, prefix ...string) *WritableLogger {
	return &WritableLogger{
		prefix: strings.Join(prefix, ``),
	}
}

func (self *WritableLogger) Write(p []byte) (int, error) {
	switch self.level {
	case Debug:
		Logger.Debugf("%v%v", self.prefix, string(p))
	case Notice:
		Logger.Noticef("%v%v", self.prefix, string(p))
	case Warning:
		Logger.Warningf("%v%v", self.prefix, string(p))
	case Error:
		Logger.Errorf("%v%v", self.prefix, string(p))
	case Fatal:
		Logger.Fatalf("%v%v", self.prefix, string(p))
	default:
		Logger.Infof("%v%v", self.prefix, string(p))
	}

	return len(p), nil
}
