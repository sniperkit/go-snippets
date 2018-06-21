package dsl
import ( 
 "command"
 )
type Instruction struct {
Name string
Switch Switch
}
func PortSlotTypeEnable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Enable"]; !ok {
      return nil 
   }
      return sw.PortSlotTypeEnable(in["Port"], in["Slot"], in["Type"], in["Enable"])
}

func PortSlotTypeDisable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Disable"]; !ok {
      return nil 
   }
      return sw.PortSlotTypeDisable(in["Port"], in["Slot"], in["Type"], in["Disable"])
}

func PortSlotTypeSpeed(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Speed"]; !ok {
      return nil 
   }
      return sw.PortSlotTypeSpeed(in["Port"], in["Slot"], in["Type"], in["Speed"])
}

func PortSlotTypePvid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Pvid"]; !ok {
      return nil 
   }
      return sw.PortSlotTypePvid(in["Port"], in["Slot"], in["Type"], in["Pvid"])
}

func VLAN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 1 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
      return sw.VLAN(in["VLAN"])
}

func NoVLAN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
      return sw.NoVLAN(in["VLAN"])
}

func VLANAddTypeSlotPort(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Add"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
      return sw.VLANAddTypeSlotPort(in["VLAN"], in["Add"], in["Type"], in["Slot"], in["Port"])
}

func VLANAddTTypeSlotPort(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["AddT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
      return sw.VLANAddTTypeSlotPort(in["VLAN"], in["AddT"], in["Type"], in["Slot"], in["Port"])
}

func VLANDelTypeSlotPort(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Del"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
      return sw.VLANDelTypeSlotPort(in["VLAN"], in["Del"], in["Type"], in["Slot"], in["Port"])
}

func VLANDelTTypeSlotPort(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["DelT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
      return sw.VLANDelTTypeSlotPort(in["VLAN"], in["DelT"], in["Type"], in["Slot"], in["Port"])
}

func VLANShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Shutdown"]; !ok {
      return nil 
   }
      return sw.VLANShutdown(in["VLAN"], in["Shutdown"])
}

func VLANNoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.VLANNoShutdown(in["VLAN"], in["NoShutdown"])
}

func VLANIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.VLANIP(in["VLAN"], in["IP"])
}

func NoVLANIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.NoVLANIP(in["VLAN"], in["IP"])
}

func VLANIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.VLANIP2(in["VLAN"], in["IP2"])
}

func NoVLANIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.NoVLANIP2(in["VLAN"], in["IP2"])
}

func NoInterfaceTypeIfname(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Ifname"]; !ok {
      return nil 
   }
      return sw.NoInterfaceTypeIfname(in["Interface"], in["Type"], in["Ifname"])
}

func VLANAddTypeSlotPortIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Add"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.VLANAddTypeSlotPortIP(in["VLAN"], in["Add"], in["Type"], in["Slot"], in["Port"], in["IP"])
}

func VLANAddTypeSlotPortIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Add"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.VLANAddTypeSlotPortIP2(in["VLAN"], in["Add"], in["Type"], in["Slot"], in["Port"], in["IP2"])
}

func VLANAddTypeSlotPortIPNoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Add"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.VLANAddTypeSlotPortIPNoShutdown(in["VLAN"], in["Add"], in["Type"], in["Slot"], in["Port"], in["IP"], in["NoShutdown"])
}

func VLANAddTypeSlotPortIP2NoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["Add"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.VLANAddTypeSlotPortIP2NoShutdown(in["VLAN"], in["Add"], in["Type"], in["Slot"], in["Port"], in["IP2"], in["NoShutdown"])
}

func VLANAddTTypeSlotPortIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["AddT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.VLANAddTTypeSlotPortIP(in["VLAN"], in["AddT"], in["Type"], in["Slot"], in["Port"], in["IP"])
}

func VLANAddTTypeSlotPortIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["AddT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.VLANAddTTypeSlotPortIP2(in["VLAN"], in["AddT"], in["Type"], in["Slot"], in["Port"], in["IP2"])
}

func VLANAddTTypeSlotPortIPNoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["AddT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.VLANAddTTypeSlotPortIPNoShutdown(in["VLAN"], in["AddT"], in["Type"], in["Slot"], in["Port"], in["IP"], in["NoShutdown"])
}

func VLANAddTTypeSlotPortIP2NoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["AddT"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Slot"]; !ok {
      return nil 
   }
if _, ok := in["Port"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.VLANAddTTypeSlotPortIP2NoShutdown(in["VLAN"], in["AddT"], in["Type"], in["Slot"], in["Port"], in["IP2"], in["NoShutdown"])
}

func VLANIP6Enable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Enable"]; !ok {
      return nil 
   }
      return sw.VLANIP6Enable(in["VLAN"], in["IP6"], in["Enable"])
}

func NoVLANIP6Enable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Enable"]; !ok {
      return nil 
   }
      return sw.NoVLANIP6Enable(in["VLAN"], in["IP6"], in["Enable"])
}

func VLANIP6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
      return sw.VLANIP6(in["VLAN"], in["IP6"])
}

func NoVLANIP6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
      return sw.NoVLANIP6(in["VLAN"], in["IP6"])
}

func VLANIP6LL(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6LL"]; !ok {
      return nil 
   }
      return sw.VLANIP6LL(in["VLAN"], in["IP6LL"])
}

func VLANIP6LLIP6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6LL"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
      return sw.VLANIP6LLIP6(in["VLAN"], in["IP6LL"], in["IP6"])
}

func NoVLANIP6LL(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["VLAN"]; !ok {
      return nil 
   }
if _, ok := in["IP6LL"]; !ok {
      return nil 
   }
      return sw.NoVLANIP6LL(in["VLAN"], in["IP6LL"])
}

func OSPF6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 1 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
      return sw.OSPF6(in["OSPF6"])
}

func NoOSPF6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
      return sw.NoOSPF6(in["OSPF6"])
}

func OSPF6Rid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Rid"]; !ok {
      return nil 
   }
      return sw.OSPF6Rid(in["OSPF6"], in["Rid"])
}

func NoOSPF6Rid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Rid"]; !ok {
      return nil 
   }
      return sw.NoOSPF6Rid(in["OSPF6"], in["Rid"])
}

func OSPF6InterfaceArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceArea(in["OSPF6"], in["Interface"], in["Area"])
}

func NoOSPF6InterfaceArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceArea(in["OSPF6"], in["Interface"], in["Area"])
}

func OSPF6InterfaceCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceCost(in["OSPF6"], in["Interface"], in["Cost"])
}

func OSPF6InterfaceDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceDeadInterval(in["OSPF6"], in["Interface"], in["DeadInterval"])
}

func OSPF6InterfaceHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceHelloInterval(in["OSPF6"], in["Interface"], in["HelloInterval"])
}

func OSPF6InterfaceRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceRetransmitInterval(in["OSPF6"], in["Interface"], in["RetransmitInterval"])
}

func OSPF6InterfaceTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceTransmitDelay(in["OSPF6"], in["Interface"], in["TransmitDelay"])
}

func OSPF6InterfacePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfacePriority(in["OSPF6"], in["Interface"], in["Priority"])
}

func OSPF6InterfaceNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.OSPF6InterfaceNetworktype(in["OSPF6"], in["Interface"], in["Networktype"])
}

func NoOSPF6InterfaceCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceCost(in["OSPF6"], in["Interface"], in["Cost"])
}

func NoOSPF6InterfaceDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceDeadInterval(in["OSPF6"], in["Interface"], in["DeadInterval"])
}

func NoOSPF6InterfaceHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceHelloInterval(in["OSPF6"], in["Interface"], in["HelloInterval"])
}

func NoOSPF6InterfaceRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceRetransmitInterval(in["OSPF6"], in["Interface"], in["RetransmitInterval"])
}

func NoOSPF6InterfaceTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceTransmitDelay(in["OSPF6"], in["Interface"], in["TransmitDelay"])
}

func NoOSPF6InterfacePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfacePriority(in["OSPF6"], in["Interface"], in["Priority"])
}

func NoOSPF6InterfaceNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6InterfaceNetworktype(in["OSPF6"], in["Interface"], in["Networktype"])
}

func OSPF6IfTypeNameArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameArea(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Area"])
}

func NoOSPF6IfTypeNameArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameArea(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Area"])
}

func OSPF6IfTypeNameCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameCost(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Cost"])
}

func OSPF6IfTypeNameDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameDeadInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["DeadInterval"])
}

func OSPF6IfTypeNameHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameHelloInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["HelloInterval"])
}

func OSPF6IfTypeNameRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameRetransmitInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["RetransmitInterval"])
}

func OSPF6IfTypeNameTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameTransmitDelay(in["OSPF6"], in["If"], in["Type"], in["Name"], in["TransmitDelay"])
}

func OSPF6IfTypeNamePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNamePriority(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Priority"])
}

func OSPF6IfTypeNameNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.OSPF6IfTypeNameNetworktype(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Networktype"])
}

func NoOSPF6IfTypeNameCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameCost(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Cost"])
}

func NoOSPF6IfTypeNameDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameDeadInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["DeadInterval"])
}

func NoOSPF6IfTypeNameHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameHelloInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["HelloInterval"])
}

func NoOSPF6IfTypeNameRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameRetransmitInterval(in["OSPF6"], in["If"], in["Type"], in["Name"], in["RetransmitInterval"])
}

func NoOSPF6IfTypeNameTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameTransmitDelay(in["OSPF6"], in["If"], in["Type"], in["Name"], in["TransmitDelay"])
}

func NoOSPF6IfTypeNamePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNamePriority(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Priority"])
}

func NoOSPF6IfTypeNameNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["If"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6IfTypeNameNetworktype(in["OSPF6"], in["If"], in["Type"], in["Name"], in["Networktype"])
}

func OSPF6ReferenceBandwidth(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["ReferenceBandwidth"]; !ok {
      return nil 
   }
      return sw.OSPF6ReferenceBandwidth(in["OSPF6"], in["ReferenceBandwidth"])
}

func NoOSPF6ReferenceBandwidth(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["ReferenceBandwidth"]; !ok {
      return nil 
   }
      return sw.NoOSPF6ReferenceBandwidth(in["OSPF6"], in["ReferenceBandwidth"])
}

func OSPF6DefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginate(in["OSPF6"], in["DefaultOriginate"])
}

func OSPF6DefaultOriginateRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Routemap"])
}

func OSPF6DefaultOriginateMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetric(in["OSPF6"], in["DefaultOriginate"], in["Metric"])
}

func OSPF6DefaultOriginateMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Metrictype"])
}

func OSPF6DefaultOriginateMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetricMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Metrictype"])
}

func OSPF6DefaultOriginateMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetricRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Routemap"])
}

