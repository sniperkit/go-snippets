package slf4go

// Logger slf4go facade interface
type Logger interface {
	GetName() string

	Trace(args ...interface{})

	TraceF(format string, args ...interface{})

	Debug(args ...interface{})

	DebugF(format string, args ...interface{})

	Info(args ...interface{})

	InfoF(format string, args ...interface{})

	Warn(args ...interface{})

	WarnF(format string, args ...interface{})

	Error(args ...interface{})

	ErrorF(format string, args ...interface{})

	Fatal(args ...interface{})

	FatalF(format string, args ...interface{})
}

// LoggerFactory logger's factory interface
type LoggerFactory interface {
	GetLogger(name string) Logger
}

// Logger level
const (
	Trace = 1 << iota
	Debug
	Info
	Warn
	Error
	Fatal
)

type loggerWrapper struct {
	impl Logger
}

func (logger *loggerWrapper) GetName() string {
	return logger.impl.GetName()
}

func (logger *loggerWrapper) Trace(args ...interface{}) {

	if (level & Trace) == Trace {
		logger.impl.Trace(args...)
	}
}

func (logger *loggerWrapper) TraceF(format string, args ...interface{}) {
	if (level & Trace) == Trace {
		logger.impl.TraceF(format, args...)
	}
}

func (logger *loggerWrapper) Debug(args ...interface{}) {
	if (level & Debug) == Debug {
		logger.impl.Debug(args...)
	}
}

func (logger *loggerWrapper) DebugF(format string, args ...interface{}) {
	if (level & Debug) == Debug {
		logger.impl.DebugF(format, args...)
	}
}

func (logger *loggerWrapper) Info(args ...interface{}) {
	if (level & Info) == Info {
		logger.impl.Info(args...)
	}
}

func (logger *loggerWrapper) InfoF(format string, args ...interface{}) {
	if (level & Info) == Info {
		logger.impl.InfoF(format, args...)
	}
}

func (logger *loggerWrapper) Warn(args ...interface{}) {
	if (level & Warn) == Warn {
		logger.impl.Warn(args...)
	}
}

func (logger *loggerWrapper) WarnF(format string, args ...interface{}) {
	if (level & Warn) == Warn {
		logger.impl.WarnF(format, args...)
	}
}

func (logger *loggerWrapper) Error(args ...interface{}) {
	if (level & Error) == Error {
		logger.impl.Error(args...)
	}
}

func (logger *loggerWrapper) ErrorF(format string, args ...interface{}) {
	if (level & Error) == Error {
		logger.impl.ErrorF(format, args...)
	}
}

func (logger *loggerWrapper) Fatal(args ...interface{}) {
	if (level & Fatal) == Fatal {
		logger.impl.Fatal(args...)
	}
}

func (logger *loggerWrapper) FatalF(format string, args ...interface{}) {
	if (level & Fatal) == Fatal {
		logger.impl.FatalF(format, args...)
	}
}
