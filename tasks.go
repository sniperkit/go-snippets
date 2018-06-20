package main

import (
	"time"
	"github.com/RichardKnop/machinery/v1/log"
	"fmt"
)

// LongRunningTask ...
func LongRunningTask() error {
	log.INFO.Print("Long running task started")
	for i := 0; i < 10; i++ {
		log.INFO.Print(10 - i)
		<-time.After(1 * time.Second)
	}
	log.INFO.Print("Long running task finished")
	return nil
}

// Add ...
func SendSingleSms(args ...string) (string, error) {
	log.INFO.Print("Send single sms task started")
	mobile := args[0]
	text := args[1]
	log.INFO.Print("Send single sms task started")
	return fmt.Sprintf("The mobile is %v and the text is %v", mobile, text), nil
}

// Multiply ...
func SendMultipleSms(args ...string) (string, error) {
	log.INFO.Print("Send multiple sms task started")
	mobile := args[0]
	text := args[1]
	log.INFO.Print("Send multiple sms task finished")
	return fmt.Sprintf("The mobile is %v and the text is %v", mobile, text), nil
}

func DelayedTask(args ...int64) error {
	log.INFO.Print("It's a delayed task")
	return nil
}
//
//func OnError(args ...string) error {
//	log.INFO.Print(args[0])
//	return nil
//}