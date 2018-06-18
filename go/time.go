package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	t2 := now.Add(time.Second * 100000)

	fmt.Println(t2)

	fmt.Println("now", now)

	then := now.Add(-3 * time.Hour)
	fmt.Println("then", then)
	fmt.Println("then mysql", then.Format("2006-01-02 15:04:05"))

	// show time diff in hours
	diff := now.Sub(t2)
	fmt.Println(diff.Hours())

	daysAgo, _ := time.Parse(time.RFC822, "07 Feb 18 10:00 UTC")
	fmt.Println("sub in hours:", int(now.Sub(daysAgo).Hours()/24))
	if int(now.Sub(daysAgo).Hours()/24) == 2 {
		fmt.Printf("%v was two days ago", daysAgo)
	}
}
