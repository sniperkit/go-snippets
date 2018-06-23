package client

import (
	"gopkg.in/validator.v2"
)

func valError(errMap error, key string) string {
	if errMap == nil {
		return ""
	}

	m := errMap.(validator.ErrorMap)
	if ea, ok := m[key]; ok {
		return ea.Error()
	}

	return ""
}

var validate = validator.Validate
