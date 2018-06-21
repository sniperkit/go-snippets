package result

type Result struct {
	Command  string
	Result   string
	Err      error
	Group    string
	SubGroup string
	Feature  string
	Case     string
	Message  string
	Pass     bool
	Task     string
}
