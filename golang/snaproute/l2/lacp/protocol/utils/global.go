// global.go
package utils

const (
	// init
	LACP_GLOBAL_INIT = iota + 1
	// allow all config
	LACP_GLOBAL_ENABLE
	// disallow all config
	LACP_GLOBAL_DISABLE
	// transition state to to allow deleting on
	// when global state changes to disable
	LACP_GLOBAL_DISABLE_PENDING
)

var LacpGlobalState int = LACP_GLOBAL_INIT

func LacpGlobalStateSet(state int) {
	LacpGlobalState = state
}

func LacpGlobalStateGet() int {
	return LacpGlobalState
}
