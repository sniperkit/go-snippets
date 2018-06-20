package job

// Job represents the model to run. Task is the executing algorithm for this job
// payload are the params for the job task and ReturnChannel is the channel to return an output
type Job struct {
	Task          Task
	Payload       Payload
	ReturnChannel chan JobResult
}

// Job represents a return type of a job either the result (Value) or Error
type JobResult struct {
	Value interface{}
	Error error
}

// Runs the Job Task
func (j *Job) Run() JobResult {
	return j.Task.Run(j.Payload)
}

func NewJob(task Task, params map[string]string, returnChannel chan JobResult) Job {
	return Job{
		Task:          task,
		Payload:       NewPayload(params),
		ReturnChannel: returnChannel,
	}
}

func NewJobResult(value interface{}, err error) JobResult {
	return JobResult{
		Value: value,
		Error: err,
	}
}

func NewJobResultChannel() chan JobResult {
	return make(chan JobResult)
}
