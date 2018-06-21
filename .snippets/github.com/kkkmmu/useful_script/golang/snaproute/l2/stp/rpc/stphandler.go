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

// lahandler
package rpc

import (
	"fmt"
	stp "l2/stp/protocol"
	"l2/stp/server"
	"models/objects"
	"reflect"
	"stpd"
	"utils/dbutils"
	//"time"
	"errors"
)

const DBName string = "UsrConfDb.db"

type STPDServiceHandler struct {
	server *server.STPServer
}

func NewSTPDServiceHandler(svr *server.STPServer) *STPDServiceHandler {
	svchdl := &STPDServiceHandler{
		server: svr,
	}
	// lets read the db and replay the data
	svchdl.ReadConfigFromDB(stp.StpGlobalStateGet())
	return svchdl
}

//
func ConvertThriftBrgConfigToStpBrgConfig(config *stpd.StpBridgeInstance, brgconfig *stp.StpBridgeConfig) {

	brgconfig.Address = config.Address
	brgconfig.Priority = uint16(config.Priority)
	brgconfig.Vlan = uint16(config.Vlan)
	brgconfig.MaxAge = uint16(config.MaxAge)
	brgconfig.HelloTime = uint16(config.HelloTime)
	brgconfig.ForwardDelay = uint16(config.ForwardDelay)
	brgconfig.ForceVersion = int32(config.ForceVersion)
	brgconfig.TxHoldCount = int32(config.TxHoldCount)
}

// converts yang true(1)/false(2) to bool
func ConvertInt32ToBool(val int32) bool {
	if val == 2 {
		return false
	}
	return true
}

// converts  bool to yang true(1)/false(2)
func ConvertBoolToInt32(val bool) int32 {
	if val {
		return 1
	}
	return 2
}

func ConvertThriftPortConfigToStpPortConfig(config *stpd.StpPort, portconfig *stp.StpPortConfig) {

	portconfig.IfIndex = stp.GetIfIndexFromIntfRef(config.IntfRef)
	portconfig.BrgIfIndex = int32(config.Vlan)
	portconfig.Priority = uint16(config.Priority)
	if config.AdminState == "UP" {
		portconfig.Enable = true
	} else {
		portconfig.Enable = false
	}
	portconfig.PathCost = int32(config.PathCost)
	portconfig.ProtocolMigration = int32(config.ProtocolMigration)
	portconfig.AdminPointToPoint = int32(config.AdminPointToPoint)
	portconfig.AdminEdgePort = ConvertInt32ToBool(config.AdminEdgePort)
	portconfig.AdminPathCost = int32(config.AdminPathCost)
	portconfig.BridgeAssurance = ConvertInt32ToBool(config.BridgeAssurance)
	portconfig.BpduGuard = ConvertInt32ToBool(config.BpduGuard)
	portconfig.BpduGuardInterval = config.BpduGuardInterval
}

func ConvertBridgeIdToString(bridgeid stp.BridgeId) string {

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x:",
		bridgeid[0],
		bridgeid[1],
		bridgeid[2],
		bridgeid[3],
		bridgeid[4],
		bridgeid[5],
		bridgeid[6],
		bridgeid[7])
}

func ConvertAddrToString(mac [6]uint8) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		mac[0],
		mac[1],
		mac[2],
		mac[3],
		mac[4],
		mac[5])
}

//NOTE—The current IETF Bridge MIB (IETF RFC 1493) uses disabled, blocking, listening, learning, forwarding, and
//broken dot1dStpPortStates. The learning and forwarding states correspond exactly to the Learning and Forwarding Port
//States specified in this standard. Disabled, blocking, listening, and broken all correspond to the Discarding Port State —
//while those dot1dStpPortStates serve to distinguish reasons for discarding frames the operation of the Forwarding and
//Learning processes is the same for all of them. The dot1dStpPortState broken represents the failure or unavailability of
//the port’s MAC as indicated by MAC_Operational FALSE; disabled represents exclusion of the port from the active
//topology by management setting of the Administrative Port State to Disabled; blocking represents exclusion of the port
//from the active topology by the spanning tree algorithm [computing an Alternate or Backup Port Role (17.7)]; listening
//represents a port that the spanning tree algorithm has selected to be part of the active topology (computing a Root Port or
//Designated Port role) but is temporarily discarding frames to guard against loops or incorrect learning.
func GetPortState(p *stp.StpPort) (state int32) {
	/* defined by model
	type enumeration {
	          enum disabled   { value 1; }
	          enum blocking   { value 2; }
	          enum listening  { value 3; }
	          enum learning   { value 4; }
	          enum forwarding { value 5; }
	          enum broken     { value 6; }
	        }
	*/
	state = 0
	//stp.StpLogger("INFO", fmt.Sprintf("PortEnabled[%t] Learning[%t] Forwarding[%t]", p.PortEnabled, p.Learning, p.Forwarding))
	if !p.PortEnabled {
		state = 1
	} else if p.Forwarding {
		state = 5
	} else if p.Learning {
		state = 4
	} else if p.PortEnabled &&
		!p.Learning &&
		!p.Forwarding {
		state = 2
	}
	// TODO need to determine how to set listening and broken states
	return state
}

func (s *STPDServiceHandler) CreateStpGlobal(config *stpd.StpGlobal) (rv bool, err error) {
	rv = true
	stp.StpLogger("INFO", fmt.Sprintf("CreateStpGlobal (server): %s", config.AdminState))

	if config.AdminState == "UP" {
		prevState := stp.StpGlobalStateGet()
		stp.StpGlobalStateSet(stp.STP_GLOBAL_ENABLE)
		s.ReadConfigFromDB(prevState)
	} else if config.AdminState == "DOWN" {
		stp.StpGlobalStateSet(stp.STP_GLOBAL_DISABLE)
	}
	return rv, err
}

func (s *STPDServiceHandler) DeleteStpGlobal(config *stpd.StpGlobal) (bool, error) {
	return true, nil
}

func (s *STPDServiceHandler) UpdateStpGlobal(origconfig *stpd.StpGlobal, updateconfig *stpd.StpGlobal, attrset []bool, op []*stpd.PatchOpInfo) (rv bool, err error) {
	stp.StpLogger("INFO", fmt.Sprintf("UpdateStpGlobal (server): %s", updateconfig.AdminState))
	rv = true
	prevState := stp.StpGlobalStateGet()

	if updateconfig.AdminState == "UP" {
		stp.StpGlobalStateSet(stp.STP_GLOBAL_ENABLE)
	} else if updateconfig.AdminState == "DOWN" {
		stp.StpGlobalStateSet(stp.STP_GLOBAL_DISABLE_PENDING)
	}
	if prevState != stp.StpGlobalStateGet() {
		s.ReadConfigFromDB(prevState)
		if updateconfig.AdminState == "DOWN" {
			stp.StpGlobalStateSet(stp.STP_GLOBAL_DISABLE)
		}
	}
	return rv, err
}

