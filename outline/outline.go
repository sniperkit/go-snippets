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

// Package outline extracts information from html pages.
package outline

// TODO: simple dump and nav tool for the web
// TODO: detect kind of document (article, index, blog index, home page, tool)

import (
	"errors"
	"io"

	"xojoc.pw/crawl/html"
)

// feed, css, javascript, etc.
/*
type Resource struct {
	URL  *url.URL
	Type string
	Rel  string
}
*/

type Outline struct {
	DocumentNode *html.Node
	HeadNode     *html.Node
	BodyNode     *html.Node
	NavNodes     []*html.Node
	MainNode     *html.Node
	ArticleNodes []*html.Node
	SidebarNode  *html.Node
	FooterNode   *html.Node
}

// TODO: check for ul

func isNav(n *html.Node) bool {
	if n.IsElement("nav") {
		return true
	}
	if n.IsElement("div") {
		switch n.Attr("id") {
		case "navigation", "nav", "menu":
			return true
		}
	}
	return false
}

func isMain(n *html.Node) bool {
	if n.IsElement("main", "article") {
		return true
	}
	if n.IsElement("div") {
		switch n.Attr("id") {
		case "main", "content":
			return true
		}
		/*
			switch attribute(n, "class") {
			case "section-content":
				return true
			}
		*/
	}
	if n.Attr("role") == "main" {
		return true
	}
	return false
}

func isSidebar(n *html.Node) bool {
	if n.IsElement("div") {
		switch n.Attr("id") {
		case "sidebar":
			return true
		}
	}
	return false
}

func isFooter(n *html.Node) bool {
	if n.IsElement("footer") {
		return true
	}
	if n.IsElement("div") {
		switch n.Attr("id") {
		case "footer":
			return true
		}
	}
	return false
}

func buildBody(n *html.Node, o *Outline) {
	/*
		if attribute(n, "lang") != "" {
			o.Language = attribute(n, "lang")
		}
	*/
	var stack []*html.Node
	n = n.FirstChild()
	for {
		if n == nil {
			if len(stack) == 0 {
				return
			}
			n = stack[0]
			stack = stack[1:]
		}

		switch {
		case isNav(n):
			o.NavNodes = append(o.NavNodes, n)
		case isMain(n):
			o.MainNode = n
		case isSidebar(n):
			o.SidebarNode = n
		case isFooter(n):
			o.FooterNode = n
		default:
			if n.FirstChild() != nil {
				stack = append(stack, n.FirstChild())
			}
		}

		n = n.NextSibling()
	}
}
func isArticle(n *html.Node) bool {
	if n.IsElement("article") {
		return true
	}
	if n.IsElement("div") {
		switch n.Attr("class") {
		case "article", "post":
			return true
		}
	}
	return false
}
func buildMain(n *html.Node, o *Outline) {
	n = n.FirstChild()
	var stack []*html.Node
	for {
		if n == nil {
			if len(stack) == 0 {
				break
			}
			n = stack[0]
			stack = stack[1:]
		}

		if isArticle(n) {
			o.ArticleNodes = append(o.ArticleNodes, n)
		} else {
			if n.FirstChild() != nil {
				stack = append(stack, n.FirstChild())
			}
		}

		n = n.NextSibling()
	}
	/*
		if len(o.ArticleNodes) == 0 {
			o.ArticleNodes = append(o.ArticleNodes, n)
		}
	*/
}

var Error = errors.New("malformed structure")

func dump(n *html.Node) {
	return
	/*
		if n == nil {
			return
		}
		cpy := *n
		cpy.Parent = nil
		cpy.FirstChild = nil
		cpy.LastChild = nil
		cpy.PrevSibling = nil
		cpy.NextSibling = nil

		spew.Println("type: ", cpy.Type, ", data: ", cpy.Data, ", attr: ", cpy.Attr)
	*/
}

func Build(r io.Reader) (*Outline, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("html.Parse: doc is nil")
	}
	o := &Outline{}
	o.DocumentNode = doc

	n := doc

	var stack []*html.Node
	for {
		if n == nil {
			if len(stack) == 0 {
				break
			}
			n = stack[0]
			stack = stack[1:]
		}

		switch {
		case n.IsElement("head"):
			if o.HeadNode != nil {
				return nil, errors.New("duplicate head")
			}
			o.HeadNode = n
		case n.IsElement("body"):
			if o.BodyNode != nil {
				return nil, errors.New("duplicate body")
			}
			o.BodyNode = n
		default:
			if n.FirstChild() != nil {
				stack = append(stack, n.FirstChild())
			}

		}

		if o.HeadNode != nil && o.BodyNode != nil {
			break
		}

		n = n.NextSibling()
	}

	if o.HeadNode == nil || o.BodyNode == nil {
		return nil, errors.New("missing head or body")
	}

	buildBody(o.BodyNode, o)

	if o.MainNode != nil {
		buildMain(o.MainNode, o)
	}

	return o, nil
}
