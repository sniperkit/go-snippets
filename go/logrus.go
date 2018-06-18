package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

func LogToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err == nil {
		logrus.SetOutput(file)
	}
	return err
}

func main() {
	if err := LogToFile("/tmp/logrustest.log"); err != nil {
		logrus.Fatal("Failed to log to file")
	}

	for i := 0; i < 3; i++ {
		logrus.Warn(fmt.Sprintf("%d. This is a test log..", i+1))
		time.Sleep(1 * time.Second)
	}
	logrus.Info("Logrus test done.")
	fmt.Println("Logrus test done.")
}
