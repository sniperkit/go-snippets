package main

import (
	"github.com/gliderlabs/com/daemon"

	_ "github.com/gl-prototypes/jq-cmd/app/jqserve"
	_ "github.com/gl-prototypes/jq-cmd/lib/ssh/init"
)

func main() {
	daemon.Run("jqserve")
}
