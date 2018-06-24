package models

import (
	"fmt"
	"math"
	"testing"
)

func generateWords(num int) map[string]Word {
	words := make([]string, num)
	for i := range words {
		words[i] = fmt.Sprint(num)
	}
	return CountTfandIdx(words)
}

func generateDocuments(num int) []*Document {
	docs := make([]*Document, 0)
	for i := 0; i < num; i++ {
		doc := &Document{
			Uri:   fmt.Sprintf("http://%d.com/", i),
			Title: fmt.Sprintf("%d", i),
		}

		// Outgoing links
		for j := 0; j < num; j++ {
			if j == i {
				continue
			}
			doc.Links = append(doc.Links, fmt.Sprintf("http://%d.com/", j))
		}

		// Words
		doc.Words = generateWords(i)
		doc.Titles = CountTfandIdx([]string{doc.Title})
		doc.MaxTf = CountMaxTf(doc.Words)
		doc.TitleMaxTf = CountMaxTf(doc.Titles)
		docs = append(docs, doc)
	}

	return docs
}

func TestNewDocumentView(t *testing.T) {
	docs := generateDocuments(10)

	for _, doc := range docs {
		docView := NewDocumentView(doc)
		if docView == nil {
			t.Log("docview is nil")
			t.Fail()
		}

		if docView.Title != doc.Title || docView.Uri != doc.Uri {
			t.Log("title and url dont match")
			t.Fail()
		}

		if len(docView.Children) != 5 {
			t.Log("children links are incorrect")
			t.Fail()
		}

		min := int(math.Min(5.0, float64(len(doc.Words))))
		if len(docView.Keywords) != min {
			t.Log("keywords are incorrect")
			t.Fail()
		}
	}
}
