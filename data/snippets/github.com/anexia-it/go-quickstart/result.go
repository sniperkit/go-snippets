package quickstart

// Result holds a message queue result
type Result struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Error string `json:"error,omitempty"`
}
