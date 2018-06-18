package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type M map[string]interface{}
type Ms []M

var ErrorPathNotFound = errors.New("Path requested is not available")

func (m M) Set(k string, v interface{}) M {
	m[k] = v
	return m
}

func (m M) Get(k string, d ...interface{}) interface{} {
	if get, b := m[k]; b {
		if IsNilOrEmpty(get) && len(d) > 0 {
			get = d[0]
		}
		return get
	} else {
		if len(d) > 0 {
			return d[0]
		} else {
			return nil
		}
	}
}

func (m M) GetRef(k string, d, out interface{}) {
	/*
		defer func() {
			if r := recover(); r != nil {
			}
		}()
	*/
	valget := reflect.Indirect(reflect.ValueOf(m.Get(k, d)))
	valout := reflect.ValueOf(out)
	valout.Elem().Set(valget)
}

func ToM(v interface{}) (M, error) {
	buffer := []byte{}
	buff := bytes.NewBuffer(buffer)
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(buff)
	decoder.UseNumber()
	res := M{}
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m M) ToBytes(encodertype string, others ...interface{}) []byte {
	encodertype = strings.ToLower(encodertype)
	if encodertype == "" {
		encodertype = "json"
	}
	/*
		if encodertype == "json" {
			bs, e := json.Marshal(m)
			if e != nil {
				return []byte{}
			}
			return bs
		}
	*/
	return ToBytes(m, encodertype)
	//return []byte{}
}

func (m *M) Cast(k string, d interface{}) error {
	var e error
	if m.Has(k) == false {
		return fmt.Errorf("No data for key %s", k)
	}
	b, e := json.Marshal(m.Get(k, nil))
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, d)
	return e
}

func (m M) GetFloat64(k string) float64 {
	i := m.Get(k, 0)
	return ToFloat64(i, 6, RoundingAuto)
}

func (m M) GetString(k string) string {
	s := m.Get(k, "")
	return ToString(s)
}

func (m M) GetInt(k string) int {
	i := m.Get(k, 0)
	return ToInt(i, RoundingAuto)
}

func (m *M) Unset(k string) {
	delete(*m, k)
}

func (m M) GetFloat32(k string) float32 {
	i := m.Get(k, 0)
	return ToFloat32(i, 4, RoundingAuto)
}

func (m M) Has(k string) bool {
	_, has := m[k]
	return has
}

func (m M) PathGet(path string) (interface{}, error) {
	pathlist := strings.Split(path, ".")
	var curobj interface{} = m
	for _, nextpath := range pathlist {
		switch curobj.(type) {
		case M:
			curM := curobj.(M)
			var found bool
			curobj, found = curM[nextpath]
			if !found {
				return nil, ErrorPathNotFound
			}

		case map[string]interface{}:
			curM := curobj.(map[string]interface{})
			var found bool
			curobj, found = curM[nextpath]
			if !found {
				return nil, ErrorPathNotFound
			}

		default:
			return nil, ErrorPathNotFound
		}
	}

	return curobj, nil
}

func (m M) Keys() []string {
	var ret []string
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

func (m M) Values() []interface{} {
	var ret []interface{}
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func CopyM(from, to *M,
	copyFieldIfNotExist bool,
	exceptFields []string) {
	//Printf("Copy from: %s to %s\n", JsonString(from), JsonString(to))
	var exceptFieldsIface []interface{}
	Serde(exceptFields, &exceptFieldsIface, "")
	fromm := *from
	tom := *to
	for f, fv := range fromm {
		if !HasMember(exceptFieldsIface, f) {
			if tom.Has(f) {
				tom.Set(f, fv)
			} else if copyFieldIfNotExist {
				tom.Set(f, fv)
			}
		}
	}
	*to = tom
}
