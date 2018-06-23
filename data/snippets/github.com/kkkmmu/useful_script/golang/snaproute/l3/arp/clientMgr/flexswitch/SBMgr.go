//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//
package flexswitch

import (
	"arpd"
	"arpdInt"
	"fmt"
	"strconv"
	"time"
	"utils/commonDefs"
	"utils/ipcutils"
)

type ArpdClient struct {
	commonDefs.ClientBase
	ClientHandle *arpd.ARPDServicesClient
}

type FSArpdClientMgr struct {
	ClientHandle *arpd.ARPDServicesClient
}

func GetArpdThriftClientHdl(paramsFile string, clients []commonDefs.ClientJson) *arpd.ARPDServicesClient {
	var clnt ArpdClient
	var err error
	fmt.Println("Connecting to Arp Client", paramsFile)
	for _, client := range clients {
		if client.Name == "arpd" {
			fmt.Println("found:", client.Name, "at port", client.Port)
			clnt.Address = "localhost:" + strconv.Itoa(client.Port)
			clnt.Transport, clnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(clnt.Address)
			if err != nil {
				fmt.Println("Failed to connect to Asicd, retrying until connection is successful")
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					clnt.Transport, clnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(clnt.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 50) == 0 {
						fmt.Println("Still can't connect to:", client.Name, "retrying..")
					}
				}
			}
			fmt.Println("Connected to", client.Name)
			clnt.ClientHandle = arpd.NewARPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
			return clnt.ClientHandle
		}
	}
	return nil
}

func (p *FSArpdClientMgr) ResolveArpIPV4(destNetIp string, vlanid int32) (err error) {
	return p.ClientHandle.ResolveArpIPv4(destNetIp, arpdInt.Int(vlanid))
}

func (p *FSArpdClientMgr) DeleteResolveArpIPv4(NbrIP string) (err error) {
	return p.ClientHandle.DeleteResolveArpIPv4(NbrIP)
}

func (p *FSArpdClientMgr) DeleteArpEntry(ipAddr string) (err error) {
	return p.ClientHandle.DeleteArpEntry(ipAddr)
}

func (p *FSArpdClientMgr) SendGarp(ifName, macAddr, ipAddr string) (err error) {
	return p.ClientHandle.SendGarp(ifName, macAddr, ipAddr)
}
