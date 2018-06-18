package main

import (
	"log/syslog"
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"fmt"
	"time"
)

func main()  {

	log       := logrus.New()
	hook, err := lSyslog.NewSyslogHook("udp", "stat.suretly.io:514", syslog.LOG_INFO, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Hooks.Add(hook)

	for i:=0;i<1000;i++ {
		log.WithFields(logrus.Fields{
			"name": "Slartibartfast11",
			"age":  42,
		}).Info("Hello Dent, Arthur Dent!")
		time.Sleep(time.Second)
	}


	hook.Writer.Close()
}
