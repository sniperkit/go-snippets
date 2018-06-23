package logger

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/fatih/color"
	"github.com/homebot/core/utils"
	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
	"google.golang.org/grpc"
)

// Logger wraps logger implementations
type Logger interface {
	// Default logging methods
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})

	// Sub logger creation
	WithService(string) Logger
	WithServiceType(string) Logger
	WithResource(string) Logger
	WithIdentity(string) Logger
	WithDetails(utils.ValueMap) Logger
}

type Forwarder chan *insightV1.LogMessage

func (f Forwarder) forward(wg *sync.WaitGroup, conn *grpc.ClientConn) {
	defer wg.Done()
	defer conn.Close()

	client := insightV1.NewLogSinkClient(conn)

	for msg := range f {
		if _, err := client.AddMessage(context.Background(), msg); err != nil {
			log.Printf("failed to forward log message: %s\n", err)
		}
	}
}

// InsightLogger is a log.Logger implementation that sends
// log messages to an insight server
type InsightLogger struct {
	minForwardLevel insightV1.SeverityLevel
	minConsoleLevel insightV1.SeverityLevel

	forwarder Forwarder
	wg        *sync.WaitGroup

	serviceType string
	service     string
	resource    string
	identity    string
	details     utils.ValueMap
}

// Option is an option for an insight logger
type Option func(i *InsightLogger) error

// NewInsightLogger creates a new insight logger instance
func NewInsightLogger(opts ...Option) (*InsightLogger, error) {
	logger := &InsightLogger{}

	for _, fn := range opts {
		if err := fn(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

// Clone clones the logger and adds additional options
func (log *InsightLogger) Clone(opts ...Option) Logger {
	n := &InsightLogger{
		minConsoleLevel: log.minConsoleLevel,
		minForwardLevel: log.minForwardLevel,
		serviceType:     log.serviceType,
		service:         log.service,
		resource:        log.resource,
		identity:        log.identity,
		details:         log.details,
	}

	for _, fn := range opts {
		fn(n)
	}

	return n
}

// Debugf logs a message at DEBUG level
func (log *InsightLogger) Debugf(msg string, args ...interface{}) {
	log.log(insightV1.SeverityLevel_DEBUG, msg, args...)
}

// Infof logs a message at INFO level
func (log *InsightLogger) Infof(msg string, args ...interface{}) {
	log.log(insightV1.SeverityLevel_INFO, msg, args...)
}

// Warnf logs a message at WARNING level
func (log *InsightLogger) Warnf(msg string, args ...interface{}) {
	log.log(insightV1.SeverityLevel_WARNING, msg, args...)
}

// Errorf logs a message at ERROR level
func (log *InsightLogger) Errorf(msg string, args ...interface{}) {
	log.log(insightV1.SeverityLevel_ERROR, msg, args...)
}

func (log *InsightLogger) log(level insightV1.SeverityLevel, msg string, args ...interface{}) {
	if level >= log.minConsoleLevel {
		log.printConsole(level, msg, args...)
	}

	if level >= log.minForwardLevel {
		log.forward(level, msg, args...)
	}
}

func (log *InsightLogger) printConsole(lvl insightV1.SeverityLevel, msg string, args ...interface{}) {
	var attr []color.Attribute

	switch lvl {
	case insightV1.SeverityLevel_DEBUG:
	case insightV1.SeverityLevel_INFO:
		attr = append(attr, color.FgHiWhite)
	case insightV1.SeverityLevel_WARNING:
		attr = append(attr, color.FgHiYellow)
	case insightV1.SeverityLevel_ERROR:
		attr = append(attr, color.FgHiRed)
	}

	printer := color.New(attr...)

	prefix := " "

	if log.identity != "" {
		prefix += "(" + log.identity + ") "
	}

	if log.resource != "" {
		prefix += log.resource + ": "
	}

	printer.Printf("%s [%s]%s%s\n", time.Now().Format(time.RFC822), level(lvl), prefix, fmt.Sprintf(msg, args...))
}

func (log *InsightLogger) forward(level insightV1.SeverityLevel, msg string, args ...interface{}) {
	if log.forwarder != nil {
		now := ptypes.TimestampNow()

		log.forwarder <- &insightV1.LogMessage{
			CreatedTime:    now,
			Details:        log.details.ToProto(),
			Identity:       log.identity,
			Service:        log.service,
			ServiceType:    log.serviceType,
			TargetResource: log.resource,
			Severity:       level,
			Message:        fmt.Sprintf(msg, args...),
		}
	}
}

// WithServiceType clones the logger and adds a service type
func (log *InsightLogger) WithServiceType(s string) Logger {
	return log.Clone(WithServiceType(s))
}

// WithService clones the logger and adds a service name
func (log *InsightLogger) WithService(s string) Logger {
	return log.Clone(WithService(s))
}

// WithResource clones the logger and adds a resource name
func (log *InsightLogger) WithResource(r string) Logger {
	return log.Clone(WithResource(r))
}

// WithIdentity clones the logger and adds an identity
func (log *InsightLogger) WithIdentity(i string) Logger {
	return log.Clone(WithIdentity(i))
}

// WithDetails clones the logger and adds additional details
func (log *InsightLogger) WithDetails(d utils.ValueMap) Logger {
	return log.Clone(WithDetails(d))
}

type NopLogger struct{}

func (NopLogger) Debugf(string, ...interface{})     {}
func (NopLogger) Infof(string, ...interface{})      {}
func (NopLogger) Warnf(string, ...interface{})      {}
func (NopLogger) Errorf(string, ...interface{})     {}
func (NopLogger) WithServiceType(string) Logger     { return NopLogger{} }
func (NopLogger) WithService(string) Logger         { return NopLogger{} }
func (NopLogger) WithResource(string) Logger        { return NopLogger{} }
func (NopLogger) WithIdentity(string) Logger        { return NopLogger{} }
func (NopLogger) WithDetails(utils.ValueMap) Logger { return NopLogger{} }
