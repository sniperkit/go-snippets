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
	"asicd/asicdCommonDefs"
	"asicd/pluginManager/pluginCommon"
	"asicdInt"
	"asicdServices"
	"encoding/json"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils/commonDefs"
	"utils/ipcutils"
	"utils/logging"
)

type AsicdClient struct {
	ClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

type CfgFileJson struct {
	SwitchMac        string            `json:"SwitchMac"`
	PluginList       []string          `json:"PluginList"`
	IfNameMap        map[string]string `json:"IfNameMap"`
	IfNamePrefix     map[string]string `json:"IfNamePrefix"`
	SysRsvdVlanRange string            `json:"SysRsvdVlanRange"`
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type ClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
}

type FSAsicdClientMgr struct {
	ClientHdl *asicdServices.ASICDServicesClient
}

// need to ensure that we are go/thread safe
var asicdmutex *sync.Mutex = &sync.Mutex{}
var Logger *logging.Writer

func (asicdClientMgr *FSAsicdClientMgr) CreateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	asicdmutex.Lock()
	val, err := asicdClientMgr.ClientHdl.CreateIPv4Neighbor(ipAddr, macAddr, vlanId, ifIdx)
	asicdmutex.Unlock()
	return val, err
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	asicdmutex.Lock()
	val, err := asicdClientMgr.ClientHdl.UpdateIPv4Neighbor(ipAddr, macAddr, vlanId, ifIdx)
	asicdmutex.Unlock()
	return val, err
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteIPv4Neighbor(ipAddr string) (int32, error) {
	return asicdClientMgr.ClientHdl.DeleteIPv4Neighbor(ipAddr, "00:00:00:00:00:00", 0, 0)
}

func (asicdClientMgr *FSAsicdClientMgr) CreateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	return asicdClientMgr.ClientHdl.CreateIPv6Neighbor(ipAddr, macAddr, vlanId, ifIdx)
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	return asicdClientMgr.ClientHdl.UpdateIPv6Neighbor(ipAddr, macAddr, vlanId, ifIdx)
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteIPv6Neighbor(ipAddr string) (int32, error) {
	return asicdClientMgr.ClientHdl.DeleteIPv6Neighbor(ipAddr, "00:00:00:00:00:00", 0, 0)
}

func (asicdClientMgr *FSAsicdClientMgr) convertAsicdIP4InfoToCommonInfo(info asicdServices.IPv4IntfState) *commonDefs.IPv4IntfState {
	entry := &commonDefs.IPv4IntfState{}
	entry.IntfRef = info.IntfRef
	entry.IfIndex = info.IfIndex
	entry.IpAddr = info.IpAddr
	entry.OperState = info.OperState
	entry.NumUpEvents = info.NumUpEvents
	entry.LastUpEventTime = info.LastUpEventTime
	entry.NumDownEvents = info.NumDownEvents
	entry.LastDownEventTime = info.LastDownEventTime
	entry.L2IntfType = info.L2IntfType
	entry.L2IntfId = info.L2IntfId
	return entry
}

func (asicdClientMgr *FSAsicdClientMgr) convertAsicdIP6InfoToCommonInfo(info asicdServices.IPv6IntfState) *commonDefs.IPv6IntfState {
	entry := &commonDefs.IPv6IntfState{}
	entry.IntfRef = info.IntfRef
	entry.IfIndex = info.IfIndex
	entry.IpAddr = info.IpAddr
	entry.OperState = info.OperState
	entry.NumUpEvents = info.NumUpEvents
	entry.LastUpEventTime = info.LastUpEventTime
	entry.NumDownEvents = info.NumDownEvents
	entry.LastDownEventTime = info.LastDownEventTime
	entry.L2IntfType = info.L2IntfType
	entry.L2IntfId = info.L2IntfId
	return entry
}

func (asicdClientMgr *FSAsicdClientMgr) convertAsicdPortStateInfoToCommonInfo(info asicdServices.PortState) *commonDefs.PortState {
	entry := &commonDefs.PortState{}
	entry.IntfRef = info.IntfRef
	entry.IfIndex = info.IfIndex
	entry.Name = info.Name
	entry.OperState = info.OperState
	entry.NumUpEvents = info.NumUpEvents
	entry.LastUpEventTime = info.LastUpEventTime
	entry.NumDownEvents = info.NumDownEvents
	entry.LastDownEventTime = info.LastDownEventTime
	entry.Pvid = info.Pvid
	entry.IfInOctets = info.IfInOctets
	entry.IfInUcastPkts = info.IfInUcastPkts
	entry.IfInDiscards = info.IfInDiscards
	entry.IfInErrors = info.IfInErrors
	entry.IfInUnknownProtos = info.IfInUnknownProtos
	entry.IfOutOctets = info.IfOutOctets
	entry.IfOutUcastPkts = info.IfOutUcastPkts
	entry.IfOutDiscards = info.IfOutDiscards
	entry.IfOutErrors = info.IfOutErrors
	entry.ErrDisableReason = info.ErrDisableReason
	return entry
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkIPv4IntfState(curMark, count int) (*commonDefs.IPv4IntfStateGetInfo, error) {
	asicdmutex.Lock()
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkIPv4IntfState(asicdServices.Int(curMark), asicdServices.Int(count))
	asicdmutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var ipv4Info commonDefs.IPv4IntfStateGetInfo
	ipv4Info.StartIdx = int32(bulkInfo.StartIdx)
	ipv4Info.EndIdx = int32(bulkInfo.EndIdx)
	ipv4Info.Count = int32(bulkInfo.Count)
	ipv4Info.More = bulkInfo.More
	ipv4Info.IPv4IntfStateList = make([]commonDefs.IPv4IntfState, int(ipv4Info.Count))
	for idx := 0; idx < int(ipv4Info.Count); idx++ {
		ipv4Info.IPv4IntfStateList[idx].IntfRef = bulkInfo.IPv4IntfStateList[idx].IntfRef
		ipv4Info.IPv4IntfStateList[idx].IfIndex = bulkInfo.IPv4IntfStateList[idx].IfIndex
		ipv4Info.IPv4IntfStateList[idx].IpAddr = bulkInfo.IPv4IntfStateList[idx].IpAddr
		ipv4Info.IPv4IntfStateList[idx].OperState = bulkInfo.IPv4IntfStateList[idx].OperState
		ipv4Info.IPv4IntfStateList[idx].NumUpEvents = bulkInfo.IPv4IntfStateList[idx].NumUpEvents
		ipv4Info.IPv4IntfStateList[idx].LastUpEventTime = bulkInfo.IPv4IntfStateList[idx].LastUpEventTime
		ipv4Info.IPv4IntfStateList[idx].NumDownEvents = bulkInfo.IPv4IntfStateList[idx].NumDownEvents
		ipv4Info.IPv4IntfStateList[idx].LastDownEventTime = bulkInfo.IPv4IntfStateList[idx].LastDownEventTime
		ipv4Info.IPv4IntfStateList[idx].L2IntfType = bulkInfo.IPv4IntfStateList[idx].L2IntfType
		ipv4Info.IPv4IntfStateList[idx].L2IntfId = bulkInfo.IPv4IntfStateList[idx].L2IntfId
	}
	return &ipv4Info, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkPort(curMark, count int) (*commonDefs.PortGetInfo, error) {
	asicdmutex.Lock()
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkPort(asicdServices.Int(curMark), asicdServices.Int(count))
	asicdmutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var portInfo commonDefs.PortGetInfo
	portInfo.StartIdx = int32(bulkInfo.StartIdx)
	portInfo.EndIdx = int32(bulkInfo.EndIdx)
	portInfo.Count = int32(bulkInfo.Count)
	portInfo.More = bulkInfo.More
	portInfo.PortList = make([]commonDefs.Port, int(portInfo.Count))
	for idx := 0; idx < int(portInfo.Count); idx++ {
		portInfo.PortList[idx].IntfRef = bulkInfo.PortList[idx].IntfRef
		portInfo.PortList[idx].IfIndex = bulkInfo.PortList[idx].IfIndex
		portInfo.PortList[idx].Description = bulkInfo.PortList[idx].Description
		portInfo.PortList[idx].PhyIntfType = bulkInfo.PortList[idx].PhyIntfType
		portInfo.PortList[idx].AdminState = bulkInfo.PortList[idx].AdminState
		portInfo.PortList[idx].MacAddr = bulkInfo.PortList[idx].MacAddr
		portInfo.PortList[idx].Speed = bulkInfo.PortList[idx].Speed
		portInfo.PortList[idx].Duplex = bulkInfo.PortList[idx].Duplex
		portInfo.PortList[idx].Autoneg = bulkInfo.PortList[idx].Autoneg
		portInfo.PortList[idx].MediaType = bulkInfo.PortList[idx].MediaType
		portInfo.PortList[idx].Mtu = bulkInfo.PortList[idx].Mtu
	}
	return &portInfo, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkPortState(curMark, count int) (*commonDefs.PortStateGetInfo, error) {
	asicdmutex.Lock()
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkPortState(asicdServices.Int(curMark), asicdServices.Int(count))
	asicdmutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var portStateInfo commonDefs.PortStateGetInfo
	portStateInfo.StartIdx = int32(bulkInfo.StartIdx)
	portStateInfo.EndIdx = int32(bulkInfo.EndIdx)
	portStateInfo.Count = int32(bulkInfo.Count)
	portStateInfo.More = bulkInfo.More
	portStateInfo.PortStateList = make([]commonDefs.PortState, int(portStateInfo.Count))
	for idx := 0; idx < int(portStateInfo.Count); idx++ {
		portStateInfo.PortStateList[idx].IntfRef = bulkInfo.PortStateList[idx].IntfRef
		portStateInfo.PortStateList[idx].IfIndex = bulkInfo.PortStateList[idx].IfIndex
		portStateInfo.PortStateList[idx].Name = bulkInfo.PortStateList[idx].Name
		portStateInfo.PortStateList[idx].OperState = bulkInfo.PortStateList[idx].OperState
		portStateInfo.PortStateList[idx].NumUpEvents = bulkInfo.PortStateList[idx].NumUpEvents
		portStateInfo.PortStateList[idx].LastUpEventTime = bulkInfo.PortStateList[idx].LastUpEventTime
		portStateInfo.PortStateList[idx].NumDownEvents = bulkInfo.PortStateList[idx].NumDownEvents
		portStateInfo.PortStateList[idx].LastDownEventTime = bulkInfo.PortStateList[idx].LastDownEventTime
		portStateInfo.PortStateList[idx].Pvid = bulkInfo.PortStateList[idx].Pvid
		portStateInfo.PortStateList[idx].IfInOctets = bulkInfo.PortStateList[idx].IfInOctets
		portStateInfo.PortStateList[idx].IfInUcastPkts = bulkInfo.PortStateList[idx].IfInUcastPkts
		portStateInfo.PortStateList[idx].IfInDiscards = bulkInfo.PortStateList[idx].IfInDiscards
		portStateInfo.PortStateList[idx].IfInErrors = bulkInfo.PortStateList[idx].IfInErrors
		portStateInfo.PortStateList[idx].IfInUnknownProtos = bulkInfo.PortStateList[idx].IfInUnknownProtos
		portStateInfo.PortStateList[idx].IfOutOctets = bulkInfo.PortStateList[idx].IfOutOctets
		portStateInfo.PortStateList[idx].IfOutUcastPkts = bulkInfo.PortStateList[idx].IfOutUcastPkts
		portStateInfo.PortStateList[idx].IfOutDiscards = bulkInfo.PortStateList[idx].IfOutDiscards
		portStateInfo.PortStateList[idx].IfOutErrors = bulkInfo.PortStateList[idx].IfOutErrors
		portStateInfo.PortStateList[idx].ErrDisableReason = bulkInfo.PortStateList[idx].ErrDisableReason
	}
	return &portStateInfo, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkVlanState(curMark, count int) (*commonDefs.VlanStateGetInfo, error) {

	asicdmutex.Lock()
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
	asicdmutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var vlanStateInfo commonDefs.VlanStateGetInfo
	vlanStateInfo.StartIdx = int32(bulkInfo.StartIdx)
	vlanStateInfo.EndIdx = int32(bulkInfo.EndIdx)
	vlanStateInfo.Count = int32(bulkInfo.Count)
	vlanStateInfo.More = bulkInfo.More
	vlanStateInfo.VlanStateList = make([]commonDefs.VlanState, int(vlanStateInfo.Count))
	for idx := 0; idx < int(vlanStateInfo.Count); idx++ {
		vlanStateInfo.VlanStateList[idx].VlanId = bulkInfo.VlanStateList[idx].VlanId
		vlanStateInfo.VlanStateList[idx].VlanName = bulkInfo.VlanStateList[idx].VlanName
		vlanStateInfo.VlanStateList[idx].OperState = bulkInfo.VlanStateList[idx].OperState
		vlanStateInfo.VlanStateList[idx].IfIndex = bulkInfo.VlanStateList[idx].IfIndex
	}

	return &vlanStateInfo, nil
}

func convertAsicdVlanStateInfoToCommonInfo(info asicdServices.VlanState) *commonDefs.VlanState {
	entry := &commonDefs.VlanState{
		VlanId:    info.VlanId,
		VlanName:  info.VlanName,
		OperState: info.OperState,
		IfIndex:   info.IfIndex,
	}
	return entry
}

func convertAsicdVlanInfoToCommonInfo(info asicdInt.Vlan) *commonDefs.Vlan {
	entry := &commonDefs.Vlan{}
	entry.VlanId = info.VlanId
	entry.IfIndexList = append(entry.IfIndexList, info.IfIndexList...)
	entry.UntagIfIndexList = append(entry.UntagIfIndexList, info.UntagIfIndexList...)
	return entry
}

func (asicdClientMgr *FSAsicdClientMgr) GetAllVlanState() ([]*commonDefs.VlanState, error) {
	curMark := 0
	count := 100
	vlanStateInfo := make([]*commonDefs.VlanState, 0)
	for {
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			vlanStateInfo = append(vlanStateInfo, convertAsicdVlanStateInfoToCommonInfo(*bulkInfo.VlanStateList[idx]))
		}
		if bulkInfo.More == false {
			break
		}
	}
	return vlanStateInfo, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetAllVlan() ([]*commonDefs.Vlan, error) {
	curMark := 0
	count := 100
	vlanInfo := make([]*commonDefs.Vlan, 0)
	for {
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			vlanInfo = append(vlanInfo, convertAsicdVlanInfoToCommonInfo(*bulkInfo.VlanList[idx]))
		}
		if bulkInfo.More == false {
			break
		}
	}
	return vlanInfo, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error) {
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
	if bulkInfo == nil {
		return nil, err
	}
	var vlanInfo commonDefs.VlanGetInfo
	vlanInfo.StartIdx = int32(bulkInfo.StartIdx)
	vlanInfo.EndIdx = int32(bulkInfo.EndIdx)
	vlanInfo.Count = int32(bulkInfo.Count)
	vlanInfo.More = bulkInfo.More
	vlanInfo.VlanList = make([]commonDefs.Vlan, int(vlanInfo.Count))
	for idx := 0; idx < int(vlanInfo.Count); idx++ {
		vlanInfo.VlanList[idx].VlanId = bulkInfo.VlanList[idx].VlanId
		vlanInfo.VlanList[idx].IfIndexList = append(vlanInfo.VlanList[idx].IfIndexList, bulkInfo.VlanList[idx].IfIndexList...)
		vlanInfo.VlanList[idx].UntagIfIndexList = append(vlanInfo.VlanList[idx].UntagIfIndexList, bulkInfo.VlanList[idx].UntagIfIndexList...)
	}
	return &vlanInfo, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetBulkLag(curMark, count int) (*commonDefs.LagGetInfo, error) {
	bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkLag(asicdInt.Int(curMark), asicdInt.Int(count))
	if bulkInfo == nil {
		return nil, err
	}
	var lagInfo commonDefs.LagGetInfo
	lagInfo.StartIdx = int32(bulkInfo.StartIdx)
	lagInfo.EndIdx = int32(bulkInfo.EndIdx)
	lagInfo.Count = int32(bulkInfo.Count)
	lagInfo.More = bulkInfo.More
	lagInfo.LagList = make([]commonDefs.Lag, int(lagInfo.Count))
	for idx := 0; idx < int(lagInfo.Count); idx++ {
		lagInfo.LagList[idx].LagIfIndex = bulkInfo.LagList[idx].LagIfIndex
		lagInfo.LagList[idx].HashType = bulkInfo.LagList[idx].HashType
		lagInfo.LagList[idx].LagName = bulkInfo.LagList[idx].LagName
		lagInfo.LagList[idx].IfIndexList = append(lagInfo.LagList[idx].IfIndexList, bulkInfo.LagList[idx].IfIndexList...)
	}
	return &lagInfo, nil
}

func GetAsicdThriftClientHdl(paramsFile string, logger *logging.Writer) *asicdServices.ASICDServicesClient {
	var asicdClient AsicdClient
	Logger = logger
	logger.Debug(fmt.Sprintln("Inside connectToServers...paramsFile", paramsFile))
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		logger.Err("Error in reading configuration file")
		return nil
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		logger.Err("Error in Unmarshalling Json")
		return nil
	}

	for _, client := range clientsList {
		if client.Name == "asicd" {
			logger.Debug(fmt.Sprintln("found asicd at port", client.Port))
			asicdClient.Address = "localhost:" + strconv.Itoa(client.Port)
			asicdClient.Transport, asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(asicdClient.Address)
			if err != nil {
				logger.Err(fmt.Sprintln("Failed to connect to Asicd, retrying until connection is successful"))
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					asicdClient.Transport, asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(asicdClient.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						logger.Err("Still can't connect to Asicd, retrying..")
					}
				}

			}
			logger.Info("Connected to Asicd")
			asicdClient.ClientHdl = asicdServices.NewASICDServicesClientFactory(asicdClient.Transport, asicdClient.PtrProtocolFactory)
			return asicdClient.ClientHdl
		}
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetAllPortState() ([]*commonDefs.PortState, error) {
	curMark := int(asicdCommonDefs.MIN_SYS_PORTS)
	count := 100
	portState := make([]*commonDefs.PortState, 0)
	for {
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkPortState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			portState = append(portState,
				asicdClientMgr.convertAsicdPortStateInfoToCommonInfo(*bulkInfo.PortStateList[idx]))
		}
		if bulkInfo.More == false {
			break
		}
	}
	return portState, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetPort(intfRef string) (*commonDefs.Port, error) {
	portInfo, err := asicdClientMgr.ClientHdl.GetPort(intfRef)
	if err != nil {
		return nil, err
	}
	port := &commonDefs.Port{
		IntfRef:     portInfo.IntfRef,
		IfIndex:     portInfo.IfIndex,
		Description: portInfo.Description,
		PhyIntfType: portInfo.PhyIntfType,
		AdminState:  portInfo.AdminState,
		MacAddr:     portInfo.MacAddr,
		Speed:       portInfo.Speed,
		Duplex:      portInfo.Duplex,
		Autoneg:     portInfo.Autoneg,
		MediaType:   portInfo.MediaType,
		Mtu:         portInfo.Mtu,
	}
	return port, nil
}

/*  API to return all ipv4 addresses created on the system... If a dameons uses this then they do not have to worry
 *  about checking is any ipv4 addresses are left on the system or not
 */
func (asicdClientMgr *FSAsicdClientMgr) GetAllIPv6IntfState() ([]*commonDefs.IPv6IntfState, error) {
	curMark := 0
	count := 100
	ipv6Info := make([]*commonDefs.IPv6IntfState, 0)
	for {
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkIPv6IntfState(asicdServices.Int(curMark),
			asicdServices.Int(count))
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			ipv6Info = append(ipv6Info,
				asicdClientMgr.convertAsicdIP6InfoToCommonInfo(*bulkInfo.IPv6IntfStateList[idx]))
		}
		if bulkInfo.More == false {
			break
		}
	}

	return ipv6Info, nil
}

/*  API to return all ipv4 addresses created on the system... If a dameons uses this then they do not have to worry
 *  about checking is any ipv4 addresses are left on the system or not
 */
func (asicdClientMgr *FSAsicdClientMgr) GetAllIPv4IntfState() ([]*commonDefs.IPv4IntfState, error) {
	curMark := 0
	count := 100
	ipv4Info := make([]*commonDefs.IPv4IntfState, 0)
	for {
		asicdmutex.Lock()
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkIPv4IntfState(asicdServices.Int(curMark),
			asicdServices.Int(count))
		asicdmutex.Unlock()
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			ipv4Info = append(ipv4Info,
				asicdClientMgr.convertAsicdIP4InfoToCommonInfo(*bulkInfo.IPv4IntfStateList[idx]))
		}
		if bulkInfo.More == false {
			break
		}
	}
	return ipv4Info, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetAllSubIPv4IntfState() ([]*commonDefs.SubIPv4IntfState, error) {
	curMark := 0
	count := 100
	ipv4Info := make([]*commonDefs.SubIPv4IntfState, 0)
	for {
		asicdmutex.Lock()
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkSubIPv4IntfState(asicdServices.Int(curMark),
			asicdServices.Int(count))
		asicdmutex.Unlock()
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			obj := bulkInfo.SubIPv4IntfStateList[idx]
			entry := &commonDefs.SubIPv4IntfState{
				IntfRef:       obj.IntfRef,
				Type:          obj.Type,
				IfIndex:       obj.IfIndex,
				IfName:        obj.IfName,
				ParentIfIndex: obj.ParentIfIndex,
				IpAddr:        obj.IpAddr,
				MacAddr:       obj.MacAddr,
				OperState:     obj.OperState,
			}
			ipv4Info = append(ipv4Info, entry)
		}
		if bulkInfo.More == false {
			break
		}
	}
	return ipv4Info, nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetAllSubIPv6IntfState() ([]*commonDefs.SubIPv6IntfState, error) {
	curMark := 0
	count := 100
	ipv6Info := make([]*commonDefs.SubIPv6IntfState, 0)
	for {
		asicdmutex.Lock()
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkSubIPv6IntfState(asicdServices.Int(curMark),
			asicdServices.Int(count))
		asicdmutex.Unlock()
		if bulkInfo == nil {
			return nil, err
		}
		curMark = int(bulkInfo.EndIdx)
		for idx := 0; idx < int(bulkInfo.Count); idx++ {
			obj := bulkInfo.SubIPv6IntfStateList[idx]
			entry := &commonDefs.SubIPv6IntfState{
				IntfRef:       obj.IntfRef,
				Type:          obj.Type,
				IfIndex:       obj.IfIndex,
				IfName:        obj.IfName,
				ParentIfIndex: obj.ParentIfIndex,
				IpAddr:        obj.IpAddr,
				MacAddr:       obj.MacAddr,
				OperState:     obj.OperState,
			}
			ipv6Info = append(ipv6Info, entry)
		}
		if bulkInfo.More == false {
			break
		}
	}
	return ipv6Info, nil
	return make([]*commonDefs.SubIPv6IntfState, 0), nil
}

/*  Library util to determine router id.
 *  Calculation Method:
 *	    1) Get all loopback interfaces on the system and return the highest value
 *		a) If no loopback configured on the system, in that case get all ipv4 interfaces and return the highest
 *		   value
 *		    b) if no ipv4 interfaces then return default router id which is 0.0.0.0
 */
func (asicdClientMgr *FSAsicdClientMgr) DetermineRouterId() string {
	rtrId := "0.0.0.0"
	asicdmutex.Lock()
	allipv4Intfs, err := asicdClientMgr.GetAllIPv4IntfState()
	asicdmutex.Unlock()
	if err != nil {
		return rtrId
	}
	loopbackIntfs := make([]string, 0)
	ipv4Intfs := make([]string, 0)
	// Get loopback interfaces & ipv4 interfaces
	for _, ipv4Intf := range allipv4Intfs {
		switch asicdCommonDefs.GetIntfTypeFromIfIndex(ipv4Intf.IfIndex) {
		case commonDefs.IfTypeLoopback:
			loopbackIntfs = append(loopbackIntfs, ipv4Intf.IpAddr)

		case commonDefs.IfTypeVlan, commonDefs.IfTypePort:
			ipv4Intfs = append(ipv4Intfs, ipv4Intf.IpAddr)
		}
	}

	for _, ipAddr := range loopbackIntfs {
		if strings.Compare(ipAddr, rtrId) > 0 {
			// current loopback Ip Addr is greater than rtrId... time to update router id
			rtrId = ipAddr
		}
	}

	if rtrId != "0.0.0.0" {
		// there was a loopback on the system which is higher then default rtrId and we are going to use that
		// ipAddr as router id
		return rtrId
	}

	for _, ipAddr := range ipv4Intfs {
		if strings.Compare(ipAddr, rtrId) > 0 {
			// current ipv4 ip addr is greater than rtrId... time to update router id
			rtrId = ipAddr
		}
	}
	return rtrId
}

// convert the lacp port names name to asic format string list
func asicDPortBmpFormatGet(distPortList []string) string {
	s := ""
	dLength := len(distPortList)

	for i := 0; i < dLength; i++ {
		num := strings.Split(distPortList[i], "-")[1]
		if i == dLength-1 {
			s += num
		} else {
			s += num + ","
		}
	}
	return s

}

func (asicdClientMgr *FSAsicdClientMgr) GetPortLinkStatus(pId int32) bool {

	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		bulkInfo, err := asicdClientMgr.ClientHdl.GetBulkPortState(asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS), asicdServices.Int(asicdCommonDefs.MAX_SYS_PORTS))
		asicdmutex.Unlock()
		if err == nil && bulkInfo.Count != 0 {
			objCount := int64(bulkInfo.Count)
			for i := int64(0); i < objCount; i++ {
				if bulkInfo.PortStateList[i].IfIndex == pId {
					return bulkInfo.PortStateList[i].OperState == pluginCommon.UpDownState[1]
				}
			}
		}
		Logger.Info(fmt.Sprintf("asicDGetPortLinkSatus: could not get status for port %d, failure in get method\n", pId))
	}
	return true

}

func (asicdClientMgr *FSAsicdClientMgr) CreateStgBridge(vlanList []uint16) int32 {

	vl := make([]int32, 0)
	if asicdClientMgr.ClientHdl != nil {
		for _, v := range vlanList {
			vl = append(vl, int32(v))
		}
		//asicdmutex.Lock()
		// default vlan is already created in opennsl
		stgid, err := asicdClientMgr.ClientHdl.CreateStg(vl)
		//asicdmutex.Unlock()
		if err == nil {
			for _, v := range vl {
				if v != 0 &&
					v != 4095 {
					protocolmac := asicdInt.RsvdProtocolMacConfig{
						MacAddr:     "01:00:0C:CC:CC:CD",
						MacAddrMask: "FF:FF:FF:FF:FF:FF",
						VlanId:      int32(v),
					}
					asicdmutex.Lock()
					asicdClientMgr.ClientHdl.EnablePacketReception(&protocolmac)
					asicdmutex.Unlock()
				}
			}
			return stgid
		} else {
			Logger.Info(fmt.Sprintf("Create Stg Group error %#v", err))
		}
	} else {
		Logger.Info(fmt.Sprintf("Create Stg Group failed asicd not connected"))
	}
	return -1
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteStgBridge(stgid int32, vlanList []uint16) error {
	vl := make([]int32, 0)

	if asicdClientMgr.ClientHdl != nil {

		for _, v := range vlanList {
			vl = append(vl, int32(v))
		}
		for _, v := range vl {
			if v != 0 &&
				v != 4095 {
				protocolmac := asicdInt.RsvdProtocolMacConfig{
					MacAddr:     "01:00:0C:CC:CC:CD",
					MacAddrMask: "FF:FF:FF:FF:FF:FF",
					VlanId:      int32(v),
				}

				Logger.Info(fmt.Sprintf("Deleting PVST MAC entry %#v", protocolmac))
				asicdmutex.Lock()
				asicdClientMgr.ClientHdl.DisablePacketReception(&protocolmac)
				asicdmutex.Unlock()
			}
		}
		Logger.Info(fmt.Sprintf("Deleting Stg Group %d with vlans %#v", stgid, vl))

		//asicdmutex.Lock()
		_, err := asicdClientMgr.ClientHdl.DeleteStg(stgid)
		//asicdmutex.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) SetStgPortState(stgid int32, ifindex int32, state int) error {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		_, err := asicdClientMgr.ClientHdl.SetPortStpState(stgid, ifindex, int32(state))
		asicdmutex.Unlock()
		return err
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) FlushStgFdb(stgid, ifindex int32) error {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		err := asicdClientMgr.ClientHdl.FlushFdbStgGroup(stgid, ifindex)
		asicdmutex.Unlock()
		return err
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) BPDUGuardDetected(ifindex int32, enable bool) error {
	if asicdClientMgr.ClientHdl != nil {
		state := "DOWN"
		if enable {
			state = "UP"
		}
		asicdmutex.Lock()
		_, err := asicdClientMgr.ClientHdl.ErrorDisablePort(ifindex, state, "STP BPDU GUARD")
		asicdmutex.Unlock()
		return err
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) GetSwitchMAC(paramsPath string) string {
	var cfgFile CfgFileJson

	asicdconffilename := paramsPath + pluginCommon.ASICD_CONFIG_FILE

	cfgFileData, err := ioutil.ReadFile(asicdconffilename)
	if err != nil {
		Logger.Err("Error reading config file -", pluginCommon.ASICD_CONFIG_FILE,
			". Using defaults (linux plugin only)")
		return "00:00:00:00:00:00"
	}
	err = json.Unmarshal(cfgFileData, &cfgFile)
	if err != nil {
		Logger.Err("Error parsing config file, using defaults (linux plugin only)")
		return "00:00:00:00:00:00"
	}

	return cfgFile.SwitchMac
}

func (asicdClientMgr *FSAsicdClientMgr) CreateLag(ifName string, hashType int32, ports string) (ifindex int32, err error) {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		ifindex, err = asicdClientMgr.ClientHdl.CreateLag(ifName, hashType, ports)
		asicdmutex.Unlock()
		return ifindex, err
	}
	return -1, err
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteLag(ifIndex int32) (err error) {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		_, err = asicdClientMgr.ClientHdl.DeleteLag(ifIndex)
		asicdmutex.Unlock()
	}
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateLag(ifIndex, hashType int32, ports string) (err error) {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		_, err = asicdClientMgr.ClientHdl.UpdateLag(ifIndex, hashType, ports)
		asicdmutex.Unlock()
	}
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateLagCfgIntfList(ifName string, ifIndexList []int32) bool {
	var ok bool
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		ok, _ = asicdClientMgr.ClientHdl.UpdateLagCfgIntfList(ifName, ifIndexList)
		asicdmutex.Unlock()
	}
	return ok
}

func (asicdClientMgr *FSAsicdClientMgr) EnablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		cfg := &asicdInt.RsvdProtocolMacConfig{
			MacAddr:     mac,
			MacAddrMask: "FF:FF:FF:FF:FF:FF",
		}
		_, err = asicdClientMgr.ClientHdl.EnablePacketReception(cfg)
		asicdmutex.Unlock()
	}
	return err

}

func (asicdClientMgr *FSAsicdClientMgr) DisablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		cfg := &asicdInt.RsvdProtocolMacConfig{
			MacAddr:     mac,
			MacAddrMask: "FF:FF:FF:FF:FF:FF",
		}
		_, err = asicdClientMgr.ClientHdl.DisablePacketReception(cfg)
		asicdmutex.Unlock()
	}
	return err

}

