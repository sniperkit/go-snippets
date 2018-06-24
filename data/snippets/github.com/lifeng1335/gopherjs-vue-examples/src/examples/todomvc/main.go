package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-json"
	"github.com/oskca/gopherjs-vue"
	"regexp"
	"strings"
	"vueutil"
)

var todoStorage *TodoStorage

type TodoStorage struct {
	*js.Object
	Uid        int
	StorageKey string
}

type Todos struct {
	*js.Object
	Todos []*Todo `js:"todos"`
}

func NewDefaultTodoStorage() *TodoStorage {
	return &TodoStorage{
		Object:     js.Global.Get("localStorage"),
		StorageKey: "todos-vuejs-2.0",
		Uid:        0,
	}
}

func (ts *TodoStorage) Fetch() *js.Object {
	var todos *js.Object
	obj := ts.Call("getItem", ts.StorageKey)
	if obj == nil {
		todos = json.Parse("[]")
		ts.Uid = 0
	} else {
		todos = json.Parse(obj.String())
		ts.Uid = todos.Length()
	}
	return todos
}

func (ts *TodoStorage) Save(todos *js.Object) {
	ts.Call("setItem", ts.StorageKey, json.Stringify(todos))
}

func NewAppData() interface{} {
	obj := js.Global.Get("Object").New()
	obj.Set("todos", todoStorage.Fetch())
	obj.Set("newTodo", "")
	obj.Set("editedTodo", nil)
	obj.Set("visibility", "all")
	return obj
}

type AppData struct {
	*js.Object
	Todos      []*Todo `js:"todos"`
	NewTodo    string  `js:"newTodo"`
	EditedTodo *Todo   `js:"editedTodo"`
	Visibility string  `js:"visibility"`
}

type AppProps struct {
	*js.Object
}

type App struct {
	Date  *AppData
	Props *AppProps
}

func NewApp(vm *vue.ViewModel) *App {
	return &App{
		Date: &AppData{
			Object: vm.Data,
		},
		Props: &AppProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (app *App) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	// vm.Data = t.Date.Object
	keys := js.Keys(app.Date.Object)
	for _, v := range keys {
		vm.Data.Set(v, app.Date.Get(v))
	}

	vm.Get("$options").Set("propsData", app.Props.Object)
}

// visibility filters
func (app *App) All() []*Todo {
	return app.Date.Todos
}

func (app *App) Active() []*Todo {
	ActiveTodos := []*Todo{}
	for _, todo := range app.Date.Todos {
		if todo.Completed == false {
			ActiveTodos = append(ActiveTodos, todo)
		}
	}
	return ActiveTodos
}

func (app *App) Completed() []*Todo {
	CompletedTodos := []*Todo{}
	for _, todo := range app.Date.Todos {
		if todo.Completed == true {
			CompletedTodos = append(CompletedTodos, todo)
		}
	}
	return CompletedTodos
}

func (app *App) FilteredTodos() []*Todo {
	switch app.Date.Visibility {
	case "active":
		return app.Active()
	case "completed":
		return app.Completed()
	default:
		return app.All()
	}
}

func (app *App) Remaining() int {
	return len(app.Active())
}

func (app *App) AllDoneGetter() bool {
	return app.Remaining() == 0
}

func (app *App) AllDoneSetter(val *js.Object) {
	for index, _ := range app.Date.Todos {
		app.Date.Todos[index].Completed = val.Bool()
	}
}

// Method
func (app *App) AddTodo() {
	title := strings.TrimSpace(app.Date.NewTodo)
	if title == "" {
		return
	}
	todoStorage.Uid = todoStorage.Uid + 1
	app.Date.Todos = append(app.Date.Todos, NewTodo(todoStorage.Uid, title, false))
	app.Date.NewTodo = ""
}

func (app *App) RemoveTodo(todo *js.Object) {
	to := &Todo{
		Object: todo,
	}
	todos := []*Todo{}
	for _, value := range app.Date.Todos {
		if value.Id != to.Id {
			todos = append(todos, value)
		}
	}
	app.Date.Todos = todos
}

func (app *App) EditTodo(vm *vue.ViewModel, todo *js.Object) {
	println("Method EditTodo")
	to := &Todo{
		Object: todo,
	}
	vm.Data.Set("beforeEditCache", to.Title)
	app.Date.EditedTodo = to
}

func (app *App) DoneEdit(todo *js.Object) {
	println("Method DoneEdit")
	to := &Todo{
		Object: todo,
	}
	// Todo: ?
	app.Date.EditedTodo = nil
	to.Title = strings.TrimSpace(to.Title)
	if to.Title == "" {
		println("Method DoneEdit => RemoveTodo")
		app.RemoveTodo(to.Object)
	}
}

func (app *App) CancelEdit(vm *vue.ViewModel, todo *js.Object) {
	println("Method CancelEdit")
	app.Date.EditedTodo = nil
	to := &Todo{
		Object: todo,
	}
	to.Title = vm.Data.Get("beforeEditCache").String()
}

func (app *App) RemoveCompleted() {
	app.Date.Todos = app.Active()
}

func WatchTodosObj() *js.Object {
	obj := js.Global.Get("Object").New()
	fn := func(todos *js.Object) {
		todoStorage.Save(todos)
	}
	obj.Set("handler", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		fn(arguments[0])
		return nil
	}))
	obj.Set("deep", true)

	objA := js.Global.Get("Object").New()
	objA.Set("todos", obj)
	return objA
}

