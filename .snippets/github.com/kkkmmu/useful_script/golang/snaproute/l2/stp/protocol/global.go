// global.go
package stp

const (
	// init
	STP_GLOBAL_INIT = iota + 1
	// allow all config
	STP_GLOBAL_ENABLE
	// disallow all config
	STP_GLOBAL_DISABLE
	// transition state to to allow deleting on
	// when global state changes to disable
	STP_GLOBAL_DISABLE_PENDING
)

var StpGlobalState int = STP_GLOBAL_INIT

func StpGlobalStateSet(state int) {
	StpGlobalState = state
}

func StpGlobalStateGet() int {
	return StpGlobalState
}
