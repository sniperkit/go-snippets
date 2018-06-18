package rest

import (
	"github.com/bradberger/context"
)

// LogFunc is a custom log function type
type LogFunc func(ctx context.Context, format string, args ...interface{})
