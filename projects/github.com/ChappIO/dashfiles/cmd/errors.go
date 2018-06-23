package cmd

import (
	"fmt"
	"os"
)

const UNKNOWN_ERROR = 1

// Environment setup issues (10+)
const (
	// The machine does not have git installed
	GIT_NOT_INSTALLED int = 11
)

// Operation failure issues (100+)
const (
	// Something went wrong while cloning the workspace
	WORKSPACE_CLONE_FAILED = 101
	INIT_FAILED            = 102
	UPDATE_FAILED          = 103
	INSTALL_PACKAGE_FAILED = 104
	EXTERNAL_SCRIPT_FAILED = 105
)

func fatal(code int, message string) {
	fmt.Println(message)
	os.Exit(code)
}
