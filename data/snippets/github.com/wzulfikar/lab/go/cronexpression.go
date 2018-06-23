package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	// gorhill()
	robfig()
}

func robfig() {
	// doc: gopkg.in/robfig/cron.v2/doc.go

	// robfig/cron's cron spec (cron expression)
	cronExpr := "0 0 8 * * *" // everyday at 8am
	// cronExpr := "0 0 */3 * * *" // everyday in 3hrs windows
	// cronExpr := "0 0 8 * * 1" // 8am every monday

	c := cron.New()
	c.AddFunc(cronExpr, func() { log.Println("Cron job executed") })
	c.Start()

	now := time.Now()

	fmt.Printf("Cron entry for `%s`:\n", cronExpr)
	entry := c.Entries()[0]
	fmt.Printf("- Prev time       : %s\n", entry.Prev)
	fmt.Printf("- Current time    : %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("- Next time       : %s\n", entry.Next)
	fmt.Printf("- Current to next : %s", diffTime(entry.Next, now))
	fmt.Println()

	// done := make(chan bool)
	// <-done
}
func gorhill() {
	// cron docs: https://godoc.org/github.com/robfig/cron
	cronExpr := "0 8 * * * *"
	nextTimes := cronexpr.MustParse(cronExpr).NextN(time.Now(), 5)

	fmt.Printf("Next 5 times of cron expression `%s`:\n", cronExpr)
	for i, next := range nextTimes {
		fmt.Printf("%d. %s\n", i+1, next)
	}
}

func diffTime(t1, t2 time.Time) string {
	diff := t1.Sub(t2)

	if diff.Hours() < 24 {
		return diff.String()
	}

	days := int(diff.Hours() / 24)
	hrs := int(diff.Hours()) % (24 * days)
	mins := int(diff.Minutes()) - (days * 24 * 60) - (hrs * 60)
	return fmt.Sprintf("%dd %dh %dm", days, hrs, mins)
}
