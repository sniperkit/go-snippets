package plugin

type Status int

// Supported exit statuses
const (
	OK Status = iota
	WARNING
	CRITICAL
	UNKNOWN
)

// ExitCode returns current status as integer
func (st Status) ExitCode() int {
	return int(st)
}

// String returns current status as string
func (st Status) String() string {
	switch st {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
