package httputil

import "strings"

type Method string

const (
	Get     Method = `GET`
	Post           = `POST`
	Put            = `PUT`
	Delete         = `DELETE`
	Head           = `HEAD`
	Options        = `OPTIONS`
	Patch          = `PATCH`
)

func IsHttpErr(err error) bool {
	if err != nil && strings.HasPrefix(err.Error(), `HTTP `) {
		return true
	}

	return false
}
