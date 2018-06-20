package main

import (
	"fmt"

	"github.com/ctava/tfcgo"
)

func main() {

	fmt.Println("SaySomething = ", tfcgo.SaySomething("Space. The final frontier."))

}
