package toolkit

import "reflect"
import "errors"

func MtoStruct(from M, to interface{}) error {
	bts := Jsonify(from)
	return Unjson(bts, &to)
}

func StructToM(from interface{}, to *M) error {
	bts := Jsonify(from)
	return Unjson(bts, to)
}

func SetPropByName(target interface{}, propName string, value interface{}) error {
	s := reflect.ValueOf(target).Elem()

	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	} else {
		return errors.New("first param should be pointer of struct object")
	}

	f := s.FieldByName(propName)

	if f.IsValid() && f.CanSet() {
		newValue := reflect.ValueOf(value)
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			newValue = newValue.Elem()
		}
		f.Set(newValue)
	}

	return nil
}
