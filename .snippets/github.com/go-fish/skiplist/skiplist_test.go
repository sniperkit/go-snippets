package skiplist

import (
	"testing"
	"unsafe"
)

var sl = NewSkiplist(32)

func TestPut(t *testing.T) {
	key := []byte("test")
	value := 111
	old, err := sl.Put(key, unsafe.Pointer(&value))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if old != nil {
		t.Fatalf("wanted nil, get %v\n", old)
	}

	val, err := sl.Get(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if val != unsafe.Pointer(&value) {
		t.Fatalf("wanted %v, get %v\n", value, val)
	}
}

func TestPutOnlyIfAbsent(t *testing.T) {
	key := []byte("test")
	value := 222
	old, err := sl.PutOnlyIfAbsent(key, unsafe.Pointer(&value))
	if err != ErrKeyExists {
		t.Fatalf("wanted %v, get %v\n", ErrKeyExists, err)
	}

	if old != nil {
		t.Fatalf("wanted nil, get %d\n", *(*int)(old))
	}

	key = []byte("test1")
	old, err = sl.PutOnlyIfAbsent(key, unsafe.Pointer(&value))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if old != nil {
		t.Fatalf("wanted nil, get %v\n", old)
	}

	val, err := sl.Get(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if val != unsafe.Pointer(&value) {
		t.Fatalf("wanted %v, get %v\n", value, val)
	}
}

func TestRemove(t *testing.T) {
	key := []byte("test")
	old, err := sl.Remove(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if *(*int)(old) != 111 {
		t.Fatalf("wanted 111, get %d\n", *(*int)(old))
	}

	val, err := sl.Get(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if val != nil {
		t.Fatalf("wanted nil, get %v\n", val)
	}

	ok, err := sl.Contains(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if ok {
		t.Fatalf("wanted false, get %v\n", ok)
	}

	val, _ = sl.Get([]byte("test1"))

	old, err = sl.CompareAndRemove([]byte("test1"), val)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if *(*int)(old) != 222 {
		t.Fatalf("wanted 222, get %d\n", *(*int)(old))
	}

	ok, err = sl.Contains([]byte("test1"))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if ok {
		t.Fatalf("wanted false, get %v\n", ok)
	}
}

func TestUpdate(t *testing.T) {
	key := []byte("test")
	action := func(oldValue unsafe.Pointer) unsafe.Pointer {
		old := *(*int)(oldValue)
		old++

		return unsafe.Pointer(&old)
	}

	old, err := sl.Update(key, action)
	if err != ErrNilValue {
		t.Fatalf("wanted %v, get %v", ErrNilValue, err)
	}

	if old != nil {
		t.Fatalf("wanted nil, get %v", old)
	}

	key = []byte("test1")
	value := 222
	old, err = sl.PutOnlyIfAbsent(key, unsafe.Pointer(&value))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if old != nil {
		t.Fatalf("wanted nil, get %d\n", *(*int)(old))
	}

	key = []byte("test1")
	old, err = sl.Update(key, action)
	if err != nil {
		t.Fatalf("wanted nil, get %v", old)
	}

	if *(*int)(old) != 222 {
		t.Fatalf("wanted 223, get %d", *(*int)(old))
	}

	val, err := sl.Get(key)
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if *(*int)(val) != 223 {
		t.Fatalf("wanted 223, get %d\n", *(*int)(val))
	}
}

func TestIterator(t *testing.T) {
	v1, err := sl.Get([]byte("test1"))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if *(*int)(v1) != 223 {
		t.Fatalf("wanted 223, get %d\n", *(*int)(v1))
	}

	val := 123
	sl.Put([]byte("test2"), unsafe.Pointer(&val))
	sl.Put([]byte("test3"), unsafe.Pointer(&val))
	sl.Put([]byte("test4"), unsafe.Pointer(&val))
	sl.Put([]byte("test5"), unsafe.Pointer(&val))
	sl.Put([]byte("test6"), unsafe.Pointer(&val))
	sl.Put([]byte("test7"), unsafe.Pointer(&val))
	sl.Put([]byte("test8"), unsafe.Pointer(&val))

	i, _ := NewIterator(sl, nil, nil)

	for i.Next() {
		k, v := i.NextNode()

		t.Log(k, v)

		if string(k) == "test5" {
			i.Remove()
		}
	}

	v, err := sl.Get([]byte("test5"))
	if err != nil {
		t.Fatalf("wanted nil, get %v\n", err)
	}

	if v != nil {
		t.Fatalf("wanted nil, get %v\n", v)
	}
}
