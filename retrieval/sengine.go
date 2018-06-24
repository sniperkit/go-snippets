package retrieval

import (
	"github.com/rsmohamad/comp4321/database"
	"github.com/rsmohamad/comp4321/models"
	"github.com/rsmohamad/comp4321/stopword"
	"github.com/surgebase/porter2"
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
)

func preprocessText(query string) (rv []string) {
	regex := regexp.MustCompile("[^a-zA-Z0-9 ]")
	query = regex.ReplaceAllString(query, " ")
	regex = regexp.MustCompile("[^\\s]+")
	words := regex.FindAllString(query, -1)
	for _, word := range words {
		cleaned := strings.ToLower(word)
		cleaned = strings.TrimSpace(cleaned)
		cleaned = porter2.Stem(cleaned)
		if !stopword.IsStopWord(cleaned) {
			rv = append(rv, cleaned)
		}
	}
	return
}

func extractPhrases(query string) []string {
	phrases := make([]string, 0)

	startPhrase := -1
	for i, char := range query {
		if char == '"' && startPhrase == -1 {
			startPhrase = i
			continue
		}

		if char == '"' && startPhrase > -1 {
			phrases = append(phrases, query[startPhrase+1:i])
			startPhrase = -1
		}
	}

	return phrases
}

type SEngine struct {
	viewer *database.Viewer
}

func NewSearchEngine(filename string) *SEngine {
	se := SEngine{}
	var err error
	se.viewer, err = database.LoadViewer(filename)
	if err != nil {
		log.Fatal("Index file not found:", filename)
	}

	return &se
}

func (e *SEngine) getDocumentViewModels(docIds []uint64, scores map[uint64]float64) []*models.DocumentView {
	rv := make([]*models.DocumentView, len(docIds))
	for i, id := range docIds {
		doc := e.viewer.GetDocument(id)
		if doc != nil {
			docView := models.NewDocumentView(doc)
			if scores == nil {
				docView.Score = 1
			} else {
				docView.Score = scores[id]
			}
			parents := e.viewer.GetParentLinks(id)
			upper := int(math.Min(float64(len(parents)), 5.0))
			docView.Parents = parents[0:upper]
			rv[i] = docView
		}
	}
	return rv
}

func (e *SEngine) RetrieveBoolean(query string) []*models.DocumentView {
	preprocessed := preprocessText(query)
	docIds := booleanFilter(preprocessed, e.viewer)
	return e.getDocumentViewModels(docIds, nil)
}

func (e *SEngine) RetrievePhrase(query string) []*models.DocumentView {
	phrases := extractPhrases(query)
	if len(phrases) == 0 {
		return e.RetrieveVSpace(query)
	}

	scores, ids := retrievePhrase(phrases, query, e.viewer)
	sort.Slice(ids, func(i, j int) bool {
		return scores[ids[i]] > scores[ids[j]]
	})

	upper := int(math.Min(50.0, float64(len(ids))))
	return e.getDocumentViewModels(ids[0:upper], scores)
}

func (e *SEngine) RetrieveVSpace(query string) []*models.DocumentView {
	preprocessed := preprocessText(query)
	scores, docIds := vspaceRetrieval(preprocessed, e.viewer)

	sort.Slice(docIds, func(i, j int) bool {
		return scores[docIds[i]] > scores[docIds[j]]
	})

	upper := int(math.Min(50.0, float64(len(docIds))))
	return e.getDocumentViewModels(docIds[0:upper], scores)
}

// Search for needle within the results of haystack
func (e *SEngine) RetrieveNested(haystack, needle string) []*models.DocumentView {
	searchKeyword := func(query string) (map[uint64]float64, []uint64) {
		phrases := extractPhrases(query)
		if len(phrases) == 0 {
			preprocessed := preprocessText(query)
			return vspaceRetrieval(preprocessed, e.viewer)
		}
		return retrievePhrase(phrases, query, e.viewer)
	}

	sortByScore := func(ids []uint64, scores map[uint64]float64) {
		sort.Slice(ids, func(i, j int) bool {
			return scores[ids[i]] > scores[ids[j]]
		})
	}

	sortByIds := func(ids []uint64) {
		sort.Slice(ids, func(i, j int) bool {
			return ids[i] < ids[j]
		})
	}

	scores, haystackIds := searchKeyword(haystack)
	_, needleIds := searchKeyword(needle)

	sortByScore(haystackIds, scores)
	upper := int(math.Min(50.0, float64(len(haystackIds))))
	haystackIds = haystackIds[0:upper]

	sortByIds(haystackIds)
	sortByIds(needleIds)
	combined := intersect(haystackIds, needleIds)
	sortByScore(combined, scores)
	return e.getDocumentViewModels(combined, scores)
}

func (e *SEngine) RetrievePageRank(query string) []*models.DocumentView {
	preprocessed := preprocessText(query)
	scores, docIds := vspaceRetrieval(preprocessed, e.viewer)

	sort.Slice(docIds, func(i, j int) bool {
		return scores[docIds[i]] > scores[docIds[j]]
	})

	upper := int(math.Min(50.0, float64(len(docIds))))
	docIds = docIds[0:upper]

	pageRanks := make(map[uint64]float64)
	for _, docId := range docIds {
		pageRanks[docId] = e.viewer.GetPageRank(docId)
	}

	sort.Slice(docIds, func(i, j int) bool {
		return pageRanks[docIds[i]] > pageRanks[docIds[j]]
	})

	return e.getDocumentViewModels(docIds, pageRanks)
}

func (e *SEngine) Close() {
	e.viewer.Close()
}