func OSPF6DefaultOriginateMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metrictype"], in["Routemap"])
}

func OSPF6DefaultOriginateMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateMetricMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPF6DefaultOriginateAlways(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlways(in["OSPF6"], in["DefaultOriginate"], in["Always"])
}

func OSPF6DefaultOriginateAlwaysRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Routemap"])
}

func OSPF6DefaultOriginateAlwaysMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetric(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"])
}

func OSPF6DefaultOriginateAlwaysMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metrictype"])
}

func OSPF6DefaultOriginateAlwaysMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetricMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"])
}

func OSPF6DefaultOriginateAlwaysMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetricRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Routemap"])
}

func OSPF6DefaultOriginateAlwaysMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metrictype"], in["Routemap"])
}

func OSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPF6DefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginate(in["OSPF6"], in["DefaultOriginate"])
}

func NoOSPF6DefaultOriginateRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Routemap"])
}

func NoOSPF6DefaultOriginateMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetric(in["OSPF6"], in["DefaultOriginate"], in["Metric"])
}

func NoOSPF6DefaultOriginateMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Metrictype"])
}

func NoOSPF6DefaultOriginateMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetricMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Metrictype"])
}

func NoOSPF6DefaultOriginateMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetricRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Routemap"])
}

func NoOSPF6DefaultOriginateMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metrictype"], in["Routemap"])
}

func NoOSPF6DefaultOriginateMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateMetricMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPF6DefaultOriginateAlways(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlways(in["OSPF6"], in["DefaultOriginate"], in["Always"])
}

func NoOSPF6DefaultOriginateAlwaysRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Routemap"])
}

func NoOSPF6DefaultOriginateAlwaysMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetric(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"])
}

func NoOSPF6DefaultOriginateAlwaysMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metrictype"])
}

func NoOSPF6DefaultOriginateAlwaysMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetricMetrictype(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"])
}

func NoOSPF6DefaultOriginateAlwaysMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetricRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Routemap"])
}

func NoOSPF6DefaultOriginateAlwaysMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metrictype"], in["Routemap"])
}

func NoOSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultOriginateAlwaysMetricMetrictypeRoutemap(in["OSPF6"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPF6Redistribute(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
      return sw.OSPF6Redistribute(in["OSPF6"], in["Redistribute"])
}

func OSPF6RedistributeMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeMetric(in["OSPF6"], in["Redistribute"], in["Metric"])
}

func OSPF6RedistributeMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeMetrictype(in["OSPF6"], in["Redistribute"], in["Metrictype"])
}

func OSPF6RedistributeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeRoutemap(in["OSPF6"], in["Redistribute"], in["Routemap"])
}

func OSPF6RedistributeMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeMetricMetrictype(in["OSPF6"], in["Redistribute"], in["Metric"], in["Metrictype"])
}

func OSPF6RedistributeMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeMetricRoutemap(in["OSPF6"], in["Redistribute"], in["Metric"], in["Routemap"])
}

func OSPF6RedistributeMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPF6RedistributeMetricMetrictypeRoutemap(in["OSPF6"], in["Redistribute"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPF6Redistribute(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
      return sw.NoOSPF6Redistribute(in["OSPF6"], in["Redistribute"])
}

func NoOSPF6RedistributeMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeMetric(in["OSPF6"], in["Redistribute"], in["Metric"])
}

func NoOSPF6RedistributeMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeMetrictype(in["OSPF6"], in["Redistribute"], in["Metrictype"])
}

func NoOSPF6RedistributeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeRoutemap(in["OSPF6"], in["Redistribute"], in["Routemap"])
}

func NoOSPF6RedistributeMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeMetricMetrictype(in["OSPF6"], in["Redistribute"], in["Metric"], in["Metrictype"])
}

func NoOSPF6RedistributeMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeMetricRoutemap(in["OSPF6"], in["Redistribute"], in["Metric"], in["Routemap"])
}

func NoOSPF6RedistributeMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPF6RedistributeMetricMetrictypeRoutemap(in["OSPF6"], in["Redistribute"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPF6Summary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
      return sw.OSPF6Summary(in["OSPF6"], in["Summary"])
}

func OSPF6SummaryNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.OSPF6SummaryNoAdvertise(in["OSPF6"], in["Summary"], in["NoAdvertise"])
}

func NoOSPF6Summary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
      return sw.NoOSPF6Summary(in["OSPF6"], in["Summary"])
}

func NoOSPF6SummaryNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.NoOSPF6SummaryNoAdvertise(in["OSPF6"], in["Summary"], in["NoAdvertise"])
}

func OSPF6DefaultMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultMetric"]; !ok {
      return nil 
   }
      return sw.OSPF6DefaultMetric(in["OSPF6"], in["DefaultMetric"])
}

func NoOSPF6DefaultMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["DefaultMetric"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DefaultMetric(in["OSPF6"], in["DefaultMetric"])
}

func OSPF6Passive(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Passive"]; !ok {
      return nil 
   }
      return sw.OSPF6Passive(in["OSPF6"], in["Passive"])
}

func NoOSPF6Passive(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Passive"]; !ok {
      return nil 
   }
      return sw.NoOSPF6Passive(in["OSPF6"], in["Passive"])
}

func OSPF6AdminDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["AdminDistance"]; !ok {
      return nil 
   }
      return sw.OSPF6AdminDistance(in["OSPF6"], in["AdminDistance"])
}

func NoOSPF6AdminDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["AdminDistance"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AdminDistance(in["OSPF6"], in["AdminDistance"])
}

func OSPF6DistanceExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceExternal(in["OSPF6"], in["Distance"], in["External"])
}

func OSPF6DistanceInter(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceInter(in["OSPF6"], in["Distance"], in["Inter"])
}

func OSPF6DistanceIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceIntra(in["OSPF6"], in["Distance"], in["Intra"])
}

func OSPF6DistanceInterIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceInterIntra(in["OSPF6"], in["Distance"], in["Inter"], in["Intra"])
}

func OSPF6DistanceInterExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceInterExternal(in["OSPF6"], in["Distance"], in["Inter"], in["External"])
}

func OSPF6DistanceInterIntraExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPF6DistanceInterIntraExternal(in["OSPF6"], in["Distance"], in["Inter"], in["Intra"], in["External"])
}

func NoOSPF6DistanceExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceExternal(in["OSPF6"], in["Distance"], in["External"])
}

func NoOSPF6DistanceInter(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceInter(in["OSPF6"], in["Distance"], in["Inter"])
}

func NoOSPF6DistanceIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceIntra(in["OSPF6"], in["Distance"], in["Intra"])
}

func NoOSPF6DistanceInterIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceInterIntra(in["OSPF6"], in["Distance"], in["Inter"], in["Intra"])
}

func NoOSPF6DistanceInterExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceInterExternal(in["OSPF6"], in["Distance"], in["Inter"], in["External"])
}

func NoOSPF6DistanceInterIntraExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistanceInterIntraExternal(in["OSPF6"], in["Distance"], in["Inter"], in["Intra"], in["External"])
}

func OSPF6DistributelistIN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["IN"]; !ok {
      return nil 
   }
      return sw.OSPF6DistributelistIN(in["OSPF6"], in["Distributelist"], in["IN"])
}

func OSPF6DistributelistOUT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["OUT"]; !ok {
      return nil 
   }
      return sw.OSPF6DistributelistOUT(in["OSPF6"], in["Distributelist"], in["OUT"])
}

func NoOSPF6DistributelistIN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["IN"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistributelistIN(in["OSPF6"], in["Distributelist"], in["IN"])
}

func NoOSPF6DistributelistOUT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["OUT"]; !ok {
      return nil 
   }
      return sw.NoOSPF6DistributelistOUT(in["OSPF6"], in["Distributelist"], in["OUT"])
}

func OSPF6AreaDefaultCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["DefaultCost"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaDefaultCost(in["OSPF6"], in["Area"], in["DefaultCost"])
}

func NoOSPF6AreaDefaultCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["DefaultCost"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaDefaultCost(in["OSPF6"], in["Area"], in["DefaultCost"])
}

func OSPF6AreaNSSA(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaNSSA(in["OSPF6"], in["Area"], in["NSSA"])
}

func OSPF6AreaNSSADefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaNSSADefaultOriginate(in["OSPF6"], in["Area"], in["NSSA"], in["DefaultOriginate"])
}

func OSPF6AreaNSSANoRedistribution(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoRedistribution"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaNSSANoRedistribution(in["OSPF6"], in["Area"], in["NSSA"], in["NoRedistribution"])
}

func OSPF6AreaNSSANoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaNSSANoSummary(in["OSPF6"], in["Area"], in["NSSA"], in["NoSummary"])
}

func OSPF6AreaNSSAStabilityInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["StabilityInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaNSSAStabilityInterval(in["OSPF6"], in["Area"], in["NSSA"], in["StabilityInterval"])
}

func OSPF6AreaTranslatorrole(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Translatorrole"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaTranslatorrole(in["OSPF6"], in["Area"], in["Translatorrole"])
}

func NoOSPF6AreaNSSA(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaNSSA(in["OSPF6"], in["Area"], in["NSSA"])
}

func NoOSPF6AreaNSSADefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaNSSADefaultOriginate(in["OSPF6"], in["Area"], in["NSSA"], in["DefaultOriginate"])
}

func NoOSPF6AreaNSSANoRedistribution(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoRedistribution"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaNSSANoRedistribution(in["OSPF6"], in["Area"], in["NSSA"], in["NoRedistribution"])
}

func NoOSPF6AreaNSSANoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaNSSANoSummary(in["OSPF6"], in["Area"], in["NSSA"], in["NoSummary"])
}

func NoOSPF6AreaNSSAStabilityInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["StabilityInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaNSSAStabilityInterval(in["OSPF6"], in["Area"], in["NSSA"], in["StabilityInterval"])
}

func NoOSPF6AreaTranslatorrole(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Translatorrole"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaTranslatorrole(in["OSPF6"], in["Area"], in["Translatorrole"])
}

func OSPF6AreaStub(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaStub(in["OSPF6"], in["Area"], in["Stub"])
}

func OSPF6AreaStubNoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaStubNoSummary(in["OSPF6"], in["Area"], in["Stub"], in["NoSummary"])
}

func NoOSPF6AreaStub(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaStub(in["OSPF6"], in["Area"], in["Stub"])
}

func NoOSPF6AreaStubNoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaStubNoSummary(in["OSPF6"], in["Area"], in["Stub"], in["NoSummary"])
}

func OSPF6AreaRange(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaRange(in["OSPF6"], in["Area"], in["Range"])
}

func OSPF6AreaRangeAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["Advertise"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaRangeAdvertise(in["OSPF6"], in["Area"], in["Range"], in["Advertise"])
}