// TODO change this to pass in the Intf
func (asicdClientMgr *FSAsicdClientMgr) IppIngressEgressDrop(srcIfIndex, dstIfIndex string) (err error) {

	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		aclName := fmt.Sprintf("IPPInOutBlockfpPort%s", dstIfIndex)
		ruleName := fmt.Sprintf("%sfpPort%s", aclName, srcIfIndex)
		filter := &asicdServices.AclIpv4Filter{
			FilterName: "blockIngressEgress",
			SrcIntf:    srcIfIndex,
			DstIntf:    dstIfIndex,
		}
		rule := &asicdServices.Acl{
			AclName:    ruleName,
			Action:     "DENY",
			Priority:   10,
			Stage:      "OUT",
			FilterName: "blockIngressEgress",
		}
		_, err = asicdClientMgr.ClientHdl.CreateAclIpv4Filter(filter)
		if err != nil {
			asicdmutex.Unlock()
			return err
		}
		_, err = asicdClientMgr.ClientHdl.CreateAcl(rule)
		if err != nil {
			asicdmutex.Unlock()
			return err
		}

		asicdmutex.Unlock()
	}

	return err
}

func (asicdClientMgr *FSAsicdClientMgr) IppIngressEgressPass(srcIfIndex, dstIfIndex string) (err error) {

	if asicdClientMgr.ClientHdl != nil {
		asicdmutex.Lock()
		aclName := fmt.Sprintf("IPPInOutBlockfpPort%s", dstIfIndex)
		ruleName := fmt.Sprintf("%sfpPort%s", aclName, srcIfIndex)
		filterName := ruleName + "Filter"
		filter := &asicdServices.AclIpv4Filter{
			FilterName: filterName,
			SrcIntf:    srcIfIndex,
			DstIntf:    dstIfIndex,
		}
		rule := &asicdServices.Acl{
			AclName:    ruleName,
			Stage:      "OUT",
			Priority:   10,
			FilterName: filterName,
		}
		_, err = asicdClientMgr.ClientHdl.CreateAclIpv4Filter(filter)
		if err != nil {
			asicdmutex.Unlock()
			return err
		}
		_, err = asicdClientMgr.ClientHdl.CreateAcl(rule)
		if err != nil {
			asicdmutex.Unlock()
			return err
		}

		asicdmutex.Unlock()
	}

	return err
}

