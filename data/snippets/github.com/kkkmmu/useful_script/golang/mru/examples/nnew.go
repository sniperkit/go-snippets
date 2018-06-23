package main

import (
	"fmt"
	"log"
	"n2x"
	"time"
)

func main() {
	n, err := n2x.New("10.71.20.231", "9001")
	if err != nil {
		panic(err)
	}

	err = n.KillSessionByName("SYSTEM")
	if err != nil {
		panic(err)
	}

	sess, err := n.GetSessionByName(n2x.DEFAULTSESSIONNAME)
	if err != nil {
		sess, err = n.OpenNewSession("101/3", "101/4")
		if err != nil {
			panic(err)
		}
	}

	for _, port := range sess.Ports {
		ips, _ := port.LegacyLinkGetSutIPAddresses()
		fmt.Printf("%q ", ips)
	}

	err = sess.StopTest()
	if err != nil {
		panic(err)
	}
	sess.StopRoutingEngine()
	ports, err := sess.GetReservedPorts()
	fmt.Printf("%q ", ports)
	for i, port := range ports {
		port.GetAvailableMediaTypes()
		port.GetMediaType()
		port.SetMediaType(n2x.SFP)
		port.GetMediaType()
		port.LaserOff()
		port.LaserOn()
		port.LegacyLinkGetSutMacAddress()
		//port.LegacyLinkAddSutIPAddresses("120.1.1.5", "30", "1000", "1")
		port.SendAllArpRequests()
		port.SendAllNeighborSolicitations()
		port.LegacyLinkGetAllAddressPools()
		var tip string
		var sip string
		if i == 0 {
			tip = fmt.Sprintf("21.1.1.2")
			sip = fmt.Sprintf("21.1.1.1")
		} else {
			tip = fmt.Sprintf("24.1.1.2")
			sip = fmt.Sprintf("24.1.1.1")
		}
		smac := fmt.Sprintf("00:1%x:00:00:1%x:00", i%255, i%255)
		err = port.LegacyLinkRemoveAllSutIPAddresses()
		if err != nil {
			panic(err)
		}
		port.LegacyLinkAddSutIPAddress(sip)
		err := port.LegacyLinkSet("0", tip, "24", smac, sip)
		if err != nil {
			panic(err)
		}
		tip6 := fmt.Sprintf("200%d:1%d::11", i, i)
		sip6 := fmt.Sprintf("200%d:1%d::12", i, i)

		err = port.LegacyLinkSet6("10", tip6, "64", smac, sip6)
		if err != nil {
			panic(err)
		}

		pools, err := port.LegacyLinkGetAllIP6AddressPools()
		if err != nil {
			panic(err)
		}

		port.GetAllOSPFs()
		err = port.DeleteAllOSPFs()
		if err != nil {
			panic(err)
		}
		err = port.DeleteAllOSPFs()
		if err != nil {
			panic(err)
		}
		port.GetAllOSPFs()
		name := fmt.Sprintf("OSPF_EMULATION_%d", i)
		ospf, err := port.AddOSPF("0.0.0.0", " 21.1.1.25", " 21.1.1.24", name)
		if err != nil {
			panic(err)
		}
		port.GetAllOSPFs()

		lsdb, err := ospf.GetLSDB()
		if err != nil {
			panic(err)
		}

		lsas, err := lsdb.GetAllLSAs()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%q\n", lsas)
		fmt.Printf("%q\n", pools)
	}

	sess.StartRoutingEngine()
	time.Sleep(time.Duration(time.Second * 40))
	for i, port := range ports {
		name := fmt.Sprintf("OSPF_EMULATION_%d", i)
		//	port.GetAllOSPFs()

		ospf, err := port.GetOSPFByName(name)
		if err != nil {
			panic(err)
		}

		lsdb, err := ospf.GetLSDB()
		if err != nil {
			panic(err)
		}

		lsas, err := lsdb.GetAllLSAs()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%q\n", lsas)

	}
	sess.ListModules()
	sess.ListAvailableModules()
	sess.ListAvailablePorts()
	sess.ListLockedPorts()
	sess.ListPorts()
	err = sess.StartTest()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", sess)
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
