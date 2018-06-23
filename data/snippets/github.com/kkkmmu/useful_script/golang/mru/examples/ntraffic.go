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
		panic(err)
		err = n.KillSessionByName(n2x.DEFAULTSESSIONNAME)
		if err != nil {
			panic(err)
		}

		sess, err = n.OpenNewSession("103/3", "103/4")
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

	ports, err := sess.GetPorts()
	fmt.Printf("%q ", ports)
	//	time.Sleep(time.Duration(time.Second * 40))
	for i, port := range ports {
		err := port.RemoveAllTraffics()
		if err != nil {
			panic(err)
		}

		t, err := port.AddTraffic("10", fmt.Sprintf("TR%d", i))
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", t)

		_, err = t.ListStreamGroups()
		if err != nil {
			panic(err)
		}

		sgs, err := t.GetAllStreamGroups()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", sgs)

		for _, sg := range sgs {
			sg.Disable()
		}

		err = sgs[0].SetIPv4UDP()
		if err != nil {
			panic(err)
		}
		err = sgs[1].SetIPv4TCP()
		if err != nil {
			panic(err)
		}

		err = sgs[2].SetIPv6TCP()
		if err != nil {
			panic(err)
		}
		err = sgs[3].SetIPv6UDP()
		if err != nil {
			panic(err)
		}

		err = sgs[4].SetIPv6ND()
		if err != nil {
			panic(err)
		}
		err = sgs[5].SetIPv4ARP()
		if err != nil {
			panic(err)
		}

		for _, sg := range sgs {
			sg.Enable()
		}

	}
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
