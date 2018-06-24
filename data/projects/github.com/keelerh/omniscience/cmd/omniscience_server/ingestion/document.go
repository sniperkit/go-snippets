package ingestion

import (
	"time"

	"github.com/olivere/elastic"
)

// Document is a structure used for serializing/deserializing data in Elasticsearch.
type Document struct {
	Id           string                `json:"id"`
	Title        string                `json:"title,omitempty"`
	Description  string                `json:"description,omitempty"`
	Service      string                `json:"service,omitempty"`
	Content      string                `json:"content,omitempty"`
	Url          string                `json:"url,omitempty"`
	LastModified time.Time             `json:"last_modified,omitempty"`
	Suggest      *elastic.SuggestField `json:"suggest_field,omitempty"`
}