func main() {
	InitFilters()
	todoStorage = NewDefaultTodoStorage()

	d := vue.NewDirective()
	d.SetUpdater(func(el *js.Object, ctx *vue.DirectiveBinding, val *js.Object) {
		//println("todo-focus", json.Stringify(ctx))
		if ctx.Value == "true" {
			el.Call("focus")
		}
	}).Register("todo-focus")

	o := vue.NewOption()
	o.Data = NewAppData()
	o = o.OnLifeCycleEvent(vue.EvtCreated, func(vm *vue.ViewModel) {
		println("OnLifeCycleEvent", "EvtCreated")
	})
	o.AddComputed("filteredTodos", func(vm *vue.ViewModel) interface{} {
		app := NewApp(vm)
		return app.FilteredTodos()
	})
	o.AddComputed("remaining", func(vm *vue.ViewModel) interface{} {
		app := NewApp(vm)
		return app.Remaining()
	})
	o.AddComputed("allDone", func(vm *vue.ViewModel) interface{} {
		// getter
		app := NewApp(vm)
		return app.AllDoneGetter()
	}, func(vm *vue.ViewModel, val *js.Object) {
		// setter
		app := NewApp(vm)
		app.AllDoneSetter(val)
		app.SyncViewModel(vm)
	})
	o.Mixin(js.M{
		"watch": WatchTodosObj(),
	})
	o.AddMethod("addTodo", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.AddTodo()
		app.SyncViewModel(vm)
	})
	o.AddMethod("removeTodo", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.RemoveTodo(args[0])
		app.SyncViewModel(vm)
	})
	o.AddMethod("editTodo", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.EditTodo(vm, args[0])
		app.SyncViewModel(vm)
	})
	o.AddMethod("doneEdit", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.DoneEdit(args[0])
		app.SyncViewModel(vm)
	})
	o.AddMethod("cancelEdit", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.CancelEdit(vm, args[0])
		app.SyncViewModel(vm)
	})
	o.AddMethod("removeCompleted", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.RemoveCompleted()
		app.SyncViewModel(vm)
	})

	v := o.NewViewModel()

	// window.addEventListener('hashchange', onHashChange)
	// onHashChange()
	js.Global.Call("addEventListener", "hashchange", func() {
		onHashChange(v)
	})
	onHashChange(v)

	v.Mount(".todoapp")
}

// handle routing
func onHashChange(vm *vue.ViewModel) {
	app := NewApp(vm)
	exp := `#\/?`
	reg := regexp.MustCompile(exp)
	visibility := reg.ReplaceAllString(js.Global.Get("location").Get("hash").String(), "")
	println("onHashChange", visibility)
	switch visibility {
	case "active":
		app.Date.Visibility = "active"
		js.Global.Get("location").Set("hash", "active")
	case "completed":
		app.Date.Visibility = "completed"
		js.Global.Get("location").Set("hash", "completed")
	case "all":
		app.Date.Visibility = "all"
		js.Global.Get("location").Set("hash", "all")
	default:
		app.Date.Visibility = "all"
		js.Global.Get("location").Set("hash", "")
	}
	app.SyncViewModel(vm)
}

func InitFilters() {
	vue.NewFilter(func(v *js.Object) interface{} {
		if v.Int() == 1 {
			return "item"
		} else {
			return "items"
		}
	}).Register("pluralize")
}

type Todo struct {
	*js.Object
	Id        int    `js:"id"`
	Title     string `js:"title"`
	Completed bool   `js:"completed"`
}

func NewTodo(id int, title string, completed bool) *Todo {
	todo := &Todo{
		Object: js.Global.Get("Object").New(),
	}
	todo.Id = id
	todo.Title = title
	todo.Completed = completed
	return todo
}
