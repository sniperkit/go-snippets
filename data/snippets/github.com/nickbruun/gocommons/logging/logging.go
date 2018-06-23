// Package logging provides caller-annotated logrus logging.
package logging

import (
	log "github.com/Sirupsen/logrus"
	"runtime"
)

// Get a caller-identifying logger.
func identified(logger *log.Logger, skip int) *log.Entry {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return log.NewEntry(logger)
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return logger.WithFields(log.Fields{
			"file": file,
			"line": line,
		})
	}

	return logger.WithFields(log.Fields{
		"file": file,
		"line": line,
		"func": fn.Name(),
	})
}

func Debug(args ...interface{}) {
	identified(log.StandardLogger(), 2).Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Debugf(format, args...)
}

func Info(args ...interface{}) {
	identified(log.StandardLogger(), 2).Info(args...)
}

func Infof(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Infof(format, args...)
}

func Print(args ...interface{}) {
	identified(log.StandardLogger(), 2).Print(args...)
}

func Printf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Printf(format, args...)
}

func Warn(args ...interface{}) {
	identified(log.StandardLogger(), 2).Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Warnf(format, args...)
}

func Warning(args ...interface{}) {
	identified(log.StandardLogger(), 2).Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Warningf(format, args...)
}

func Error(args ...interface{}) {
	identified(log.StandardLogger(), 2).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Errorf(format, args...)
}

func Panic(args ...interface{}) {
	identified(log.StandardLogger(), 2).Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	identified(log.StandardLogger(), 2).Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	identified(log.StandardLogger(), 2).Fatalf(format, args...)
}