// CreateDot1dStpBridgeConfig
func (s *STPDServiceHandler) CreateStpBridgeInstance(config *stpd.StpBridgeInstance) (rv bool, err error) {

	brgconfig := &stp.StpBridgeConfig{}
	ConvertThriftBrgConfigToStpBrgConfig(config, brgconfig)

	if brgconfig.Vlan == 0 {
		brgconfig.Vlan = stp.DEFAULT_STP_BRIDGE_VLAN
	}

	err = stp.StpBrgConfigParamCheck(brgconfig, true)
	if err == nil {
		if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {
			stp.StpLogger("INFO", "CreateStpBridgeInstance (server): created ")
			stp.StpLogger("INFO", fmt.Sprintf("addr:", config.Address))
			stp.StpLogger("INFO", fmt.Sprintf("prio:", config.Priority))
			stp.StpLogger("INFO", fmt.Sprintf("vlan:", config.Vlan))
			stp.StpLogger("INFO", fmt.Sprintf("age:", config.MaxAge))
			stp.StpLogger("INFO", fmt.Sprintf("hello:", config.HelloTime))        // int32
			stp.StpLogger("INFO", fmt.Sprintf("fwddelay:", config.ForwardDelay))  // int32
			stp.StpLogger("INFO", fmt.Sprintf("version:", config.ForceVersion))   // int32
			stp.StpLogger("INFO", fmt.Sprintf("txHoldCount", config.TxHoldCount)) //

			cfg := server.STPConfig{
				Msgtype: server.STPConfigMsgCreateBridge,
				Msgdata: brgconfig,
			}
			s.server.ConfigCh <- cfg
		}
		return true, err
	}

	return rv, err
}

