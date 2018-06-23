package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

type TreeData struct {
	*js.Object
	Open bool `js:"open"`
}

type TreeProps struct {
	*js.Object
	Model *Node `js:"model"`
}

type Tree struct {
	Data  *TreeData
	Props *TreeProps
}

func NewTree(vm *vue.ViewModel) *Tree {
	return &Tree{
		Data: &TreeData{
			Object: vm.Data,
		},
		Props: &TreeProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func NewTreeData() interface{} {
	g := &TreeData{
		Object: js.Global.Get("Object").New(),
	}
	g.Open = false
	return g
}

func (t *Tree) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	//vm.Data = t.Date.Object
	keys := js.Keys(t.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, t.Data.Get(v))
	}

	vm.Get("$options").Set("propsData", t.Props.Object)
}

func (t *Tree) IsFolder() bool {
	if len(t.Props.Model.Children) != 0 {
		return true
	}
	return false
}

func (t *Tree) AddChild() {
	t.Props.Model.Children = append(t.Props.Model.Children, NewNode("new stuff", nil))
}

func (t *Tree) ChangeType() {
	if !t.IsFolder() {
		t.Data.Open = true
		t.AddChild()
	}
}

func (t *Tree) Toggle() {
	if t.IsFolder() {
		t.Data.Open = !t.Data.Open
	}
}

func InitTree() {
	o := vue.NewOption()
	o.Name = "item"
	o.Template = vueutil.GetTemplateById("item-template")
	o.AddProp("model")
	o.Data = NewTreeData
	o.AddMethod("Toggle", func(vm *vue.ViewModel, args []*js.Object) {
		tree := NewTree(vm)
		tree.Toggle()
		tree.SyncViewModel(vm)
	})
	o.AddMethod("AddChild", func(vm *vue.ViewModel, args []*js.Object) {
		tree := NewTree(vm)
		tree.AddChild()
		tree.SyncViewModel(vm)
	})
	o.AddMethod("ChangeType", func(vm *vue.ViewModel, args []*js.Object) {
		tree := NewTree(vm)
		tree.ChangeType()
		tree.SyncViewModel(vm)
	})
	o.AddComputed("isFolder", func(vm *vue.ViewModel) interface{} {
		tree := NewTree(vm)
		return tree.IsFolder()
	})
	o.NewComponent().Register("item")
}
