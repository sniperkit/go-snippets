package responses

//Error contains a error message
type Error struct {
	Error string `json:",required"`
}

func NewError(errorString string) *Error {
	return &Error{
		Error: errorString,
	}
}
