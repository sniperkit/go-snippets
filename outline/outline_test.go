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

package outline_test

import (
	"fmt"
	"reflect"
	"testing"

	"xojoc.pw/crawl/html"
	"xojoc.pw/must"

	"xojoc.pw/crawl/outline"
)

var htmls = []struct {
	file string
	out  outline.Outline
}{
	{"basic.html", outline.Outline{
		Doctype:     "",
		Language:    "it",
		Title:       "Title",
		Description: "Description",
		Author:      "Author",

		Menu: []outline.Anchor{
			{
				URL:   must.URL("/a"),
				Label: "a",
			},
		},

		HeadNode: &html.Node{
			Data: "head",
		},
		BodyNode: &html.Node{
			Data: "body",
		},
		MenuNode: &html.Node{
			Data: "nav",
			Attr: []html.Attribute{{Key: "id", Val: "navigation"}},
		},
	}},
	{"free-sw.html", outline.Outline{
		Doctype: `html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
	       "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"`,
		Language: "en",
		Title: `What is free software?
- GNU Project - Free Software Foundation`,
		Description: `Since 1983, developing the free Unix style operating system GNU, so that computer users can have the freedom to share and improve the software they use.`,

		Menu: []outline.Anchor{
			{URL: must.URL("/gnu/gnu.html")},
		},

		HeadNode: &html.Node{
			Data: "head",
		},
		BodyNode: &html.Node{
			Data: "body",
		},
		MenuNode: &html.Node{
			Data: "div",
			Attr: []html.Attribute{{Key: "id", Val: "navigation"}},
		},
	}},
}

func cmp(t *testing.T, want, got string) {
	if want != got {
		t.Errorf("Want %q, got %q", want, got)
	}
}

func TestBuild(t *testing.T) {
	for _, h := range htmls {
		r := must.Open("testfiles/articles/" + h.file)
		defer r.Close()
		got, err := outline.Build(r)
		must.OK(err)
		// got.Doctype != h.out.Doctype ||
		cmp(t, h.out.Language, got.Language)
		cmp(t, h.out.Title, got.Title)
		cmp(t, h.out.Description, got.Description)
		cmp(t, h.out.Author, got.Author)

		if !reflect.DeepEqual(h.out.Menu, got.Menu) {
			t.Errorf(`
		# got
		%#v
		# want
		%#v
		`, got.Menu, h.out.Menu)
		}

		cmp(t, h.out.HeadNode.Data, got.HeadNode.Data)
		cmp(t, h.out.BodyNode.Data, got.BodyNode.Data)
		cmp(t, h.out.MenuNode.Data, got.MenuNode.Data)
		fmt.Printf("%+v\n", got.MenuNode)
	}
}
func TestOutline_Extract(t *testing.T) {

}
