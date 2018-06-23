package toolkit

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func Html2text(source string) (string, error) {
	// ========= Parse the HTML

	doc, err := html.Parse(strings.NewReader(source))
	if err != nil {
		return "", err
	}

	res := ""
	var f func(n *html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case "li", "br", "p", "div", "hr":
				res = fmt.Sprintf("%s\n", res)
			}
		case html.TextNode:
			res = fmt.Sprintf("%s%s ", res, n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// ========= Replace double spaces

	resFinal := ""
	for _, each := range strings.Split(res, "\n") {
		each := strings.TrimSpace(each)

		if each == "" {
			continue
		}

		resFinal = fmt.Sprintf("%s\n%s", resFinal, each)
	}

	return strings.TrimSpace(resFinal), nil
}
