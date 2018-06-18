package loggers

//LoggerRFC5424Simple is the same as LoggerRFC5424 but without params and still follows the RFC 5424 format: https://tools.ietf.org/html/rfc5424
type LoggerRFC5424Simple interface {
	Emergency(s string)
	Alert(s string)
	Critical(s string)
	Error(s string)
	Warn(s string)
	Notice(s string)
	Info(s string)
	Debug(s string)
}
