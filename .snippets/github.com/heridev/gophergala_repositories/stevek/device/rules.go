package device

type rule struct {
	apply  func(message, config) bool // TODO - need to stop passing config around
	action map[bool]action            // action to take on true/false return from apply
}
