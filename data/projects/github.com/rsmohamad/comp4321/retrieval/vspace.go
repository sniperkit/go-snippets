package retrieval

import (
	"github.com/rsmohamad/comp4321/database"
	"math"
)

type CosSimResult struct {
	score float64
	docId uint64
}

func cosSim(query []string, docId uint64, viewer *database.Viewer, res *chan *CosSimResult) *CosSimResult {
	var textInnerProduct float64 = 0
	var titleInnerProduct float64 = 0
	queryMag := math.Sqrt(float64(len(query)))
	docMag := viewer.GetMagnitude(docId, false)
	titleMag := viewer.GetMagnitude(docId, true)

	for _, word := range query {
		textInnerProduct += viewer.GetTfIdf(docId, word, false)
		titleInnerProduct += viewer.GetTfIdf(docId, word, true)
	}

	textScore := textInnerProduct / (queryMag * docMag)
	titleScore := titleInnerProduct / (queryMag * titleMag)
	score := textScore
	if !math.IsNaN(titleScore) {
		score += titleScore * 1.5
	}

	rv := &CosSimResult{score, docId}
	*res <- rv
	return rv
}

func getDocumentScores(query []string, viewer *database.Viewer, docsToSearch []uint64) (map[uint64]float64, []uint64) {
	documentScores := make(map[uint64]float64)
	documentIds := make([]uint64, 0)
	res := make(chan *CosSimResult)
	defer close(res)

	for _, id := range docsToSearch {
		_, exist := documentScores[id]
		if !exist {
			go cosSim(query, id, viewer, &res)
			documentScores[id] = 0
			documentIds = append(documentIds, id)
		}
	}

	for range documentIds {
		result := <-res
		documentScores[result.docId] = result.score
	}

	return documentScores, documentIds
}

func vspaceRetrieval(query []string, viewer *database.Viewer) (map[uint64]float64, []uint64) {
	docsToSearch := make([]uint64, 0)
	res := make(chan []uint64)

	for _, word := range query {
		go func(word string) {
			ids := booleanFilter([]string{word}, viewer)
			res <- ids
		}(word)
	}

	for range query {
		docsToSearch = append(docsToSearch, <-res...)
	}

	return getDocumentScores(query, viewer, docsToSearch)
}
