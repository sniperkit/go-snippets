package models

type TimestampVersion struct {
	Version string `json:"version"`
}

type MetadataField struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type Metadata []MetadataField

type InRequest struct {
	Source  Source  `json:"source"`
	Version TimestampVersion `json:"version"`
}

type InResponse struct {
	Version  TimestampVersion  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version  TimestampVersion  `json:"version"`
}

type CheckResponse []TimestampVersion

type Source struct {}