func (asicdClientMgr *FSAsicdClientMgr) IppVlanConversationSet(vlan uint16, ifindex int32) error {

	// get the vlan info
	vlanInfo, err := asicdClientMgr.ClientHdl.GetBulkVlan(asicdInt.Int(vlan), asicdInt.Int(vlan))
	if err != nil {
		return err
	}

	// append the port to the vlan
	for i := vlanInfo.StartIdx; i < (vlanInfo.StartIdx + vlanInfo.Count); i++ {
		vlancfg := vlanInfo.VlanList[i]
		for _, ifndx := range vlancfg.IfIndexList {
			if ifndx == ifindex {
				return nil
			}
		}
		oldvlancfg := &asicdServices.Vlan{
			VlanId: int32(vlan),
		}
		newvlancfg := &asicdServices.Vlan{
			VlanId: int32(vlan),
		}
		patchList := []*asicdServices.PatchOpInfo{&asicdServices.PatchOpInfo{
			Op:    "add",
			Path:  "IntfList",
			Value: fmt.Sprintf("%s", ifindex),
		}}
		_, err = asicdClientMgr.ClientHdl.UpdateVlan(oldvlancfg, newvlancfg, nil, patchList)
		//(1: Vlan origconfig, 2: Vlan newconfig, 3: list<bool> attrset, 4: list<PatchOpInfo> op)
		if err != nil {
			return err
		}
	}
	return nil
}

