package responses

import (
	"net/http"
)

//Status is the generic Status which for now just maps to a HTTP Status
type Status int

//Int just casts it to the native int
func (s Status) Int() int {
	return int(s)
}

const (
	_ Status = iota

	//StatusOK refer to http.StatusOK comments
	StatusOK = http.StatusOK

	//StatusCreated refer to http.StatusCreated comments
	StatusCreated = http.StatusCreated

	//StatusNoContent refer to http.StatusNoContent comments
	StatusNoContent = http.StatusNoContent

	//StatusBadRequest refer to http.StatusBadRequest comments
	StatusBadRequest = http.StatusBadRequest

	//StatusUnauthorized refer to http.StatusUnauthorized comments
	StatusUnauthorized = http.StatusUnauthorized

	//StatusForbidden refer to http.StatusForbidden comments
	StatusForbidden = http.StatusForbidden

	//StatusNotFound refer to http.StatusNotFound comments
	StatusNotFound = http.StatusNotFound

	//StatusMethodNotAllowed refer to http.StatusMethodNotAllowed comments
	StatusMethodNotAllowed = http.StatusMethodNotAllowed

	//StatusConflict refer to http.StatusConflict comments
	StatusConflict = http.StatusConflict

	//StatusInternalServerError refer to http.StatusInternalServerError comments
	StatusInternalServerError = http.StatusInternalServerError
)
