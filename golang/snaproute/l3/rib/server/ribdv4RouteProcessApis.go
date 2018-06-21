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
	"errors"
	"fmt"
	defs "l3/rib/ribdCommonDefs"
	//"models/objects"
	"net"
	"reflect"
	"ribd"
	"ribdInt"
	"strconv"
	"strings"
	netutils "utils/netUtils"
	"utils/patriciaDB"
	//"utils/policy/policyCommonDefs"
)

var V4RouteInfoMap *patriciaDB.Trie //Routes are stored in patricia trie

/*
   Returns the longest prefix match route to reach the destination network destNet
*/
func (m RIBDServer) GetV4RouteReachabilityInfo(destNet string, ifIndex ribdInt.Int) (nextHopIntf *ribdInt.NextHopInfo, err error) {
	logger.Debug("GetV4RouteReachabilityInfo of ", destNet, " ifIndex:", ifIndex)
	//t1 := time.Now()
	var retnextHopIntf ribdInt.NextHopInfo
	nextHopIntf = &retnextHopIntf
	var found bool
	destNetIp, err := getIP(destNet)
	if err != nil {
		logger.Err("getIP returned Invalid dest ip address for ", destNet)
		return nextHopIntf, errors.New("Invalid dest ip address")
	}
	lookupIp := destNetIp.To4()
	if lookupIp == nil {
		//logger.Err("Incorrect ip type lookup")
		return nextHopIntf, errors.New("Incorrect ip type lookup")
	}
	destNetIp = lookupIp
	rmapInfoListItem := V4RouteInfoMap.GetLongestPrefixNode(patriciaDB.Prefix(destNetIp))
	if rmapInfoListItem != nil {
		//fmt.Println("Madhavi!! GetV4RouteReachabilityInfo:, rmapInfoList not nil for ", destNetIp)
		rmapInfoList := rmapInfoListItem.(RouteInfoRecordList)
		if rmapInfoList.selectedRouteProtocol != "INVALID" {
			//fmt.Println("Madhavi!! GetV4RouteReachabilityInfo: setting found = true")
			found = true
			routeInfoList, ok := rmapInfoList.routeInfoProtocolMap[rmapInfoList.selectedRouteProtocol]
			if !ok || len(routeInfoList) == 0 {
				//fmt.Println("Madhavi!! GetV4RouteReachabilityInfo: ok:", ok, " len(routeInfoList):", len(routeInfoList))
				logger.Err("Selected route not found because len(routeInfoList) = 0")
				return nil, errors.New("dest ip address not reachable")
			}
			//v := routeInfoList[0]
			nhFound, v, _ := findRouteWithNextHop(routeInfoList, defs.IPv4, "", ribd.Int(ifIndex))
			//fmt.Println("Madhavi!! GetV4RouteReachabilityInfo: nhFound:", nhFound, " v:", v)
			if !nhFound {
				logger.Err("Next hop for ifIndex:", ifIndex, " not found")
				return nil, errors.New(fmt.Sprintln("dest ip address not reachable via ifIndex", ifIndex))
			}
			//logger.Info("v.networkAddr", v.networkAddr)
			prefix := strings.SplitAfter(v.networkAddr, "/")
			if len(prefix) == 0 {
				logger.Err(v.networkAddr, " not in cidr format")
				return nextHopIntf, errors.New(fmt.Sprintln(v.networkAddr, " not in cidr format"))
			}
			cidrLen, _ := strconv.Atoi(prefix[(len(prefix) - 1)])
			//logger.Info("cidr len := ", cidrLen, " prefix", prefix)
			reachable := netutils.CheckIfInRange(v.networkAddr, destNet+"/"+prefix[(len(prefix)-1)], cidrLen, 32)
			if !reachable {
				logger.Info("CheckIfInRange reachability check returned false for ", v.networkAddr, ",", destNet+"/"+prefix[(len(prefix)-1)], ",", cidrLen, ",", 32)
				found = false
			}
			nextHopIntf.NextHopIp = v.nextHopIp.String()
			nextHopIntf.NextHopIfIndex = ribdInt.Int(v.nextHopIfIndex)
			nextHopIntf.Metric = ribdInt.Int(v.metric)
			nextHopIntf.Ipaddr = v.destNetIp.String()
			nextHopIntf.Mask = v.networkMask.String()
			nextHopIntf.IsReachable = true
			//nextHopIntf.IsReachable = v.resolvedNextHopIpIntf.IsReachable
		}
	}

	if found == false {
		logger.Err("dest IP", destNetIp, " not reachable ")
		err = errors.New("dest ip address not reachable")
		return nextHopIntf, err
	}
	//duration := time.Since(t1)
	//logger.Debug("time to get longestPrefixLen = ", duration.Nanoseconds(), " ipAddr of the route: ", nextHopIntf.Ipaddr, " next hop ip of the route = ", nextHopIntf.NextHopIp, " ifIndex: ", nextHopIntf.NextHopIfIndex)
	return nextHopIntf, err
}

