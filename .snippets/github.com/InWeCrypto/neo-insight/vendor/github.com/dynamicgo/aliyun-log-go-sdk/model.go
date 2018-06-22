package sls

// GetHistogramsResponse defines response from GetHistograms call
type SingleHistogram struct {
	Progress string              `json:"progress"`
	Count    int64               `json:"count"`
	From     int64               `json:"from"`
	To       int64               `json:"to"`
}

type GetHistogramsResponse struct {
	Progress string              `json:"progress"`
	Count    int64               `json:"count"`
	Histograms     []SingleHistogram `json:"histograms"`
}

// GetLogsResponse defines response from GetLogs call
type GetLogsResponse struct {
	Progress string              `json:"progress"`
	Count    int64               `json:"count"`
	Logs     []map[string]string `json:"logs"`
}

// IndexKey ...
type IndexKey struct {
	Token         []string `json:"token"` // tokens that split the log line.
	CaseSensitive bool     `json:"caseSensitive"`
	Type          string   `json:"type"` // text, long, double
}

type IndexLine struct {
	Token         []string `json:"token"`
	CaseSensitive bool     `json:"caseSensitive"`
	IncludeKeys   []string `json:"include_keys,omitempty"`
	ExcludeKeys   []string `json:"exclude_keys,omitempty"`
}

// Index is an index config for a log store.
type Index struct {
	TTL  int                 `json:"ttl"`
	Keys map[string]IndexKey `json:"keys,omitempty"`
	Line *IndexLine          `json:"line,omitempty"`
}
