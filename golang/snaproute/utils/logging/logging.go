//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/sysdCommonDefs"
	"log"
	"log/syslog"
	"models/objects"
	"os"
	"sysd"
	"time"
)

const (
	DB_CONNECT_TIME_INTERVAL   = 2
	DB_CONNECT_RETRY_LOG_COUNT = 100
)

type LoggerIntf interface {
	Crit(...interface{}) error
	Err(...interface{}) error
	Warning(...interface{}) error
	Alert(...interface{}) error
	Emerg(...interface{}) error
	Notice(...interface{}) error
	Info(...interface{}) error
	Println(...interface{}) error
	Debug(...interface{}) error
}

func ConvertLevelStrToVal(str string) sysdCommonDefs.SRDebugLevel {
	var val sysdCommonDefs.SRDebugLevel
	switch str {
	case "off":
		val = sysdCommonDefs.OFF
	case "crit":
		val = sysdCommonDefs.CRIT
	case "err":
		val = sysdCommonDefs.ERR
	case "warn":
		val = sysdCommonDefs.WARN
	case "alert":
		val = sysdCommonDefs.ALERT
	case "emerg":
		val = sysdCommonDefs.EMERG
	case "notice":
		val = sysdCommonDefs.NOTICE
	case "info":
		val = sysdCommonDefs.INFO
	case "debug":
		val = sysdCommonDefs.DEBUG
	case "trace":
		val = sysdCommonDefs.TRACE
	}
	return val
}

func ConvertLevelValToStr(val sysdCommonDefs.SRDebugLevel) string {
	var str string
	switch val {
	case sysdCommonDefs.OFF:
		str = "off"
	case sysdCommonDefs.CRIT:
		str = "crit"
	case sysdCommonDefs.ERR:
		str = "err"
	case sysdCommonDefs.WARN:
		str = "warn"
	case sysdCommonDefs.ALERT:
		str = "alert"
	case sysdCommonDefs.EMERG:
		str = "emerg"
	case sysdCommonDefs.NOTICE:
		str = "notice"
	case sysdCommonDefs.INFO:
		str = "info"
	case sysdCommonDefs.DEBUG:
		str = "debug"
	case sysdCommonDefs.TRACE:
		str = "trace"
	}
	return str
}

type Writer struct {
	SysLogger       *syslog.Writer
	nullLogger      *log.Logger
	dbHdl           redis.Conn
	MyComponentName string
	MyLogLevel      sysdCommonDefs.SRDebugLevel
	initialized     bool
	subSocket       *nanomsg.SubSocket
	socketCh        chan []byte
}

func NewLogger(name string, tag string, listenToConfig bool) (*Writer, error) {
	var err error
	srLogger := new(Writer)
	srLogger.MyComponentName = name
	srLogger.initialized = false

	srLogger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}
	// if SysLogger can't be initialized then send all logs to /dev/null
	devNull, err := os.Open(os.DevNull)
	if err == nil {
		srLogger.nullLogger = log.New(devNull, tag, log.Ldate|log.Ltime|log.Lshortfile)
	}

	srLogger.MyLogLevel = sysdCommonDefs.INFO
	// Read logging level from DB
	srLogger.readLogLevelFromDb()
	srLogger.initialized = true
	fmt.Println("Logging level ", srLogger.MyLogLevel, " set for ", srLogger.MyComponentName)
	if listenToConfig {
		go srLogger.ListenForLoggingNotifications()
	}
	return srLogger, err
}

func (logger *Writer) ReadComponentLoggingFromDb(dbHdl redis.Conn) error {
	logger.Info("Reading ComponentLogging")
	var dbObj objects.ComponentLogging
	objList, err := dbObj.GetAllObjFromDb(dbHdl)
	if err != nil {
		logger.Err("DB query failed for ComponentLogging config")
		return err
	}
	if objList != nil {
		for idx := 0; idx < len(objList); idx++ {
			obj := sysd.NewComponentLogging()
			dbObject := objList[idx].(objects.ComponentLogging)
			objects.ConvertsysdComponentLoggingObjToThrift(&dbObject, obj)
			if obj.Module == logger.MyComponentName {
				logger.MyLogLevel = ConvertLevelStrToVal(obj.Level)
				return nil
			}
		}
	}
	return nil
}

func (logger *Writer) UpdateComponentLoggingInDb() error {
	logger.Info("Reading ComponentLogging")
	dbObj := objects.ComponentLogging{
		Module: logger.MyComponentName,
		Level:  ConvertLevelValToStr(logger.MyLogLevel),
	}
	if logger.dbHdl == nil {
		logger.Info("Failed to update logging level for", logger.MyComponentName, "nil dbHdl")
		return errors.New("Uninitialized dbHdl")
	}
	err := dbObj.StoreObjectInDb(logger.dbHdl)
	if err != nil {
		logger.Info("Failed to update logging level for", logger.MyComponentName, "in db", err.Error())
		return err
	}
	return nil
}

func (logger *Writer) readLogLevelFromDb() error {
	var err error
	retryCount := 0
	ticker := time.NewTicker(DB_CONNECT_TIME_INTERVAL * time.Second)
	logger.dbHdl, err = redis.Dial("tcp", ":6379")
	if err != nil {
		for _ = range ticker.C {
			retryCount += 1
			logger.dbHdl, err = redis.Dial("tcp", ":6379")
			if err != nil {
				if retryCount%DB_CONNECT_RETRY_LOG_COUNT == 0 {
					logger.Err(fmt.Sprintln("Failed to dial out to Redis server. Ret    rying connection. Num retries = ", retryCount))
				}
			} else {
				break
			}
		}
	}

	if logger.dbHdl != nil {
		logger.ReadComponentLoggingFromDb(logger.dbHdl)
	}
	return nil
}

