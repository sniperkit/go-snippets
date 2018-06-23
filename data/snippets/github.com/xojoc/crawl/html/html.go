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

package html

import (
	"io"
	"strings"

	rawhtml "golang.org/x/net/html"
)

// todo: tree traversal like ftw

type Node struct {
	node *rawhtml.Node
}

func Parse(r io.Reader) (*Node, error) {
	n, err := rawhtml.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Node{node: n}, nil
}

func (n *Node) Render(w io.Writer) error {
	return rawhtml.Render(w, n.node)
}

func (n *Node) Attr(keys ...string) string {
	if n == nil {
		return ""
	}
	if n.node == nil {
		return ""
	}
	for _, attr := range n.node.Attr {
		for _, key := range keys {
			if attr.Key == key {
				return attr.Val
			}
		}
	}
	return ""
}

func hasClass(n *Node, class string) bool {
	if n == nil {
		return false
	}
	if n.node == nil {
		return false
	}
	if n.node.Type != rawhtml.ElementNode {
		return false
	}
	for _, attr := range n.node.Attr {
		if attr.Key == "class" {
			fs := strings.Fields(attr.Val)
			for _, f := range fs {
				if class == f {
					return true
				}
			}
		}
	}
	return false
}

func (n *Node) Classes(class string) []*Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	var classes []*Node

	ok := true
	cs := strings.Fields(class)
	for _, c := range cs {
		if !hasClass(n, c) {
			ok = false
			break
		}
	}
	if ok {
		classes = append(classes, n)
	}

	for p := n.node.FirstChild; p != nil; p = p.NextSibling {
		classes = append(classes, (&Node{node: p}).Classes(class)...)
	}
	return classes
}

func (n *Node) PlainText() string {
	if n == nil {
		return ""
	}
	if n.node == nil {
		return ""
	}

	if n.node.Type == rawhtml.TextNode {
		return n.node.Data
	}

	t := ""

	for p := n.node.FirstChild; p != nil; p = p.NextSibling {
		t += (&Node{node: p}).PlainText()
	}

	return t
}

func (n *Node) NextElement(name ...string) *Node {
	for p := n.node.FirstChild; p != nil; p = p.NextSibling {
		if p.Type == rawhtml.ElementNode {
			for _, na := range name {
				if p.Data == na {
					return &Node{node: p}
				}
			}
		}
	}
	return nil
}

func (n *Node) Elements(elms ...string) []*Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	var elements []*Node
	if n.node.Type == rawhtml.ElementNode {
		for _, elm := range elms {
			if n.node.Data == elm {
				elements = append(elements, n)
				break
			}
		}
	}

	for p := n.FirstChild(); p != nil; p = p.NextSibling() {
		elements = append(elements, p.Elements(elms...)...)
	}

	return elements
}

func (n *Node) IsElement(names ...string) bool {
	if n == nil {
		return false
	}
	if n.node == nil {
		return false
	}
	if n.node.Type != rawhtml.ElementNode {
		return false
	}
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if n.node.Data == name {
			return true
		}
	}
	return false
}
func (n *Node) NextSibling() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	if n.node.NextSibling == nil {
		return nil
	}
	return &Node{node: n.node.NextSibling}
}
func (n *Node) NextSiblingElement() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	p := n.node.NextSibling
	for {
		if p == nil {
			return nil
		}
		if p.Type == rawhtml.ElementNode {
			return &Node{node: p}
		}
		p = p.NextSibling
	}
}

func (n *Node) FirstChild() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	if n.node.FirstChild == nil {
		return nil
	}
	return &Node{node: n.node.FirstChild}
}
func (n *Node) FirstChildElement() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	p := n.node.FirstChild
	for {
		if p == nil {
			return nil
		}
		if p.Type == rawhtml.ElementNode {
			return &Node{node: p}
		}
		p = p.NextSibling
	}
}

func (n *Node) ID(id string) *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	if n.Attr("id") == id {
		return n
	}
	for p := n.FirstChild(); p != nil; p = p.NextSibling() {
		t := p.ID(id)
		if t != nil {
			return t
		}
	}
	return nil
}

func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	if n.node.Parent == nil {
		return nil
	}
	return &Node{node: n.node.Parent}
}
func (n *Node) ParentElement() *Node {
	if n == nil {
		return nil
	}
	if n.node == nil {
		return nil
	}
	p := n.node.Parent
	for {
		if p == nil {
			return nil
		}
		if p.Type == rawhtml.ElementNode {
			return &Node{node: p}
		}
		p = p.Parent
	}
}

func (n *Node) FakeParent() *Node {
	return &Node{
		node: &rawhtml.Node{
			FirstChild: n.node,
			LastChild:  n.node,
		}}
}

func (n *Node) Elements2(f func(*Node), elms ...string) {
	if n == nil {
		return
	}
	if n.node == nil {
		return
	}
	if n.node.Type == rawhtml.ElementNode {
		for _, elm := range elms {
			if n.node.Data == elm {
				f(n)
				break
			}
		}
	}
	for p := n.FirstChild(); p != nil; p = p.NextSibling() {
		p.Elements2(f, elms...)
	}
}

/*
func printNode(n *html.Node, w io.Writer) error {
	if n == nil {
		return nil
	}
	switch n.Type {
	case html.ErrorNode:
		return nil
	case html.TextNode:
		_, err := w.Write([]byte(n.Data))
		if err != nil {
			return err
		}
	case html.ElementNode:
		_, err := w.Write([]byte(`<`))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(n.Data + " "))
		if err != nil {
			return err
		}

		for _, attr := range n.Attr {
			_, err = w.Write([]byte(fmt.Sprintf("%s=%q ", attr.Key, attr.Val)))
			if err != nil {
				return err
			}
		}

		_, err = w.Write([]byte(`>`))
		if err != nil {
			return err
		}
	}
	for p := n.FirstChild; p != nil; p = p.NextSibling {
		err := printNode(p, w)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
