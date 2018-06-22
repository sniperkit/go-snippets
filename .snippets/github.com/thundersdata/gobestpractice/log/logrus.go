/*
Package log show best practice about logging. We choose logrus as our logger.
*/
package log

import (
	"gobestpractice/conf"

	"fmt"
	"runtime"
	"strings"

	"io/ioutil"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	kl "github.com/tracer0tong/kafkalogrus"
)

var logger *log.Entry

func init() {
	logger = log.WithFields(log.Fields{
		"project": "gobestpractice",
	})

	if conf.Config.Profile == "production" {
		// Only log the Info severity or above to stdout.
		log.SetLevel(log.InfoLevel)

		// Add logrusKafka hook to output log to kafka.
		hook, err := kl.NewKafkaLogrusHook(
			conf.Config.AppName,
			[]log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel},
			&log.JSONFormatter{},
			conf.Config.KafkaBroker,
			"test",
			true)
		if err != nil {
			log.Error(errors.Annotate(err, "Cannot create KafkaHook."))
		} else {
			log.AddHook(hook)
			// Disable all output to stdout/err when hook is ready.
			log.SetOutput(ioutil.Discard)
		}
	}
}

// Logger is the log entry.
func Logger() *log.Entry {
	if conf.Config.Profile == "production" {
		return logger
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	return log.WithField("source", fmt.Sprintf("%s:%d", file, line))
}
