package database

import (
	"fmt"
	"github.com/rsmohamad/comp4321/models"
	"testing"
)

func generateWords(num int) map[string]models.Word {
	words := make([]string, num)
	for i := range words {
		words[i] = fmt.Sprint(num)
	}
	return models.CountTfandIdx(words)
}

func generateDocuments(num int) []*models.Document {
	docs := make([]*models.Document, 0)
	for i := 0; i < num; i++ {
		doc := &models.Document{
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
		doc.Titles = models.CountTfandIdx([]string{doc.Title})
		doc.MaxTf = models.CountMaxTf(doc.Words)
		doc.TitleMaxTf = models.CountMaxTf(doc.Titles)
		docs = append(docs, doc)
	}

	return docs
}

func TestInsertion(t *testing.T) {
	indexer, _ := LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(10)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.Close()

	viewer, _ := LoadViewer("index_test.db")
	viewer.ForEachDocument(func(p *models.Document, i int) {
		if p.Title != fmt.Sprint(i) {
			t.Fail()
		}
	})
	viewer.Close()
}

func TestContains(t *testing.T) {
	indexer, _ := LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(10)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.Close()

	viewer, _ := LoadViewer("index_test.db")

	for i := 0; i < 10; i++ {
		if !viewer.ContainsUrl(fmt.Sprintf("http://%d.com/", i)) {
			t.Fail()
		}
	}

	if viewer.ContainsUrl(fmt.Sprintf("http://%d.com/", 10)) {
		t.Fail()
	}

	viewer.Close()
}

func TestPageRank(t *testing.T) {
	indexer, _ := LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(10)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.UpdateAdjList()
	indexer.UpdatePageRank()
	indexer.Close()

	viewer, _ := LoadViewer("index_test.db")

	pr := viewer.GetPageRank(uint64(1))
	for i := 0; i < 10; i++ {
		if pr != viewer.GetPageRank(uint64(i+1)) {
			t.Fail()
		}
	}

	viewer.Close()
}

func TestAdjList(t *testing.T) {
	indexer, _ := LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(10)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.UpdateAdjList()
	indexer.Close()

	viewer, _ := LoadViewer("index_test.db")

	for i := 0; i < 10; i++ {
		parents := viewer.GetParentLinks(uint64(i + 1))

		if len(parents) != 9 {
			t.Log(len(parents))
			t.Fail()
		}

		for _, parent := range parents {
			if parent == fmt.Sprintf("http://%d.com/", i) {
				t.Fail()
			}
		}
	}

	viewer.Close()
}

func TestTermWeights(t *testing.T) {
	indexer, _ := LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(10)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.UpdateTermWeights()
	indexer.Close()

	viewer, _ := LoadViewer("index_test.db")

	for i := 0; i < 10; i++ {
		if i == 0 {
			continue
		}

		mag := viewer.GetMagnitude(uint64(i+1), true)
		score := viewer.GetTfIdf(uint64(i+1), fmt.Sprint(i), true)
		score /= mag

		if score != 1 {
			t.Log("Score", score)
			t.Fail()
		}
	}

	viewer.Close()
}
