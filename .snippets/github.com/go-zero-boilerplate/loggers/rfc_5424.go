package loggers

//LoggerRFC5424 follows the RFC 5424 format: https://tools.ietf.org/html/rfc5424
type LoggerRFC5424 interface {
	Emergency(format string, params ...interface{})
	Alert(format string, params ...interface{})
	Critical(format string, params ...interface{})
	Error(format string, params ...interface{})
	Warn(format string, params ...interface{})
	Notice(format string, params ...interface{})
	Info(format string, params ...interface{})
	Debug(format string, params ...interface{})
}