func OSPF6AreaRangeNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaRangeNoAdvertise(in["OSPF6"], in["Area"], in["Range"], in["NoAdvertise"])
}

func NoOSPF6AreaRange(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaRange(in["OSPF6"], in["Area"], in["Range"])
}

func NoOSPF6AreaRangeAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["Advertise"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaRangeAdvertise(in["OSPF6"], in["Area"], in["Range"], in["Advertise"])
}

func NoOSPF6AreaRangeNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaRangeNoAdvertise(in["OSPF6"], in["Area"], in["Range"], in["NoAdvertise"])
}

func OSPF6AreaVirtuallink(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallink(in["OSPF6"], in["Area"], in["Virtuallink"])
}

func OSPF6AreaVirtuallinkDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallinkDeadInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["DeadInterval"])
}

func OSPF6AreaVirtuallinkHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallinkHelloInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["HelloInterval"])
}

func OSPF6AreaVirtuallinkInstanceid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["Instanceid"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallinkInstanceid(in["OSPF6"], in["Area"], in["Virtuallink"], in["Instanceid"])
}

func OSPF6AreaVirtuallinkRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallinkRetransmitInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["RetransmitInterval"])
}

func OSPF6AreaVirtuallinkTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.OSPF6AreaVirtuallinkTransmitDelay(in["OSPF6"], in["Area"], in["Virtuallink"], in["TransmitDelay"])
}

func NoOSPF6AreaVirtuallink(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallink(in["OSPF6"], in["Area"], in["Virtuallink"])
}

func NoOSPF6AreaVirtuallinkDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallinkDeadInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["DeadInterval"])
}

func NoOSPF6AreaVirtuallinkHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallinkHelloInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["HelloInterval"])
}

func NoOSPF6AreaVirtuallinkInstanceid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["Instanceid"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallinkInstanceid(in["OSPF6"], in["Area"], in["Virtuallink"], in["Instanceid"])
}

func NoOSPF6AreaVirtuallinkRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallinkRetransmitInterval(in["OSPF6"], in["Area"], in["Virtuallink"], in["RetransmitInterval"])
}

func NoOSPF6AreaVirtuallinkTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF6"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.NoOSPF6AreaVirtuallinkTransmitDelay(in["OSPF6"], in["Area"], in["Virtuallink"], in["TransmitDelay"])
}

func OSPF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 1 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
      return sw.OSPF(in["OSPF"])
}

func NoOSPF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
      return sw.NoOSPF(in["OSPF"])
}

func OSPFRid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Rid"]; !ok {
      return nil 
   }
      return sw.OSPFRid(in["OSPF"], in["Rid"])
}

func NoOSPFRid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Rid"]; !ok {
      return nil 
   }
      return sw.NoOSPFRid(in["OSPF"], in["Rid"])
}

func OSPFNetworkArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Network"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.OSPFNetworkArea(in["OSPF"], in["Network"], in["Area"])
}

func NoOSPFNetworkArea(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Network"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
      return sw.NoOSPFNetworkArea(in["OSPF"], in["Network"], in["Area"])
}

func OSPFInterfaceCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceCost(in["OSPF"], in["Interface"], in["Cost"])
}

func OSPFInterfaceDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceDeadInterval(in["OSPF"], in["Interface"], in["DeadInterval"])
}

func OSPFInterfaceHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceHelloInterval(in["OSPF"], in["Interface"], in["HelloInterval"])
}

func OSPFInterfaceRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceRetransmitInterval(in["OSPF"], in["Interface"], in["RetransmitInterval"])
}

func OSPFInterfaceTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceTransmitDelay(in["OSPF"], in["Interface"], in["TransmitDelay"])
}

func OSPFInterfacePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.OSPFInterfacePriority(in["OSPF"], in["Interface"], in["Priority"])
}

func OSPFInterfaceNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.OSPFInterfaceNetworktype(in["OSPF"], in["Interface"], in["Networktype"])
}

func NoOSPFInterfaceCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Cost"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceCost(in["OSPF"], in["Interface"], in["Cost"])
}

func NoOSPFInterfaceDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceDeadInterval(in["OSPF"], in["Interface"], in["DeadInterval"])
}

func NoOSPFInterfaceHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceHelloInterval(in["OSPF"], in["Interface"], in["HelloInterval"])
}

func NoOSPFInterfaceRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceRetransmitInterval(in["OSPF"], in["Interface"], in["RetransmitInterval"])
}

func NoOSPFInterfaceTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceTransmitDelay(in["OSPF"], in["Interface"], in["TransmitDelay"])
}

func NoOSPFInterfacePriority(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Priority"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfacePriority(in["OSPF"], in["Interface"], in["Priority"])
}

func NoOSPFInterfaceNetworktype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Networktype"]; !ok {
      return nil 
   }
      return sw.NoOSPFInterfaceNetworktype(in["OSPF"], in["Interface"], in["Networktype"])
}

func OSPFReferenceBandwidth(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["ReferenceBandwidth"]; !ok {
      return nil 
   }
      return sw.OSPFReferenceBandwidth(in["OSPF"], in["ReferenceBandwidth"])
}

func NoOSPFReferenceBandwidth(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["ReferenceBandwidth"]; !ok {
      return nil 
   }
      return sw.NoOSPFReferenceBandwidth(in["OSPF"], in["ReferenceBandwidth"])
}

func OSPFDefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginate(in["OSPF"], in["DefaultOriginate"])
}

func OSPFDefaultOriginateRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateRoutemap(in["OSPF"], in["DefaultOriginate"], in["Routemap"])
}

func OSPFDefaultOriginateMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetric(in["OSPF"], in["DefaultOriginate"], in["Metric"])
}

func OSPFDefaultOriginateMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetrictype(in["OSPF"], in["DefaultOriginate"], in["Metrictype"])
}

func OSPFDefaultOriginateMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetricMetrictype(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Metrictype"])
}

func OSPFDefaultOriginateMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetricRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Routemap"])
}

func OSPFDefaultOriginateMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metrictype"], in["Routemap"])
}

func OSPFDefaultOriginateMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateMetricMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPFDefaultOriginateAlways(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlways(in["OSPF"], in["DefaultOriginate"], in["Always"])
}

func OSPFDefaultOriginateAlwaysRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Routemap"])
}

func OSPFDefaultOriginateAlwaysMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetric(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"])
}

func OSPFDefaultOriginateAlwaysMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetrictype(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metrictype"])
}

func OSPFDefaultOriginateAlwaysMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetricMetrictype(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"])
}

func OSPFDefaultOriginateAlwaysMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetricRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Routemap"])
}

func OSPFDefaultOriginateAlwaysMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metrictype"], in["Routemap"])
}

func OSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPFDefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginate(in["OSPF"], in["DefaultOriginate"])
}

func NoOSPFDefaultOriginateRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateRoutemap(in["OSPF"], in["DefaultOriginate"], in["Routemap"])
}

func NoOSPFDefaultOriginateMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetric(in["OSPF"], in["DefaultOriginate"], in["Metric"])
}

func NoOSPFDefaultOriginateMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetrictype(in["OSPF"], in["DefaultOriginate"], in["Metrictype"])
}

func NoOSPFDefaultOriginateMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetricMetrictype(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Metrictype"])
}

func NoOSPFDefaultOriginateMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetricRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Routemap"])
}

func NoOSPFDefaultOriginateMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metrictype"], in["Routemap"])
}

func NoOSPFDefaultOriginateMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateMetricMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPFDefaultOriginateAlways(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlways(in["OSPF"], in["DefaultOriginate"], in["Always"])
}

func NoOSPFDefaultOriginateAlwaysRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Routemap"])
}

func NoOSPFDefaultOriginateAlwaysMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetric(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"])
}

func NoOSPFDefaultOriginateAlwaysMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetrictype(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metrictype"])
}

func NoOSPFDefaultOriginateAlwaysMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetricMetrictype(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"])
}

func NoOSPFDefaultOriginateAlwaysMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetricRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Routemap"])
}

func NoOSPFDefaultOriginateAlwaysMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metrictype"], in["Routemap"])
}

func NoOSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
if _, ok := in["Always"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultOriginateAlwaysMetricMetrictypeRoutemap(in["OSPF"], in["DefaultOriginate"], in["Always"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPFRedistribute(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
      return sw.OSPFRedistribute(in["OSPF"], in["Redistribute"])
}

func OSPFRedistributeMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeMetric(in["OSPF"], in["Redistribute"], in["Metric"])
}

func OSPFRedistributeMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeMetrictype(in["OSPF"], in["Redistribute"], in["Metrictype"])
}

func OSPFRedistributeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeRoutemap(in["OSPF"], in["Redistribute"], in["Routemap"])
}

func OSPFRedistributeMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeMetricMetrictype(in["OSPF"], in["Redistribute"], in["Metric"], in["Metrictype"])
}

func OSPFRedistributeMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeMetricRoutemap(in["OSPF"], in["Redistribute"], in["Metric"], in["Routemap"])
}

func OSPFRedistributeMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.OSPFRedistributeMetricMetrictypeRoutemap(in["OSPF"], in["Redistribute"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func NoOSPFRedistribute(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistribute(in["OSPF"], in["Redistribute"])
}

func NoOSPFRedistributeMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeMetric(in["OSPF"], in["Redistribute"], in["Metric"])
}

func NoOSPFRedistributeMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeMetrictype(in["OSPF"], in["Redistribute"], in["Metrictype"])
}

func NoOSPFRedistributeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeRoutemap(in["OSPF"], in["Redistribute"], in["Routemap"])
}

func NoOSPFRedistributeMetricMetrictype(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeMetricMetrictype(in["OSPF"], in["Redistribute"], in["Metric"], in["Metrictype"])
}

func NoOSPFRedistributeMetricRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeMetricRoutemap(in["OSPF"], in["Redistribute"], in["Metric"], in["Routemap"])
}

func NoOSPFRedistributeMetricMetrictypeRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Redistribute"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
if _, ok := in["Metrictype"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoOSPFRedistributeMetricMetrictypeRoutemap(in["OSPF"], in["Redistribute"], in["Metric"], in["Metrictype"], in["Routemap"])
}

func OSPFSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
      return sw.OSPFSummary(in["OSPF"], in["Summary"])
}

func OSPFSummaryNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.OSPFSummaryNoAdvertise(in["OSPF"], in["Summary"], in["NoAdvertise"])
}

func NoOSPFSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
      return sw.NoOSPFSummary(in["OSPF"], in["Summary"])
}

func NoOSPFSummaryNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Summary"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.NoOSPFSummaryNoAdvertise(in["OSPF"], in["Summary"], in["NoAdvertise"])
}

func OSPFDefaultMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultMetric"]; !ok {
      return nil 
   }
      return sw.OSPFDefaultMetric(in["OSPF"], in["DefaultMetric"])
}

func NoOSPFDefaultMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["DefaultMetric"]; !ok {
      return nil 
   }
      return sw.NoOSPFDefaultMetric(in["OSPF"], in["DefaultMetric"])
}

func OSPFPassive(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Passive"]; !ok {
      return nil 
   }
      return sw.OSPFPassive(in["OSPF"], in["Passive"])
}

func NoOSPFPassive(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Passive"]; !ok {
      return nil 
   }
      return sw.NoOSPFPassive(in["OSPF"], in["Passive"])
}

func OSPFAdminDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["AdminDistance"]; !ok {
      return nil 
   }
      return sw.OSPFAdminDistance(in["OSPF"], in["AdminDistance"])
}

func NoOSPFAdminDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["AdminDistance"]; !ok {
      return nil 
   }
      return sw.NoOSPFAdminDistance(in["OSPF"], in["AdminDistance"])
}

func OSPFDistanceExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceExternal(in["OSPF"], in["Distance"], in["External"])
}

func OSPFDistanceInter(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceInter(in["OSPF"], in["Distance"], in["Inter"])
}

func OSPFDistanceIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceIntra(in["OSPF"], in["Distance"], in["Intra"])
}

func OSPFDistanceInterIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceInterIntra(in["OSPF"], in["Distance"], in["Inter"], in["Intra"])
}

func OSPFDistanceInterExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceInterExternal(in["OSPF"], in["Distance"], in["Inter"], in["External"])
}

func OSPFDistanceInterIntraExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.OSPFDistanceInterIntraExternal(in["OSPF"], in["Distance"], in["Inter"], in["Intra"], in["External"])
}

func NoOSPFDistanceExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceExternal(in["OSPF"], in["Distance"], in["External"])
}

func NoOSPFDistanceInter(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceInter(in["OSPF"], in["Distance"], in["Inter"])
}

func NoOSPFDistanceIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceIntra(in["OSPF"], in["Distance"], in["Intra"])
}

func NoOSPFDistanceInterIntra(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceInterIntra(in["OSPF"], in["Distance"], in["Inter"], in["Intra"])
}

func NoOSPFDistanceInterExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceInterExternal(in["OSPF"], in["Distance"], in["Inter"], in["External"])
}

func NoOSPFDistanceInterIntraExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
if _, ok := in["Inter"]; !ok {
      return nil 
   }
if _, ok := in["Intra"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistanceInterIntraExternal(in["OSPF"], in["Distance"], in["Inter"], in["Intra"], in["External"])
}

func OSPFDistributelistIN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["IN"]; !ok {
      return nil 
   }
      return sw.OSPFDistributelistIN(in["OSPF"], in["Distributelist"], in["IN"])
}

func OSPFDistributelistOUT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["OUT"]; !ok {
      return nil 
   }
      return sw.OSPFDistributelistOUT(in["OSPF"], in["Distributelist"], in["OUT"])
}

func NoOSPFDistributelistIN(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["IN"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistributelistIN(in["OSPF"], in["Distributelist"], in["IN"])
}

func NoOSPFDistributelistOUT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Distributelist"]; !ok {
      return nil 
   }
if _, ok := in["OUT"]; !ok {
      return nil 
   }
      return sw.NoOSPFDistributelistOUT(in["OSPF"], in["Distributelist"], in["OUT"])
}

func OSPFAreaDefaultCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["DefaultCost"]; !ok {
      return nil 
   }
      return sw.OSPFAreaDefaultCost(in["OSPF"], in["Area"], in["DefaultCost"])
}

func NoOSPFAreaDefaultCost(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["DefaultCost"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaDefaultCost(in["OSPF"], in["Area"], in["DefaultCost"])
}

func OSPFAreaNSSA(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
      return sw.OSPFAreaNSSA(in["OSPF"], in["Area"], in["NSSA"])
}

func OSPFAreaNSSADefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.OSPFAreaNSSADefaultOriginate(in["OSPF"], in["Area"], in["NSSA"], in["DefaultOriginate"])
}

func OSPFAreaNSSANoRedistribution(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoRedistribution"]; !ok {
      return nil 
   }
      return sw.OSPFAreaNSSANoRedistribution(in["OSPF"], in["Area"], in["NSSA"], in["NoRedistribution"])
}

func OSPFAreaNSSANoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.OSPFAreaNSSANoSummary(in["OSPF"], in["Area"], in["NSSA"], in["NoSummary"])
}

func OSPFAreaNSSAStabilityInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["StabilityInterval"]; !ok {
      return nil 
   }
      return sw.OSPFAreaNSSAStabilityInterval(in["OSPF"], in["Area"], in["NSSA"], in["StabilityInterval"])
}

func OSPFAreaTranslatorrole(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Translatorrole"]; !ok {
      return nil 
   }
      return sw.OSPFAreaTranslatorrole(in["OSPF"], in["Area"], in["Translatorrole"])
}

func NoOSPFAreaNSSA(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaNSSA(in["OSPF"], in["Area"], in["NSSA"])
}

func NoOSPFAreaNSSADefaultOriginate(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["DefaultOriginate"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaNSSADefaultOriginate(in["OSPF"], in["Area"], in["NSSA"], in["DefaultOriginate"])
}

func NoOSPFAreaNSSANoRedistribution(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoRedistribution"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaNSSANoRedistribution(in["OSPF"], in["Area"], in["NSSA"], in["NoRedistribution"])
}

func NoOSPFAreaNSSANoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaNSSANoSummary(in["OSPF"], in["Area"], in["NSSA"], in["NoSummary"])
}

func NoOSPFAreaNSSAStabilityInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["NSSA"]; !ok {
      return nil 
   }
if _, ok := in["StabilityInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaNSSAStabilityInterval(in["OSPF"], in["Area"], in["NSSA"], in["StabilityInterval"])
}

func NoOSPFAreaTranslatorrole(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Translatorrole"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaTranslatorrole(in["OSPF"], in["Area"], in["Translatorrole"])
}

func OSPFAreaStub(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
      return sw.OSPFAreaStub(in["OSPF"], in["Area"], in["Stub"])
}

func OSPFAreaStubNoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.OSPFAreaStubNoSummary(in["OSPF"], in["Area"], in["Stub"], in["NoSummary"])
}

func NoOSPFAreaStub(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaStub(in["OSPF"], in["Area"], in["Stub"])
}

func NoOSPFAreaStubNoSummary(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Stub"]; !ok {
      return nil 
   }
if _, ok := in["NoSummary"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaStubNoSummary(in["OSPF"], in["Area"], in["Stub"], in["NoSummary"])
}

func OSPFAreaRange(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
      return sw.OSPFAreaRange(in["OSPF"], in["Area"], in["Range"])
}

func OSPFAreaRangeAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["Advertise"]; !ok {
      return nil 
   }
      return sw.OSPFAreaRangeAdvertise(in["OSPF"], in["Area"], in["Range"], in["Advertise"])
}

func OSPFAreaRangeNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.OSPFAreaRangeNoAdvertise(in["OSPF"], in["Area"], in["Range"], in["NoAdvertise"])
}

func NoOSPFAreaRange(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaRange(in["OSPF"], in["Area"], in["Range"])
}

func NoOSPFAreaRangeAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["Advertise"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaRangeAdvertise(in["OSPF"], in["Area"], in["Range"], in["Advertise"])
}

func NoOSPFAreaRangeNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Range"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaRangeNoAdvertise(in["OSPF"], in["Area"], in["Range"], in["NoAdvertise"])
}

func OSPFAreaVirtuallink(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallink(in["OSPF"], in["Area"], in["Virtuallink"])
}

func OSPFAreaVirtuallinkDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallinkDeadInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["DeadInterval"])
}

func OSPFAreaVirtuallinkHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallinkHelloInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["HelloInterval"])
}

func OSPFAreaVirtuallinkInstanceid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["Instanceid"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallinkInstanceid(in["OSPF"], in["Area"], in["Virtuallink"], in["Instanceid"])
}

func OSPFAreaVirtuallinkRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallinkRetransmitInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["RetransmitInterval"])
}

func OSPFAreaVirtuallinkTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.OSPFAreaVirtuallinkTransmitDelay(in["OSPF"], in["Area"], in["Virtuallink"], in["TransmitDelay"])
}

func NoOSPFAreaVirtuallink(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallink(in["OSPF"], in["Area"], in["Virtuallink"])
}

func NoOSPFAreaVirtuallinkDeadInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["DeadInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallinkDeadInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["DeadInterval"])
}

func NoOSPFAreaVirtuallinkHelloInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["HelloInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallinkHelloInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["HelloInterval"])
}

func NoOSPFAreaVirtuallinkInstanceid(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["Instanceid"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallinkInstanceid(in["OSPF"], in["Area"], in["Virtuallink"], in["Instanceid"])
}

func NoOSPFAreaVirtuallinkRetransmitInterval(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["RetransmitInterval"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallinkRetransmitInterval(in["OSPF"], in["Area"], in["Virtuallink"], in["RetransmitInterval"])
}

func NoOSPFAreaVirtuallinkTransmitDelay(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["OSPF"]; !ok {
      return nil 
   }
if _, ok := in["Area"]; !ok {
      return nil 
   }
if _, ok := in["Virtuallink"]; !ok {
      return nil 
   }
if _, ok := in["TransmitDelay"]; !ok {
      return nil 
   }
      return sw.NoOSPFAreaVirtuallinkTransmitDelay(in["OSPF"], in["Area"], in["Virtuallink"], in["TransmitDelay"])
}

func RoutemapDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.RoutemapDeny(in["Routemap"], in["Deny"])
}

func RoutemapDenyMatchASPath(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchASPath(in["Routemap"], in["Deny"], in["Match"], in["ASPath"])
}

func RoutemapDenyMatchCLNSAddress(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["CLNS"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchCLNSAddress(in["Routemap"], in["Deny"], in["Match"], in["CLNS"], in["Address"])
}

func RoutemapDenyMatchCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchCommunity(in["Routemap"], in["Deny"], in["Match"], in["Community"])
}

func RoutemapDenyMatchCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchCommunityExactMatch(in["Routemap"], in["Deny"], in["Match"], in["Community"], in["ExactMatch"])
}

func RoutemapDenyMatchExtCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchExtCommunity(in["Routemap"], in["Deny"], in["Match"], in["ExtCommunity"])
}

func RoutemapDenyMatchExtCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchExtCommunityExactMatch(in["Routemap"], in["Deny"], in["Match"], in["ExtCommunity"], in["ExactMatch"])
}

func RoutemapDenyMatchInterfaceTypeName(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchInterfaceTypeName(in["Routemap"], in["Deny"], in["Match"], in["Interface"], in["Type"], in["Name"])
}

func RoutemapDenyMatchIPAddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIPAddressAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Address"], in["AccessList"])
}

func RoutemapDenyMatchIPAddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIPAddressPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Address"], in["PrefixList"])
}

func RoutemapDenyMatchIPNexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIPNexthopAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Nexthop"], in["AccessList"])
}

func RoutemapDenyMatchIPNexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIPNexthopPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Nexthop"], in["PrefixList"])
}

func RoutemapDenyMatchIP6AddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIP6AddressAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Address"], in["AccessList"])
}

func RoutemapDenyMatchIP6AddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIP6AddressPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Address"], in["PrefixList"])
}

func RoutemapDenyMatchIP6NexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIP6NexthopAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Nexthop"], in["AccessList"])
}

func RoutemapDenyMatchIP6NexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchIP6NexthopPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Nexthop"], in["PrefixList"])
}

func RoutemapDenyMatchMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchMetric(in["Routemap"], in["Deny"], in["Match"], in["Metric"])
}

func RoutemapDenyMatchOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchOrigin(in["Routemap"], in["Deny"], in["Match"], in["Origin"])
}

func RoutemapDenyMatchRouteTypeExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["RouteType"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchRouteTypeExternal(in["Routemap"], in["Deny"], in["Match"], in["RouteType"], in["External"])
}

func RoutemapDenyMatchTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyMatchTag(in["Routemap"], in["Deny"], in["Match"], in["Tag"])
}

func RoutemapDenyNoMatchASPath(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchASPath(in["Routemap"], in["Deny"], in["Match"], in["ASPath"])
}

func RoutemapDenyNoMatchCLNSAddress(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["CLNS"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchCLNSAddress(in["Routemap"], in["Deny"], in["Match"], in["CLNS"], in["Address"])
}

func RoutemapDenyNoMatchCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchCommunity(in["Routemap"], in["Deny"], in["Match"], in["Community"])
}

func RoutemapDenyNoMatchCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchCommunityExactMatch(in["Routemap"], in["Deny"], in["Match"], in["Community"], in["ExactMatch"])
}

func RoutemapDenyNoMatchExtCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchExtCommunity(in["Routemap"], in["Deny"], in["Match"], in["ExtCommunity"])
}

func RoutemapDenyNoMatchExtCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchExtCommunityExactMatch(in["Routemap"], in["Deny"], in["Match"], in["ExtCommunity"], in["ExactMatch"])
}

func RoutemapDenyNoMatchInterfaceTypeName(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchInterfaceTypeName(in["Routemap"], in["Deny"], in["Match"], in["Interface"], in["Type"], in["Name"])
}

func RoutemapDenyNoMatchIPAddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIPAddressAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Address"], in["AccessList"])
}

func RoutemapDenyNoMatchIPAddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIPAddressPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Address"], in["PrefixList"])
}

func RoutemapDenyNoMatchIPNexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIPNexthopAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Nexthop"], in["AccessList"])
}

func RoutemapDenyNoMatchIPNexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIPNexthopPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP"], in["Nexthop"], in["PrefixList"])
}

func RoutemapDenyNoMatchIP6AddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIP6AddressAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Address"], in["AccessList"])
}

func RoutemapDenyNoMatchIP6AddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIP6AddressPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Address"], in["PrefixList"])
}

func RoutemapDenyNoMatchIP6NexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIP6NexthopAccessList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Nexthop"], in["AccessList"])
}

func RoutemapDenyNoMatchIP6NexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchIP6NexthopPrefixList(in["Routemap"], in["Deny"], in["Match"], in["IP6"], in["Nexthop"], in["PrefixList"])
}

func RoutemapDenyNoMatchMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchMetric(in["Routemap"], in["Deny"], in["Match"], in["Metric"])
}

func RoutemapDenyNoMatchOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchOrigin(in["Routemap"], in["Deny"], in["Match"], in["Origin"])
}

func RoutemapDenyNoMatchRouteTypeExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["RouteType"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchRouteTypeExternal(in["Routemap"], in["Deny"], in["Match"], in["RouteType"], in["External"])
}

func RoutemapDenyNoMatchTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoMatchTag(in["Routemap"], in["Deny"], in["Match"], in["Tag"])
}

func RoutemapDenySetAggregatorASIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetAggregatorASIP(in["Routemap"], in["Deny"], in["Set"], in["Aggregator"], in["AS"], in["IP"])
}

func RoutemapDenySetASPathPrepend(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
if _, ok := in["Prepend"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetASPathPrepend(in["Routemap"], in["Deny"], in["Set"], in["ASPath"], in["Prepend"])
}

func RoutemapDenySetAtomicAggregator(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["AtomicAggregator"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetAtomicAggregator(in["Routemap"], in["Deny"], in["Set"], in["AtomicAggregator"])
}

func RoutemapDenySetCommListDelete(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["CommList"]; !ok {
      return nil 
   }
if _, ok := in["Delete"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommListDelete(in["Routemap"], in["Deny"], in["Set"], in["CommList"], in["Delete"])
}

func RoutemapDenySetCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunity(in["Routemap"], in["Deny"], in["Set"], in["Community"])
}

func RoutemapDenySetCommunityNone(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["None"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunityNone(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["None"])
}

func RoutemapDenySetCommunityInternet(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["Internet"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunityInternet(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["Internet"])
}

func RoutemapDenySetCommunityLocalAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["LocalAS"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunityLocalAS(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["LocalAS"])
}

func RoutemapDenySetCommunityNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunityNoAdvertise(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["NoAdvertise"])
}

func RoutemapDenySetCommunityNoExport(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoExport"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetCommunityNoExport(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["NoExport"])
}

func RoutemapDenySetDampening(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Dampening"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetDampening(in["Routemap"], in["Deny"], in["Set"], in["Dampening"])
}

func RoutemapDenySetExtCommunityRT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["RT"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetExtCommunityRT(in["Routemap"], in["Deny"], in["Set"], in["ExtCommunity"], in["RT"])
}

func RoutemapDenySetExtCommunitySOO(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["SOO"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetExtCommunitySOO(in["Routemap"], in["Deny"], in["Set"], in["ExtCommunity"], in["SOO"])
}

func RoutemapDenySetIPNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetIPNexthop(in["Routemap"], in["Deny"], in["Set"], in["IP"], in["Nexthop"])
}

func RoutemapDenySetIP6Nexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetIP6Nexthop(in["Routemap"], in["Deny"], in["Set"], in["IP6"], in["Nexthop"])
}

func RoutemapDenySetLevel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Level"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetLevel(in["Routemap"], in["Deny"], in["Set"], in["Level"])
}

func RoutemapDenySetLocalPreference(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["LocalPreference"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetLocalPreference(in["Routemap"], in["Deny"], in["Set"], in["LocalPreference"])
}

func RoutemapDenySetMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetMetric(in["Routemap"], in["Deny"], in["Set"], in["Metric"])
}

func RoutemapDenySetMetricType(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["MetricType"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetMetricType(in["Routemap"], in["Deny"], in["Set"], in["MetricType"])
}

func RoutemapDenySetOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetOrigin(in["Routemap"], in["Deny"], in["Set"], in["Origin"])
}

func RoutemapDenySetOriginatorID(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["OriginatorID"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetOriginatorID(in["Routemap"], in["Deny"], in["Set"], in["OriginatorID"])
}

func RoutemapDenySetTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetTag(in["Routemap"], in["Deny"], in["Set"], in["Tag"])
}

func RoutemapDenySetWeight(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Weight"]; !ok {
      return nil 
   }
      return sw.RoutemapDenySetWeight(in["Routemap"], in["Deny"], in["Set"], in["Weight"])
}

func RoutemapDenyNoSetAggregatorAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetAggregatorAS(in["Routemap"], in["Deny"], in["Set"], in["Aggregator"], in["AS"])
}

func RoutemapDenyNoSetAggregatorASIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetAggregatorASIP(in["Routemap"], in["Deny"], in["Set"], in["Aggregator"], in["AS"], in["IP"])
}

func RoutemapDenyNoSetASPathPrepend(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
if _, ok := in["Prepend"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetASPathPrepend(in["Routemap"], in["Deny"], in["Set"], in["ASPath"], in["Prepend"])
}

func RoutemapDenyNoSetAtomicAggregator(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["AtomicAggregator"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetAtomicAggregator(in["Routemap"], in["Deny"], in["Set"], in["AtomicAggregator"])
}

func RoutemapDenyNoSetCommListDelete(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["CommList"]; !ok {
      return nil 
   }
if _, ok := in["Delete"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommListDelete(in["Routemap"], in["Deny"], in["Set"], in["CommList"], in["Delete"])
}

func RoutemapDenyNoSetCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunity(in["Routemap"], in["Deny"], in["Set"], in["Community"])
}

func RoutemapDenyNoSetCommunityNone(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["None"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunityNone(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["None"])
}

func RoutemapDenyNoSetCommunityInternet(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["Internet"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunityInternet(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["Internet"])
}

func RoutemapDenyNoSetCommunityLocalAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["LocalAS"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunityLocalAS(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["LocalAS"])
}

func RoutemapDenyNoSetCommunityNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunityNoAdvertise(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["NoAdvertise"])
}

func RoutemapDenyNoSetCommunityNoExport(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoExport"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetCommunityNoExport(in["Routemap"], in["Deny"], in["Set"], in["Community"], in["NoExport"])
}

func RoutemapDenyNoSetDampening(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Dampening"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetDampening(in["Routemap"], in["Deny"], in["Set"], in["Dampening"])
}

func RoutemapDenyNoSetExtCommunityRT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["RT"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetExtCommunityRT(in["Routemap"], in["Deny"], in["Set"], in["ExtCommunity"], in["RT"])
}

func RoutemapDenyNoSetExtCommunitySOO(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["SOO"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetExtCommunitySOO(in["Routemap"], in["Deny"], in["Set"], in["ExtCommunity"], in["SOO"])
}

func RoutemapDenyNoSetIPNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetIPNexthop(in["Routemap"], in["Deny"], in["Set"], in["IP"], in["Nexthop"])
}

func RoutemapDenyNoSetIP6Nexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetIP6Nexthop(in["Routemap"], in["Deny"], in["Set"], in["IP6"], in["Nexthop"])
}

func RoutemapDenyNoSetLevel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Level"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetLevel(in["Routemap"], in["Deny"], in["Set"], in["Level"])
}

func RoutemapDenyNoSetLocalPreference(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["LocalPreference"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetLocalPreference(in["Routemap"], in["Deny"], in["Set"], in["LocalPreference"])
}

func RoutemapDenyNoSetMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetMetric(in["Routemap"], in["Deny"], in["Set"], in["Metric"])
}

func RoutemapDenyNoSetMetricType(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["MetricType"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetMetricType(in["Routemap"], in["Deny"], in["Set"], in["MetricType"])
}

func RoutemapDenyNoSetOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetOrigin(in["Routemap"], in["Deny"], in["Set"], in["Origin"])
}

func RoutemapDenyNoSetOriginatorID(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["OriginatorID"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetOriginatorID(in["Routemap"], in["Deny"], in["Set"], in["OriginatorID"])
}

func RoutemapDenyNoSetTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetTag(in["Routemap"], in["Deny"], in["Set"], in["Tag"])
}

func RoutemapDenyNoSetWeight(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Weight"]; !ok {
      return nil 
   }
      return sw.RoutemapDenyNoSetWeight(in["Routemap"], in["Deny"], in["Set"], in["Weight"])
}

func RoutemapPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.RoutemapPermit(in["Routemap"], in["Permit"])
}

func RoutemapPermitMatchASPath(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchASPath(in["Routemap"], in["Permit"], in["Match"], in["ASPath"])
}

func RoutemapPermitMatchCLNSAddress(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["CLNS"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchCLNSAddress(in["Routemap"], in["Permit"], in["Match"], in["CLNS"], in["Address"])
}

func RoutemapPermitMatchCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchCommunity(in["Routemap"], in["Permit"], in["Match"], in["Community"])
}

func RoutemapPermitMatchCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchCommunityExactMatch(in["Routemap"], in["Permit"], in["Match"], in["Community"], in["ExactMatch"])
}

func RoutemapPermitMatchExtCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchExtCommunity(in["Routemap"], in["Permit"], in["Match"], in["ExtCommunity"])
}

func RoutemapPermitMatchExtCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchExtCommunityExactMatch(in["Routemap"], in["Permit"], in["Match"], in["ExtCommunity"], in["ExactMatch"])
}

func RoutemapPermitMatchInterfaceTypeName(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchInterfaceTypeName(in["Routemap"], in["Permit"], in["Match"], in["Interface"], in["Type"], in["Name"])
}

func RoutemapPermitMatchIPAddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIPAddressAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Address"], in["AccessList"])
}

func RoutemapPermitMatchIPAddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIPAddressPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Address"], in["PrefixList"])
}

func RoutemapPermitMatchIPNexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIPNexthopAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Nexthop"], in["AccessList"])
}

func RoutemapPermitMatchIPNexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIPNexthopPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Nexthop"], in["PrefixList"])
}

func RoutemapPermitMatchIP6AddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIP6AddressAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Address"], in["AccessList"])
}

func RoutemapPermitMatchIP6AddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIP6AddressPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Address"], in["PrefixList"])
}

func RoutemapPermitMatchIP6NexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIP6NexthopAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Nexthop"], in["AccessList"])
}

func RoutemapPermitMatchIP6NexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchIP6NexthopPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Nexthop"], in["PrefixList"])
}

func RoutemapPermitMatchMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchMetric(in["Routemap"], in["Permit"], in["Match"], in["Metric"])
}

func RoutemapPermitMatchOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchOrigin(in["Routemap"], in["Permit"], in["Match"], in["Origin"])
}

func RoutemapPermitMatchRouteTypeExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["RouteType"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchRouteTypeExternal(in["Routemap"], in["Permit"], in["Match"], in["RouteType"], in["External"])
}

func RoutemapPermitMatchTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitMatchTag(in["Routemap"], in["Permit"], in["Match"], in["Tag"])
}

func RoutemapPermitNoMatchASPath(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchASPath(in["Routemap"], in["Permit"], in["Match"], in["ASPath"])
}

func RoutemapPermitNoMatchCLNSAddress(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["CLNS"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchCLNSAddress(in["Routemap"], in["Permit"], in["Match"], in["CLNS"], in["Address"])
}

func RoutemapPermitNoMatchCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchCommunity(in["Routemap"], in["Permit"], in["Match"], in["Community"])
}

func RoutemapPermitNoMatchCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchCommunityExactMatch(in["Routemap"], in["Permit"], in["Match"], in["Community"], in["ExactMatch"])
}

func RoutemapPermitNoMatchExtCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchExtCommunity(in["Routemap"], in["Permit"], in["Match"], in["ExtCommunity"])
}

func RoutemapPermitNoMatchExtCommunityExactMatch(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["ExactMatch"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchExtCommunityExactMatch(in["Routemap"], in["Permit"], in["Match"], in["ExtCommunity"], in["ExactMatch"])
}

func RoutemapPermitNoMatchInterfaceTypeName(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Type"]; !ok {
      return nil 
   }
if _, ok := in["Name"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchInterfaceTypeName(in["Routemap"], in["Permit"], in["Match"], in["Interface"], in["Type"], in["Name"])
}

func RoutemapPermitNoMatchIPAddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIPAddressAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Address"], in["AccessList"])
}

func RoutemapPermitNoMatchIPAddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIPAddressPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Address"], in["PrefixList"])
}

func RoutemapPermitNoMatchIPNexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIPNexthopAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Nexthop"], in["AccessList"])
}

func RoutemapPermitNoMatchIPNexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIPNexthopPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP"], in["Nexthop"], in["PrefixList"])
}

func RoutemapPermitNoMatchIP6AddressAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIP6AddressAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Address"], in["AccessList"])
}

func RoutemapPermitNoMatchIP6AddressPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Address"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIP6AddressPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Address"], in["PrefixList"])
}

func RoutemapPermitNoMatchIP6NexthopAccessList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["AccessList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIP6NexthopAccessList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Nexthop"], in["AccessList"])
}

func RoutemapPermitNoMatchIP6NexthopPrefixList(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchIP6NexthopPrefixList(in["Routemap"], in["Permit"], in["Match"], in["IP6"], in["Nexthop"], in["PrefixList"])
}

func RoutemapPermitNoMatchMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchMetric(in["Routemap"], in["Permit"], in["Match"], in["Metric"])
}

func RoutemapPermitNoMatchOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchOrigin(in["Routemap"], in["Permit"], in["Match"], in["Origin"])
}

func RoutemapPermitNoMatchRouteTypeExternal(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["RouteType"]; !ok {
      return nil 
   }
if _, ok := in["External"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchRouteTypeExternal(in["Routemap"], in["Permit"], in["Match"], in["RouteType"], in["External"])
}

func RoutemapPermitNoMatchTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Match"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoMatchTag(in["Routemap"], in["Permit"], in["Match"], in["Tag"])
}

func RoutemapPermitSetAggregatorASIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetAggregatorASIP(in["Routemap"], in["Permit"], in["Set"], in["Aggregator"], in["AS"], in["IP"])
}

func RoutemapPermitSetASPathPrepend(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
if _, ok := in["Prepend"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetASPathPrepend(in["Routemap"], in["Permit"], in["Set"], in["ASPath"], in["Prepend"])
}

func RoutemapPermitSetAtomicAggregator(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["AtomicAggregator"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetAtomicAggregator(in["Routemap"], in["Permit"], in["Set"], in["AtomicAggregator"])
}

func RoutemapPermitSetCommListDelete(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["CommList"]; !ok {
      return nil 
   }
if _, ok := in["Delete"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommListDelete(in["Routemap"], in["Permit"], in["Set"], in["CommList"], in["Delete"])
}

func RoutemapPermitSetCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunity(in["Routemap"], in["Permit"], in["Set"], in["Community"])
}

func RoutemapPermitSetCommunityNone(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["None"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunityNone(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["None"])
}

func RoutemapPermitSetCommunityInternet(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["Internet"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunityInternet(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["Internet"])
}

func RoutemapPermitSetCommunityLocalAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["LocalAS"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunityLocalAS(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["LocalAS"])
}

func RoutemapPermitSetCommunityNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunityNoAdvertise(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["NoAdvertise"])
}

func RoutemapPermitSetCommunityNoExport(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoExport"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetCommunityNoExport(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["NoExport"])
}

func RoutemapPermitSetDampening(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Dampening"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetDampening(in["Routemap"], in["Permit"], in["Set"], in["Dampening"])
}

func RoutemapPermitSetExtCommunityRT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["RT"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetExtCommunityRT(in["Routemap"], in["Permit"], in["Set"], in["ExtCommunity"], in["RT"])
}

func RoutemapPermitSetExtCommunitySOO(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["SOO"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetExtCommunitySOO(in["Routemap"], in["Permit"], in["Set"], in["ExtCommunity"], in["SOO"])
}

func RoutemapPermitSetIPNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetIPNexthop(in["Routemap"], in["Permit"], in["Set"], in["IP"], in["Nexthop"])
}

func RoutemapPermitSetIP6Nexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetIP6Nexthop(in["Routemap"], in["Permit"], in["Set"], in["IP6"], in["Nexthop"])
}

func RoutemapPermitSetLevel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Level"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetLevel(in["Routemap"], in["Permit"], in["Set"], in["Level"])
}

func RoutemapPermitSetLocalPreference(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["LocalPreference"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetLocalPreference(in["Routemap"], in["Permit"], in["Set"], in["LocalPreference"])
}

func RoutemapPermitSetMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetMetric(in["Routemap"], in["Permit"], in["Set"], in["Metric"])
}

func RoutemapPermitSetMetricType(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["MetricType"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetMetricType(in["Routemap"], in["Permit"], in["Set"], in["MetricType"])
}

func RoutemapPermitSetOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetOrigin(in["Routemap"], in["Permit"], in["Set"], in["Origin"])
}

func RoutemapPermitSetOriginatorID(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["OriginatorID"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetOriginatorID(in["Routemap"], in["Permit"], in["Set"], in["OriginatorID"])
}

func RoutemapPermitSetTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetTag(in["Routemap"], in["Permit"], in["Set"], in["Tag"])
}

func RoutemapPermitSetWeight(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Weight"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitSetWeight(in["Routemap"], in["Permit"], in["Set"], in["Weight"])
}

func RoutemapPermitNoSetAggregatorAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetAggregatorAS(in["Routemap"], in["Permit"], in["Set"], in["Aggregator"], in["AS"])
}

func RoutemapPermitNoSetAggregatorASIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Aggregator"]; !ok {
      return nil 
   }
if _, ok := in["AS"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetAggregatorASIP(in["Routemap"], in["Permit"], in["Set"], in["Aggregator"], in["AS"], in["IP"])
}

func RoutemapPermitNoSetASPathPrepend(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ASPath"]; !ok {
      return nil 
   }
if _, ok := in["Prepend"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetASPathPrepend(in["Routemap"], in["Permit"], in["Set"], in["ASPath"], in["Prepend"])
}

