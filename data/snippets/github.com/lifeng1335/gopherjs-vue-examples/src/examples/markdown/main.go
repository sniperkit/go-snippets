package main

import (
	"vueutil"

	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

type AppData struct {
	*js.Object
	Input string `js:"input"`
}

type AppProps struct {
	*js.Object
}

type App struct {
	Data  *AppData
	Props *AppProps
}

func NewApp(vm *vue.ViewModel) *App {
	return &App{
		Data: &AppData{
			Object: vm.Data,
		},
		Props: &AppProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (app *App) SyncViewModel(vm *vue.ViewModel) {
	keys := js.Keys(app.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, app.Data.Get(v))
	}
	vm.Get("$options").Set("propsData", app.Props.Object)
}

func main() {
	o := vue.NewOption()
	o.Data = NewAppData()
	o.AddMethod("Update", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.Update(args[0])
		app.SyncViewModel(vm)
	})
	o.AddComputed("compiledMarkdown", func(vm *vue.ViewModel) interface{} {
		app := NewApp(vm)
		return app.CompiledMarkdown()
	})

	v := o.NewViewModel()
	v.Mount("#editor")
}

func NewAppData() interface{} {
	ad := &AppData{
		Object: js.Global.Get("Object").New(),
	}
	ad.Input = `# hello`
	return ad
}

func (app *App) Update(event *js.Object) {
	//lodash := js.Global.Get("_")
	//lodash.Call("debounce", func() {
	//	m.Input = event.Get("target").Get("value").String()
	//}, 300)
	app.Data.Input = event.Get("target").Get("value").String()
}

func (app *App) CompiledMarkdown() *js.Object {
	return js.Global.Call("marked", app.Data.Input, `{ sanitize: true }`)
}