/*
    Function updates the route reachability status of a network. When a route is created/deleted/state changes,
	we traverse the entire route map and call this function for each of the destination network with :
	    prefix = route prefix of route being visited
		handle = routeInfoList data stored at this node
		item - reachabilityInfo data formed with route that is modified and the state
*/
func UpdateV4RouteReachabilityStatus(prefix patriciaDB.Prefix, //prefix of the node being traversed
	handle patriciaDB.Item, //data interface (routeInforRecordList) for this node
	item patriciaDB.Item) /*RouteReachabilityStatusInfo data */ (err error) {

	if handle == nil {
		logger.Err("nil handle")
		return err
	}
	routeReachabilityStatusInfo := item.(RouteReachabilityStatusInfo)
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(routeReachabilityStatusInfo.destNet)
	if err != nil {
		logger.Err("Error getting IP from cidr: ", routeReachabilityStatusInfo.destNet)
		return err
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	destIpPrefix, err := getNetowrkPrefixFromStrings(ipAddrStr, ipMaskStr)
	if err != nil {
		logger.Err("Error getting ip prefix for ip:", ipAddrStr, " mask:", ipMaskStr)
		return err
	}
	logger.Debug("UpdateRouteReachabilityStatus network: ", routeReachabilityStatusInfo.destNet, " status:", routeReachabilityStatusInfo.status, "ip: ", ip.String(), " destIPPrefix: ", destIpPrefix, " ipMaskStr:", ipMaskStr)
	rmapInfoRecordList := handle.(RouteInfoRecordList)
	//for each of the routes for this destination, check if the nexthop ip matches destPrefix - which is the route being modified
	for k, v := range rmapInfoRecordList.routeInfoProtocolMap {
		//logger.Debug("UpdateRouteReachabilityStatus - protocol: ", k)
		for i := 0; i < len(v); i++ {
			if v[i].nextHopIpType != routeReachabilityStatusInfo.ipType {
				logger.Debug("Skipping nexthop:", v[i].nextHopIp.String(), " since the nextHopIpType ", v[i].nextHopIpType, " not the same as ipType:", routeReachabilityStatusInfo.ipType)
				continue
			}
			vPrefix, err := getNetowrkPrefixFromStrings(v[i].nextHopIp.String(), ipMaskStr)
			if err != nil {
				logger.Err("Error getting ip prefix for v[i].nextHopIp:", v[i].nextHopIp.String(), " mask:", ipMaskStr)
				return err
			}
			nextHopIntf := ribdInt.NextHopInfo{
				NextHopIp:      v[i].nextHopIp.String(),
				NextHopIfIndex: ribdInt.Int(v[i].nextHopIfIndex),
			}
			//is the next hop same as the modified route
			if bytes.Equal(vPrefix, destIpPrefix) {
				if routeReachabilityStatusInfo.status == "Down" && v[i].resolvedNextHopIpIntf.IsReachable == true {
					v[i].resolvedNextHopIpIntf.IsReachable = false
					rmapInfoRecordList.routeInfoProtocolMap[k] = v
					V4RouteInfoMap.Set(prefix, rmapInfoRecordList)
					//logger.Debug("Adding to DBRouteCh from updateRouteReachability case 1")
					RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
						OrigConfigObject: RouteDBInfo{v[i], rmapInfoRecordList},
						Op:               defs.Add,
					}
					//RouteServiceHandler.WriteIPv4RouteStateEntryToDB(RouteDBInfo{v[i], rmapInfoRecordList})
					//logger.Debug("Bringing down route : ip: ", v[i].networkAddr)
					RouteReachabilityStatusUpdate(k, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Down", k, nextHopIntf})
					/*
					   The reachability status for this network has been updated, now check if there are routes dependent on
					   this prefix and call reachability status
					*/
					if RouteServiceHandler.NextHopInfoMap[NextHopInfoKey{string(prefix)}].refCount > 0 {
						//logger.Debug("There are dependent routes for this ip ", v[i].networkAddr)
						V4RouteInfoMap.VisitAndUpdate(UpdateV4RouteReachabilityStatus, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Down", k, nextHopIntf})
					}
				} else if routeReachabilityStatusInfo.status == "Up" && v[i].resolvedNextHopIpIntf.IsReachable == false {
					//logger.Debug("Bringing up route : ip: ", v[i].networkAddr)
					v[i].resolvedNextHopIpIntf.IsReachable = true
					rmapInfoRecordList.routeInfoProtocolMap[k] = v
					V4RouteInfoMap.Set(prefix, rmapInfoRecordList)
					//logger.Debug("Adding to DBRouteCh from updateRouteReachability case 2")
					RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
						OrigConfigObject: RouteDBInfo{v[i], rmapInfoRecordList},
						Op:               defs.Add,
					}
					//RouteServiceHandler.WriteIPv4RouteStateEntryToDB(RouteDBInfo{v[i], rmapInfoRecordList})
					RouteReachabilityStatusUpdate(k, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Up", k, nextHopIntf})
					/*
					   The reachability status for this network has been updated, now check if there are routes dependent on
					   this prefix and call reachability status
					*/
					if RouteServiceHandler.NextHopInfoMap[NextHopInfoKey{string(prefix)}].refCount > 0 {
						//logger.Debug("There are dependent routes for this ip ", v[i].networkAddr)
						V4RouteInfoMap.VisitAndUpdate(UpdateV4RouteReachabilityStatus, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Up", k, nextHopIntf})
					}
				}
			}
		}
	}
	return err
}

/*
    This function performs config parameters validation for Route update operation.
	Key validations performed by this fucntion include:
	   - Validate destinationNw. If provided in CIDR notation, convert to ip addr and mask values
*/
func (m RIBDServer) RouteConfigValidationCheckForUpdate(oldcfg *ribd.IPv4Route, cfg *ribd.IPv4Route, attrset []bool) (err error) {
	//logger.Info("RouteConfigValidationCheckForUpdate")
	if netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot update ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API")
		return errors.New(fmt.Sprintln("Cannot update ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API"))
	}
	isCidr := strings.Contains(cfg.DestinationNw, "/")
	if isCidr {
		/*
		   the given address is in CIDR format
		*/
		ip, ipNet, err := net.ParseCIDR(cfg.DestinationNw)
		if err != nil {
			logger.Err("Invalid Destination IP address")
			return errors.New("Invalid Desitnation IP address")
		}
		_, err = getNetworkPrefixFromCIDR(cfg.DestinationNw)
		if err != nil {
			return errors.New("Invalid destination ip/network Mask")
		}
		cfg.DestinationNw = ip.String()
		oldcfg.DestinationNw = ip.String()
		ipMask := make(net.IP, 4)
		copy(ipMask, ipNet.Mask)
		ipMaskStr := net.IP(ipMask).String()
		cfg.NetworkMask = ipMaskStr
		oldcfg.NetworkMask = ipMaskStr
	}
	_, err = validateNetworkPrefix(cfg.DestinationNw, cfg.NetworkMask)
	if err != nil {
		logger.Err(" getNetowrkPrefixFromStrings returned err ", err)
		return errors.New("Invalid destination ip address")
	}
	/*
		    Default operation for update function is to update route Info. The following
			logic deals with updating route attributes
	*/
	if attrset != nil {
		//logger.Debug("attr set not nil, set individual attributes")
		objTyp := reflect.TypeOf(*cfg)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				//logger.Debug("ProcessRouteUpdateConfig (server): changed ", objName)
				if objName == "Protocol" {
					/*
					   Updating route protocol type is not allowed
					*/
					logger.Err("Cannot update Protocol value of a route")
					return errors.New("Cannot set Protocol field")
				}
				if objName == "NullRoute" {
					logger.Err("Cannot update null route attribute, please delete and create the route with the correct value")
					return errors.New("Cannot update null route attribute, please delete and create the route with the correct value")
				}
				if objName == "NextHop" {
					/*
					   Next hop info is being updated
					*/
					if len(cfg.NextHop) == 0 {
						/*
						   Expects non-zero nexthop info
						*/
						logger.Err("Must specify next hop")
						return errors.New("Next hop ip not specified")
					}
					/*
					   Check if next hop IP is valid
					*/
					for i := 0; i < len(cfg.NextHop); i++ {
						_, err = getIP(cfg.NextHop[i].NextHopIp)
						if err != nil {
							logger.Err("nextHopIpAddr invalid")
							return errors.New("Invalid next hop ip address")
						}
						/*
						   Check if next hop intf is valid L3 interface
						*/
						if cfg.NextHop[i].NextHopIntRef != "" {
							//logger.Debug("IntRef before : ", cfg.NextHop[i].NextHopIntRef)
							cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(cfg.NextHop[i].NextHopIntRef)
							if err != nil {
								logger.Err("Invalid NextHop IntRef ", cfg.NextHop[i].NextHopIntRef)
								return errors.New("Invalid Nexthop Intref")
							}
							//logger.Debug("IntRef after : ", cfg.NextHop[0].NextHopIntRef)
						} else {
							if len(oldcfg.NextHop) == 0 || len(oldcfg.NextHop) < i {
								logger.Err("Number of nextHops for old cfg < new cfg")
								return errors.New("number of nexthops not correct for update replace operation")
							}
							//logger.Debug("IntRef not provided, take the old value", oldcfg.NextHop[i].NextHopIntRef)
							cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(oldcfg.NextHop[i].NextHopIntRef)
							if err != nil {
								logger.Err("Invalid NextHop IntRef ", oldcfg.NextHop[i].NextHopIntRef)
								return errors.New("Invalid Nexthop Intref")
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (m RIBDServer) RouteConfigValidationCheckForPatchUpdate(oldcfg *ribd.IPv4Route, cfg *ribd.IPv4Route, op []*ribd.PatchOpInfo) (err error) {
	//logger.Info("RouteConfigValidationCheckForPatchUpdate")
	if netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot patch update ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API")
		return errors.New(fmt.Sprintln("Cannot add/remove from ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API"))
	}
	isCidr := strings.Contains(cfg.DestinationNw, "/")
	if isCidr {
		/*
		   the given address is in CIDR format
		*/
		ip, ipNet, err := net.ParseCIDR(cfg.DestinationNw)
		if err != nil {
			logger.Err("Invalid Destination IP address")
			return errors.New("Invalid Desitnation IP address")
		}
		_, err = getNetworkPrefixFromCIDR(cfg.DestinationNw)
		if err != nil {
			return errors.New("Invalid destination ip/network Mask")
		}
		cfg.DestinationNw = ip.String()
		oldcfg.DestinationNw = ip.String()
		ipMask := make(net.IP, 4)
		copy(ipMask, ipNet.Mask)
		ipMaskStr := net.IP(ipMask).String()
		cfg.NetworkMask = ipMaskStr
		oldcfg.NetworkMask = ipMaskStr
	}
	_, err = validateNetworkPrefix(cfg.DestinationNw, cfg.NetworkMask)
	if err != nil {
		logger.Err(" getNetowrkPrefixFromStrings returned err ", err)
		return errors.New("Invalid destination ip address")
	}
	for idx := 0; idx < len(op); idx++ {
		//logger.Debug("patch update")
		switch op[idx].Path {
		case "NextHop":
			//logger.Debug("Patch update for next hop")
			if len(op[idx].Value) == 0 {
				/*
					If route update is trying to add next hop, non zero nextHop info is expected
				*/
				logger.Err("Must specify next hop")
				return errors.New("Next hop ip not specified")
			}
			//logger.Debug("value = ", op[idx].Value)
			valueObjArr := []ribd.NextHopInfo{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				logger.Err("error unmarshaling value:", err)
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			//logger.Debug("Number of nextHops:", len(valueObjArr))
			for _, val := range valueObjArr {
				/*
				   Check if the next hop ip valid
				*/
				//logger.Debug("nextHop info: ip - ", val.NextHopIp, " intf: ", val.NextHopIntRef, " wt:", val.Weight)
				_, err = getIP(val.NextHopIp)
				if err != nil {
					logger.Err("nextHopIpAddr invalid")
					return errors.New("Invalid next hop ip address")
				}

				switch op[idx].Op {
				case "add":
					/*
					   Check if the next hop ref is valid L3 interface for add operation
					*/
					//logger.Debug("IntRef before : ", val.NextHopIntRef)
					val.NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(val.NextHopIntRef)
					if err != nil {
						logger.Err("Invalid NextHop IntRef ", val.NextHopIntRef)
						return errors.New("Invalid NextHop Intref")
					}
					//logger.Debug("IntRef after : ", val.NextHopIntRef)
				case "remove":
					//logger.Debug("remove op"))
				default:
					//logger.Err("operation ", op[idx].Op, " not supported")
					return errors.New(fmt.Sprintln("operation ", op[idx].Op, " not supported"))
				}
			}
		default:
			logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
			return errors.New("Invalid attribute for patch update")
		}
	}

	return nil
}

/*
    This function performs config parameters validation for op = "add" and "del" values.
	Key validations performed by this fucntion include:
	   - if the Protocol specified is valid (STATIC/CONNECTED/EBGP/OSPF)
	   - Validate destinationNw. If provided in CIDR notation, convert to ip addr and mask values
	   - In case of op == "del", check if the route is present in the DB
	   - for each of the nextHop info, check:
	       - if the next hop ip is valid
		   - if the nexthopIntf is valid L3 intf and if so, convert to string value
*/
func (m RIBDServer) RouteConfigValidationCheck(cfg *ribd.IPv4Route, op string) (err error) {
	if netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot create/delete ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API")
		return errors.New(fmt.Sprintln("Cannot create/delete ipv6 route (destination:", cfg.DestinationNw, ") using Ipv4Route API"))
	}
	isCidr := strings.Contains(cfg.DestinationNw, "/")
	if isCidr {
		/*
		   the given address is in CIDR format
		*/
		ip, ipNet, err := net.ParseCIDR(cfg.DestinationNw)
		if err != nil {
			logger.Err("Invalid Destination IP address")
			return errors.New("Invalid Desitnation IP address")
		}
		_, err = getNetworkPrefixFromCIDR(cfg.DestinationNw)
		if err != nil {
			return errors.New(fmt.Sprintln("RouteConfigValidationCheck for route:", cfg, " Invalid destination ip/network Mask"))
		}
		/*
		   Convert the CIDR format address to IP and mask strings
		*/
		cfg.DestinationNw = ip.String()
		ipMask := make(net.IP, 4)
		copy(ipMask, ipNet.Mask)
		ipMaskStr := net.IP(ipMask).String()
		cfg.NetworkMask = ipMaskStr
		/*
			In case where user provides CIDR address, the DB cannot verify if the route is present, so check here
		*/
		/*		if m.DbHdl != nil {
				var dbObjCfg objects.IPv4Route
				dbObjCfg.DestinationNw = cfg.DestinationNw
				dbObjCfg.NetworkMask = cfg.NetworkMask
				key := "IPv4Route#" + cfg.DestinationNw + "#" + cfg.NetworkMask
				_, err := m.DbHdl.GetObjectFromDb(dbObjCfg, key)
				if err == nil {
					logger.Err("Duplicate entry")
					return errors.New("Duplicate entry")
				}
			}*/
	}
	_, err = validateNetworkPrefix(cfg.DestinationNw, cfg.NetworkMask)
	if err != nil {
		logger.Err("RouteConfigValidationCheck for route:", cfg, " validateNetworkPrefix() returned err ", err)
		return err
	}
	/*
	   op is to add new route
	*/
	if op == "add" {
		/*
		   check if route protocol type is valid
		*/
		_, ok := RouteProtocolTypeMapDB[cfg.Protocol]
		if !ok {
			logger.Err("route type ", cfg.Protocol, " invalid")
			err = errors.New("Invalid route protocol type")
			return err
		}
		if cfg.NullRoute == true {
			logger.Debug("this is a null route, so dont validate nexthop attribute")
			if cfg.NextHop == nil || len(cfg.NextHop) == 0 {
				cfg.NextHop = make([]*ribd.NextHopInfo, 0)
				cfg.NextHop = append(cfg.NextHop, &ribd.NextHopInfo{
					NextHopIp: "255.255.255.255",
				})
			}
			return nil
		}
		//logger.Debug("Number of nexthops = ", len(cfg.NextHop))
		if len(cfg.NextHop) == 0 {
			/*
				Expects non-zero nexthop info
			*/
			logger.Err("Must specify next hop")
			return errors.New("Next hop ip not specified")
		}
		for i := 0; i < len(cfg.NextHop); i++ {
			/*
			   Check if the NextHop IP valid
			*/
			_, err = getIP(cfg.NextHop[i].NextHopIp)
			if err != nil {
				logger.Err("nextHopIpAddr invalid")
				return errors.New("Invalid next hop ip address")
			}
			//logger.Debug("IntRef before : ", cfg.NextHop[i].NextHopIntRef)
			/*
			   Validate if nextHopIntRef is a valid L3 interface
			*/
			if cfg.NextHop[i].NextHopIntRef == "" {
				//logger.Info("RouteConfigValidationCheck for route:", cfg, "NextHopIntRef not set")
				nhIntf, err := RouteServiceHandler.GetV4RouteReachabilityInfo(cfg.NextHop[i].NextHopIp, -1)
				if err != nil {
					logger.Err("RouteConfigValidationCheck for route:", cfg, "next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable")
					return errors.New(fmt.Sprintln("next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable"))
				}
				cfg.NextHop[i].NextHopIntRef = strconv.Itoa(int(nhIntf.NextHopIfIndex))
			} else {
				nhIntf := cfg.NextHop[i].NextHopIntRef
				cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(cfg.NextHop[i].NextHopIntRef)
				if err != nil {
					logger.Err("RouteConfigValidationCheck for route:", cfg, "Invalid NextHop IntRef ", cfg.NextHop[i].NextHopIntRef)
					return err
				}
				nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[i].NextHopIntRef)
				_, err := RouteServiceHandler.GetV4RouteReachabilityInfo(cfg.NextHop[i].NextHopIp, ribdInt.Int(nextHopIntRef))
				if err != nil {
					logger.Err("RouteConfigValidationCheck for route:", cfg, "next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable via interface ", nhIntf)
					return errors.New(fmt.Sprintln("next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable via ", nhIntf))
				}
			}
			//logger.Debug("IntRef after : ", cfg.NextHop[i].NextHopIntRef)
		}
	} else {
		if cfg.NullRoute == true {
			if cfg.NextHop == nil || len(cfg.NextHop) == 0 {
				cfg.NextHop = make([]*ribd.NextHopInfo, 0)
				cfg.NextHop = append(cfg.NextHop, &ribd.NextHopInfo{
					NextHopIp: "255.255.255.255",
				})
			}
			return nil
		}

	}
	return nil
}
func Getv4RoutesPerProtocol(protocol string) []*ribd.RouteInfoSummary {
	routes := make([]*ribd.RouteInfoSummary, 0)
	routemapInfo := ProtocolRouteMap[protocol]
	if routemapInfo.v4routeMap == nil {
		return routes
	}
	for destNetIp, val := range routemapInfo.v4routeMap {
		if val.totalcount == 0 {
			continue
		}
		v4Item := V4RouteInfoMap.Get(patriciaDB.Prefix(destNetIp))
		if v4Item == nil {
			continue
		}
		v4routeInfoRecordList := v4Item.(RouteInfoRecordList)
		v4protocolRouteList, ok := v4routeInfoRecordList.routeInfoProtocolMap[protocol]
		if !ok || len(v4protocolRouteList) == 0 {
			logger.Info("Unexpected: no route for destNet:", destNetIp, " found in routeMap of type:", protocol)
			continue
		}
		isInstalledinHw := true
		if v4routeInfoRecordList.selectedRouteProtocol != protocol {
			isInstalledinHw = false
		}
		destNet := ""
		nextHopList := make([]*ribd.NextHopInfo, 0)
		nextHopInfo := make([]ribd.NextHopInfo, len(v4protocolRouteList))
		i := 0
		for sel := 0; sel < len(v4protocolRouteList); sel++ {
			//logger.Info("protocol:", protocol, " sel:", sel, " v4protocolRouteList[sel].nextHopIp.String():", v4protocolRouteList[sel].nextHopIp.String(), " v4protocolRouteList[sel].nextHopIfIndex:", v4protocolRouteList[sel].nextHopIfIndex)
			destNet = v4protocolRouteList[sel].networkAddr
			nextHopInfo[i].NextHopIp = v4protocolRouteList[sel].nextHopIp.String()
			nextHopInfo[i].NextHopIntRef = strconv.Itoa(int(v4protocolRouteList[sel].nextHopIfIndex))
			intfEntry, ok := IntfIdNameMap[int32(v4protocolRouteList[sel].nextHopIfIndex)]
			if ok {
				nextHopInfo[i].NextHopIntRef = intfEntry.name
			}
			nextHopInfo[i].Weight = int32(v4protocolRouteList[sel].weight)
			nextHopList = append(nextHopList, &nextHopInfo[i])
			i++
		}
		routes = append(routes, &ribd.RouteInfoSummary{
			DestinationNw:   destNet,
			IsInstalledInHw: isInstalledinHw,
			NextHopList:     nextHopList,
		})
	}
	return routes
}
func Getv4RoutesPerInterface(intfref string) []string { //*ribd.RouteInfoSummary {
	routes := make([]string, 0) //[]*ribd.RouteInfoSummary, 0)
	routemapInfo := InterfaceRouteMap[intfref]
	if routemapInfo.v4routeMap == nil {
		return routes
	}
	for destNetIp, val := range routemapInfo.v4routeMap {
		if val.totalcount == 0 {
			continue
		}
		v4Item := V4RouteInfoMap.Get(patriciaDB.Prefix(destNetIp))
		if v4Item == nil {
			continue
		}
		v4routeInfoRecordList := v4Item.(RouteInfoRecordList)
		v4protocolRouteList, ok := v4routeInfoRecordList.routeInfoProtocolMap[v4routeInfoRecordList.selectedRouteProtocol]
		if !ok || len(v4protocolRouteList) == 0 {
			logger.Info("Unexpected: no route for destNet:", destNetIp, " found in routeMap of type:", intfref)
			continue
		}
		//		isInstalledinHw := true
		destNet := ""
		//	nextHopList := make([]*ribd.NextHopInfo, 0)
		//	nextHopInfo := make([]ribd.NextHopInfo, len(v4protocolRouteList))
		//	i := 0
		for sel := 0; sel < len(v4protocolRouteList); sel++ {
			destNet = v4protocolRouteList[sel].networkAddr
			/*			nextHopInfo[i].NextHopIp = v4protocolRouteList[sel].nextHopIp.String()
						nextHopInfo[i].NextHopIntRef = strconv.Itoa(int(v4protocolRouteList[sel].nextHopIfIndex))
						intfEntry, ok := IntfIdNameMap[int32(v4protocolRouteList[sel].nextHopIfIndex)]
						if ok {
							nextHopInfo[i].NextHopIntRef = intfEntry.name
						}
						nextHopInfo[i].Weight = int32(v4protocolRouteList[sel].weight)
						nextHopList = append(nextHopList, &nextHopInfo[i])
						i++*/
		}
		/*		routes = append(routes, &ribd.RouteInfoSummary{
				DestinationNw:   destNet,
				IsInstalledInHw: isInstalledinHw,
				NextHopList:     nextHopList,
			})*/
		routes = append(routes, destNet)
	}
	return routes
}

func (m RIBDServer) GetBulkIPv4RouteState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.IPv4RouteStateGetInfo, err error) { //(routes []*ribdInt.Routes, err error) {
	//logger.Debug("GetBulkIPv4RouteState")
	var i, validCount ribd.Int
	var toIndex ribd.Int
	var temproute []ribd.IPv4RouteState = make([]ribd.IPv4RouteState, rcount)
	var nextRoute *ribd.IPv4RouteState
	var returnRoutes []*ribd.IPv4RouteState
	var returnRouteGetInfo ribd.IPv4RouteStateGetInfo
	var prefixNodeRouteList RouteInfoRecordList
	var prefixNodeRoute RouteInfoRecord
	i = 0
	sel := 0
	found := false
	routes = &returnRouteGetInfo
	moreRoutes := true
	if destNetSlice == nil {
		//logger.Debug("destNetSlice not initialized: No Routes installed in RIB")
		return routes, err
	}
	for ; ; i++ {
		found = false
		if i+fromIndex >= ribd.Int(len(destNetSlice)) {
			//logger.Debug("All the routes fetched")
			moreRoutes = false
			break
		}
		/*		if destNetSlice[i+fromIndex].isValid == false {
				//logger.Debug("Invalid route")
				continue
			}*/
		if validCount == rcount {
			//logger.Debug("Enough routes fetched")
			break
		}
		prefixNode := V4RouteInfoMap.Get(destNetSlice[i+fromIndex].prefix)
		if prefixNode != nil {
			prefixNodeRouteList = prefixNode.(RouteInfoRecordList)
			if prefixNodeRouteList.isPolicyBasedStateValid == false {
				continue
			}
			//logger.Debug("selectedRouteProtocol = ", prefixNodeRouteList.selectedRouteProtocol)
			if prefixNodeRouteList.routeInfoProtocolMap == nil || prefixNodeRouteList.selectedRouteProtocol == "INVALID" || prefixNodeRouteList.routeInfoProtocolMap[prefixNodeRouteList.selectedRouteProtocol] == nil {
				//logger.Debug("selected route not valid")
				continue
			}
			routeInfoList := prefixNodeRouteList.routeInfoProtocolMap[prefixNodeRouteList.selectedRouteProtocol]
			for sel = 0; sel < len(routeInfoList); sel++ {
				if routeInfoList[sel].nextHopIp.String() == destNetSlice[i+fromIndex].nextHopIp {
					//logger.Debug("Found the entry corresponding to the nextHop ip")
					found = true
					break
				}
			}
			if !found {
				//logger.Debug("The corresponding route with nextHopIP was not found in the record DB")
				continue
			}
			prefixNodeRoute = routeInfoList[sel] //prefixNodeRouteList.routeInfoList[prefixNodeRouteList.selectedRouteIdx]
			nextRoute = &temproute[validCount]
			nextRoute.DestinationNw = prefixNodeRoute.networkAddr
			nextRoute.RouteCreatedTime = prefixNodeRoute.routeCreatedTime
			nextRoute.RouteUpdatedTime = prefixNodeRoute.routeUpdatedTime
			nextRoute.IsNetworkReachable = prefixNodeRoute.resolvedNextHopIpIntf.IsReachable
			nextRoute.PolicyList = make([]string, 0)
			routePolicyListInfo := ""
			if prefixNodeRouteList.policyList != nil {
				for k := 0; k < len(prefixNodeRouteList.policyList); k++ {
					routePolicyListInfo = "policy " + prefixNodeRouteList.policyList[k] + "["
					policyRouteIndex := PolicyRouteIndex{destNetIP: prefixNodeRoute.networkAddr, policy: prefixNodeRouteList.policyList[k]}
					policyStmtMap, ok := PolicyEngineDB.PolicyEntityMap[policyRouteIndex]
					if !ok || policyStmtMap.PolicyStmtMap == nil {
						continue
					}
					routePolicyListInfo = routePolicyListInfo + " stmtlist[["
					for stmt, conditionsAndActionsList := range policyStmtMap.PolicyStmtMap {
						routePolicyListInfo = routePolicyListInfo + stmt + ":[conditions:"
						for c := 0; c < len(conditionsAndActionsList.ConditionList); c++ {
							routePolicyListInfo = routePolicyListInfo + conditionsAndActionsList.ConditionList[c].Name + ","
						}
						routePolicyListInfo = routePolicyListInfo + "],[actions:"
						for a := 0; a < len(conditionsAndActionsList.ActionList); a++ {
							routePolicyListInfo = routePolicyListInfo + conditionsAndActionsList.ActionList[a].Name + ","
						}
						routePolicyListInfo = routePolicyListInfo + "]]"
					}
					routePolicyListInfo = routePolicyListInfo + "]"
					nextRoute.PolicyList = append(nextRoute.PolicyList, routePolicyListInfo)
				}
			}
			toIndex = ribd.Int(i + fromIndex)
			if len(returnRoutes) == 0 {
				returnRoutes = make([]*ribd.IPv4RouteState, 0)
			}
			returnRoutes = append(returnRoutes, nextRoute)
			validCount++
		}
	}
	routes.IPv4RouteStateList = returnRoutes
	routes.StartIdx = fromIndex
	routes.EndIdx = toIndex + 1
	routes.More = moreRoutes
	routes.Count = validCount
	return routes, err
}

func (m RIBDServer) Getv4Route(destNetIp string) (route *ribdInt.IPv4RouteState, err error) {
	var returnRoute ribdInt.IPv4RouteState
	route = &returnRoute
	/*
	   the given address is in CIDR format
	*/
	destNet, err := getNetworkPrefixFromCIDR(destNetIp)
	if err != nil {
		return route, errors.New("Invalid destination ip/network Mask")
	}
	routeInfoRecordListItem := V4RouteInfoMap.Get(destNet)
	if routeInfoRecordListItem == nil {
		logger.Err("No such route")
		err = errors.New("Route does not exist")
		return route, err
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList) //RouteInfoMap.Get(destNet).(RouteInfoRecordList)
	if routeInfoRecordList.selectedRouteProtocol == "INVALID" {
		logger.Err("No selected route for this network")
		err = errors.New("No selected route for this network")
		return route, err
	}
	routeInfoList := routeInfoRecordList.routeInfoProtocolMap[routeInfoRecordList.selectedRouteProtocol]
	nextHopInfo := make([]ribdInt.RouteNextHopInfo, len(routeInfoList))
	route.NextHopList = make([]*ribdInt.RouteNextHopInfo, 0)
	i := 0
	for _, nh := range routeInfoList {
		routeInfoRecord := nh
		nextHopInfo[i].NextHopIp = routeInfoRecord.nextHopIp.String()
		nextHopInfo[i].NextHopIntRef = strconv.Itoa(int(routeInfoRecord.nextHopIfIndex))
		intfEntry, ok := IntfIdNameMap[int32(routeInfoRecord.nextHopIfIndex)]
		if ok {
			//logger.Debug("Map found for ifndex : ", routeInfoRecord.nextHopIfIndex, "Name = ", intfEntry.name)
			nextHopInfo[i].NextHopIntRef = intfEntry.name
		}
		//logger.Debug("IntfRef = ", nextHopInfo[i].NextHopIntRef)
		nextHopInfo[i].Weight = int32(routeInfoRecord.weight)
		route.NextHopList = append(route.NextHopList, &nextHopInfo[i])
		i++

	}
	routeInfoRecord := routeInfoList[0]
	route.DestinationNw = routeInfoRecord.networkAddr
	route.Protocol = routeInfoRecordList.selectedRouteProtocol
	route.RouteCreatedTime = routeInfoRecord.routeCreatedTime
	route.RouteUpdatedTime = routeInfoRecord.routeUpdatedTime
	route.NextBestRoute = &ribdInt.NextBestRouteInfo{}
	route.NextBestRoute.Protocol = SelectNextBestRoute(routeInfoRecordList, routeInfoRecordList.selectedRouteProtocol)
	nextbestrouteInfoList := routeInfoRecordList.routeInfoProtocolMap[route.NextBestRoute.Protocol]
	//logger.Info("len of routeInfoList - ", len(routeInfoList), "selected route protocol = ", routeList.selectedRouteProtocol, " route Protocol: ", entry.protocol, " route nwAddr: ", entry.networkAddr)
	nextBestRouteNextHopInfo := make([]ribdInt.RouteNextHopInfo, len(nextbestrouteInfoList))
	i1 := 0
	for sel1 := 0; sel1 < len(nextbestrouteInfoList); sel1++ {
		//logger.Info("nextHop ", sel, " weight = ", routeInfoList[sel].weight, " ip ", routeInfoList[sel].nextHopIp, " intref ", routeInfoList[sel].nextHopIfIndex)
		nextBestRouteNextHopInfo[i1].NextHopIp = nextbestrouteInfoList[sel1].nextHopIp.String()
		nextBestRouteNextHopInfo[i1].NextHopIntRef = strconv.Itoa(int(nextbestrouteInfoList[sel1].nextHopIfIndex))
		intfEntry, ok := IntfIdNameMap[int32(nextbestrouteInfoList[sel1].nextHopIfIndex)]
		if ok {
			//logger.Debug("Map foud for ifndex : ", routeInfoList[sel].nextHopIfIndex, "Name = ", intfEntry.name)
			nextBestRouteNextHopInfo[i1].NextHopIntRef = intfEntry.name
		}
		//logger.Debug("IntfRef = ", nextHopInfo[i].NextHopIntRef)
		nextBestRouteNextHopInfo[i1].Weight = int32(nextbestrouteInfoList[sel1].weight)
		route.NextBestRoute.NextHopList = append(route.NextBestRoute.NextHopList, &nextBestRouteNextHopInfo[i1])
		i1++
	}
	return route, err
}
func (m RIBDServer) GetTotalv4RouteCount() (number int, err error) {
	return v4rtCount, err
}
func (m RIBDServer) Getv4RouteCreatedTime(number int) (time string, err error) {
	_, ok := v4routeCreatedTimeMap[number]
	if !ok {
		logger.Info(number, " number of  v4 routes not created yet")
		return "", errors.New("Not enough v4 routes")
	}
	return v4routeCreatedTimeMap[number], err
}

func (m RIBDServer) ProcessV4RouteCreateConfig(cfg *ribd.IPv4Route, addType int, sliceIdx ribd.Int) (val bool, err error) {
	logger.Debug("ProcessV4RouteCreateConfig: Received create route request for ip ", cfg.DestinationNw, " mask ", cfg.NetworkMask, " number of next hops: ", len(cfg.NextHop), " null Route:", cfg.NullRoute, " sliceIdx:", sliceIdx)
	newCfg := ribd.IPv4Route{
		DestinationNw: cfg.DestinationNw,
		NetworkMask:   cfg.NetworkMask,
		Protocol:      cfg.Protocol,
		Cost:          cfg.Cost,
		NullRoute:     cfg.NullRoute,
	}
	for i := 0; i < len(cfg.NextHop); i++ {
		logger.Debug("nexthop info: ip: ", cfg.NextHop[i].NextHopIp, " intref: ", cfg.NextHop[i].NextHopIntRef)
		nh := ribd.NextHopInfo{
			NextHopIp:     cfg.NextHop[i].NextHopIp,
			NextHopIntRef: cfg.NextHop[i].NextHopIntRef,
			Weight:        cfg.NextHop[i].Weight,
		}
		newCfg.NextHop = make([]*ribd.NextHopInfo, 0)
		newCfg.NextHop = append(newCfg.NextHop, &nh)
		//policyRoute := BuildPolicyRouteFromribdIPv4Route(&newCfg)
		params := BuildRouteParamsFromribdIPv4Route(&newCfg, addType, Invalid, sliceIdx)
		_, err = createRoute(params)
	}

	return true, err
}

func (m RIBDServer) ProcessBulkRouteCreateConfig(bulkCfg []*ribdInt.IPv4RouteConfig) (val bool, err error) {
	logger.Debug("ProcessBulkRouteCreateConfig: Received create route request for  ", len(bulkCfg), " number of routes")
	index := 0
	for _, cfg := range bulkCfg {

		newCfg := ribd.IPv4Route{
			DestinationNw: cfg.DestinationNw,
			NetworkMask:   cfg.NetworkMask,
			Protocol:      cfg.Protocol,
			Cost:          cfg.Cost,
			NullRoute:     cfg.NullRoute,
		}
		for i := 0; i < len(cfg.NextHop); i++ {
			logger.Debug("nexthop info: ip: ", cfg.NextHop[i].NextHopIp, " intref: ", cfg.NextHop[i].NextHopIntRef)
			nh := ribd.NextHopInfo{
				NextHopIp:     cfg.NextHop[i].NextHopIp,
				NextHopIntRef: cfg.NextHop[i].NextHopIntRef,
				Weight:        cfg.NextHop[i].Weight,
			}
			newCfg.NextHop = make([]*ribd.NextHopInfo, 0)
			newCfg.NextHop = append(newCfg.NextHop, &nh)
		}

		//policyRoute := BuildPolicyRouteFromribdIPv4Route(&newCfg)
		params := BuildRouteParamsFromribdIPv4Route(&newCfg, FIBAndRIB, Invalid, ribd.Int(len(destNetSlice)))
		params.bulk = true
		index++
		if index == len(bulkCfg) {
			params.bulkEnd = true
		}
		//logger.Debug("createType = ", params.createType, "deleteType = ", params.deleteType, "index:", index, " bulk:", params.bulk, " bulkEnd:", params.bulkEnd))
		//PolicyEngineFilter(policyRoute, policyCommonDefs.PolicyPath_Import, params)
		_, err = createRoute(params)
	}

	return true, err
}

func (m RIBDServer) ProcessV4RouteDeleteConfig(cfg *ribd.IPv4Route, delType int) (val bool, err error) {
	logger.Debug("ProcessV4RouteDeleteConfig:Received Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "number of nextHops:", len(cfg.NextHop), "Protocol ", cfg.Protocol)
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return 0,err
	}
	var nextHopIfIndex ribd.Int
	for i := 0; i < len(cfg.NextHop); i++ {
		if cfg.NullRoute == true { //commonDefs.IfTypeNull {
			logger.Info("null route create request")
			cfg.NextHop[i].NextHopIp = "255.255.255.255"
		}
		logger.Debug("nexthop info: ip: ", cfg.NextHop[i].NextHopIp, " intref: ", cfg.NextHop[i].NextHopIntRef)
		nextHopIfIndex = -1
		if cfg.NextHop[i].NextHopIntRef != "" {
			cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(cfg.NextHop[i].NextHopIntRef)
			if err != nil {
				logger.Err(fmt.Sprintln("Invalid NextHop IntRef ", cfg.NextHop[i].NextHopIntRef))
				return false, err
			}
			nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[i].NextHopIntRef)
			nextHopIfIndex = ribd.Int(nextHopIntRef)
		}
		_, err = deleteIPRoute(cfg.DestinationNw, defs.IPv4, cfg.NetworkMask, cfg.Protocol, cfg.NextHop[i].NextHopIp, nextHopIfIndex, ribd.Int(delType), defs.RoutePolicyStateChangetoInValid)
	}
	return true, err
}

func (m RIBDServer) Processv4RoutePatchUpdateConfig(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, op []*ribd.PatchOpInfo) (ret bool, err error) {
	logger.Debug("Processv4RoutePatchUpdateConfig:Received update route request with number of patch ops: ", len(op))
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return err
	}
	destNet, err := getNetowrkPrefixFromStrings(origconfig.DestinationNw, origconfig.NetworkMask)
	if err != nil {
		logger.Err(" getNetowrkPrefixFromStrings returned err ", err)
		return ret, err
	}
	ok := V4RouteInfoMap.Match(destNet)
	if !ok {
		err = errors.New("Processv4RoutePatchUpdateConfig:No route found")
		return ret, err
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "NextHop":
			//logger.Debug("Patch update for next hop")
			/*newconfig should only have the next hops that have to be added or deleted*/
			newconfig.NextHop = make([]*ribd.NextHopInfo, 0)
			//logger.Debug("value = ", op[idx].Value)
			valueObjArr := []ribd.NextHopInfo{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				//logger.Debug("error unmarshaling value:", err))
				return ret, errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			//logger.Debug("Number of nextHops:", len(valueObjArr)))
			for _, val := range valueObjArr {
				//logger.Debug("nextHop info: ip - ", val.NextHopIp, " intf: ", val.NextHopIntRef, " wt:", val.Weight)
				//wt,_ := strconv.Atoi((op[idx].Value[j]["Weight"]))
				//logger.Debug("IntRef before : ", val.NextHopIntRef)
				val.NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(val.NextHopIntRef)
				if err != nil {
					logger.Err("Invalid NextHop IntRef ", val.NextHopIntRef)
					return ret, errors.New("Invalid NextHop Intref")
				}
				//logger.Debug("IntRef after : ", val.NextHopIntRef)
				nh := ribd.NextHopInfo{
					NextHopIp:     val.NextHopIp,
					NextHopIntRef: val.NextHopIntRef,
					Weight:        val.Weight,
				}
				newconfig.NextHop = append(newconfig.NextHop, &nh)
			}
			switch op[idx].Op {
			case "add":
				m.ProcessV4RouteCreateConfig(newconfig, FIBAndRIB, ribd.Int(len(destNetSlice)))
			case "remove":
				m.ProcessV4RouteDeleteConfig(newconfig, FIBAndRIB)
			default:
				logger.Err("Operation ", op[idx].Op, " not supported")
			}
		default:
			logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
			err = errors.New(fmt.Sprintln("Operation ", op[idx].Op, " not supported"))
		}
	}
	return ret, err
}

func (m RIBDServer) Processv4RouteUpdateConfig(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, attrset []bool) (val bool, err error) {
	logger.Debug("Processv4RouteUpdateConfig:Received update route request ")
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return err
	}
	destNet, err := getNetowrkPrefixFromStrings(origconfig.DestinationNw, origconfig.NetworkMask)
	if err != nil {
		logger.Err(" getNetowrkPrefixFromStrings returned err ", err)
		return val, err
	}
	ok := V4RouteInfoMap.Match(destNet)
	if !ok {
		err = errors.New(fmt.Sprintln("No route found for ip ", destNet))
		return val, err
	}
	routeInfoRecordListItem := V4RouteInfoMap.Get(destNet)
	if routeInfoRecordListItem == nil {
		logger.Err("No route for destination network", destNet)
		return val, err
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	callUpdate := true
	if attrset != nil {
		found, routeInfoRecord, index := findRouteWithNextHop(routeInfoRecordList.routeInfoProtocolMap[origconfig.Protocol], defs.IPv4, origconfig.NextHop[0].NextHopIp, -1)
		if !found || index == -1 {
			logger.Err("Invalid nextHopIP")
			return val, errors.New(fmt.Sprintln("Invalid Next Hop IP:", origconfig.NextHop[0].NextHopIp))
		}
		objTyp := reflect.TypeOf(*origconfig)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				//logger.Debug(fmt.Sprintf("ProcessRouteUpdateConfig (server): changed ", objName)
				if objName == "NextHop" {
					if len(newconfig.NextHop) == 0 {
						logger.Err("Must specify next hop")
						return val, err
					} else {
						nextHopIpAddr, err := getIP(newconfig.NextHop[0].NextHopIp)
						if err != nil {
							logger.Err("nextHopIpAddr invalid")
							return val, errors.New("Invalid next hop")
						}
						//logger.Debug("Update the next hop info old ip: ", origconfig.NextHop[0].NextHopIp, " new value: ", newconfig.NextHop[0].NextHopIp, " weight : ", newconfig.NextHop[0].Weight)
						routeInfoRecord.nextHopIp = nextHopIpAddr
						routeInfoRecord.weight = ribd.Int(newconfig.NextHop[0].Weight)
						if newconfig.NextHop[0].NextHopIntRef != "" {
							nextHopIntRef, _ := strconv.Atoi(newconfig.NextHop[0].NextHopIntRef)
							routeInfoRecord.nextHopIfIndex = ribd.Int(nextHopIntRef)
						}
					}
				}
				if objName == "Cost" {
					routeInfoRecord.metric = ribd.Int(newconfig.Cost)
				}
			}
		}
		routeInfoRecordList.routeInfoProtocolMap[origconfig.Protocol][index] = routeInfoRecord
		V4RouteInfoMap.Set(destNet, routeInfoRecordList)
		//logger.Debug("Adding to DBRouteCh from processRouteUpdateConfig")
		RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
			OrigConfigObject: RouteDBInfo{routeInfoRecord, routeInfoRecordList},
			Op:               defs.Add,
		}
		if callUpdate == false {
			return val, err
		}
	}
	updateBestRoute(destNet, routeInfoRecordList)
	return val, err
}
