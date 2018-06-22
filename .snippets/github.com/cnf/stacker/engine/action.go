package engine

import (
	"strings"

	"github.com/golang/glog"
)

// Action interface. because even interfaces need some action once in a while.
type Action interface {
	Setup
	StartStop
	Add() error
	Remove() error
	RemoveAll() error
}

// NewAction is a function definition each Action must provide
type NewAction func() (a Action)

var actionList = make(map[string]NewAction)

// RegisterAction is called to register an Action object
func RegisterAction(name string, creator NewAction) {
	name = strings.ToLower(name)
	if glog.V(1) {
		glog.Infof("registering Action: %s", name)
	}

	if actionList[name] != nil {
		glog.Fatalf("an Action named `%s` already exists", name)
	}
	actionList[name] = creator
}
