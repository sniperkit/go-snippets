package job

type Transformations interface {
	Handler // Submit and Confirm are imported from here.
	Error(func(Job)) // Error function whenever corner cases are not met.
}