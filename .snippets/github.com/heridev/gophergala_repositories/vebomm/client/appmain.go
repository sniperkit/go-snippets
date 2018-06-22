package client

import (
	"encoding/json"
	"time"
	"fmt"
	"github.com/gophergala/vebomm/core"

	"github.com/phaikawl/wade/utils"
	"github.com/gopherjs/gopherjs/js"
	"github.com/phaikawl/wade/app"
	"github.com/phaikawl/wade/page"
)

const (
	HomePage     = "pg-home"
	RegisterPage = "pg-register"
	LoginPage    = "pg-login"
	RatePage = "pg-rate"
)

func update() {
	app.App().Render()
}

type AppMain struct {
	*app.Application
	User *core.User
}

func (am AppMain) Main(app *app.Application) {
	app.PageMgr.SetTemplate(Tmpl_main)
	am.Application = app

	r := app.Router()
	r.Handle("/", page.Page{
		Id:         HomePage,
		Title:      "Home",
		Controller: am.HomeHandler,
	})

	r.Handle("/register", page.Page{
		Id:         RegisterPage,
		Title:      "Register",
		Controller: am.RegisterHandler,
	})

	r.Handle("/login", page.Page{
		Id:         LoginPage,
		Title:      "Login",
		Controller: am.LoginHandler,
	})
	
	r.Handle("/rate/:id", page.Page{
		Id:         RatePage,
		Title:      "Rate",
		Controller: am.RateHandler,
	})

	page.GlobalDisplayScope.AddController(func(ctx page.Context) {
		if am.User == nil && isClientSide {
			am.User = am.checkAuth()
		}
	})
}

type Mode int

const (
	ModeIdle Mode = iota
	ModeGetMic
	ModeSearching
	ModeTalking
)

var _vma *RateVM
type RateVM struct {
	AppMain
	rateUserId string
	widget js.Object
}

func (vm *RateVM) Rate() {
	value := vm.widget.Call("getValue").Int()
	go func() {
		vm.Http.GET(fmt.Sprintf("/api/rate/%v/%v", vm.rateUserId, value))
		vm.PageMgr.GoToPage(HomePage)
	}()
}

func (am AppMain) RateHandler(ctx page.Context) {
	_vma = &RateVM{
		AppMain: am,
	}
	
	go func() {
		time.AfterFunc(2*time.Second, func() {
			if isClientSide {
				_vma.widget = js.Global.Get("Slider").New("#ex1")
			}
		})
	}()
	
	_vma.rateUserId, _ = ctx.NamedParams.Get("id")
}

var _vmh *HomeVM
type HomeVM struct {
	AppMain
	Cons core.Constraints
	Mode Mode
	Gender string
	stream js.Object
}

func (vm *HomeVM) Stop() {
	vm.stream.Call("close")
	loc := js.Global.Get("window").Get("location")
	loc.Call("replace", "http://"+loc.Get("hostname").String()+":"+loc.Get("port").String()+fmt.Sprintf("/web/rate/%v",vm.User.Id))
}

func (vm *HomeVM) Start() {
	vm.Cons.Gender = -1
	switch vm.Gender {
	case "M": vm.Cons.Gender = 1
	case "F": vm.Cons.Gender = 0
	}
	
	vm.Mode = ModeGetMic
	user := vm.AppMain.User
	cons, _ := json.Marshal(vm.Cons)
	jsGetMic(func(stream js.Object) {
		vm.Mode = ModeSearching
		jsSocket(utils.ToString(user.Id), string(cons), stream, func(remoteStream js.Object) {
			vm.Mode = ModeTalking
			println("talking")
			go update()
		}, func() {
			vm.PageMgr.GoToPage(RatePage, vm.User.Id)
		})
		go update()
	}, func(err js.Object) {
		
	})
}

func (am AppMain) HomeHandler(ctx page.Context) {
	if am.User == nil {
		ctx.GoToPage(RegisterPage)
		return
	}
	
	_vmh = &HomeVM{
		AppMain: am,
	}
	_vmh.Cons.MaxAge = 100
}

func (am AppMain) checkAuth() *core.User {
	idI := jsLocalStorage.Get("userid")
	if idI != nil {
		id := idI.Int64()
		resp, err := am.Http.POST("/api/login", core.User{Id: id})
		if err == nil && !resp.Failed() {
			var res core.LoginResult
			err = resp.ParseJSON(&res)
			if err == nil && res.Ok {
				return res.User
			}
		}
	}

	return nil
}

func (am AppMain) saveCred(user *core.User) {
	am.User = user
	if isClientSide {
		jsLocalStorage.Set("userid", user.Id)
	}
}

var _vml *LoginVM

type LoginVM struct {
	AppMain
	User   core.User
	Result core.LoginResult
}

func (vm *LoginVM) DoLogin() {
	go func() {
		resp, err := vm.Http.POST("/api/login", vm.User)
		if err != nil || resp.Failed() {
			jsAlert.Error("Connection problem or server error.")
			return
		}

		err = resp.ParseJSON(&vm.Result)
		if err != nil || !vm.Result.Ok {
			jsAlert.Error("Login failed.")
			return
		}

		jsAlert.Success("Login successful.")
		vm.AppMain.saveCred(vm.Result.User)
		vm.PageMgr.GoToPage(HomePage)
	}()
}

func (am AppMain) LoginHandler(ctx page.Context) {
	_vml = &LoginVM{
		AppMain: am,
	}
}
