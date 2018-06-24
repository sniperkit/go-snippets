package retrieval

import (
	"fmt"
	"github.com/rsmohamad/comp4321/database"
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

func insertIntoIndex(num int) {
	indexer, _ := database.LoadIndexer("index_test.db")
	indexer.DropAll()

	docs := generateDocuments(num)
	for _, doc := range docs {
		indexer.UpdateOrAddPage(doc)
	}

	indexer.FlushInverted()
	indexer.UpdatePageRank()
	indexer.UpdateAdjList()
	indexer.UpdateTermWeights()
	indexer.Close()
}

func TestNewSearchEngine(t *testing.T) {
	insertIntoIndex(10)
	se := NewSearchEngine("index_test.db")
	if se == nil {
		t.Fail()
	} else {
		se.Close()
	}
}

func TestSEngine_RetrieveBoolean(t *testing.T) {
	insertIntoIndex(10)
	se := NewSearchEngine("index_test.db")
	defer se.Close()

	for i := 0; i < 10; i++ {
		res := se.RetrieveBoolean(fmt.Sprint(i))

		if len(res) != 1 {
			t.Fail()
		}

		if res[0].Title != fmt.Sprint(i) {
			t.Fail()
		}
	}
}

func TestSEngine_RetrievePhrase(t *testing.T) {
	insertIntoIndex(10)
	se := NewSearchEngine("index_test.db")
	defer se.Close()

	for i := 0; i < 10; i++ {
		res := se.RetrievePhrase(fmt.Sprintf("\"%d\"", i))

		if len(res) != 1 {
			t.Fail()
		}

		if res[0].Title != fmt.Sprint(i) {
			t.Fail()
		}
	}
}

func TestSEngine_RetrieveVSpace(t *testing.T) {
	insertIntoIndex(10)
	se := NewSearchEngine("index_test.db")
	defer se.Close()

	for i := 0; i < 10; i++ {
		res := se.RetrieveVSpace(fmt.Sprint(i))

		if len(res) != 1 {
			t.Fail()
		}

		if res[0].Title != fmt.Sprint(i) {
			t.Fail()
		}
	}
}
