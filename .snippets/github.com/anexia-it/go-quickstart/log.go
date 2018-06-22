package quickstart

import "go.uber.org/zap"

// rootLogger holds the root logger instance and is a package-private variable
var rootLogger *zap.Logger

// GetLogger returns a named logger
func GetLogger(name string) *zap.Logger {
	return rootLogger.Named(name)
}

func init() {
	// init is called during program initialization

	// Create a new logger and ignore possible errors
	rootLogger, _ = zap.NewDevelopment()
}
