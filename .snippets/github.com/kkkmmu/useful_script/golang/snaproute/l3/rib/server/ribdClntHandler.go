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
	"arpd"
	"asicdServices"
	"git.apache.org/thrift.git/lib/go/thrift"
	"infra/sysd/sysdCommonDefs"
	"l3/rib/ribdCommonDefs"
	"utils/keepalive"
	"utils/patriciaDB"
)

type RIBClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

type AsicdClient struct {
	baseClient
	RIBClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

type ArpdClient struct {
	baseClient
	RIBClientBase
	ClientHdl *arpd.ARPDServicesClient
}
type BGPdClient struct {
	baseClient
}
type OSPFdClient struct {
	baseClient
}
type ClientIf interface {
	DmnDownHandler()
	DmnUpHandler()
	ConnectToClient()
}

type baseClient struct {
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

var asicdclnt AsicdClient
var arpdclnt ArpdClient
var bgpdclnt BGPdClient
var ospfdclnt OSPFdClient

func deleteV4RoutesOfType(protocol string, destNet string) {
	var testroutes []RouteInfoRecord
	testroutes = make([]RouteInfoRecord, 0)

	routeInfoRecordListItem := V4RouteInfoMap.Get(patriciaDB.Prefix(destNet))
	if routeInfoRecordListItem == nil {
		logger.Info("Unexpected: no route for destNet:", destNet, " found in routeMap")
		return
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	protocolRouteList, ok := routeInfoRecordList.routeInfoProtocolMap[protocol]
	if !ok || len(protocolRouteList) == 0 {
		logger.Info("Unexpected: no route for destNet:", destNet, " found in routeMap of type:", protocol)
		return
	}
	for _, testroute := range protocolRouteList {
		//logger.Info("will call delete for route with ip:", testroute.destNetIp.String(), " nexthop:", testroute.nextHopIp.String())
		testroutes = append(testroutes, testroute)
	}
	//logger.Info("found ", len(testroutes), " number of ", protocol, " routes in routemap:", testroutes)
	for _, protoroute := range testroutes { //protocolRouteList {
		//logger.Info(len(testroutes), " number of ", protocol, " routes in routemap:", testroutes, " remaining")
		//logger.Info("protoroute:", protoroute, " nexthop:", protoroute.nextHopIp.String())
		_, err := deleteIPRoute(protoroute.destNetIp.String(), ribdCommonDefs.IPv4, protoroute.networkMask.String(), protocol, protoroute.nextHopIp.String(), protoroute.nextHopIfIndex, FIBAndRIB, ribdCommonDefs.RoutePolicyStateChangetoInValid)
		logger.Info("err :", err, " while deleting ", protocol, " route with destNet:", protoroute.destNetIp.String(), " nexthopIP:", protoroute.nextHopIp.String())
	}
}
func deleteV6RoutesOfType(protocol string, destNet string) {
	var testroutes []RouteInfoRecord
	testroutes = make([]RouteInfoRecord, 0)

	routeInfoRecordListItem := V6RouteInfoMap.Get(patriciaDB.Prefix(destNet))
	if routeInfoRecordListItem == nil {
		logger.Info("Unexpected: no route for destNet:", destNet, " found in routeMap")
		return
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	protocolRouteList, ok := routeInfoRecordList.routeInfoProtocolMap[protocol]
	if !ok || len(protocolRouteList) == 0 {
		logger.Info("Unexpected: no route for destNet:", destNet, " found in routeMap of type:", protocol)
		return
	}
	for _, testroute := range protocolRouteList {
		//logger.Info("will call delete for route with ip:", testroute.destNetIp.String(), " nexthop:", testroute.nextHopIp.String())
		testroutes = append(testroutes, testroute)
	}
	//logger.Info("found ", len(testroutes), " number of ", protocol, " routes in routemap:", testroutes)
	for _, protoroute := range testroutes { //protocolRouteList {
		//logger.Info(len(testroutes), " number of ", protocol, " routes in routemap:", testroutes, " remaining")
		//logger.Info("protoroute:", protoroute, " nexthop:", protoroute.nextHopIp.String())
		_, err := deleteIPRoute(protoroute.destNetIp.String(), ribdCommonDefs.IPv6, protoroute.networkMask.String(), protocol, protoroute.nextHopIp.String(), protoroute.nextHopIfIndex, FIBAndRIB, ribdCommonDefs.RoutePolicyStateChangetoInValid)
		logger.Info("err :", err, " while deleting ", protocol, " route with destNet:", protoroute.destNetIp.String(), " nexthopIP:", protoroute.nextHopIp.String())
	}
}
func DeleteRoutesOfType(protocol string) {
	func_mesg := "DeleteRoutesOfType of type:" + protocol
	protocolRouteMap, ok := ProtocolRouteMap[protocol]
	if !ok {
		logger.Info(func_mesg, "No routes of ", protocol, " type configured")
		return
	}
	if protocolRouteMap.v4routeMap != nil {
		logger.Info(func_mesg, " number of v4 routes:", len(protocolRouteMap.v4routeMap))
		for destNet, count := range protocolRouteMap.v4routeMap {
			if count.totalcount > 0 {
				logger.Info(func_mesg, ":", count, " number of v4 routes for destNet IP:", string(destNet))
				deleteV4RoutesOfType(protocol, destNet)
				//deleteV6RoutesOfType(protocol, destNet)
				protocolRouteMap.totalcount.totalcount = protocolRouteMap.totalcount.totalcount - count.totalcount
				protocolRouteMap.totalcount.ecmpcount = protocolRouteMap.totalcount.ecmpcount - count.ecmpcount
				totalCount := protocolRouteMap.v4routeMap[destNet]
				totalCount.totalcount = 0
				totalCount.ecmpcount = 0
				protocolRouteMap.v4routeMap[destNet] = totalCount
				//			protocolRouteMap.routeMap[destNet].ecmpcount = 0
				ProtocolRouteMap[protocol] = protocolRouteMap
			}
		}
	}
	if protocolRouteMap.v6routeMap != nil {
		logger.Info(func_mesg, " number of v6 routes:", len(protocolRouteMap.v6routeMap))
		for destNet, count := range protocolRouteMap.v6routeMap {
			if count.totalcount > 0 {
				logger.Info(count, " number of v6 routes for destNet IP:", string(destNet))
				//deleteV4RoutesOfType(protocol, destNet)
				deleteV6RoutesOfType(protocol, destNet)
				protocolRouteMap.totalcount.totalcount = protocolRouteMap.totalcount.totalcount - count.totalcount
				protocolRouteMap.totalcount.ecmpcount = protocolRouteMap.totalcount.ecmpcount - count.ecmpcount
				totalCount := protocolRouteMap.v6routeMap[destNet]
				totalCount.totalcount = 0
				totalCount.ecmpcount = 0
				protocolRouteMap.v6routeMap[destNet] = totalCount
				//			protocolRouteMap.routeMap[destNet].ecmpcount = 0
				ProtocolRouteMap[protocol] = protocolRouteMap
			}
		}
	}
}

//Daemon DOWN handler functions
func (clnt *ArpdClient) DmnDownHandler() {
	logger.Info("DmnDownHandler for ArpdClient")
	clnt.IsConnected = false
}
func (clnt *AsicdClient) DmnDownHandler() {
	logger.Info("DmnDownHandler for AsicdClient")
	clnt.IsConnected = false
}
func (clnt *baseClient) DmnDownHandler() {
	logger.Info("DmnDownHandler for baseClient")
}
func (clnt *BGPdClient) DmnDownHandler() {
	logger.Info("DmnDownHandler for BGPd")
	//uninstall all BGP routes
	DeleteRoutesOfType("EBGP")
}
func (clnt *OSPFdClient) DmnDownHandler() {
	logger.Info("DmnDownHandler for OSPFd")
	//uninstall all OSPF routes
	DeleteRoutesOfType("OSPF")
}
func (mgr *RIBDServer) DmnDownHandler(name string) error {
	logger.Info("In DmnDownHandler call DmnDownHandler for client: ", name)
	client, exist := mgr.Clients[name]
	if exist {
		client.DmnDownHandler()
	}
	return nil
}

//Daemon UP handler functions
func (clnt *ArpdClient) DmnUpHandler() {
	logger.Info("DmnUpHandler for ArpdClient")
	if arpdclnt.IsConnected {
		logger.Info("RIBD already connected to arpd")
		return
	}
	//	go clnt.ConnectToClient()
}
func (clnt *AsicdClient) DmnUpHandler() {
	logger.Info("DmnUpHandler for AsicdClient")
	if asicdclnt.IsConnected {
		logger.Info("RIBD already connected to asicd")
		return
	}
	//	go clnt.ConnectToClient()
}
func (clnt *BGPdClient) DmnUpHandler() {
	logger.Info("DmnUpHandler for BGPd")
	//no op here since BGP calls GetBulkRoutesForProtocol
}
func (clnt *baseClient) DmnUpHandler() {
	logger.Info("DmnUpHandler for baseClient")
}
func (mgr *RIBDServer) DmnUpHandler(name string) error {
	client, exist := mgr.Clients[name]
	if exist {
		client.DmnUpHandler()
	}
	return nil
}

func (mgr *RIBDServer) ListenToClientStateChanges() {
	logger.Info("ListenToClientStateChanges")
	clientStatusListener := keepalive.InitDaemonStatusListener()
	if clientStatusListener != nil {
		go clientStatusListener.StartDaemonStatusListner()
		for {
			select {
			case clientStatus := <-clientStatusListener.DaemonStatusCh:
				logger.Info("Received client status: ", clientStatus.Name, clientStatus.Status)
				switch clientStatus.Status {
				case sysdCommonDefs.STOPPED, sysdCommonDefs.RESTARTING:
					logger.Info(clientStatus.Name, " stopped or restarting")
					mgr.DmnDownHandler(clientStatus.Name)
				case sysdCommonDefs.UP:
					logger.Info(clientStatus.Name, " up now")
					mgr.DmnUpHandler(clientStatus.Name)
				}
			}
		}
	}
}

/*
//connect to client functions
func (clnt *ArpdClient) ConnectToClient() {
	var timer *time.Timer
	logger.Info("in go routine ConnectToClient for connecting to ARPd")
	for {
		logger.Info("in for loop of go routine ConnectToClient for connecting to ARPd")
		timer = time.NewTimer(time.Second * 1)
		<-timer.C
		logger.Info("Connecting to arpd at address ", arpdclnt.Address)
		//arpdclnt.Address = "localhost:" + strconv.Itoa(port)
		arpdclnt.Transport, arpdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(arpdclnt.Address)
		if arpdclnt.Transport != nil && arpdclnt.PtrProtocolFactory != nil {
			logger.Info("connecting to arpd,asicdclnt.IsConnected:", asicdclnt.IsConnected)
			arpdclnt.ClientHdl = arpd.NewARPDServicesClientFactory(arpdclnt.Transport, arpdclnt.PtrProtocolFactory)
			arpdclnt.IsConnected = true
			RouteServiceHandler.Clients["arpd"] = &arpdclnt
			if asicdclnt.IsConnected == true {
				logger.Info(" Connected to all clients: call AcceptConfigActions")
				RouteServiceHandler.AcceptConfigActions()
			}
			timer.Stop()
			return
		}
	}
}
func (clnt *AsicdClient) ConnectToClient() {
	var timer *time.Timer
	logger.Info("in go routine ConnectToClient for connecting to ASICd")
	for {
		timer = time.NewTimer(time.Second * 10)
		<-timer.C
		logger.Info("Connecting to asicd at address ", asicdclnt.Address)
		//asicdclnt.Address = "localhost:" + strconv.Itoa(port)
		asicdclnt.Transport, asicdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(asicdclnt.Address)
		if asicdclnt.Transport != nil && asicdclnt.PtrProtocolFactory != nil {
			logger.Info("connecting to asicd,arpdclnt.IsConnected:", arpdclnt.IsConnected)
			asicdclnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(asicdclnt.Transport, asicdclnt.PtrProtocolFactory)
			asicdclnt.IsConnected = true
			RouteServiceHandler.Clients["asicd"] = &asicdclnt
			if arpdclnt.IsConnected == true {
				logger.Info(" Connected to all clients: call AcceptConfigActions")
				RouteServiceHandler.AcceptConfigActions()
			}
			timer.Stop()
			return
		}
	}
}
func (clnt *baseClient) ConnectToClient() {
}


func (ribdServiceHandler *RIBDServer) connectToClient(name string) {
	var timer *time.Timer
	logger.Info("in go routine ConnectToClient for connecting to %s\n", name)
	for {
		timer = time.NewTimer(time.Second * 10)
		<-timer.C
		if name == "asicd" {
			logger.Info("Connecting to asicd at address ", asicdclnt.Address)
			//asicdclnt.Address = "localhost:" + strconv.Itoa(port)
			asicdclnt.Transport, asicdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(asicdclnt.Address)
			if asicdclnt.Transport != nil && asicdclnt.PtrProtocolFactory != nil {
				logger.Info("connecting to asicd,arpdclnt.IsConnected:", arpdclnt.IsConnected)
				asicdclnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(asicdclnt.Transport, asicdclnt.PtrProtocolFactory)
				asicdclnt.IsConnected = true
				ribdServiceHandler.Clients["asicd"] = asicdclnt
				if arpdclnt.IsConnected == true {
					logger.Info(" Connected to all clients: call AcceptConfigActions")
					ribdServiceHandler.AcceptConfigActions()
				}
				timer.Stop()
				return
			}
		}
		if name == "arpd" {
			logger.Info("Connecting to arpd at address ", arpdclnt.Address)
			//arpdclnt.Address = "localhost:" + strconv.Itoa(port)
			arpdclnt.Transport, arpdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(arpdclnt.Address)
			if arpdclnt.Transport != nil && arpdclnt.PtrProtocolFactory != nil {
				logger.Info("connecting to arpd,asicdclnt.IsConnected:", asicdclnt.IsConnected)
				arpdclnt.ClientHdl = arpd.NewARPDServicesClientFactory(arpdclnt.Transport, arpdclnt.PtrProtocolFactory)
				arpdclnt.IsConnected = true
				ribdServiceHandler.Clients["arpd"] = arpdclnt
				if asicdclnt.IsConnected == true {
					logger.Info(" Connected to all clients: call AcceptConfigActions")
					ribdServiceHandler.AcceptConfigActions()
				}
				timer.Stop()
				return
			}
		}
	}
}

func (ribdServiceHandler *RIBDServer) ConnectToClients(paramsFile string) {
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		logger.Err("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		logger.Err("Error in Unmarshalling Json")
		return
	}

	for _, client := range clientsList {
		logger.Info("#### Client name is ", client.Name)
		if client.Name == "bgpd" {
			ribdServiceHandler.Clients["bgpd"] = &bgpdclnt
		}
		if client.Name == "ospfd" {
			ribdServiceHandler.Clients["ospfd"] = &ospfdclnt
		}
		/*
			if client.Name == "asicd" {
				logger.Info("found asicd at port ", client.Port)
				asicdclnt.Address = "localhost:" + strconv.Itoa(client.Port)
				asicdclnt.Transport, asicdclnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(asicdclnt.Address)
				if err != nil {
					logger.Info("Failed to connect to Asicd, retrying until connection is successful")
					count := 0
					ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
					for _ = range ticker.C {
						asicdclnt.Transport, asicdclnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(asicdclnt.Address)
						if err == nil {
							ticker.Stop()
							break
						}
						count++
						if (count % 10) == 0 {
							logger.Info("Still can't connect to Asicd, retrying...")
						}
					}
				}
				if asicdclnt.Transport != nil && asicdclnt.PtrProtocolFactory != nil {
					logger.Info("connecting to asicd,arpdclnt.IsConnected:", arpdclnt.IsConnected)
					asicdclnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(asicdclnt.Transport, asicdclnt.PtrProtocolFactory)
					asicdclnt.IsConnected = true
					ribdServiceHandler.Clients["asicd"] = &asicdclnt
					if arpdclnt.IsConnected == true {
						logger.Info(" Connected to all clients: call AcceptConfigActions")
						ribdServiceHandler.AcceptConfigActions()
					}
				} else {
					logger.Info("asicd clnt nil even after err is nil with createipchandles")
					//go ribdServiceHandler.connectToClient(client.Name)
					//go asicdclnt.ConnectToClient()
				}
			}
			if client.Name == "arpd" {
				logger.Info("RIBD: found arpd at port ", client.Port)
				arpdclnt.Address = "localhost:" + strconv.Itoa(client.Port)
				logger.Info("arpdclnt.Address:", arpdclnt.Address)
				arpdclnt.Transport, arpdclnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(arpdclnt.Address)
				logger.Info("arpdclnt.transport:", arpdclnt.Transport, " arpdclnt.ProtocolFactory:", arpdclnt.PtrProtocolFactory, " err:", err)
				if err != nil {
					logger.Info("Failed to connect to Arpd, retrying until connection is successful")
					count := 0
					ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
					for _ = range ticker.C {
						arpdclnt.Transport, arpdclnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(arpdclnt.Address)
						if err == nil {
							ticker.Stop()
							break
						}
						count++
						if (count % 10) == 0 {
							logger.Info("Still can't connect to Arpd, retrying...")
						}
					}
				}
				if arpdclnt.Transport != nil && arpdclnt.PtrProtocolFactory != nil {
					logger.Info("connecting to arpd,asicdclnt.IsConnected:", asicdclnt.IsConnected)
					arpdclnt.ClientHdl = arpd.NewARPDServicesClientFactory(arpdclnt.Transport, arpdclnt.PtrProtocolFactory)
					arpdclnt.IsConnected = true
					ribdServiceHandler.Clients["arpd"] = &arpdclnt
					if asicdclnt.IsConnected == true {
						logger.Info(" Connected to all clients: call AcceptConfigActions")
						ribdServiceHandler.AcceptConfigActions()
					}
				} else {
					logger.Info("arpd clnt nil even after err is nil with createipchandles")
					//go ribdServiceHandler.connectToClient(client.Name)
					//go arpdclnt.ConnectToClient()
				}
			}
	}
}*/
