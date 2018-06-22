package gopherWall

import (
	"fmt"
	"os/exec"
)

// list the firewall rules
func List() {
	output, err := exec.Command("ipfw", "list").Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(output[:]))
}