func RoutemapPermitNoSetAtomicAggregator(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["AtomicAggregator"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetAtomicAggregator(in["Routemap"], in["Permit"], in["Set"], in["AtomicAggregator"])
}

func RoutemapPermitNoSetCommListDelete(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["CommList"]; !ok {
      return nil 
   }
if _, ok := in["Delete"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommListDelete(in["Routemap"], in["Permit"], in["Set"], in["CommList"], in["Delete"])
}

func RoutemapPermitNoSetCommunity(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunity(in["Routemap"], in["Permit"], in["Set"], in["Community"])
}

func RoutemapPermitNoSetCommunityNone(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["None"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunityNone(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["None"])
}

func RoutemapPermitNoSetCommunityInternet(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["Internet"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunityInternet(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["Internet"])
}

func RoutemapPermitNoSetCommunityLocalAS(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["LocalAS"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunityLocalAS(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["LocalAS"])
}

func RoutemapPermitNoSetCommunityNoAdvertise(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoAdvertise"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunityNoAdvertise(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["NoAdvertise"])
}

func RoutemapPermitNoSetCommunityNoExport(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Community"]; !ok {
      return nil 
   }
if _, ok := in["NoExport"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetCommunityNoExport(in["Routemap"], in["Permit"], in["Set"], in["Community"], in["NoExport"])
}

func RoutemapPermitNoSetDampening(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Dampening"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetDampening(in["Routemap"], in["Permit"], in["Set"], in["Dampening"])
}

func RoutemapPermitNoSetExtCommunityRT(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["RT"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetExtCommunityRT(in["Routemap"], in["Permit"], in["Set"], in["ExtCommunity"], in["RT"])
}

func RoutemapPermitNoSetExtCommunitySOO(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["ExtCommunity"]; !ok {
      return nil 
   }
if _, ok := in["SOO"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetExtCommunitySOO(in["Routemap"], in["Permit"], in["Set"], in["ExtCommunity"], in["SOO"])
}

func RoutemapPermitNoSetIPNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetIPNexthop(in["Routemap"], in["Permit"], in["Set"], in["IP"], in["Nexthop"])
}

func RoutemapPermitNoSetIP6Nexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetIP6Nexthop(in["Routemap"], in["Permit"], in["Set"], in["IP6"], in["Nexthop"])
}

func RoutemapPermitNoSetLevel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Level"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetLevel(in["Routemap"], in["Permit"], in["Set"], in["Level"])
}

func RoutemapPermitNoSetLocalPreference(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["LocalPreference"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetLocalPreference(in["Routemap"], in["Permit"], in["Set"], in["LocalPreference"])
}

func RoutemapPermitNoSetMetric(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Metric"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetMetric(in["Routemap"], in["Permit"], in["Set"], in["Metric"])
}

func RoutemapPermitNoSetMetricType(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["MetricType"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetMetricType(in["Routemap"], in["Permit"], in["Set"], in["MetricType"])
}

func RoutemapPermitNoSetOrigin(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Origin"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetOrigin(in["Routemap"], in["Permit"], in["Set"], in["Origin"])
}

func RoutemapPermitNoSetOriginatorID(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["OriginatorID"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetOriginatorID(in["Routemap"], in["Permit"], in["Set"], in["OriginatorID"])
}

func RoutemapPermitNoSetTag(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Tag"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetTag(in["Routemap"], in["Permit"], in["Set"], in["Tag"])
}

func RoutemapPermitNoSetWeight(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Set"]; !ok {
      return nil 
   }
if _, ok := in["Weight"]; !ok {
      return nil 
   }
      return sw.RoutemapPermitNoSetWeight(in["Routemap"], in["Permit"], in["Set"], in["Weight"])
}

func NoRoutemap(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
      return sw.NoRoutemap(in["Routemap"])
}

func NoRoutemapDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.NoRoutemapDeny(in["Routemap"], in["Deny"])
}

func NoRoutemapPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Routemap"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.NoRoutemapPermit(in["Routemap"], in["Permit"])
}

func IPPrefixListSequenceNumber(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["SequenceNumber"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSequenceNumber(in["IP"], in["PrefixList"], in["SequenceNumber"])
}

func NoIPPrefixListSequenceNumber(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["SequenceNumber"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSequenceNumber(in["IP"], in["PrefixList"], in["SequenceNumber"])
}

func IPPrefixListDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDeny(in["IP"], in["PrefixList"], in["Deny"])
}

func IPPrefixListDescription(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Description"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDescription(in["IP"], in["PrefixList"], in["Description"])
}

func IPPrefixListDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDenyGe(in["IP"], in["PrefixList"], in["Deny"], in["Ge"])
}

func IPPrefixListDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDenyGeLe(in["IP"], in["PrefixList"], in["Deny"], in["Ge"], in["Le"])
}

func IPPrefixListDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDenyLe(in["IP"], in["PrefixList"], in["Deny"], in["Le"])
}

func IPPrefixListDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListDenyLeGe(in["IP"], in["PrefixList"], in["Deny"], in["Le"], in["Ge"])
}

func IPPrefixListSeqDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqDeny(in["IP"], in["PrefixList"], in["Seq"], in["Deny"])
}

func IPPrefixListSeqDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqDenyGe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"])
}

func IPPrefixListSeqDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqDenyGeLe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"], in["Le"])
}

func IPPrefixListSeqDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqDenyLe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"])
}

func IPPrefixListSeqDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqDenyLeGe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"], in["Ge"])
}

func NoIPPrefixListDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDeny(in["IP"], in["PrefixList"], in["Deny"])
}

func NoIPPrefixListDescription(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Description"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDescription(in["IP"], in["PrefixList"], in["Description"])
}

func NoIPPrefixListDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDenyGe(in["IP"], in["PrefixList"], in["Deny"], in["Ge"])
}

func NoIPPrefixListDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDenyGeLe(in["IP"], in["PrefixList"], in["Deny"], in["Ge"], in["Le"])
}

func NoIPPrefixListDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDenyLe(in["IP"], in["PrefixList"], in["Deny"], in["Le"])
}

func NoIPPrefixListDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListDenyLeGe(in["IP"], in["PrefixList"], in["Deny"], in["Le"], in["Ge"])
}

func NoIPPrefixListSeqDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqDeny(in["IP"], in["PrefixList"], in["Seq"], in["Deny"])
}

func NoIPPrefixListSeqDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqDenyGe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"])
}

func NoIPPrefixListSeqDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqDenyGeLe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"], in["Le"])
}

func NoIPPrefixListSeqDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqDenyLe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"])
}

func NoIPPrefixListSeqDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqDenyLeGe(in["IP"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"], in["Ge"])
}

func IPPrefixListPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.IPPrefixListPermit(in["IP"], in["PrefixList"], in["Permit"])
}

func IPPrefixListPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListPermitGe(in["IP"], in["PrefixList"], in["Permit"], in["Ge"])
}

func IPPrefixListPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListPermitGeLe(in["IP"], in["PrefixList"], in["Permit"], in["Ge"], in["Le"])
}

func IPPrefixListPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListPermitLe(in["IP"], in["PrefixList"], in["Permit"], in["Le"])
}

func IPPrefixListPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListPermitLeGe(in["IP"], in["PrefixList"], in["Permit"], in["Le"], in["Ge"])
}

func IPPrefixListSeqPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqPermit(in["IP"], in["PrefixList"], in["Seq"], in["Permit"])
}

func IPPrefixListSeqPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqPermitGe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"])
}

func IPPrefixListSeqPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqPermitGeLe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"], in["Le"])
}

func IPPrefixListSeqPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqPermitLe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"])
}

func IPPrefixListSeqPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IPPrefixListSeqPermitLeGe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"], in["Ge"])
}

func NoIPPrefixListPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListPermit(in["IP"], in["PrefixList"], in["Permit"])
}

func NoIPPrefixListPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListPermitGe(in["IP"], in["PrefixList"], in["Permit"], in["Ge"])
}

func NoIPPrefixListPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListPermitGeLe(in["IP"], in["PrefixList"], in["Permit"], in["Ge"], in["Le"])
}

func NoIPPrefixListPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListPermitLe(in["IP"], in["PrefixList"], in["Permit"], in["Le"])
}

func NoIPPrefixListPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListPermitLeGe(in["IP"], in["PrefixList"], in["Permit"], in["Le"], in["Ge"])
}

func NoIPPrefixListSeqPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqPermit(in["IP"], in["PrefixList"], in["Seq"], in["Permit"])
}

func NoIPPrefixListSeqPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqPermitGe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"])
}

func NoIPPrefixListSeqPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqPermitGeLe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"], in["Le"])
}

func NoIPPrefixListSeqPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqPermitLe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"])
}

func NoIPPrefixListSeqPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIPPrefixListSeqPermitLeGe(in["IP"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"], in["Ge"])
}

func IP6PrefixListSequenceNumber(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["SequenceNumber"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSequenceNumber(in["IP6"], in["PrefixList"], in["SequenceNumber"])
}

func NoIP6PrefixListSequenceNumber(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["SequenceNumber"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSequenceNumber(in["IP6"], in["PrefixList"], in["SequenceNumber"])
}

func IP6PrefixListDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDeny(in["IP6"], in["PrefixList"], in["Deny"])
}

func IP6PrefixListDescription(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Description"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDescription(in["IP6"], in["PrefixList"], in["Description"])
}

func IP6PrefixListDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDenyGe(in["IP6"], in["PrefixList"], in["Deny"], in["Ge"])
}

func IP6PrefixListDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDenyGeLe(in["IP6"], in["PrefixList"], in["Deny"], in["Ge"], in["Le"])
}

func IP6PrefixListDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDenyLe(in["IP6"], in["PrefixList"], in["Deny"], in["Le"])
}

func IP6PrefixListDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListDenyLeGe(in["IP6"], in["PrefixList"], in["Deny"], in["Le"], in["Ge"])
}

func IP6PrefixListSeqDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqDeny(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"])
}

func IP6PrefixListSeqDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqDenyGe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"])
}

func IP6PrefixListSeqDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqDenyGeLe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"], in["Le"])
}

func IP6PrefixListSeqDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqDenyLe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"])
}

func IP6PrefixListSeqDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqDenyLeGe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"], in["Ge"])
}

func NoIP6PrefixListDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDeny(in["IP6"], in["PrefixList"], in["Deny"])
}

func NoIP6PrefixListDescription(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Description"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDescription(in["IP6"], in["PrefixList"], in["Description"])
}

func NoIP6PrefixListDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDenyGe(in["IP6"], in["PrefixList"], in["Deny"], in["Ge"])
}

func NoIP6PrefixListDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDenyGeLe(in["IP6"], in["PrefixList"], in["Deny"], in["Ge"], in["Le"])
}

func NoIP6PrefixListDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDenyLe(in["IP6"], in["PrefixList"], in["Deny"], in["Le"])
}

func NoIP6PrefixListDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListDenyLeGe(in["IP6"], in["PrefixList"], in["Deny"], in["Le"], in["Ge"])
}

func NoIP6PrefixListSeqDeny(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqDeny(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"])
}

func NoIP6PrefixListSeqDenyGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqDenyGe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"])
}

