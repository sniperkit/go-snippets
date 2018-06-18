package merge

import (
	"reflect"
	"testing"
)

type Foo struct {
	A string
	B *string
	C *int
	D []byte
	e int
}

func EqualPublicFields(l, r *Foo) bool {
	// all tests l and r have different e
	if l.e == r.e {
		return false
	}
	if l.A != r.A || l.B != r.B || l.C != r.C || !reflect.DeepEqual(l.D, r.D) {
		return false
	}
	return true
}

func TestMerge1(t *testing.T) {
	b := "world"
	c := 123
	f := &Foo{
		A: "hello",
		B: &b,
		C: &c,
		D: []byte("test"),
		e: 0,
	}

	nb := "WORLD"
	nc := 321
	new := &Foo{
		A: "HELLO",
		B: &nb,
		C: &nc,
		D: []byte("pass"),
		e: 1,
	}
	if err := MergeStruct(f, new); err != nil {
		t.Fatal(err)
	}
	if !EqualPublicFields(f, new) {
		t.Errorf("source and destination should be equal")
	}
}

func TestMerge2(t *testing.T) {
	f := &Foo{
		e: 0,
	}

	nb := "WORLD"
	nc := 321
	new := &Foo{
		A: "HELLO",
		B: &nb,
		C: &nc,
		D: []byte("pass"),
		e: 1,
	}
	if err := MergeStruct(f, new); err != nil {
		t.Fatal(err)
	}
	if !EqualPublicFields(f, new) {
		t.Errorf("source and destination should be equal")
	}
}

func TestMerge3(t *testing.T) {
	f := &Foo{
		e: 0,
	}

	nc := 321
	new := &Foo{
		C: &nc,
		D: []byte("pass"),
		e: 1,
	}
	if err := MergeStruct(f, new); err != nil {
		t.Fatal(err)
	}
	if !EqualPublicFields(f, new) {
		t.Errorf("source and destination should be equal")
	}
}
func TestMergeWrongTypes(t *testing.T) {
	b := 1
	if err := MergeStruct(&b, &Foo{}); err != ErrorMergeSameStruct {
		t.Errorf("different merge types should not do it")
	}

	if err := MergeStruct(1, 2); err != ErrorMergePointers {
		t.Errorf("both arguments should be pointers")
	}
}
