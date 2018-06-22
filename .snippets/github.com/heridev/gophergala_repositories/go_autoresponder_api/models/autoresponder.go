package models

type Autoresponder struct {
  Id           int64
  Title        string `db:"title" json:"title"`
  Description  string `db:"description" json:"description"`
  Lists        []EmailList `json:"lists"`
}
