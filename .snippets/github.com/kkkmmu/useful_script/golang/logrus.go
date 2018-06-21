package main

import (
	"github.com/Sirupsen/logrus"
	"os"
)

func main() {
	logrus.Debug("Hello world!")
	logrus.Info("Hello world!")
	logrus.Warn("Hello world!")
	logrus.Error("Hello world!")
	//logrus.Fatal("Hello world!")

	loger := logrus.New()
	loger.Out = os.Stdout
	loger.Debug("Hello world!")
	loger.Info("Hello world!")
	loger.Print("Hello world!")
	loger.Warn("Hello world!")
	loger.Formatter = &logrus.JSONFormatter{}
	loger.Debug("Hello world!")
	loger.Info("Hello world!")
	loger.Print("Hello world!")
	loger.Warn("Hello world!")
	//loger.Panic("Hello world!")
}

func init() {
	//logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
