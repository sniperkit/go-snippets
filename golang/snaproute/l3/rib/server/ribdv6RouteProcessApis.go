// ribdv6RouteProcessApis.go
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	defs "l3/rib/ribdCommonDefs"
	// "models/objects"
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

var V6RouteInfoMap *patriciaDB.Trie //Routes are stored in patricia trie

/*
   Returns the longest prefix match route to reach the destination network destNet
*/
func (m RIBDServer) GetV6RouteReachabilityInfo(destNet string, ifIndex ribdInt.Int) (nextHopIntf *ribdInt.NextHopInfo, err error) {
	//logger.Debug("GetRouteReachabilityInfo of ", destNet)
	//t1 := time.Now()
	var retnextHopIntf ribdInt.NextHopInfo
	nextHopIntf = &retnextHopIntf
	var found bool
	destNetIp, err := getIP(destNet)
	if err != nil {
		logger.Err("getIP returned Invalid dest ip address for ", destNet)
		return nextHopIntf, errors.New("Invalid dest ip address")
	}
	lookupIp := destNetIp.To16()
	if lookupIp == nil {
		logger.Err("getIP returned Invalid dest ip address for ", destNet)
		return nextHopIntf, errors.New("Invalid dest ip address")
	}
	destNetIp = lookupIp
	rmapInfoListItem := V6RouteInfoMap.GetLongestPrefixNode(patriciaDB.Prefix(destNetIp))
	if rmapInfoListItem != nil {
		rmapInfoList := rmapInfoListItem.(RouteInfoRecordList)
		if rmapInfoList.selectedRouteProtocol != "INVALID" {
			found = true
			routeInfoList, ok := rmapInfoList.routeInfoProtocolMap[rmapInfoList.selectedRouteProtocol]
			if !ok || len(routeInfoList) == 0 {
				logger.Err("Selected route not found")
				return nil, errors.New("dest ip address not reachable")
			}
			nhFound, v, _ := findRouteWithNextHop(routeInfoList, defs.IPv6, "", ribd.Int(ifIndex))
			if !nhFound {
				logger.Err("Selected route not found")
				return nil, errors.New(fmt.Sprintln("dest ip address not reachable via ifIndex", ifIndex))
			}
			//		v := routeInfoList[0]
			nextHopIntf.NextHopIp = v.nextHopIp.String()
			nextHopIntf.NextHopIfIndex = ribdInt.Int(v.nextHopIfIndex)
			nextHopIntf.Metric = ribdInt.Int(v.metric)
			nextHopIntf.Ipaddr = v.destNetIp.String()
			nextHopIntf.Mask = v.networkMask.String()
			nextHopIntf.IsReachable = true
		}
	}

	if found == false {
		//logger.Err("dest IP", destNetIp, " not reachable ")
		err = errors.New("dest ip address not reachable")
		return nextHopIntf, err
	}
	//	duration := time.Since(t1)
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
func UpdateV6RouteReachabilityStatus(prefix patriciaDB.Prefix, //prefix of the node being traversed
	handle patriciaDB.Item, //data interface (routeInforRecordList) for this node
	item patriciaDB.Item) /*RouteReachabilityStatusInfo data of the v6 route that is being checked with*/ (err error) {

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
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	destIpPrefix, err := getNetowrkPrefixFromStrings(ipAddrStr, ipMaskStr)
	if err != nil {
		logger.Err("Error getting ip prefix for ip:", ipAddrStr, " mask:", ipMaskStr)
		return err
	}
	//logger.Debug("UpdateRouteReachabilityStatus network: ", routeReachabilityStatusInfo.destNet, " status:", routeReachabilityStatusInfo.status, "ip: ", ip.String(), " destIPPrefix: ", destIpPrefix, " ipMaskStr:", ipMaskStr)
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
					V6RouteInfoMap.Set(prefix, rmapInfoRecordList)
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
						V6RouteInfoMap.VisitAndUpdate(UpdateV6RouteReachabilityStatus, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Down", k, nextHopIntf})
					}
				} else if routeReachabilityStatusInfo.status == "Up" && v[i].resolvedNextHopIpIntf.IsReachable == false {
					//logger.Debug("Bringing up route : ip: ", v[i].networkAddr)
					v[i].resolvedNextHopIpIntf.IsReachable = true
					rmapInfoRecordList.routeInfoProtocolMap[k] = v
					V6RouteInfoMap.Set(prefix, rmapInfoRecordList)
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
						V6RouteInfoMap.VisitAndUpdate(UpdateV6RouteReachabilityStatus, RouteReachabilityStatusInfo{v[i].networkAddr, v[i].ipType, "Up", k, nextHopIntf})
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
func (m RIBDServer) IPv6RouteConfigValidationCheckForUpdate(oldcfg *ribd.IPv6Route, cfg *ribd.IPv6Route, attrset []bool) (err error) {
	logger.Info("IPv6RouteConfigValidationCheckForUpdate")
	if !netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot update ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API")
		return errors.New(fmt.Sprintln("Cannot update ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API"))
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
		ipMask := make(net.IP, 16)
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
		logger.Debug("attr set not nil, set individual attributes")
		objTyp := reflect.TypeOf(*cfg)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				logger.Debug(fmt.Sprintf("Processv6RouteUpdateConfig (server validation): changed ", objName))
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
							logger.Err(fmt.Sprintln("nextHopIpAddr invalid"))
							return errors.New("Invalid next hop ip address")
						}
						/*
						   Check if next hop intf is valid L3 interface
						*/
						if cfg.NextHop[i].NextHopIntRef != "" {
							logger.Debug(fmt.Sprintln("IntRef before : ", cfg.NextHop[i].NextHopIntRef))
							cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(cfg.NextHop[i].NextHopIntRef)
							if err != nil {
								logger.Err(fmt.Sprintln("Invalid NextHop IntRef ", cfg.NextHop[i].NextHopIntRef))
								return errors.New("Invalid Nexthop Intref")
							}
							logger.Debug(fmt.Sprintln("IntRef after : ", cfg.NextHop[0].NextHopIntRef))
						} else {
							if len(oldcfg.NextHop) == 0 || len(oldcfg.NextHop) < i {
								logger.Err("Number of nextHops for old cfg < new cfg")
								return errors.New("number of nexthops not correct for update replace operation")
							}
							logger.Debug(fmt.Sprintln("IntRef not provided, take the old value", oldcfg.NextHop[i].NextHopIntRef))
							cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(oldcfg.NextHop[i].NextHopIntRef)
							if err != nil {
								logger.Err(fmt.Sprintln("Invalid NextHop IntRef ", oldcfg.NextHop[i].NextHopIntRef))
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

func (m RIBDServer) IPv6RouteConfigValidationCheckForPatchUpdate(oldcfg *ribd.IPv6Route, cfg *ribd.IPv6Route, op []*ribd.PatchOpInfo) (err error) {
	logger.Info(fmt.Sprintln("IPv6RouteConfigValidationCheckForPatchUpdate"))
	if !netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot patch update ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API")
		return errors.New(fmt.Sprintln("Cannot add/remove from ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API"))
	}
	isCidr := strings.Contains(cfg.DestinationNw, "/")
	if isCidr {
		logger.Debug("cidr address")
		/*
		   the given address is in CIDR format
		*/
		ip, ipNet, err := net.ParseCIDR(cfg.DestinationNw)
		if err != nil {
			logger.Err(fmt.Sprintln("Invalid Destination IP address"))
			return errors.New("Invalid Desitnation IP address")
		}
		_, err = getNetworkPrefixFromCIDR(cfg.DestinationNw)
		if err != nil {
			return errors.New("Invalid destination ip/network Mask")
		}
		logger.Debug(fmt.Sprintln("At the beginning oldcfg.destinationnw:", oldcfg.DestinationNw, " cfg.DesinationNw:", cfg.DestinationNw))
		cfg.DestinationNw = ip.String()
		oldcfg.DestinationNw = ip.String()
		ipMask := make(net.IP, 16)
		copy(ipMask, ipNet.Mask)
		ipMaskStr := net.IP(ipMask).String()
		cfg.NetworkMask = ipMaskStr
		oldcfg.NetworkMask = ipMaskStr
		logger.Debug(fmt.Sprintln("After conversion oldcfg.destinationnw:", oldcfg.DestinationNw, " cfg.DesinationNw:", cfg.DestinationNw))
	}
	_, err = validateNetworkPrefix(cfg.DestinationNw, cfg.NetworkMask)
	if err != nil {
		logger.Info(fmt.Sprintln(" getNetowrkPrefixFromStrings returned err ", err))
		return errors.New("Invalid destination ip address")
	}
	for idx := 0; idx < len(op); idx++ {
		logger.Debug(fmt.Sprintln("patch update"))
		switch op[idx].Path {
		case "NextHop":
			logger.Debug("Patch update for next hop")
			if len(op[idx].Value) == 0 {
				/*
					If route update is trying to add next hop, non zero nextHop info is expected
				*/
				logger.Err("Must specify next hop")
				return errors.New("Next hop ip not specified")
			}
			logger.Debug(fmt.Sprintln("value = ", op[idx].Value))
			valueObjArr := []ribd.NextHopInfo{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				logger.Debug(fmt.Sprintln("error unmarshaling value:", err))
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug(fmt.Sprintln("Number of nextHops:", len(valueObjArr)))
			for _, val := range valueObjArr {
				/*
				   Check if the next hop ip valid
				*/
				logger.Debug(fmt.Sprintln("nextHop info: ip - ", val.NextHopIp, " intf: ", val.NextHopIntRef, " wt:", val.Weight))
				_, err = getIP(val.NextHopIp)
				if err != nil {
					logger.Err(fmt.Sprintln("nextHopIpAddr invalid"))
					return errors.New("Invalid next hop ip address")
				}

				switch op[idx].Op {
				case "add":
					/*
					   Check if the next hop ref is valid L3 interface for add operation
					*/
					logger.Debug(fmt.Sprintln("IntRef before : ", val.NextHopIntRef))
					if val.NextHopIntRef == "" {
						logger.Info(fmt.Sprintln("NextHopIntRef not set"))
						nhIntf, err := RouteServiceHandler.GetV6RouteReachabilityInfo(val.NextHopIp, -1)
						if err != nil {
							logger.Err(fmt.Sprintln("next hop ip ", val.NextHopIp, " not reachable"))
							return errors.New(fmt.Sprintln("next hop ip ", val.NextHopIp, " not reachable"))
						}
						val.NextHopIntRef = strconv.Itoa(int(nhIntf.NextHopIfIndex))
					} else {
						val.NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(val.NextHopIntRef)
						if err != nil {
							logger.Err(fmt.Sprintln("Invalid NextHop IntRef ", val.NextHopIntRef))
							return err
						}
					}
					logger.Debug(fmt.Sprintln("IntRef after : ", val.NextHopIntRef))
				case "remove":
					logger.Debug(fmt.Sprintln("remove op"))
				default:
					logger.Err(fmt.Sprintln("operation ", op[idx].Op, " not supported"))
					return errors.New(fmt.Sprintln("operation ", op[idx].Op, " not supported"))
				}
			}
		default:
			logger.Err(fmt.Sprintln("Patch update for attribute:", op[idx].Path, " not supported"))
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
func (m RIBDServer) IPv6RouteConfigValidationCheck(cfg *ribd.IPv6Route, op string) (err error) {
	logger.Debug(fmt.Sprintln("IPv6RouteConfigValidationCheck"))
	if !netutils.IsIPv6Addr(cfg.DestinationNw) {
		logger.Err("Cannot create/delete ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API")
		return errors.New(fmt.Sprintln("Cannot create/delete ipv4 route (destination:", cfg.DestinationNw, ") using Ipv6Route API"))
	}
	isCidr := strings.Contains(cfg.DestinationNw, "/")
	if isCidr {
		/*
		   the given address is in CIDR format
		*/
		ip, ipNet, err := net.ParseCIDR(cfg.DestinationNw)
		if err != nil {
			logger.Err(fmt.Sprintln("Invalid Destination IP address"))
			return errors.New("Invalid Desitnation IP address")
		}
		_, err = getNetworkPrefixFromCIDR(cfg.DestinationNw)
		if err != nil {
			return errors.New("Invalid destination ip/network Mask")
		}
		/*
		   Convert the CIDR format address to IP and mask strings
		*/
		cfg.DestinationNw = ip.String()
		ipMask := make(net.IP, 16)
		copy(ipMask, ipNet.Mask)
		ipMaskStr := net.IP(ipMask).String()
		cfg.NetworkMask = ipMaskStr
		/*
			In case where user provides CIDR address, the DB cannot verify if the route is present, so check here
		*/
		/*		if m.DbHdl != nil {
				var dbObjCfg objects.IPv6Route
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
		logger.Info(fmt.Sprintln(" getNetowrkPrefixFromStrings returned err ", err))
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
			logger.Err(fmt.Sprintln("route type ", cfg.Protocol, " invalid"))
			err = errors.New("Invalid route protocol type")
			return err
		}
		logger.Debug(fmt.Sprintln("Number of nexthops = ", len(cfg.NextHop)))
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
				logger.Err(fmt.Sprintln("nextHopIpAddr invalid"))
				return errors.New("Invalid next hop ip address")
			}
			logger.Debug(fmt.Sprintln("IntRef before : ", cfg.NextHop[i].NextHopIntRef))
			/*
			   Validate if nextHopIntRef is a valid L3 interface
			*/
			if cfg.NextHop[i].NextHopIntRef == "" {
				logger.Info(fmt.Sprintln("NextHopIntRef not set"))
				nhIntf, err := RouteServiceHandler.GetRouteReachabilityInfo(cfg.NextHop[i].NextHopIp, -1)
				if err != nil {
					logger.Err(fmt.Sprintln("next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable"))
					return errors.New(fmt.Sprintln("next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable"))
				}
				cfg.NextHop[i].NextHopIntRef = strconv.Itoa(int(nhIntf.NextHopIfIndex))
			} else {
				nhIntf := cfg.NextHop[i].NextHopIntRef
				cfg.NextHop[i].NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(cfg.NextHop[i].NextHopIntRef)
				if err != nil {
					logger.Err(fmt.Sprintln("Invalid NextHop IntRef ", cfg.NextHop[i].NextHopIntRef))
					return err
				}
				nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[i].NextHopIntRef)
				_, err := RouteServiceHandler.GetRouteReachabilityInfo(cfg.NextHop[i].NextHopIp, ribdInt.Int(nextHopIntRef))
				if err != nil {
					logger.Err("RouteConfigValidationCheck for route:", cfg, "next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable via interface ", nhIntf)
					return errors.New(fmt.Sprintln("next hop ip ", cfg.NextHop[i].NextHopIp, " not reachable via ", nhIntf))
				}
			}
			//logger.Debug(fmt.Sprintln("IntRef after : ", cfg.NextHop[i].NextHopIntRef))
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
func (m RIBDServer) GetTotalv6RouteCount() (number int, err error) {
	return v6rtCount, err
}
func (m RIBDServer) Getv6RouteCreatedTime(number int) (time string, err error) {
	_, ok := v6routeCreatedTimeMap[number]
	if !ok {
		logger.Info(number, " number of  v6 routes not created yet")
		return "", errors.New("Not enough v6 routes")
	}
	return v6routeCreatedTimeMap[number], err
}
func Getv6RoutesPerProtocol(protocol string) []*ribd.RouteInfoSummary {
	v6routes := make([]*ribd.RouteInfoSummary, 0)
	routemapInfo := ProtocolRouteMap[protocol]
	if routemapInfo.v6routeMap == nil {
		return v6routes
	}
	for destNetIp, val := range routemapInfo.v6routeMap {
		if val.totalcount == 0 {
			continue
		}
		v6Item := V6RouteInfoMap.Get(patriciaDB.Prefix(destNetIp))
		if v6Item == nil {
			continue
		}
		v6routeInfoRecordList := v6Item.(RouteInfoRecordList)
		v6protocolRouteList, ok := v6routeInfoRecordList.routeInfoProtocolMap[protocol]
		if !ok || len(v6protocolRouteList) == 0 {
			logger.Info("Unexpected: no route for destNet:", destNetIp, " found in routeMap of type:", protocol)
			continue
		}
		v6isInstalledinHw := true
		if v6routeInfoRecordList.selectedRouteProtocol != protocol {
			v6isInstalledinHw = false
		}
		v6destNet := ""
		v6nextHopList := make([]*ribd.NextHopInfo, 0)
		v6nextHopInfo := make([]ribd.NextHopInfo, len(v6protocolRouteList))
		j := 0
		for sel1 := 0; sel1 < len(v6protocolRouteList); sel1++ {
			v6destNet = v6protocolRouteList[sel1].networkAddr
			v6nextHopInfo[j].NextHopIp = v6protocolRouteList[sel1].nextHopIp.String()
			v6nextHopInfo[j].NextHopIntRef = strconv.Itoa(int(v6protocolRouteList[sel1].nextHopIfIndex))
			intfEntry, ok := IntfIdNameMap[int32(v6protocolRouteList[sel1].nextHopIfIndex)]
			if ok {
				v6nextHopInfo[j].NextHopIntRef = intfEntry.name
			}
			v6nextHopInfo[j].Weight = int32(v6protocolRouteList[sel1].weight)
			v6nextHopList = append(v6nextHopList, &v6nextHopInfo[j])
			j++
		}
		v6routes = append(v6routes, &ribd.RouteInfoSummary{
			DestinationNw:   v6destNet,
			IsInstalledInHw: v6isInstalledinHw,
			NextHopList:     v6nextHopList,
		})
	}
	return v6routes
}
func Getv6RoutesPerInterface(intfref string) []string { //*ribd.RouteInfoSummary {
	v6routes := make([]string, 0) //make([]*ribd.RouteInfoSummary, 0)
	routemapInfo := InterfaceRouteMap[intfref]
	if routemapInfo.v6routeMap == nil {
		return v6routes
	}
	for destNetIp, val := range routemapInfo.v6routeMap {
		if val.totalcount == 0 {
			continue
		}
		v6Item := V6RouteInfoMap.Get(patriciaDB.Prefix(destNetIp))
		if v6Item == nil {
			continue
		}
		v6routeInfoRecordList := v6Item.(RouteInfoRecordList)
		v6protocolRouteList, ok := v6routeInfoRecordList.routeInfoProtocolMap[v6routeInfoRecordList.selectedRouteProtocol]
		if !ok || len(v6protocolRouteList) == 0 {
			logger.Info("Unexpected: no route for destNet:", destNetIp, " found in routeMap of type:", intfref)
			continue
		}
		//		v6isInstalledinHw := true
		v6destNet := ""
		//		v6nextHopList := make([]*ribd.NextHopInfo, 0)
		//		v6nextHopInfo := make([]ribd.NextHopInfo, len(v6protocolRouteList))
		//		j := 0
		for sel1 := 0; sel1 < len(v6protocolRouteList); sel1++ {
			v6destNet = v6protocolRouteList[sel1].networkAddr
			/*		v6nextHopInfo[j].NextHopIp = v6protocolRouteList[sel1].nextHopIp.String()
					v6nextHopInfo[j].NextHopIntRef = strconv.Itoa(int(v6protocolRouteList[sel1].nextHopIfIndex))
					intfEntry, ok := IntfIdNameMap[int32(v6protocolRouteList[sel1].nextHopIfIndex)]
					if ok {
						v6nextHopInfo[j].NextHopIntRef = intfEntry.name
					}
					v6nextHopInfo[j].Weight = int32(v6protocolRouteList[sel1].weight)
					v6nextHopList = append(v6nextHopList, &v6nextHopInfo[j])
					j++*/
		}
		/*		v6routes = append(v6routes, &ribd.RouteInfoSummary{
				DestinationNw: v6destNet,
				IsInstalledInHw: v6isInstalledinHw,
				NextHopList:     v6nextHopList,
			})*/
		v6routes = append(v6routes, v6destNet)
	}
	return v6routes
}
func (m RIBDServer) Getv6Route(destNetIp string) (route *ribdInt.IPv6RouteState, err error) {
	var returnRoute ribdInt.IPv6RouteState
	route = &returnRoute
	/*
	   the given address is in CIDR format
	*/
	destNet, err := getNetworkPrefixFromCIDR(destNetIp)
	if err != nil {
		return route, errors.New("Invalid destination ip/network Mask")
	}
	routeInfoRecordListItem := V6RouteInfoMap.Get(destNet)
	if routeInfoRecordListItem == nil {
		logger.Debug("No such route")
		err = errors.New("Route does not exist")
		return route, err
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList) //RouteInfoMap.Get(destNet).(RouteInfoRecordList)
	if routeInfoRecordList.selectedRouteProtocol == "INVALID" {
		logger.Debug("No selected route for this network")
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
			logger.Debug(fmt.Sprintln("Map found for ifndex : ", routeInfoRecord.nextHopIfIndex, "Name = ", intfEntry.name))
			nextHopInfo[i].NextHopIntRef = intfEntry.name
		}
		logger.Debug(fmt.Sprintln("IntfRef = ", nextHopInfo[i].NextHopIntRef))
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

func (m RIBDServer) ProcessV6RouteCreateConfig(cfg *ribd.IPv6Route, addType int, sliceIdx ribd.Int) (val bool, err error) {
	logger.Debug("ProcessV6RouteCreate: Received create route request for ip: ", cfg.DestinationNw, " mask ", cfg.NetworkMask, " number of next hops: ", len(cfg.NextHop), " sliceIdx:", sliceIdx)
	newCfg := ribd.IPv6Route{
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

	//	policyRoute := BuildPolicyRouteFromribdIPv6Route(&newCfg)
	params := BuildRouteParamsFromribdIPv6Route(&newCfg, addType, Invalid, sliceIdx)

	logger.Debug("createType = ", params.createType, "deleteType = ", params.deleteType)
	//	PolicyEngineFilter(policyRoute, policyCommonDefs.PolicyPath_Import, params)
	_, err = createRoute(params)

	return true, err
}

func (m RIBDServer) ProcessV6RouteDeleteConfig(cfg *ribd.IPv6Route, delType int) (val bool, err error) {
	logger.Debug("ProcessRoutev6DeleteConfig:Received Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "number of nextHops:", len(cfg.NextHop), "Protocol ", cfg.Protocol)
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return 0,err
	}
	var nextHopIfIndex ribd.Int
	for i := 0; i < len(cfg.NextHop); i++ {
		nextHopIfIndex = -1
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
		_, err = deleteIPRoute(cfg.DestinationNw, defs.IPv6, cfg.NetworkMask, cfg.Protocol, cfg.NextHop[i].NextHopIp, nextHopIfIndex, ribd.Int(delType), defs.RoutePolicyStateChangetoInValid)
	}
	return true, err
}

func (m RIBDServer) Processv6RoutePatchUpdateConfig(origconfig *ribd.IPv6Route, newconfig *ribd.IPv6Route, op []*ribd.PatchOpInfo) (ret bool, err error) {
	logger.Debug("Processv6RoutePatchUpdateConfig:Received update route request with number of patch ops: ", len(op))
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return err
	}
	destNet, err := getNetowrkPrefixFromStrings(origconfig.DestinationNw, origconfig.NetworkMask)
	if err != nil {
		logger.Debug(" getNetowrkPrefixFromStrings returned err ", err)
		return ret, err
	}
	ok := V6RouteInfoMap.Match(destNet)
	if !ok {
		err = errors.New("No route found")
		return ret, err
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "NextHop":
			logger.Debug("Patch update for next hop")
			/*newconfig should only have the next hops that have to be added or deleted*/
			newconfig.NextHop = make([]*ribd.NextHopInfo, 0)
			logger.Debug("value = ", op[idx].Value)
			valueObjArr := []ribd.NextHopInfo{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				logger.Debug(fmt.Sprintln("error unmarshaling value:", err))
				return ret, errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug("Number of nextHops:", len(valueObjArr))
			for _, val := range valueObjArr {
				logger.Debug("nextHop info: ip - ", val.NextHopIp, " intf: ", val.NextHopIntRef, " wt:", val.Weight)
				//wt,_ := strconv.Atoi((op[idx].Value[j]["Weight"]))
				logger.Debug("IntRef before : ", val.NextHopIntRef)
				if val.NextHopIntRef == "" {
					logger.Info("NextHopIntRef not set")
					nhIntf, err := RouteServiceHandler.GetRouteReachabilityInfo(val.NextHopIp, -1)
					if err != nil {
						logger.Err("next hop ip ", val.NextHopIp, " not reachable")
						return ret, errors.New(fmt.Sprintln("next hop ip ", val.NextHopIp, " not reachable"))
					}
					val.NextHopIntRef = strconv.Itoa(int(nhIntf.NextHopIfIndex))
				} else {
					val.NextHopIntRef, err = m.ConvertIntfStrToIfIndexStr(val.NextHopIntRef)
					if err != nil {
						logger.Err("Invalid NextHop IntRef ", val.NextHopIntRef)
						return ret, err
					}
				}
				logger.Debug("IntRef after : ", val.NextHopIntRef)
				nh := ribd.NextHopInfo{
					NextHopIp:     val.NextHopIp,
					NextHopIntRef: val.NextHopIntRef,
					Weight:        val.Weight,
				}
				newconfig.NextHop = append(newconfig.NextHop, &nh)
			}
			switch op[idx].Op {
			case "add":
				m.ProcessV6RouteCreateConfig(newconfig, FIBAndRIB, ribd.Int(len(destNetSlice)))
			case "remove":
				m.ProcessV6RouteDeleteConfig(newconfig, FIBAndRIB)
			default:
				logger.Err("Operation ", op[idx].Op, " not supported")
			}
		default:
			logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
		}
	}
	return ret, err
}

func (m RIBDServer) Processv6RouteUpdateConfig(origconfig *ribd.IPv6Route, newconfig *ribd.IPv6Route, attrset []bool) (val bool, err error) {
	logger.Debug(fmt.Sprintln("Processv6RouteUpdateConfig:Received update route request origconfig.DestinationNw:", origconfig.DestinationNw, " newconfig.DestinationNw:", newconfig.DestinationNw))
	if !RouteServiceHandler.AcceptConfig {
		logger.Debug("Not ready to accept config")
		//return err
	}
	destNet, err := getNetowrkPrefixFromStrings(origconfig.DestinationNw, origconfig.NetworkMask)
	if err != nil {
		logger.Debug(fmt.Sprintln(" getNetowrkPrefixFromStrings returned err ", err))
		return val, err
	}
	ok := V6RouteInfoMap.Match(destNet)
	if !ok {
		err = errors.New(fmt.Sprintln("No route found for ip ", destNet))
		return val, err
	}
	routeInfoRecordListItem := V6RouteInfoMap.Get(destNet)
	if routeInfoRecordListItem == nil {
		logger.Debug(fmt.Sprintln("No route for destination network", destNet))
		return val, err
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	callUpdate := true
	if attrset != nil {
		found, routeInfoRecord, index := findRouteWithNextHop(routeInfoRecordList.routeInfoProtocolMap[origconfig.Protocol], defs.IPv6, origconfig.NextHop[0].NextHopIp, -1)
		if !found || index == -1 {
			logger.Debug("Invalid nextHopIP")
			return val, err
		}
		objTyp := reflect.TypeOf(*origconfig)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				logger.Debug(fmt.Sprintf("Processv6RouteUpdateConfig (server): changed ", objName))
				if objName == "NextHop" {
					if len(newconfig.NextHop) == 0 {
						logger.Err("Must specify next hop")
						return val, err
					} else {
						nextHopIpAddr, err := getIP(newconfig.NextHop[0].NextHopIp)
						if err != nil {
							logger.Debug("nextHopIpAddr invalid")
							return val, errors.New("Invalid next hop")
						}
						logger.Debug(fmt.Sprintln("Update the next hop info old ip: ", origconfig.NextHop[0].NextHopIp, " new value: ", newconfig.NextHop[0].NextHopIp, " weight : ", newconfig.NextHop[0].Weight))
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
				/*				if objName == "OutgoingInterface" {
								nextHopIfIndex, _ := strconv.Atoi(newconfig.OutgoingInterface)
								routeInfoRecord.nextHopIfIndex = ribd.Int(nextHopIfIndex)
								callUpdate = false
							}*/
			}
		}
		routeInfoRecordList.routeInfoProtocolMap[origconfig.Protocol][index] = routeInfoRecord
		V6RouteInfoMap.Set(destNet, routeInfoRecordList)
		logger.Debug("Adding to DBRouteCh from processRouteUpdateConfig")
		RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
			OrigConfigObject: RouteDBInfo{routeInfoRecord, routeInfoRecordList},
			Op:               defs.Add,
		}
		//RouteServiceHandler.WriteIPv4RouteStateEntryToDB(RouteDBInfo{routeInfoRecord, routeInfoRecordList})
		if callUpdate == false {
			return val, err
		}
	}
	updateBestRoute(destNet, routeInfoRecordList)
	return val, err
}
