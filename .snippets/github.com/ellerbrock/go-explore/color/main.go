package main

import "github.com/fatih/color"

func main() {
	okMsg("all good, keep going ...")
	errMsg("damn something went wrong ...", 1)
}

func okMsg(msg string) {
	color.Green("[ ✓ ] " + msg)
}

func errMsg(msg string, errCode int) int {
	color.Red("[ ✘ ] " + msg)
	return errCode
}
