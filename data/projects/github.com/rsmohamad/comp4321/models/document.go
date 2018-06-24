package models

import (
	"code.cloudfoundry.org/bytefmt"
	"time"
)

type Word struct {
	Tf        int
	Positions []int
}

// Document class for representation inside the system.
// Has fields relevant to search execution but not presentation.
type Document struct {
	Title      string
	Uri        string
	Links      []string
	Words      map[string]Word
	Titles     map[string]Word
	Len        int
	MaxTf      int
	TitleMaxTf int
	Modtime    int64
}

func (d Document) GetSizeStr() string {
	if d.Len == 0 {
		return "Not available"
	}

	return bytefmt.ByteSize(uint64(d.Len))
}

func (d Document) GetTimeStr() string {
	if d.Modtime < 0 {
		return "No date available"
	}

	t := time.Unix(d.Modtime, 0)
	return t.Format("02 Jan 2006")
}

func CountTfandIdx(words []string) map[string]Word {
	m := make(map[string]Word)
	for position, word := range words {
		wordModel := m[word]
		wordModel.Tf++
		wordModel.Positions = append(wordModel.Positions, position)
		m[word] = wordModel
	}
	return m
}

func CountMaxTf(words map[string]Word) int {
	max := 0
	for _, word := range words {
		if word.Tf > max {
			max = word.Tf
		}
	}
	return max
}
