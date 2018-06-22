package tools

import "gopkg.in/bblfsh/sdk.v1/uast"

type Tokenizer struct{}

func (t Tokenizer) Exec(node *uast.Node) error {
	for _, token := range Tokens(node) {
		print(token)
	}
	return nil
}

// Tokens returns a slice of tokens contained in the node.
func Tokens(n *uast.Node) []string {
	var tokens []string
	iter := uast.NewOrderPathIter(uast.NewPath(n))
	for {
		p := iter.Next()
		if p.IsEmpty() {
			break
		}

		n := p.Node()
		if n.Token != "" {
			tokens = append(tokens, n.Token)
		}
	}
	return tokens
}
