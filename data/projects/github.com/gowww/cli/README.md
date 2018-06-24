# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) cli [![GoDoc](https://godoc.org/github.com/gowww/cli?status.svg)](https://godoc.org/github.com/gowww/cli) [![Build](https://travis-ci.org/gowww/cli.svg?branch=master)](https://travis-ci.org/gowww/cli) [![Coverage](https://coveralls.io/repos/github/gowww/cli/badge.svg?branch=master)](https://coveralls.io/github/gowww/cli?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/cli)](https://goreportcard.com/report/github.com/gowww/cli) ![Status Unstable](https://img.shields.io/badge/status-unstable-red.svg)

Package [cli](https://godoc.org/github.com/gowww/cli) wraps the standard [flag](https://golang.org/pkg/flag/) package for a cleaner command line interface.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/cli
	```

2. Import it in your code:

	```Go
	import "github.com/gowww/cli"
	```

## Usage

Work in progress...

### Example

```Go
package main

import "github.com/gowww/cli"

var (
	address    string
	production bool
)

func main() {
	cli.Description = "Command line interface example."

	cli.Command("run", run, "Run app.").
		Bool(&production, "docker", false, `Run the server in production environment.`)

	cli.Command("watch", watch, "Detect changes and rerun app.")

	cli.String(&address, "address", ":8080", "The address to listen and serve on.")

	cli.Parse()
}

func run() {
	// Run app.
}

func watch() {
	// Detect changes and rerun app.
}
```

#### Usage output

```Shell
Command line interface example.

Usage:

        example [command]

Commands:

        run    Run app.
        watch  Detect changes and rerun app.

Flags:

        -address=:8080  The address to listen and serve on.

```
