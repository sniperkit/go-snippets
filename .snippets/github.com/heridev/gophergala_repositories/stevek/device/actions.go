package device

// Should a packet continue onwards to the next stage
type action int

const (
	permit action = iota
	deny
	next
)

// define some common actions
var permitAction = make(map[bool]action)
var denyAction = make(map[bool]action)

func init() {
	permitAction[true] = permit
	permitAction[false] = next
	denyAction[true] = deny
	denyAction[false] = next
}
