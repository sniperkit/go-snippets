package main

import "syscall"
import "os"
import "os/exec"

func main() {
	// For our example we’ll exec ls. Go requires an absolute path to the binary we want to execute, so we’ll use exec.LookPath to find it (probably /bin/ls).
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}
	// Exec requires arguments in slice form (as apposed to one big string). We’ll give ls a few common arguments. Note that the first argument should be the program name.
	args := []string{"ls", "-a", "-l", "-h"}
	// Exec also needs a set of environment variables to use. Here we just provide our current environment.
	env := os.Environ()
	// Here’s the actual syscall.Exec call. If this call is successful, the execution of our process will end here and be replaced by the /bin/ls -a -l -h process. If there is an error we’ll get a return value.
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
