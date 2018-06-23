package engine

import (
	"strings"

	"github.com/golang/glog"
)

// Reaction interface. It interfaces reactions! to things...
type Reaction interface {
	Setup
	StartStop
	Commander
}

// NewReaction is a function definition each Reaction must provide
type NewReaction func() (r Reaction)

var reactionList = make(map[string]NewReaction)

// RegisterReaction is called to register a Reaction object
func RegisterReaction(name string, creator NewReaction) {
	name = strings.ToLower(name)
	if glog.V(1) {
		glog.Infof("registering Reaction: %s", name)
	}

	if actionList[name] != nil {
		glog.Fatalf("a Reaction named `%s` already exists", name)
	}
	reactionList[name] = creator
}
