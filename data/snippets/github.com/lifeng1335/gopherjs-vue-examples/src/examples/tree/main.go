package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

type Node struct {
	*js.Object
	Name     string  `js:"name"`
	Children []*Node `js:"children"`
}

func NewNode(Name string, Children []*Node) *Node {
	d := &Node{
		Object: js.Global.Get("Object").New(),
	}
	d.Name = Name
	d.Children = Children
	return d
}

func NewAppData() *js.Object {
	obj := js.Global.Get("Object").New()
	m := NewNode("My Tree", []*Node{
		NewNode("hello", nil),
		NewNode("wat", nil),
		NewNode("child folder", []*Node{
			NewNode("child folder", []*Node{
				NewNode("hello", nil),
				NewNode("wat", nil),
			}),
			NewNode("hello", nil),
			NewNode("wat", nil),
			NewNode("child folder", []*Node{
				NewNode("hello", nil),
				NewNode("wat", nil),
			}),
		}),
	})
	obj.Set("treeData", m)
	return obj
}

func main() {
	InitTree()

	o := vue.NewOption()
	o.Data = NewAppData()

	v := o.NewViewModel()
	v.Mount("#demo")
}
