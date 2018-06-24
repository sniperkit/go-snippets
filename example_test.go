package cli_test

import "github.com/gowww/cli"

var (
	address    string
	production bool
)

func Example() {
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
