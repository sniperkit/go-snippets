// This is a quick cli for copying directories to verify that
// fileutil.CopyDirectory works

package main

import (
	"fmt"
	"os"

	"github.com/uoregon-libraries/gopkg/fileutil"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %q <source directory> <destination directory>\n\n", os.Args[0])
		os.Exit(1)
	}

	var err = fileutil.CopyDirectory(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Fail: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}
