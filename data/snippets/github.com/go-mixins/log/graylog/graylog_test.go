// +build integration

package graylog_test

import (
	"testing"
	"time"

	"github.com/go-mixins/log/graylog"
)

func TestLogOutput(t *testing.T) {
	logger := graylog.New("localhost:12201")
	const delay = 10 * time.Millisecond
	for i := 0; i < 1000; i++ {
		time.Sleep(delay)
		logger.Info("test message")
	}
}
