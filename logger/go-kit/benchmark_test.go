package logger

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/albert-widi/go_common/errors"
	klog "github.com/go-kit/kit/log"
)

var log *Logger
var llog klog.Logger

func init() {
	log = fake()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	f, _ := os.Open("/dev/null")
	logrus.StandardLogger().Out = f
	llog = klog.NewJSONLogger(f)
}

func BenchmarkSimpleLogger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Error("Something is wrong", "this is so wrong")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.WithFields(Fields{"field1": "value1"}).Info("This is a info with fields")
	}
}

func BenchmarkLoggerWithLongFields(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.WithFields(Fields{"field1": "value1", "field2": "value2", "field3": "value4", "field5": "value5", "field6": "value6",
			"field7": "value7", "field8": "value8", "field9": "value9", "field10": "value10",
			"field11": "this is probably a very long text that need to be logged and need some consideration how should we log this",
			"field12": 13.10, "field13": 100000000}).
			Info("This is a info with fields")
	}
}

func BenchmarkErrors(b *testing.B) {
	err := errors.New("This is new errors", errors.Fields{"err1": "err_value_1"})
	for n := 0; n < b.N; n++ {
		log.Errors(err)
	}
}

func BenchmarkLongErrorsFields(b *testing.B) {
	err := errors.New("This is new errors", errors.Fields{"err1": "err_value_1", "order_id": "XWYZ012312831",
		"text": "this is a very long text that need to be avoided when we are writing error context", "transaction_id": 128392,
		"request_id": 12382931123819, "float": 102392.1293291, "text2": "this is another text that is long enough",
		"err8": "err_value8", "err9": "err_value9", "another_text": "another text that is long enough for an error context"})
	for n := 0; n < b.N; n++ {
		log.Errors(err)
	}
}

func BenchmarkErrorsWithFields(b *testing.B) {
	err := errors.New("This is new errors", errors.Fields{"err1": "err_value_1"})
	for n := 0; n < b.N; n++ {
		log.WithFields(Fields{"field1": "value1"}).Errors(err)
	}
}

func BenchmarkGokitLog(b *testing.B) {
	for n := 0; n < b.N; n++ {
		llog.Log("Something is wrong")
	}
}

// func BenchmarkSimpleLogrus(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		logrus.Error("Something is wrong")
// 	}
// }

// func BenchmarkLogrusWithFields(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		logrus.WithFields(logrus.Fields{"field1": "value1"}).Info("This is a info with fields")
// 	}
// }
