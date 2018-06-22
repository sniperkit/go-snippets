package main

import "github.com/spf13/cobra"

// cmdRoot is the root of our command tree
var cmdRoot = &cobra.Command{
	Use:   "quickstart",
	Short: "Go quickstart sample application",
}

func main() {
	// cmdRoot.Execute acts as our entrypoint: call it
	cmdRoot.Execute()
}
