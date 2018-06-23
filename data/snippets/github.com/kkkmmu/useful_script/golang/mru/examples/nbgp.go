package main

import (
	"fmt"
	"log"
	"n2x"
	//	"time"
)

func main() {
	n, err := n2x.New("10.71.20.231", "9001")
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

	ports, err := sess.GetReservedPorts()
	fmt.Printf("%q ", ports)
	sess.StopRoutingEngine()
	//	time.Sleep(time.Duration(time.Second * 40))
	for i, port := range ports {
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

		_, err = port.LegacyLinkGetAllIP6AddressPools()
		if err != nil {
			panic(err)
		}

		err = port.RemoveAllBGPs()
		if err != nil {
			panic(err)
		}

		b, err := port.AddBGP(n2x.IPV4_EBGP, "123", "456", "test")
		if err != nil {
			panic(err)
		}

		rp, err := b.AddRoutes("156.1.1.1", "32", "10000", "1")
		if err != nil {
			panic(err)
		}

		fmt.Printf("%q\n", b)
		fmt.Printf("%q\n", rp)

	}

	sess.StartRoutingEngine()
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
