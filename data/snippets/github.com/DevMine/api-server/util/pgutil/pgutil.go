// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pgutil provides utilities to deal with PostgreSQL specific data
// types.
package pgutil

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

// construct a regexp to extract values:
var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	// escaped double quotes
	escapedDoubleQuotesExp = regexp.MustCompile(`\\"`)

	valueIndex int
)

// ParsePgArray parses the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
// Thanks to adharris: https://gist.github.com/adharris/4163702
func ParsePgArray(array string) []string {
	var results []string
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]

		// trim remaining comma, if any
		s = strings.TrimSuffix(s, ",")

		// the string _might_ be wrapped in quotes, so trim them:
		s, err := strconv.Unquote(s)
		if err != nil {
			glog.Error(err)
		}

		results = append(results, s)
	}
	return results
}

// ParsePgRow a pg row which looks like "(12,\"foo, bar\", 12)" (unquoted)
func ParsePgRow(row string) []string {
	// trim parenthesis
	row = strings.Trim(row, "()")

	// replace escaped quotes, if any
	row = escapedDoubleQuotesExp.ReplaceAllLiteralString(row, `"`)

	sr := strings.NewReader(row)
	cr := csv.NewReader(sr)

	results, err := cr.Read()
	if err != nil {
		glog.Error(err)
	}

	return results
}

// TimestampTZToTime converts a pg timestamp with timezone datetime
// formatted string to a time.Time
func TimestampTZToTime(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.Replace(t, " ", "T", -1)+":00")
}
