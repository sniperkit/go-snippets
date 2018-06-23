package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

	"gopkg.in/bblfsh/sdk.v2/protocol"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
	"gopkg.in/bblfsh/sdk.v2/uast"
)

func ParseString(code string) (*ast.File, error) {
	fs := token.NewFileSet()
	tree, err := parser.ParseFile(fs, "input.go", code, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func Parse(code string) (uast.Node, error) {
	f, err := ParseString(code)
	if err != nil {
		return nil, err
	}
	return convValue(reflect.ValueOf(f))
}

var (
	scopeType   = reflect.TypeOf((*ast.Scope)(nil))
	objectType  = reflect.TypeOf((*ast.Object)(nil))
	nodeType    = reflect.TypeOf((*ast.Node)(nil)).Elem()
	posType     = reflect.TypeOf(token.Pos(0))
	uastIntType = reflect.TypeOf(uast.Int(0))
)

func convPos(p token.Pos) uast.Node {
	return uast.Position{Offset: uint32(p)}.ToObject()
}

// convValue takes an AST node/value and converts it to a tree of uast types
// like Object and List. In this case we have a full control of json encoding
// and can annotate the tree with native AST type names.
func convValue(v reflect.Value) (uast.Node, error) {
	if !v.IsValid() {
		return nil, nil
	}
	t := v.Type()
	switch t {
	case posType:
		return convPos(v.Interface().(token.Pos)), nil
	}
	switch t.Kind() {
	case reflect.Slice:
		if v.Len() == 0 {
			return nil, nil
		}
		arr := make(uast.Array, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			el, err := convValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			arr = append(arr, el)
		}
		return arr, nil
	case reflect.Struct:
		m := make(uast.Object, t.NumField())
		m[uast.KeyType] = uast.String(t.Name()) // annotate nodes with type names
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Type == scopeType || f.Type == objectType {
				// do not follow scope and object pointers - need a graph structure for it
				continue
			}
			el, err := convValue(v.Field(i))
			if err != nil {
				return nil, err
			}
			m[f.Name] = el
		}
		return m, nil
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil, nil
		}
		o, err := convValue(v.Elem())
		if err != nil {
			return nil, err
		}
		if m, ok := o.(uast.Object); ok && v.Type().Implements(nodeType) {
			n := v.Interface().(ast.Node)
			m[uast.KeyStart] = convPos(n.Pos())
			m[uast.KeyEnd] = convPos(n.End())
		}
		return o, nil
	}
	o := v.Interface()
	if s, ok := o.(interface {
		String() string
	}); ok {
		return uast.String(s.String()), nil
	} else if t.ConvertibleTo(uastIntType) {
		return v.Convert(uastIntType).Interface().(uast.Int), nil
	}
	return uast.ToNode(o)
}

func NewDriver() *Driver {
	return &Driver{}
}

type Driver struct{}

func (Driver) Start() error {
	return nil
}
func (Driver) Close() error {
	return nil
}
func (Driver) Parse(req *driver.InternalParseRequest) (*driver.InternalParseResponse, error) {
	code, err := req.Encoding.Decode(req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contents: %v", err)
	}
	n, err := Parse(code)
	if err != nil {
		return nil, err
	}
	return &driver.InternalParseResponse{
		Status: driver.Status(protocol.Ok),
		AST:    n,
	}, nil
}
