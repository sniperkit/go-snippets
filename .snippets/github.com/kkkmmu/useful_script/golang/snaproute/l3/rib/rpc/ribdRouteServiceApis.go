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

// ribdRouteServiceApis.go
package rpc

import (
	"errors"
	defs "l3/rib/ribdCommonDefs"
	"l3/rib/server"
	"models/objects"
	"ribd"
	"ribdInt"
)

/* Create route API
 */

func (m RIBDServicesHandler) CreateIPv4Route(cfg *ribd.IPv4Route) (val bool, err error) {
	logger.Info("Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask)
	/* Validate Route config parameters for "add" operation
	 */
	err = m.server.RouteConfigValidationCheck(cfg, "add")
	if err != nil {
		logger.Err("validation check failed with error ", err)
		return false, err
	}
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.Add,
	}
	return true, nil
}

/*
   OnewayCreate API for route
*/
func (m RIBDServicesHandler) OnewayCreateIPv4Route(cfg *ribd.IPv4Route) (err error) {
	logger.Info("OnewayCreateIPv4Route - Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask, "cfg.NextHopIntRef: ", cfg.NextHop[0].NextHopIntRef)
	m.CreateIPv4Route(cfg)
	/*m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.Add,
	}*/
	return err
}

func (m RIBDServicesHandler) CreateIPv6Route(cfg *ribd.IPv6Route) (val bool, err error) {
	logger.Info("Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask)
	/* Validate Route config parameters for "add" operation
	 */
	err = m.server.IPv6RouteConfigValidationCheck(cfg, "add")
	if err != nil {
		logger.Err("validation check failed with error ", err)
		return false, err
	}
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.Addv6,
	}
	return true, nil
}

/*
   OnewayCreate API for route
*/
func (m RIBDServicesHandler) OnewayCreateIPv6Route(cfg *ribd.IPv6Route) (err error) {
	logger.Info("OnewayCreateIPv6Route - Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask, "cfg.NextHopIntRef: ", cfg.NextHop[0].NextHopIntRef)
	m.CreateIPv6Route(cfg)
	return err
}

/*
   Create Routes in Bulk using Oneway create API
*/
func (m RIBDServicesHandler) OnewayCreateBulkIPv4Route(cfg []*ribdInt.IPv4RouteConfig) (err error) {
	logger.Info("OnewayCreateBulkIPv4Route for ", len(cfg), " routes")
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigBulkRouteConfigObject: cfg,
		Op: defs.AddBulk,
	}
	/*for i := 0; i < len(cfg); i++ {
		newCfg := ribd.IPv4Route{
			DestinationNw: cfg[i].DestinationNw,
			NetworkMask:   cfg[i].NetworkMask,
			Cost:          cfg[i].Cost,
			Protocol:      cfg[i].Protocol,
		}
		newCfg.NextHop = make([]*ribd.NextHopInfo, 0)
		nextHop := ribd.NextHopInfo{
			NextHopIp:     cfg[i].NextHopIp,
			NextHopIntRef: cfg[i].NextHopIntRef,
			Weight:        cfg[i].Weight,
		}
		newCfg.NextHop = append(newCfg.NextHop, &nextHop)
		m.CreateIPv4Route(&newCfg)
	}*/
	return err
}

/*
   Delete Route
*/
func (m RIBDServicesHandler) DeleteIPv4Route(cfg *ribd.IPv4Route) (val bool, err error) {
	logger.Info("DeleteIPv4Route:Received Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "Protocol ", cfg.Protocol, "number of nextHops: ", len(cfg.NextHop))
	/*
	   Validate route config parameters for "del" operation
	*/
	err = m.server.RouteConfigValidationCheck(cfg, "del")
	if err != nil {
		logger.Err("validation check failed with error ", err)
		return false, err
	}
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.Del,
	}
	return true, nil
}

/*
   Delete route using Oneway Api
*/
func (m RIBDServicesHandler) OnewayDeleteIPv4Route(cfg *ribd.IPv4Route) (err error) {
	logger.Info("OnewayDeleteIPv4Route:RouteReceived Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "nextHopIP:", cfg.NextHop[0].NextHopIp, "Protocol ", cfg.Protocol)
	m.DeleteIPv4Route(cfg)
	return err
}

/*
   Delete Route
*/
func (m RIBDServicesHandler) DeleteIPv6Route(cfg *ribd.IPv6Route) (val bool, err error) {
	logger.Info("DeleteIPv6Route:Received Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "Protocol ", cfg.Protocol, "number of nextHops: ", len(cfg.NextHop))
	/*
	   Validate route config parameters for "del" operation
	*/
	err = m.server.IPv6RouteConfigValidationCheck(cfg, "del")
	if err != nil {
		logger.Err("validation check failed with error ", err)
		return false, err
	}
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.Delv6,
	}
	return true, nil
}

/*
   Delete route using Oneway Api
*/
func (m RIBDServicesHandler) OnewayDeleteIPv6Route(cfg *ribd.IPv6Route) (err error) {
	logger.Info("OnewayDeleteIPv6Route:RouteReceived Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "nextHopIP:", cfg.NextHop[0].NextHopIp, "Protocol ", cfg.Protocol)
	m.DeleteIPv6Route(cfg)
	return err
}

/*
   Update route
*/
func (m RIBDServicesHandler) UpdateIPv4Route(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) { //[]*ribd.PatchOpInfo) (val bool, err error) {
	logger.Info("UpdateIPv4Route: Received update route request")
	/*
	   validate route config parameters for update operation
	*/
	if op == nil || len(op) == 0 {
		err = m.server.RouteConfigValidationCheckForUpdate(origconfig, newconfig, attrset)
		if err != nil {
			logger.Err("validation check failed with error ", err)
			return false, err
		}
	} else {
		err = m.server.RouteConfigValidationCheckForPatchUpdate(origconfig, newconfig, op)
		if err != nil {
			logger.Err("validation check failed with error ", err)
			return false, err
		}
	}
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		Op:               defs.Update,
		PatchOp:          op,
	}

	return true, nil
}

/*
   one way update route function
*/
func (m RIBDServicesHandler) OnewayUpdateIPv4Route(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, attrset []bool) (err error) {
	logger.Info("OneWayUpdateIPv4Route: Received update route request")
	m.UpdateIPv4Route(origconfig, newconfig, attrset, nil)
	return err
}

/*
   Update route
*/
func (m RIBDServicesHandler) UpdateIPv6Route(origconfig *ribd.IPv6Route, newconfig *ribd.IPv6Route, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) { //[]*ribd.PatchOpInfo) (val bool, err error) {
	logger.Info("UpdateIPv6Route: Received update route request")
	/*
	   validate route config parameters for update operation
	*/
	if op == nil || len(op) == 0 {
		logger.Debug("UpdateIPv6Route:At the beginning origconfig.destinationnw:", origconfig.DestinationNw, " newconfig.DesinationNw:", newconfig.DestinationNw)
		err = m.server.IPv6RouteConfigValidationCheckForUpdate(origconfig, newconfig, attrset)
		if err != nil {
			logger.Err("validation check failed with error ", err)
			return false, err
		}
		logger.Debug("UpdateIPv6Route:At the end origconfig.destinationnw:", origconfig.DestinationNw, " newconfig.DesinationNw:", newconfig.DestinationNw)
	} else {
		err = m.server.IPv6RouteConfigValidationCheckForPatchUpdate(origconfig, newconfig, op)
		if err != nil {
			logger.Err("validation check failed with error ", err)
			return false, err
		}
	}
	logger.Debug("UpdateIPv6Route:Call routeconfch origconfig.destinationnw:", origconfig.DestinationNw, " newconfig.DesinationNw:", newconfig.DestinationNw)
	m.server.RouteConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		Op:               defs.Updatev6,
		PatchOp:          op,
	}

	return true, nil
}

/*
   one way update route function
*/
func (m RIBDServicesHandler) OnewayUpdateIPv6Route(origconfig *ribd.IPv6Route, newconfig *ribd.IPv6Route, attrset []bool) (err error) {
	logger.Info("OneWayUpdateIPv6Route: Received update route request")
	m.UpdateIPv6Route(origconfig, newconfig, attrset, nil)
	return err
}

/*
   Applications call this function to fetch all the routes that need to be redistributed into them.
*/
func (m RIBDServicesHandler) GetBulkRoutesForProtocol(srcProtocol string, fromIndex ribdInt.Int, rcount ribdInt.Int) (routes *ribdInt.RoutesGetInfo, err error) {
	ret, err := m.server.GetBulkRoutesForProtocol(srcProtocol, fromIndex, rcount)
	return ret, err
}

/*
   Api to track a route's reachability status
*/
func (m RIBDServicesHandler) TrackReachabilityStatus(ipAddr string, protocol string, op string) (err error) {
	m.server.TrackReachabilityCh <- server.TrackReachabilityInfo{ipAddr, protocol, op}
	return nil
}

func (m RIBDServicesHandler) GetIPv4RouteState(destNw string) (*ribd.IPv4RouteState, error) {
	logger.Info("Get state for IPv4Route")
	route := ribd.NewIPv4RouteState()
	if m.server.DbHdl == nil {
		logger.Err("DbHdl not initialized")
		return route, errors.New("DBHdl not initialized")
	}
	var routeObj objects.IPv4RouteState
	var routeObjtemp objects.IPv4RouteState
	obj, err := m.server.DbHdl.GetObjectFromDb(routeObj, destNw)
	if err == nil {
		routeObjtemp = obj.(objects.IPv4RouteState)
		objects.ConvertribdIPv4RouteStateObjToThrift(&routeObjtemp, route)
	}
	return route, nil
}

func (m RIBDServicesHandler) GetBulkIPv4RouteState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.IPv4RouteStateGetInfo, err error) {
	logger.Debug("GetBulkIPv4RouteState")
	returnRoutes := make([]*ribd.IPv4RouteState, 0)
	var returnRouteGetInfo ribd.IPv4RouteStateGetInfo
	routes = &returnRouteGetInfo
	if m.server.DbHdl == nil {
		logger.Err("DbHdl not initialized")
		return routes, errors.New("DBHdl not initialized")
	}
	var routeObj objects.IPv4RouteState
	var routeObjtemp objects.IPv4RouteState
	err, objCount, nextMarker, more, objs := m.server.DbHdl.GetBulkObjFromDb(routeObj, int64(fromIndex), int64(rcount))
	logger.Debug("objCount = ", objCount, " len(obj) ", len(objs), " more ", more, " nextMarker: ", nextMarker)
	var tempRoute []ribd.IPv4RouteState = make([]ribd.IPv4RouteState, len(objs))
	if err == nil {
		for i := 0; i < len(objs); i++ {
			routeObjtemp = objs[i].(objects.IPv4RouteState)
			logger.Debug("obj ", i, routeObjtemp.DestinationNw, " ", routeObjtemp.NextHopList)
			objects.ConvertribdIPv4RouteStateObjToThrift(&routeObjtemp, &tempRoute[i])
			returnRoutes = append(returnRoutes, &tempRoute[i])
		}
		routes.IPv4RouteStateList = returnRoutes
		routes.StartIdx = fromIndex
		routes.EndIdx = ribd.Int(nextMarker)
		routes.More = more
		routes.Count = ribd.Int(objCount)
		/*		if routes.Count > 0 {
					fmt.Println(" DestinationNw  NextHop")
				}
				for _,rt := range routes.IPv4RouteStateList {
					fmt.Println(rt.DestinationNw , " ", rt.NextHopList)
				}*/
		return routes, err
	}
	return routes, err
}

