package logger

import (
	"sync"

	"github.com/homebot/core/utils"
	"google.golang.org/grpc"
)

// WithForwarder configues a log message forwarder
func WithForwarder(conn *grpc.ClientConn) Option {
	return func(i *InsightLogger) error {
		forwarder := make(Forwarder, 1000)

		i.forwarder = forwarder
		i.wg = &sync.WaitGroup{}
		i.wg.Add(1)

		go i.forwarder.forward(i.wg, conn)

		return nil
	}
}

// WithServiceType adds the service type to log messages
func WithServiceType(s string) Option {
	return func(i *InsightLogger) error {
		i.serviceType = s
		return nil
	}
}

// WithService adds the service name to log messages
func WithService(s string) Option {
	return func(i *InsightLogger) error {
		i.service = s
		return nil
	}
}

// WithResource adds a resource name to log messages
func WithResource(s string) Option {
	return func(i *InsightLogger) error {
		i.resource = s
		return nil
	}
}

// WithIdentity adds an identity name to log messages
func WithIdentity(s string) Option {
	return func(i *InsightLogger) error {
		i.identity = s
		return nil
	}
}

// WithDetails adds additional key-value pairs to log messages
func WithDetails(details utils.ValueMap) Option {
	return func(i *InsightLogger) error {
		i.details = details
		return nil
	}
}
