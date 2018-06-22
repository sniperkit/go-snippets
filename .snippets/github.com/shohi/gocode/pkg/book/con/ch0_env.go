package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func printEnv() {
	env := os.Environ()
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	fmt.Fprintf(tw, "Key\tValue\n")
	fmt.Println(os.Hostname())
	for _, raw := range env {
		pair := strings.Split(raw, "=")
		key := pair[0]
		value := pair[1]
		fmt.Fprintf(tw, "%s\t%s\n", key, value)
	}
}

type ST struct {
	Value int
}

func emptyStructCompare() {
	var a, b struct{}
	var c, d ST
	var t time.Time
	fmt.Printf("a ==> %p, %v\n", &a, a)
	fmt.Printf("b ==> %p, %v\n", &b, b)

	fmt.Printf("c ==> %p, %v\n", &c, d)
	fmt.Printf("d ==> %p, %v\n", &d, d)
	fmt.Printf("time ==> %v", t)

}

func main() {
	// printEnv()
	// fmt.Println(runtime.Caller(1))
	emptyStructCompare()
}
