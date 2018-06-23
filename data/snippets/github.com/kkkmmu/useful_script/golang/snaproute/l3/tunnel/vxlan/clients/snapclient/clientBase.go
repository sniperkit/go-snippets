// hw.go
package snapclient

import (
	"arpd"
	"asicd/pluginManager/pluginCommon"
	"asicdServices"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"ribd"
	"strconv"
	"strings"
	"utils/commonDefs"
	"utils/ipcutils"

	"git.apache.org/thrift.git/lib/go/thrift"
)

//var softswitch *vxlan_linux.VxlanLinux

type VXLANClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

/*
func ConvertVxlanConfigToVxlanLinuxConfig(c *VxlanConfig) *vxlan_linux.VxlanConfig {

	return &vxlan_linux.VxlanConfig{
		VNI:    c.VNI,
		VlanId: c.VlanId,
		Group:  c.Group,
		MTU:    c.MTU,
	}
}


func ConvertVtepToVxlanLinuxConfig(vtep *VtepDbEntry) *vxlan_linux.VtepConfig {
	return &vxlan_linux.VtepConfig{
		Vni:          vtep.Vni,
		VtepName:     vtep.VtepName,
		SrcIfName:    vtep.SrcIfName,
		UDP:          vtep.UDP,
		TTL:          vtep.TTL,
		TunnelSrcIp:  vtep.SrcIp,
		TunnelDstIp:  vtep.DstIp,
		VlanId:       vtep.VlanId,
		TunnelSrcMac: vtep.SrcMac,
		TunnelDstMac: vtep.DstMac,
	}
}
*/

// look up the various other daemons based on c string
func GetClientPort(paramsFile string, c string) int {
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		return 0
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		return 0
	}

	for _, client := range clientsList {
		if client.Name == c {
			return client.Port
		}
	}
	return 0
}

func (intf VXLANSnapClient) DisconnectFromClients(name string) error {
	clientList := [3]string{"asicd", "ribd", "arpd"}
	if name != "" {
		clientList = [3]string{name, "", ""}
	}
	for i := 0; i < len(clientList); i++ {
		if asicdclnt.IsConnected && name == "asicd" {
			intf.deleteASICdSubscriber()
			asicdclnt.PtrProtocolFactory = nil
			if err := asicdclnt.Transport.Close(); err != nil {
				logger.Err("Error closing thrift channel to ASICD", err)
				return err
			}
			asicdclnt.IsConnected = false
			asicdclnt.ClientHdl = nil
			asicdclnt.Transport = nil
			asicdclnt.PtrProtocolFactory = nil

			logger.Info("VXLAN -> ASICD diconnected")
		} else if ribdclnt.IsConnected && name == "ribd" {
			intf.deleteRIBdSubscriber()
			ribdclnt.PtrProtocolFactory = nil
			if err := ribdclnt.Transport.Close(); err != nil {
				logger.Err("Error closing thrift channel to RIBD", err)
				return err
			}
			ribdclnt.IsConnected = false
			ribdclnt.ClientHdl = nil
			ribdclnt.Transport = nil
			ribdclnt.PtrProtocolFactory = nil

			logger.Info("VXLAN -> RIBD diconnected")
		} else if arpdclnt.IsConnected && name == "arpd" {
			arpdclnt.PtrProtocolFactory = nil
			if err := arpdclnt.Transport.Close(); err != nil {
				logger.Err("Error closing thrift channel to ARPD", err)
				return err
			}
			arpdclnt.IsConnected = false
			arpdclnt.ClientHdl = nil
			arpdclnt.Transport = nil
			arpdclnt.PtrProtocolFactory = nil

			logger.Info("VXLAN -> ARPD diconnected")
		}
	}
	return nil
}

