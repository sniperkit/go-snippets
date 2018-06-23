package toolkit

import (
	"bytes"
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"net"
	"reflect"
	//"strconv"
	//"fmt"
	//. "strconv"
	"strings"
	"time"
)

func init() {
	//Printf("Registering M")
	RegisterGobObject(&M{})
}

/*
func ToInt(i interface{}) int {
	switch i.(type) {
	case string:
		iv, e := Atoi(i.(string))
		if e != nil {
			return 0
		}
		return iv
	case int:
		return i.(int)
	case int32, int64:
		return int(i.(int64))
	case float32:
		return int(i.(float32))
	case float64:
		return int(i.(float64))
	default:
		return 0
	}
}

func ToFloat32(i interface{}) float32 {
	switch i.(type) {
	case string:
		f32, e := ParseFloat(i.(string), 32)
		if e == nil {
			return 0
		}
		return float32(f32)
	case int, int32, int64:
		return float32(i.(int))
	case float32:
		return i.(float32)
	case float64:
		return float32(i.(float64))
	default:
		return 0
	}
}

func ToFloat64(i interface{}) float64 {
	switch i.(type) {
	case string:
		f64, e := ParseFloat(i.(string), 64)
		if e == nil {
			return 0
		}
		return f64
	case int, int32, int64:
		return float64(i.(int))
	case float32:
		return float64(i.(float32))
	case float64:
		return i.(float64)
	default:
		return 0
	}
}
*/

func MakeDate(layout string, value string) time.Time {
	t, e := time.Parse(layout, value)
	if e != nil {
		t, _ = time.Parse("2-Jan-2006", "1-Jan-1900")
		return t
	} else {
		return t
	}
}

func AddTime(dt0 time.Time, dt1 time.Time) time.Time {
	dtx := dt0
	return dtx.Add(dt1.Sub(MakeDate("03:04", "00:00")))
}

/*
func Value(i interface{}, fieldName string, def interface{}) interface{} {
	rv := reflect.ValueOf(i)
	var ret interface{}
	found := false
	if rv.Kind() == reflect.Map {
		mapkeys := rv.MapKeys()
		for i := 0; i < len(mapkeys) && !found; i++ {
			mapkey := mapkeys[i]
			mapkeyname := mapkey.String()
			if mapkeyname == fieldName {
				found = true
				mapvalue := rv.MapIndex(mapkey)
				if mapvalue.IsNil() {
					ret = def
				} else {
					ret = mapvalue.Interface()
				}
			}
		}
	} else if rv.Kind() == reflect.Struct {
		fv := rv.FieldByName(fieldName)
		if fv.IsValid() {
			found = true
			if (fv.Kind() == reflect.Struct || fv.Kind() == reflect.Map) && fv.IsNil() {
				ret = def
			} else {
				ret = fv.Interface()
			}
		}
	}

	if !found {
		return def
	} else {
		return ret
	}
}
*/

func Field(o interface{}, fieldName string) (reflect.Value, bool) {
	ref := reflect.ValueOf(o)
	if !ref.IsValid() {
		return ref, false
	}

	var es reflect.Value
	if ref.Kind() == reflect.Ptr {
		es = ref.Elem()
	} else {
		es = ref
	}
	fi := es.FieldByName(fieldName)
	if fi.IsValid() {
		return fi, true
	}
	return fi, false
}

func Jsonify(o interface{}) []byte {
	bs, e := json.Marshal(o)
	if e != nil {
		bs, _ = json.Marshal(struct{}{})
	}
	return bs
}

func JsonString(o interface{}) string {
	bs := Jsonify(o)
	return string(bs)
}

var _indentChar string

func SetIndentChar(c string) {
	_indentChar = c
}

func IndentChar() string {
	if _indentChar == "" {
		_indentChar = " "
	}
	return _indentChar
}

func JsonStringIndent(o interface{}, indentChar string) string {
	bs, e := json.MarshalIndent(o, "", IndentChar())
	if e != nil {
		return "{}"
	}
	return string(bs)
}

func Unjson(b []byte, result interface{}) error {
	buff := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buff)
	decoder.UseNumber()
	e := decoder.Decode(&result)
	return e
}

func UnjsonFromString(s string, result interface{}) error {
	b := bytes.NewBufferString(s)
	decoder := json.NewDecoder(b)
	decoder.UseNumber()
	e := decoder.Decode(&result)
	return e
}

func GetIP() ([]string, error) {

	ret := make([]string, 0)
	he := func(err error) ([]string, error) {
		return ret, err
	}

	ifaces, err := net.Interfaces()
	he(err)

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		he(err)

		for _, addr := range addrs {
			interfaceTxt := addr.String()
			if strings.HasSuffix(interfaceTxt, "24") {
				interfaceTxt = interfaceTxt[0 : len(interfaceTxt)-3]
				ret = append(ret, interfaceTxt)
			}
		}
	}
	if len(ret) == 0 {
		ret = append(ret, "127.0.0.1")
	}
	return ret, nil
}
