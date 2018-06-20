package job

// Interface to represent the abstractions of the JobHandler
type Handler interface {
	Confirm(func(Job))
	Process(func(Job))
	Submit(Job)
}