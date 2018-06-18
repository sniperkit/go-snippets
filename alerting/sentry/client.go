package sentry

import (
	raven "github.com/getsentry/raven-go"

	"github.com/go-zero-boilerplate/facades/logging"
)

var client *raven.Client

//GetClient returns the singleton sentry/raven instance
func GetClient() *raven.Client {
	return client
}

//CaptureErrorPacketAndWait will capture an error with stack trace and allow modification of the Packet before sending
func CaptureErrorPacketAndWait(message string, err error, modifyPacketBeforeSend func(p *raven.Packet), sentryTags map[string]string) {
	newPacket := raven.NewPacket(message, raven.NewException(err, raven.NewStacktrace(1, 3, nil)))
	if modifyPacketBeforeSend != nil {
		modifyPacketBeforeSend(newPacket)
	}
	_, ch := client.Capture(newPacket, sentryTags)
	if err = <-ch; err != nil {
		logging.Get().WithError(err).Error("Failed to deliver Sentry Packet")
	}
}

//InitClient creates the sentry/raven singleton
func InitClient(dsn string, tags map[string]string) {
	ravenClient, err := raven.NewClient(dsn, tags)
	if err != nil {
		panic(err)
	}

	client = ravenClient
}
