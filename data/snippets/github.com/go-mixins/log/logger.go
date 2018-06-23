package log

// M defines log entry context
type M map[string]interface{}

// Logger provides generic logging functions
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
}

// ContextLogger adds some context to log entries
type ContextLogger interface {
	WithContext(M) ContextLogger
	Logger
}

//go:generate moq -out mock/logger.go -pkg mock . ContextLogger
