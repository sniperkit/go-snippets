package dsl

import (
	"command"
	"log"
)

type V8500 struct {
	Name string
}

func (v8 V8500) Port(typ, slot, port string) string {
	return " " + typ + " " + slot + "/" + port
}

func (v8 V8500) InterfaceFromType(typ, name string) string {
	if typ == "vlan" {
		return " " + "vlan" + " " + name
	} else {
		return " " + typ + " " + name
	}
}

var V8 = V8500{
	Name: "V8",
}

func (v8 V8500) PortSlotTypeEnable(Port, Slot, Type, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
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

func (v8 V8500) PortSlotTypeDisable(Port, Slot, Type, Disable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})

	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
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

func (v8 V8500) PortSlotTypeSpeed(Port, Slot, Type, Speed string) []*command.Command {

	var local string
	if Speed == "1000" {
		local = "1g"
	} else if Speed == "100" {
		local = "100m"
	} else if Speed == "10" {
		local = "10m"
	} else {
		log.Printf("Invalid spped set % set port :%s speed to 1g ", Speed, Port)
		local = "1g"
	}

	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})
	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "bandwidth " + local,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v8 V8500) PortSlotTypePvid(Port, Slot, Type, Pvid string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})
	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "pvid " + Pvid,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLAN(VLAN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})

	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "vlan database",
	})

	res = append(res, &command.Command{
		Mode: "config-vlan",
		CMD:  "vlan " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-vlan",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoVLAN(VLAN string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})

	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "vlan database",
	})

	res = append(res, &command.Command{
		Mode: "config-vlan",
		CMD:  "no vlan " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-vlan",
		CMD:  "exit",
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "switchport mode access",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "switchport access vlan " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "switchport mode trunk",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "switchport trunk allowed vlan add " + VLAN,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANDelTypeSlotPort(VLAN, Del, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no switchport access vlan",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANDelTTypeSlotPort(VLAN, DelT, Type, Slot, Port string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.Port(Type, Slot, Port),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "no switchport access trunk",
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANShutdown(VLAN, Shutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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
	res = append(res, v8.NoInterfaceTypeIfname("", "vlan", VLAN)...)
	return res
}

func (v8 V8500) VLANNoShutdown(VLAN, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) VLANIP(VLAN, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) NoVLANIP(VLAN, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) VLANIP2(VLAN, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) NoVLANIP2(VLAN, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) NoInterfaceTypeIfname(Interface, Type, Ifname string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "no interface " + Type + " " + Ifname,
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) VLANAddTypeSlotPortIP(VLAN, Add, Type, Slot, Port, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v8.VLANIP(VLAN, IP)...)
	return res
}

func (v8 V8500) VLANAddTypeSlotPortIP2(VLAN, Add, Type, Slot, Port, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v8.VLANIP2(VLAN, IP2)...)
	return res
}

func (v8 V8500) VLANAddTypeSlotPortIPNoShutdown(VLAN, Add, Type, Slot, Port, IP, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v8.VLANIP(VLAN, IP)...)
	res = append(res, v8.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v8 V8500) VLANAddTypeSlotPortIP2NoShutdown(VLAN, Add, Type, Slot, Port, IP2, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTypeSlotPort(VLAN, Add, Type, Slot, Port)...)
	res = append(res, v8.VLANIP2(VLAN, IP2)...)
	res = append(res, v8.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v8 V8500) VLANAddTTypeSlotPortIP(VLAN, AddT, Type, Slot, Port, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v8.VLANIP(VLAN, IP)...)
	return res
}

func (v8 V8500) VLANAddTTypeSlotPortIP2(VLAN, AddT, Type, Slot, Port, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v8.VLANIP2(VLAN, IP2)...)
	return res
}

