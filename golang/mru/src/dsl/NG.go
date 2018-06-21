package dsl

import (
	"command"
)

type NG struct {
}

func (ng *NG) EnablePort(port string) []*command.Command {
	return nil
}

func (ng *NG) DisablePort(port string) []*command.Command {
	return nil
}

func (ng *NG) PortSpeed(port, speed string) []*command.Command {
	return nil
}

func (ng *NG) PortPvid(port, pvid string) []*command.Command {
	return nil
}

func (ng *NG) IsVlanExist(vid string) []*command.Command {
	return nil
}

func (ng *NG) CreateVlan(vid string) []*command.Command {
	return nil
}

func (ng *NG) DestroyVlan(vid string) []*command.Command {
	return nil
}

func (ng *NG) VlanAddUTPort(vid, port string) []*command.Command {
	return nil
}

func (ng *NG) VlanAddTPort(vid, port string) []*command.Command {
	return nil
}

func (ng *NG) VlanDelPort(vid, port string) []*command.Command {
	return nil
}

func (ng *NG) IsVlanInterfaceExist(vid string) []*command.Command {
	return nil
}

func (ng *NG) IsInterfaceExist(name string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceAddIP(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceDelIP(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceAddIP2(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceDelIP2(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceShutdown(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceNoShutdown(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceEnableIPv6(vid string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceDisableIPv6(vid string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceAddIPv6LL(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceDelIPv6LL(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceAddIPv6(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) VlanInterfaceDelIPv6(vid, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceAddIP(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceDelIP(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceEnableIPv6(name string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceDisableIPv6(name string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceAddIPv6(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceDelIPv6(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceAddIPv6LL(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceDelIPv6LL(name, ip string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceShutdown(name string) []*command.Command {
	return nil
}

func (ng *NG) L3InterfaceNoShutdown(name string) []*command.Command {
	return nil
}

func (ng *NG) IsOSPF6ProcessExist(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6CreateProcess(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6DestroyProcess(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6ClearProcess(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6ClearAllProcess() []*command.Command {
	return nil
}

func (ng *NG) OSPF6ShutdownProcess(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6NoShutdownProcess(id string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetProcessRID(id, rid string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetProcessRID(id, rid string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetAreaStub(id, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetAreaStub(id, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetAreaNSSA(id, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetAreaNSSA(id, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6EnableOSPFOnL3Interface(ifname, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6DisableOSPFOnL3Interface(ifname, area string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceCost(ifname, cost string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceCost(ifname, cost string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfacePriority(ifname, priority string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfacePriority(ifname, priority string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceHelloInterval(ifname, hi string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceHelloInterval(ifname, hi string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceDeadInterval(ifname, di string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceDeadInterval(ifname, di string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceRetranmitInterval(ifname, ri string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceRetransmitInterval(ifname, ri string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceTranmitDelay(ifname, td string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceTransmitDelay(ifname, td string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetL3InterfaceNetworkType(ifname, network string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetL3InterfaceNetworkType(ifname, network string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetRedistibution(id, rtype string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetRedistibution(id, rtype string) []*command.Command {
	return nil
}

func (ng *NG) OSPFSetRedistibutionRouteMap(id, rtype, mname string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetRedistibutionRouteMap(id, rtype, mname string) []*command.Command {

	return nil
}

func (ng *NG) OSPF6SetRedistibutionMetricType(id, rtype, mtype string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetRedistibutionMetricType(id, rtype, mytpe string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetRedistibutionMetric(id, rtype, metric string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetRedistibutionMetric(id, rtype, metric string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetRedistibutionMetricTypeMetric(id, rtype, mtype, metric string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetRedistibutionMetricTypeMetric(id, rtype, mtype, metric string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLink(area, neigh string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLink(area, neigh string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLinkHelloInterval(area, neigh, hi string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLinkHelloInterval(area, neigh, hi string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLinkDeadInterval(area, neigh, di string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLinkDeadInterval(area, neigh, di string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLinkRetransmitInterval(area, neigh, ri string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLinkRetransmitInterval(area, neigh, ri string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLinkTransmitDelay(area, neigh, td string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLinkTransmitDelay(area, neigh, td string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6SetVirtualLinkInstancID(area, neigh, instid string) []*command.Command {
	return nil
}

func (ng *NG) OSPF6UnSetVirtualLinkInstanceID(area, neigh, instid string) []*command.Command {
	return nil
}