func (m RIBDServicesHandler) GetIPv6RouteState(destNw string) (*ribd.IPv6RouteState, error) {
	logger.Info("Get state for IPv6Route")
	route := ribd.NewIPv6RouteState()
	if m.server.DbHdl == nil {
		logger.Err("DbHdl not initialized")
		return route, errors.New("DBHdl not initialized")
	}
	var routeObj objects.IPv6RouteState
	var routeObjtemp objects.IPv6RouteState
	obj, err := m.server.DbHdl.GetObjectFromDb(routeObj, destNw)
	if err == nil {
		routeObjtemp = obj.(objects.IPv6RouteState)
		objects.ConvertribdIPv6RouteStateObjToThrift(&routeObjtemp, route)
	}
	return route, nil
}

func (m RIBDServicesHandler) GetBulkIPv6RouteState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.IPv6RouteStateGetInfo, err error) {
	logger.Debug("GetBulkIPv6RouteState")
	returnRoutes := make([]*ribd.IPv6RouteState, 0)
	var returnRouteGetInfo ribd.IPv6RouteStateGetInfo
	routes = &returnRouteGetInfo
	if m.server.DbHdl == nil {
		logger.Err("DbHdl not initialized")
		return routes, errors.New("DBHdl not initialized")
	}
	var routeObj objects.IPv6RouteState
	var routeObjtemp objects.IPv6RouteState
	err, objCount, nextMarker, more, objs := m.server.DbHdl.GetBulkObjFromDb(routeObj, int64(fromIndex), int64(rcount))
	logger.Debug("objCount = ", objCount, " len(obj) ", len(objs), " more ", more, " nextMarker: ", nextMarker)
	var tempRoute []ribd.IPv6RouteState = make([]ribd.IPv6RouteState, len(objs))
	if err == nil {
		for i := 0; i < len(objs); i++ {
			routeObjtemp = objs[i].(objects.IPv6RouteState)
			logger.Debug("obj ", i, routeObjtemp.DestinationNw, " ", routeObjtemp.NextHopList)
			objects.ConvertribdIPv6RouteStateObjToThrift(&routeObjtemp, &tempRoute[i])
			returnRoutes = append(returnRoutes, &tempRoute[i])
		}
		routes.IPv6RouteStateList = returnRoutes
		routes.StartIdx = fromIndex
		routes.EndIdx = ribd.Int(nextMarker)
		routes.More = more
		routes.Count = ribd.Int(objCount)
		/*		if routes.Count > 0 {
					fmt.Println(" DestinationNw  NextHop")
				}
				for _,rt := range routes.IPv4RouteStateList {
					fmt.Println(rt.DestinationNw , " ", rt.NextHopList)
				}*/
		return routes, err
	}
	return routes, err
}
func (m RIBDServicesHandler) GetBulkRouteStatsPerProtocolState(fromIndex ribd.Int, count ribd.Int) (stats *ribd.RouteStatsPerProtocolStateGetInfo, err error) {
	ret, err := m.server.GetBulkRouteStatsPerProtocolState(fromIndex, count)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteStatsPerProtocolState(Protocol string) (stats *ribd.RouteStatsPerProtocolState, err error) {
	stats = ribd.NewRouteStatsPerProtocolState()
	stats, err = m.server.GetRouteStatsPerProtocolState(Protocol)
	return stats, err
}
func (m RIBDServicesHandler) GetBulkRouteStatsPerInterfaceState(fromIndex ribd.Int, count ribd.Int) (stats *ribd.RouteStatsPerInterfaceStateGetInfo, err error) {
	ret, err := m.server.GetBulkRouteStatsPerInterfaceState(fromIndex, count)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteStatsPerInterfaceState(Intfref string) (stats *ribd.RouteStatsPerInterfaceState, err error) {
	stats = ribd.NewRouteStatsPerInterfaceState()
	stats, err = m.server.GetRouteStatsPerInterfaceState(Intfref)
	return stats, err
}

func (m RIBDServicesHandler) GetBulkRouteStatState(fromIndex ribd.Int, count ribd.Int) (stats *ribd.RouteStatStateGetInfo, err error) {
	if fromIndex != 0 {
		err := errors.New("Invalid range")
		return nil, err
	}
	tempstats := make([]*ribd.RouteStatState, 1)
	var ret_stats ribd.RouteStatStateGetInfo
	stats = &ret_stats
	tempstats[0] = &ribd.RouteStatState{}
	tempstats[0].PerProtocolRouteCountList = m.server.GetPerProtocolRouteCountList()
	for _, v := range tempstats[0].PerProtocolRouteCountList {
		tempstats[0].TotalRouteCount = tempstats[0].TotalRouteCount + v.RouteCount
		tempstats[0].ECMPRouteCount = tempstats[0].ECMPRouteCount + v.EcmpCount
	}
	v4Count, _ := m.GetTotalv4RouteCount()
	tempstats[0].V4RouteCount = int32(v4Count)
	v6Count, _ := m.GetTotalv6RouteCount()
	tempstats[0].V6RouteCount = int32(v6Count)
	stats.RouteStatStateList = tempstats
	stats.StartIdx = fromIndex
	stats.EndIdx = fromIndex
	stats.Count = 1
	stats.More = false
	return stats, err
}
func (m RIBDServicesHandler) GetRouteStatState(vrf string) (*ribd.RouteStatState, error) {
	stat := ribd.NewRouteStatState()
	v4Count, _ := m.GetTotalv4RouteCount()
	v6Count, _ := m.GetTotalv6RouteCount()
	stat.PerProtocolRouteCountList = m.server.GetPerProtocolRouteCountList()
	for _, v := range stat.PerProtocolRouteCountList {
		stat.TotalRouteCount = stat.TotalRouteCount + v.RouteCount
		stat.ECMPRouteCount = stat.ECMPRouteCount + v.EcmpCount
	}
	stat.V4RouteCount = int32(v4Count)
	stat.V6RouteCount = int32(v6Count)
	return stat, nil
}

func (m RIBDServicesHandler) GetRIBEventState(index int32) (*ribd.RIBEventState, error) {
	logger.Info("Get state for RIBEventState")
	route := ribd.NewRIBEventState()
	return route, nil
}

func (m RIBDServicesHandler) GetBulkRIBEventState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.RIBEventStateGetInfo, err error) {
	ret, err := m.server.GetBulkRIBEventState(fromIndex, rcount)
	return ret, err
}

func (m RIBDServicesHandler) GetBulkRouteDistanceState(fromIndex ribd.Int, rcount ribd.Int) (routeDistanceStates *ribd.RouteDistanceStateGetInfo, err error) {
	ret, err := m.server.GetBulkRouteDistanceState(fromIndex, rcount)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteDistanceState(protocol string) (*ribd.RouteDistanceState, error) {
	logger.Info("Get state for RouteDistanceState")
	state := ribd.NewRouteDistanceState()
	state, err := m.server.GetRouteDistanceState(protocol)
	return state, err
}
func (m RIBDServicesHandler) Getv4Route(destNetIp string) (route *ribdInt.IPv4RouteState, err error) {
	ret, err := m.server.Getv4Route(destNetIp)
	return ret, err
}
func (m RIBDServicesHandler) Getv6Route(destNetIp string) (route *ribdInt.IPv6RouteState, err error) {
	ret, err := m.server.Getv6Route(destNetIp)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteReachabilityInfo(destNet string, ifIndex ribdInt.Int) (nextHopIntf *ribdInt.NextHopInfo, err error) {
	nh, err := m.server.GetRouteReachabilityInfo(destNet, ifIndex)
	return nh, err
}
func (m RIBDServicesHandler) GetTotalv4RouteCount() (number ribdInt.Int, err error) {
	num, err := m.server.GetTotalv4RouteCount()
	return ribdInt.Int(num), err
}
func (m RIBDServicesHandler) GetTotalv6RouteCount() (number ribdInt.Int, err error) {
	num, err := m.server.GetTotalv6RouteCount()
	return ribdInt.Int(num), err
}
func (m RIBDServicesHandler) Getv4RouteCreatedTime(number ribdInt.Int) (time string, err error) {
	time, err = m.server.Getv4RouteCreatedTime(int(number))
	return time, err
}
func (m RIBDServicesHandler) Getv6RouteCreatedTime(number ribdInt.Int) (time string, err error) {
	time, err = m.server.Getv6RouteCreatedTime(int(number))
	return time, err
}