func (logger *Writer) SetLevel(level sysdCommonDefs.SRDebugLevel) error {
	logger.MyLogLevel = level
	logger.Debug(fmt.Sprintln("Changed logging level to: ", logger.MyLogLevel, " for ", logger.MyComponentName))
	return nil
}

func (logger *Writer) Crit(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.CRIT {
			return logger.SysLogger.Crit(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Err(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.ERR {
			return logger.SysLogger.Err(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Warning(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.WARN {
			return logger.SysLogger.Warning(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Alert(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.ALERT {
			return logger.SysLogger.Alert(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Emerg(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.EMERG {
			return logger.SysLogger.Emerg(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Notice(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.NOTICE {
			return logger.SysLogger.Notice(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Info(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.INFO {
			return logger.SysLogger.Info(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Println(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.INFO {
			return logger.SysLogger.Info(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Debug(message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.DEBUG {
			return logger.SysLogger.Debug(fmt.Sprintln(message))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintln(message))
	}
	return nil
}

func (logger *Writer) Write(message string) (int, error) {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.TRACE {
			n, err := logger.SysLogger.Write([]byte(message))
			return n, err
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(message)
	}
	return 0, nil
}

func (logger *Writer) Critf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.CRIT {
			return logger.SysLogger.Crit(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Errf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.ERR {
			return logger.SysLogger.Err(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Warningf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.WARN {
			return logger.SysLogger.Warning(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Alertf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.ALERT {
			return logger.SysLogger.Alert(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Emergf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.EMERG {
			return logger.SysLogger.Emerg(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Noticef(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.NOTICE {
			return logger.SysLogger.Notice(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Infof(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.INFO {
			return logger.SysLogger.Info(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Printf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.INFO {
			return logger.SysLogger.Info(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Debugf(format string, message ...interface{}) error {
	if logger.initialized {
		if logger.MyLogLevel >= sysdCommonDefs.DEBUG {
			return logger.SysLogger.Debug(fmt.Sprintf(format, message...))
		}
	} else if logger.nullLogger != nil {
		logger.nullLogger.Println(fmt.Sprintf(format, message))
	}
	return nil
}

func (logger *Writer) Close() error {
	var err error
	if logger.initialized {
		err = logger.SysLogger.Close()
	}
	logger.dbHdl.Close()
	logger = nil
	return err
}

func (logger *Writer) SetupSubSocket() error {
	var err error
	var socket *nanomsg.SubSocket
	if socket, err = nanomsg.NewSubSocket(); err != nil {
		logger.Err(fmt.Sprintf("Failed to create subscribe socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	if err = socket.Subscribe(""); err != nil {
		logger.Err(fmt.Sprintf("Failed to subscribe to \"\" on subscribe socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	if _, err = socket.Connect(sysdCommonDefs.PUB_SOCKET_ADDR); err != nil {
		logger.Err(fmt.Sprintf("Failed to connect to publisher socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	logger.Info(fmt.Sprintf("Connected to publisher socket %s", sysdCommonDefs.PUB_SOCKET_ADDR))
	if err = socket.SetRecvBuffer(1024 * 1024); err != nil {
		logger.Err(fmt.Sprintln("Failed to set the buffer size for subsriber socket %s, error:", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}
	logger.subSocket = socket
	logger.socketCh = make(chan []byte)
	return nil
}

func (logger *Writer) ProcessLoggingNotification(rxBuf []byte) error {
	var msg sysdCommonDefs.Notification
	err := json.Unmarshal(rxBuf, &msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmarshal logging notification: ", rxBuf))
		return err
	}
	if msg.Type == sysdCommonDefs.G_LOG {
		var gLog sysdCommonDefs.GlobalLogging
		err = json.Unmarshal(msg.Payload, &gLog)
		if err != nil {
			logger.Err(fmt.Sprintln("Unable to unmarshal global logging notification: ", msg.Payload))
			return err
		}
		logger.SetLevel(gLog.Level)
		logger.UpdateComponentLoggingInDb()
	}
	if msg.Type == sysdCommonDefs.C_LOG {
		var cLog sysdCommonDefs.ComponentLogging
		err = json.Unmarshal(msg.Payload, &cLog)
		if err != nil {
			logger.Err(fmt.Sprintln("Unable to unmarshal component logging notification: ", msg.Payload))
			return err
		}
		if cLog.Name == logger.MyComponentName {
			logger.SetLevel(cLog.Level)
		}
	}
	return nil
}

func (logger *Writer) ProcessLogNotifications() error {
	for {
		select {
		case rxBuf := <-logger.socketCh:
			if rxBuf != nil {
				logger.ProcessLoggingNotification(rxBuf)
			}
		}
	}
	return nil
}

func (logger *Writer) ListenForLoggingNotifications() error {
	err := logger.SetupSubSocket()
	if err != nil {
		logger.Err(fmt.Sprintln("Failed to subscribe to logging notifications"))
		return err
	}
	go logger.ProcessLogNotifications()
	for {
		rxBuf, err := logger.subSocket.Recv(0)
		if err != nil {
			logger.Err(fmt.Sprintln("Recv on logging subscriber socket failed with error:", err))
			continue
		}
		logger.socketCh <- rxBuf
	}
	logger.Info(fmt.Sprintln("Existing logging config lister"))
	return nil
}
