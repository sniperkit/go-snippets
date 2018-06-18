// 复制对象

package helper

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type fieldMap map[string]reflect.StructField

// CopySlice 复制数组
func CopySlice(src, dest interface{}) error {
	from := reflect.ValueOf(src)
	if from.Kind() != reflect.Ptr {
		return errors.New("src must be point to slice")
	}
	to := reflect.ValueOf(dest)
	if to.Kind() != reflect.Ptr {
		return errors.New("dest must be point to slice")
	}
	fromSlice := from.Elem()
	if fromSlice.Kind() != reflect.Slice {
		return errors.New("src must be point to slice")
	}
	toSlice := to.Elem()
	if toSlice.Kind() != reflect.Slice {
		return errors.New("dest must be point to slice")
	}

	//  fromType := fromSlice.Elem()
	toType := toSlice.Type().Elem()
	isptr := false
	if toType.Kind() == reflect.Ptr {
		isptr = true
		toType = toType.Elem()
	}
	for i := 0; i < fromSlice.Len(); i++ {
		fromObj := fromSlice.Index(i)
		toObj := reflect.New(toType).Elem()
		if fromObj.Type().ConvertibleTo(toObj.Type()) {
			toObj.Set(fromObj.Convert(toObj.Type()))
		} else {
			if fromObj.Kind() == reflect.Ptr {
				fromObj = fromObj.Elem()
			}
			copyValue(fromObj, toObj)
		}
		if isptr {
			toSlice.Set(reflect.Append(toSlice, toObj.Addr()))
		} else {
			toSlice.Set(reflect.Append(toSlice, toObj))
		}
	}
	return nil
}

// CopyObject 复制单个对象
func CopyObject(src, dest interface{}) error {
	from := reflect.ValueOf(src)
	if from.Kind() != reflect.Ptr {
		return errors.New("src must be point to struct")
	}
	to := reflect.ValueOf(dest)
	if to.Kind() != reflect.Ptr {
		return errors.New("dest must be point to struct")
	}
	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return nil
	}
	if from.Type().ConvertibleTo(to.Type()) {
		to.Elem().Set(from.Convert(to.Type()).Elem())
		return nil
	}

	from, to = from.Elem(), to.Elem()
	copyValue(from, to)
	return nil
}

func copyValue(from, to reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	fromFields := getFields(from.Type())
	toFields := getFields(to.Type())
	for name, toType := range toFields {
		fromType, ok := fromFields[name]
		if !ok {
			continue
		}
		toFV := to.FieldByIndex(toType.Index)
		if !toFV.CanSet() {
			continue
		}

		fromFV := from.FieldByIndex(fromType.Index)
		if fromType.Type.ConvertibleTo(toType.Type) {
			if toType.Type.Kind() == reflect.String && fromType.Type.Kind() != reflect.String {
				toFV.SetString(fmt.Sprintf("%s", fromFV.Interface()))
			} else {
				toFV.Set(fromFV.Convert(toType.Type))
			}
		}
	}
}

func getFields(t reflect.Type) fieldMap {
	fields := make(fieldMap, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields[strings.ToLower(t.Field(i).Name)] = t.Field(i)
	}
	return fields
}