func (s *STPDServiceHandler) HandleDbReadStpGlobal(dbHdl *dbutils.DBUtil) error {
	if dbHdl != nil {
		var dbObj objects.StpGlobal
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			stp.StpLogger("ERROR", "DB Query failed when retrieving StpPort objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := stpd.NewStpGlobal()
			dbObject := objList[idx].(objects.StpGlobal)
			objects.ConvertstpdStpGlobalObjToThrift(&dbObject, obj)
			_, err = s.CreateStpGlobal(obj)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (s *STPDServiceHandler) HandleDbReadStpBridgeInstance(dbHdl *dbutils.DBUtil, del bool) error {
	if dbHdl != nil {
		var dbObj objects.StpBridgeInstance
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			stp.StpLogger("ERROR", "DB Query failed when retrieving StpBridgeInstance objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := stpd.NewStpBridgeInstance()
			dbObject := objList[idx].(objects.StpBridgeInstance)
			objects.ConvertstpdStpBridgeInstanceObjToThrift(&dbObject, obj)
			if !del {
				_, err = s.CreateStpBridgeInstance(obj)
			} else {
				_, err = s.DeleteStpBridgeInstance(obj)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *STPDServiceHandler) HandleDbReadStpPort(dbHdl *dbutils.DBUtil) error {
	if dbHdl != nil {
		var dbObj objects.StpPort
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			stp.StpLogger("ERROR", "DB Query failed when retrieving StpPort objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := stpd.NewStpPort()
			dbObject := objList[idx].(objects.StpPort)
			objects.ConvertstpdStpPortObjToThrift(&dbObject, obj)
			_, err = s.CreateStpPort(obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *STPDServiceHandler) ReadConfigFromDB(prevState int) error {

	dbHdl := dbutils.NewDBUtil(stp.GetStpLogger())
	err := dbHdl.Connect()
	if err != nil {
		stp.StpLogger("ERROR", fmt.Sprintf("Failed to open connection to DB with error %s", err))
		return err
	}
	defer dbHdl.Disconnect()

	// only need to call on bootup
	if prevState == stp.STP_GLOBAL_INIT {
		if err := s.HandleDbReadStpGlobal(dbHdl); err != nil {
			stp.StpLogger("ERROR", fmt.Sprintf("Error getting All StpGlobal objects %s", err))
			return err
		}
	}
	currState := stp.StpGlobalStateGet()

	// going from enable to disable or stp is in enable state
	if (prevState != currState && currState == stp.STP_GLOBAL_ENABLE) ||
		currState == stp.STP_GLOBAL_ENABLE {
		if err := s.HandleDbReadStpBridgeInstance(dbHdl, false); err != nil {
			stp.StpLogger("ERROR", fmt.Sprintf("Error getting All StpBridgeInstance objects %s", err))
			return err
		}

		if err = s.HandleDbReadStpPort(dbHdl); err != nil {
			stp.StpLogger("ERROR", fmt.Sprintf("Error getting All StpPort objects %s", err))
			return err
		}
	} else if currState == stp.STP_GLOBAL_DISABLE_PENDING ||
		prevState == stp.STP_GLOBAL_ENABLE {
		// only need to delete the bridge instance
		// this will trigger a delete of the ports within the server
		if err := s.HandleDbReadStpBridgeInstance(dbHdl, true); err != nil {
			stp.StpLogger("ERROR", fmt.Sprintf("Error getting All StpBridgeInstance objects", err))
			return err
		}
	}
	return nil
}

func (s *STPDServiceHandler) DeleteStpBridgeInstance(config *stpd.StpBridgeInstance) (rv bool, err error) {
	rv = true
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE ||
		stp.StpGlobalStateGet() == stp.STP_GLOBAL_DISABLE_PENDING {

		// Aggregation found now lets delete
		//lacp.DeleteLaAgg(GetIdByName(config.NameKey))
		stp.StpLogger("INFO", "DeleteStpBridgeInstance (server): deleted ")
		brgconfig := &stp.StpBridgeConfig{}
		ConvertThriftBrgConfigToStpBrgConfig(config, brgconfig)
		cfg := server.STPConfig{
			Msgtype: server.STPConfigMsgDeleteBridge,
			Msgdata: brgconfig,
		}
		s.server.ConfigCh <- cfg
	}
	return rv, err
}

func (s *STPDServiceHandler) UpdateStpBridgeInstance(origconfig *stpd.StpBridgeInstance, updateconfig *stpd.StpBridgeInstance, attrset []bool, op []*stpd.PatchOpInfo) (rv bool, err error) {
	rv = true

	var b *stp.Bridge
	brgconfig := &stp.StpBridgeConfig{}
	objTyp := reflect.TypeOf(*origconfig)

	// convert thrift struct to stp struct
	ConvertThriftBrgConfigToStpBrgConfig(updateconfig, brgconfig)
	// perform paramater checks to validate the config coming down
	err = stp.StpBrgConfigParamCheck(brgconfig, false)
	if err != nil {
		return false, err
	}
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		key := stp.BridgeKey{
			Vlan: uint16(origconfig.Vlan),
		}
		// see if the bridge instance exists
		if !stp.StpFindBridgeById(key, &b) {
			return false, errors.New("Unknown Bridge in update config")
		}

		// config message data
		cfg := server.STPConfig{
			Msgdata: brgconfig,
		}

		// attribute that user is allowed to update
		attrMap := map[string]server.STPConfigMsgType{
			"MaxAge":       server.STPConfigMsgUpdateBridgeMaxAge,
			"HelloTime":    server.STPConfigMsgUpdateBridgeHelloTime,
			"ForwardDelay": server.STPConfigMsgUpdateBridgeForwardDelay,
			"TxHoldCount":  server.STPConfigMsgUpdateBridgeTxHoldCount,
			"Priority":     server.STPConfigMsgUpdateBridgePriority,
			"ForceVersion": server.STPConfigMsgUpdateBridgeForceVersion,
		}

		// important to note that the attrset starts at index 0 which is the BaseObj
		// which is not the first element on the thrift obj, thus we need to skip
		// this attribute
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			//fmt.Println("UpdateStpBridgeInstance (server): (index, objName) ", i, objName)
			if attrset[i] {
				stp.StpLogger("INFO", fmt.Sprintf("UpdateStpBridgeInstance (server): changed ", objName))

				if msgtype, ok := attrMap[objName]; ok {
					// set message type
					cfg.Msgtype = msgtype
					// send config message to server
					s.server.ConfigCh <- cfg
				}
			}
		}
	}
	return rv, err
}

func (s *STPDServiceHandler) CreateStpPort(config *stpd.StpPort) (rv bool, err error) {
	rv = true
	portconfig := &stp.StpPortConfig{}
	ConvertThriftPortConfigToStpPortConfig(config, portconfig)
	err = stp.StpPortConfigParamCheck(portconfig, false, true)
	// only create the instance if it up
	if config.AdminState == "UP" {
		if err == nil {
			if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {
				stp.StpLogger("INFO", fmt.Sprintf("CreateStpPort (server): created %#v", config))

				cfg := server.STPConfig{
					Msgtype: server.STPConfigMsgCreatePort,
					Msgdata: portconfig,
				}
				s.server.ConfigCh <- cfg
			}
			return rv, err
		} else {
			return false, err
		}
	}
	return rv, err
}

func (s *STPDServiceHandler) DeleteStpPort(config *stpd.StpPort) (rv bool, err error) {

	rv = true
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		stp.StpLogger("INFO", "DeleteStpPort (server): deleted")
		portconfig := &stp.StpPortConfig{}
		ConvertThriftPortConfigToStpPortConfig(config, portconfig)

		// assume that if the port exists that you can delete it otherwise
		// it may have been a port which Admin DOWN
		ifIndex := stp.GetIfIndexFromIntfRef(config.IntfRef)
		brgIfIndex := int32(config.Vlan)
		var p *stp.StpPort
		if stp.StpFindPortByIfIndex(ifIndex, brgIfIndex, &p) {

			cfg := server.STPConfig{
				Msgtype: server.STPConfigMsgDeletePort,
				Msgdata: config,
			}
			s.server.ConfigCh <- cfg
		}
	}
	return rv, err
}

func (s *STPDServiceHandler) UpdateStpPort(origconfig *stpd.StpPort, updateconfig *stpd.StpPort, attrset []bool, op []*stpd.PatchOpInfo) (rv bool, err error) {
	rv = true

	var p *stp.StpPort
	portconfig := &stp.StpPortConfig{}
	objTyp := reflect.TypeOf(*origconfig)
	//objVal := reflect.ValueOf(origconfig)
	//updateObjVal := reflect.ValueOf(*updateconfig)

	ConvertThriftPortConfigToStpPortConfig(updateconfig, portconfig)
	err = stp.StpPortConfigParamCheck(portconfig, true, false)
	if err != nil {
		return false, err
	}
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {
		ifIndex := stp.GetIfIndexFromIntfRef(origconfig.IntfRef)
		brgIfIndex := int32(origconfig.Vlan)
		if !stp.StpFindPortByIfIndex(ifIndex, brgIfIndex, &p) {
			if updateconfig.AdminState == "DOWN" {
				return rv, nil
			}
		}

		err = stp.StpPortConfigSave(portconfig, true)
		if err != nil {
			return false, err
		}

		// config message data
		cfg := server.STPConfig{
			Msgdata: portconfig,
		}

		attrMap := map[string]server.STPConfigMsgType{
			"Priority":          server.STPConfigMsgUpdatePortPriority,
			"AdminState":        server.STPConfigMsgUpdatePortEnable,
			"PathCost":          server.STPConfigMsgUpdatePortPathCost,
			"ProtocolMigration": server.STPConfigMsgUpdatePortProtocolMigration,
			"AdminPointToPoint": server.STPConfigMsgUpdatePortAdminPointToPoint,
			"AdminEdge":         server.STPConfigMsgUpdatePortAdminEdge,
			"AdminPathCost":     server.STPConfigMsgUpdatePortAdminPathCost,
			"BpduGuard":         server.STPConfigMsgUpdatePortBpduGuard,
			"BridgeAssurance":   server.STPConfigMsgUpdatePortBridgeAssurance,
		}

		// important to note that the attrset starts at index 0 which is the BaseObj
		// which is not the first element on the thrift obj, thus we need to skip
		// this attribute

		// check to see if AdminState is being changed
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if objName == "AdminState" {

				if updateconfig.AdminState == "UP" {
					cfg.Msgtype = server.STPConfigMsgCreatePort
				} else { // DOWN
					cfg.Msgtype = server.STPConfigMsgDeletePort
				}
				// send config message to server
				s.server.ConfigCh <- cfg
				return rv, err
			}
		}

		// handle all other attribute updates
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			//fmt.Println("UpdateDot1dStpBridgeConfig (server): (index, objName) ", i, objName)
			if attrset[i] {
				stp.StpLogger("INFO", fmt.Sprintf("StpPort (server): changed ", objName))

				if msgtype, ok := attrMap[objName]; ok {
					// set message type
					cfg.Msgtype = msgtype
					// send config message to server
					s.server.ConfigCh <- cfg
				}
			}
		}
	}
	return rv, err
}

func (s *STPDServiceHandler) GetStpBridgeInstanceState(vlan int16) (*stpd.StpBridgeInstanceState, error) {
	sbs := &stpd.StpBridgeInstanceState{}

	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		key := stp.BridgeKey{
			Vlan: uint16(vlan),
		}
		var b *stp.Bridge
		if stp.StpFindBridgeById(key, &b) {
			sbs.BridgeHelloTime = int32(b.BridgeTimes.HelloTime)
			sbs.TxHoldCount = stp.TransmitHoldCountDefault
			sbs.BridgeForwardDelay = int32(b.BridgeTimes.ForwardingDelay)
			sbs.BridgeMaxAge = int32(b.BridgeTimes.MaxAge)
			sbs.Address = ConvertAddrToString(stp.GetBridgeAddrFromBridgeId(b.BridgePriority.DesignatedBridgeId))
			sbs.Priority = int32(stp.GetBridgePriorityFromBridgeId(b.BridgePriority.DesignatedBridgeId))
			sbs.Vlan = int16(b.BrgIfIndex)
			sbs.ProtocolSpecification = 2
			//nextStpBridgeInstanceState.TimeSinceTopologyChange uint32 //The time (in hundredths of a second) since the last time a topology change was detected by the bridge entity. For RSTP, this reports the time since the tcWhile timer for any port on this Bridge was nonzero.
			//nextStpBridgeInstanceState.TopChanges              uint32 //The total number of topology changes detected by this bridge since the management entity was last reset or initialized.
			sbs.DesignatedRoot = ConvertBridgeIdToString(b.BridgePriority.RootBridgeId)
			sbs.RootCost = int32(b.BridgePriority.RootPathCost)
			sbs.RootPort = int32(b.BridgePriority.DesignatedPortId)
			sbs.MaxAge = int32(b.RootTimes.MaxAge)
			sbs.HelloTime = int32(b.RootTimes.HelloTime)
			sbs.HoldTime = int32(b.TxHoldCount)
			sbs.ForwardDelay = int32(b.RootTimes.ForwardingDelay)
			sbs.Vlan = int16(b.Vlan)
			sbs.IfIndex = b.BrgIfIndex
		} else {
			return sbs, errors.New(fmt.Sprintf("STP: Error could not find bridge vlan %d", vlan))
		}
	}
	return sbs, nil
}

// GetBulkStpBridgeInstanceState will return the status of all the stp bridges
func (s *STPDServiceHandler) GetBulkStpBridgeInstanceState(fromIndex stpd.Int, count stpd.Int) (obj *stpd.StpBridgeInstanceStateGetInfo, err error) {
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		var StpBridgeInstanceStateList []stpd.StpBridgeInstanceState = make([]stpd.StpBridgeInstanceState, count)
		var nextStpBridgeInstanceState *stpd.StpBridgeInstanceState
		var returnStpBridgeInstanceStates []*stpd.StpBridgeInstanceState
		var returnStpBridgeInstanceStateGetInfo stpd.StpBridgeInstanceStateGetInfo
		var b *stp.Bridge
		validCount := stpd.Int(0)
		toIndex := fromIndex
		obj = &returnStpBridgeInstanceStateGetInfo
		brgListLen := stpd.Int(len(stp.BridgeListTable))
		for currIndex := fromIndex; validCount != count && currIndex < brgListLen; currIndex++ {

			b = stp.BridgeListTable[currIndex]
			nextStpBridgeInstanceState = &StpBridgeInstanceStateList[validCount]
			nextStpBridgeInstanceState.BridgeHelloTime = int32(b.BridgeTimes.HelloTime)
			nextStpBridgeInstanceState.TxHoldCount = stp.TransmitHoldCountDefault
			nextStpBridgeInstanceState.BridgeForwardDelay = int32(b.BridgeTimes.ForwardingDelay)
			nextStpBridgeInstanceState.BridgeMaxAge = int32(b.BridgeTimes.MaxAge)
			nextStpBridgeInstanceState.Address = ConvertAddrToString(stp.GetBridgeAddrFromBridgeId(b.BridgePriority.DesignatedBridgeId))
			nextStpBridgeInstanceState.Priority = int32(stp.GetBridgePriorityFromBridgeId(b.BridgePriority.DesignatedBridgeId))
			nextStpBridgeInstanceState.Vlan = int16(b.BrgIfIndex)
			nextStpBridgeInstanceState.ProtocolSpecification = 2
			//nextStpBridgeInstanceState.TimeSinceTopologyChange uint32 //The time (in hundredths of a second) since the last time a topology change was detected by the bridge entity. For RSTP, this reports the time since the tcWhile timer for any port on this Bridge was nonzero.
			//nextStpBridgeInstanceState.TopChanges              uint32 //The total number of topology changes detected by this bridge since the management entity was last reset or initialized.
			nextStpBridgeInstanceState.DesignatedRoot = ConvertBridgeIdToString(b.BridgePriority.RootBridgeId)
			nextStpBridgeInstanceState.RootCost = int32(b.BridgePriority.RootPathCost)
			nextStpBridgeInstanceState.RootPort = int32(b.BridgePriority.DesignatedPortId)
			nextStpBridgeInstanceState.MaxAge = int32(b.RootTimes.MaxAge)
			nextStpBridgeInstanceState.HelloTime = int32(b.RootTimes.HelloTime)
			nextStpBridgeInstanceState.HoldTime = int32(b.TxHoldCount)
			nextStpBridgeInstanceState.ForwardDelay = int32(b.RootTimes.ForwardingDelay)
			nextStpBridgeInstanceState.Vlan = int16(b.Vlan)
			nextStpBridgeInstanceState.IfIndex = b.BrgIfIndex

			if len(returnStpBridgeInstanceStates) == 0 {
				returnStpBridgeInstanceStates = make([]*stpd.StpBridgeInstanceState, 0)
			}
			returnStpBridgeInstanceStates = append(returnStpBridgeInstanceStates, nextStpBridgeInstanceState)
			validCount++
			toIndex++
		}
		// lets try and get the next agg if one exists then there are more routes
		moreRoutes := false
		if fromIndex+count < brgListLen {
			moreRoutes = true
		}

		// lets try and get the next agg if one exists then there are more routes
		obj.StpBridgeInstanceStateList = returnStpBridgeInstanceStates
		obj.StartIdx = fromIndex
		obj.EndIdx = toIndex + 1
		obj.More = moreRoutes
		obj.Count = validCount
	}
	return obj, nil
}

func (s *STPDServiceHandler) GetStpPortState(vlan int32, intfRef string) (*stpd.StpPortState, error) {
	sps := &stpd.StpPortState{}
	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		var p *stp.StpPort
		ifIndex := stp.GetIfIndexFromIntfRef(intfRef)
		if stp.StpFindPortByIfIndex(ifIndex, vlan, &p) {

			sps.OperPointToPoint = ConvertBoolToInt32(p.OperPointToPointMAC)
			sps.Vlan = p.BrgIfIndex
			sps.OperEdgePort = ConvertBoolToInt32(p.OperEdge)
			sps.DesignatedPort = fmt.Sprintf("%d", p.PortPriority.DesignatedPortId)
			sps.AdminEdgePort = ConvertBoolToInt32(p.AdminEdge)
			sps.ForwardTransitions = int32(p.ForwardingTransitions)
			//nextStpPortState.ProtocolMigration  int32  //When operating in RSTP (version 2) mode, writing true(1) to this object forces this port to transmit RSTP BPDUs. Any other operation on this object has no effect and it always returns false(2) when read.
			sps.IntfRef = stp.GetPortNameFromIfIndex(p.IfIndex)
			//nextStpPortState.PathCost = int32(p.PortPathCost) //The contribution of this port to the path cost of paths towards the spanning tree root which include this port.  802.1D-1998 recommends that the default value of this parameter be in inverse proportion to    the speed of the attached LAN.  New implementations should support PathCost32. If the port path costs exceeds the maximum value of this object then this object should report the maximum value, namely 65535.  Applications should try to read the PathCost32 object if this object reports the maximum value.
			sps.Priority = int32(p.Priority) //The value of the priority field that is contained in the first (in network byte order) octet of the (2 octet long) Port ID.  The other octet of the Port ID is given by the value of IfIndex. On bridges supporting IEEE 802.1t or IEEE 802.1w, permissible values are 0-240, in steps of 16.
			sps.DesignatedBridge = stp.CreateBridgeIdStr(p.PortPriority.DesignatedBridgeId)
			//nextStpPortState.AdminPointToPoint  int32(p.)  //The administrative point-to-point status of the LAN segment attached to this port, using the enumeration values of the IEEE 802.1w clause.  A value of forceTrue(0) indicates that this port should always be treated as if it is connected to a point-to-point link.  A value of forceFalse(1) indicates that this port should be treated as having a shared media connection.  A value of auto(2) indicates that this port is considered to have a point-to-point link if it is an Aggregator and all of its    members are aggregatable, or if the MAC entity is configured for full duplex operation, either through auto-negotiation or by management means.  Manipulating this object changes the underlying adminPortToPortMAC.  The value of this object MUST be retained across reinitializations of the management system.
			sps.State = GetPortState(p)
			sps.Enable = ConvertBoolToInt32(p.PortEnabled)
			sps.DesignatedRoot = stp.CreateBridgeIdStr(p.PortPriority.RootBridgeId)
			sps.DesignatedCost = int32(p.PortPathCost)
			//nextStpPortState.AdminPathCost = p.AdminPathCost
			//nextStpPortState.PathCost32 = int32(p.PortPathCost)
			// Bridge Assurance
			sps.BridgeAssuranceInconsistant = ConvertBoolToInt32(p.BridgeAssuranceInconsistant)
			sps.BridgeAssurance = ConvertBoolToInt32(p.BridgeAssurance)
			// Bpdu Guard
			sps.BpduGuard = ConvertBoolToInt32(p.BpduGuard)
			sps.BpduGuardDetected = ConvertBoolToInt32(p.BPDUGuardTimer.GetCount() != 0)
			// root timers
			sps.MaxAge = int32(p.PortTimes.MaxAge)
			sps.ForwardDelay = int32(p.PortTimes.ForwardingDelay)
			sps.HelloTime = int32(p.PortTimes.HelloTime)
			// counters
			sps.StpInPkts = int64(p.StpRx)
			sps.StpOutPkts = int64(p.StpTx)
			sps.RstpInPkts = int64(p.RstpRx)
			sps.RstpOutPkts = int64(p.RstpTx)
			sps.TcInPkts = int64(p.TcRx)
			sps.TcOutPkts = int64(p.TcTx)
			sps.TcAckInPkts = int64(p.TcAckRx)
			sps.TcAckOutPkts = int64(p.TcAckTx)
			sps.PvstInPkts = int64(p.PvstRx)
			sps.PvstOutPkts = int64(p.PvstTx)
			sps.BpduInPkts = int64(p.BpduRx)
			sps.BpduOutPkts = int64(p.BpduTx)
			// fsm-states
			sps.PimPrevState = p.PimMachineFsm.GetPrevStateStr()
			sps.PimCurrState = p.PimMachineFsm.GetCurrStateStr()
			sps.PrtmPrevState = p.PrtMachineFsm.GetPrevStateStr()
			sps.PrtmCurrState = p.PrtMachineFsm.GetCurrStateStr()
			sps.PrxmPrevState = p.PrxmMachineFsm.GetPrevStateStr()
			sps.PrxmCurrState = p.PrxmMachineFsm.GetCurrStateStr()
			sps.PstmPrevState = p.PstMachineFsm.GetPrevStateStr()
			sps.PstmCurrState = p.PstMachineFsm.GetCurrStateStr()
			sps.TcmPrevState = p.TcMachineFsm.GetPrevStateStr()
			sps.TcmCurrState = p.TcMachineFsm.GetCurrStateStr()
			sps.PpmPrevState = p.PpmmMachineFsm.GetPrevStateStr()
			sps.PpmCurrState = p.PpmmMachineFsm.GetCurrStateStr()
			sps.PtxmPrevState = p.PtxmMachineFsm.GetPrevStateStr()
			sps.PtxmCurrState = p.PtxmMachineFsm.GetCurrStateStr()
			sps.PtimPrevState = p.PtmMachineFsm.GetPrevStateStr()
			sps.PtimCurrState = p.PtmMachineFsm.GetCurrStateStr()
			sps.BdmPrevState = p.BdmMachineFsm.GetPrevStateStr()
			sps.BdmCurrState = p.BdmMachineFsm.GetCurrStateStr()
			// current counts
			sps.EdgeDelayWhile = p.EdgeDelayWhileTimer.GetCount()
			sps.FdWhile = p.FdWhileTimer.GetCount()
			sps.HelloWhen = p.HelloWhenTimer.GetCount()
			sps.MdelayWhile = p.MdelayWhiletimer.GetCount()
			sps.RbWhile = p.RbWhileTimer.GetCount()
			sps.RcvdInfoWhile = p.RcvdInfoWhiletimer.GetCount()
			sps.RrWhile = p.RrWhileTimer.GetCount()
			sps.TcWhile = p.TcWhileTimer.GetCount()
			sps.BaWhile = p.BAWhileTimer.GetCount()

		} else {
			return sps, errors.New(fmt.Sprintf("STP: Error unabled to locate bridge vlan %d stp port intfref %s", vlan, intfRef))
		}
	}
	return sps, nil
}

// GetBulkAggregationLacpMemberStateCounters will return the status of all
// the lag members.
func (s *STPDServiceHandler) GetBulkStpPortState(fromIndex stpd.Int, count stpd.Int) (obj *stpd.StpPortStateGetInfo, err error) {

	if stp.StpGlobalStateGet() == stp.STP_GLOBAL_ENABLE {

		var stpPortStateList []stpd.StpPortState = make([]stpd.StpPortState, count)
		var nextStpPortState *stpd.StpPortState
		var returnStpPortStates []*stpd.StpPortState
		var returnStpPortStateGetInfo stpd.StpPortStateGetInfo
		//var a *lacp.LaAggregator
		validCount := stpd.Int(0)
		toIndex := fromIndex
		obj = &returnStpPortStateGetInfo
		stpPortListLen := stpd.Int(len(stp.PortListTable))
		stp.StpLogger("INFO", fmt.Sprintf("Total configured ports %d fromIndex %d count %d", stpPortListLen, fromIndex, count))
		for currIndex := fromIndex; validCount != count && currIndex < stpPortListLen; currIndex++ {

			//stp.StpLogger("INFO", fmt.Sprintf("CurrIndex %d stpPortListLen %d", currIndex, stpPortListLen))

			p := stp.PortListTable[currIndex]
			nextStpPortState = &stpPortStateList[validCount]

			nextStpPortState.OperPointToPoint = ConvertBoolToInt32(p.OperPointToPointMAC)
			nextStpPortState.Vlan = p.BrgIfIndex
			nextStpPortState.OperEdgePort = ConvertBoolToInt32(p.OperEdge)
			nextStpPortState.DesignatedPort = fmt.Sprintf("%d", p.PortPriority.DesignatedPortId)
			nextStpPortState.AdminEdgePort = ConvertBoolToInt32(p.AdminEdge)
			nextStpPortState.ForwardTransitions = int32(p.ForwardingTransitions)
			//nextStpPortState.ProtocolMigration  int32  //When operating in RSTP (version 2) mode, writing true(1) to this object forces this port to transmit RSTP BPDUs. Any other operation on this object has no effect and it always returns false(2) when read.
			nextStpPortState.IntfRef = stp.GetPortNameFromIfIndex(p.IfIndex)
			//nextStpPortState.PathCost = int32(p.PortPathCost) //The contribution of this port to the path cost of paths towards the spanning tree root which include this port.  802.1D-1998 recommends that the default value of this parameter be in inverse proportion to    the speed of the attached LAN.  New implementations should support PathCost32. If the port path costs exceeds the maximum value of this object then this object should report the maximum value, namely 65535.  Applications should try to read the PathCost32 object if this object reports the maximum value.
			nextStpPortState.Priority = int32(p.Priority) //The value of the priority field that is contained in the first (in network byte order) octet of the (2 octet long) Port ID.  The other octet of the Port ID is given by the value of IfIndex. On bridges supporting IEEE 802.1t or IEEE 802.1w, permissible values are 0-240, in steps of 16.
			nextStpPortState.DesignatedBridge = stp.CreateBridgeIdStr(p.PortPriority.DesignatedBridgeId)
			//nextStpPortState.AdminPointToPoint  int32(p.)  //The administrative point-to-point status of the LAN segment attached to this port, using the enumeration values of the IEEE 802.1w clause.  A value of forceTrue(0) indicates that this port should always be treated as if it is connected to a point-to-point link.  A value of forceFalse(1) indicates that this port should be treated as having a shared media connection.  A value of auto(2) indicates that this port is considered to have a point-to-point link if it is an Aggregator and all of its    members are aggregatable, or if the MAC entity is configured for full duplex operation, either through auto-negotiation or by management means.  Manipulating this object changes the underlying adminPortToPortMAC.  The value of this object MUST be retained across reinitializations of the management system.
			nextStpPortState.State = GetPortState(p)
			nextStpPortState.Enable = ConvertBoolToInt32(p.PortEnabled)
			nextStpPortState.DesignatedRoot = stp.CreateBridgeIdStr(p.PortPriority.RootBridgeId)
			nextStpPortState.DesignatedCost = int32(p.PortPathCost)
			//nextStpPortState.AdminPathCost = p.AdminPathCost
			//nextStpPortState.PathCost32 = int32(p.PortPathCost)
			// Bridge Assurance
			nextStpPortState.BridgeAssuranceInconsistant = ConvertBoolToInt32(p.BridgeAssuranceInconsistant)
			nextStpPortState.BridgeAssurance = ConvertBoolToInt32(p.BridgeAssurance)
			// Bpdu Guard
			nextStpPortState.BpduGuard = ConvertBoolToInt32(p.BpduGuard)
			nextStpPortState.BpduGuardDetected = ConvertBoolToInt32(p.BPDUGuardTimer.GetCount() != 0)
			// root timers
			nextStpPortState.MaxAge = int32(p.PortTimes.MaxAge)
			nextStpPortState.ForwardDelay = int32(p.PortTimes.ForwardingDelay)
			nextStpPortState.HelloTime = int32(p.PortTimes.HelloTime)
			// counters
			nextStpPortState.StpInPkts = int64(p.StpRx)
			nextStpPortState.StpOutPkts = int64(p.StpTx)
			nextStpPortState.RstpInPkts = int64(p.RstpRx)
			nextStpPortState.RstpOutPkts = int64(p.RstpTx)
			nextStpPortState.TcInPkts = int64(p.TcRx)
			nextStpPortState.TcOutPkts = int64(p.TcTx)
			nextStpPortState.TcAckInPkts = int64(p.TcAckRx)
			nextStpPortState.TcAckOutPkts = int64(p.TcAckTx)
			nextStpPortState.PvstInPkts = int64(p.PvstRx)
			nextStpPortState.PvstOutPkts = int64(p.PvstTx)
			nextStpPortState.BpduInPkts = int64(p.BpduRx)
			nextStpPortState.BpduOutPkts = int64(p.BpduTx)
			// fsm-states
			nextStpPortState.PimPrevState = p.PimMachineFsm.GetPrevStateStr()
			nextStpPortState.PimCurrState = p.PimMachineFsm.GetCurrStateStr()
			nextStpPortState.PrtmPrevState = p.PrtMachineFsm.GetPrevStateStr()
			nextStpPortState.PrtmCurrState = p.PrtMachineFsm.GetCurrStateStr()
			nextStpPortState.PrxmPrevState = p.PrxmMachineFsm.GetPrevStateStr()
			nextStpPortState.PrxmCurrState = p.PrxmMachineFsm.GetCurrStateStr()
			nextStpPortState.PstmPrevState = p.PstMachineFsm.GetPrevStateStr()
			nextStpPortState.PstmCurrState = p.PstMachineFsm.GetCurrStateStr()
			nextStpPortState.TcmPrevState = p.TcMachineFsm.GetPrevStateStr()
			nextStpPortState.TcmCurrState = p.TcMachineFsm.GetCurrStateStr()
			nextStpPortState.PpmPrevState = p.PpmmMachineFsm.GetPrevStateStr()
			nextStpPortState.PpmCurrState = p.PpmmMachineFsm.GetCurrStateStr()
			nextStpPortState.PtxmPrevState = p.PtxmMachineFsm.GetPrevStateStr()
			nextStpPortState.PtxmCurrState = p.PtxmMachineFsm.GetCurrStateStr()
			nextStpPortState.PtimPrevState = p.PtmMachineFsm.GetPrevStateStr()
			nextStpPortState.PtimCurrState = p.PtmMachineFsm.GetCurrStateStr()
			nextStpPortState.BdmPrevState = p.BdmMachineFsm.GetPrevStateStr()
			nextStpPortState.BdmCurrState = p.BdmMachineFsm.GetCurrStateStr()
			// current counts
			nextStpPortState.EdgeDelayWhile = p.EdgeDelayWhileTimer.GetCount()
			nextStpPortState.FdWhile = p.FdWhileTimer.GetCount()
			nextStpPortState.HelloWhen = p.HelloWhenTimer.GetCount()
			nextStpPortState.MdelayWhile = p.MdelayWhiletimer.GetCount()
			nextStpPortState.RbWhile = p.RbWhileTimer.GetCount()
			nextStpPortState.RcvdInfoWhile = p.RcvdInfoWhiletimer.GetCount()
			nextStpPortState.RrWhile = p.RrWhileTimer.GetCount()
			nextStpPortState.TcWhile = p.TcWhileTimer.GetCount()
			nextStpPortState.BaWhile = p.BAWhileTimer.GetCount()

			if len(returnStpPortStates) == 0 {
				returnStpPortStates = make([]*stpd.StpPortState, 0)
			}
			returnStpPortStates = append(returnStpPortStates, nextStpPortState)
			validCount++
			toIndex++
		}
		// lets try and get the next agg if one exists then there are more routes
		moreRoutes := false
		if fromIndex+count < stpPortListLen {
			moreRoutes = true
		}
		// lets try and get the next agg if one exists then there are more routes
		obj.StpPortStateList = returnStpPortStates
		obj.StartIdx = fromIndex
		obj.EndIdx = toIndex + 1
		obj.More = moreRoutes
		obj.Count = validCount
	}
	return obj, nil
}

// UNUSED: Actual call from user should be getting info from CONFD directly
// as it will read the info from the DB
func (s *STPDServiceHandler) GetStpBridgeInstance(vlan int16) (*stpd.StpBridgeInstance, error) {
	sps := &stpd.StpBridgeInstance{}
	return sps, nil
}

func (s *STPDServiceHandler) GetBulkStpBridgeInstance(fromIndex stpd.Int, count stpd.Int) (obj *stpd.StpBridgeInstanceGetInfo, err error) {
	var stpBridgeInstanceList []stpd.StpBridgeInstance = make([]stpd.StpBridgeInstance, count)
	var nextStpBridgeInstance *stpd.StpBridgeInstance
	var returnStpBridgeInstances []*stpd.StpBridgeInstance
	var returnStpBridgeInstanceGetInfo stpd.StpBridgeInstanceGetInfo
	//var a *lacp.LaAggregator
	validCount := stpd.Int(0)
	toIndex := fromIndex
	obj = &returnStpBridgeInstanceGetInfo
	stpDefaultBridgeListLen := stpd.Int(1)
	stp.StpLogger("INFO", fmt.Sprintf("GetBulkStpBridgeInstance (server):"))
	stp.StpLogger("INFO", fmt.Sprintf("Total default bridge instances %d fromIndex %d count %d", stpDefaultBridgeListLen, fromIndex, count))
	for currIndex := fromIndex; validCount != count && currIndex < stpDefaultBridgeListLen; currIndex++ {

		//stp.StpLogger("INFO", fmt.Sprintf("CurrIndex %d stpPortListLen %d", currIndex, stpPortListLen))

		nextStpBridgeInstance = &stpBridgeInstanceList[validCount]

		/*
			1 : i16 Vlan
			2 : string Address
			3 : i32 Priority
			4 : i32 MaxAge
			5 : i32 HelloTime
			6 : i32 ForwardDelay
			7 : i32 ForceVersion
			8 : i32 TxHoldCount

			Vlan         uint16 `SNAPROUTE: "KEY",  ACCESS:"w",  MULTIPLICITY:"*", AUTODISCOVER: "true", DEFAULT: "4095", DESCRIPTION: Each bridge is associated with a domain.  Typically this domain is represented as the vlan; The default domain is 4095, MIN: "1" ,  MAX: "4095"`
			Address      string `DESCRIPTION: The bridge identifier of the root of the spanning tree as determined by the Spanning Tree Protocol as executed by this node.  This value is used as the Root Identifier parameter in all Configuration Bridge PDUs originated by this node.,  DEFAULT: "00:00:00:00:00:00"`
			Priority     int32  `DESCRIPTION: The value of the write-able portion of the Bridge ID i.e. the first two octets of the 8 octet long Bridge ID.  The other last 6 octets of the Bridge ID are given by the value of Address. On bridges supporting IEEE 802.1t or IEEE 802.1w permissible values are 0-61440 in steps of 4096.  Extended Priority is enabled when the lower 12 bits are set using the Bridges VLAN id, MIN: "0" ,  MAX: "65535", DEFAULT: 32768`
			MaxAge       int32  `DESCRIPTION: The value that all bridges use for MaxAge when this bridge is acting as the root.  Note that 802.1D-1998 specifies that the range for this parameter is related to the value of HelloTime.  The granularity of this timer is specified by 802.1D-1998 to be 1 second.  An agent may return a badValue error if a set is attempted to a value that is not a whole number of seconds., MIN: "6" ,  MAX: "40", DEFAULT: 20`
			HelloTime    int32  `DESCRIPTION: The value that all bridges use for HelloTime when this bridge is acting as the root.  The granularity of this timer is specified by 802.1D-1998 to be 1 second.  An agent may return a badValue error if a set is attempted    to a value that is not a whole number of seconds., MIN: "1" ,  MAX: "2", DEFAULT: 2`
			ForwardDelay int32  `DESCRIPTION: The value that all bridges use for ForwardDelay when this bridge is acting as the root.  Note that 802.1D-1998 specifies that the range for this parameter is related to the value of MaxAge.  The granularity of this timer is specified by 802.1D-1998 to be 1 second.  An agent may return a badValue error if a set is attempted to a value that is not a whole number of seconds., MIN: "3" ,  MAX: "30", DEFAULT: 15`
			ForceVersion int32  `DESCRIPTION: Stp Version, SELECTION: stp(1)/rstp-pvst(2)/mstp(3), DEFAULT: 2`
			TxHoldCount  int32  `DESCRIPTION: Configures the number of BPDUs that can be sent before pausing for 1 second., MIN: "1" ,  MAX: "10", DEFAULT: 6`

		*/
		// defaults according to model values, thus if model changes these should change as well
		// perhaps we should consider opening the JSON genObjConfig.json file and filling
		// in the values this way.  For now going to hard code.
		nextStpBridgeInstance.Vlan = int16(stp.DEFAULT_STP_BRIDGE_VLAN)
		nextStpBridgeInstance.Priority = 32768
		nextStpBridgeInstance.Address = "00-00-00-00-00-00" // use switch mac
		nextStpBridgeInstance.MaxAge = int32(20)
		nextStpBridgeInstance.HelloTime = int32(2)
		nextStpBridgeInstance.ForwardDelay = int32(15)
		nextStpBridgeInstance.ForceVersion = int32(2)
		nextStpBridgeInstance.TxHoldCount = int32(6)
		// lets create the object in the stack now
		// we are going to create based on CONFD creating StpGlobal
		//s.CreateStpPort(nextStpPort)

		if len(returnStpBridgeInstances) == 0 {
			returnStpBridgeInstances = make([]*stpd.StpBridgeInstance, 0)
		}
		returnStpBridgeInstances = append(returnStpBridgeInstances, nextStpBridgeInstance)
		validCount++
		toIndex++
	}
	// lets try and get the next agg if one exists then there are more routes
	moreRoutes := false
	if fromIndex+count < stpDefaultBridgeListLen {
		moreRoutes = true
	}
	// lets try and get the next agg if one exists then there are more routes
	obj.StpBridgeInstanceList = returnStpBridgeInstances
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount
	return obj, nil
}

// GetBulkStpPort used for Auto-Discovery
func (s *STPDServiceHandler) GetBulkStpPort(fromIndex stpd.Int, count stpd.Int) (obj *stpd.StpPortGetInfo, err error) {

	var stpPortList []stpd.StpPort = make([]stpd.StpPort, count)
	var nextStpPort *stpd.StpPort
	var returnStpPorts []*stpd.StpPort
	var returnStpPortGetInfo stpd.StpPortGetInfo
	//var a *lacp.LaAggregator
	validCount := stpd.Int(0)
	toIndex := fromIndex
	obj = &returnStpPortGetInfo
	stpPortListLen := stpd.Int(len(stp.PortConfigMap))
	stp.StpLogger("INFO", fmt.Sprintf("GetBulkStpPort (server):"))
	stp.StpLogger("INFO", fmt.Sprintf("Total default ports %d fromIndex %d count %d", stpPortListLen, fromIndex, count))
	for currIndex := fromIndex; validCount != count && currIndex < stpPortListLen; currIndex++ {

		//stp.StpLogger("INFO", fmt.Sprintf("CurrIndex %d stpPortListLen %d", currIndex, stpPortListLen))

		p := stp.PortConfigMap[int32(currIndex)]
		nextStpPort = &stpPortList[validCount]

		/*
					1 : i32 Vlan
					2 : string IntfRef
					3 : i32 Priority
					4 : string AdminState
					5 : i32 PathCost
					6 : i32 PathCost32
					7 : i32 ProtocolMigration
					8 : i32 AdminPointToPoint
					9 : i32 AdminEdgePort
					10 : i32 AdminPathCost
					11 : i32 BpduGuard
					12 : i32 BpduGuardInterval
					13 : i32 BridgeAssurance

						Vlan              int32  `SNAPROUTE: "KEY", ACCESS:"rw", MULTIPLICITY:"*", AUTODISCOVER:"true", DESCRIPTION: The value of instance of the vlan object,  for the bridge corresponding to this port., MIN: "0" ,  MAX: "4094"`
			IntfRef           string `SNAPROUTE: "KEY", ACCESS:"rw", DESCRIPTION: The port number of the port for which this entry contains Spanning Tree Protocol management information. `
			Priority          int32  `DESCRIPTION: The value of the priority field that is contained in the first in network byte order octet of the 2 octet long Port ID.  The other octet of the Port ID is given by the value of StpPort. On bridges supporting IEEE 802.1t or IEEE 802.1w, permissible values are 0-240, in steps of 16., MIN: "0" ,  MAX: "255", DEFAULT: 128`
			AdminState        string `DESCRIPTION: The enabled/disabled status of the port., SELECTION: UP/DOWN, DEFAULT: UP`
			PathCost          int32  `DESCRIPTION: The contribution of this port to the path cost of paths towards the spanning tree root which include this port.  802.1D-1998 recommends that the default value of this parameter be in inverse proportion to the speed of the attached LAN.  New implementations should support PathCost32. If the port path costs exceeds the maximum value of this object then this object should report the maximum value; namely 65535.  Applications should try to read the PathCost32 object if this object reports the maximum value.  Value of 1 will force node to auto discover the value        based on the ports capabilities., MIN: "1" ,  MAX: "65535", DEFAULT: 1`
			PathCost32        int32  `DESCRIPTION: The contribution of this port to the path cost of paths towards the spanning tree root which include this port.  802.1D-1998 recommends that the default value of this parameter be in inverse proportion to the speed of the attached LAN.  This object replaces PathCost to support IEEE 802.1t. Value of 1 will force node to auto discover the value        based on the ports capabilities., MIN: "1" ,  MAX: "200000000", DEFAULT: 1`
			ProtocolMigration int32  `DESCRIPTION: When operating in RSTP (version 2) mode writing true(1) to this object forces this port to transmit RSTP BPDUs. Any other operation on this object has no effect and it always returns false(2) when read., SELECTION: false(2)/true(1), DEFAULT: 1`
			AdminPointToPoint int32  `DESCRIPTION: The administrative point-to-point status of the LAN segment attached to this port using the enumeration values of the IEEE 802.1w clause.  A value of forceTrue(0) indicates that this port should always be treated as if it is connected to a point-to-point link.  A value of forceFalse(1) indicates that this port should be treated as having a shared media connection.  A value of auto(2) indicates that this port is considered to have a point-to-point link if it is an Aggregator and all of its    members are aggregatable or if the MAC entity is configured for full duplex operation, either through auto-negotiation or by management means.  Manipulating this object changes the underlying adminPortToPortMAC.  The value of this object MUST be retained across reinitializations of the management system., SELECTION: forceTrue(0)/forceFalse(1)/auto(2), DEFAULT: 2`
			AdminEdgePort     int32  `DESCRIPTION: The administrative value of the Edge Port parameter.  A value of true(1) indicates that this port should be assumed as an edge-port and a value of false(2) indicates that this port should be assumed as a non-edge-port.  Setting this object will also cause the corresponding instance of OperEdgePort to change to the same value.  Note that even when this object's value is true the value of the corresponding instance of OperEdgePort can be false if a BPDU has been received.  The value of this object MUST be retained across reinitializations of the management system., SELECTION: false(2)/true(1), DEFAULT: 2`
			AdminPathCost     int32  `DESCRIPTION: The administratively assigned value for the contribution of this port to the path cost of paths toward the spanning tree root.  Writing a value of '0' assigns the automatically calculated default Path Cost value to the port.  If the default Path Cost is being used this object returns '0' when read.  This complements the object PathCost or PathCost32 which returns the operational value of the path cost.    The value of this object MUST be retained across reinitializations of the management system., MIN: "0" ,  MAX: "200000000", DEFAULT: 200000`
			BpduGuard         int32  `DESCRIPTION: A Port as OperEdge which receives BPDU with BpduGuard enabled will shut the port down., SELECTION: false(2)/true(1), DEFAULT: 2`
			BpduGuardInterval int32  `DESCRIPTION: The interval time to which a port will try to recover from BPDU Guard err-disable state.  If no BPDU frames are detected after this timeout plus 3 Times Hello Time then the port will transition back to Up state.  If condition is cleared manually then this operation is ignored.  If set to zero then timer is inactive and recovery is based on manual intervention. DEFAULT: 15`
			BridgeAssurance   int32  `DESCRIPTION: When enabled BPDUs will be transmitted out of all stp ports regardless of state.  When an stp port fails to receive a BPDU the port should  transition to a Blocked state.  Upon reception of BDPU after shutdown  should transition port into the bridge., SELECTION: false(2)/true(1), DEFAULT: 2`

		*/
		// defaults according to model values, thus if model changes these should change as well
		// perhaps we should consider opening the JSON genObjConfig.json file and filling
		// in the values this way.  For now going to hard code.
		nextStpPort.Vlan = int32(stp.DEFAULT_STP_BRIDGE_VLAN)
		nextStpPort.IntfRef = p.Name
		nextStpPort.Priority = int32(128)
		nextStpPort.AdminState = "UP"
		nextStpPort.PathCost = int32(1)
		nextStpPort.PathCost32 = int32(1)
		nextStpPort.ProtocolMigration = int32(1)
		nextStpPort.AdminPointToPoint = int32(2)
		nextStpPort.AdminEdgePort = int32(2)
		nextStpPort.AdminPathCost = int32(200000)
		nextStpPort.BpduGuard = int32(2)
		nextStpPort.BpduGuardInterval = int32(15)
		nextStpPort.BridgeAssurance = int32(2)

		// lets create the object in the stack now
		// we are going to create based on CONFD creating StpGlobal
		//s.CreateStpPort(nextStpPort)

		if len(returnStpPorts) == 0 {
			returnStpPorts = make([]*stpd.StpPort, 0)
		}
		returnStpPorts = append(returnStpPorts, nextStpPort)
		validCount++
		toIndex++
	}
	// lets try and get the next agg if one exists then there are more routes
	moreRoutes := false
	if fromIndex+count < stpPortListLen {
		moreRoutes = true
	}
	// lets try and get the next agg if one exists then there are more routes
	obj.StpPortList = returnStpPorts
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount
	return obj, nil
}

// UNUSED: Actual call from user should be getting info from CONFD directly
// as it will read the info from the DB
func (s *STPDServiceHandler) GetStpPort(vlan int32, intfRef string) (*stpd.StpPort, error) {
	sps := &stpd.StpPort{}
	return sps, nil
}
