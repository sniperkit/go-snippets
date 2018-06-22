package engine

import (
	"strings"

	"github.com/golang/glog"
)

// ConfStore is the interface for every config store module.
type ConfStore interface {
	Setup
	StartStop
}

// NewConfStore is a function definition each ConfStore must provide
type NewConfStore func() (c ConfStore)

var confStoreList = make(map[string]NewConfStore)

// RegisterConfStore registers a ConfStore with the engine.
func RegisterConfStore(name string, creator NewConfStore) {
	name = strings.ToLower(name)
	if glog.V(1) {
		glog.Infof("registering ConfStore: %s", name)
	}

	if confStoreList[name] != nil {
		glog.Fatalf("a ConfStore named `%s` already exists", name)
	}
	confStoreList[name] = creator
}
