package router

import (
	"node/model"
	"node/protocol"
)

const (
	navigateDefaultHashValue = 64
)

type (
	// FuncActionCall is used as description a type of router
	FuncActionCall func([]byte) []byte
	// Router is used as
	Router struct {
		navigate []map[int]FuncActionCall
	}
)

// CreateRouter is used as create router
func CreateRouter() *Router {
	pRouter := &Router{
		navigate: make([]map[int]FuncActionCall, navigateDefaultHashValue),
	}
	for i := 0; i < navigateDefaultHashValue; i++ {
		pRouter.navigate[i] = make(map[int]FuncActionCall)
	}

	pRouter.initRouter()

	return pRouter
}
func (r *Router) initRouter() {
	r.addAction(protocol.GUnknownProtoCode, model.UnknownProtoCode)
}

// addAction is used as add action for router
func (r *Router) addAction(cmd int, call FuncActionCall) {
	r.navigate[cmd%navigateDefaultHashValue][cmd] = call
}

// getAction is used as get  a call function which is belong to cmd
func (r *Router) getAction(cmd int) (callback FuncActionCall, found bool) {
	callback, found = r.navigate[cmd%navigateDefaultHashValue][cmd]
	return callback, found
}

// Execute is used as execute action by cmd
func (r *Router) Execute(cmd int, data []byte) (ret []byte) {
	if callback, found := r.navigate[cmd%navigateDefaultHashValue][cmd]; found {
		return callback(data)
	}
	r.navigate[protocol.GUnknownProtoCode%navigateDefaultHashValue][protocol.GUnknownProtoCode](data)
	return ret
}
