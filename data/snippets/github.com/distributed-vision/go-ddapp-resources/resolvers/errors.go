package resolvers

import "fmt"

type EntityNotFound struct {
	reason string
	cause  error
}

func (e *EntityNotFound) Error() string {
	if e.cause != nil {
		return fmt.Sprint(e.reason, e.cause.Error())
	}

	return e.reason
}

func NewEntityNotFound(reason string, cause error) *EntityNotFound {
	return &EntityNotFound{reason, cause}
}
