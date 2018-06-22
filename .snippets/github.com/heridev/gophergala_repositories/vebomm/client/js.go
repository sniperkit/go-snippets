package client

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/phaikawl/wade/app"
)

var (
	jsUndefined    = js.Undefined
	jsAlert        AlertJs
	jsLocalStorage js.Object
	isClientSide   = app.ClientSide
)

func init() {
	if isClientSide {
		jsAlert = AlertJs{js.Global.Get("alertify")}
		jsLocalStorage = js.Global.Get("localStorage")
	}
}

func jsGetMic(successFn, errFn func(js.Object)) {
	if isClientSide {
		js.Global.Call("veboGetMic", successFn, errFn)
	}
}

func jsSocket(id string, cons string, stream js.Object, statFn func(remoteObject js.Object), dcFn func()) {
	if isClientSide {
		js.Global.Call("socketConn", id, cons, stream, statFn, dcFn)
	}
}

type AlertJs struct{ obj js.Object }

func (ajs AlertJs) Error(msg string) {
	if isClientSide {
		ajs.obj.Call("error", msg)
	}
}

func (ajs AlertJs) Success(msg string) {
	if isClientSide {
		ajs.obj.Call("success", msg)
	}
}

func (ajs AlertJs) Log(msg string) {
	if isClientSide {
		ajs.obj.Call("log", msg)
	}
}
