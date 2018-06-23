package server

import (
	"net"
	"testing"
)

func (server *ARPServer) SetPortPropertyMap() {
	portEnt, _ := server.portPropMap[20]
	portEnt.IfName = "fpPort20"
	portEnt.MacAddr = "11:22:33:44:55:66"
	portEnt.UntagVlanId = -1
	portEnt.L3PortPropMap = make(map[int]L3PortProp)
	l3PortPropEnt, _ := portEnt.L3PortPropMap[-1]
	l3PortPropEnt.IpAddr = "10.10.10.20"
	l3PortPropEnt.Netmask = net.IPMask([]byte{0xff, 0xff, 0xff, 0})
	l3PortPropEnt.L3IfIdx = 20
	l3PortPropEnt.LagIfIdx = -1
	portEnt.L3PortPropMap[-1] = l3PortPropEnt
	server.portPropMap[20] = portEnt
}

func (server *ARPServer) SetVlanPropertyMap() {

}

func (server *ARPServer) SetLagPropertyMap() {

}

func (server *ARPServer) SetL3PropertyMap() {
	l3Ent, _ := server.l3IntfPropMap[20]
	l3Ent.Netmask = net.IPMask([]byte{0xff, 0xff, 0xff, 0})
	l3Ent.IpAddr = "10.10.10.20"
	l3Ent.IfName = "fpPort20"
	server.l3IntfPropMap[20] = l3Ent
}

func TestProcessArpEntryUpdateMsg(t *testing.T) {
	t.Log("Testing initArpParams()")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}
	ser := NewARPServer(logger)

	msg := UpdateArpEntryMsg{
		PortIfIdx: 10,
		IpAddr:    "10.10.10.10",
		MacAddr:   "00:11:22:33:44:55",
		Type:      false,
		VlanId:    -1,
		L3IfIdx:   10,
		LagIfIdx:  -1,
	}

	ser.processArpEntryUpdateMsg(msg)
	_, exist := ser.arpCache["10.10.10.10"]
	if !exist {
		t.Log("Successfully test processArpEntryUpdateMsg() with invalid port")
	} else {
		t.Errorf("Successfully test processArpEntryUpdateMsg() with invalid port")
	}
	ser.SetPortPropertyMap()
	msg = UpdateArpEntryMsg{
		PortIfIdx: 20,
		IpAddr:    "10.10.10.10",
		MacAddr:   "00:11:22:33:44:55",
		Type:      false,
		VlanId:    -1,
		L3IfIdx:   10,
		LagIfIdx:  -1,
	}

	ser.processArpEntryUpdateMsg(msg)
	_, exist = ser.arpCache["10.10.10.10"]
	if !exist {
		t.Log("Successfully test processArpEntryUpdateMsg() with valid port but no l3 interface")
	} else {
		t.Errorf("Successfully test processArpEntryUpdateMsg() with valid port and no l3 interface")
	}

	ser.SetL3PropertyMap()
	msg = UpdateArpEntryMsg{
		PortIfIdx: 10,
		IpAddr:    "10.10.10.10",
		MacAddr:   "00:11:22:33:44:55",
		Type:      false,
		VlanId:    -1,
		L3IfIdx:   10,
		LagIfIdx:  -1,
	}

	ser.processArpEntryUpdateMsg(msg)
	_, exist = ser.arpCache["10.10.10.10"]
	if !exist {
		t.Log("Successfully test processArpEntryUpdateMsg() with valid port and invalid l3 interface")
	} else {
		t.Errorf("Successfully test processArpEntryUpdateMsg() with invalid port and invalid l3 interface")
	}
}
