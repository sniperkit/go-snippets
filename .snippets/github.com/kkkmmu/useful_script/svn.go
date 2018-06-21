package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	c := context.Background()
	tc, f := context.WithTimeout(c, time.Second*5)
	command := exec.CommandContext(tc, "git", "log")
	resp, _ := command.Output()
	fmt.Println(string(resp))
	f()
}
