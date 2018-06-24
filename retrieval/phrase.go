package retrieval

import (
	"github.com/rsmohamad/comp4321/database"
	"log"
	"sort"
)

type Bigram struct {
	n1, n2 string
}

// Split a set of words into bigrams.
func splitToBigrams(query []string) (bigrams []Bigram) {
	for i := 0; (i + 1) < len(query); i++ {
		bg := Bigram{query[i], query[i+1]}
		bigrams = append(bigrams, bg)
	}
	return
}

func hasBigram(id uint64, bigram Bigram, viewer *database.Viewer, title bool) bool {
	pos1 := viewer.GetPositionIndices(id, bigram.n1, title)
	pos2 := viewer.GetPositionIndices(id, bigram.n2, title)

	for i := range pos2 {
		pos2[i]--
	}

	common := intersect(pos1, pos2)
	return len(common) > 0
}

// Returns docIds that contain the bigram phrase.
func hasPhrase(bigram Bigram, viewer *database.Viewer) []uint64 {
	docIds := booleanFilter([]string{bigram.n1, bigram.n2}, viewer)
	rv := make([]uint64, 0)

	for _, id := range docIds {
		inBody := hasBigram(id, bigram, viewer, false)
		inTitle := hasBigram(id, bigram, viewer, true)

		if inBody || inTitle {
			rv = append(rv, id)
		}
	}

	return rv
}

// Treat the query as a phrase and returns docIds containing the phrase.
// Changes the query into bigrams and find documents containing all bigrams.
func filterPhrase(query []string, viewer *database.Viewer) []uint64 {
	if len(query) <= 1 {
		return booleanFilter(query, viewer)
	}

	bigrams := splitToBigrams(query)
	log.Println("Bigrams:", bigrams)
	docWithBigrams := make([][]uint64, 0)

	for _, bigram := range bigrams {
		docWithBigrams = append(docWithBigrams, hasPhrase(bigram, viewer))
	}

	sort.Slice(docWithBigrams, func(i, j int) bool {
		return len(docWithBigrams[i]) < len(docWithBigrams[j])
	})

	docIDs := make([]uint64, 0)
	docIDs = append(docIDs, docWithBigrams[0]...)
	for i, docs := range docWithBigrams {
		if i == 0 {
			continue
		}
		docIDs = intersect(docIDs, docs)
	}

	return docIDs
}

func retrievePhrase(phrases []string, query string, viewer *database.Viewer) (map[uint64]float64, []uint64) {
	docIds := make([]uint64, 0)
	for _, phrase := range phrases {
		preprocessed := preprocessText(phrase)
		docIds = append(docIds, filterPhrase(preprocessed, viewer)...)
	}

	preprocessed := preprocessText(query)
	return getDocumentScores(preprocessed, viewer, docIds)
}
