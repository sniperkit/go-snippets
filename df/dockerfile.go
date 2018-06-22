// Package df is a wrapper around
// github.com/docker/docker/builder/parser to make it easier to parse
// Dockerfile syntax.
package df

import (
	"io"
	"strings"

	dfparser "github.com/docker/docker/builder/parser"
)

// Dockerfile represents a Dockerfile.
//
// TODO(sqs): The fields present are all we care about right now. This
// is obviously incomplete.
type Dockerfile struct {
	From string
}

// Decode reads a Dockerfile and returns a struct representation of
// it. An error is returned if parsing fails.
func Decode(r io.Reader) (*Dockerfile, error) {
	ast, err := dfparser.Parse(r)
	if err != nil {
		return nil, err
	}

	var df Dockerfile
	for _, n := range ast.Children {
		switch strings.ToUpper(n.Value) {
		case "FROM":
			df.From = strings.TrimSpace(n.Original[len(n.Value):])
		}
	}

	return &df, nil
}
