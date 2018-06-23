package gorgonzola

// HTTPError represents custom HTTP error
type HTTPError struct {
	Err     error
	Message string
	Code    int
}

// Error implements the error interface
func (e HTTPError) Error() string {
	return e.Err.Error()
}
