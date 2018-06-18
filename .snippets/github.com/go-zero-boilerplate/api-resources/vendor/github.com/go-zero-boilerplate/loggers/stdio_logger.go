package loggers

//LoggerStdIO implements the two StdIO levels of processes (Stdout and Stderr)
type LoggerStdIO interface {
	Err(format string, params ...interface{})
	Out(format string, params ...interface{})
}
