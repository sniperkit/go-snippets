package slicerjob

import "code.google.com/p/go-uuid/uuid"

type Job struct {
	ID       string  `json:"id"`
	Status   Status  `json:"status"`
	Progress float64 `json:"progress"`
	URL      string  `json:"url"`
	GCodeURL string  `json:"gcode_url"`
}

type SlicerPreset struct {
	Slicer  string   `json:slicer`
	Presets []string `json:presets`
}

// New creates a new Job with a random UUID for an ID.  If urlformat is
// non-empty the URL of the returned job is computed as
// fmt.Sprintf(urlformat,job.ID).
func New() *Job {
	job := new(Job)
	job.ID = uuid.New()
	return job
}
