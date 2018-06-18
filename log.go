package toolkit

import (
	"fmt"
	//"io"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eaciit/cast"
)

type LogEngine struct {
	LogToStdOut     bool
	LogToFile       bool
	Path            string
	FileNamePattern string
	UseDateFormat   string
	logInfo         *log.Logger
	logWarn         *log.Logger
	logError        *log.Logger
	//logFile         *log.Logger
	//logFileHandler  *os.File
}

func NewLog(toStdOut bool, toFile bool, path string, fileNamePattern string, useDateFormate string) (*LogEngine, error) {
	var e error
	l := new(LogEngine)
	l.LogToStdOut = toStdOut
	l.LogToFile = toFile
	l.Path = path
	l.FileNamePattern = fileNamePattern
	l.UseDateFormat = useDateFormate
	//l.logger = log.New(out, prefix, flag)

	e = l.initLogger()
	if e != nil {
		return nil, e
	}
	return l, nil
}

func (l *LogEngine) initLogger() error {
	//var e error = nil
	if l.LogToStdOut {
		l.logError = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime)
		l.logInfo = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
		l.logWarn = log.New(os.Stdout, "WARNING ", log.Ldate|log.Ltime)
	}
	return nil
}

func (l *LogEngine) AddLog(msg string, logtype string) error {
	var e error
	logtype = strings.ToUpper(logtype) + " "

	if l.LogToStdOut {
		if logtype == "ERROR " {
			l.logError.Println(msg)
		} else if logtype == "WARNING " {
			l.logWarn.Println(msg)
		} else {
			l.logInfo.Println(msg)
		}
		if e != nil {
			return errors.New("Log.AddLog Error: " + e.Error())
		}
	}

	if l.LogToFile {
		filename := l.FileNamePattern
		if l.UseDateFormat != "" && strings.Contains(l.FileNamePattern, "%s") {
			filename = fmt.Sprintf(l.FileNamePattern, cast.Date2String(time.Now(), l.UseDateFormat))
		}
		filename = filepath.Join(l.Path, filename)
		f, e := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if e != nil {
			return errors.New("Log.AddLog Error: " + e.Error())
		}
		defer f.Close()
		logFile := log.New(f, logtype, log.Ldate|log.Ltime)
		logFile.Println(msg)
	}

	return nil
}

func (l *LogEngine) Info(msg string) error {
	return l.AddLog(msg, "INFO")
}

func (l *LogEngine) Error(msg string) error {
	return l.AddLog(msg, "ERROR")
}

func (l *LogEngine) Warning(msg string) error {
	return l.AddLog(msg, "WARNING")
}

func (l *LogEngine) Infof(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "INFO")
}

func (l *LogEngine) Errorf(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "ERROR")
}

func (l *LogEngine) Warningf(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "WARNING")
}

func (l *LogEngine) Close() {
	//l.logFileHandler.Close()
}
