package maptranslator

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func FromMap(from map[string]interface{}, toType reflect.Type) (interface{}, error) {
	toValue := reflect.New(toType)
	return toValue, mapstructure.Decode(from, toValue)
}

func ToMap(from interface{}) (map[string]interface{}, error) {
	return structs.Map(from), nil
}
