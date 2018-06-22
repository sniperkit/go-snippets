package lib

import (
	"os/exec"
)

// GoFmt ...
// @desc Code formatting
func GoFmt(filename string) {
	cmd := exec.Command("go", "fmt", filename)
	cmd.Start()
	cmd.Wait()
}
