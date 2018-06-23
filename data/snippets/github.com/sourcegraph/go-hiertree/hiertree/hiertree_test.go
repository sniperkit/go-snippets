package hiertree

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var only = flag.Int("test.only", -1, "only run the TestList test row at this index")

type elem string

func (e elem) PathComponents() []string {
	return strings.Split(string(e), "/")
}

func mkElems(in []elem) (out []Elem) {
	out = make([]Elem, len(in))
	for i, e := range in {
		out[i] = Elem(e)
	}
	return
}

func TestList(t *testing.T) {
	tests := []struct {
		elems []Elem   // input
		paths []string // output
		error string   // substring of error, if error is expected
	}{
		{
			elems: mkElems([]elem{}),
			paths: []string{},
		},
		{
			elems: mkElems([]elem{"bar", "foo"}),
			paths: []string{"bar*", "foo*"},
		},
		{
			elems: mkElems([]elem{"foo", "foo/bar"}),
			paths: []string{"foo*>", "[foo/]bar*"},
		},
		{
			elems: mkElems([]elem{"foo/bar", "foo/baz"}),
			paths: []string{"foo>", "[foo/]bar*", "[foo/]baz*"},
		},
		{
			elems: mkElems([]elem{"foo", "foo/bar", "foo/baz"}),
			paths: []string{"foo*>", "[foo/]bar*", "[foo/]baz*"},
		},
		{
			elems: mkElems([]elem{"foo/bar/baz"}),
			paths: []string{"foo>", "[foo/]bar>", "[foo/bar/]baz*"},
		},
		{
			elems: mkElems([]elem{"foo/bar", "foo/bar/baz"}),
			paths: []string{"foo>", "[foo/]bar*>", "[foo/bar/]baz*"},
		},
		{
			elems: mkElems([]elem{"foo/bar/baz", "foo/bar/qux"}),
			paths: []string{"foo>", "[foo/]bar>", "[foo/bar/]baz*", "[foo/bar/]qux*"},
		},
		{
			elems: mkElems([]elem{"foo/bar/baz/bud", "foo/bar/qux/qup"}),
			paths: []string{"foo>", "[foo/]bar>", "[foo/bar/]baz>", "[foo/bar/baz/]bud*", "[foo/bar/]qux>", "[foo/bar/qux/]qup*"},
		},
		{
			elems: mkElems([]elem{"baz/bud", "baz/bud/eee/fff", "qux/qup", "qux/qup/eee/fff"}),
			paths: []string{"baz>", "[baz/]bud*>", "[baz/bud/]eee>", "[baz/bud/eee/]fff*", "qux>", "[qux/]qup*>", "[qux/qup/]eee>", "[qux/qup/eee/]fff*"},
		},
		{
			elems: mkElems([]elem{"foo/bar/baz/qux", "foo/bar"}),
			paths: []string{"foo>", "[foo/]bar*>", "[foo/bar/]baz>", "[foo/bar/baz/]qux*"},
		},
		{
			elems: mkElems([]elem{"foo/bar", "baz/qux"}),
			paths: []string{"baz>", "[baz/]qux*", "foo>", "[foo/]bar*"},
		},
		{
			elems: mkElems([]elem{"foo/bar", "baz/qux", "foo/baz"}),
			paths: []string{"baz>", "[baz/]qux*", "foo>", "[foo/]bar*", "[foo/]baz*"},
		},

		// errors
		{
			elems: mkElems([]elem{"bar", "bar"}),
			error: "duplicate",
		},
	}

	for i, test := range tests {
		if *only != -1 && *only != i {
			continue
		}
		entries, err := List(test.elems)
		label := fmt.Sprintf("elems#%d %v:", i, test.elems)
		if test.error == "" {
			if err != nil {
				t.Errorf("%s: unexpected error %q", label, err)
			} else if paths := Inspect(entries); !reflect.DeepEqual(test.paths, paths) {
				t.Errorf("%s:\nwant %v\n got %v", label, test.paths, paths)
			}
		} else {
			if err == nil {
				t.Errorf("%s: want error containing %q, got no error", label, test.error)
			} else if !strings.Contains(err.Error(), test.error) {
				t.Errorf("%s: want error containing %q, got %q", label, test.error, err)
			}
		}
	}
}
