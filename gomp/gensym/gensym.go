// Package gensym is used to generate unique symbol names.
package gensym

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"strconv"
)

const defaultname = "gompsym"

type TrieNode struct {
	// TODO: rune?
	move  map[uint8]*TrieNode
	final bool
}

func MkGenDummy(s string) func() string {
	var count int = 0
	return func() string {
		buffer := bytes.NewBufferString("__sym")
		fmt.Fprint(buffer, strconv.Itoa(count))
		count++
		return buffer.String()
	}
}

// MkGen returns a function that upon each call will provide a string.
// This string can be safely used as a name for a new variable in the
// file scope of the file provided in source.
// Each call returns a new function. Functions returned
// from different calls using the same source string generate
// identical sequences of names.
func MkGen(src string) func() string {
	known := extractSymbols(src)
	vocabPos := len(vocab)
	if VocabOn {
		vocabPos = 0
	}
	id := 0
	root := &TrieNode{move: make(map[uint8]*TrieNode), final: true}
	for _, s := range known {
		addWord(root, s)
	}
	return func() string {
		for vocabPos < len(vocab) && !addWord(root, vocab[vocabPos]) {
			vocabPos++
		}
		if vocabPos >= len(vocab) {
			// ran out of vocab: invent a new symbol
			return nextName(root, &id)
		} else {
			vocabPos++
			return vocab[vocabPos-1]
		}
	}
}

// nextName returns a string that is not present in the
// trie represented by root (although it may be a prefix of such string)
func nextName(root *TrieNode, id *int) string {
	for {
		name := defaultname + strconv.Itoa(*id)
		(*id)++
		if addWord(root, name) {
			return name
		}
	}
}

// addWord attempts to add the word s to the trie and returns
// true if s has not already been there
func addWord(root *TrieNode, s string) bool {
	nd := root
	for i := range s {
		if _, ok := nd.move[s[i]]; !ok {
			nd.move[s[i]] = &TrieNode{move: make(map[uint8]*TrieNode)}
		}
		nd = nd.move[s[i]]
	}
	ok := !nd.final
	nd.final = true
	return ok
}

func extractSymbols(src string) []string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, []byte(src), nil, scanner.ScanComments)
	res := make([]string, 0)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.IDENT {
			res = append(res, lit)
		}
	}
	return res
}
