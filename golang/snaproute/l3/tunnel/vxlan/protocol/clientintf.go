package vxlan

import (
	"net"
)

const (
	VXLANBaseClientStr = "BaseClient"
	VXLANSnapClientStr = "SnapClient"
	VXLANMockClientStr = "SnapMockTestClient"
)

type VtepCreateCfgData struct {
	VtepName string
	IfIndex  int32
}

type PortEvtCb func(ifindex int32)

// interface class is used to store the communication methods
// for the various daemon communications
type VXLANClientIntf interface {
	IsClientIntfType(client VXLANClientIntf, clientStr string) bool
	// used to notify server of updates
	SetServerChannels(s *VxLanConfigChannels)
	ConnectToClients(clientFile string, name string) error
	DisconnectFromClients(name string) error
	ConstructPortConfigMap()
	// create/delete
	CreateVtep(vtep *VtepDbEntry, vteplistener chan<- MachineEvent)
	DeleteVtep(vtep *VtepDbEntry)
	UpdateVtepAttr(vtepName string, vni uint32, tos, ttl uint8, mtu uint16)
	CreateVxlan(vxlan *VxlanDbEntry)
	DeleteVxlan(vxlan *VxlanDbEntry)
	UpdateVxlan(vni uint32, addvlanlist, delvlanlist, addUntaggedVlan, delUntaggedVlan []uint16)

	// access ports
	//AddHostToVxlan(vni int32, intfreflist, untagintfreflist []string)
	//DelHostFromVxlan(vni int32, intfreflist, untagintfreflist []string)
	// vtep fsm
	GetIntfInfo(name string, intfchan chan<- MachineEvent)
	GetNextHopInfo(ip net.IP, nexthopchan chan<- MachineEvent) bool
	UnRegisterReachability(ip net.IP)
	RegisterReachability(ip net.IP)
	UnresolveNextHopMac(nextHopIp net.IP, nexthopif string)
	ResolveNextHopMac(nextHopIp net.IP, nexthopif string, nexthopmacchan chan<- MachineEvent)

	GetLinkState(ifname string) string
	GetAllVlans() []uint16
	RegisterLinkUpDownEvents(ifindex int32, upcb PortEvtCb, downdb PortEvtCb)
}

// SetIntf:
// The user may implement mulitple interfaces for uses
// by the server.  This was created to avoid import cycle
func RegisterClients(intf VXLANClientIntf) {
	if ClientIntf == nil {
		ClientIntf = make([]VXLANClientIntf, 0)
	}
	ClientIntf = append(ClientIntf, intf)
}

func DeRegisterClients() {
	for _, client := range ClientIntf {
		client.DisconnectFromClients("")
	}
	ClientIntf = nil
}

type BaseClientIntf struct {
}

func (b BaseClientIntf) IsClientIntfType(client VXLANClientIntf, clientStr string) bool {
	switch client.(type) {
	case *BaseClientIntf:
		if clientStr == "BaseClient" {
			return true
		}
	}
	return false
}

func (b BaseClientIntf) SetServerChannels(s *VxLanConfigChannels) {

}
func (b BaseClientIntf) ConnectToClients(clientFile string, name string) error {
	return nil
}
func (b BaseClientIntf) DisconnectFromClients(name string) error {
	return nil
}
func (b BaseClientIntf) ConstructPortConfigMap() {

}
func (b BaseClientIntf) GetIntfInfo(name string, intfchan chan<- MachineEvent) {

}
func (b BaseClientIntf) CreateVtep(vtep *VtepDbEntry, vteplistener chan<- MachineEvent) {

}
func (b BaseClientIntf) UpdateVtepAttr(vtepName string, vni uint32, tos, ttl uint8, mtu uint16) {

}
func (b BaseClientIntf) DeleteVtep(vtep *VtepDbEntry) {

}
func (b BaseClientIntf) CreateVxlan(vxlan *VxlanDbEntry) {

}
func (b BaseClientIntf) DeleteVxlan(vxlan *VxlanDbEntry) {

}
func (b BaseClientIntf) UpdateVxlan(
	vni uint32, addvlanlist, delvlanlist, addUntaggedVlan, delUntaggedVlan []uint16) {

}
func (b BaseClientIntf) CreateVxlanAccess() {

}
func (b BaseClientIntf) DeleteVxlanAccess() {

}
func (b BaseClientIntf) GetNextHopInfo(ip net.IP, nexthopchan chan<- MachineEvent) bool {
	return true
}
func (b BaseClientIntf) UnRegisterReachability(ip net.IP) {

}
func (b BaseClientIntf) RegisterReachability(ip net.IP) {

}
func (b BaseClientIntf) ResolveNextHopMac(nextHopIp net.IP, nextHopIfName string, nexthopmacchan chan<- MachineEvent) {

}
func (b BaseClientIntf) UnresolveNextHopMac(nextHopIp net.IP, nexthopif string) {

}
func (b BaseClientIntf) AddHostToVxlan(vni int32, intfreflist, untagintfreflist []string) {

}
func (b BaseClientIntf) DelHostFromVxlan(vni int32, intfreflist, untagintfreflist []string) {

}
func (b BaseClientIntf) GetAllVlans() []uint16 {
	return []uint16{}
}
func (b BaseClientIntf) RegisterLinkUpDownEvents(ifindex int32, upcb PortEvtCb, downdb PortEvtCb) {

}
func (b BaseClientIntf) GetLinkState(ifname string) string {
	return "UP"
}
