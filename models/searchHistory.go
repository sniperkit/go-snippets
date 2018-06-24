package models

import "time"

// Only exported fields are serialized.
// Therefore all fields are serialized.
type SearchHistory struct {
	Query string
	Time  int64
}

func NewSearchHistory(query string) SearchHistory {
	tm := time.Now()
	return SearchHistory{Query: query, Time: tm.Unix()}
}

func (d SearchHistory) GetQuery() string {
	return d.Query
}

func (d SearchHistory) GetTime() string {
	if d.Time < 0 {
		return "No date available"
	}

	t := time.Unix(d.Time, 0)
	return t.Format(time.RFC822)
}
