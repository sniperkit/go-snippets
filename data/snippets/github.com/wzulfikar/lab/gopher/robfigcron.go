package main

import (
	"fmt"
	"os"
	"os/signal"

	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	for {
		c := cron.New()
		c.AddFunc("* * * * * *", func() {
			fmt.Println("Every second")
		})

		c.AddFunc("@daily", func() {
			fmt.Println("Every day")
		})

		go c.Start()
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
	}
	// Inspect the cron job entries' next and previous run times.
	// spew.Dump(c.Entries())
}
