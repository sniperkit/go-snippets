package redis

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kdar/cachei"
	"reflect"
	"testing"

	// "github.com/vmihailenco/msgpack"
)

var float64_ float64
var string_ string

var cacheTests = []struct {
	value interface{}
	into  interface{}
}{
	{map[string]interface{}{
		"hey":   "there",
		"float": float64(4),
		//"bytes": []byte("caca"),
	}, &map[string]interface{}{}},
	{float64(6), &float64_},
	{"hey", &string_},
}

// func TestCaca(t *testing.T) {
//   in := map[string]interface{}{"foo": 1, "hello": "world"}
//   b, err := msgpack.Marshal(in)
//   _ = err

//   var out interface{}
//   err = msgpack.Unmarshal(b, &out)
//   fmt.Printf("%v %#v\n", err, out)
// }

func TestInto(t *testing.T) {
	cache, err := cachei.Open("redis", cachei.DataSource{})
	if err != nil {
		t.Fatal(err)

	}

	for i, tt := range cacheTests {
		cacheMiss := false
		missFunc := func() (interface{}, error) {
			cacheMiss = true
			return tt.value, nil
		}

		verr, cerr := cache.OutSetFn(fmt.Sprintf("__kdarcacheinto_redis_test_%d", i), 1, tt.into, missFunc)
		if verr != nil {
			t.Fatal(verr)
		}

		if cerr != nil {
			t.Fatal(cerr)
		}

		if !cacheMiss {
			t.Fatalf("%d-1. Expected a cache miss, but instead found cache", i)
		}

		if !reflect.DeepEqual(tt.value, reflect.ValueOf(tt.into).Elem().Interface()) {
			t.Fatalf("%d-1. Expected:\n%s,\ngot:\n%s", i, spew.Sprintf("%#v", &tt.value), spew.Sprintf("%#v", tt.into))
		}

		cacheMiss = false

		verr, cerr = cache.OutSetFn(fmt.Sprintf("__kdarcacheinto_redis_test_%d", i), 1, tt.into, missFunc)
		if verr != nil {
			t.Fatal(verr)
		}

		if cerr != nil {
			t.Fatal(cerr)
		}

		if cacheMiss {
			t.Fatalf("%d-2. Expected cache, but instead got a cache miss.", i)
		}

		if !reflect.DeepEqual(tt.value, reflect.ValueOf(tt.into).Elem().Interface()) {
			t.Fatalf("%d-2. Expected: %s, got: %s", i, spew.Sprintf("%#v", tt.value), spew.Sprintf("%#v", tt.into))
		}
	}
}

func TestGet(t *testing.T) {
	cache, err := cachei.Open("redis", cachei.DataSource{})
	if err != nil {
		t.Fatal(err)
	}

	for i, tt := range cacheTests {
		cacheMiss := false
		missFunc := func() (interface{}, error) {
			cacheMiss = true
			return tt.value, nil
		}

		v1, verr, cerr := cache.GetSetFn(fmt.Sprintf("__kdarcacheget_redis_test_%d", i), 1, missFunc)
		if verr != nil {
			t.Fatal(verr)
		}

		if cerr != nil {
			t.Fatal(cerr)
		}

		if !cacheMiss {
			t.Fatalf("%d-1. Expected a cache miss, but instead found cache", i)
		}

		if !reflect.DeepEqual(tt.value, v1) {
			t.Fatalf("%d-1. Expected: %s, got: %s", i, spew.Sprintf("%#v", tt.value), spew.Sprintf("%#v", v1))
		}

		cacheMiss = false

		v2, verr, cerr := cache.GetSetFn(fmt.Sprintf("__kdarcacheget_redis_test_%d", i), 1, missFunc)
		if verr != nil {
			t.Fatal(verr)
		}

		if cerr != nil {
			t.Fatal(cerr)
		}

		if cacheMiss {
			t.Fatalf("%d-2. Expected cache, but instead got a cache miss.", i)
		}

		if !reflect.DeepEqual(v1, v2) {
			t.Fatalf("%d-2. Expected: %s, got: %s", i, spew.Sprintf("%#v", v1), spew.Sprintf("%#v", v2))
		}
	}
}
