package device

type message struct {
	Port   int    `json:port`
	Action string `json:action`
	Ghost  bool   `json:ghost`
	Device string `json:device`
}
