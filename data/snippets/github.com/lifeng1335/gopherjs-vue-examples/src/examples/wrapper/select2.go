package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-json"
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

var JQuery *js.Object = js.Global.Get("$")

type Select2Data struct {
	*js.Object
}

type Select2Props struct {
	*js.Object
	Options []*Option `js:"options"`
	Value   int       `js:"value"`
}

type Select2 struct {
	Data  *Select2Data
	Props *Select2Props
}

func NewSelect2(vm *vue.ViewModel) *Select2 {
	return &Select2{
		Data: &Select2Data{
			Object: vm.Data,
		},
		Props: &Select2Props{
			Object: vueutil.PropsData(vm),
		},
	}
}

func InitSelect2() {
	o := vue.NewOption()
	o.Name = "select2"
	o.Template = vueutil.GetTemplateById("select2-template")
	o.AddProp("options", "value")
	o.Data = vueutil.EmptyDataFunc

	o = vueutil.AddWatch(o, "value", func(vm *vue.ViewModel, newVal *js.Object, oldVal *js.Object) {
		println("update value", json.Stringify(newVal), json.Stringify(oldVal))
		// $(this.$el).val(value).trigger('change');
		sel := NewSelect2(vm)
		JQuery.New(vm.Get("$el")).Call("val", sel.Props.Value).Call("trigger", "change")
	})
	o = vueutil.AddWatch(o, "options", func(vm *vue.ViewModel, newVal *js.Object, oldVal *js.Object) {
		println("update options", json.Stringify(newVal), json.Stringify(oldVal))
		// $(this.$el).select2({ data: options })
		sel := NewSelect2(vm)
		arg := js.Global.Get("Object").New()
		arg.Set("data", sel.Props.Options)
		JQuery.New(vm.Get("$el")).Call("select2", arg)
	})
	o = vueutil.AddMounted(o, func(vm *vue.ViewModel) {
		println("Mounted Call", vm.Get("$el"))
		sel := NewSelect2(vm)
		println("on Mounted", sel.Props.Value)
		//println(json.Stringify(sel.Props.Object))

		objA := JQuery.New(vm.Get("$el"))
		//println(objA)

		argA := js.Global.Get("Object").New()
		argA.Set("data", sel.Props.Options)
		objB := objA.Call("select2", argA)
		objC := objB.Call("val", sel.Props.Value)

		//objC.Call("trigger", "change")
		objD := objC.Call("trigger", "change")
		objD.Call("on", "change", func() {
			// vm.$emit('input', this.value)
			// Todo: Notice "this" refer to JQuery Object
			println("on trigger change", JQuery.New(vm.Get("$el")).Index(0).Get("value"))
			value := JQuery.New(vm.Get("$el")).Index(0).Get("value").Int()
			vm.Emit("input", value)
		})
	})
	o = vueutil.AddDestroyed(o, func(vm *vue.ViewModel) {
		println("Destroyed Call", vm.Get("$el"))
		// $(this.$el).off().select2('destroy')
		JQuery.New(vm.Get("$el")).Call("off").Call("select2", "destroy")
	})

	o.NewComponent().Register("select2")
}

type Option struct {
	*js.Object
	Id   int    `js:"id"`
	Text string `js:"text"`
}

func NewOption(id int, text string) *Option {
	opt := &Option{
		Object: js.Global.Get("Object").New(),
	}
	opt.Id = id
	opt.Text = text
	return opt
}
