/*
Command gopathexec executes program with arguments, while expanding $GOPATH with multiple workspaces into multiple arguments.

Usage: gopathexec program [args]

It is helpful for writing go generate directives for programs that do not understand import paths, but you need to specify include
paths that are other Go package folders.

For example, consider a go generate directive like this:

	//go:generate protoc -I$GOPATH/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

That will only work if your GOPATH env var happens to contain one workspace. By prepending the above command with gopathexec:

	//go:generate gopathexec protoc -I$GOPATH/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

It will effectively execute:

	protoc -I/workspace1/src/github.com/gogo/protobuf/protobuf -I/workspace2/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

If you have 1 or no GOPATH workspaces, gopathexec executes the given program and arguments without modification.
*/
package main // import "sourcegraph.com/sourcegraph/gopathexec"

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("gopathexec: no arguments provided, nothing to exec")
	}
	program := os.Args[1]
	inArgs := os.Args[2:]

	var outArgs []string

	gopath := os.Getenv("GOPATH")
	workspaces := filepath.SplitList(gopath)
	if len(workspaces) > 1 {
		for _, arg := range inArgs {
			if strings.Contains(arg, gopath) {
				for _, workspace := range workspaces {
					outArgs = append(outArgs, strings.Replace(arg, gopath, workspace, -1))
				}
			} else {
				outArgs = append(outArgs, arg)
			}
		}
	} else {
		outArgs = inArgs
	}

	cmd := exec.Command(program, outArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
