package log

import "strings"

type Level int

const (
	PANIC Level = iota
	FATAL
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

func (self Level) String() string {
	switch self {
	case PANIC:
		return `panic`
	case CRITICAL:
		return `critical`
	case ERROR:
		return `error`
	case WARNING:
		return `warning`
	case NOTICE:
		return `notice`
	case INFO:
		return `info`
	case DEBUG:
		return `debug`
	default:
		return ``
	}
}

func GetLevel(level string) Level {
	switch strings.ToLower(level) {
	case `panic`:
		return PANIC
	case `fatal`:
		return FATAL
	case `critical`, `crit`:
		return CRITICAL
	case `error`, `err`:
		return ERROR
	case `warning`, `warn`:
		return WARNING
	case `notice`:
		return NOTICE
	case `info`:
		return INFO
	case `debug`:
		return DEBUG
	default:
		return DEBUG
	}
}
