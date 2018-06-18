package main

import (
	"fmt"
	"time"

	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	j := &CronJob{PrintText: "Hello world"}

	// execute cron job manually (ie. testing purpose).
	// A `Job` is an interface with single function `Run()`
	// j.Run()

	// or add job to robfig/cron and execute:

	// 1: add cron job
	c := cron.New()
	c.AddJob("* * * * * *", j)
	c.Start()

	// 2: simulate a running system with channel
	// for robfig/cron to operate.
	// use `ctrl+c` to stop.
	done := make(chan bool)
	<-done
}

type CronJob struct {
	PrintText string
}

func (j *CronJob) Run() {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("=== CRON JOB STARTED: %s ===\n", now)
	fmt.Println(j.PrintText)
	fmt.Println("=== CRON JOB DONE ===")
}
