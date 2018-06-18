package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"regexp"
	"strings"
	"vueutil"
)

const RegStr = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`

func NewAppData() interface{} {
	obj := js.Global.Get("Object").New()
	obj.Set("newUser", NewUser("", ""))
	return obj
}

type AppData struct {
	*js.Object
	NewUser *User `js:"newUser"`
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
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	// vm.Data = t.Date.Object
	keys := js.Keys(app.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, app.Data.Get(v))
	}

	vm.Get("$options").Set("propsData", app.Props.Object)
}

func (app *App) Validation() *js.Object {
	obj := js.Global.Get("Object").New()
	ValidName := false
	if strings.TrimSpace(app.Data.NewUser.Name) != "" {
		ValidName = true
	}
	reg := regexp.MustCompile(RegStr)
	ValidEmail := reg.Match([]byte(app.Data.NewUser.Email))

	obj.Set("name", ValidName)
	obj.Set("email", ValidEmail)
	return obj
}

func (app *App) IsValid() bool {
	obj := app.Validation()
	for _, value := range js.Keys(obj) {
		if obj.Get(value).Bool() == false {
			return false
		}
	}
	return true
}

func (app *App) AddUser(vm *vue.ViewModel) {
	if app.IsValid() {
		usersRef := vm.Get("$options").Get("firebase").Get("users")
		usersRef.Call("push", app.Data.NewUser)
		app.Data.NewUser.Name = ""
		app.Data.NewUser.Email = ""
	}
}

func (app *App) RemoveUser(vm *vue.ViewModel, user *js.Object) {
	// usersRef.child(user['.key']).remove()
	usersRef := vm.Get("$options").Get("firebase").Get("users")
	usersRef.Call("child", user.Get(".key")).Call("remove")
}

func main() {
	o := vue.NewOption()
	o.Data = NewAppData()
	o.Mixin(js.M{
		"firebase": NewFirebase(),
	})
	o = o.OnLifeCycleEvent(vue.EvtCreated, func(vm *vue.ViewModel) {
		println("OnLifeCycleEvent", "EvtCreated")
	})
	o.AddComputed("validation", func(vm *vue.ViewModel) interface{} {
		//println("Call validation")
		app := NewApp(vm)
		return app.Validation()
	})
	o.AddMethod("addUser", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.AddUser(vm)
		app.SyncViewModel(vm)
	})
	o.AddMethod("removeUser", func(vm *vue.ViewModel, args []*js.Object) {
		app := NewApp(vm)
		app.RemoveUser(vm, args[0])
		app.SyncViewModel(vm)
	})

	v := o.NewViewModel()
	v.Mount("#app")
}

type User struct {
	*js.Object
	Name  string `js:"name"`
	Email string `js:"email"`
}

func NewUser(Name string, Email string) *User {
	u := &User{
		Object: js.Global.Get("Object").New(),
	}
	u.Name = Name
	u.Email = Email
	return u
}

func InitFirebase() *js.Object {
	// Notice!!!: You always need to call Vue.use
	VueFire := js.Global.Get("VueFire")
	vue.Use(VueFire)

	config := js.Global.Get("Object").New()
	config.Set("apiKey", "AIzaSyAi_yuJciPXLFr_PYPeU3eTvtXf8jbJ8zw")
	config.Set("authDomain", "vue-demo-537e6.firebaseapp.com")
	config.Set("databaseURL", "https://vue-demo-537e6.firebaseio.com")

	// firebase.initializeApp(config)
	firebase := js.Global.Get("firebase")
	return firebase.Call("initializeApp", config)
}

func NewFirebase() *js.Object {
	InitFirebase()

	// var usersRef = firebase.database().ref('users')
	firebase := js.Global.Get("firebase")
	usersRef := firebase.Call("database").Call("ref", "users")

	fb := js.Global.Get("Object").New()
	fb.Set("users", usersRef)
	return fb
}
