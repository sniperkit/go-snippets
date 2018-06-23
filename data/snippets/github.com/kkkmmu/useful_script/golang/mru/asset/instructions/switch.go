package dsl
import ( 
 "command"
 )
type Switch interface {
PortSlotTypeEnable(Port, Slot, Type, Enable string) []*command.Command
PortSlotTypeDisable(Port, Slot, Type, Disable string) []*command.Command
PortSlotTypeSpeed(Port, Slot, Type, Speed string) []*command.Command
PortSlotTypePvid(Port, Slot, Type, Pvid string) []*command.Command
VLAN(VLAN string) []*command.Command
NoVLAN(VLAN string) []*command.Command
VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port string) []*command.Command
VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port string) []*command.Command
VLANDelTypeSlotPort(VLAN, Del, Type, Slot, Port string) []*command.Command
VLANDelTTypeSlotPort(VLAN, DelT, Type, Slot, Port string) []*command.Command
VLANShutdown(VLAN, Shutdown string) []*command.Command
VLANNoShutdown(VLAN, NoShutdown string) []*command.Command
VLANIP(VLAN, IP string) []*command.Command
NoVLANIP(VLAN, IP string) []*command.Command
VLANIP2(VLAN, IP2 string) []*command.Command
NoVLANIP2(VLAN, IP2 string) []*command.Command
NoInterfaceTypeIfname(Interface, Type, Ifname string) []*command.Command
VLANAddTypeSlotPortIP(VLAN, Add, Type, Slot, Port, IP string) []*command.Command
VLANAddTypeSlotPortIP2(VLAN, Add, Type, Slot, Port, IP2 string) []*command.Command
VLANAddTypeSlotPortIPNoShutdown(VLAN, Add, Type, Slot, Port, IP, NoShutdown string) []*command.Command
VLANAddTypeSlotPortIP2NoShutdown(VLAN, Add, Type, Slot, Port, IP2, NoShutdown string) []*command.Command
VLANAddTTypeSlotPortIP(VLAN, AddT, Type, Slot, Port, IP string) []*command.Command
VLANAddTTypeSlotPortIP2(VLAN, AddT, Type, Slot, Port, IP2 string) []*command.Command
VLANAddTTypeSlotPortIPNoShutdown(VLAN, AddT, Type, Slot, Port, IP, NoShutdown string) []*command.Command
VLANAddTTypeSlotPortIP2NoShutdown(VLAN, AddT, Type, Slot, Port, IP2, NoShutdown string) []*command.Command
VLANIP6Enable(VLAN, IP6, Enable string) []*command.Command
NoVLANIP6Enable(VLAN, IP6, Enable string) []*command.Command
VLANIP6(VLAN, IP6 string) []*command.Command
NoVLANIP6(VLAN, IP6 string) []*command.Command
VLANIP6LL(VLAN, IP6LL string) []*command.Command
VLANIP6LLIP6(VLAN, IP6LL, IP6 string) []*command.Command
NoVLANIP6LL(VLAN, IP6LL string) []*command.Command
OSPF6(OSPF6 string) []*command.Command
NoOSPF6(OSPF6 string) []*command.Command
OSPF6Rid(OSPF6, Rid string) []*command.Command
NoOSPF6Rid(OSPF6, Rid string) []*command.Command
OSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command
NoOSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command
OSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command
OSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command
OSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command
OSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command
OSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command
OSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command
OSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command
NoOSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command
NoOSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command
NoOSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command
NoOSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command
NoOSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command
NoOSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command
NoOSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command
OSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command
NoOSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command
OSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command
OSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command
OSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command
OSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command
OSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command
OSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command
OSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command
NoOSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command
NoOSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command
NoOSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command
NoOSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command
NoOSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command
NoOSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command
NoOSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command
OSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command
NoOSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command
OSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command
OSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command
OSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command
OSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command
OSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command
OSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command
OSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command
OSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command
OSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command
OSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command
OSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command
OSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command
OSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command
OSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command
OSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command
OSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command
NoOSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command
NoOSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command
NoOSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command
NoOSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command
NoOSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command
NoOSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command
NoOSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command
NoOSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command
NoOSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command
NoOSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command
NoOSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command
OSPF6Redistribute(OSPF6, Redistribute string) []*command.Command
OSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command
OSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command
OSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command
OSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command
OSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command
OSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command
NoOSPF6Redistribute(OSPF6, Redistribute string) []*command.Command
NoOSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command
NoOSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command
NoOSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command
NoOSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command
NoOSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command
NoOSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command
OSPF6Summary(OSPF6, Summary string) []*command.Command
OSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command
NoOSPF6Summary(OSPF6, Summary string) []*command.Command
NoOSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command
OSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command
NoOSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command
OSPF6Passive(OSPF6, Passive string) []*command.Command
NoOSPF6Passive(OSPF6, Passive string) []*command.Command
OSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command
NoOSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command
OSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command
OSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command
OSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command
OSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command
OSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command
OSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command
NoOSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command
NoOSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command
NoOSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command
NoOSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command
NoOSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command
NoOSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command
OSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command
OSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command
NoOSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command
NoOSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command
OSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command
NoOSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command
OSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command
OSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command
OSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, NoRedistribution string) []*command.Command
OSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, NoSummary string) []*command.Command
OSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command
OSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command
NoOSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command
NoOSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command
NoOSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, NoRedistribution string) []*command.Command
NoOSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, NoSummary string) []*command.Command
NoOSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command
NoOSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command
OSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command
OSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command
NoOSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command
NoOSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command
OSPF6AreaRange(OSPF6, Area, Range string) []*command.Command
OSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command
OSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command
NoOSPF6AreaRange(OSPF6, Area, Range string) []*command.Command
NoOSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command
NoOSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command
OSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command
OSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command
OSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command
OSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command
OSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command
OSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command
NoOSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command
NoOSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command
NoOSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command
NoOSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command
NoOSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command
NoOSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command
OSPF(OSPF string) []*command.Command
NoOSPF(OSPF string) []*command.Command
OSPFRid(OSPF, Rid string) []*command.Command
NoOSPFRid(OSPF, Rid string) []*command.Command
OSPFNetworkArea(OSPF, Network, Area string) []*command.Command
NoOSPFNetworkArea(OSPF, Network, Area string) []*command.Command
OSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command
OSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command
OSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command
OSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command
OSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command
OSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command
OSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command
NoOSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command
NoOSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command
NoOSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command
NoOSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command
NoOSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command
NoOSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command
NoOSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command
OSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command
NoOSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command
OSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command
OSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command
OSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command
OSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command
OSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command
OSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command
OSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command
OSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command
OSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command
OSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command
OSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command
OSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command
OSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command
OSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command
OSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command
OSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command
NoOSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command
NoOSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command
NoOSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command
NoOSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command
NoOSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command
NoOSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command
NoOSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command
NoOSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command
NoOSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command
NoOSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command
NoOSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command
OSPFRedistribute(OSPF, Redistribute string) []*command.Command
OSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command
OSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command
OSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command
OSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command
OSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command
OSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command
NoOSPFRedistribute(OSPF, Redistribute string) []*command.Command
NoOSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command
NoOSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command
NoOSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command
NoOSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command
NoOSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command
NoOSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command
OSPFSummary(OSPF, Summary string) []*command.Command
OSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command
NoOSPFSummary(OSPF, Summary string) []*command.Command
NoOSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command
OSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command
NoOSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command
OSPFPassive(OSPF, Passive string) []*command.Command
NoOSPFPassive(OSPF, Passive string) []*command.Command
OSPFAdminDistance(OSPF, AdminDistance string) []*command.Command
NoOSPFAdminDistance(OSPF, AdminDistance string) []*command.Command
OSPFDistanceExternal(OSPF, Distance, External string) []*command.Command
OSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command
OSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command
OSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command
OSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command
OSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command
NoOSPFDistanceExternal(OSPF, Distance, External string) []*command.Command
NoOSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command
NoOSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command
NoOSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command
NoOSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command
NoOSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command
OSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command
OSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command
NoOSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command
NoOSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command
OSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command
NoOSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command
OSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command
OSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command
OSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, NoRedistribution string) []*command.Command
OSPFAreaNSSANoSummary(OSPF, Area, NSSA, NoSummary string) []*command.Command
OSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command
OSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command
NoOSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command
NoOSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command
NoOSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, NoRedistribution string) []*command.Command
NoOSPFAreaNSSANoSummary(OSPF, Area, NSSA, NoSummary string) []*command.Command
NoOSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command
NoOSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command
OSPFAreaStub(OSPF, Area, Stub string) []*command.Command
OSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command
NoOSPFAreaStub(OSPF, Area, Stub string) []*command.Command
NoOSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command
OSPFAreaRange(OSPF, Area, Range string) []*command.Command
OSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command
OSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command
NoOSPFAreaRange(OSPF, Area, Range string) []*command.Command
NoOSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command
NoOSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command
OSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command
OSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command
OSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command
OSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command
OSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command
OSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command
NoOSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command
NoOSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command
NoOSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command
NoOSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command
NoOSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command
NoOSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command
RoutemapDeny(Routemap, Deny string) []*command.Command
RoutemapDenyMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command
RoutemapDenyMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command
RoutemapDenyMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command
RoutemapDenyMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command
RoutemapDenyMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command
RoutemapDenyMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command
RoutemapDenyMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command
RoutemapDenyMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command
RoutemapDenyMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command
RoutemapDenyMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command
RoutemapDenyMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command
RoutemapDenyMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command
RoutemapDenyMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command
RoutemapDenyMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command
RoutemapDenyMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command
RoutemapDenyMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command
RoutemapDenyMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command
RoutemapDenyMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command
RoutemapDenyMatchTag(Routemap, Deny, Match, Tag string) []*command.Command
RoutemapDenyNoMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command
RoutemapDenyNoMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command
RoutemapDenyNoMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command
RoutemapDenyNoMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command
RoutemapDenyNoMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command
RoutemapDenyNoMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command
RoutemapDenyNoMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command
RoutemapDenyNoMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command
RoutemapDenyNoMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command
RoutemapDenyNoMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command
RoutemapDenyNoMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command
RoutemapDenyNoMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command
RoutemapDenyNoMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command
RoutemapDenyNoMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command
RoutemapDenyNoMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command
RoutemapDenyNoMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command
RoutemapDenyNoMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command
RoutemapDenyNoMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command
RoutemapDenyNoMatchTag(Routemap, Deny, Match, Tag string) []*command.Command
RoutemapDenySetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command
RoutemapDenySetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command
RoutemapDenySetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command
RoutemapDenySetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command
RoutemapDenySetCommunity(Routemap, Deny, Set, Community string) []*command.Command
RoutemapDenySetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command
RoutemapDenySetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command
RoutemapDenySetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command
RoutemapDenySetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command
RoutemapDenySetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command
RoutemapDenySetDampening(Routemap, Deny, Set, Dampening string) []*command.Command
RoutemapDenySetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command
RoutemapDenySetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command
RoutemapDenySetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command
RoutemapDenySetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command
RoutemapDenySetLevel(Routemap, Deny, Set, Level string) []*command.Command
RoutemapDenySetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command
RoutemapDenySetMetric(Routemap, Deny, Set, Metric string) []*command.Command
RoutemapDenySetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command
RoutemapDenySetOrigin(Routemap, Deny, Set, Origin string) []*command.Command
RoutemapDenySetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command
RoutemapDenySetTag(Routemap, Deny, Set, Tag string) []*command.Command
RoutemapDenySetWeight(Routemap, Deny, Set, Weight string) []*command.Command
RoutemapDenyNoSetAggregatorAS(Routemap, Deny, Set, Aggregator, AS string) []*command.Command
RoutemapDenyNoSetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command
RoutemapDenyNoSetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command
RoutemapDenyNoSetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command
RoutemapDenyNoSetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command
RoutemapDenyNoSetCommunity(Routemap, Deny, Set, Community string) []*command.Command
RoutemapDenyNoSetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command
RoutemapDenyNoSetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command
RoutemapDenyNoSetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command
RoutemapDenyNoSetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command
RoutemapDenyNoSetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command
RoutemapDenyNoSetDampening(Routemap, Deny, Set, Dampening string) []*command.Command
RoutemapDenyNoSetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command
RoutemapDenyNoSetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command
RoutemapDenyNoSetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command
RoutemapDenyNoSetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command
RoutemapDenyNoSetLevel(Routemap, Deny, Set, Level string) []*command.Command
RoutemapDenyNoSetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command
RoutemapDenyNoSetMetric(Routemap, Deny, Set, Metric string) []*command.Command
RoutemapDenyNoSetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command
RoutemapDenyNoSetOrigin(Routemap, Deny, Set, Origin string) []*command.Command
RoutemapDenyNoSetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command
RoutemapDenyNoSetTag(Routemap, Deny, Set, Tag string) []*command.Command
RoutemapDenyNoSetWeight(Routemap, Deny, Set, Weight string) []*command.Command
RoutemapPermit(Routemap, Permit string) []*command.Command
RoutemapPermitMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command
RoutemapPermitMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command
RoutemapPermitMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command
RoutemapPermitMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command
RoutemapPermitMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command
RoutemapPermitMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command
RoutemapPermitMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command
RoutemapPermitMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command
RoutemapPermitMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command
RoutemapPermitMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command
RoutemapPermitMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command
RoutemapPermitMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command
RoutemapPermitMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command
RoutemapPermitMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command
RoutemapPermitMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command
RoutemapPermitMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command
RoutemapPermitMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command
RoutemapPermitMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command
RoutemapPermitMatchTag(Routemap, Permit, Match, Tag string) []*command.Command
RoutemapPermitNoMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command
RoutemapPermitNoMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command
RoutemapPermitNoMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command
RoutemapPermitNoMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command
RoutemapPermitNoMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command
RoutemapPermitNoMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command
RoutemapPermitNoMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command
RoutemapPermitNoMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command
RoutemapPermitNoMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command
RoutemapPermitNoMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command
RoutemapPermitNoMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command
RoutemapPermitNoMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command
RoutemapPermitNoMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command
RoutemapPermitNoMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command
RoutemapPermitNoMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command
RoutemapPermitNoMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command
RoutemapPermitNoMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command
RoutemapPermitNoMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command
RoutemapPermitNoMatchTag(Routemap, Permit, Match, Tag string) []*command.Command
RoutemapPermitSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command
RoutemapPermitSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command
RoutemapPermitSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command
RoutemapPermitSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command
RoutemapPermitSetCommunity(Routemap, Permit, Set, Community string) []*command.Command
RoutemapPermitSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command
RoutemapPermitSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command
RoutemapPermitSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command
RoutemapPermitSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command
RoutemapPermitSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command
RoutemapPermitSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command
RoutemapPermitSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command
RoutemapPermitSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command
RoutemapPermitSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command
RoutemapPermitSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command
RoutemapPermitSetLevel(Routemap, Permit, Set, Level string) []*command.Command
RoutemapPermitSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command
RoutemapPermitSetMetric(Routemap, Permit, Set, Metric string) []*command.Command
RoutemapPermitSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command
RoutemapPermitSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command
RoutemapPermitSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command
RoutemapPermitSetTag(Routemap, Permit, Set, Tag string) []*command.Command
RoutemapPermitSetWeight(Routemap, Permit, Set, Weight string) []*command.Command
RoutemapPermitNoSetAggregatorAS(Routemap, Permit, Set, Aggregator, AS string) []*command.Command
RoutemapPermitNoSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command
RoutemapPermitNoSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command
RoutemapPermitNoSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command
RoutemapPermitNoSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command
RoutemapPermitNoSetCommunity(Routemap, Permit, Set, Community string) []*command.Command
RoutemapPermitNoSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command
RoutemapPermitNoSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command
RoutemapPermitNoSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command
RoutemapPermitNoSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command
RoutemapPermitNoSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command
RoutemapPermitNoSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command
RoutemapPermitNoSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command
RoutemapPermitNoSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command
RoutemapPermitNoSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command
RoutemapPermitNoSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command
RoutemapPermitNoSetLevel(Routemap, Permit, Set, Level string) []*command.Command
RoutemapPermitNoSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command
RoutemapPermitNoSetMetric(Routemap, Permit, Set, Metric string) []*command.Command
RoutemapPermitNoSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command
RoutemapPermitNoSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command
RoutemapPermitNoSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command
RoutemapPermitNoSetTag(Routemap, Permit, Set, Tag string) []*command.Command
RoutemapPermitNoSetWeight(Routemap, Permit, Set, Weight string) []*command.Command
NoRoutemap(Routemap string) []*command.Command
NoRoutemapDeny(Routemap, Deny string) []*command.Command
NoRoutemapPermit(Routemap, Permit string) []*command.Command
IPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command
NoIPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command
IPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command
IPPrefixListDescription(IP, PrefixList, Description string) []*command.Command
IPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command
IPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command
IPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command
IPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command
IPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command
IPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command
IPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command
IPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command
IPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command
NoIPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command
NoIPPrefixListDescription(IP, PrefixList, Description string) []*command.Command
NoIPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command
NoIPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command
NoIPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command
NoIPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command
NoIPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command
NoIPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command
NoIPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command
NoIPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command
NoIPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command
IPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command
IPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command
IPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command
IPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command
IPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command
IPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command
IPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command
IPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command
IPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command
IPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command
NoIPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command
NoIPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command
NoIPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command
NoIPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command
NoIPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command
NoIPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command
NoIPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command
NoIPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command
NoIPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command
NoIPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command
IP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command
NoIP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command
IP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command
IP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command
IP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command
IP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command
IP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command
IP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command
IP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command
IP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command
IP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command
IP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command
IP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command
NoIP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command
NoIP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command
NoIP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command
NoIP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command
NoIP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command
NoIP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command
NoIP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command
NoIP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command
NoIP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command
NoIP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command
NoIP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command
IP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command
IP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command
IP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command
IP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command
IP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command
IP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command
IP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command
IP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command
IP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command
IP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command
NoIP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command
NoIP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command
NoIP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command
NoIP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command
NoIP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command
NoIP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command
NoIP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command
NoIP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command
NoIP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command
NoIP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command
ClearIPRoute(Clear, IP, Route string) []*command.Command
ClearIPRouteKernel(Clear, IP, Route, Kernel string) []*command.Command
IPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command
IPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command
IPRoutePrefixNexthopSrc(IP, Route, Prefix, Nexthop, Src string) []*command.Command
IPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command
IPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command
IPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command
IPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command
NoIPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command
NoIPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command
NoIPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command
NoIPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command
NoIPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command
NoIPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command
IP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command
IP6RoutePrefixNexthopDistance(IP6, Route, Prefix, Nexthop, Distance string) []*command.Command
IP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command
IP6RoutePrefixNexthopOIFDistance(IP6, Route, Prefix, Nexthop, OIF, Distance string) []*command.Command
IP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command
IP6RoutePrefixTunnelDistance(IP6, Route, Prefix, Tunnel, Distance string) []*command.Command
IP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command
IP6RouteVRFPrefixNexthopDistance(IP6, Route, VRF, Prefix, Nexthop, Distance string) []*command.Command
IP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command
IP6RouteVRFPrefixNexthopOIFDistance(IP6, Route, VRF, Prefix, Nexthop, OIF, Distance string) []*command.Command
IP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command
IP6RouteVRFPrefixTunnelDistance(IP6, Route, VRF, Prefix, Tunnel, Distance string) []*command.Command
NoIP6RoutePrefix(IP6, Route, Prefix string) []*command.Command
NoIP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command
NoIP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command
NoIP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command
NoIP6RouteVRFPrefix(IP6, Route, VRF, Prefix string) []*command.Command
NoIP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command
NoIP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command
NoIP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command
InterfaceIP(Interface, IP string) []*command.Command
InterfaceIP2(Interface, IP2 string) []*command.Command
NoInterfaceIP(Interface, IP string) []*command.Command
NoInterfaceIP2(Interface, IP2 string) []*command.Command
InterfaceShutdown(Interface, Shutdown string) []*command.Command
InterfaceNoShutdown(Interface, NoShutdown string) []*command.Command
InterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command
NoInterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command
InterfaceIP6(Interface, IP6 string) []*command.Command
InterfaceIP6LL(Interface, IP6LL string) []*command.Command
NoInterfaceIP6(Interface, IP6 string) []*command.Command
NoInterfaceIP6LL(Interface, IP6LL string) []*command.Command
}
