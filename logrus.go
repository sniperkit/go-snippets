package main

import (
    "github.com/Sirupsen/logrus"
    "gopkg.in/yaml.v2"
    "github.com/rifflock/lfshook"

    "os"
    "fmt"
    "runtime"
    "strings"
    "io"
    "io/ioutil"
//    "path/filepath"
    "reflect"
)

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		//slash := strings.LastIndex(file, "/")
		//if slash >= 0 {
		//	file = file[slash+1:]
		//}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

type Config struct {
    Logging struct {
        Level string
        Format string
    }
}

var log = logrus.New()

func init() {
    file, err := os.OpenFile("./config.yaml", os.O_RDONLY, 0666)
    defer file.Close()
    if err != nil {
        fmt.Println("Cannot open the file")
        os.Exit(1)
    }
    chunks := make([]byte, 0, 1024)
    buffer := make([]byte, 1024)
    for {
        n, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if err != nil && err == io.EOF || 0 == n {
            fmt.Println(err)
            break
        }
        chunks = append(chunks, buffer[:n]...)
        //fmt.Println(string(chunks))
    }
    fmt.Println(len(chunks))
    fmt.Println(string(chunks))
    data, err := ioutil.ReadFile("./config.yaml")
    fmt.Println(reflect.TypeOf(data))
    config := Config{}
    err = yaml.Unmarshal(chunks, &config)
    if err != nil {
        fmt.Printf("error: %v", err)
        os.Exit(1)
    }
    fmt.Printf("--- logger:\n%v\n\n", config.Logging.Level)
    fmt.Printf("--- logger:\n%v\n\n", config.Logging.Format)
    switch strings.ToLower(config.Logging.Format) {
    case "json":
        log.Formatter = &logrus.JSONFormatter{}
        break
    case "text":
        log.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: false, DisableSorting: false}
        break
    default:
        log.Formatter = &logrus.TextFormatter{DisableTimestamp: true, FullTimestamp: true}
    }

    switch strings.ToLower(config.Logging.Level) {
    case "debug":
        log.Level = logrus.DebugLevel
        break
    case "info":
        log.Level = logrus.InfoLevel
        break
    default:
        log.Level = logrus.InfoLevel
        break
    }

    logFile, _ := os.OpenFile("./test.log", os.O_RDWR | os.O_APPEND | os.O_CREATE | os.O_SYNC, 0666)
    log.Out = logFile
    log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
        logrus.DebugLevel : "/var/log/info.log",
        logrus.ErrorLevel: "/var/log/error.log",
    }))
}

func main() {
    log.WithFields(logrus.Fields{
        "animal": "walrus",
        "number": 8,
        "file": fileInfo(1),
    }).Debugf("Started observing beach")

    log.WithFields(logrus.Fields{
       "animal": "hello",
       "number": 10,
       "file": fileInfo(1),
    }).Debugf("Started2")
}
