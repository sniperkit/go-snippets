package toolkit

import (
	"encoding/gob"
	"errors"
	"reflect"
	"strings"
	"time"
)

var gobs []string

type G interface{}

func RegisterGobObject(o interface{}) {
	name := reflect.ValueOf(o).Type().Name()
	if HasMember(gobs, name) {
		return
	}

	gob.Register(o)
	gobs = append(gobs, name)
}

func TypeName(o interface{}) string {
	v := reflect.ValueOf(o)
	return v.Type().String()
	/*
		v := reflect.ValueOf(o)
		if !v.IsValid() {
			return ""
		}

		var t reflect.Type
		if v.Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
			if v.Kind() == reflect.Ptr {
				return "*" + TypeName(v.Elem().Interface())
			} else {
				t = v.Type()
			}
		} else {
			t = reflect.TypeOf(o)
		}

		//t := v.Type()
		name := t.Name()
		pkg := t.PkgPath()
		if pkg != "" {
			return pkg + "." + name
		} else {
			return name
		}
	*/
}

func IsNil(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.Indirect(reflect.ValueOf(x))
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}

func IsNilOrEmpty(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.Indirect(reflect.ValueOf(x))
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()

		totalPublicProps := 0
		totalPublicPropsNotEmpty := 0

		for i := v.NumField() - 1; i >= 0; i-- {
			if vt.Field(i).PkgPath != "" {
				continue // Private field
			}

			totalPublicProps++
			if !IsNilOrEmpty(v.Field(i).Interface()) {
				totalPublicPropsNotEmpty++
			}
		}

		// has few public properties, but all of them is empty.
		// example: we store time.Time data in session/M, when we trying to get the data the returned value is always nil, because in previous commit, only if there is at least 1 property which is not empty, the value marked as not nil. but time.Time type doesn't have any public properties, this condition causing returned value always nil.
		if totalPublicProps > 0 && totalPublicPropsNotEmpty == 0 {
			return true
		}
	}

	return false
}

func IsNumber(o interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(o))
	ts := strings.ToLower(v.Type().String())
	if (strings.Contains(ts, "int") || strings.Contains(ts, "float")) && !strings.HasPrefix(ts, "interface{}") {
		return true
	}
	return false
}

func IsStringNumber(txt string, decsep string) (f float64, e error) {
	hasDes := false
	newtxt := "0"
	decPoint := 0
	//Printf("%v ", txt)
	for _, c := range txt {
		s := string(c)
		if strings.Compare(s, "0") >= 0 && strings.Compare(s, "9") <= 0 {
			newtxt += s
			if hasDes {
				decPoint += 1
			}
		} else if s == decsep {
			if !hasDes {
				newtxt += "."
				hasDes = true
			} else {
				e = errors.New("IsStringNumber: Multiple decimal separator found")
				return
			}
		} else {
			//Printfn("%v %v", txt, s)
			e = errors.New("IsStringNumber: Wrong character " + txt)
			return
		}
	}
	if strings.HasSuffix(newtxt, ".") {
		newtxt += "0"
	}
	//Printfn("%v",newtxt)
	f = ToFloat64(newtxt, decPoint, RoundingAuto)
	return
}

func IsPointer(o interface{}) bool {
	v := reflect.ValueOf(o)
	return v.Kind() == reflect.Ptr
}

func IsSlice(o interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(o))
	return v.Kind() == reflect.Slice
}

func GetEmptySliceElement(o interface{}) (interface{}, error) {
	rv := reflect.Indirect(reflect.ValueOf(o))
	if rv.Kind() != reflect.Slice {
		return nil, errors.New("GetEmptySliceElement: " + TypeName(o) + " is not a slice")
	}
	sliceType := rv.Type().Elem()
	newelem := reflect.New(sliceType)
	//Println(newelem.Type().String())
	if string(sliceType.String()[0]) == "*" {
		return Value2Interface(newelem), nil
	} else {
		return Value2Interface(newelem.Elem()), nil
	}
}

func AppendSlice(o interface{}, v interface{}) error {
	rvPtr := reflect.ValueOf(o)
	if rvPtr.Kind() != reflect.Ptr {
		return errors.New("AppendSlice: Is not a pointer of slice")
	}
	rv := reflect.Indirect(rvPtr)
	if rv.Kind() != reflect.Slice {
		return errors.New("AppendSlice: " + rv.Kind().String() + " is not a pointer of slice")
	}
	rv = reflect.Append(rv, reflect.ValueOf(v))
	rvPtr.Elem().Set(rv)
	return nil
}

