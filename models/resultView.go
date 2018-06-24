package models

type ResultView struct {
	Query        string
	Results      []*DocumentView
	TotalResults int
}
