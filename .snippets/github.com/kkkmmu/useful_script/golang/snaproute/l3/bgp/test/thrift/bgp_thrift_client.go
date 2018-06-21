//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

// bgp_client.go
package main

import (
	"bgpd"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const CONF_IP string = "localhost" //"10.0.2.15"
const CONF_PORT string = "4050"

func main() {
	fmt.Println("Starting the BGP thrift client...")
	var clientTransport thrift.TTransport

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	clientTransport, err := thrift.NewTSocket("localhost:" + CONF_PORT)
	if err != nil {
		fmt.Println("NewTSocket failed with error:", err)
		return
	}

	clientTransport = transportFactory.GetTransport(clientTransport)
	if err = clientTransport.Open(); err != nil {
		fmt.Println("Failed to open the socket, error:", err)
	}

	client := bgpd.NewBGPDServicesClientFactory(clientTransport, protocolFactory)

	globalConfigArgs := bgpd.NewBGPGlobal()
	globalConfigArgs.ASNum = "5000"
	globalConfigArgs.RouterId = "localhost"
	fmt.Println("calling CreateBgpGlobal with attr:", globalConfigArgs)
	ret, err := client.CreateBGPGlobal(globalConfigArgs)
	if !ret {
		fmt.Println("CreateBgpGlobal FAILED, ret:", ret, "err:", err)
	}
	fmt.Println("Created BGP global conf")

	peerConfigArgs := bgpd.NewBGPv4Neighbor()
	peerConfigArgs.NeighborAddress = "11.1.11.203"
	peerConfigArgs.LocalAS = "5000"
	peerConfigArgs.PeerAS = "5000"
	peerConfigArgs.Description = "IBGP Peer"
	fmt.Println("calling CreateBgpPeer with attr:", peerConfigArgs)
	ret, err = client.CreateBGPv4Neighbor(peerConfigArgs)
	if !ret {
		fmt.Println("CreateBgpPeer FAILED, ret:", ret, "err:", err)
	}
	fmt.Println("Created BGP peer conf")

	//	peerCommandArgs := &server.PeerConfigCommands{net.ParseIP("11.1.11.203"), 1}
	//	err = client.Call("ConfigInterface.PeerCommand", peerCommandArgs, &reply)
	//	if err != nil {
	//		fmt.Println("ConfigInterface.AddPeer FAILED with err:", err)
	//	}

}
