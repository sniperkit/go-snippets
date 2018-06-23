// global.go
package vxlan

const (
	// init
	VXLAN_GLOBAL_INIT = iota + 1
	// allow all config
	VXLAN_GLOBAL_ENABLE
	// disallow all config
	VXLAN_GLOBAL_DISABLE
	// transition state to to allow deleting on
	// when global state changes to disable
	VXLAN_GLOBAL_DISABLE_PENDING
)

var VxlanGlobalState int = VXLAN_GLOBAL_INIT

func VxlanGlobalStateSet(state int) {
	VxlanGlobalState = state
}

func VxlanGlobalStateGet() int {
	return VxlanGlobalState
}
