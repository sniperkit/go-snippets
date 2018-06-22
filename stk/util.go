package main

import (
	"os"
	"regexp"
	"strings"
)

var nohtml *regexp.Regexp = regexp.MustCompile("<[^>]*>")
var wrongQuote *regexp.Regexp = regexp.MustCompile("&#39;")

func stripHTML(htmlContent string) string {
	return nohtml.ReplaceAllString(htmlContent, "")
}

func fixQuotes(content string) string {
	return wrongQuote.ReplaceAllString(content, "'")
}

func xterm(code string) func(s string) string {
	env := os.Getenv("TERM")
	isXterm := strings.Contains(env, "xterm")

	return func(text string) (output string) {
		if isXterm {
			output = code + text + "\033[0m"
		} else {
			output = text
		}
		return
	}
}

func bold(text string) string {
	return xterm("\033[1m")(text)
}

func underline(text string) string {
	return xterm("\033[4m")(text)
}
