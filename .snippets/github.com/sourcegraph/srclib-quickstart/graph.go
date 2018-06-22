// graph.go is responsible for finding references and definitions in a source unit.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"sourcegraph.com/sourcegraph/srclib/graph"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

func init() {
	_, err := flagParser.AddCommand("graph",
		// TODO: Update this description with your own.
		"graph source units",
		"Reads source unit descriptors from STDIN and writes references and definitions to STDOUT.",
		&graphCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type GraphCmd struct{}

var graphCmd GraphCmd

func graphFile(path string, output *graph.Output) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %s", path, err)
	}
	defer f.Close()

	// TODO: Add references to `output.Refs` and definitions to `output.Defs`.

	return nil
}

func graphUnits(units unit.SourceUnits) (*graph.Output, error) {
	output := graph.Output{}

	for _, u := range units {
		for _, f := range u.Files {
			graphFile(f, &output)
		}
	}

	return &output, nil
}

// Execute is called by the command line interface when invoked with the `graph` subcommand.
func (c *GraphCmd) Execute(args []string) error {
	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("Failed to read STDIN: %s", err)
	}
	var units unit.SourceUnits
	if err := json.NewDecoder(bytes.NewReader(inputBytes)).Decode(&units); err != nil {
		// Legacy API: Try parsing input as a single source unit.
		var u *unit.SourceUnit
		if err := json.NewDecoder(bytes.NewReader(inputBytes)).Decode(&u); err != nil {
			return fmt.Errorf("Failed to parse source units from input: %s", err)
		}
		units = unit.SourceUnits{u}
	}
	if err := os.Stdin.Close(); err != nil {
		return fmt.Errorf("Failed to close STDIN: %s", err)
	}

	if len(units) == 0 {
		log.Fatal("Input contains no source unit data.")
	}

	out, err := graphUnits(units)
	if err != nil {
		return fmt.Errorf("Failed to graph source units: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		return fmt.Errorf("Failed to output graph data: %s", err)
	}
	return nil
}
