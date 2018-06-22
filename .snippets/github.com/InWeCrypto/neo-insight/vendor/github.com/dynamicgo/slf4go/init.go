package slf4go

var loggerFactory LoggerFactory = newNativeLoggerFactory()
var level = Trace | Debug | Info | Warn | Error | Fatal

// Backend set new slf4go backend logger factory
func Backend(factory LoggerFactory) {
	if factory == nil {
		panic("factory can't be nil")
	}

	loggerFactory = factory
}

// Get get/create new logger by name
func Get(name string) Logger {
	return &loggerWrapper{impl: loggerFactory.GetLogger(name)}
}

// SetLevel set logger level
func SetLevel(l int) {
	level = l
}

// GetLevel get logger level
func GetLevel() int {
	return level
}
