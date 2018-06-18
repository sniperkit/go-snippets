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

package outline

import (
	"net/url"

	"golang.org/x/text/language"
	"xojoc.pw/crawl/html"
)

// TODO: markdown?
/*
func toText(n *html.Node, buf *bytes.Buffer) {
	for {
		if n == nil {
			break
		}
		if n.Type == html.TextNode {
			buf.WriteString(strings.TrimSpace(n.Data))
		}
		if n.Type == html.ElementNode {
			toText(n.FirstChild, buf)
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				buf.WriteString("\n\n")
			case "p":
				buf.WriteString("\n\n")

			}
		}

		n = n.NextSibling
	}
}
*/

func ToText(n *html.Node) string {
	/*
		if n == nil {
			return ""
		}
		var buf bytes.Buffer
		toText(n.FirstChild, &buf)
		return buf.String()
	*/
	return n.PlainText()
}

type Anchor struct {
	Title string
	// <a>label</a> or <a><img alt="label"/></a>
	Label string
	URL   *url.URL
	Rel   string
}

// Document Kind
type DocumentType int

const (
	UnknownType DocumentType = iota
	ArticleType
	BlogIndexType
	FAQType
)

type Document struct {
	Type DocumentType

	Language    language.Tag
	Title       string
	Description string
	Author      string

	Nav []Anchor
}

func (o *Outline) Extract() *Document {
	if o.DocumentNode == nil ||
		o.HeadNode == nil ||
		o.BodyNode == nil ||
		o.MainNode == nil {
		return nil
	}

	i := &Document{}
	extractHead(o.HeadNode, i)
	for _, n := range o.NavNodes {
		extractNav(n, i)
	}
	return i
}

func extractHead(n *html.Node, i *Document) {
	n = n.FirstChild()
	for {
		if n == nil {
			return
		}
		//		if n.Type == html.ElementNode {
		//			switch n.Data {
		switch {
		//			case "link":
		case n.IsElement("link"):
			//			case "meta":
		case n.IsElement("meta"):
			//			if attribute(n, "content") != "" {
			if n.Attr("content") != "" {
				name := n.Attr("name", "http-equiv")
				switch name {
				case "description":
					i.Description = n.Attr("content")
				case "author":
					i.Author = n.Attr("content")
				}
			}
			//			case "title":
		case n.IsElement("title"):
			if n.FirstChild() != nil {
				i.Title = n.FirstChild().PlainText()
			}
		}
		//		}
		n = n.NextSibling()
	}
}

func extractNav(n *html.Node, i *Document) {
	n = n.FirstChild()
	var stack []*html.Node
	for {
		if n == nil {
			if len(stack) == 0 {
				return
			}
			n = stack[0]
			stack = stack[1:]
		}

		//		if n.Type == html.ElementNode && n.Data == "a" {
		if n.IsElement("a") {
			// TODO: img alt label
			var err error
			a := Anchor{}
			if n.Attr("href") != "" {
				a.URL, err = url.Parse(n.Attr("href"))
				if err != nil {
					continue
				}
			}
			a.Title = n.Attr("title")
			a.Rel = n.Attr("rel")
			/*
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					a.Label = n.FirstChild.Data
				}
			*/
			a.Label = n.PlainText()
			i.Nav = append(i.Nav, a)
		} else {
			if n.FirstChild() != nil {
				stack = append(stack, n.FirstChild())
			}
		}

		n = n.NextSibling()
	}
}

func extractMain(n *html.Node, i *Document) {
	n = n.FirstChild()
	var stack []*html.Node
	for {
		if n == nil {
			if len(stack) == 0 {
				return
			}
			n = stack[0]
			stack = stack[1:]
		}
	}
}
