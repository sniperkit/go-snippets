// Validator provide a validator to validate fields of struct.
//
// Author: Mephis Pheies
package protocols

import (
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

func init() {
	binding.Validator.RegisterValidation("serversion", serviceVersion)
}

// Validator is used to validate struct field type
var Validator = binding.Validator

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

func serviceVersion(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if version, ok := field.Interface().(string); ok {
		return reVersion.MatchString(version)
	}
	return false
}