func (v8 V8500) VLANAddTTypeSlotPortIPNoShutdown(VLAN, AddT, Type, Slot, Port, IP, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v8.VLANIP(VLAN, IP)...)
	res = append(res, v8.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v8 V8500) VLANAddTTypeSlotPortIP2NoShutdown(VLAN, AddT, Type, Slot, Port, IP2, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, v8.VLAN(VLAN)...)
	res = append(res, v8.VLANAddTTypeSlotPort(VLAN, AddT, Type, Slot, Port)...)
	res = append(res, v8.VLANIP2(VLAN, IP2)...)
	res = append(res, v8.VLANNoShutdown(VLAN, NoShutdown)...)
	return res
}

func (v8 V8500) VLANIP6Enable(VLAN, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) NoVLANIP6Enable(VLAN, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) VLANIP6(VLAN, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) NoVLANIP6(VLAN, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) VLANIP6LL(VLAN, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) VLANIP6LLIP6(VLAN, IP6LL, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, v8.VLANIP6LL(VLAN, IP6LL)...)
	res = append(res, v8.VLANIP6(VLAN, IP6)...)
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoVLANIP6LL(VLAN, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + VLAN,
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

func (v8 V8500) OSPF6(OSPF6 string) []*command.Command {
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

func (v8 V8500) NoOSPF6(OSPF6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "no router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) OSPF6Rid(OSPF6, Rid string) []*command.Command {
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

func (v8 V8500) NoOSPF6Rid(OSPF6, Rid string) []*command.Command {
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

func (v8 V8500) OSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceArea(OSPF6, Interface, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command {

	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceCost(OSPF6, Interface, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceDeadInterval(OSPF6, Interface, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceHelloInterval(OSPF6, Interface, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceRetransmitInterval(OSPF6, Interface, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceTransmitDelay(OSPF6, Interface, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfacePriority(OSPF6, Interface, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) NoOSPF6InterfaceNetworktype(OSPF6, Interface, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface vlan " + Interface,
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

func (v8 V8500) OSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.InterfaceFromType(Type, Name),
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "ipv6 router ospf tag " + OSPF6 + " area " + Area,
	})

	res = append(res, &command.Command{
		Mode: "config-if",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoOSPF6IfTypeNameArea(OSPF6, If, Type, Name, Area string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameCost(OSPF6, If, Type, Name, Cost string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameDeadInterval(OSPF6, If, Type, Name, DeadInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameHelloInterval(OSPF6, If, Type, Name, HelloInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameRetransmitInterval(OSPF6, If, Type, Name, RetransmitInterval string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameTransmitDelay(OSPF6, If, Type, Name, TransmitDelay string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNamePriority(OSPF6, If, Type, Name, Priority string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) NoOSPF6IfTypeNameNetworktype(OSPF6, If, Type, Name, Networktype string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "interface  " + v8.InterfaceFromType(Type, Name),
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

func (v8 V8500) OSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command {
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

func (v8 V8500) NoOSPF6ReferenceBandwidth(OSPF6, ReferenceBandwidth string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginate(OSPF6, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateRoutemap(OSPF6, DefaultOriginate, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetric(OSPF6, DefaultOriginate, Metric string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetrictype(OSPF6, DefaultOriginate, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetricMetrictype(OSPF6, DefaultOriginate, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetricRoutemap(OSPF6, DefaultOriginate, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetrictypeRoutemap(OSPF6, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlways(OSPF6, DefaultOriginate, Always string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysRoutemap(OSPF6, DefaultOriginate, Always, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetric(OSPF6, DefaultOriginate, Always, Metric string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetrictype(OSPF6, DefaultOriginate, Always, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetricMetrictype(OSPF6, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetricRoutemap(OSPF6, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF6, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6Redistribute(OSPF6, Redistribute string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPF6Redistribute(OSPF6, Redistribute string) []*command.Command {
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

func (v8 V8500) NoOSPF6RedistributeMetric(OSPF6, Redistribute, Metric string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) NoOSPF6RedistributeMetrictype(OSPF6, Redistribute, Metrictype string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) NoOSPF6RedistributeRoutemap(OSPF6, Redistribute, Routemap string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) NoOSPF6RedistributeMetricMetrictype(OSPF6, Redistribute, Metric, Metrictype string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) NoOSPF6RedistributeMetricRoutemap(OSPF6, Redistribute, Metric, Routemap string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) NoOSPF6RedistributeMetricMetrictypeRoutemap(OSPF6, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	return v8.NoOSPF6Redistribute(OSPF6, Redistribute)
}

func (v8 V8500) OSPF6Summary(OSPF6, Summary string) []*command.Command {
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

func (v8 V8500) OSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command {
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

func (v8 V8500) NoOSPF6Summary(OSPF6, Summary string) []*command.Command {
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

func (v8 V8500) NoOSPF6SummaryNoAdvertise(OSPF6, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "no summary-address " + Summary + " " + "not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) OSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command {
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

func (v8 V8500) NoOSPF6DefaultMetric(OSPF6, DefaultMetric string) []*command.Command {
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

func (v8 V8500) OSPF6Passive(OSPF6, Passive string) []*command.Command {
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

func (v8 V8500) NoOSPF6Passive(OSPF6, Passive string) []*command.Command {
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

func (v8 V8500) OSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command {
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

func (v8 V8500) NoOSPF6AdminDistance(OSPF6, AdminDistance string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command {
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

func (v8 V8500) OSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceExternal(OSPF6, Distance, External string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceInter(OSPF6, Distance, Inter string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceIntra(OSPF6, Distance, Intra string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceInterIntra(OSPF6, Distance, Inter, Intra string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceInterExternal(OSPF6, Distance, Inter, External string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistanceInterIntraExternal(OSPF6, Distance, Inter, Intra, External string) []*command.Command {
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

func (v8 V8500) OSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command {
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

func (v8 V8500) OSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistributelistIN(OSPF6, Distributelist, IN string) []*command.Command {
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

func (v8 V8500) NoOSPF6DistributelistOUT(OSPF6, Distributelist, OUT string) []*command.Command {
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

func (v8 V8500) OSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaDefaultCost(OSPF6, Area, DefaultCost string) []*command.Command {
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

func (v8 V8500) OSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command {
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

func (v8 V8500) OSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) OSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, Redistribution string) []*command.Command {
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

func (v8 V8500) OSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, Summary string) []*command.Command {
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

func (v8 V8500) OSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command {
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

func (v8 V8500) OSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaNSSA(OSPF6, Area, NSSA string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaNSSADefaultOriginate(OSPF6, Area, NSSA, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaNSSANoRedistribution(OSPF6, Area, NSSA, Redistribution string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaNSSANoSummary(OSPF6, Area, NSSA, Summary string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaNSSAStabilityInterval(OSPF6, Area, NSSA, StabilityInterval string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaTranslatorrole(OSPF6, Area, Translatorrole string) []*command.Command {
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

func (v8 V8500) OSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command {
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

func (v8 V8500) OSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaStub(OSPF6, Area, Stub string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaStubNoSummary(OSPF6, Area, Stub, NoSummary string) []*command.Command {
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

func (v8 V8500) OSPF6AreaRange(OSPF6, Area, Range string) []*command.Command {
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

func (v8 V8500) OSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command {
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

func (v8 V8500) OSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ipv6 ospf " + OSPF6,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "area " + Area + " range " + Range + " " + "not-advertise ",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoOSPF6AreaRange(OSPF6, Area, Range string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaRangeAdvertise(OSPF6, Area, Range, Advertise string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaRangeNoAdvertise(OSPF6, Area, Range, NoAdvertise string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) OSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallink(OSPF6, Area, Virtuallink string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallinkDeadInterval(OSPF6, Area, Virtuallink, DeadInterval string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallinkHelloInterval(OSPF6, Area, Virtuallink, HelloInterval string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallinkInstanceid(OSPF6, Area, Virtuallink, Instanceid string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallinkRetransmitInterval(OSPF6, Area, Virtuallink, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) NoOSPF6AreaVirtuallinkTransmitDelay(OSPF6, Area, Virtuallink, TransmitDelay string) []*command.Command {
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
func (v8 V8500) OSPF(OSPF string) []*command.Command {
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

func (v8 V8500) NoOSPF(OSPF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "no router ospf " + OSPF,
	})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) OSPFRid(OSPF, Rid string) []*command.Command {
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

func (v8 V8500) NoOSPFRid(OSPF, Rid string) []*command.Command {
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

func (v8 V8500) OSPFNetworkArea(OSPF, Network, Area string) []*command.Command {
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

func (v8 V8500) NoOSPFNetworkArea(OSPF, Network, Area string) []*command.Command {
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

func (v8 V8500) OSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command {
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

func (v8 V8500) OSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command {
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

func (v8 V8500) OSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command {

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

func (v8 V8500) OSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) OSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command {
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

func (v8 V8500) OSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command {
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

func (v8 V8500) OSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceCost(OSPF, Interface, Cost string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceDeadInterval(OSPF, Interface, DeadInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceHelloInterval(OSPF, Interface, HelloInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceRetransmitInterval(OSPF, Interface, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceTransmitDelay(OSPF, Interface, TransmitDelay string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfacePriority(OSPF, Interface, Priority string) []*command.Command {
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

func (v8 V8500) NoOSPFInterfaceNetworktype(OSPF, Interface, Networktype string) []*command.Command {
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

func (v8 V8500) OSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command {
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

func (v8 V8500) NoOSPFReferenceBandwidth(OSPF, ReferenceBandwidth string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginate(OSPF, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateRoutemap(OSPF, DefaultOriginate, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetric(OSPF, DefaultOriginate, Metric string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetrictype(OSPF, DefaultOriginate, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetricMetrictype(OSPF, DefaultOriginate, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetricRoutemap(OSPF, DefaultOriginate, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetrictypeRoutemap(OSPF, DefaultOriginate, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlways(OSPF, DefaultOriginate, Always string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysRoutemap(OSPF, DefaultOriginate, Always, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetric(OSPF, DefaultOriginate, Always, Metric string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetrictype(OSPF, DefaultOriginate, Always, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetricMetrictype(OSPF, DefaultOriginate, Always, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetricRoutemap(OSPF, DefaultOriginate, Always, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(OSPF, DefaultOriginate, Always, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFRedistribute(OSPF, Redistribute string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command {
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

func (v8 V8500) OSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
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

func (v8 V8500) NoOSPFRedistribute(OSPF, Redistribute string) []*command.Command {
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

func (v8 V8500) NoOSPFRedistributeMetric(OSPF, Redistribute, Metric string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) NoOSPFRedistributeMetrictype(OSPF, Redistribute, Metrictype string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) NoOSPFRedistributeRoutemap(OSPF, Redistribute, Routemap string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) NoOSPFRedistributeMetricMetrictype(OSPF, Redistribute, Metric, Metrictype string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) NoOSPFRedistributeMetricRoutemap(OSPF, Redistribute, Metric, Routemap string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) NoOSPFRedistributeMetricMetrictypeRoutemap(OSPF, Redistribute, Metric, Metrictype, Routemap string) []*command.Command {
	return v8.NoOSPFRedistribute(OSPF, Redistribute)
}

func (v8 V8500) OSPFSummary(OSPF, Summary string) []*command.Command {
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

func (v8 V8500) OSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{
		Mode: "config",
		CMD:  "router ospf " + OSPF,
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "summary-address " + Summary + " not-advertise",
	})

	res = append(res, &command.Command{
		Mode: "config-router",
		CMD:  "exit",
	})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoOSPFSummary(OSPF, Summary string) []*command.Command {
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

func (v8 V8500) NoOSPFSummaryNoAdvertise(OSPF, Summary, NoAdvertise string) []*command.Command {
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

func (v8 V8500) OSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command {
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

func (v8 V8500) NoOSPFDefaultMetric(OSPF, DefaultMetric string) []*command.Command {
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

func (v8 V8500) OSPFPassive(OSPF, Passive string) []*command.Command {
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

func (v8 V8500) NoOSPFPassive(OSPF, Passive string) []*command.Command {
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

func (v8 V8500) OSPFAdminDistance(OSPF, AdminDistance string) []*command.Command {
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

func (v8 V8500) NoOSPFAdminDistance(OSPF, AdminDistance string) []*command.Command {
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

func (v8 V8500) OSPFDistanceExternal(OSPF, Distance, External string) []*command.Command {
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

func (v8 V8500) OSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command {
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

func (v8 V8500) OSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command {
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

func (v8 V8500) OSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command {
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

func (v8 V8500) OSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command {
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

func (v8 V8500) OSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceExternal(OSPF, Distance, External string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceInter(OSPF, Distance, Inter string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceIntra(OSPF, Distance, Intra string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceInterIntra(OSPF, Distance, Inter, Intra string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceInterExternal(OSPF, Distance, Inter, External string) []*command.Command {
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

func (v8 V8500) NoOSPFDistanceInterIntraExternal(OSPF, Distance, Inter, Intra, External string) []*command.Command {
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

func (v8 V8500) OSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command {
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

func (v8 V8500) OSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command {
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

func (v8 V8500) NoOSPFDistributelistIN(OSPF, Distributelist, IN string) []*command.Command {
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

func (v8 V8500) NoOSPFDistributelistOUT(OSPF, Distributelist, OUT string) []*command.Command {
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

func (v8 V8500) OSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaDefaultCost(OSPF, Area, DefaultCost string) []*command.Command {
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

func (v8 V8500) OSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command {
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

func (v8 V8500) OSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) OSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, Redistribution string) []*command.Command {
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

func (v8 V8500) OSPFAreaNSSANoSummary(OSPF, Area, NSSA, Summary string) []*command.Command {
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

func (v8 V8500) OSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command {
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

func (v8 V8500) OSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaNSSA(OSPF, Area, NSSA string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaNSSADefaultOriginate(OSPF, Area, NSSA, DefaultOriginate string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaNSSANoRedistribution(OSPF, Area, NSSA, Redistribution string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaNSSANoSummary(OSPF, Area, NSSA, Summary string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaNSSAStabilityInterval(OSPF, Area, NSSA, StabilityInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaTranslatorrole(OSPF, Area, Translatorrole string) []*command.Command {
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

func (v8 V8500) OSPFAreaStub(OSPF, Area, Stub string) []*command.Command {
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

func (v8 V8500) OSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaStub(OSPF, Area, Stub string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaStubNoSummary(OSPF, Area, Stub, NoSummary string) []*command.Command {
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

func (v8 V8500) OSPFAreaRange(OSPF, Area, Range string) []*command.Command {
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

func (v8 V8500) OSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command {
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

func (v8 V8500) OSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaRange(OSPF, Area, Range string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaRangeAdvertise(OSPF, Area, Range, Advertise string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaRangeNoAdvertise(OSPF, Area, Range, NoAdvertise string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) OSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallink(OSPF, Area, Virtuallink string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallinkDeadInterval(OSPF, Area, Virtuallink, DeadInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallinkHelloInterval(OSPF, Area, Virtuallink, HelloInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallinkInstanceid(OSPF, Area, Virtuallink, Instanceid string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallinkRetransmitInterval(OSPF, Area, Virtuallink, RetransmitInterval string) []*command.Command {
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

func (v8 V8500) NoOSPFAreaVirtuallinkTransmitDelay(OSPF, Area, Virtuallink, TransmitDelay string) []*command.Command {
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

func (v8 V8500) RoutemapDeny(Routemap, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyMatchTag(Routemap, Deny, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchASPath(Routemap, Deny, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchCLNSAddress(Routemap, Deny, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchCommunity(Routemap, Deny, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchCommunityExactMatch(Routemap, Deny, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchExtCommunity(Routemap, Deny, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchExtCommunityExactMatch(Routemap, Deny, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchInterfaceTypeName(Routemap, Deny, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIPAddressAccessList(Routemap, Deny, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIPAddressPrefixList(Routemap, Deny, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIPNexthopAccessList(Routemap, Deny, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIPNexthopPrefixList(Routemap, Deny, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIP6AddressAccessList(Routemap, Deny, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIP6AddressPrefixList(Routemap, Deny, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIP6NexthopAccessList(Routemap, Deny, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchIP6NexthopPrefixList(Routemap, Deny, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchMetric(Routemap, Deny, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchOrigin(Routemap, Deny, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchRouteTypeExternal(Routemap, Deny, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoMatchTag(Routemap, Deny, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunity(Routemap, Deny, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetDampening(Routemap, Deny, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetLevel(Routemap, Deny, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetMetric(Routemap, Deny, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetOrigin(Routemap, Deny, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetTag(Routemap, Deny, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenySetWeight(Routemap, Deny, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetAggregatorAS(Routemap, Deny, Set, Aggregator, AS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetAggregatorASIP(Routemap, Deny, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetASPathPrepend(Routemap, Deny, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetAtomicAggregator(Routemap, Deny, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommListDelete(Routemap, Deny, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunity(Routemap, Deny, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunityNone(Routemap, Deny, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunityInternet(Routemap, Deny, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunityLocalAS(Routemap, Deny, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunityNoAdvertise(Routemap, Deny, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetCommunityNoExport(Routemap, Deny, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetDampening(Routemap, Deny, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetExtCommunityRT(Routemap, Deny, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetExtCommunitySOO(Routemap, Deny, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetIPNexthop(Routemap, Deny, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetIP6Nexthop(Routemap, Deny, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetLevel(Routemap, Deny, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetLocalPreference(Routemap, Deny, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetMetric(Routemap, Deny, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetMetricType(Routemap, Deny, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetOrigin(Routemap, Deny, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetOriginatorID(Routemap, Deny, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetTag(Routemap, Deny, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapDenyNoSetWeight(Routemap, Deny, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Permit
func (v8 V8500) RoutemapPermit(Routemap, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})

	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitMatchTag(Routemap, Permit, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchASPath(Routemap, Permit, Match, ASPath string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match as-path " + ASPath})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchCLNSAddress(Routemap, Permit, Match, CLNS, Address string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match clns address " + Address})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchCommunity(Routemap, Permit, Match, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchCommunityExactMatch(Routemap, Permit, Match, Community, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match community " + Community + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchExtCommunity(Routemap, Permit, Match, ExtCommunity string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchExtCommunityExactMatch(Routemap, Permit, Match, ExtCommunity, ExactMatch string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match extcommunity " + ExtCommunity + " exact-match"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchInterfaceTypeName(Routemap, Permit, Match, Interface, Type, Name string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match interface " + Type + " " + Name})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIPAddressAccessList(Routemap, Permit, Match, IP, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIPAddressPrefixList(Routemap, Permit, Match, IP, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIPNexthopAccessList(Routemap, Permit, Match, IP, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIPNexthopPrefixList(Routemap, Permit, Match, IP, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ip nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIP6AddressAccessList(Routemap, Permit, Match, IP6, Address, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIP6AddressPrefixList(Routemap, Permit, Match, IP6, Address, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 address prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIP6NexthopAccessList(Routemap, Permit, Match, IP6, Nexthop, AccessList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop " + AccessList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchIP6NexthopPrefixList(Routemap, Permit, Match, IP6, Nexthop, PrefixList string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match ipv6 nexthop prefix-list " + PrefixList})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchMetric(Routemap, Permit, Match, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchOrigin(Routemap, Permit, Match, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchRouteTypeExternal(Routemap, Permit, Match, RouteType, External string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match route-type external " + External})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoMatchTag(Routemap, Permit, Match, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no match tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunity(Routemap, Permit, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetLevel(Routemap, Permit, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetMetric(Routemap, Permit, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetTag(Routemap, Permit, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitSetWeight(Routemap, Permit, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetAggregatorAS(Routemap, Permit, Set, Aggregator, AS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetAggregatorASIP(Routemap, Permit, Set, Aggregator, AS, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set aggregator as " + AS + " " + IP})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetASPathPrepend(Routemap, Permit, Set, ASPath, Prepend string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set as-path prepend " + Prepend})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetAtomicAggregator(Routemap, Permit, Set, AtomicAggregator string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set atomic-aggregate "})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommListDelete(Routemap, Permit, Set, CommList, Delete string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set comm-list " + CommList + " " + "delete"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunity(Routemap, Permit, Set, Community string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunityNone(Routemap, Permit, Set, Community, None string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community none"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunityInternet(Routemap, Permit, Set, Community, Internet string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "internet"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunityLocalAS(Routemap, Permit, Set, Community, LocalAS string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "local-AS"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunityNoAdvertise(Routemap, Permit, Set, Community, NoAdvertise string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-advertise"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetCommunityNoExport(Routemap, Permit, Set, Community, NoExport string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set community " + Community + " " + "no-export"})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetDampening(Routemap, Permit, Set, Dampening string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set dampening " + Dampening})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetExtCommunityRT(Routemap, Permit, Set, ExtCommunity, RT string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity rt " + RT})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetExtCommunitySOO(Routemap, Permit, Set, ExtCommunity, SOO string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set extcommunity soo " + SOO})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetIPNexthop(Routemap, Permit, Set, IP, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ip nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetIP6Nexthop(Routemap, Permit, Set, IP6, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set ipv6 nexthop " + Nexthop})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetLevel(Routemap, Permit, Set, Level string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set level " + Level})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetLocalPreference(Routemap, Permit, Set, LocalPreference string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set local-preference " + LocalPreference})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetMetric(Routemap, Permit, Set, Metric string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric " + Metric})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetMetricType(Routemap, Permit, Set, MetricType string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set metric-type " + MetricType})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetOrigin(Routemap, Permit, Set, Origin string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set origin " + Origin})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetOriginatorID(Routemap, Permit, Set, OriginatorID string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set originator-id " + OriginatorID})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetTag(Routemap, Permit, Set, Tag string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set tag " + Tag})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) RoutemapPermitNoSetWeight(Routemap, Permit, Set, Weight string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "no set weight " + Weight})
	res = append(res, &command.Command{Mode: "config-route-map", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoRoutemap(Routemap string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoRoutemapDeny(Routemap, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoRoutemapPermit(Routemap, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no route-map " + Routemap + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSequenceNumber(IP, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDescription(IP, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDeny(IP, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDescription(IP, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDenyGe(IP, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDenyGeLe(IP, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDenyLe(IP, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListDenyLeGe(IP, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqDeny(IP, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqDenyGe(IP, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqDenyGeLe(IP, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqDenyLe(IP, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqDenyLeGe(IP, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Prefix List Permit

func (v8 V8500) IPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListPermit(IP, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListPermitGe(IP, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListPermitGeLe(IP, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListPermitLe(IP, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListPermitLeGe(IP, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqPermit(IP, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqPermitGe(IP, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqPermitGeLe(IP, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqPermitLe(IP, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPPrefixListSeqPermitLeGe(IP, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

//IPv6 Prefix List

func (v8 V8500) IP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSequenceNumber(IP6, PrefixList, SequenceNumber string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list sequence-number"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDeny(IP6, PrefixList, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDescription(IP6, PrefixList, Description string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " description " + Description})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDenyGe(IP6, PrefixList, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDenyGeLe(IP6, PrefixList, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDenyLe(IP6, PrefixList, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListDenyLeGe(IP6, PrefixList, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqDeny(IP6, PrefixList, Seq, Deny string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqDenyGe(IP6, PrefixList, Seq, Deny, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqDenyGeLe(IP6, PrefixList, Seq, Deny, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqDenyLe(IP6, PrefixList, Seq, Deny, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqDenyLeGe(IP6, PrefixList, Seq, Deny, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " deny " + Deny + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

///Prefix List Permit

func (v8 V8500) IP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListPermit(IP6, PrefixList, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListPermitGe(IP6, PrefixList, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListPermitGeLe(IP6, PrefixList, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListPermitLe(IP6, PrefixList, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListPermitLeGe(IP6, PrefixList, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqPermit(IP6, PrefixList, Seq, Permit string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqPermitGe(IP6, PrefixList, Seq, Permit, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqPermitGeLe(IP6, PrefixList, Seq, Permit, Ge, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " ge " + Ge + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqPermitLe(IP6, PrefixList, Seq, Permit, Le string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6PrefixListSeqPermitLeGe(IP6, PrefixList, Seq, Permit, Le, Ge string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 prefix-list " + PrefixList + " + seq " + Seq + " permit " + Permit + " le " + Le + " ge " + Ge})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) ClearIPRoute(Clear, IP, Route string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "clear ip route"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) ClearIPRouteKernel(Clear, IP, Route, Kernel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "clear ip route kernel"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRoutePrefixNexthopSrc(IP, Route, Prefix, Nexthop, Src string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route " + Prefix + " " + Nexthop + " src " + Src})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route default " + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route default " + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ip route vrf " + VRF + " " + Prefix + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRoutePrefixNexthop(IP, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRoutePrefixNexthopDistance(IP, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRouteDefaultNexthop(IP, Route, Default, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route default " + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRouteDefaultNexthopDistance(IP, Route, Default, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route default " + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRouteVRFPrefixNexthopOIF(IP, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIPRouteVRFPrefixOIF(IP, Route, VRF, Prefix, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ip route vrf " + VRF + " " + Prefix + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixNexthopDistance(IP6, Route, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixNexthopOIFDistance(IP6, Route, Prefix, Nexthop, OIF, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " " + Nexthop + " " + OIF + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RoutePrefixTunnelDistance(IP6, Route, Prefix, Tunnel, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route " + Prefix + " tunnel " + Tunnel + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixNexthopDistance(IP6, Route, VRF, Prefix, Nexthop, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixNexthopOIFDistance(IP6, Route, VRF, Prefix, Nexthop, OIF, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) IP6RouteVRFPrefixTunnelDistance(IP6, Route, VRF, Prefix, Tunnel, Distance string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel + " " + Distance})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RoutePrefix(IP6, Route, Prefix string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RoutePrefixNexthop(IP6, Route, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RoutePrefixNexthopOIF(IP6, Route, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RoutePrefixTunnel(IP6, Route, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RouteVRFPrefix(IP6, Route, VRF, Prefix string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoIP6RouteVRFPrefixNexthop(IP6, Route, VRF, Prefix, Nexthop string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v8 V8500) NoIP6RouteVRFPrefixNexthopOIF(IP6, Route, VRF, Prefix, Nexthop, OIF string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " " + Nexthop + " " + OIF})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res

}

func (v8 V8500) NoIP6RouteVRFPrefixTunnel(IP6, Route, VRF, Prefix, Tunnel string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "no ipv6 route vrf " + VRF + " " + Prefix + " tunnel " + Tunnel})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceIP(Interface, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ip address " + IP})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceIP2(Interface, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ip address " + IP2 + " secondary"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoInterfaceIP(Interface, IP string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ip address " + IP})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoInterfaceIP2(Interface, IP2 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ip address " + IP2 + " secondary"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceShutdown(Interface, Shutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "shutdown"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceNoShutdown(Interface, NoShutdown string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no shutdown"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 enable "})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoInterfaceIP6Enable(Interface, IP6, Enable string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 enable "})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceIP6(Interface, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 address " + IP6})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) InterfaceIP6LL(Interface, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "ipv6 address link-local " + IP6LL})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoInterfaceIP6(Interface, IP6 string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 address " + IP6})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}

func (v8 V8500) NoInterfaceIP6LL(Interface, IP6LL string) []*command.Command {
	res := make([]*command.Command, 0, 1)
	res = append(res, &command.Command{Mode: "normal", CMD: "configure terminal"})
	res = append(res, &command.Command{Mode: "config", CMD: "interface " + Interface})
	res = append(res, &command.Command{Mode: "config-if", CMD: "exit"})
	res = append(res, &command.Command{Mode: "config-if", CMD: "no ipv6 address link-local " + IP6LL})
	res = append(res, &command.Command{Mode: "config", CMD: "exit"})
	return res
}
