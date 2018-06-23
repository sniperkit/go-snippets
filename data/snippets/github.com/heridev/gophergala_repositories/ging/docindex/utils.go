package docindex

import (
	"bytes"

	"golang.org/x/net/html"
)

func removeDocSourcecode(text string) string {
	buf := bytes.NewBufferString(text)
	root, err := html.Parse(buf)
	if err != nil {
		return text
	}
	var visit func(n *html.Node)
	visit = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "pre" {
			if n.Parent != nil {
				next := n.NextSibling
				n.Parent.RemoveChild(n)
				visit(next)
			}
		}
		if n.FirstChild != nil {
			visit(n.FirstChild)
		}
		if n.NextSibling != nil {
			visit(n.NextSibling)
		}
	}
	visit(root)
	rbuf := new(bytes.Buffer)
	html.Render(rbuf, root)
	return rbuf.String()
}