func (asicdClientMgr *FSAsicdClientMgr) IppVlanConversationClear(vlan uint16, ifindex int32) (err error) {
	// get the vlan info
	vlanInfo, err := asicdClientMgr.ClientHdl.GetBulkVlan(asicdInt.Int(vlan), asicdInt.Int(vlan))
	if err != nil {
		return err
	}

	// append the port to the vlan
	for i := vlanInfo.StartIdx; i < (vlanInfo.StartIdx + vlanInfo.Count); i++ {
		vlancfg := vlanInfo.VlanList[i]
		for _, ifndx := range vlancfg.IfIndexList {
			if ifndx == ifindex {
				oldvlancfg := &asicdServices.Vlan{
					VlanId: int32(vlan),
				}
				newvlancfg := &asicdServices.Vlan{
					VlanId: int32(vlan),
				}
				patchList := []*asicdServices.PatchOpInfo{&asicdServices.PatchOpInfo{
					Op:    "remove",
					Path:  "IntfList",
					Value: fmt.Sprintf("%s", ifindex),
				}}
				_, err = asicdClientMgr.ClientHdl.UpdateVlan(oldvlancfg, newvlancfg, nil, patchList)
				break
			}
		}
	}
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) IsLoopbackType(ifIndex int32) bool {
	if pluginCommon.GetTypeFromIfIndex(ifIndex) == commonDefs.IfTypeLoopback {
		return true
	}
	return false
}

