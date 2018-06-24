// TODO - remove this in favour of either returning a generic "Not Found" string
// or putting not found in a common package
package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	HttpCode int    `json:"code"`
	Message  string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.HttpCode, e.Message)
}

// NewFmtError creates a new error using fmt.Sprintf
func NewFmtError(httpCode int, format string, a ...interface{}) error {
	return &Error{
		HttpCode: httpCode,
		Message:  fmt.Sprintf(format, a...),
	}
}

// New500Error sets the http code to 500 & makes a new fmt error
func New500Error(format string, a ...interface{}) error {
	return &Error{
		HttpCode: http.StatusInternalServerError,
		Message:  fmt.Sprintf(format, a...),
	}
}

// Error500IfErr wraps generic error responses, turning
// error into *Error with HttpCode : 500 if they exist,
// or passing nil through if no error exists
func Error500IfErr(err error) error {
	if err == nil {
		return nil
	}

	// if this error is already *Error, let it pass through
	if _, ok := err.(*Error); ok {
		return err
	}

	return &Error{
		HttpCode: http.StatusInternalServerError,
		Message:  err.Error(),
	}
}

var (
	ErrNotFound     = &Error{HttpCode: http.StatusNotFound, Message: "Not Found"}
	ErrAccessDenied = &Error{HttpCode: http.StatusForbidden, Message: "Access Denied"}

	Err500                  = &Error{HttpCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrNoConnection         = &Error{HttpCode: http.StatusInternalServerError, Message: "no valid database connection"}
	ErrDatabaseDoesNotExist = &Error{HttpCode: http.StatusInternalServerError, Message: "database does not exist"}
	ErrDatabaseExists       = &Error{HttpCode: http.StatusInternalServerError, Message: "database already exists"}

	ErrEmailDoesntExist             = &Error{HttpCode: http.StatusBadRequest, Message: "Email does not exist"}
	ErrEmailRequired                = &Error{HttpCode: http.StatusBadRequest, Message: "Email is required"}
	ErrEmailTaken                   = &Error{HttpCode: http.StatusBadRequest, Message: "Email already exists"}
	ErrInvalidEmail                 = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Email"}
	ErrInvalidKey                   = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Key"}
	ErrInvalidName                  = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Name"}
	ErrInvalidParent                = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Parent"}
	ErrInvalidPassword              = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Password"}
	ErrInvalidUser                  = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid User"}
	ErrInvalidUsername              = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Username"}
	ErrInvalidUserNamePasswordCombo = &Error{HttpCode: http.StatusBadRequest, Message: "Invalid Username / Password Combination"}
	ErrNameRequired                 = &Error{HttpCode: http.StatusBadRequest, Message: "Name is Required"}
	ErrNoIdentifier                 = &Error{HttpCode: http.StatusBadRequest, Message: "No Identifier provided"}
	ErrOwnerRequired                = &Error{HttpCode: http.StatusBadRequest, Message: "Owner is Required"}
	ErrPasswordRequired             = &Error{HttpCode: http.StatusBadRequest, Message: "Password is Required"}
	ErrPasswordTooShort             = &Error{HttpCode: http.StatusBadRequest, Message: "Password is too short"}
	ErrTokenAlreadyUsed             = &Error{HttpCode: http.StatusBadRequest, Message: "This token has already been used"}
	ErrTokenExpired                 = &Error{HttpCode: http.StatusBadRequest, Message: "This token has expired"}
	ErrUsernameRequired             = &Error{HttpCode: http.StatusBadRequest, Message: "Username is Required"}
	ErrUsernameTaken                = &Error{HttpCode: http.StatusBadRequest, Message: "Username already exists"}
	ErrUserNotFound                 = &Error{HttpCode: http.StatusNotFound, Message: "user not found"}
	ErrUserRequired                 = &Error{HttpCode: http.StatusBadRequest, Message: "User is Required"}
)
