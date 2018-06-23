package main

import (
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	hook, err := logrustash.NewHook("tcp", "stat.suretly.io:5044", "test")
	if err != nil {
		log.Fatal(err)
	}
	log.Hooks.Add(hook)
	ctx := log.WithFields(logrus.Fields{
		"method": "main",
	})

	for i:=0;i<100000;i++ {
		ctx.Info("Hello World!")
	}
}