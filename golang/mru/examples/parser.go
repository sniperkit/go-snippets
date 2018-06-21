package main

import (
	"command"
	"dsl"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Printf("%#v\n", dsl.ParserV8)
	file, err := ioutil.ReadFile("asset/instructions/instruction.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	ins := make([]string, 0, 10)

	for _, l := range lines {
		fmt.Println(l)
		if l == "" {
			continue
		}
		ins = append(ins, l)
	}

	for _, i := range ins {
		fmt.Println(string(dsl.MapInstructionToKey(i)))
	}

	for _, i := range ins {
		fmt.Printf("%#v\n", dsl.GetInstructionParamsFromCMD(i))
	}

	for _, i := range ins {
		fmt.Println(i)
		cmds, err := dsl.Engine.Parse("V8", &command.Command{CMD: "$" + i})
		if err == nil {
			for _, c := range cmds {
				fmt.Printf("%#v", c)
			}
		} else {
			fmt.Printf("%#v", err)
		}
	}

	cmds, err := dsl.Engine.Parse("V8", &command.Command{
		Mode: "config",
		CMD:  "show running-config",
	})

	if err != nil {
		panic(err)
	}

	for _, c := range cmds {
		fmt.Printf("%#v\n", c)
	}

	TestInstruction("V8", "$OSPF6(1234)")
	TestInstruction("V8", "$OSPF6(1234).Interface(100)")
	TestInstruction("V8", "$OSPF6(1234).Interface(100).Cost(100)")
	TestInstruction("V8", "$No().OSPF6(1234).Area(789).Stub()")
	TestInstruction("V8", "$No().OSPF6(1234).Area(456).Stub().NoSummary()")

	TestInstruction("V8", "$OSPF6(3456).Area(1111).NSSA()")
	TestInstruction("V8", "$OSPF6(3456).Area(1111).NSSA().DefaultOriginate()")
	TestInstruction("V8", "$OSPF6(3456).Area(1111).NSSA().NoRedistribution()")
	TestInstruction("V8", "$OSPF6(3456).Area(1111).NSSA().NoSummary()")
	TestInstruction("V8", "$OSPF6(3456).Area(1111).NSSA().StabilityInterval(4321)")
	TestInstruction("V8", "$OSPF6(3456).Area(1111).Translatorrole(candidate)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).NSSA()")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).NSSA().DefaultOriginate()")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).NSSA().NoRedistribution()")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).NSSA().NoSummary()")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).NSSA().StabilityInterval(4321)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(1111).Translatorrole(always)")

	TestInstruction("V8", "$OSPF6(1234).Area(4567).Virtuallink(10.10.10.10)")
	TestInstruction("V8", "$OSPF6(1234).Area(4567).Virtuallink(10.1.1.1).DeadInterval(40)")
	TestInstruction("V8", "$OSPF6(1234).Distance().External(10000)")
	TestInstruction("V8", "$OSPF6(1234).Distance().Inter(10000)")
	TestInstruction("V8", "$OSPF6(1234).Distance().Intra(10000)")
	TestInstruction("V8", "$OSPF6(1234).Distance().Inter(10000).Intra(10000)")
	TestInstruction("V8", "$OSPF6(1234).Distance().Inter(10000).External(10000)")
	TestInstruction("V8", "$OSPF6(1234).Distance().Inter(10000).Intra(10000).External(10000)")
	TestInstruction("V8", "$No().OSPF6(1234).Distance().External(10000)")
	TestInstruction("V8", "$No().OSPF6(1234).Distance().Inter(10000)")
	TestInstruction("V8", "$No().OSPF6(1234).Distance().Intra(10000)")
	TestInstruction("V8", "$No().OSPF6(1234).Distance().Inter(10000).Intra(10000)")
	TestInstruction("V8", "$No().OSPF6(1234).Distance().Inter(10000).External(10000)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).DeadInterval(#Interval)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).HelloInterval(4321)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).Instanceid(3456)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).RetransmitInterval(4321)")
	TestInstruction("V8", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).TransmitDelay(3333)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).DeadInterval(4321)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).HelloInterval(4321)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).Instanceid(3456)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).RetransmitInterval(4321)")
	TestInstruction("V8", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).TransmitDelay(3333)")

	TestInstruction("V8", "$Port(3).Enable()")
	TestInstruction("V8", "$Port(3).Disable()")
	TestInstruction("V8", "$Port(3).Pvid()")
	TestInstruction("V8", "$Port(3).Speed()")
	TestInstruction("V8", "$VLAN(3)")
	TestInstruction("V8", "$No().VLAN(3)")
	TestInstruction("V8", "$VLAN(124).Add(4)")
	TestInstruction("V8", "$VLAN(124).AddT(4)")
	TestInstruction("V8", "$VLAN(124).Del(4)")
	TestInstruction("V8", "$VLAN(124).DelT(4)")
	TestInstruction("V8", "$VLAN(124).Shutdown()")
	TestInstruction("V8", "$VLAN(124).NoShutdown()")
	TestInstruction("V8", "$VLAN(124).IP(100.1.1.1/24)")
	TestInstruction("V8", "$No().VLAN(124).IP(100.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).IP2(200.1.1.1/24)")
	TestInstruction("V8", "$No().VLAN(124).IP2(200.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).AddT(4).IP(100.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).AddT(4).IP2(200.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).AddT(4).IP(100.1.1.1/24).NoShutdown()")
	TestInstruction("V8", "$VLAN(124).AddT(4).IP2(200.1.1.1/24).NoShutdown()")
	TestInstruction("V8", "$VLAN(124).AddUT(4).IP(100.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).AddUT(4).IP2(200.1.1.1/24)")
	TestInstruction("V8", "$VLAN(124).AddUT(4).IP(100.1.1.1/24).NoShutdown()")
	TestInstruction("V8", "$VLAN(124).AddUT(4).IP2(200.1.1.1/24).NoShutdown()")
	TestInstruction("V8", "$VLAN(124).IP6().Enable()")
	TestInstruction("V8", "$No().VLAN(124).IP6().Enable()")
	TestInstruction("V8", "$VLAN(124).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V8", "$No().VLAN(124).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V8", "$VLAN(124).IP6LL(fe80:0001:0002:0003::1)")
	TestInstruction("V8", "$VLAN(124).IP6LL(fe80:0001:0002:0003::1).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V8", "$No().VLAN(124).IP6LL(fe80:0001:0002:0003::1)")

	TestInstruction("V5", "$OSPF6(1234)")
	TestInstruction("V5", "$OSPF6(1234).Interface(100)")
	TestInstruction("V5", "$OSPF6(1234).Interface(100).Cost(100)")
	TestInstruction("V5", "$No().OSPF6(1234).Area(789).Stub()")
	TestInstruction("V5", "$No().OSPF6(1234).Area(456).Stub().NoSummary()")

	TestInstruction("V5", "$OSPF6(3456).Area(1111).NSSA()")
	TestInstruction("V5", "$OSPF6(3456).Area(1111).NSSA().DefaultOriginate()")
	TestInstruction("V5", "$OSPF6(3456).Area(1111).NSSA().NoRedistribution()")
	TestInstruction("V5", "$OSPF6(3456).Area(1111).NSSA().NoSummary()")
	TestInstruction("V5", "$OSPF6(3456).Area(1111).NSSA().StabilityInterval(4321)")
	TestInstruction("V5", "$OSPF6(3456).Area(1111).Translatorrole(candidate)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).NSSA()")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).NSSA().DefaultOriginate()")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).NSSA().NoRedistribution()")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).NSSA().NoSummary()")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).NSSA().StabilityInterval(4321)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(1111).Translatorrole(always)")

	TestInstruction("V5", "$OSPF6(1234).Area(4567).Virtuallink(10.10.10.10)")
	TestInstruction("V5", "$OSPF6(1234).Area(4567).Virtuallink(10.1.1.1).DeadInterval(40)")
	TestInstruction("V5", "$OSPF6(1234).Distance().External(10000)")
	TestInstruction("V5", "$OSPF6(1234).Distance().Inter(10000)")
	TestInstruction("V5", "$OSPF6(1234).Distance().Intra(10000)")
	TestInstruction("V5", "$OSPF6(1234).Distance().Inter(10000).Intra(10000)")
	TestInstruction("V5", "$OSPF6(1234).Distance().Inter(10000).External(10000)")
	TestInstruction("V5", "$OSPF6(1234).Distance().Inter(10000).Intra(10000).External(10000)")
	TestInstruction("V5", "$No().OSPF6(1234).Distance().External(10000)")
	TestInstruction("V5", "$No().OSPF6(1234).Distance().Inter(10000)")
	TestInstruction("V5", "$No().OSPF6(1234).Distance().Intra(10000)")
	TestInstruction("V5", "$No().OSPF6(1234).Distance().Inter(10000).Intra(10000)")
	TestInstruction("V5", "$No().OSPF6(1234).Distance().Inter(10000).External(10000)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).DeadInterval(#Interval)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).HelloInterval(4321)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).Instanceid(3456)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).RetransmitInterval(4321)")
	TestInstruction("V5", "$OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).TransmitDelay(3333)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).DeadInterval(4321)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).HelloInterval(4321)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).Instanceid(3456)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).RetransmitInterval(4321)")
	TestInstruction("V5", "$No().OSPF6(3456).Area(22222).Virtuallink(12.1.1.1).TransmitDelay(3333)")

	TestInstruction("V5", "$Port(3).Enable()")
	TestInstruction("V5", "$Port(3).Disable()")
	TestInstruction("V5", "$Port(3).Pvid()")
	TestInstruction("V5", "$Port(3).Speed()")
	TestInstruction("V5", "$VLAN(3)")
	TestInstruction("V5", "$No().VLAN(3)")
	TestInstruction("V5", "$VLAN(124).Add(4)")
	TestInstruction("V5", "$VLAN(124).AddT(4)")
	TestInstruction("V5", "$VLAN(124).Del(4)")
	TestInstruction("V5", "$VLAN(124).DelT(4)")
	TestInstruction("V5", "$VLAN(124).Shutdown()")
	TestInstruction("V5", "$VLAN(124).NoShutdown()")
	TestInstruction("V5", "$VLAN(124).IP(100.1.1.1/24)")
	TestInstruction("V5", "$No().VLAN(124).IP(100.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).IP2(200.1.1.1/24)")
	TestInstruction("V5", "$No().VLAN(124).IP2(200.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).AddT(4).IP(100.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).AddT(4).IP2(200.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).AddT(4).IP(100.1.1.1/24).NoShutdown()")
	TestInstruction("V5", "$VLAN(124).AddT(4).IP2(200.1.1.1/24).NoShutdown()")
	TestInstruction("V5", "$VLAN(124).AddUT(4).IP(100.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).AddUT(4).IP2(200.1.1.1/24)")
	TestInstruction("V5", "$VLAN(124).AddUT(4).IP(100.1.1.1/24).NoShutdown()")
	TestInstruction("V5", "$VLAN(124).AddUT(4).IP2(200.1.1.1/24).NoShutdown()")
	TestInstruction("V5", "$VLAN(124).IP6().Enable()")
	TestInstruction("V5", "$No().VLAN(124).IP6().Enable()")
	TestInstruction("V5", "$VLAN(124).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V5", "$No().VLAN(124).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V5", "$VLAN(124).IP6LL(fe80:0001:0002:0003::1)")
	TestInstruction("V5", "$VLAN(124).IP6LL(fe80:0001:0002:0003::1).IP6(2001:db8:1000:1000::1/64)")
	TestInstruction("V5", "$No().VLAN(124).IP6LL(fe80:0001:0002:0003::1)")

}

func TestInstruction(set, ins string) {
	cmds, err := dsl.Engine.Parse(set, &command.Command{
		CMD: ins,
	})

	if err != nil {
		panic(err)
	}

	for _, c := range cmds {
		fmt.Printf("%#v\n", c)
	}
}
