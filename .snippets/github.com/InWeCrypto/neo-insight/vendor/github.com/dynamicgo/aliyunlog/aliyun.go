package aliyunlog

import (
	"fmt"
	"sync"
	"time"

	sls "github.com/dynamicgo/aliyun-log-go-sdk"
	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/gogo/protobuf/proto"
)

type aliyunLog struct {
	topic    string
	source   string
	mq       chan []*sls.LogContent
	logstore *sls.LogStore
}

func newAliyunLog(topic string, source string, logstore *sls.LogStore, cached int) *aliyunLog {
	logger := &aliyunLog{
		topic:    topic,
		source:   source,
		logstore: logstore,
		mq:       make(chan []*sls.LogContent, cached),
	}

	go logger.runLoop()

	return logger
}

func (logger *aliyunLog) runLoop() {
	for content := range logger.mq {

		group := &sls.LogGroup{
			Topic:  proto.String(logger.topic),
			Source: proto.String(logger.source),
			Logs: []*sls.Log{
				&sls.Log{
					Contents: content,
					Time:     proto.Uint32(uint32(time.Now().Unix())),
				},
			},
		}

		if err := logger.logstore.PutLogs(group); err != nil {
			fmt.Printf("logstore put logs err, %s\n", err)
			continue
		}
	}
}

func (logger *aliyunLog) GetName() string {
	return logger.topic
}

func (logger *aliyunLog) Trace(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Trace"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) TraceF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Trace"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

func (logger *aliyunLog) Debug(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Debug"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) DebugF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Debug"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

func (logger *aliyunLog) Info(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Info"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) InfoF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Info"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

func (logger *aliyunLog) Warn(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Warn"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) WarnF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Warn"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

func (logger *aliyunLog) Error(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Error"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) ErrorF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Error"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

func (logger *aliyunLog) Fatal(args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Fatal"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprint(args...)),
		},
	}
}

func (logger *aliyunLog) FatalF(format string, args ...interface{}) {
	logger.mq <- []*sls.LogContent{

		&sls.LogContent{
			Key:   proto.String("Level"),
			Value: proto.String("Fatal"),
		},
		&sls.LogContent{
			Key:   proto.String("Content"),
			Value: proto.String(fmt.Sprintf(format, args...)),
		},
	}
}

type aliyunLogHub struct {
	project   *sls.LogProject
	logstore  *sls.LogStore
	source    string
	loggermap map[string]*aliyunLog
	mutex     sync.RWMutex
	cached    int
}

// NewAliyunBackend create new aliyun log-hub backend
func NewAliyunBackend(cnf *config.Config) (slf4go.LoggerFactory, error) {

	project := &sls.LogProject{
		Name:            cnf.GetString("slf4go.aliyun.project", "xxxx"),
		Endpoint:        cnf.GetString("slf4go.aliyun.endpoint", "xxxxx"),
		AccessKeyID:     cnf.GetString("slf4go.aliyun.accesskey.id", "xxxxx"),
		AccessKeySecret: cnf.GetString("slf4go.aliyun.accesskey.secret", "xxxxx"),
	}

	logstore, err := project.GetLogStore(cnf.GetString("slf4go.aliyun.logstore", "xxxx"))

	if err != nil {
		return nil, err
	}

	return &aliyunLogHub{
		project:   project,
		logstore:  logstore,
		source:    cnf.GetString("slf4go.aliyun.source", "xxxxx"),
		loggermap: make(map[string]*aliyunLog),
		cached:    int(cnf.GetInt64("slf4go.aliyun.cached", 0)),
	}, nil
}

func (loghub *aliyunLogHub) GetLogger(name string) slf4go.Logger {

	loghub.mutex.RLock()

	logger, ok := loghub.loggermap[name]

	loghub.mutex.RUnlock()

	if !ok {
		loghub.mutex.Lock()
		logger = newAliyunLog(name, loghub.source, loghub.logstore, loghub.cached)
		loghub.loggermap[name] = logger
		loghub.mutex.Unlock()
	}

	return logger
}
