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

package sitemap

import (
	"bufio"
	"bytes"
	"io"

	"xojoc.pw/must"
)

type Location struct {
	URL string
}

type Sitemap struct {
	Locations []*Location
	Sitemaps  []string
}

func Parse(r io.Reader) (*Sitemap, error) {
	s := Sitemap{}
	buf := bufio.NewReader(r)
	for {
		l, err := buf.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte(`<loc>`)) {
			s.Locations = append(s.Locations,
				&Location{
					URL: string(
						l[len(`<loc>`) : len(l)-len(`</loc>`)],
					),
				})
		}
		if bytes.HasPrefix(l, []byte(`<sitemap><loc>`)) {

			s.Sitemaps = append(s.Sitemaps,
				string(l[len(`<sitemap><loc>`):len(l)-len(`</loc></sitemap>`)]))
		}

	}
	return &s, nil
}

func MustParse(r io.Reader) *Sitemap {
	s, err := Parse(r)
	must.OK(err)
	return s
}
