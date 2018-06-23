package script

import (
	"command"
)

type Script struct {
	Name     string
	Count    int
	Timer    int
	Commands []*command.Command
}

/*
2017/06/14 22:29:17 server.go:174: Data[4][Command] [show running-config]
2017/06/14 22:29:17 server.go:174: Data[2][Command] [show running-config]
2017/06/14 22:29:17 server.go:174: Data[1][Mode] [enable]
2017/06/14 22:29:17 server.go:174: Data[3][Delay] [0]
2017/06/14 22:29:17 server.go:174: Data[0][Mode] [normal]
2017/06/14 22:29:17 server.go:174: Data[3][Command] [show running-config]
2017/06/14 22:29:17 server.go:174: Data[3][Expected] [#]
2017/06/14 22:29:17 server.go:174: Data[0][Command] [show running-config]
2017/06/14 22:29:17 server.go:174: Data[1][Command] [show running-config]
2017/06/14 22:29:17 server.go:174: Data[4][Delay] [0]
2017/06/14 22:29:17 server.go:174: Data[2][Expected] [#]
2017/06/14 22:29:17 server.go:174: Data[4][Mode] [enable]
2017/06/14 22:29:17 server.go:174: Data[0][Delay] [0]
2017/06/14 22:29:17 server.go:174: Data[3][Mode] [enable]
2017/06/14 22:29:17 server.go:174: Data[4][Expected] [#]
2017/06/14 22:29:17 server.go:174: Name [Scripts]
2017/06/14 22:29:17 server.go:174: Data[0][Expected] [vlan 100]
2017/06/14 22:29:17 server.go:174: Data[1][Delay] [0]
2017/06/14 22:29:17 server.go:174: Data[2][Mode] [enable]
2017/06/14 22:29:17 server.go:174: Data[2][Delay] [0]
2017/06/14 22:29:17 server.go:174: Data[1][Expected] [#]
*/
