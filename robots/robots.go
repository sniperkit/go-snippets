/*  Copyright (C) 2018 Alexandru Cojocaru

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

// Package robots parses a robots.txt file as specified by Wikipedia https://en.wikipedia.org/wiki/Robots.txt
package robots // import "xojoc.pw/crawl/robots"

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

// Txt contains the robots.txt rules.
type Txt struct {
	CrawlDelay map[string]int
	Allow      map[string][]string
	Disallow   map[string][]string

	Sitemaps []string
}

// Allowed returns true if user agent ua can access path.
// False otherwise.
func (t *Txt) Allowed(ua string, path string) bool {
	allowed := true

	disallowedPath := ""

	for _, d := range append(t.Disallow["*"], t.Disallow[ua]...) {
		if d == "" {
			allowed = true
			disallowedPath = ""
			continue
		}
		if strings.HasPrefix(path, d) {
			if len(d) < len(disallowedPath) || disallowedPath == "" {
				disallowedPath = d
			}
			allowed = false
		}

	}
	for _, a := range append(t.Allow["*"], t.Allow[ua]...) {
		if strings.HasPrefix(path, a) {
			if len(a) > len(disallowedPath) {
				allowed = true
			}
		}
	}

	return allowed
}

// Delay returns the number of seconds to wait between successive accesses to
// the same host.
// Returns 0 if no delay is specified.
func (t *Txt) Delay(ua string) int {
	if i, ok := t.CrawlDelay[ua]; ok {
		return i
	}
	return t.CrawlDelay["*"]
}

// Parse parses a robots.txt file.
func Parse(r io.Reader) (*Txt, error) {
	txt := &Txt{}
	txt.CrawlDelay = map[string]int{}
	txt.Disallow = map[string][]string{}
	txt.Allow = map[string][]string{}

	buf := bufio.NewReader(r)
	var ua string

	lua := []byte("User-agent: ")
	ldis := []byte("Disallow:")
	lall := []byte("Allow: ")
	lsm := []byte("Sitemap: ")
	lcd := []byte("Crawl-delay: ")

	for {
		l, err := buf.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		l = l[:len(l)-1]

		switch {
		case bytes.HasPrefix(l, lua):
			ua = string(l[len(lua):])
		case bytes.HasPrefix(l, ldis):
			dis := ""
			if len(l) > len(ldis) {
				dis = string(l[len(ldis)+1:])
			}
			txt.Disallow[ua] = append(txt.Disallow[ua], dis)
		case bytes.HasPrefix(l, lall):
			txt.Allow[ua] = append(txt.Allow[ua], string(l[len(lall):]))
		case bytes.HasPrefix(l, lcd):
			i, err := strconv.Atoi(string(l[len(lcd):]))
			if err == nil {
				txt.CrawlDelay[ua] = i
			}
		case bytes.HasPrefix(l, lsm):
			txt.Sitemaps = append(txt.Sitemaps, string(l[len(lsm):]))
		default:
			// skip line
		}
	}

	return txt, nil
}
