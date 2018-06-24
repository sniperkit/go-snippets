package retrieval

import (
	"github.com/rsmohamad/comp4321/database"
	"sort"
)

func intersect(list1, list2 []uint64) (answer []uint64) {
	i := 0
	j := 0

	for i != len(list1) && j != len(list2) {
		if list1[i] == list2[j] {
			answer = append(answer, list1[i])
			i++
			j++
		} else if list1[i] < list2[j] {
			i++
		} else {
			j++
		}
	}
	return
}

func booleanFilter(query []string, viewer *database.Viewer) (docIDs []uint64) {
	if len(query) == 0 {
		return
	}

	wordDoc := make(map[string][]uint64)
	for _, word := range query {
		wordDoc[word] = viewer.GetContainingPages(word)
	}

	sort.Slice(query, func(i, j int) bool {
		return len(wordDoc[query[i]]) < len(wordDoc[query[j]])
	})

	docIDs = append(docIDs, wordDoc[query[0]]...)
	for i, word := range query {
		if i == 0 {
			continue
		}
		docIDs = intersect(docIDs, wordDoc[word])
	}

	return
}
