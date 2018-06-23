package dsl

import (
	"command"
)

type V5624G struct {
	Name string
}

func (v5 V5624G) Port(typ, slot, port string) string {
	return port
}

func (v8 V5624G) InterfaceFromType(typ, name string) string {
	if typ == "vlan" {
		return " " + "br" + name
	} else {
		return " " + typ + " " + name
	}
}

var V5 = V5624G{
	Name: "V5",
}

func (v5 V5624G) PortSlotTypeEnable(Port, Slot, Type, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})
	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "port enable " + v5.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) PortSlotTypeDisable(Port, Slot, Type, Disable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})
	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "port disable " + v5.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) PortSlotTypeSpeed(Port, Slot, Type, Speed string) []*command.Command {

	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})
	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "port speed  " + Port + " " + Speed,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) PortSlotTypePvid(Port, Slot, Type, Pvid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})
	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan pvid " + Port + " " + Pvid,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLAN(VLAN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})

	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan create " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLAN(VLAN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})

	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "no vlan " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan add " + VLAN + " " + Port + " untagged",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan add " + VLAN + " " + Port + " tagged",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANDelTypeSlotPort(VLAN, Del, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan del " + VLAN + " " + Port,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANDelTTypeSlotPort(VLAN, DelT, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "bridge",
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "vlan del " + VLAN + " " + Port,
	})

	res = append(res, &command.Command{
		Mode: "bridge",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANShutdown(VLAN, Shutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "shutdown",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANNoShutdown(VLAN, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no shutdown",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANIP(VLAN, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip address " + IP,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLANIP(VLAN, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip address " + IP,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANIP2(VLAN, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip address " + IP2 + " secondary",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLANIP2(VLAN, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip address " + IP2 + " secondary",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceTypeIfname(Interface, Type, Ifname string) []*command.Command {
	return nil
}

func (v5 V5624G) VLANAddTypeSlotPortIP(VLAN, Add, Type, Slot, Port, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v5.VLANIP(VLAN, IP)...)
	return res
}

func (v5 V5624G) VLANAddTypeSlotPortIP2(VLAN, Add, Type, Slot, Port, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v5.VLANIP2(VLAN, IP2)...)
	return res
}

func (v5 V5624G) VLANAddTypeSlotPortIPNoShutdown(VLAN, Add, Type, Slot, Port, IP, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v5.VLANIP(VLAN, IP)...)
	res = append(res, v5.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v5 V5624G) VLANAddTypeSlotPortIP2NoShutdown(VLAN, Add, Type, Slot, Port, IP2, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v5.VLANIP2(VLAN, IP2)...)
	res = append(res, v5.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v5 V5624G) VLANAddTTypeSlotPortIP(VLAN, AddT, Type, Slot, Port, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v5.VLANIP(VLAN, IP)...)
	return res
}

func (v5 V5624G) VLANAddTTypeSlotPortIP2(VLAN, AddT, Type, Slot, Port, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v5.VLANIP2(VLAN, IP2)...)
	return res
}

func (v5 V5624G) VLANAddTTypeSlotPortIPNoShutdown(VLAN, AddT, Type, Slot, Port, IP, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v5.VLANIP(VLAN, IP)...)
	res = append(res, v5.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v5 V5624G) VLANAddTTypeSlotPortIP2NoShutdown(VLAN, AddT, Type, Slot, Port, IP2, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v5.VLAN(VLAN)...)
	res = append(res, v5.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v5.VLANIP2(VLAN, IP2)...)
	res = append(res, v5.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v5 V5624G) VLANIP6Enable(VLAN, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 enable",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLANIP6Enable(VLAN, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 enable",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANIP6(VLAN, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 address " + IP6,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLANIP6(VLAN, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 address " + IP6,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANIP6LL(VLAN, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 address link-local " + IP6LL,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) VLANIP6LLIP6(VLAN, IP6LL, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, v5.VLANIP6LL(VLAN, IP6LL)...)
	res = append(res, v5.VLANIP6(VLAN, IP6)...)
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoVLANIP6LL(VLAN, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 address link-local " + IP6LL,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6(OSPF6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6(OSPF6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "no router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6Rid(OSPF6, Rid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "router-id  " + Rid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6Rid(OSPF6, Rid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no router-id",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 router ospf area " + Area,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 router ospf area " + Area,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf cost " + Cost,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command {

	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf priority " + Priority,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf cost",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf dead-interval ",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf hello-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf retransmit-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf transmit-delay",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf priority",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface br" + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 router ospf area " + Area + " tag " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 router ospf area " + Area + " tag " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf cost " + Cost,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf priority " + Priority,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf cost",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf dead-interval ",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf hello-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf retransmit-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf transmit-delay",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf priority",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v5.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ipv6 ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "auto-cost reference-bandwidth " + ReferenceBandwidth,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no auto-cost reference-bandwidth",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always  metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6Redistribute(OSPF6, Redistribute string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no edistribute " + Redistribute + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6Redistribute(OSPF6, Redistribute string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no redistribute " + Redistribute,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) NoOSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) NoOSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) NoOSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) NoOSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) NoOSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	return v5.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v5 V5624G) OSPF6Summary(OSPF6, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "summary-address " + Summary,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "summary-address " + Summary + " " + "not-advertise",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6Summary(OSPF6, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no summary-address " + Summary,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no summary-address " + Summary + " " + " not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-metric " + DefaultMetric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-metric",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6Passive(OSPF6, Passive string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "passive-interface " + Passive,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6Passive(OSPF6, Passive string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no passive-interfaced " + Passive,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance " + AdminDistance,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance " + AdminDistance,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 inter-area " + Inter,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 intra-area " + Intra,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 inter-area " + Inter + " intra-area " + Intra,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 inter-area " + Inter + " external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospfv3 inter-area " + Inter + " intra-area " + Intra + " external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospfv3",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distribute-list " + Distributelist + " in ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distribute-list " + Distributelist + " out ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distribute-list " + Distributelist + " in",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distribute-list " + Distributelist + " out",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " default-cost " + DefaultCost,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " default-cost ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa default-information-originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, Redistribution string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa no-redistribution",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa stability-interval " + StabilityInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " translator-role " + Translatorrole,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa  default-information-originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, Redistribution string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa no-redistribution",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa stability-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa translator-role ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " stub ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " stub no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " stub",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " stub no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaRange(OSPF6, Area, Range string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range + " advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range + " not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaRange(OSPF6, Area, Range string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " instance-id " + Instanceid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + "   transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " instance-id " + Instanceid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + "   transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

//OSPF
func (v5 V5624G) OSPF(OSPF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPF(OSPF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "no router ospf " + OSPF,
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRid(OSPF, Rid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "router-id  " + Rid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFRid(OSPF, Rid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no router-id",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFNetworkArea(OSPF, Network, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "network " + Network + " area " + Area,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFNetworkArea(OSPF, Network, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "network " + Network + " area " + Area,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf cost " + Cost,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command {

	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf priority " + Priority,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ip ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf cost",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf dead-interval ",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf hello-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf retransmit-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf transmit-delay",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf priority",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no ip ospf network " + Networktype,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "auto-cost reference-bandwidth " + ReferenceBandwidth,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no auto-cost reference-bandwidth",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate route-map " + Routemap + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-information originate always route-map " + Routemap + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always  metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-information originate always metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistribute(OSPF, Redistribute string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no edistribute " + Redistribute + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype + " metric " + Metric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "redistribute " + Redistribute + " metric-type " + Metrictype + " metric " + Metric + " route-map " + Routemap,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFRedistribute(OSPF, Redistribute string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no redistribute " + Redistribute,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) NoOSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) NoOSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) NoOSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) NoOSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) NoOSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	return v5.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v5 V5624G) OSPFSummary(OSPF, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "summary-address " + Summary,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "summary-address " + Summary + " " + "not-advertise",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFSummary(OSPF, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no summary-address " + Summary,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no summary-address " + Summary + " not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "default-metric " + DefaultMetric,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no default-metric",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFPassive(OSPF, Passive string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "passive-interface " + Passive,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFPassive(OSPF, Passive string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no passive-interfaced " + Passive,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAdminDistance(OSPF, AdminDistance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance " + AdminDistance,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAdminDistance(OSPF, AdminDistance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance " + AdminDistance,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceExternal(OSPF, Distance, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf inter-area " + Inter,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf intra-area " + Intra,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf inter-area " + Inter + " intra-area " + Intra,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf inter-area " + Inter + " external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distance ospf inter-area " + Inter + " intra-area " + Intra + " external " + External,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceExternal(OSPF, Distance, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distance ospf",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distribute-list " + Distributelist + " in ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "distribute-list " + Distributelist + " out ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distribute-list " + Distributelist + " in",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no distribute-list " + Distributelist + " out",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " default-cost " + DefaultCost,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " default-cost ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa default-information-originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, Redistribution string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa no-redistribution",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaNSSANoSummary(OSPF, Area, NSSA, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " nssa stability-interval " + StabilityInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " translator-role " + Translatorrole,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa  default-information-originate",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, Redistribution string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa no-redistribution",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaNSSANoSummary(OSPF, Area, NSSA, Summary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa stability-interval",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " nssa translator-role ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaStub(OSPF, Area, Stub string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " stub ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " stub no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaStub(OSPF, Area, Stub string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " stub",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " stub no-summary",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaRange(OSPF, Area, Range string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range + " advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range + " not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaRange(OSPF, Area, Range string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " range " + Range,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " instance-id " + Instanceid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) OSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + "   transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " dead-interval " + DeadInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " hello-interval " + HelloInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + " instance-id " + Instanceid,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " virtual-link " + Virtuallink + " retransmit-interval " + RetransmitInterval,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoOSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no area " + Area + " virtual-link " + Virtuallink + "   transmit-delay " + TransmitDelay,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDeny(Routemap, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyMatchTag(Routemap, Deny, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoMatchTag(Routemap, Deny, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunity(Routemap, Deny, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetDampening(Routemap, Deny, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetLevel(Routemap, Deny, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetMetric(Routemap, Deny, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetOrigin(Routemap, Deny, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetTag(Routemap, Deny, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenySetWeight(Routemap, Deny, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetAggregatorAS(Routemap, Deny, Set, Aggregator, AS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunity(Routemap, Deny, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetDampening(Routemap, Deny, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetLevel(Routemap, Deny, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetMetric(Routemap, Deny, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetOrigin(Routemap, Deny, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetTag(Routemap, Deny, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapDenyNoSetWeight(Routemap, Deny, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Permit
func (v5 V5624G) RoutemapPermit(Routemap, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitMatchTag(Routemap, Permit, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoMatchTag(Routemap, Permit, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunity(Routemap, Permit, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetLevel(Routemap, Permit, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetMetric(Routemap, Permit, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetTag(Routemap, Permit, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitSetWeight(Routemap, Permit, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetAggregatorAS(Routemap, Permit, Set, Aggregator, AS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunity(Routemap, Permit, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetLevel(Routemap, Permit, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetMetric(Routemap, Permit, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetTag(Routemap, Permit, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) RoutemapPermitNoSetWeight(Routemap, Permit, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoRoutemap(Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoRoutemapDeny(Routemap, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoRoutemapPermit(Routemap, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDescription(IP, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDescription(IP, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Prefix List Permit

func (v5 V5624G) IPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

//IPv6 Prefix List

func (v5 V5624G) IP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Prefix List Permit

func (v5 V5624G) IP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) ClearIPRoute(Clear, IP, Route string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "clear ip route"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) ClearIPRouteKernel(Clear, IP, Route, Kernel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "clear ip route kernel"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRoutePrefixNexthopSrc(IP, Route, Prefix, Nexthop, Src string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop + " src " + Src})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route default " + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route default " + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route vrf " + VRF + " " + Prefix + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route default " + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route default " + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route vrf " + VRF + " " + Prefix + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixNexthopDistance(IP6, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixNexthopOIFDistance(IP6, Route, Prefix, Nexthop, OIF, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + OIF + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RoutePrefixTunnelDistance(IP6, Route, Prefix, Tunnel, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " tunnel " + Tunnel + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixNexthopDistance(IP6, Route, VRF, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixNexthopOIFDistance(IP6, Route, VRF, Prefix, Nexthop, OIF, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) IP6RouteVRFPrefixTunnelDistance(IP6, Route, VRF, Prefix, Tunnel, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RoutePrefix(IP6, Route, Prefix string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RouteVRFPrefix(IP6, Route, VRF, Prefix string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoIP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoIP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v5 V5624G) NoIP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceIP(Interface, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ip address " + IP})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceIP2(Interface, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ip address " + IP2 + " secondary"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceIP(Interface, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ip address " + IP})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceIP2(Interface, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ip address " + IP2 + " secondary"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceShutdown(Interface, Shutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "shutdown"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceNoShutdown(Interface, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no shutdown"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 enable "})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 enable "})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceIP6(Interface, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 address " + IP6})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) InterfaceIP6LL(Interface, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 address link-local " + IP6LL})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceIP6(Interface, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 address " + IP6})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v5 V5624G) NoInterfaceIP6LL(Interface, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 address link-local " + IP6LL})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}
