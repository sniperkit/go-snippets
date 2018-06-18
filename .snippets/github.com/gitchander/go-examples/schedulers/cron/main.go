package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 39-42 * * * *", func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})
	c.Start()
	defer c.Stop()

	time.Sleep(8 * time.Minute)
}
