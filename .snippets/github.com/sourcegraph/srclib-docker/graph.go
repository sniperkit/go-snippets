package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"strings"

	"sourcegraph.com/sourcegraph/srclib-docker/df"
	"sourcegraph.com/sourcegraph/srclib/ann"
	"sourcegraph.com/sourcegraph/srclib/graph"
	"sourcegraph.com/sourcegraph/srclib/grapher"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

// The Def.Path and Ref.DefPath of a Dockerfile. There is 1 Dockerfile
// per source unit, so each Dockerfile is at that source unit's "."
// path.
const DockerfileDefPath = "."

func init() {
	_, err := parser.AddCommand("graph",
		"graph a Dockerfile",
		"Graph a Dockerfile.",
		&graphCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type GraphCmd struct{}

var graphCmd GraphCmd

func (c *GraphCmd) Execute(args []string) error {
	var unit *unit.SourceUnit
	if err := json.NewDecoder(os.Stdin).Decode(&unit); err != nil {
		return err
	}
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	dfpath := unit.Files[0]
	data, err := ioutil.ReadFile(dfpath)
	if err != nil {
		return err
	}
	dfJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var o grapher.Output

	// Add def.
	o.Defs = append(o.Defs, &graph.Def{
		DefKey:   graph.DefKey{Path: DockerfileDefPath},
		TreePath: ".",
		Name:     "Dockerfile",
		File:     dfpath,
		Exported: true,

		DefStart: 0,
		DefEnd:   len(data) - 1,

		Data: dfJSON,
	})

	// Add link annotations.
	dataLCStr := string(bytes.ToLower(data))
	d, err := df.Decode(bytes.NewReader(data))
	if err == nil {
		// Add link to base image.
		if d.From != "" {
			_, _, dockerRepo, _, err := resolveImageRef(d.From)
			if err == nil {
				start, end := findStartEnd(dataLCStr, d.From)
				if start != -1 {
					ann := &ann.Ann{
						File:  dfpath,
						Start: start,
						End:   end,
					}
					if err := ann.SetLinkURL("https://registry.hub.docker.com/_/" + dockerRepo); err != nil {
						return err
					}
					o.Anns = append(o.Anns, ann)
				}
			} else {
				log.Printf("Error parsing base image FROM link %q: %s. Skipping link.", d.From, err)
			}
		}

		// Add links to things that look like repo URIs.
		uriIdxs := repoURIPat.FindAllStringIndex(dataLCStr, -1)
		for _, m := range uriIdxs {
			start, end := m[0], m[1]
			uri := string(data[start:end])
			ann := &ann.Ann{
				File:  dfpath,
				Start: start,
				End:   end,
			}
			if err := ann.SetLinkURL("https://sourcegraph.com/" + uri); err != nil {
				return err
			}
			o.Anns = append(o.Anns, ann)
		}

		// Add links to docs for Dockerfile instructions.
		instrIdxs := instructionPat.FindAllStringIndex(dataLCStr, -1)
		for _, m := range instrIdxs {
			start, end := m[0], m[1]
			instruction := string(data[start:end])
			ann := &ann.Ann{
				File:  dfpath,
				Start: start,
				End:   end,
			}
			if err := ann.SetLinkURL("https://docs.docker.com/reference/builder/#" + strings.ToLower(instruction)); err != nil {
				return err
			}
			o.Anns = append(o.Anns, ann)
		}
	} else {
		log.Printf("Error parsing Dockerfile %q: %s. Skipping links.", dfpath, err)
	}

	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

// findStartEnd fuzzily (case-insensitively) finds the start and end
// character positions of substr in s. If substr is not in s, -1 is returned for both start and end.
func findStartEnd(s, substr string) (start, end int) {
	substr = strings.ToLower(substr)
	i := strings.Index(s, substr)
	if i == -1 {
		return -1, -1
	}
	return i, i + len(substr)
}

var (
	repoURIPat     = regexp.MustCompile(`(?:github\.com|sourcegraph\.com)/[\w.-]+/[\w.-]+`)
	instructionPat = regexp.MustCompile(`(?m:^(from|maintainer|run|cmd|expose|env|add|copy|entrypoint|volume|user|workdir|onbuild))`)
)
