package nsq

import (
	"strings"

	nsq "github.com/nsqio/go-nsq"
)

// Logger should implement basic log output functions for all levels used in NSQ
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

// Output implements stdlib log.Logger.Output using the underlying logger
func (l *config) Output(_ int, s string) error {
	s = strings.TrimSpace(s)
	var level string
	if len(s) >= 3 {
		level, s = s[:3], s[3:]
	}
	switch level {
	case nsq.LogLevelError.String():
		if !l.keep404 && strings.Contains(s, "TOPIC_NOT_FOUND") {
			return nil
		}
		l.logger.Error(s)
	case nsq.LogLevelWarning.String():
		l.logger.Warn(s)
	case nsq.LogLevelDebug.String():
		l.logger.Debug(s)
	default:
		l.logger.Info(s)
	}
	return nil
}