func createVirtualV4Obj(intRef, ipAddr, macAddr string, enable bool) *asicdServices.SubIPv4Intf {
	obj := asicdServices.NewSubIPv4Intf()
	obj.IntfRef = intRef
	obj.Type = pluginCommon.SUB_INTF_VIRTUAL_TYPE
	obj.IpAddr = ipAddr
	obj.MacAddr = macAddr
	obj.Enable = enable

	return obj
}

func createVirtualV6Obj(intRef, ipAddr, macAddr string, enable bool) *asicdServices.SubIPv6Intf {
	obj := asicdServices.NewSubIPv6Intf()
	obj.IntfRef = intRef
	obj.Type = pluginCommon.SUB_INTF_VIRTUAL_TYPE
	obj.IpAddr = ipAddr
	obj.MacAddr = macAddr
	obj.Enable = enable

	return obj
}

func (asicdClientMgr *FSAsicdClientMgr) CreateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.CreateSubIPv4Intf(createVirtualV4Obj(intRef, ipAddr, macAddr, enable))
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) CreateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.CreateSubIPv6Intf(createVirtualV6Obj(intRef, ipAddr, macAddr, enable))
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.UpdateVirtualIPv4Intf(intRef, pluginCommon.SUB_INTF_VIRTUAL_TYPE, ipAddr, macAddr, enable)
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) UpdateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.UpdateVirtualIPv6Intf(intRef, pluginCommon.SUB_INTF_VIRTUAL_TYPE, ipAddr, macAddr, enable)
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.DeleteSubIPv4Intf(createVirtualV4Obj(intRef, ipAddr, macAddr, enable))
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) DeleteVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	_, err = asicdClientMgr.ClientHdl.DeleteSubIPv6Intf(createVirtualV6Obj(intRef, ipAddr, macAddr, enable))
	return err
}

func (asicdClientMgr *FSAsicdClientMgr) IsLinuxOnlyPlugin() (bool, error) {
	asicdmutex.Lock()
	ret, err := asicdClientMgr.ClientHdl.IsLinuxOnlyPlugin()
	asicdmutex.Unlock()
	return ret, err
}
