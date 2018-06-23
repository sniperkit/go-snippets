package sentry

import (
	"fmt"
	"runtime"
	"strings"

	apex "github.com/francoishill/log"
	"github.com/francoishill/log/handlers/level"
	raven "github.com/getsentry/raven-go"
)

//NewLogHandler returns a new logger that logs to sentry
func NewLogHandler(client *raven.Client, logLevel apex.Level) apex.Handler {
	sentryLogger := &logger{client: client}
	return level.New(sentryLogger, logLevel)
}

type logger struct {
	client *raven.Client
}

var severityMapping = [...]raven.Severity{
	apex.DebugLevel:     raven.DEBUG,
	apex.InfoLevel:      raven.INFO,
	apex.NoticeLevel:    raven.WARNING,
	apex.WarnLevel:      raven.WARNING,
	apex.ErrorLevel:     raven.ERROR,
	apex.CriticalLevel:  raven.FATAL,
	apex.AlertLevel:     raven.FATAL,
	apex.EmergencyLevel: raven.FATAL,
}

func (l *logger) HandleLog(e *apex.Entry) error {
	tags := map[string]string{}

	var interfaces []raven.Interface
	if errInterface, hasErr := e.Fields["error"]; hasErr {
		if err, isErrType := errInterface.(error); isErrType {
			// interfaces = append(interfaces, raven.NewException(err, raven.NewStacktrace(0, 3, l.client.IncludePaths())))
			interfaces = append(interfaces, raven.NewException(err, newAutoStacktrace(3, l.client.IncludePaths())))
		}
	}

	for key, val := range e.Fields {
		if key == "error" || key == "stack" {
			//error is added above and stack is not reported to sentry here, the newAutoStacktrace is used above
			continue
		}
		tags[key] = fmt.Sprintf("%v", val)
	}

	packet := raven.NewPacket(e.Message, interfaces...)
	packet.Init("") //will setup default fields like project, culprit, etc
	packet.Level = severityMapping[e.Level]

	go l.client.Capture(packet, tags)

	return nil
}

//newAutoStacktrace is a slight variation of raven.NewStacktrace (skipping some frames based on file path)
func newAutoStacktrace(context int, appPackagePrefixes []string) *raven.Stacktrace {
	var frames []*raven.StacktraceFrame
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		ignoredPaths := []string{
			"github.com/go-zero-boilerplate/facades/alerting/sentry/logger.go",
			"github.com/francoishill/log/handlers/multi/multi.go",
			"github.com/francoishill/log/handlers/level/level.go",
			"github.com/francoishill/log/logger.go",
			"github.com/francoishill/log/entry.go",
		}

		mustIgnore := false
		for _, i := range ignoredPaths {
			if strings.Contains(strings.ToLower(file), strings.ToLower(i)) {
				mustIgnore = true
				break
			}
		}
		if mustIgnore {
			continue
		}

		frame := raven.NewStacktraceFrame(pc, file, line, context, appPackagePrefixes)
		if frame != nil {
			frames = append(frames, frame)
		}
	}
	// If there are no frames, the entire stacktrace is nil
	if len(frames) == 0 {
		return nil
	}
	// Optimize the path where there's only 1 frame
	if len(frames) == 1 {
		return &raven.Stacktrace{frames}
	}
	// Sentry wants the frames with the oldest first, so reverse them
	for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
		frames[i], frames[j] = frames[j], frames[i]
	}
	return &raven.Stacktrace{frames}
}
