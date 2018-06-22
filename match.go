package main

import (
	"strings"

	"github.com/bmatcuk/doublestar"
)

func Match(pattern string, filenames []string) []string {
	var matched []string

	globbed, _ := doublestar.Glob(pattern)
	for _, filename := range filenames {
		if contains(filename, globbed) {
			matched = append(matched, filename)
		}
	}

	return matched
}

func contains(filename string, filenames []string) bool {
	for _, f := range filenames {
		if strings.EqualFold(filename, f) {
			return true
		}
	}
	return false
}
