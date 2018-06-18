package logging

type LogTracer interface {
	Stop(errPtr *error)
}