func MakeSlice(o interface{}) interface{} {
	t := reflect.TypeOf(o)
	//fmt.Printf("Type: %s \n", t.String())
	return reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
}

func SliceLen(o interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return 0
	}
	return v.Len()
}

func SliceSubset(o interface{}, lbound, ubound int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	l := v.Len()
	if lbound < l && ubound < l {
		var arrays reflect.Value
		for i := lbound; i <= ubound; i++ {
			elem := v.Index(i)
			if i == lbound {
				arrays = reflect.MakeSlice(elem.Type(), 0, 0)
			}
			arrays = reflect.Append(arrays, elem)
		}
		return arrays.Interface()
	}
	return nil
}

func MapKeys(o interface{}) []interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Map {
		return []interface{}{}
	}

	var ret []interface{}
	keyvalues := v.MapKeys()
	for _, kv := range keyvalues {
		ret = append(ret, kv.Interface())
	}
	return ret
}

func MapLen(o interface{}) int {
	return len(MapKeys(o))
}

func MapItem(o interface{}, key interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Map {
		return nil
	}
	item := v.MapIndex(reflect.ValueOf(key))
	return item.Interface()
}

func SliceItem(o interface{}, index int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return nil
	}
	if v.Len()-1 < index {
		return nil
	}
	return Value2Interface(v.Index(index))
}

func SliceSetItem(o interface{}, i int, d interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		err := "SliceSetItem: object returned is not a slice. It is " + v.Kind().String() + " " + v.Type().String()
		//Println(err)
		return errors.New(err)
	}
	currentLen := v.Len()
	if i >= currentLen {
		//i = currentLen + 1
		//Println("Set capacity to ", i+1)
		//v.SetCap(i + 1)
		v1 := reflect.Append(v, reflect.ValueOf(d))
		v.Set(v1)
	} else {
		v.Index(i).Set(reflect.ValueOf(d))
	}
	return nil
}

func Serde(o interface{}, dest interface{}, serdeType string) error {
	vt1 := reflect.TypeOf(o)
	vt2 := reflect.TypeOf(dest)
	vt1kind := vt1.Kind()
	if vt1kind == reflect.Ptr {
		vt1kind = reflect.ValueOf(o).Elem().Type().Kind()
	}

	vt1name := vt1.String()
	vt2name := vt2.String()
	if vt1name == vt2name && vt1kind != reflect.Map && vt1kind != reflect.Slice {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(o).Elem())
	}

	bs, e := ToBytesWithError(o, serdeType)
	if len(bs) == 0 {
		return errors.New("toolkit.Serde: Serialization Fail " + e.Error())
	}

	e = FromBytes(bs, serdeType, dest)
	if e != nil {
		return errors.New("toolkit.Serde: Deserialization fail " + e.Error())
	}

	return nil
}

func Value2Interface(vi reflect.Value) interface{} {
	vik := vi.Type().String()
	if vik == "string" || vik == "String" {
		return vi.String()
	} else if (vik == "int" || vik == "int8" || vik == "int16" || vik == "int32" || strings.Contains(vik, "uint")) && !strings.Contains(vik, "interface") {
		return int(vi.Int())
	} else if strings.Contains(vik, "float") {
		return vi.Float()
	} else if strings.Contains(vik, "bool") {
		return vi.Bool()
	} else {
		//Printfn("data: %s", JsonString(vi.Interface()))
		return vi.Interface()
	}
}

func ExecFunc(fn interface{}, ins ...interface{}) (outs []reflect.Value, e error) {
	rvfn := reflect.ValueOf(fn)
	if rvfn.Kind() != reflect.Func {
		e = errors.New("Execfunc: First parameter should be a function")
		return
	}
	var rvins []reflect.Value
	for _, in := range ins {
		rvins = append(rvins, reflect.ValueOf(in))
	}
	outs = rvfn.Call(rvins)
	return
}

func ExecuteBlockWithTimeout(callback func() interface{}, timeout time.Duration) (interface{}, bool) {
	ch := make(chan interface{}, 1)

	go func() {
		ch <- callback()
	}()

	select {
	case res := <-ch:
		return res, true
	case <-time.After(timeout * time.Second):
		return nil, false
	}
}
