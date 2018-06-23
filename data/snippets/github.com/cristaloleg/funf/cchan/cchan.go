package cchan

import (
	"reflect"
)

// Convert ...
func Convert(in interface{}, f interface{}) (reflect.Value, error) {
	inVal := reflect.ValueOf(in)
	if inVal.Kind() != reflect.Chan {
		panic("in parameter must be a channel")
	}

	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		panic("f parameter must be a function")
	}

	rtype := fVal.Type().Out(0)
	ctype := reflect.ChanOf(reflect.BothDir, rtype)
	outVal := reflect.MakeChan(ctype, 0)

	go func() {
		var elem reflect.Value
		for ok := true; ok; {
			if elem, ok = inVal.Recv(); ok {
				result := fVal.Call([]reflect.Value{elem})
				outVal.Send(result[0])
			}
		}
		outVal.Close()
	}()

	return outVal, nil
}
