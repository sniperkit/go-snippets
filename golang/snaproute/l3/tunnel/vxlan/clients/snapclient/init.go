package snapclient

import ()

func init() {
	PortVlanDb = make(map[uint16][]*portVlanValue, 0)

}
