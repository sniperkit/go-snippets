package server

import (
	"testing"
)

func TestProcessResolveIPv4(t *testing.T) {
	t.Log("Testing process resolve IPv4 ")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}

	ser := NewARPServer(logger)
	ser.initArpParams()
	ser.SetPortPropertyMap()
	conf := ResolveIPv4{
		TargetIP: "10.10.10.10",
		IfId:     20,
	}

	go ser.recvConf(t)
	ser.processResolveIPv4(conf)
}

func (server *ARPServer) recvConf(t *testing.T) {
	conf := <-server.arpEntryUpdateCh
	if conf.PortIfIdx != 20 ||
		conf.IpAddr != "10.10.10.10" ||
		conf.MacAddr != "incomplete" ||
		conf.Type != true {
		t.Errorf("Test unsuccessful")
	} else {
		t.Log("Test successful")
	}
}

func TestProcessArpConf(t *testing.T) {
	t.Log("Testing processArpConf()")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}
	ser := NewARPServer(logger)
	ser.initArpParams()
	conf := ArpConf{
		RefTimeout: 20,
	}

	ser.processArpConf(conf)
}

func TestProcessDeleteResolvedIPv4(t *testing.T) {
	t.Log("Testing processArpConf()")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}
	ser := NewARPServer(logger)
	ser.initArpParams()
	go ser.processDelReslovedIPv4(t)
	ser.processDeleteResolvedIPv4("10.10.10.10")
}

func (server *ARPServer) processDelReslovedIPv4(t *testing.T) {
	msg := <-server.arpDeleteArpEntryFromRibCh
	if msg != "10.10.10.10" {
		t.Errorf("Test unsuccessful")
	} else {
		t.Log("Test successful")
	}
}

func TestProcessArpAction(t *testing.T) {
	t.Log("Testing processArpAction()")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}
	ser := NewARPServer(logger)
	ser.initArpParams()
	go ser.processArpAct(t)
	msg := ArpActionMsg{
		Type: DeleteByIPAddr,
		Obj:  "10.10.10.10",
	}
	ser.processArpAction(msg)
}

func (server *ARPServer) processArpAct(t *testing.T) {
	msg := <-server.arpActionProcessCh
	if msg.Type != DeleteByIPAddr ||
		msg.Obj != "10.10.10.10" {
		t.Errorf("Test unsuccessful")
	} else {
		t.Log("Test successful")
	}
}
