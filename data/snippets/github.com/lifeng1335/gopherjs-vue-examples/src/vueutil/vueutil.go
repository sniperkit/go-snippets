package vueutil

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

func PropsData(vm *vue.ViewModel) *js.Object {
	obj := js.Global.Get("Object").New()
	keys := js.Keys(vm.Get("$options").Get("propsData"))
	for _, value := range keys {
		obj.Set(value, vm.Get(value))
	}
	return obj
}

func GetTemplateById(ElementId string) string {
	return js.Global.Get("document").Call("getElementById", ElementId).Get("innerText").String()
}

func EmptyDataFunc() interface{} {
	return js.Global.Get("Object").New()
}

func AddWatch(o *vue.Option, name string, fn func(vm *vue.ViewModel, newVal *js.Object, oldVal *js.Object)) *vue.Option {
	obj := js.Global.Get("Object").New()
	obj.Set("handler", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		vm := &vue.ViewModel{
			Object: this,
		}
		fn(vm, arguments[0], arguments[1])
		return nil
	}))
	obj.Set("deep", false)

	watch := js.Global.Get("Object").New()
	watch.Set(name, obj)
	return o.Mixin(
		js.M{
			"watch": watch,
		},
	)
}

func AddMounted(o *vue.Option, fn func(vm *vue.ViewModel)) *vue.Option {
	mounted := js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		vm := &vue.ViewModel{
			Object: this,
		}
		fn(vm)
		return nil
	})
	return o.Mixin(
		js.M{
			"mounted": mounted,
		},
	)
}

func AddDestroyed(o *vue.Option, fn func(vm *vue.ViewModel)) *vue.Option {
	destroyed := js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		vm := &vue.ViewModel{
			Object: this,
		}
		fn(vm)
		return nil
	})
	return o.Mixin(
		js.M{
			"destroyed": destroyed,
		},
	)
}
