package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// args := []string{"test.cpp", "-o", "test.o"}
	// args = append(args, "--std=c++11")

	cmd := "echo 'Think123!' | kinit dhcp4dns@HQ.TPN.THINKINGPHONES.NET"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: " + cmd)
	}

	fmt.Printf("Command ran: %s\n", string(out))
}
