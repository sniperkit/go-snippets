package quickstart

import "fmt"

// VersionMajor defines the major version number
const VersionMajor = 0

// VersionMinor defines the minor version number
const VersionMinor = 1

// VersionString returns the complete version number as a string
func VersionString() string {
	return fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
}
