package main

import "fmt"
import "bufio"
import "os/exec"

func main() {
	cmd0 := exec.Command("echo", "-n", "My first command from golang")
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}
	outputBuf0 := bufio.NewReader(stdout0)
	output0, _, err := outputBuf0.ReadLine()
	if err != nil {
		fmt.Printf("Error: Can not read data the pipe: %s\n", err)
		return
	}
	fmt.Printf("%s\n", string(output0))

}
