package main

import (
	"bufio"
	"github.com/koron/beni"
	"os"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	beni.Highlight(stdin, stdout, "go", "base16", "Terminal256")
}
