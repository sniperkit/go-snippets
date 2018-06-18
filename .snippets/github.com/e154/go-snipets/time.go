package main

import (
	"fmt"
	"time"
	"strings"
)

func main() {
	p := fmt.Println
	p("Welcome to the playground!")

	t := time.Now()
	fmt.Printf("The time is %s\n", t)

	var h time.Duration = 6
	t = t.Add(h*time.Hour)

	fmt.Printf("The time plus %d Hours -> %s\n", h, t.Local())

	// сравнение времени
	t2 := t.Add(time.Second * 60)
	p("==", t2.Sub(t).Seconds())

	if strings.Contains(t.Sub(t2).String(), "-") {
		p("-")
	} else {
		p("+")
	}

	// время из строки
	str := "2015-04-08 21:27:52 +0700 +0700"
	t3, _ := time.Parse(time.RFC1123Z, str)
	p(t3)

	fmt.Println(fmt.Sprintf("%d:%d %d/%d/%d", t.Hour(), t.Minute(), t.Day(), t.Month(), t.Year()))

	// Ticker
	c := time.Tick(1 * time.Second)
	for now := range c {
		fmt.Printf("%v %s\n", now, "qwe")
		break
	}

	t4 := time.Date(2020, 03, 1, 0, 0, 0, 0, time.UTC)
	p(t4)

	fmt.Println(len("3"))
	fmt.Println(len("12"))
	t5, _ := time.Parse("2006-01-02", "2020-04-01")
	p(t5)
}