// ConnectToClients:
// connect the clients to which vxland will need send/receive information
// For this client we will need to connect to asicd, ribd, and arpd
// 1) asicd - provision: vtep/vxlan,
//           notifications: vlan port membership, link up/down
// 2) ribd - provision: next hop ip retreival
//           notifications: next hop reachability changes
// 3) arpd - provision: resolve next hop ip
// return nil means all clients are connected
func (intf VXLANSnapClient) ConnectToClients(clientFile string, name string) error {
	clientList := [3]string{"asicd", "ribd", "arpd"}
	if name != "" {
		clientList = [3]string{name, "", ""}
	}
	//for allclientsconnect < len(clientList) {
	for _, client := range clientList {
		port := GetClientPort(clientFile, client)
		//logger.Info(fmt.Sprintf("VXLAN -> looking to file %s connect %s port %d isconnected asicd[%t] ribd[%t] arpd[%t]",
		//	clientFile, client, port, asicdclnt.IsConnected, ribdclnt.IsConnected, arpdclnt.IsConnected))

		if !asicdclnt.IsConnected && client == "asicd" {
			asicdclnt.Address = "localhost:" + strconv.Itoa(port)
			asicdclnt.Transport, asicdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(asicdclnt.Address)
			if asicdclnt.Transport != nil && asicdclnt.PtrProtocolFactory != nil {
				asicdclnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(asicdclnt.Transport, asicdclnt.PtrProtocolFactory)
				asicdclnt.IsConnected = true
				// lets gather all info needed from asicd such as the port
				logger.Info("VXLAN -> ASICD connected")
				// need to listen for por vlan membership notifications
				intf.createASICdSubscriber()
			}
		} else if !ribdclnt.IsConnected && client == "ribd" {
			ribdclnt.Address = "localhost:" + strconv.Itoa(port)
			ribdclnt.Transport, ribdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(ribdclnt.Address)
			if ribdclnt.Transport != nil && asicdclnt.PtrProtocolFactory != nil {
				ribdclnt.ClientHdl = ribd.NewRIBDServicesClientFactory(ribdclnt.Transport, ribdclnt.PtrProtocolFactory)
				ribdclnt.IsConnected = true
				logger.Info("VXLAN -> RIBD connected")
				// need to listen for next hop notifications
				intf.createRIBdSubscriber()
			}
		} else if !arpdclnt.IsConnected && client == "arpd" {
			arpdclnt.Address = "localhost:" + strconv.Itoa(port)
			arpdclnt.Transport, arpdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(arpdclnt.Address)
			if arpdclnt.Transport != nil && arpdclnt.PtrProtocolFactory != nil {
				arpdclnt.ClientHdl = arpd.NewARPDServicesClientFactory(arpdclnt.Transport, arpdclnt.PtrProtocolFactory)
				arpdclnt.IsConnected = true
				logger.Info("VXLAN -> ARPD connected")
			}
		}
	}

	if arpdclnt.IsConnected &&
		ribdclnt.IsConnected &&
		asicdclnt.IsConnected {
		logger.Info("All Clients connected")

		// lets setup the port map needed for configuration
		intf.ConstructPortConfigMap()

	} else {
		logger.Info(fmt.Sprintf("Not all Clients are connected asicd connected[%t] ribd connected[%t] arpd connected[%t",
			asicdclnt.IsConnected, ribdclnt.IsConnected, asicdclnt.IsConnected))
		return errors.New(fmt.Sprintf("Not all Clients are connected asicd connected[%t] ribd connected[%t] arpd connected[%t",
			asicdclnt.IsConnected, ribdclnt.IsConnected, asicdclnt.IsConnected))

	}
	logger.Info("Returning nil")
	return nil
}

func asicDGetLoopbackInfo() (success bool, lbname string, mac net.HardwareAddr, ip net.IP) {
	// TODO this logic only assumes one loopback interface.  More logic is needed
	// to handle multiple  loopbacks configured.  The idea should be
	// that the lowest IP address is used.
	if asicdclnt.ClientHdl != nil {
		more := true
		for more {
			currMarker := asicdServices.Int(0)
			bulkInfo, err := asicdclnt.ClientHdl.GetBulkLogicalIntfState(currMarker, 5)
			if err == nil {
				objCount := int(bulkInfo.Count)
				more = bool(bulkInfo.More)
				currMarker = asicdServices.Int(bulkInfo.EndIdx)
				for i := 0; i < objCount; i++ {
					ifindex := bulkInfo.LogicalIntfStateList[i].IfIndex
					lbname = bulkInfo.LogicalIntfStateList[i].Name
					if pluginCommon.GetTypeFromIfIndex(ifindex) == commonDefs.IfTypeLoopback {
						mac, _ = net.ParseMAC(bulkInfo.LogicalIntfStateList[i].SrcMac)
						ipV4ObjMore := true
						ipV4ObjCurrMarker := asicdServices.Int(0)
						for ipV4ObjMore {
							ipV4BulkInfo, _ := asicdclnt.ClientHdl.GetBulkIPv4IntfState(ipV4ObjCurrMarker, 20)
							ipV4ObjCount := int(ipV4BulkInfo.Count)
							ipV4ObjCurrMarker = asicdServices.Int(bulkInfo.EndIdx)
							ipV4ObjMore = bool(ipV4BulkInfo.More)
							for j := 0; j < ipV4ObjCount; j++ {
								if ipV4BulkInfo.IPv4IntfStateList[j].IfIndex == ifindex {
									success = true
									ip = net.ParseIP(strings.Split(ipV4BulkInfo.IPv4IntfStateList[j].IpAddr, "/")[0])
									return success, lbname, mac, ip
								}
							}
						}
					}
				}
			}
		}
	}
	return success, lbname, mac, ip
}

func asicDLearnFwdDbEntry(mac net.HardwareAddr, vtepName string, ifindex int32) {
	//macstr := mac.String()
	// convert a vxland config to hw config
	//if asicdclnt.ClientHdl != nil {
	//asicdclnt.ClientHdl.DeleteVxlanVtep(ConvertVtepConfigToVxlanAsicdConfig(vtep))
	//}
	/* Add as another interface
	else {
		// run standalone
		if softswitch != nil {
			softswitch.LearnFdbVtep(macstr, vtepName, ifindex)
		}
	}
	*/
}
