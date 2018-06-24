package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/hiro4bbh/sticker"
)

// TemplateStandardFunctions is the function list used by templates.
var TemplateStandardFunctions = template.FuncMap{
	// interval returns the int slice between start and end.
	"interval": func(start, end int) []int {
		a := make([]int, 0, end-start)
		for i := start; i < end; i++ {
			a = append(a, i)
		}
		return a
	},
	// map returns a any-typed map with string key.
	// For example, (map "key1" value1 "key2" value2 ...).
	//
	// This function returns an error unless this cannot create a legal map.
	"map": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("arguments of map must be (key value)*")
		}
		m := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values)/2; i++ {
			key, ok := values[2*i+0].(string)
			if !ok {
				return nil, fmt.Errorf("#%d key must be string", i)
			}
			m[key] = values[2*i+1]
		}
		return m, nil
	},
	// mapSet sets the specified key and values to the map, and returns the empty string.
	// For example, (mapSet m "key1" value1 "key2" value2 ...).
	//
	// This function returns an error unless this cannot set correctly.
	"mapSet": func(m map[string]interface{}, values ...interface{}) (string, error) {
		if len(values)%2 != 0 {
			return "", fmt.Errorf("arguments of mapSet must be m (key value)*")
		}
		for i := 0; i < len(values)/2; i++ {
			key, ok := values[2*i+0].(string)
			if !ok {
				return "", fmt.Errorf("#%d key must be string", i)
			}
			m[key] = values[2*i+1]
		}
		return "", nil
	},
	// slice returns a any-typed slice.
	"slice": func(values ...interface{}) []interface{} {
		return values
	},

	// intToUint64 returns uint64 converted from int.
	"intToUint64": func(v int) uint64 {
		return uint64(v)
	},
	// intToFloat32 returns float32 converted from int.
	"intToFloat32": func(v int) float32 {
		return float32(v)
	},

	// decInt returns the previous int of the given int.
	"decInt": func(i int) int {
		return i - 1
	},
	// divInt returns the divided int of the given int slice.
	"divInt": func(values ...int) int {
		if len(values) == 0 {
			return 0
		}
		z := values[0]
		for i := 1; i < len(values); i++ {
			z /= values[i]
		}
		return z
	},
	// divFloat32 returns the divided float32 of the given float32 slice.
	"divFloat32": func(values ...float32) float32 {
		if len(values) == 0 {
			return 0
		}
		z := values[0]
		for i := 1; i < len(values); i++ {
			z /= values[i]
		}
		return z
	},
	// incInt returns the next integer of the given int.
	"incInt": func(i int) int {
		return i + 1
	},
	// mulInt returns the product integer of the given int slice.
	"mulInt": func(values ...int) int {
		z := 1
		for _, value := range values {
			z *= value
		}
		return z
	},
	// mulFloat32 returns the product integer of the given float32 slice.
	"mulFloat32": func(values ...float32) float32 {
		z := float32(1.0)
		for _, value := range values {
			z *= value
		}
		return z
	},
	// subInt returns the subtracted integer of the given int slice.
	"subInt": func(values ...int) int {
		if len(values) == 0 {
			return 0
		}
		z := values[0]
		for i := 1; i < len(values); i++ {
			z -= values[i]
		}
		return z
	},

	// concat returns the concatenated string.
	"concat": func(strs ...string) string {
		str := ""
		for _, s := range strs {
			str += s
		}
		return str
	},
	// toUpperFirst returns the string whose first character changes into upper case.
	"toUpperFirst": func(str string) string {
		return strings.ToUpper(str[0:1]) + str[1:]
	},

	// sortStrings returns a sorted string slice in ascending/decreasing order.
	"sortStrings": func(a []string, inDesc bool) []string {
		aSorted := make([]string, len(a))
		copy(aSorted, a)
		if inDesc {
			sort.Sort(sort.Reverse(sort.StringSlice(aSorted)))
		} else {
			sort.Strings(aSorted)
		}
		return aSorted
	},
	// sortSparseVector returns a sorted uint32 key in ascending/decreasing order by key/value.
	"sortSparseVector": func(m sticker.SparseVector, byKey bool, inDesc bool) []uint32 {
		kvs := make([]sticker.KeyValue32, 0, len(m))
		for key, value := range m {
			kvs = append(kvs, sticker.KeyValue32{key, value})
		}
		if byKey {
			if inDesc {
				sort.Sort(sort.Reverse(sticker.KeyValues32OrderedByKey(kvs)))
			} else {
				sort.Sort(sticker.KeyValues32OrderedByKey(kvs))
			}
		} else {
			if inDesc {
				sort.Sort(sort.Reverse(sticker.KeyValues32OrderedByValue(kvs)))
			} else {
				sort.Sort(sticker.KeyValues32OrderedByValue(kvs))
			}
		}
		keys := make([]uint32, 0, len(kvs))
		for _, kv := range kvs {
			keys = append(keys, kv.Key)
		}
		return keys
	},

	// cutUint32Slice returns a cutted sub-slice.
	// If start < 0, then it is treated as len(s) - start.
	"cutUint32Slice": func(s []uint32, start int, n int) []uint32 {
		if start >= 0 {
			if start < 0 {
				start = 0
			}
			if start+n >= len(s) {
				return s[start:]
			}
			return s[start : start+n]
		}
		if len(s)+start < 0 {
			start = -len(s)
		}
		if len(s)+start+n >= len(s) {
			return s[len(s)+start:]
		}
		return s[len(s)+start : len(s)+start+n]
	},

	// annotateUint32 returns a string annotated with m.
	"annotateUint32": func(m []string, v uint32) string {
		if v < uint32(len(m)) {
			return m[v]
		}
		return fmt.Sprintf("%d", v)
	},
	// annotateInts returns a string slice whose value is annotated with m.
	"annotateInts": func(m []string, a []int) []string {
		as := make([]string, len(a))
		for i, v := range a {
			var av string
			if v < len(m) {
				av = m[v]
			} else {
				av = fmt.Sprintf("%d", v)
			}
			as[i] = av
		}
		return as
	},

	// toJSON returns JSON string encoding the given object.
	//
	// This function returns an error in encoding.
	"toJSON": func(o interface{}) (string, error) {
		s, err := json.Marshal(o)
		if err != nil {
			return "", err
		}
		return string(s), nil
	},
}
