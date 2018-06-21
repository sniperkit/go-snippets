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

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"l3/rib/ribdCommonDefs"
	"ribd"
	"ribdInt"
)

type RibdClient struct {
	OspfClientBase
	ClientHdl *ribd.RIBDServicesClient
}

type RouteMdata struct {
	metric uint32
	ipaddr uint32
	mask   uint32
	isDel  bool
}

func (server *OSPFServer) startRibdUpdates() error {
	server.logger.Info("ASBR: Listen for RIBd updates")
	server.listenForRIBUpdates(ribdCommonDefs.PUB_SOCKET_OSPFD_ADDR)

	go server.createRIBSubscriber()
	return nil
}

func (server *OSPFServer) listenForRIBUpdates(address string) error {
	var err error
	if server.ribSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		server.logger.Err(fmt.Sprintln("ERR: Failed to create RIB subscribe socket, error:", err))
		return err
	}

	if err = server.ribSubSocket.Subscribe(""); err != nil {
		server.logger.Err(fmt.Sprintln("ERR: Failed to subscribe to \"\" on RIB subscribe socket, error:", err))
		return err
	}

	if _, err = server.ribSubSocket.Connect(address); err != nil {
		server.logger.Err(fmt.Sprintln("ERR: Failed to connect to RIB publisher socket, address:", address, "error:", err))
		return err
	}

	server.logger.Info(fmt.Sprintln("Connected to RIB publisher at address:", address))
	if err = server.ribSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		server.logger.Err(fmt.Sprintln("ERR: Failed to set the buffer size for RIB publisher socket, error:", err))
		return err
	}
	return nil
}

func (server *OSPFServer) createRIBSubscriber() {
	for {
		server.logger.Info("Read on RIB subscriber socket...")
		ribrxBuf, err := server.ribSubSocket.Recv(0)
		if err != nil {
			server.logger.Err(fmt.Sprintln("ERR: Recv on RIB subscriber socket failed with error:", err))
			server.ribSubSocketErrCh <- err
			continue
		}
		server.ribSubSocketCh <- ribrxBuf
	}
}

func (server *OSPFServer) processRibdNotification(ribrxBuf []byte) {
	server.logger.Info(fmt.Sprintln("ASBR: Ribd notification received."))
	var route ribdCommonDefs.RoutelistInfo

	reader := bytes.NewReader(ribrxBuf)
	decoder := json.NewDecoder(reader)
	msg := ribdCommonDefs.RibdNotifyMsg{}
	for err := decoder.Decode(&msg); err == nil; err = decoder.Decode(&msg) {
		err = json.Unmarshal(msg.MsgBuf, &route)
		if err != nil {
			server.logger.Err("ASBR: Err in processing routes from RIB")
		}
		server.logger.Info(fmt.Sprintln("ASBR: Receive  route, dest:", route.RouteInfo.Ipaddr, "netmask:", route.RouteInfo.Mask, "nexthop:", route.RouteInfo.NextHopIp))
		server.ProcessRibdRoutes(route.RouteInfo, msg.MsgType)
	}

}

/* @fn getRibdRoutes
Getbulk for RIBD routes before listening to RIBD updates.
*/
func (server *OSPFServer) getRibdRoutes() {
	var currMarker ribdInt.Int
	var count ribdInt.Int
	count = 100
	for {
		server.logger.Info(fmt.Sprintln("ASBR: Getting ", count, " objects from currMarker", currMarker))
		if server.ribdClient.ClientHdl == nil {
			server.logger.Err(fmt.Sprintln("RIBD: Null ribdClient "))
			return
		}
		getBulkInfo, err := server.ribdClient.ClientHdl.GetBulkRoutesForProtocol("OSPF", currMarker, count)
		if err != nil {
			server.logger.Info(fmt.Sprintln("GetBulkRoutesForProtocol with err ", err))
			return
		}
		if getBulkInfo.Count == 0 {
			server.logger.Info("ASBR: 0 objects returned from GetBulkRoutesForProtocol")
			return
		}
		server.logger.Info(fmt.Sprintln("ASBR: len(getBulkInfo.RouteList)  = ", len(getBulkInfo.RouteList), " num objects returned = ", getBulkInfo.Count))

		for _, route := range getBulkInfo.RouteList {

			server.logger.Info(fmt.Sprintln("Receive  route, dest:", route.DestNetIp, "netmask:", route.Mask, "nexthop:", route.NextHopIp))
			server.ProcessRibdRoutes(*route, ribdCommonDefs.NOTIFY_ROUTE_CREATED)
		}
		if getBulkInfo.More == false {
			server.logger.Info("more returned as false, so no more get bulks")
			return
		}
		currMarker = ribdInt.Int(getBulkInfo.EndIdx)
	}

}

/*@fn ProcessRibdRoutes
Send notif to LSDB to generate/delete AS external LSA
*/
func (server *OSPFServer) ProcessRibdRoutes(route ribdInt.Routes, msgType uint16) {
	server.logger.Info("ASBR: Process Ribd routes. msg ")
	ipaddr := convertAreaOrRouterIdUint32(route.Ipaddr)
	mask := convertAreaOrRouterIdUint32(route.Mask)
	metric := uint32(route.Metric)
	isDel := false
	if msgType == ribdCommonDefs.NOTIFY_ROUTE_DELETED {
		isDel = true
	}
	routemdata := RouteMdata{
		ipaddr: ipaddr,
		mask:   mask,
		metric: metric,
		isDel:  isDel,
	}
	ignore := server.verifyOspfRoute(ipaddr, mask)
	if !ignore {
		server.logger.Info(fmt.Sprintln("ASBR: Generate As external for ", route.Ipaddr, route.Mask))
		/* send message to LSDB to generate AS ext LSA */
		server.ExternalRouteNotif <- routemdata
	}
}

/*@fn verifyOspfRoute
Verify if the RIBD route exists in the OSPF routes
*/
func (server *OSPFServer) verifyOspfRoute(ipaddr uint32, mask uint32) bool {
	/* HACK */
	if ipaddr == 855703813 {
		return false
	}
	for key, _ := range server.IntfConfMap {
		intf, _ := server.IntfConfMap[key]
		ip_str := intf.IfIpAddr.String()
		ip := convertAreaOrRouterIdUint32(ip_str)
		server.logger.Info(fmt.Sprintln("ASBR: verify OSPF routes  ip ", ip, " mask ", mask, " ipaddr ", ipaddr))
		if ip == ipaddr {
			server.logger.Info(fmt.Sprintln("ASBR: Ignore route from RIB as ospf is configured to the IF ", ipaddr))
			return true
		}
	}
	return false
}
