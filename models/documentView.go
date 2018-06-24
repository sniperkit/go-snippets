package models

import (
	"math"
	"sort"
)

type kw struct {
	Word string
	Tf   int
}

// Class for presenting the search results.
type DocumentView struct {
	Title    string
	Uri      string
	Date     string
	Size     string
	Parents  []string
	Children []string
	Keywords []kw
	Tf       []int
	Score    float64
}

func NewDocumentView(d *Document) *DocumentView {
	dv := DocumentView{}
	dv.Title = d.Title
	dv.Uri = d.Uri
	dv.Date = d.GetTimeStr()
	dv.Size = d.GetSizeStr()

	words := make([]string, 0, len(d.Words))
	for k := range d.Words {
		words = append(words, k)
	}

	sort.Slice(words, func(i, j int) bool {
		if d.Words[words[i]].Tf == d.Words[words[j]].Tf {
			return words[i] < words[j]
		}
		return d.Words[words[i]].Tf > d.Words[words[j]].Tf
	})

	upper := int(math.Min(float64(len(d.Links)), 5.0))
	dv.Children = d.Links[0:upper]
	upper = int(math.Min(float64(len(words)), 5.0))
	words = words[0:upper]

	dv.Keywords = make([]kw, 0)
	for _, w := range words {
		keyword := kw{Word: w, Tf: d.Words[w].Tf}
		dv.Keywords = append(dv.Keywords, keyword)
	}
	return &dv
}
