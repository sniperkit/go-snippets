package dsl

import (
	"command"
)

type IRIS struct {
}

func (iris *IRIS) EnablePort(port string) []*command.Command {
	return nil
}

func (iris *IRIS) DisablePort(port string) []*command.Command {
	return nil
}

func (iris *IRIS) PortSpeed(port, speed string) []*command.Command {
	return nil
}

func (iris *IRIS) PortPvid(port, pvid string) []*command.Command {
	return nil
}

func (iris *IRIS) IsVlanExist(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) CreateVlan(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) DestroyVlan(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanAddUTPort(vid, port string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanAddTPort(vid, port string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanDelPort(vid, port string) []*command.Command {
	return nil
}

func (iris *IRIS) IsVlanInterfaceExist(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) IsInterfaceExist(name string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceAddIP(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceDelIP(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceAddIP2(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceDelIP2(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceShutdown(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceNoShutdown(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceEnableIPv6(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceDisableIPv6(vid string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceAddIPv6LL(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceDelIPv6LL(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceAddIPv6(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) VlanInterfaceDelIPv6(vid, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceAddIP(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceDelIP(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceEnableIPv6(name string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceDisableIPv6(name string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceAddIPv6(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceDelIPv6(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceAddIPv6LL(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceDelIPv6LL(name, ip string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceShutdown(name string) []*command.Command {
	return nil
}

func (iris *IRIS) L3InterfaceNoShutdown(name string) []*command.Command {
	return nil
}

func (iris *IRIS) IsOSPF6ProcessExist(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6CreateProcess(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6DestroyProcess(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6ClearProcess(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6ClearAllProcess() []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6ShutdownProcess(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6NoShutdownProcess(id string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetProcessRID(id, rid string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetProcessRID(id, rid string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetAreaStub(id, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetAreaStub(id, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetAreaNSSA(id, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetAreaNSSA(id, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6EnableOSPFOnL3Interface(ifname, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6DisableOSPFOnL3Interface(ifname, area string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceCost(ifname, cost string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceCost(ifname, cost string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfacePriority(ifname, priority string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfacePriority(ifname, priority string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceHelloInterval(ifname, hi string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceHelloInterval(ifname, hi string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceDeadInterval(ifname, di string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceDeadInterval(ifname, di string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceRetranmitInterval(ifname, ri string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceRetransmitInterval(ifname, ri string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceTranmitDelay(ifname, td string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceTransmitDelay(ifname, td string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetL3InterfaceNetworkType(ifname, network string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetL3InterfaceNetworkType(ifname, network string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetRedistibution(id, rtype string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetRedistibution(id, rtype string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPFSetRedistibutionRouteMap(id, rtype, mname string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetRedistibutionRouteMap(id, rtype, mname string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetRedistibutionMetricType(id, rtype, mtype string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetRedistibutionMetricType(id, rtype, mytpe string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetRedistibutionMetric(id, rtype, metric string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetRedistibutionMetric(id, rtype, metric string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetRedistibutionMetricTypeMetric(id, rtype, mtype, metric string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetRedistibutionMetricTypeMetric(id, rtype, mtype, metric string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLink(area, neigh string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLink(area, neigh string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLinkHelloInterval(area, neigh, hi string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLinkHelloInterval(area, neigh, hi string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLinkDeadInterval(area, neigh, di string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLinkDeadInterval(area, neigh, di string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLinkRetransmitInterval(area, neigh, ri string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLinkRetransmitInterval(area, neigh, ri string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLinkTransmitDelay(area, neigh, td string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLinkTransmitDelay(area, neigh, td string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6SetVirtualLinkInstancID(area, neigh, instid string) []*command.Command {
	return nil
}

func (iris *IRIS) OSPF6UnSetVirtualLinkInstanceID(area, neigh, instid string) []*command.Command {
	return nil
}
