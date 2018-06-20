package job

// represents a task handler that receives a payload and returns a generic type response or error
type Task interface {
	Run(payload Payload) JobResult
}
