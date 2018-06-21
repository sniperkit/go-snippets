// utils.go
package vxlan

import (
	"reflect"
)

func CompareVNI(vni uint32, netvni [3]byte) int {
	v := [3]byte{byte(vni >> 16 & 0xff),
		byte(vni >> 8 & 0xff),
		byte(vni >> 0 & 0xff),
	}
	if v == netvni {
		return 0
	} else if v[0] > netvni[0] ||
		v[1] > netvni[1] ||
		v[2] > netvni[1] {
		return 1
	}
	return -1
}

func CopyStruct(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
