package insight

import "fmt"

// Errors .
const (
	JSONRPCParserError    = -32700
	JSONRPCInvalidRequest = -32600
	JSONRPCMethodNotFound = -32601
	JSONRPCInvalidParams  = -32602
	JSONRPCInnerError     = -32603
)

// JSONRPCError .
type JSONRPCError struct {
	ID      int
	Message string
}

func errorf(id int, fmtstr string, args ...interface{}) *JSONRPCError {
	return &JSONRPCError{
		ID:      id,
		Message: fmt.Sprintf(fmtstr, args...),
	}
}