func NoIP6PrefixListSeqDenyGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqDenyGeLe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Ge"], in["Le"])
}

func NoIP6PrefixListSeqDenyLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqDenyLe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"])
}

func NoIP6PrefixListSeqDenyLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Deny"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqDenyLeGe(in["IP6"], in["PrefixList"], in["Seq"], in["Deny"], in["Le"], in["Ge"])
}

func IP6PrefixListPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListPermit(in["IP6"], in["PrefixList"], in["Permit"])
}

func IP6PrefixListPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListPermitGe(in["IP6"], in["PrefixList"], in["Permit"], in["Ge"])
}

func IP6PrefixListPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListPermitGeLe(in["IP6"], in["PrefixList"], in["Permit"], in["Ge"], in["Le"])
}

func IP6PrefixListPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListPermitLe(in["IP6"], in["PrefixList"], in["Permit"], in["Le"])
}

func IP6PrefixListPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListPermitLeGe(in["IP6"], in["PrefixList"], in["Permit"], in["Le"], in["Ge"])
}

func IP6PrefixListSeqPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqPermit(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"])
}

func IP6PrefixListSeqPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqPermitGe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"])
}

func IP6PrefixListSeqPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqPermitGeLe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"], in["Le"])
}

func IP6PrefixListSeqPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqPermitLe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"])
}

func IP6PrefixListSeqPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.IP6PrefixListSeqPermitLeGe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"], in["Ge"])
}

func NoIP6PrefixListPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListPermit(in["IP6"], in["PrefixList"], in["Permit"])
}

func NoIP6PrefixListPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListPermitGe(in["IP6"], in["PrefixList"], in["Permit"], in["Ge"])
}

func NoIP6PrefixListPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListPermitGeLe(in["IP6"], in["PrefixList"], in["Permit"], in["Ge"], in["Le"])
}

func NoIP6PrefixListPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListPermitLe(in["IP6"], in["PrefixList"], in["Permit"], in["Le"])
}

func NoIP6PrefixListPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListPermitLeGe(in["IP6"], in["PrefixList"], in["Permit"], in["Le"], in["Ge"])
}

func NoIP6PrefixListSeqPermit(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqPermit(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"])
}

func NoIP6PrefixListSeqPermitGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqPermitGe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"])
}

func NoIP6PrefixListSeqPermitGeLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqPermitGeLe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Ge"], in["Le"])
}

func NoIP6PrefixListSeqPermitLe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqPermitLe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"])
}

func NoIP6PrefixListSeqPermitLeGe(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["PrefixList"]; !ok {
      return nil 
   }
if _, ok := in["Seq"]; !ok {
      return nil 
   }
if _, ok := in["Permit"]; !ok {
      return nil 
   }
if _, ok := in["Le"]; !ok {
      return nil 
   }
if _, ok := in["Ge"]; !ok {
      return nil 
   }
      return sw.NoIP6PrefixListSeqPermitLeGe(in["IP6"], in["PrefixList"], in["Seq"], in["Permit"], in["Le"], in["Ge"])
}

func ClearIPRoute(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["Clear"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
      return sw.ClearIPRoute(in["Clear"], in["IP"], in["Route"])
}

func ClearIPRouteKernel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["Clear"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Kernel"]; !ok {
      return nil 
   }
      return sw.ClearIPRouteKernel(in["Clear"], in["IP"], in["Route"], in["Kernel"])
}

func IPRoutePrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.IPRoutePrefixNexthop(in["IP"], in["Route"], in["Prefix"], in["Nexthop"])
}

func IPRoutePrefixNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IPRoutePrefixNexthopDistance(in["IP"], in["Route"], in["Prefix"], in["Nexthop"], in["Distance"])
}

func IPRoutePrefixNexthopSrc(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Src"]; !ok {
      return nil 
   }
      return sw.IPRoutePrefixNexthopSrc(in["IP"], in["Route"], in["Prefix"], in["Nexthop"], in["Src"])
}

func IPRouteDefaultNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Default"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.IPRouteDefaultNexthop(in["IP"], in["Route"], in["Default"], in["Nexthop"])
}

func IPRouteDefaultNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Default"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IPRouteDefaultNexthopDistance(in["IP"], in["Route"], in["Default"], in["Nexthop"], in["Distance"])
}

func IPRouteVRFPrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.IPRouteVRFPrefixNexthopOIF(in["IP"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func IPRouteVRFPrefixOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.IPRouteVRFPrefixOIF(in["IP"], in["Route"], in["VRF"], in["Prefix"], in["OIF"])
}

func NoIPRoutePrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.NoIPRoutePrefixNexthop(in["IP"], in["Route"], in["Prefix"], in["Nexthop"])
}

func NoIPRoutePrefixNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.NoIPRoutePrefixNexthopDistance(in["IP"], in["Route"], in["Prefix"], in["Nexthop"], in["Distance"])
}

func NoIPRouteDefaultNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Default"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.NoIPRouteDefaultNexthop(in["IP"], in["Route"], in["Default"], in["Nexthop"])
}

func NoIPRouteDefaultNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Default"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.NoIPRouteDefaultNexthopDistance(in["IP"], in["Route"], in["Default"], in["Nexthop"], in["Distance"])
}

func NoIPRouteVRFPrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.NoIPRouteVRFPrefixNexthopOIF(in["IP"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func NoIPRouteVRFPrefixOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.NoIPRouteVRFPrefixOIF(in["IP"], in["Route"], in["VRF"], in["Prefix"], in["OIF"])
}

func IP6RoutePrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixNexthop(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"])
}

func IP6RoutePrefixNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixNexthopDistance(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"], in["Distance"])
}

func IP6RoutePrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixNexthopOIF(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func IP6RoutePrefixNexthopOIFDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixNexthopOIFDistance(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"], in["OIF"], in["Distance"])
}

func IP6RoutePrefixTunnel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixTunnel(in["IP6"], in["Route"], in["Prefix"], in["Tunnel"])
}

func IP6RoutePrefixTunnelDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RoutePrefixTunnelDistance(in["IP6"], in["Route"], in["Prefix"], in["Tunnel"], in["Distance"])
}

func IP6RouteVRFPrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixNexthop(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"])
}

func IP6RouteVRFPrefixNexthopDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixNexthopDistance(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["Distance"])
}

func IP6RouteVRFPrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixNexthopOIF(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func IP6RouteVRFPrefixNexthopOIFDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixNexthopOIFDistance(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["OIF"], in["Distance"])
}

func IP6RouteVRFPrefixTunnel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixTunnel(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Tunnel"])
}

func IP6RouteVRFPrefixTunnelDistance(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
if _, ok := in["Distance"]; !ok {
      return nil 
   }
      return sw.IP6RouteVRFPrefixTunnelDistance(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Tunnel"], in["Distance"])
}

func NoIP6RoutePrefix(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
      return sw.NoIP6RoutePrefix(in["IP6"], in["Route"], in["Prefix"])
}

func NoIP6RoutePrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.NoIP6RoutePrefixNexthop(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"])
}

func NoIP6RoutePrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.NoIP6RoutePrefixNexthopOIF(in["IP6"], in["Route"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func NoIP6RoutePrefixTunnel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
      return sw.NoIP6RoutePrefixTunnel(in["IP6"], in["Route"], in["Prefix"], in["Tunnel"])
}

func NoIP6RouteVRFPrefix(sw Switch, in map[string]string) []*command.Command{
if len(in) != 5 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
      return sw.NoIP6RouteVRFPrefix(in["IP6"], in["Route"], in["VRF"], in["Prefix"])
}

func NoIP6RouteVRFPrefixNexthop(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
      return sw.NoIP6RouteVRFPrefixNexthop(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"])
}

func NoIP6RouteVRFPrefixNexthopOIF(sw Switch, in map[string]string) []*command.Command{
if len(in) != 7 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Nexthop"]; !ok {
      return nil 
   }
if _, ok := in["OIF"]; !ok {
      return nil 
   }
      return sw.NoIP6RouteVRFPrefixNexthopOIF(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Nexthop"], in["OIF"])
}

func NoIP6RouteVRFPrefixTunnel(sw Switch, in map[string]string) []*command.Command{
if len(in) != 6 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Route"]; !ok {
      return nil 
   }
if _, ok := in["VRF"]; !ok {
      return nil 
   }
if _, ok := in["Prefix"]; !ok {
      return nil 
   }
if _, ok := in["Tunnel"]; !ok {
      return nil 
   }
      return sw.NoIP6RouteVRFPrefixTunnel(in["IP6"], in["Route"], in["VRF"], in["Prefix"], in["Tunnel"])
}

func InterfaceIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.InterfaceIP(in["Interface"], in["IP"])
}

func InterfaceIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.InterfaceIP2(in["Interface"], in["IP2"])
}

func NoInterfaceIP(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP"]; !ok {
      return nil 
   }
      return sw.NoInterfaceIP(in["Interface"], in["IP"])
}

func NoInterfaceIP2(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP2"]; !ok {
      return nil 
   }
      return sw.NoInterfaceIP2(in["Interface"], in["IP2"])
}

func InterfaceShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["Shutdown"]; !ok {
      return nil 
   }
      return sw.InterfaceShutdown(in["Interface"], in["Shutdown"])
}

func InterfaceNoShutdown(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["NoShutdown"]; !ok {
      return nil 
   }
      return sw.InterfaceNoShutdown(in["Interface"], in["NoShutdown"])
}

func InterfaceIP6Enable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Enable"]; !ok {
      return nil 
   }
      return sw.InterfaceIP6Enable(in["Interface"], in["IP6"], in["Enable"])
}

func NoInterfaceIP6Enable(sw Switch, in map[string]string) []*command.Command{
if len(in) != 4 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
if _, ok := in["Enable"]; !ok {
      return nil 
   }
      return sw.NoInterfaceIP6Enable(in["Interface"], in["IP6"], in["Enable"])
}

func InterfaceIP6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
      return sw.InterfaceIP6(in["Interface"], in["IP6"])
}

func InterfaceIP6LL(sw Switch, in map[string]string) []*command.Command{
if len(in) != 2 {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6LL"]; !ok {
      return nil 
   }
      return sw.InterfaceIP6LL(in["Interface"], in["IP6LL"])
}

func NoInterfaceIP6(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6"]; !ok {
      return nil 
   }
      return sw.NoInterfaceIP6(in["Interface"], in["IP6"])
}

func NoInterfaceIP6LL(sw Switch, in map[string]string) []*command.Command{
if len(in) != 3 {
      return nil 
   }
if _, ok := in["No"]; !ok {
      return nil 
   }
if _, ok := in["Interface"]; !ok {
      return nil 
   }
if _, ok := in["IP6LL"]; !ok {
      return nil 
   }
      return sw.NoInterfaceIP6LL(in["Interface"], in["IP6LL"])
}

