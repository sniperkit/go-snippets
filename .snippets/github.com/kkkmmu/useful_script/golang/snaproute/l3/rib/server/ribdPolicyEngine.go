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

// ribdPolicyEngine.go
package server

import (
	"l3/rib/ribdCommonDefs"
	"ribd"
	"ribdInt"
	"strconv"
	"strings"
	netUtils "utils/netUtils"
	"utils/patriciaDB"
	"utils/policy"
	"utils/policy/policyCommonDefs"
)

/*
   Index of policy entity
*/
type PolicyRouteIndex struct {
	destNetIP string //CIDR format
	policy    string
}

/*
   data-structure used to communicate with policy engine
*/
type RouteParams struct {
	ipType         ribdCommonDefs.IPType
	destNetIp      string
	networkMask    string
	nextHopIp      string
	nextHopIfIndex ribd.Int
	metric         ribd.Int
	sliceIdx       ribd.Int
	routeType      ribd.Int
	createType     ribd.Int
	deleteType     ribd.Int
	weight         ribd.Int
	bulk           bool
	bulkEnd        bool
}

type TraverseAndApplyPolicyData struct {
	data       interface{}
	updatefunc policy.PolicyApplyfunc
}

func policyEngineActionRejectRoute(params interface{}) {
	routeInfo := params.(RouteParams)
	logger.Info("policyEngineActionRejectRoute for route ", routeInfo.destNetIp, " ", routeInfo.networkMask)
	cfg := ribd.IPv4Route{
		DestinationNw: routeInfo.destNetIp,
		Protocol:      ReverseRouteProtoTypeMapDB[int(routeInfo.routeType)],
		Cost:          int32(routeInfo.metric),
		NetworkMask:   routeInfo.networkMask,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     routeInfo.nextHopIp,
		NextHopIntRef: strconv.Itoa(int(routeInfo.nextHopIfIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)
	_, err := RouteServiceHandler.ProcessV4RouteDeleteConfig(&cfg, FIBAndRIB) //routeInfo.destNetIp, routeInfo.networkMask, ReverseRouteProtoTypeMapDB[int(routeInfo.routeType)], routeInfo.nextHopIp) // FIBAndRIB)//,ribdCommonDefs.RoutePolicyStateChangetoInValid)
	if err != nil {
		logger.Info("deleting v4 route failed with err ", err)
		return
	}
}

/*func policyEngineActionUndoRejectRoute(conditionsList []string, params interface{}, policyStmt policy.PolicyStmt) {
	routeInfo := params.(RouteParams)
	logger.Info(("policyEngineActionUndoRejectRoute - route: ", routeInfo.destNetIp, ":", routeInfo.networkMask, " type ", routeInfo.routeType))
	var tempRoute ribdInt.Routes
	if routeInfo.routeType == ribdCommonDefs.STATIC {
		logger.Info(("this is a static route, fetch it from the DB"))
		DbName := PARAMSDIR + "/UsrConfDb.db"
		logger.Info(("DB Location: ", DbName))
		dbHdl, err := sql.Open("sqlite3", DbName)
		if err != nil {
			logger.Info(("Failed to create the handle with err ", err))
			return
		}

		if err = dbHdl.Ping(); err != nil {
			logger.Info(("Failed to keep DB connection alive"))
			return
		}
		dbCmd := "select * from IPV4Route"
		rows, err := dbHdl.Query(dbCmd)
		if err != nil {
			logger.Info(fmt.Sprintf("DB Query failed for %s with err %s\n", dbCmd, err))
			return
		}
		var ipRoute IPRoute
		for rows.Next() {
			if err = rows.Scan(&ipRoute.DestinationNw, &ipRoute.NetworkMask, &ipRoute.Cost, &ipRoute.NextHopIp, &ipRoute.OutgoingIntfType, &ipRoute.OutgoingInterface, &ipRoute.Protocol); err != nil {
				logger.Info(fmt.Sprintf("DB Scan failed when iterating over IPV4Route rows with error %s\n", err))
				return
			}
			outIntf, _ := strconv.Atoi(ipRoute.OutgoingInterface)
			var outIntfType ribd.Int
			if ipRoute.OutgoingIntfType == "VLAN" {
				outIntfType = commonDefs.IfTypeVlan
			} else {
				outIntfType = commonDefs.IfTypePort
			}
			proto, _ := strconv.Atoi(ipRoute.Protocol)
			tempRoute.Ipaddr = ipRoute.DestinationNw
			tempRoute.Mask = ipRoute.NetworkMask
			tempRoute.NextHopIp = ipRoute.NextHopIp
			tempRoute.NextHopIfType = ribdInt.Int(outIntfType)
			tempRoute.IfIndex = ribdInt.Int(outIntf)
			tempRoute.Prototype = ribdInt.Int(proto)
			tempRoute.Metric = ribdInt.Int(ipRoute.Cost)

			entity, err := buildPolicyEntityFromRoute(tempRoute, params)
			if err != nil {
				logger.Err(("Error builiding policy entity params"))
				return
			}
			if !PolicyEngineDB.ConditionCheckValid(entity, conditionsList, policyStmt) {
				logger.Info(("This route does not qualify for reversing reject route"))
				continue
			}
			cfg := ribd.IPv4Route{
				DestinationNw:     tempRoute.Ipaddr,
				Protocol:          "STATIC",
				OutgoingInterface: ipRoute.OutgoingInterface,
				OutgoingIntfType:  ipRoute.OutgoingIntfType,
				Cost:              int32(tempRoute.Metric),
				NetworkMask:       tempRoute.Mask,
				NextHopIp:         tempRoute.NextHopIp}

			_, err = RouteServiceHandler.ProcessRouteCreateConfig(&cfg) //tempRoute.Ipaddr, tempRoute.Mask, tempRoute.Metric, tempRoute.NextHopIp, tempRoute.NextHopIfType, tempRoute.IfIndex, "STATIC") //tempRoute.Prototype)
			if err != nil {
				logger.Info(fmt.Sprintf("Route create failed with err %s\n", err))
				return
			}
		}
	} else if routeInfo.routeType == ribdCommonDefs.CONNECTED {
		logger.Info(("this is a connected route, fetch it from ASICD"))
		if !asicdclnt.IsConnected {
			logger.Info(("Not connected to ASICD"))
			return
		}
		var currMarker asicdServices.Int
		var count asicdServices.Int
		count = 100
		for {
			logger.Info(fmt.Sprintf("Getting %d objects from currMarker %d\n", count, currMarker))
			IPIntfBulk, err := asicdclnt.ClientHdl.GetBulkIPv4IntfState(currMarker, count)
			if err != nil {
				logger.Info(("GetBulkIPv4IntfState with err ", err))
				return
			}
			if IPIntfBulk.Count == 0 {
				logger.Info(("0 objects returned from GetBulkIPv4IntfState"))
				return
			}
			logger.Info(fmt.Sprintf("len(IPIntfBulk.IPv4IntfStateList)  = %d, num objects returned = %d\n", len(IPIntfBulk.IPv4IntfStateList), IPIntfBulk.Count))
			for i := 0; i < int(IPIntfBulk.Count); i++ {
				var ipMask net.IP
				ip, ipNet, err := net.ParseCIDR(IPIntfBulk.IPv4IntfStateList[i].IpAddr)
				if err != nil {
					return
				}
				ipMask = make(net.IP, 4)
				copy(ipMask, ipNet.Mask)
				ipAddrStr := ip.String()
				ipMaskStr := net.IP(ipMask).String()
				tempRoute.Ipaddr = ipAddrStr
				tempRoute.Mask = ipMaskStr
				tempRoute.NextHopIp = "0.0.0.0"
				tempRoute.NextHopIfType = ribdInt.Int(asicdCommonDefs.GetIntfTypeFromIfIndex(IPIntfBulk.IPv4IntfStateList[i].IfIndex))
				nextHopIfTypeStr := ""
				switch tempRoute.NextHopIfType {
				case commonDefs.IfTypePort:
					nextHopIfTypeStr = "PHY"
					break
				case commonDefs.IfTypeVlan:
					nextHopIfTypeStr = "VLAN"
					break
				case commonDefs.IfTypeNull:
					nextHopIfTypeStr = "NULL"
					break
				}
				tempRoute.IfIndex = ribdInt.Int(asicdCommonDefs.GetIntfIdFromIfIndex(IPIntfBulk.IPv4IntfStateList[i].IfIndex))
				tempRoute.Prototype = ribdCommonDefs.CONNECTED
				tempRoute.Metric = 0
				entity, err := buildPolicyEntityFromRoute(tempRoute, params)
				if err != nil {
					logger.Err(("Error builiding policy entity params"))
					return
				}
				if !PolicyEngineDB.ConditionCheckValid(entity, conditionsList, policyStmt) {
					logger.Info(("This route does not qualify for reversing reject route"))
					continue
				}
				logger.Info(fmt.Sprintf("Calling createv4Route with ipaddr %s mask %s\n", ipAddrStr, ipMaskStr))
				cfg := ribd.IPv4Route{
					DestinationNw:     tempRoute.Ipaddr,
					Protocol:          "CONNECTED",
					OutgoingInterface: strconv.Itoa(int(tempRoute.IfIndex)),
					OutgoingIntfType:  nextHopIfTypeStr,
					Cost:              0,
					NetworkMask:       tempRoute.Mask,
					NextHopIp:         "0.0.0.0"}
				_, err = RouteServiceHandler.ProcessRouteCreateConfig(&cfg) //ipAddrStr, ipMaskStr, 0, "0.0.0.0", ribd.Int(asicdCommonDefs.GetIntfTypeFromIfIndex(IPIntfBulk.IPv4IntfList[i].IfIndex)), ribd.Int(asicdCommonDefs.GetIntfIdFromIfIndex(IPIntfBulk.IPv4IntfList[i].IfIndex)), "CONNECTED") // FIBAndRIB, ribd.Int(len(destNetSlice)))
				if err != nil {
					logger.Info(fmt.Sprintf("Failed to create connected route for ip Addr %s/%s intfType %d intfId %d\n", ipAddrStr, ipMaskStr, ribd.Int(asicdCommonDefs.GetIntfTypeFromIfIndex(IPIntfBulk.IPv4IntfStateList[i].IfIndex)), ribd.Int(asicdCommonDefs.GetIntfIdFromIfIndex(IPIntfBulk.IPv4IntfStateList[i].IfIndex))))
				}
			}
			if IPIntfBulk.More == false {
				logger.Info(("more returned as false, so no more get bulks"))
				return
			}
			currMarker = asicdServices.Int(IPIntfBulk.EndIdx)
		}
	}
}*/
func policyEngineUndoRouteDispositionAction(action interface{}, conditionList []interface{},
	policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineUndoRouteDispositionAction")
	if action.(string) == "Reject" {
		logger.Info("Reject action")
		conditionNameList := make([]string, len(conditionList))
		for i := 0; i < len(conditionList); i++ {
			condition := conditionList[i].(policy.PolicyCondition)
			conditionNameList[i] = condition.Name
		}
		//policyEngineActionUndoRejectRoute(conditionNameList, params, policyStmt)
	} else if action.(string) == "Accept" {
		policyEngineActionRejectRoute(params)
	}
}
func policyEngineActionUndoNetworkStatemenAdvertiseAction(actionItem interface{}, conditionsList []interface{}, policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineActionUndoNetworkStatemenAdvertiseAction")
	RouteInfo := params.(RouteParams)
	var route ribdInt.Routes
	networkStatementTargetProtocol := actionItem.(string)
	//Send a event based on target protocol
	var evt int
	evt = ribdCommonDefs.NOTIFY_ROUTE_DELETED
	switch RouteProtocolTypeMapDB[networkStatementTargetProtocol] {
	case ribdCommonDefs.BGP:
		logger.Info("Undo network statement advertise to BGP")
		route = ribdInt.Routes{Ipaddr: RouteInfo.destNetIp, Mask: RouteInfo.networkMask, NextHopIp: RouteInfo.nextHopIp, IPAddrType: ribdInt.Int(RouteInfo.ipType), IfIndex: ribdInt.Int(RouteInfo.nextHopIfIndex), Metric: ribdInt.Int(RouteInfo.metric), Prototype: ribdInt.Int(RouteInfo.routeType)}
		route.NetworkStatement = true
		publisherInfo, ok := PublisherInfoMap["BGP"]
		if ok {
			RedistributionNotificationSend(publisherInfo.pub_socket, route, evt, networkStatementTargetProtocol)
		}
		break
	default:
		logger.Info("Unknown target protocol")
	}
	UpdateRedistributeTargetMap(evt, networkStatementTargetProtocol, route)
}
func policyEngineActionUndoRedistribute(actionItem interface{}, conditionsList []interface{}, policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineActionUndoRedistribute")
	RouteInfo := params.(RouteParams)
	var route ribdInt.Routes
	redistributeActionInfo := actionItem.(policy.RedistributeActionInfo)
	//Send a event based on target protocol
	var evt int
	logger.Info("redistributeAction set to ", redistributeActionInfo.Redistribute)
	if redistributeActionInfo.Redistribute == true {
		logger.Info("evt = NOTIFY_ROUTE_DELETED")
		evt = ribdCommonDefs.NOTIFY_ROUTE_DELETED
	} else {
		logger.Info("evt = NOTIFY_ROUTE_CREATED")
		evt = ribdCommonDefs.NOTIFY_ROUTE_CREATED
	}
	route = ribdInt.Routes{Ipaddr: RouteInfo.destNetIp, Mask: RouteInfo.networkMask, NextHopIp: RouteInfo.nextHopIp, IPAddrType: ribdInt.Int(RouteInfo.ipType), IfIndex: ribdInt.Int(RouteInfo.nextHopIfIndex), Metric: ribdInt.Int(RouteInfo.metric), Prototype: ribdInt.Int(RouteInfo.routeType)}
	route.RouteOrigin = ReverseRouteProtoTypeMapDB[int(RouteInfo.routeType)]
	publisherInfo, ok := PublisherInfoMap[redistributeActionInfo.RedistributeTargetProtocol]
	if ok {
		logger.Info("ReditributeNotificationSend event called for target protocol - ", redistributeActionInfo.RedistributeTargetProtocol)
		RedistributionNotificationSend(publisherInfo.pub_socket, route, evt, redistributeActionInfo.RedistributeTargetProtocol)
	} else {
		logger.Info("Unknown target protocol")
	}
	UpdateRedistributeTargetMap(evt, redistributeActionInfo.RedistributeTargetProtocol, route)
}
func policyEngineUpdateRoute(prefix patriciaDB.Prefix, item patriciaDB.Item, handle patriciaDB.Item) (err error) {
	logger.Info("policyEngineUpdateRoute for ", prefix)

	rmapInfoRecordList := item.(RouteInfoRecordList)
	if rmapInfoRecordList.routeInfoProtocolMap == nil {
		logger.Info("No routes configured for this prefix")
		return err
	}
	routeInfoList := rmapInfoRecordList.routeInfoProtocolMap[rmapInfoRecordList.selectedRouteProtocol]
	if len(routeInfoList) == 0 {
		logger.Info("len(routeInfoList) == 0")
		return err
	}
	logger.Info("Selected route protocol = ", rmapInfoRecordList.selectedRouteProtocol)
	selectedRouteInfoRecord := routeInfoList[0]
	//route := ribdInt.Routes{Ipaddr:selectedRouteInfoRecord.destNetIp.String() , Mask: selectedRouteInfoRecord.networkMask.String(), NextHopIp: selectedRouteInfoRecord.nextHopIp.String(), NextHopIfType: ribdInt.Int(selectedRouteInfoRecord.nextHopIfType), IfIndex: ribdInt.Int(selectedRouteInfoRecord.nextHopIfIndex), Metric: ribdInt.Int(selectedRouteInfoRecord.metric), Prototype: ribdInt.Int(selectedRouteInfoRecord.protocol), IsPolicyBasedStateValid: rmapInfoRecordList.isPolicyBasedStateValid}
	nextHopIf := strconv.Itoa(int(selectedRouteInfoRecord.nextHopIfIndex))
	//Even though we could potentially have multiple selected routes, calling update once for this prefix should suffice
	//routeServiceHandler.UpdateIPv4Route(&cfg, nil, nil)
	if selectedRouteInfoRecord.ipType == ribdCommonDefs.IPv4 {
		cfg := ribd.IPv4Route{
			DestinationNw: selectedRouteInfoRecord.destNetIp.String(),
			Protocol:      ReverseRouteProtoTypeMapDB[int(selectedRouteInfoRecord.protocol)],
			Cost:          int32(selectedRouteInfoRecord.metric),
			NetworkMask:   selectedRouteInfoRecord.networkMask.String(),
		}
		nextHop := ribd.NextHopInfo{
			NextHopIp:     selectedRouteInfoRecord.nextHopIp.String(),
			NextHopIntRef: nextHopIf,
		}
		cfg.NextHop = make([]*ribd.NextHopInfo, 0)
		cfg.NextHop = append(cfg.NextHop, &nextHop)
		RouteServiceHandler.Processv4RouteUpdateConfig(&cfg, &cfg, nil)
	} else if selectedRouteInfoRecord.ipType == ribdCommonDefs.IPv6 {
		cfg := ribd.IPv6Route{
			DestinationNw: selectedRouteInfoRecord.destNetIp.String(),
			Protocol:      ReverseRouteProtoTypeMapDB[int(selectedRouteInfoRecord.protocol)],
			Cost:          int32(selectedRouteInfoRecord.metric),
			NetworkMask:   selectedRouteInfoRecord.networkMask.String(),
		}
		nextHop := ribd.NextHopInfo{
			NextHopIp:     selectedRouteInfoRecord.nextHopIp.String(),
			NextHopIntRef: nextHopIf,
		}
		cfg.NextHop = make([]*ribd.NextHopInfo, 0)
		cfg.NextHop = append(cfg.NextHop, &nextHop)
		RouteServiceHandler.Processv6RouteUpdateConfig(&cfg, &cfg, nil)
	}
	return err
}
func policyEngineTraverseAndUpdate() {
	logger.Info("policyEngineTraverseAndUpdate")
	V4RouteInfoMap.VisitAndUpdate(policyEngineUpdateRoute, nil)
	V6RouteInfoMap.VisitAndUpdate(policyEngineUpdateRoute, nil)
}
func policyEngineActionAcceptRoute(params interface{}) {
	routeInfo := params.(RouteParams)
	logger.Info("policyEngineActionAcceptRoute for ip ", routeInfo.destNetIp, " and mask ", routeInfo.networkMask)
	//	_, err := createRoute(routeInfo.ipType, routeInfo.destNetIp, routeInfo.networkMask, routeInfo.metric, routeInfo.weight, routeInfo.nextHopIp, routeInfo.nextHopIfIndex, routeInfo.routeType, routeInfo.createType, ribdCommonDefs.RoutePolicyStateChangetoValid, routeInfo.sliceIdx)
	//_, err := routeServiceHandler.InstallRoute(routeInfo)
	_, err := createRoute(routeInfo)
	if err != nil {
		logger.Info("creating v4 route failed with err ", err)
		return
	}
}
func policyEngineActionUndoSetAdminDistance(actionItem interface{}, conditionsList []interface{}, policyDef policy.Policy, conditionItem interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineActionUndoSetAdminDistance")
	if ProtocolAdminDistanceMapDB == nil {
		logger.Info("ProtocolAdminDistanceMap nil")
		return
	}
	if conditionItem == nil {
		logger.Info("No valid condition provided for set admin distance action")
		return
	}
	conditionInfo := conditionItem.(policy.PolicyCondition).ConditionInfo
	conditionProtocol := conditionInfo.(string)
	//case policyCommonDefs.PolicyConditionTypeProtocolMatch:
	routeDistanceConfig, ok := ProtocolAdminDistanceMapDB[conditionProtocol]
	if !ok {
		logger.Info("Invalid protocol provided for undo set admin distance")
		return
	}
	routeDistanceConfig.configuredDistance = -1
	ProtocolAdminDistanceMapDB[conditionProtocol] = routeDistanceConfig
	logger.Info("Setting configured distance of prototype ", conditionProtocol, " to value ", 0, " default distance of this protocol is ", routeDistanceConfig.defaultDistance)
	policyEngineTraverseAndUpdate()
}

func policyEngineActionSetAdminDistance(actionItem interface{}, conditionList []interface{},
	policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngipolicyEngineActionSetAdminDistance")
	actionInfo := actionItem.(int)
	logger.Info("PoilcyActionTypeSetAdminDistance action to be applied")
	if ProtocolAdminDistanceMapDB == nil {
		logger.Info("ProtocolAdminDistanceMap nil")
		return
	}
	if conditionList == nil {
		logger.Info("No valid condition provided for set admin distance action")
		return
	}
	for i := 0; i < len(conditionList); i++ {
		//case policyCommonDefs.PolicyConditionTypeProtocolMatch:
		conditionProtocol := conditionList[i].(string)
		routeDistanceConfig, ok := ProtocolAdminDistanceMapDB[conditionProtocol]
		if !ok {
			logger.Info("Invalid protocol provided for set admin distance")
			return
		}
		routeDistanceConfig.configuredDistance = actionInfo
		ProtocolAdminDistanceMapDB[conditionProtocol] = routeDistanceConfig
		logger.Info("Setting distance of prototype ", conditionProtocol, " to value ", actionInfo)
	}
	policyEngineTraverseAndUpdate()
	return
}

func policyEngineRouteDispositionAction(action interface{}, conditionInfo []interface{}, policyDef policy.Policy,
	params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineRouteDispositionAction")
	if action.(string) == "Reject" {
		logger.Info("Reject action")
		policyEngineActionRejectRoute(params)
	} else if action.(string) == "Accept" {
		policyEngineActionAcceptRoute(params)
	}
}

func defaultImportPolicyEngineActionFunc(actionInfo interface{}, conditionInfo []interface{},
	policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("defaultImportPolicyEngineAction")
	policyEngineActionAcceptRoute(params)
}

func defaultExportPolicyEngineActionFunc(actionInfo interface{}, conditionInfo []interface{},
	policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("defaultExportPolicyEngineActionFunc")
}

func policyEngineActionNetworkStatementAdvertise(actionInfo interface{}, conditionInfo []interface{},
	policyDef policy.Policy, params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineActionNetworkStatementAdvertise")
	var route ribdInt.Routes
	networkStatementAdvertiseTargetProtocol := actionInfo.(string)
	//Send a event based on target protocol
	RouteInfo := params.(RouteParams)
	var evt int
	if RouteInfo.createType != Invalid {
		logger.Info("Create type not invalid")
		evt = ribdCommonDefs.NOTIFY_ROUTE_CREATED
	} else if RouteInfo.deleteType != Invalid {
		logger.Info("Delete type not invalid")
		evt = ribdCommonDefs.NOTIFY_ROUTE_DELETED
	} else {
		logger.Info("Create/Delete invalid,  so evt = NOTIFY_ROUTE_CREATED")
		evt = ribdCommonDefs.NOTIFY_ROUTE_CREATED
	}
	switch RouteProtocolTypeMapDB[networkStatementAdvertiseTargetProtocol] {
	case ribdCommonDefs.BGP:
		logger.Info("NetworkStatemtnAdvertise to BGP")
		route = ribdInt.Routes{Ipaddr: RouteInfo.destNetIp, Mask: RouteInfo.networkMask, NextHopIp: RouteInfo.nextHopIp, IPAddrType: ribdInt.Int(RouteInfo.ipType), IfIndex: ribdInt.Int(RouteInfo.nextHopIfIndex), Metric: ribdInt.Int(RouteInfo.metric), Prototype: ribdInt.Int(RouteInfo.routeType)}
		route.NetworkStatement = true
		publisherInfo, ok := PublisherInfoMap["BGP"]
		if ok {
			RedistributionNotificationSend(publisherInfo.pub_socket, route, evt, networkStatementAdvertiseTargetProtocol)
		}
		break
	default:
		logger.Info("Unknown target protocol")
	}
	UpdateRedistributeTargetMap(evt, networkStatementAdvertiseTargetProtocol, route)
}

func policyEngineActionRedistribute(actionInfo interface{}, conditionInfo []interface{}, policyDef policy.Policy,
	params interface{}, policyStmt policy.PolicyStmt) {
	logger.Info("policyEngineActionRedistribute")
	var route ribdInt.Routes
	redistributeActionInfo := actionInfo.(policy.RedistributeActionInfo)
	//Send a event based on target protocol
	RouteInfo := params.(RouteParams)
	if (RouteInfo.createType != Invalid || RouteInfo.deleteType != Invalid) && redistributeActionInfo.Redistribute == false {
		logger.Info("Don't redistribute action set for a route create/delete, return")
		return
	}
	var evt int
	if RouteInfo.createType != Invalid {
		logger.Info("Create type not invalid")
		evt = ribdCommonDefs.NOTIFY_ROUTE_CREATED
	} else if RouteInfo.deleteType != Invalid {
		logger.Info("Delete type not invalid")
		evt = ribdCommonDefs.NOTIFY_ROUTE_DELETED
	} else {
		logger.Info("Create/Delete invalid, redistributeAction set to ", redistributeActionInfo.Redistribute)
		if redistributeActionInfo.Redistribute == true {
			logger.Info("evt = NOTIFY_ROUTE_CREATED")
			evt = ribdCommonDefs.NOTIFY_ROUTE_CREATED
		} else {
			logger.Info("evt = NOTIFY_ROUTE_DELETED")
			evt = ribdCommonDefs.NOTIFY_ROUTE_DELETED
		}
	}
	if strings.Contains(ReverseRouteProtoTypeMapDB[int(RouteInfo.routeType)], redistributeActionInfo.RedistributeTargetProtocol) {
		logger.Info("Redistribute target protocol same as route source, do nothing more here")
		return
	}
	if RouteInfo.ipType == ribdCommonDefs.IPv6 {
		testIp := RouteInfo.destNetIp + "/128"
		logger.Info("Redistribute: route dest ip info:", RouteInfo.destNetIp)
		inRange := netUtils.CheckIfInRange(testIp, "fe80::/10", 10, 128)
		if inRange {
			//link local ip , dont redistribute
			return
		}
	}
	route = ribdInt.Routes{Ipaddr: RouteInfo.destNetIp, Mask: RouteInfo.networkMask, NextHopIp: RouteInfo.nextHopIp, IPAddrType: ribdInt.Int(RouteInfo.ipType), IfIndex: ribdInt.Int(RouteInfo.nextHopIfIndex), Metric: ribdInt.Int(RouteInfo.metric), Prototype: ribdInt.Int(RouteInfo.routeType)}
	route.RouteOrigin = ReverseRouteProtoTypeMapDB[int(RouteInfo.routeType)]
	publisherInfo, ok := PublisherInfoMap[redistributeActionInfo.RedistributeTargetProtocol]
	if ok {
		logger.Info("ReditributeNotificationSend event called for target protocol - ", redistributeActionInfo.RedistributeTargetProtocol)
		RedistributionNotificationSend(publisherInfo.pub_socket, route, evt, redistributeActionInfo.RedistributeTargetProtocol)
	} else {
		logger.Info("Unknown target protocol")
	}
	UpdateRedistributeTargetMap(evt, redistributeActionInfo.RedistributeTargetProtocol, route)
}

func UpdateRouteAndPolicyDB(policyDetails policy.PolicyDetails, params interface{}) {
	routeInfo := params.(RouteParams)
	route := ribdInt.Routes{Ipaddr: routeInfo.destNetIp, Mask: routeInfo.networkMask, IPAddrType: ribdInt.Int(routeInfo.ipType), NextHopIp: routeInfo.nextHopIp, IfIndex: ribdInt.Int(routeInfo.nextHopIfIndex), Metric: ribdInt.Int(routeInfo.metric), Prototype: ribdInt.Int(routeInfo.routeType)}
	var op int
	if routeInfo.deleteType != Invalid {
		op = del
	} else {
		if policyDetails.EntityDeleted == false {
			logger.Info("Reject action was not applied, so add this policy to the route")
			op = add
			updateRoutePolicyState(route, op, policyDetails.Policy, policyDetails.PolicyStmt)
		}
		route.PolicyHitCounter++
	}
	updatePolicyRouteMap(route, policyDetails.Policy, op)

}
func DoesRouteExist(params interface{}) (exists bool) {
	//check if the route still exists - it may have been deleted by the previous statement action
	routeDeleted := false
	routeInfo := params.(RouteParams)
	ipPrefix, err := getNetowrkPrefixFromStrings(routeInfo.destNetIp, routeInfo.networkMask)
	if err != nil {
		logger.Info("Error when getting ipPrefix, err= ", err)
		return
	}
	routeInfoRecordList := RouteInfoMapGet(routeInfo.ipType, ipPrefix)
	if routeInfoRecordList == nil {
		logger.Info("Route for type ", routeInfo.ipType, " and prefix", ipPrefix, " no longer exists")
		routeDeleted = true
	} else {
		if routeInfoRecordList.(RouteInfoRecordList).selectedRouteProtocol != ReverseRouteProtoTypeMapDB[int(routeInfo.routeType)] {
			logger.Info("this protocol is not the selected route anymore", err)
			routeDeleted = true
		} else {
			routeInfoList := routeInfoRecordList.(RouteInfoRecordList).routeInfoProtocolMap[routeInfoRecordList.(RouteInfoRecordList).selectedRouteProtocol]
			if routeInfoList == nil {
				logger.Info("Route no longer exists for this protocol")
				routeDeleted = true
			} else {
				routeFound := false
				route := ribdInt.Routes{Ipaddr: routeInfo.destNetIp, Mask: routeInfo.networkMask, NextHopIp: routeInfo.nextHopIp, IfIndex: ribdInt.Int(routeInfo.nextHopIfIndex), Metric: ribdInt.Int(routeInfo.metric), Prototype: ribdInt.Int(routeInfo.routeType)}
				for i := 0; i < len(routeInfoList); i++ {
					testRoute := ribdInt.Routes{Ipaddr: routeInfoList[i].destNetIp.String(), Mask: routeInfoList[i].networkMask.String(), NextHopIp: routeInfoList[i].nextHopIp.String(), IfIndex: ribdInt.Int(routeInfoList[i].nextHopIfIndex), Metric: ribdInt.Int(routeInfoList[i].metric), Prototype: ribdInt.Int(routeInfoList[i].protocol), IsPolicyBasedStateValid: routeInfoList[i].isPolicyBasedStateValid}
					if isSameRoute(testRoute, route) {
						logger.Info("Route still exists")
						routeFound = true
					}
				}
				if !routeFound {
					logger.Info("This specific route no longer exists")
					routeDeleted = true
				}
			}
		}
	}
	exists = !routeDeleted
	return exists
}
func PolicyEngineFilter(route ribdInt.Routes, policyPath int, params interface{}) {
	logger.Info("PolicyEngineFilter")
	var policyPath_Str string
	if policyPath == policyCommonDefs.PolicyPath_Import {
		policyPath_Str = "Import"
	} else if policyPath == policyCommonDefs.PolicyPath_Export {
		policyPath_Str = "Export"
	} else if policyPath == policyCommonDefs.PolicyPath_All {
		policyPath_Str = "ALL"
		logger.Err("policy path ", policyPath_Str, " unexpected in this function")
		return
	}
	routeInfo := params.(RouteParams)
	//if the policy type if ipv6, check if it is link local
	if routeInfo.ipType == ribdCommonDefs.IPv6 {
		testIp := routeInfo.destNetIp + "/128"
		logger.Info("Redistribute: route dest ip info:", routeInfo.destNetIp)
		inRange := netUtils.CheckIfInRange(testIp, "fe80::/10", 10, 128)
		if inRange {
			//link local ip , dont redistribute
			return
		}
	}
	if destNetSlice[routeInfo.sliceIdx].isValid == false && routeInfo.createType != Invalid && policyPath == policyCommonDefs.PolicyPath_Export {
		logger.Info("route down, return from policyenginefilter for deletetype and export path")
		return
	}
	logger.Info("PolicyEngineFilter for policypath ", policyPath_Str, "createType = ", routeInfo.createType, " deleteType = ", routeInfo.deleteType, " route: ", route.Ipaddr, ":", route.Mask, " protocol type: ", route.Prototype, " addrtype:", route.IPAddrType)
	entity, err := buildPolicyEntityFromRoute(route, params)
	if err != nil {
		logger.Info(("Error building policy params"))
		return
	}
	entity.PolicyList = make([]string, 0)
	for j := 0; j < len(route.PolicyList); j++ {
		entity.PolicyList = append(entity.PolicyList, route.PolicyList[j])
	}
	PolicyEngineDB.PolicyEngineFilter(entity, policyPath, params)
	var op int
	if routeInfo.deleteType != Invalid {
		op = delAll //wipe out the policyList
		updateRoutePolicyState(route, op, "", "")
	}
}

func policyEngineApplyForRoute(prefix patriciaDB.Prefix, item patriciaDB.Item, traverseAndApplyPolicyDataInfo patriciaDB.Item) (err error) {
	logger.Info("policyEngineApplyForRoute for route:", item)
	traverseAndApplyPolicyData := traverseAndApplyPolicyDataInfo.(TraverseAndApplyPolicyData)
	rmapInfoRecordList := item.(RouteInfoRecordList)
	if rmapInfoRecordList.routeInfoProtocolMap == nil {
		logger.Info(("rmapInfoRecordList.routeInfoProtocolMap) = nil"))
		return err
	}
	logger.Debug("rmapInfoRecordList:", rmapInfoRecordList, " Selected route protocol = ", rmapInfoRecordList.selectedRouteProtocol)
	selectedRouteList := rmapInfoRecordList.routeInfoProtocolMap[rmapInfoRecordList.selectedRouteProtocol]
	if len(selectedRouteList) == 0 {
		logger.Info("len(selectedRouteList) == 0")
		return err
	}
	for i := 0; i < len(selectedRouteList); i++ {
		selectedRouteInfoRecord := selectedRouteList[i]
		if selectedRouteInfoRecord.sliceIdx == -1 || selectedRouteInfoRecord.sliceIdx >= len(destNetSlice) || destNetSlice[selectedRouteInfoRecord.sliceIdx].isValid == false {
			logger.Info("route ", selectedRouteInfoRecord, " not valid, continue, sliceIdx:", selectedRouteInfoRecord.sliceIdx, " len(destNetSlice):", len(destNetSlice))
			continue
		}
		policyRoute := ribdInt.Routes{Ipaddr: selectedRouteInfoRecord.destNetIp.String(), Mask: selectedRouteInfoRecord.networkMask.String(), NextHopIp: selectedRouteInfoRecord.nextHopIp.String(), IfIndex: ribdInt.Int(selectedRouteInfoRecord.nextHopIfIndex), Metric: ribdInt.Int(selectedRouteInfoRecord.metric), Prototype: ribdInt.Int(selectedRouteInfoRecord.protocol), IsPolicyBasedStateValid: rmapInfoRecordList.isPolicyBasedStateValid}
		params := RouteParams{destNetIp: policyRoute.Ipaddr, networkMask: policyRoute.Mask, routeType: ribd.Int(policyRoute.Prototype), nextHopIp: selectedRouteInfoRecord.nextHopIp.String(), sliceIdx: ribd.Int(policyRoute.SliceIdx), createType: Invalid, deleteType: Invalid}
		entity, err := buildPolicyEntityFromRoute(policyRoute, params)
		if err != nil {
			logger.Err("Error builiding policy entity params")
			return err
		}
		entity.PolicyList = make([]string, 0)
		for j := 0; j < len(rmapInfoRecordList.policyList); j++ {
			entity.PolicyList = append(entity.PolicyList, rmapInfoRecordList.policyList[j])
		}
		traverseAndApplyPolicyData.updatefunc(entity, traverseAndApplyPolicyData.data, params)
	}
	return err
}
func policyEngineTraverseAndApply(data interface{}, updatefunc policy.PolicyApplyfunc) {
	logger.Info("PolicyEngineTraverseAndApply - traverse routing table and apply policy ")
	traverseAndApplyPolicyData := TraverseAndApplyPolicyData{data: data, updatefunc: updatefunc}
	V4RouteInfoMap.VisitAndUpdate(policyEngineApplyForRoute, traverseAndApplyPolicyData)
	V6RouteInfoMap.VisitAndUpdate(policyEngineApplyForRoute, traverseAndApplyPolicyData)
}
func policyEngineTraverseAndReverse(applyPolicyItem interface{}) {
	updateInfo := applyPolicyItem.(policy.PolicyEngineApplyInfo)
	applyPolicyInfo := updateInfo.ApplyPolicy //.(policy.ApplyPolicyInfo)
	policy := applyPolicyInfo.ApplyPolicy     //policyItem.(policy.Policy)
	logger.Info("PolicyEngineTraverseAndReverse - traverse routing table and inverse policy actions", policy.Name)
	ext := policy.Extensions.(PolicyExtensions)
	if ext.routeList == nil {
		logger.Info("No route affected by this policy, so nothing to do")
		return
	}
	var policyRoute ribdInt.Routes
	var params RouteParams
	for idx := 0; idx < len(ext.routeInfoList); idx++ {
		policyRoute = ext.routeInfoList[idx]
		params = RouteParams{destNetIp: policyRoute.Ipaddr, networkMask: policyRoute.Mask, routeType: ribd.Int(policyRoute.Prototype), sliceIdx: ribd.Int(policyRoute.SliceIdx), createType: Invalid, deleteType: Invalid}
		ipPrefix, err := getNetowrkPrefixFromStrings(ext.routeInfoList[idx].Ipaddr, ext.routeInfoList[idx].Mask)
		if err != nil {
			logger.Info("Invalid route ", ext.routeList[idx])
			continue
		}
		entity, err := buildPolicyEntityFromRoute(policyRoute, params)
		if err != nil {
			logger.Err("Error builiding policy entity params")
			return
		}
		//PolicyEngineDB.PolicyEngineUndoPolicyForEntity(entity, policy, params)
		success := PolicyEngineDB.PolicyEngineUndoApplyPolicyForEntity(entity, updateInfo, params)
		if success {
			deleteRoutePolicyState(params.ipType, ipPrefix, policy.Name)
			PolicyEngineDB.DeletePolicyEntityMapEntry(entity, policy.Name)
		}
	}
}
