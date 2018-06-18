package tasks

//NewError creates a new Error
func NewError(err error, statusCode int) *Error {
	return &Error{
		err:        err,
		statusCode: statusCode,
	}
}

//NewErrorDefaultStatus creates a new Error with a default Status Code
func NewErrorDefaultStatus(err error) *Error {
	return NewError(err, 500)
}

//Error contains an Error and an optional StatusCode
type Error struct {
	err        error
	statusCode int
}

//GetError returns the private err field
func (e *Error) GetError() error {
	return e.err
}

//GetStatusCode returns the private statusCode field
func (e *Error) GetStatusCode() int {
	return e.statusCode
}
