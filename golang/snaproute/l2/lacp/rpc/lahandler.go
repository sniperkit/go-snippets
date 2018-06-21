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
// This file contains the Thrift server handle for config objects.
// This file also contains logic for restart, and global state change for
// config objects

// lahandler
package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"l2/lacp/protocol/drcp"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"l2/lacp/server"
	"lacpd"
	"models/objects"
	//"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"utils/dbutils"
)

const DBName string = "UsrConfDb.db"

type LACPDServiceHandler struct {
	svr *server.LAServer
}

func NewLACPDServiceHandler(svr *server.LAServer) *LACPDServiceHandler {
	lacp.LacpStartTime = time.Now()
	// link up/down events for now
	handle := &LACPDServiceHandler{
		svr: svr,
	}
	prevState := utils.LacpGlobalStateGet()
	handle.ReadConfigFromDB(prevState)
	return handle
}

func ConvertStringToUint8Array(s string) [6]uint8 {
	var arr [6]uint8
	x, _ := hex.DecodeString(s)
	for k, v := range x {
		arr[k] = uint8(v)
	}
	return arr
}

func ConvertModelLagTypeToLaAggType(yangLagType int32) uint32 {
	var LagType uint32
	if yangLagType == 0 {
		LagType = lacp.LaAggTypeLACP
	} else {
		LagType = lacp.LaAggTypeSTATIC
	}
	return LagType
}

func ConvertLaAggTypeToModelLagType(aggType uint32) int32 {
	var LagType int32
	switch aggType {
	case lacp.LaAggTypeLACP:
		LagType = 0
		break
	case lacp.LaAggTypeSTATIC:
		LagType = 1
		break
	default:
		fmt.Println("ERROR: unknown LagType %d", aggType)
	}
	return LagType
}

func ConvertModelSpeedToLaAggSpeed(yangSpeed string) int {
	speedMap := map[string]int{
		"SPEED_100Gb":   1, //EthernetSpeedSPEED100Gb,
		"SPEED_10Gb":    2, //EthernetSpeedSPEED10Gb,
		"SPEED_40Gb":    3, //EthernetSpeedSPEED40Gb,
		"SPEED_25Gb":    4, //EthernetSpeedSPEED25Gb,
		"SPEED_1Gb":     5, //EthernetSpeedSPEED1Gb,
		"SPEED_100Mb":   6, //EthernetSpeedSPEED100Mb,
		"SPEED_10Mb":    7, //EthernetSpeedSPEED10Mb,
		"SPEED_UNKNOWN": 8, //EthernetSpeedSPEEDUNKNOWN
	}

	speed, err := speedMap[yangSpeed]
	if err {
		return 8 //EthernetSpeedSPEEDUNKNOWN
	}
	return speed
}

func ConvertModelLacpModeToLaAggMode(yangLacpMode int32) uint32 {
	var mode uint32
	if yangLacpMode == 0 {
		// ACTIVE
		mode = lacp.LacpModeActive
	} else {
		// PASSIVE
		mode = lacp.LacpModePassive
	}

	return mode
}

func ConvertLaAggModeToModelLacpMode(lacpMode uint32) int32 {
	var mode int32
	if lacpMode == lacp.LacpModeActive {
		// ACTIVE
		mode = 0
	} else {
		// PASSIVE
		mode = 1
	}

	return mode
}

func ConvertModelLacpPeriodToLaAggInterval(yangInterval int32) time.Duration {
	var interval time.Duration
	if yangInterval == 1 {
		interval = lacp.LacpSlowPeriodicTime
	} else {
		interval = lacp.LacpFastPeriodicTime
	}
	return interval
}

func ConvertLaAggIntervalToLacpPeriod(interval time.Duration) int32 {
	var period int32
	switch interval {
	case lacp.LacpSlowPeriodicTime:
		period = 1
		break
	case lacp.LacpFastPeriodicTime:
		period = 0
		break
	default:
		period = 0
	}
	return period
}

func ConvertSqlBooleanToBool(sqlbool string) bool {
	if sqlbool == "true" {
		return true
	} else if sqlbool == "True" {
		return true
	} else if sqlbool == "1" {
		return true
	}
	return false
}

func ConvertAdminStateStringToBool(s string) bool {
	if s == "UP" {
		return true
	} else if s == "ON" {
		return true
	} else if s == "ENABLE" {
		return true
	}
	return false
}

func ConvertRxMachineStateToYangState(state int) int32 {
	var yangstate int32
	switch state {
	case lacp.LacpRxmStateInitialize:
		yangstate = 3
		break
	case lacp.LacpRxmStatePortDisabled:
		yangstate = 5
		break
	case lacp.LacpRxmStateExpired:
		yangstate = 1
		break
	case lacp.LacpRxmStateLacpDisabled:
		yangstate = 4
		break
	case lacp.LacpRxmStateDefaulted:
		yangstate = 2
		break
	case lacp.LacpRxmStateCurrent:
		yangstate = 0
		break
	}
	return yangstate
}

func ConvertMuxMachineStateToYangState(state int) int32 {
	var yangstate int32
	switch state {
	case lacp.LacpMuxmStateDetached, lacp.LacpMuxmStateCDetached:
		yangstate = 0
		break
	case lacp.LacpMuxmStateWaiting, lacp.LacpMuxmStateCWaiting:
		yangstate = 1
		break
	case lacp.LacpMuxmStateAttached, lacp.LacpMuxmStateCAttached:
		yangstate = 2
		break
	case lacp.LacpMuxmStateCollecting:
		yangstate = 3
		break
	case lacp.LacpMuxmStateDistributing:
		yangstate = 4
		break
	case lacp.LacpMuxStateCCollectingDistributing:
		yangstate = 5
		break
	}
	return yangstate
}

func ConvertCdmMachineStateToYangState(state int) int32 {
	var yangstate int32
	switch state {
	case lacp.LacpCdmStateNoActorChurn:
		yangstate = 0
		break
	case lacp.LacpCdmStateActorChurn:
		yangstate = 1
		break
	}
	return yangstate
}

var gAggKeyMap map[string]uint16
var gAggKeyVal uint16
var gAggKeyFreeList []uint16

func GenerateKeyByAggName(AggName string) uint16 {
	var rKey uint16
	if len(gAggKeyFreeList) == 0 {
		gAggKeyVal += 1
		rKey = gAggKeyVal
	} else {
		rKey = gAggKeyFreeList[0]
		// remove element from list
		gAggKeyFreeList = append(gAggKeyFreeList[:0], gAggKeyFreeList[1:]...)
	}
	return rKey
}

func GetKeyByAggName(AggName string) uint16 {

	var Key uint16
	if gAggKeyMap == nil {
		gAggKeyMap = make(map[string]uint16)
		gAggKeyFreeList = make([]uint16, 0)
	}

	if _, ok := gAggKeyMap[AggName]; ok {
		Key = gAggKeyMap[AggName]
	} else {
		Key = GenerateKeyByAggName(AggName)
		// store the newly generated Key
		gAggKeyMap[AggName] = Key
	}
	return Key
}

func (la *LACPDServiceHandler) HandleDbReadLacpGlobal(dbHdl *dbutils.DBUtil) error {
	if dbHdl != nil {
		var dbObj objects.LacpGlobal
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			fmt.Println("DB Query failed when retrieving LacpGlobal objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := lacpd.NewLacpGlobal()
			dbObject := objList[idx].(objects.LacpGlobal)
			objects.ConvertlacpdLacpGlobalObjToThrift(&dbObject, obj)
			_, err = la.CreateLacpGlobal(obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (la *LACPDServiceHandler) HandleDbReadDistributedRelay(dbHdl *dbutils.DBUtil, del bool) error {
	if dbHdl != nil {
		var dbObj objects.DistributedRelay
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			fmt.Println("DB Query failed when retrieving LaPortChannel objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := lacpd.NewDistributedRelay()
			dbObject := objList[idx].(objects.DistributedRelay)
			objects.ConvertlacpdDistributedRelayObjToThrift(&dbObject, obj)
			if !del {
				_, err = la.CreateDistributedRelay(obj)
			} else {
				_, err = la.DeleteDistributedRelay(obj)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (la *LACPDServiceHandler) HandleDbReadLaPortChannel(dbHdl *dbutils.DBUtil, del bool) error {
	if dbHdl != nil {
		var dbObj objects.LaPortChannel
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			fmt.Println("DB Query failed when retrieving LaPortChannel objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := lacpd.NewLaPortChannel()
			dbObject := objList[idx].(objects.LaPortChannel)
			objects.ConvertlacpdLaPortChannelObjToThrift(&dbObject, obj)
			if !del {
				_, err = la.CreateLaPortChannel(obj)
			} else {
				_, err = la.DeleteLaPortChannel(obj)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (la *LACPDServiceHandler) ReadConfigFromDB(prevState int) error {
	dbHdl := dbutils.NewDBUtil(utils.GetLaLogger())
	err := dbHdl.Connect()
	if err != nil {
		fmt.Printf("Failed to open connection to the DB with error %s", err)
		return err
	}
	defer dbHdl.Disconnect()

	logger := utils.GetLaLogger()

	if prevState == utils.LACP_GLOBAL_INIT {

		if err := la.HandleDbReadLacpGlobal(dbHdl); err != nil {
			fmt.Println("Error getting All LacpGlobal objects")
			return err
		}
	}
	currState := utils.LacpGlobalStateGet()

	logger.Info(fmt.Sprintf("Global State prev %d curr %d", prevState, currState))

	if currState == utils.LACP_GLOBAL_DISABLE_PENDING ||
		prevState == utils.LACP_GLOBAL_ENABLE {

		// lets delete the Aggregator first
		if err := la.HandleDbReadLaPortChannel(dbHdl, true); err != nil {
			fmt.Println("Error getting All LaPortChannel objects")
			return err
		}

		if err := la.HandleDbReadDistributedRelay(dbHdl, true); err != nil {
			fmt.Println("Error getting All DistributedRelay objects")
			return err
		}
	} else if prevState != currState {

		if err := la.HandleDbReadDistributedRelay(dbHdl, false); err != nil {
			fmt.Println("Error getting All DistributedRelay objects")
			return err
		}

		if err := la.HandleDbReadLaPortChannel(dbHdl, false); err != nil {
			fmt.Println("Error getting All LaPortChannel objects")
			return err
		}
	}
	return nil
}

func (la *LACPDServiceHandler) CreateLacpGlobal(config *lacpd.LacpGlobal) (bool, error) {
	if config.AdminState == "UP" {
		prevState := utils.LacpGlobalStateGet()
		utils.LacpGlobalStateSet(utils.LACP_GLOBAL_ENABLE)
		la.ReadConfigFromDB(prevState)
	} else if config.AdminState == "DOWN" {
		utils.LacpGlobalStateSet(utils.LACP_GLOBAL_DISABLE)
	}
	return true, nil
}

// can't delete an autocreated object
func (la *LACPDServiceHandler) DeleteLacpGlobal(config *lacpd.LacpGlobal) (bool, error) {
	return true, nil
}

func (la *LACPDServiceHandler) UpdateLacpGlobal(origconfig *lacpd.LacpGlobal, updateconfig *lacpd.LacpGlobal, attrset []bool, op []*lacpd.PatchOpInfo) (bool, error) {
	prevState := utils.LacpGlobalStateGet()

	if updateconfig.AdminState == "UP" {
		utils.LacpGlobalStateSet(utils.LACP_GLOBAL_ENABLE)
	} else if updateconfig.AdminState == "DOWN" {
		utils.LacpGlobalStateSet(utils.LACP_GLOBAL_DISABLE_PENDING)
	}
	logger := utils.GetLaLogger()
	logger.Info(fmt.Sprintf("Global State Update AdminState %s prev %d curr %d", updateconfig.AdminState, prevState, utils.LacpGlobalStateGet()))

	if prevState != utils.LacpGlobalStateGet() {
		la.ReadConfigFromDB(prevState)
		if updateconfig.AdminState == "DOWN" {
			utils.LacpGlobalStateSet(utils.LACP_GLOBAL_DISABLE)
		}
	}
	return true, nil
}

// CreateLaPortChannel will create an lacp lag
//	1 : i32 	LagType  (0 == LACP, 1 == STATIC)
//	2 : string 	Description
//	3 : bool 	Enabled
//	4 : i16 	Mtu
//	5 : i16 	MinLinks
//	6 : string 	Type
//	7 : string 	NameKey
//	8 : i32 	Interval (0 == LONG, 1 == SHORT)
//	9 : i32 	LacpMode (0 == ACTIVE, 1 == PASSIVE)
//	10 : string SystemIdMac
//	11 : i16 	SystemPriority
func (la *LACPDServiceHandler) CreateLaPortChannel(config *lacpd.LaPortChannel) (bool, error) {

	aggModeMap := map[uint32]uint32{
		//LacpActivityTypeACTIVE:  "ACTIVE",
		//LacpActivityTypeSTANDBY: "STANDBY",
		lacp.LacpModeOn:      lacp.LacpModeOn,
		lacp.LacpModeActive:  lacp.LacpModeActive,
		lacp.LacpModePassive: lacp.LacpModePassive,
	}
	aggIntervalToTimeoutMap := map[time.Duration]time.Duration{
		//LacpPeriodTypeSLOW: "LONG",
		//LacpPeriodTypeFAST: "SHORT",
		//lacp.LacpSlowPeriodicTime: "LONG",
		//lacp.LacpFastPeriodicTime: "SHORT",
		lacp.LacpSlowPeriodicTime: lacp.LacpLongTimeoutTime,
		lacp.LacpFastPeriodicTime: lacp.LacpShortTimeoutTime,
	}

	nameKey := config.IntfRef
	switchIdMac := config.SystemIdMac
	if config.SystemIdMac == "00:00:00:00:00:00" ||
		config.SystemIdMac == "00-00-00-00-00-00" ||
		config.SystemIdMac == "" {
		tmpmac := utils.GetSwitchMac()
		switchIdMac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", tmpmac[0], tmpmac[1], tmpmac[2], tmpmac[3], tmpmac[4], tmpmac[5])
	}

	var a *lacp.LaAggregator
	if lacp.LaFindAggByName(nameKey, &a) {

		return false, errors.New(fmt.Sprintf("LACP: Error trying to create Lag %s that already exists", config.IntfRef))

	} else {
		id := GetKeyByAggName(nameKey)
		conf := &lacp.LaAggConfig{
			Id:  int(id),
			Key: id,
			// Identifier of the lag
			Name: nameKey,
			// Type of LAG STATIC or LACP
			Type:     ConvertModelLagTypeToLaAggType(config.LagType),
			MinLinks: uint16(config.MinLinks),
			Enabled:  ConvertAdminStateStringToBool(config.AdminState),
			// lacp config
			Lacp: lacp.LacpConfigInfo{
				Interval:       ConvertModelLacpPeriodToLaAggInterval(config.Interval),
				Mode:           ConvertModelLacpModeToLaAggMode(config.LacpMode),
				SystemIdMac:    switchIdMac,
				SystemPriority: uint16(config.SystemPriority),
			},
			HashMode: uint32(config.LagHash),
		}
		for _, intfref := range config.IntfRefList {
			ifindex := utils.GetIfIndexFromName(intfref)
			conf.LagMembers = append(conf.LagMembers, uint16(ifindex))
		}
		err1 := lacp.LaAggConfigAggCreateCheck(conf)
		err2 := lacp.LaAggConfigParamCheck(conf)
		if err1 != nil {
			return false, err1
		} else if err2 != nil {
			return false, err2
		} else {
			if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {

				cfg := server.LAConfig{
					Msgtype: server.LAConfigMsgCreateLaPortChannel,
					Msgdata: conf,
				}
				la.svr.ConfigCh <- cfg
				//lacp.CreateLaAgg(conf)

				for _, intfref := range config.IntfRefList {
					mode, ok := aggModeMap[uint32(conf.Lacp.Mode)]
					if !ok || conf.Type == lacp.LaAggTypeSTATIC {
						mode = lacp.LacpModeOn
					}

					timeout, ok := aggIntervalToTimeoutMap[conf.Lacp.Interval]
					if !ok {
						timeout = lacp.LacpLongTimeoutTime
					}
					ifindex := utils.GetIfIndexFromName(intfref)
					conf := &lacp.LaAggPortConfig{
						Id:       uint16(ifindex),
						Prio:     uint16(conf.Lacp.SystemPriority),
						Key:      uint16(conf.Key),
						AggId:    int(conf.Id),
						Enable:   conf.Enabled,
						Mode:     int(mode),
						Timeout:  timeout,
						TraceEna: true,
					}

					cfg := server.LAConfig{
						Msgtype: server.LAConfigMsgCreateLaAggPort,
						Msgdata: conf,
					}
					la.svr.ConfigCh <- cfg
				}
			}
		}
	}
	return true, nil
}

func (la *LACPDServiceHandler) DeleteLaPortChannel(config *lacpd.LaPortChannel) (bool, error) {
	err := lacp.LaAggConfigDeleteCheck(config.IntfRef)
	if err == nil {
		if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE ||
			utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_DISABLE_PENDING {

			logger := utils.GetLaLogger()
			logger.Info(fmt.Sprintln("Deleting La PortChannel", config.IntfRef))
			//nameKey := fmt.Sprintf("agg-%d", config.LagId)
			id := GetKeyByAggName(config.IntfRef)
			conf := &lacp.LaAggConfig{
				Id: int(id),
			}
			cfg := server.LAConfig{
				Msgtype: server.LAConfigMsgDeleteLaPortChannel,
				Msgdata: conf,
			}
			la.svr.ConfigCh <- cfg
		}

		return true, nil
	}
	return false, err
}

func GetAddDelMembers(orig []uint16, update []int32) (add, del []int32) {

	origMap := make(map[int32]bool, len(orig))
	// make the origional list a map
	for _, m := range orig {
		origMap[int32(m)] = false
	}

	// iterate through list, mark found entries
	// if entry is not found then added to add slice
	for _, m := range update {
		if _, ok := origMap[m]; ok {
			// found entry, don't want to just
			// remove from map in case the user supplied duplicate values
			origMap[m] = true
		} else {
			add = append(add, m)
		}
	}

	// iterate through map to find entries wither were not found in update
	for m, v := range origMap {
		if !v {
			del = append(del, m)
		}
	}
	return
}

func (la *LACPDServiceHandler) UpdateLaPortChannel(origconfig *lacpd.LaPortChannel, updateconfig *lacpd.LaPortChannel, attrset []bool, op []*lacpd.PatchOpInfo) (bool, error) {

	aggIntervalToTimeoutMap := map[time.Duration]time.Duration{
		//LacpPeriodTypeSLOW: "LONG",
		//LacpPeriodTypeFAST: "SHORT",
		lacp.LacpSlowPeriodicTime: lacp.LacpLongTimeoutTime,
		lacp.LacpFastPeriodicTime: lacp.LacpShortTimeoutTime,
	}

	objTyp := reflect.TypeOf(*origconfig)
	//objVal := reflect.ValueOf(origconfig)
	//updateObjVal := reflect.ValueOf(*updateconfig)

	nameKey := updateconfig.IntfRef

	id := GetKeyByAggName(nameKey)
	conf := &lacp.LaAggConfig{
		Id:  int(id),
		Key: id,
		// Identifier of the lag
		Name: nameKey,
		// Type of LAG STATIC or LACP
		Type:     ConvertModelLagTypeToLaAggType(updateconfig.LagType),
		MinLinks: uint16(updateconfig.MinLinks),
		Enabled:  ConvertAdminStateStringToBool(updateconfig.AdminState),
		// lacp config
		Lacp: lacp.LacpConfigInfo{
			Interval:       ConvertModelLacpPeriodToLaAggInterval(updateconfig.Interval),
			Mode:           ConvertModelLacpModeToLaAggMode(updateconfig.LacpMode),
			SystemIdMac:    updateconfig.SystemIdMac,
			SystemPriority: uint16(updateconfig.SystemPriority),
		},
		HashMode: uint32(updateconfig.LagHash),
	}

	ifindexList := make([]int32, 0)
	for _, intfref := range updateconfig.IntfRefList {
		ifindex := utils.GetIfIndexFromName(intfref)
		if ifindex != 0 {
			ifindexList = append(ifindexList, ifindex)
		}
	}
	addList, delList := GetAddDelMembers(lacp.ConfigAggMap[updateconfig.IntfRef].LagMembers, ifindexList)
	addPorts := make([]uint16, len(addList))
	delPorts := make([]uint16, len(delList))
	for _, p := range addList {
		addPorts = append(addPorts, uint16(p))
	}
	for _, p := range delList {
		delPorts = append(delPorts, uint16(p))
	}
	err1 := lacp.LaAggConfigAggPortUpdateCheck(updateconfig.IntfRef, addPorts, delPorts)
	err2 := lacp.LaAggConfigParamCheck(conf)
	if err1 != nil {
		return false, err1
	} else if err2 != nil {
		return false, err2
	} else {
		if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {

			// lets deal with Members attribute first
			for i := 0; i < objTyp.NumField(); i++ {
				objName := objTyp.Field(i).Name
				//fmt.Println("UpdateAggregationLacpConfig (server): (index, objName) ", i, objName)
				if attrset[i] {
					fmt.Println("UpdateAggregationLacpConfig (server): changed ", objName)
					if objName == "IntfRefList" {
						var a *lacp.LaAggregator

						if lacp.LaFindAggById(int(conf.Id), &a) {
							for _, m := range delList {
								conf := &lacp.LaAggPortConfig{
									Id: uint16(m),
								}

								cfg := server.LAConfig{
									Msgtype: server.LAConfigMsgDeleteLaAggPort,
									Msgdata: conf,
								}
								la.svr.ConfigCh <- cfg
							}
							for _, ifindex := range addList {
								mode := int(a.Config.Mode)
								if a.AggType == lacp.LaAggTypeSTATIC {
									mode = lacp.LacpModeOn
								}

								timeout, ok := aggIntervalToTimeoutMap[a.Config.Interval]
								if !ok {
									timeout = lacp.LacpLongTimeoutTime
								}
								id := GetKeyByAggName(nameKey)
								conf := &lacp.LaAggPortConfig{
									Id:       uint16(ifindex),
									Prio:     uint16(a.Config.SystemPriority),
									Key:      id,
									AggId:    int(id),
									Enable:   conf.Enabled,
									Mode:     mode,
									Timeout:  timeout,
									TraceEna: true,
								}

								cfg := server.LAConfig{
									Msgtype: server.LAConfigMsgCreateLaAggPort,
									Msgdata: conf,
								}
								la.svr.ConfigCh <- cfg
							}
						}
					}
				}
			}

			attrMap := map[string]server.LaConfigMsgType{
				"AdminState":     server.LAConfigMsgUpdateLaPortChannelAdminState,
				"LagType":        server.LAConfigMsgUpdateLaPortChannelLagType,
				"LagHash":        server.LAConfigMsgUpdateLaPortChannelLagHash,
				"LacpMode":       server.LAConfigMsgUpdateLaPortChannelAggMode,
				"Interval":       server.LAConfigMsgUpdateLaPortChannelPeriod,
				"SystemIdMac":    server.LAConfigMsgUpdateLaPortChannelSystemIdMac,
				"SystemPriority": server.LAConfigMsgUpdateLaPortChannelSystemPriority,
			}

			// important to note that the attrset starts at index 0 which is the BaseObj
			// which is not the first element on the thrift obj, thus we need to skip
			// this attribute
			for i := 0; i < objTyp.NumField(); i++ {
				objName := objTyp.Field(i).Name
				//fmt.Println("UpdateAggregationLacpConfig (server): (index, objName) ", i, objName)
				if attrset[i] {
					fmt.Println("UpdateAggregationLacpConfig (server): changed ", objName)
					if msgtype, ok := attrMap[objName]; ok {
						// set message type
						cfg := server.LAConfig{
							Msgdata: conf,
							Msgtype: msgtype,
						}
						if objName == "AdminState" {
							lacp.SaveLaAggConfig(conf)
							la.svr.ConfigCh <- cfg
							return true, nil

						} else {
							lacp.SaveLaAggConfig(conf)
							la.svr.ConfigCh <- cfg
						}
					}
				}
			}
		}
	}
	return true, nil
}

// SetLaAggType will set whether the agg is static or lacp enabled
func SetLaAggType(conf *lacp.LaAggConfig) error {
	return SetLaAggMode(conf)
}

// SetLaAggPortMode will set the lacp mode of the port based
// on the model values
func SetLaAggMode(conf *lacp.LaAggConfig) error {

	var a *lacp.LaAggregator
	var p *lacp.LaAggPort
	if lacp.LaFindAggById(conf.Id, &a) {

		if conf.Type == lacp.LaAggTypeSTATIC {
			// configured ports
			for _, pId := range a.PortNumList {
				if lacp.LaFindPortById(uint16(pId), &p) {
					lacp.SetLaAggPortLacpMode(uint16(pId), lacp.LacpModeOn)
				}
			}
		} else {
			for _, pId := range a.PortNumList {
				if lacp.LaFindPortById(uint16(pId), &p) {
					lacp.SetLaAggPortLacpMode(uint16(pId), int(conf.Lacp.Mode))
				}
			}
		}
	}
	lacp.SaveLaAggConfig(conf)

	return nil
}

func SetLaAggHashMode(conf *lacp.LaAggConfig) error {
	lacp.SetLaAggHashMode(conf.Id, conf.HashMode)
	return nil
}

func SetLaAggPeriod(conf *lacp.LaAggConfig) error {
	var a *lacp.LaAggregator
	if lacp.LaFindAggById(conf.Id, &a) {
		// configured ports
		for _, pId := range a.PortNumList {
			lacp.SetLaAggPortLacpPeriod(uint16(pId), conf.Lacp.Interval)
		}
	}
	return nil
}

func SetLaAggSystemInfo(conf *lacp.LaAggConfig) error {
	var a *lacp.LaAggregator
	if lacp.LaFindAggById(conf.Id, &a) {
		// configured ports
		for _, pId := range a.PortNumList {
			lacp.SetLaAggPortSystemInfo(uint16(pId), conf.Lacp.SystemIdMac, conf.Lacp.SystemPriority)
		}
	}
	return nil
}

// SetPortLacpLogEnable will enable on a per port basis logging
// modStr - PORT, RXM, TXM, PTXM, TXM, CDM, ALL
// modStr can be a string containing one or more of the above
func (la *LACPDServiceHandler) SetPortLacpLogEnable(Id lacpd.Uint16, modStr string, ena bool) (lacpd.Int, error) {
	modules := make(map[string]chan bool)
	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(Id), &p) {
		modules["RXM"] = p.RxMachineFsm.RxmLogEnableEvent
		modules["TXM"] = p.TxMachineFsm.TxmLogEnableEvent
		modules["PTXM"] = p.PtxMachineFsm.PtxmLogEnableEvent
		modules["TXM"] = p.TxMachineFsm.TxmLogEnableEvent
		modules["CDM"] = p.CdMachineFsm.CdmLogEnableEvent
		modules["MUXM"] = p.MuxMachineFsm.MuxmLogEnableEvent

		for k, v := range modules {
			if strings.Contains(k, "PORT") || strings.Contains(k, "ALL") {
				p.EnableLogging(ena)
			}
			if strings.Contains(k, modStr) || strings.Contains(k, "ALL") {
				v <- ena
			}
		}
		return 0, nil
	}
	return 1, errors.New(fmt.Sprintf("LACP: LOG set failed,  Unable to find Port", Id))
}

func (la *LACPDServiceHandler) GetLaPortChannelState(IntfRef string) (*lacpd.LaPortChannelState, error) {
	pcs := &lacpd.LaPortChannelState{}

	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {
		var a *lacp.LaAggregator
		id := GetKeyByAggName(IntfRef)
		if lacp.LaFindAggById(int(id), &a) {
			pcs.IntfRef = a.AggName
			pcs.IfIndex = int32(a.HwAggId)
			pcs.PresentInHw = a.PresentInHw
			pcs.LagType = ConvertLaAggTypeToModelLagType(a.AggType)
			pcs.AdminState = "DOWN"
			if a.AdminState {
				pcs.AdminState = "UP"
			}
			pcs.OperState = "DOWN"
			if a.OperState {
				pcs.OperState = "UP"
			}
			pcs.MinLinks = int16(a.AggMinLinks)
			pcs.Interval = ConvertLaAggIntervalToLacpPeriod(a.Config.Interval)
			pcs.LacpMode = ConvertLaAggModeToModelLacpMode(a.Config.Mode)
			pcs.SystemIdMac = a.Config.SystemIdMac
			pcs.SystemPriority = int16(a.Config.SystemPriority)
			pcs.LagHash = int32(a.LagHash)
			//pcs.Ifindex = int32(a.HwAggId)
			for _, m := range a.PortNumList {
				name := utils.GetNameFromIfIndex(int32(m))
				if name != "" {
					pcs.IntfRefList = append(pcs.IntfRefList, name)
				}
				var p *lacp.LaAggPort
				if lacp.LaFindPortById(m, &p) {
					if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDistributingBit) {
						distName := utils.GetNameFromIfIndex(int32(m))
						if distName != "" {
							pcs.IntfRefListUpInBundle = append(pcs.IntfRefListUpInBundle, distName)
						}
					}
				}
			}
		} else {
			return nil, errors.New(fmt.Sprintln("LACP: Unable to find port channel from LagId", IntfRef))
		}
	} else {
		//fmt.Println("Lacp Global Disabled, returning saved config map=%v intf=%s\n", lacp.ConfigAggMap, IntfRef)
		if ac, ok := lacp.ConfigAggMap[IntfRef]; ok {
			fmt.Println("Found", IntfRef)
			/*
						// Aggregator name
				Name string
				// Aggregator_MAC_address
				Mac [6]uint8
				// Aggregator_Identifier
				Id int
				// Actor_Admin_Aggregator_Key
				Key uint16
				// Aggregator Type, LACP or STATIC
				Type uint32
				// Minimum number of links
				MinLinks uint16
				// Enabled
				Enabled bool
				// LAG_ports
				LagMembers []uint16

				// System to attach this agg to
				Lacp LacpConfigInfo

				// mau properties of each link
				Properties PortProperties

				// hash config
				HashMode uint32
			*/
			pcs.IntfRef = IntfRef
			pcs.IfIndex = 0
			pcs.PresentInHw = false
			pcs.LagType = int32(ac.Type)

			pcs.AdminState = "DOWN"
			if ac.Enabled {
				pcs.AdminState = "UP"
			}
			/*
				Global state is down so all aggs should be down as they are not
				provisioned
			*/
			pcs.OperState = "DOWN"
			pcs.MinLinks = int16(ac.MinLinks)
			pcs.Interval = ConvertLaAggIntervalToLacpPeriod(ac.Lacp.Interval)
			pcs.LacpMode = ConvertLaAggModeToModelLacpMode(ac.Lacp.Mode)
			pcs.SystemIdMac = ac.Lacp.SystemIdMac
			pcs.SystemPriority = int16(ac.Lacp.SystemPriority)
			pcs.LagHash = int32(ac.HashMode)
			//pcs.Ifindex = int32(a.HwAggId)
			for _, m := range ac.LagMembers {
				name := utils.GetNameFromIfIndex(int32(m))
				if name != "" {
					pcs.IntfRefList = append(pcs.IntfRefList, name)
				}
			}
		} else {
			return nil, errors.New(fmt.Sprintln("LACP: Unable to find port channel from LagId", IntfRef, "and LACP Global is disabled"))
		}
	}
	return pcs, nil
}

// GetBulkLaAggrGroupState will return the status of all the lag groups
// All lag groups are stored in a map, thus we will assume that the order
// at which a for loop iterates over the map is preserved.  It is assumed
// that at the time of this operation that no new aggragators are added,
// otherwise can get inconsistent results
func (la *LACPDServiceHandler) GetBulkLaPortChannelState(fromIndex lacpd.Int, count lacpd.Int) (obj *lacpd.LaPortChannelStateGetInfo, err error) {

	var lagStateList []lacpd.LaPortChannelState = make([]lacpd.LaPortChannelState, count)
	var nextLagState *lacpd.LaPortChannelState
	var returnLagStates []*lacpd.LaPortChannelState
	var returnLagStateGetInfo lacpd.LaPortChannelStateGetInfo
	var a *lacp.LaAggregator
	var ac *lacp.LaAggConfig
	validCount := lacpd.Int(0)
	toIndex := fromIndex
	moreRoutes := false
	obj = &returnLagStateGetInfo

	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_DISABLE {
		var currIndex lacpd.Int
		for currIndex = 0; validCount != count && lacp.LaAggConfigGetByIndex(int(currIndex), &ac); currIndex++ {

			if currIndex < fromIndex {
				continue
			} else {

				nextLagState = &lagStateList[validCount]
				nextLagState.IntfRef = ac.Name
				nextLagState.IfIndex = 0
				nextLagState.PresentInHw = false
				nextLagState.LagType = int32(ac.Type)
				nextLagState.AdminState = "DOWN"
				if ac.Enabled {
					nextLagState.AdminState = "UP"
				}
				nextLagState.OperState = "DOWN"
				nextLagState.MinLinks = int16(ac.MinLinks)
				nextLagState.Interval = ConvertLaAggIntervalToLacpPeriod(ac.Lacp.Interval)
				nextLagState.LacpMode = ConvertLaAggModeToModelLacpMode(ac.Lacp.Mode)
				nextLagState.SystemIdMac = ac.Lacp.SystemIdMac
				nextLagState.SystemPriority = int16(ac.Lacp.SystemPriority)
				nextLagState.LagHash = int32(ac.HashMode)
				for _, m := range ac.LagMembers {
					name := utils.GetNameFromIfIndex(int32(m))
					if name != "" {
						nextLagState.IntfRefList = append(nextLagState.IntfRefList, name)
					}

				}
				if len(returnLagStates) == 0 {
					returnLagStates = make([]*lacpd.LaPortChannelState, 0)
				}
				returnLagStates = append(returnLagStates, nextLagState)
				validCount++
				toIndex++
			}
		}
		// lets try and get the next agg if one exists then there are more routes
		if ac != nil {
			moreRoutes = lacp.LaAggConfigGetByIndex(int(currIndex), &ac)
		}
	} else {
		for currIndex := lacpd.Int(0); validCount != count && lacp.LaGetAggNext(&a); currIndex++ {

			if currIndex < fromIndex {
				continue
			} else {

				nextLagState = &lagStateList[validCount]
				nextLagState.IntfRef = a.AggName
				nextLagState.IfIndex = int32(a.HwAggId)
				nextLagState.PresentInHw = a.PresentInHw
				nextLagState.LagType = ConvertLaAggTypeToModelLagType(a.AggType)
				nextLagState.AdminState = "DOWN"
				if a.AdminState {
					nextLagState.AdminState = "UP"
				}
				nextLagState.OperState = "DOWN"
				if a.OperState {
					nextLagState.OperState = "UP"
				}
				nextLagState.MinLinks = int16(a.AggMinLinks)
				nextLagState.Interval = ConvertLaAggIntervalToLacpPeriod(a.Config.Interval)
				nextLagState.LacpMode = ConvertLaAggModeToModelLacpMode(a.Config.Mode)
				nextLagState.SystemIdMac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", a.AggMacAddr[0],
					a.AggMacAddr[1],
					a.AggMacAddr[2],
					a.AggMacAddr[3],
					a.AggMacAddr[4],
					a.AggMacAddr[5])
				nextLagState.SystemPriority = int16(a.AggPriority)
				nextLagState.LagHash = int32(a.LagHash)
				if len(a.PortNumList) > 0 {
					nextLagState.IntfRefList = make([]string, 0)
				}
				if len(a.DistributedPortNumList) > 0 {
					nextLagState.IntfRefListUpInBundle = make([]string, 0)
				}
				for _, m := range a.PortNumList {
					name := utils.GetNameFromIfIndex(int32(m))
					if name != "" {
						nextLagState.IntfRefList = append(nextLagState.IntfRefList, name)
					}
					var p *lacp.LaAggPort
					if lacp.LaFindPortById(m, &p) {
						if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDistributingBit) {
							distName := utils.GetNameFromIfIndex(int32(m))
							if distName != "" {
								nextLagState.IntfRefListUpInBundle = append(nextLagState.IntfRefListUpInBundle, distName)
							}
						}
					}
				}

				if len(returnLagStates) == 0 {
					returnLagStates = make([]*lacpd.LaPortChannelState, 0)
				}
				returnLagStates = append(returnLagStates, nextLagState)
				validCount++
				toIndex++
			}
		}
		// lets try and get the next agg if one exists then there are more routes
		if a != nil {
			moreRoutes = lacp.LaGetAggNext(&a)
		}
	}
	obj.LaPortChannelStateList = returnLagStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount

	return obj, nil
}

func (la *LACPDServiceHandler) GetLaPortChannelIntfRefListState(intfref string) (*lacpd.LaPortChannelIntfRefListState, error) {
	pcms := &lacpd.LaPortChannelIntfRefListState{}
	var p *lacp.LaAggPort
	id := utils.GetIfIndexFromName(intfref)
	aggTimeoutToIntervalMap := map[time.Duration]int32{
		//LacpPeriodTypeSLOW: "LONG",
		//LacpPeriodTypeFAST: "SHORT",
		//lacp.LacpSlowPeriodicTime: "LONG",
		//lacp.LacpFastPeriodicTime: "SHORT",
		lacp.LacpLongTimeoutTime:  0,
		lacp.LacpShortTimeoutTime: 1,
	}

	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_DISABLE {
		var ac *lacp.LaAggConfig
		if lacp.LaAggConfigDoesIntfRefListMemberExist(intfref, &ac) {
			pcms.Aggregatable = true
			pcms.Collecting = false
			pcms.Distributing = false
			pcms.Defaulted = true
			pcms.OperState = "DOWN"

			// out of sync
			pcms.Synchronization = 1
			timeout, ok := aggTimeoutToIntervalMap[ac.Lacp.Interval]
			if !ok {
				timeout = 0
			}

			pcms.Timeout = timeout
			pcms.Activity = ConvertLaAggModeToModelLacpMode(ac.Lacp.Mode)

			pcms.OperKey = int16(0)
			pcms.IntfRef = intfref
			pcms.LagIntfRef = ac.Name

			pcms.DrniName = ""
			pcms.DrniSynced = false

			/*
				// partner info
				pcms.PartnerId = p.PartnerOper.System.LacpSystemConvertSystemIdToString()
				//ms.PartnerKey = int16(p.PartnerOper.Key)

				// System
				pcms.SystemIdMac = p.ActorOper.System.LacpSystemConvertSystemIdToString()[6:]
				pcms.LagType = ConvertLaAggTypeToModelLagType(p.AggAttached.AggType)
				pcms.SystemId = p.ActorOper.System.LacpSystemConvertSystemIdToString()
				pcms.Interval = ConvertLaAggIntervalToLacpPeriod(p.AggAttached.Config.Interval)
				pcms.Enabled = p.PortEnabled

				// stats

				pcms.LacpInPkts = int64(p.LacpCounter.AggPortStatsLACPDUsRx)
				pcms.LacpOutPkts = int64(p.LacpCounter.AggPortStatsLACPDUsTx)
				pcms.LacpRxErrors = int64(p.LacpCounter.AggPortStatsIllegalRx)
				pcms.LacpTxErrors = 0
				pcms.LacpUnknownErrors = int64(p.LacpCounter.AggPortStatsUnknownRx)
				pcms.LacpErrors = int64(p.LacpCounter.AggPortStatsIllegalRx) + int64(p.LacpCounter.AggPortStatsUnknownRx)
				pcms.LacpInMissMatchPkts = int64(p.LacpCounter.AggPortStateMissMatchInfoRx)
				pcms.LampInPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsRx)
				pcms.LampInResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsRx)
				pcms.LampOutPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsTx)
				pcms.LampOutResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsTx)

				// debug
				pcms.DebugId = int32(p.AggPortDebug.AggPortDebugInformationID)
				pcms.RxMachine = ConvertRxMachineStateToYangState(p.AggPortDebug.AggPortDebugRxState)
				pcms.RxTime = int32(p.AggPortDebug.AggPortDebugLastRxTime)
				pcms.MuxMachine = ConvertMuxMachineStateToYangState(p.AggPortDebug.AggPortDebugMuxState)
				pcms.MuxReason = string(p.AggPortDebug.AggPortDebugMuxReason)
				pcms.ActorChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugActorChurnState)
				pcms.PartnerChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugPartnerChurnState)
				pcms.ActorChurnCount = int64(p.AggPortDebug.AggPortDebugActorChurnCount)
				pcms.PartnerChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerChurnCount)
				pcms.ActorSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugActorSyncTransitionCount)
				pcms.PartnerSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugPartnerSyncTransitionCount)
				pcms.ActorChangeCount = int64(p.AggPortDebug.AggPortDebugActorChangeCount)
				pcms.PartnerChangeCount = int64(p.AggPortDebug.AggPortDebugPartnerChangeCount)
				pcms.ActorCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugActorCDSChurnState)
				pcms.PartnerCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugPartnerCDSChurnState)
				pcms.ActorCdsChurnCount = int64(p.AggPortDebug.AggPortDebugActorCDSChurnCount)
				pcms.PartnerCdsChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerCDSChurnCount)
			*/
		} else {
			return pcms, errors.New(fmt.Sprintf("LACP: Unabled to find config port by IntfRef %s", intfref))
		}
	} else {
		if lacp.LaFindPortById(uint16(id), &p) {
			// actor info
			pcms.Aggregatable = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateAggregationBit)
			pcms.Collecting = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateCollectingBit)
			pcms.Distributing = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDistributingBit)
			pcms.Defaulted = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDefaultedBit)

			if pcms.Distributing {
				pcms.OperState = "UP"
			} else {
				pcms.OperState = "DOWN"
			}

			if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateSyncBit) {
				// in sync
				pcms.Synchronization = 0
			} else {
				// out of sync
				pcms.Synchronization = 1
			}
			// short 1, long 0
			if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateTimeoutBit) {
				// short
				pcms.Timeout = 1
			} else {
				// long
				pcms.Timeout = 0
			}

			if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateActivityBit) {
				// active
				pcms.Activity = 0
			} else {
				// passive
				pcms.Activity = 1
			}

			pcms.OperKey = int16(p.ActorOper.Key)
			pcms.IntfRef = p.IntfNum
			if p.AggAttached != nil {
				pcms.LagIntfRef = p.AggAttached.AggName
				//		pcms.Mode = ConvertLaAggModeToModelLacpMode(p.AggAttached.Config.Mode)
			}

			pcms.DrniName = p.DrniName
			pcms.DrniSynced = p.DrniSynced

			// partner info
			pcms.PartnerId = p.PartnerOper.System.LacpSystemConvertSystemIdToString()
			pcms.PartnerKey = int16(p.PartnerOper.Key)

			// System
			//pcms.SystemIdMac = p.ActorOper.System.LacpSystemConvertSystemIdToString()[6:]
			//pcms.LagType = ConvertLaAggTypeToModelLagType(p.AggAttached.AggType)
			pcms.SystemId = p.ActorOper.System.LacpSystemConvertSystemIdToString()
			//pcms.Interval = ConvertLaAggIntervalToLacpPeriod(p.AggAttached.Config.Interval)
			//pcms.Enabled = p.PortEnabled

			// stats
			pcms.LacpInPkts = int64(p.LacpCounter.AggPortStatsLACPDUsRx)
			pcms.LacpOutPkts = int64(p.LacpCounter.AggPortStatsLACPDUsTx)
			pcms.LacpRxErrors = int64(p.LacpCounter.AggPortStatsIllegalRx)
			pcms.LacpTxErrors = 0
			pcms.LacpUnknownErrors = int64(p.LacpCounter.AggPortStatsUnknownRx)
			pcms.LacpErrors = int64(p.LacpCounter.AggPortStatsIllegalRx) + int64(p.LacpCounter.AggPortStatsUnknownRx)
			pcms.LacpInMissMatchPkts = int64(p.LacpCounter.AggPortStateMissMatchInfoRx)
			pcms.LampInPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsRx)
			pcms.LampInResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsRx)
			pcms.LampOutPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsTx)
			pcms.LampOutResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsTx)

			// debug
			pcms.DebugId = int32(p.AggPortDebug.AggPortDebugInformationID)
			pcms.RxMachine = ConvertRxMachineStateToYangState(p.AggPortDebug.AggPortDebugRxState)
			pcms.RxTime = int32(p.AggPortDebug.AggPortDebugLastRxTime)
			pcms.MuxMachine = ConvertMuxMachineStateToYangState(p.AggPortDebug.AggPortDebugMuxState)
			pcms.MuxReason = string(p.AggPortDebug.AggPortDebugMuxReason)
			pcms.ActorChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugActorChurnState)
			pcms.PartnerChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugPartnerChurnState)
			pcms.ActorChurnCount = int64(p.AggPortDebug.AggPortDebugActorChurnCount)
			pcms.PartnerChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerChurnCount)
			pcms.ActorSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugActorSyncTransitionCount)
			pcms.PartnerSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugPartnerSyncTransitionCount)
			pcms.ActorChangeCount = int64(p.AggPortDebug.AggPortDebugActorChangeCount)
			pcms.PartnerChangeCount = int64(p.AggPortDebug.AggPortDebugPartnerChangeCount)
			pcms.ActorCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugActorCDSChurnState)
			pcms.PartnerCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugPartnerCDSChurnState)
			pcms.ActorCdsChurnCount = int64(p.AggPortDebug.AggPortDebugActorCDSChurnCount)
			pcms.PartnerCdsChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerCDSChurnCount)
		} else {
			return pcms, errors.New(fmt.Sprintf("LACP: Unabled to find port by IntfRef %s", intfref))
		}
	}
	return pcms, nil
}

// GetBulkAggregationLacpMemberStateCounters will return the status of all
// the lag members.
func (la *LACPDServiceHandler) GetBulkLaPortChannelIntfRefListState(fromIndex lacpd.Int, count lacpd.Int) (obj *lacpd.LaPortChannelIntfRefListStateGetInfo, err error) {

	var lagMemberStateList []lacpd.LaPortChannelIntfRefListState = make([]lacpd.LaPortChannelIntfRefListState, count)
	var nextLagMemberState *lacpd.LaPortChannelIntfRefListState
	var returnLagMemberStates []*lacpd.LaPortChannelIntfRefListState
	var returnLagMemberStateGetInfo lacpd.LaPortChannelIntfRefListStateGetInfo
	var p *lacp.LaAggPort
	var ac *lacp.LaAggConfig
	validCount := lacpd.Int(0)
	toIndex := fromIndex
	moreRoutes := false

	aggTimeoutToIntervalMap := map[time.Duration]int32{
		//LacpPeriodTypeSLOW: "LONG",
		//LacpPeriodTypeFAST: "SHORT",
		//lacp.LacpSlowPeriodicTime: "LONG",
		//lacp.LacpFastPeriodicTime: "SHORT",
		lacp.LacpLongTimeoutTime:  0,
		lacp.LacpShortTimeoutTime: 1,
	}

	obj = &returnLagMemberStateGetInfo

	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_DISABLE {
		lagIndex := lacpd.Int(0)
		currIndex := lacpd.Int(0)
		for lagIndex = 0; validCount != count && lacp.LaAggConfigGetByIndex(int(lagIndex), &ac); lagIndex++ {
			for _, ifindex := range ac.LagMembers {
				intfref := utils.GetNameFromIfIndex(int32(ifindex))
				if lacp.LaAggConfigDoesIntfRefListMemberExist(intfref, &ac) {
					if currIndex < fromIndex {
						currIndex++
						continue
					} else {

						nextLagMemberState = &lagMemberStateList[validCount]
						nextLagMemberState.Aggregatable = true
						nextLagMemberState.Collecting = false
						nextLagMemberState.Distributing = false
						nextLagMemberState.Defaulted = true
						nextLagMemberState.OperState = "DOWN"

						// out of sync
						nextLagMemberState.Synchronization = 1
						timeout, ok := aggTimeoutToIntervalMap[ac.Lacp.Interval]
						if !ok {
							timeout = 0
						}

						nextLagMemberState.Timeout = timeout

						nextLagMemberState.Activity = ConvertLaAggModeToModelLacpMode(ac.Lacp.Mode)

						nextLagMemberState.OperKey = int16(0)
						nextLagMemberState.IntfRef = intfref
						nextLagMemberState.LagIntfRef = ac.Name

						nextLagMemberState.DrniName = ""
						nextLagMemberState.DrniSynced = false

						// partner info
						//pcms.PartnerId = p.PartnerOper.System.LacpSystemConvertSystemIdToString()
						//pcms.PartnerKey = int16(p.PartnerOper.Key)

						// System
						//pcms.SystemIdMac = p.ActorOper.System.LacpSystemConvertSystemIdToString()[6:]
						//pcms.LagType = ConvertLaAggTypeToModelLagType(p.AggAttached.AggType)
						//pcms.SystemId = p.ActorOper.System.LacpSystemConvertSystemIdToString()
						//pcms.Interval = ConvertLaAggIntervalToLacpPeriod(p.AggAttached.Config.Interval)
						//pcms.Enabled = p.PortEnabled

						// stats
						/*
							pcms.LacpInPkts = int64(p.LacpCounter.AggPortStatsLACPDUsRx)
							pcms.LacpOutPkts = int64(p.LacpCounter.AggPortStatsLACPDUsTx)
							pcms.LacpRxErrors = int64(p.LacpCounter.AggPortStatsIllegalRx)
							pcms.LacpTxErrors = 0
							pcms.LacpUnknownErrors = int64(p.LacpCounter.AggPortStatsUnknownRx)
							pcms.LacpErrors = int64(p.LacpCounter.AggPortStatsIllegalRx) + int64(p.LacpCounter.AggPortStatsUnknownRx)
							pcms.LacpInMissMatchPkts = int64(p.LacpCounter.AggPortStateMissMatchInfoRx)
							pcms.LampInPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsRx)
							pcms.LampInResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsRx)
							pcms.LampOutPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsTx)
							pcms.LampOutResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsTx)

							// debug
							pcms.DebugId = int32(p.AggPortDebug.AggPortDebugInformationID)
							pcms.RxMachine = ConvertRxMachineStateToYangState(p.AggPortDebug.AggPortDebugRxState)
							pcms.RxTime = int32(p.AggPortDebug.AggPortDebugLastRxTime)
							pcms.MuxMachine = ConvertMuxMachineStateToYangState(p.AggPortDebug.AggPortDebugMuxState)
							pcms.MuxReason = string(p.AggPortDebug.AggPortDebugMuxReason)
							pcms.ActorChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugActorChurnState)
							pcms.PartnerChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugPartnerChurnState)
							pcms.ActorChurnCount = int64(p.AggPortDebug.AggPortDebugActorChurnCount)
							pcms.PartnerChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerChurnCount)
							pcms.ActorSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugActorSyncTransitionCount)
							pcms.PartnerSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugPartnerSyncTransitionCount)
							pcms.ActorChangeCount = int64(p.AggPortDebug.AggPortDebugActorChangeCount)
							pcms.PartnerChangeCount = int64(p.AggPortDebug.AggPortDebugPartnerChangeCount)
							pcms.ActorCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugActorCDSChurnState)
							pcms.PartnerCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugPartnerCDSChurnState)
							pcms.ActorCdsChurnCount = int64(p.AggPortDebug.AggPortDebugActorCDSChurnCount)
							pcms.PartnerCdsChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerCDSChurnCount)
						*/
					}
					if len(returnLagMemberStates) == 0 {
						returnLagMemberStates = make([]*lacpd.LaPortChannelIntfRefListState, 0)
					}
					returnLagMemberStates = append(returnLagMemberStates, nextLagMemberState)
					validCount++
					toIndex++
					currIndex++
				}
			}
		}
		if ac != nil {
			moreRoutes = lacp.LaAggConfigGetByIndex(int(currIndex), &ac)
		}

	} else {
		for currIndex := lacpd.Int(0); validCount != count && lacp.LaGetPortNext(&p); currIndex++ {

			if currIndex < fromIndex {
				continue
			} else {

				nextLagMemberState = &lagMemberStateList[validCount]

				// actor info
				nextLagMemberState.Aggregatable = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateAggregationBit)
				nextLagMemberState.Collecting = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateCollectingBit)
				nextLagMemberState.Distributing = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDistributingBit)
				nextLagMemberState.Defaulted = lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateDefaultedBit)

				if nextLagMemberState.Distributing {
					nextLagMemberState.OperState = "UP"
				} else {
					nextLagMemberState.OperState = "DOWN"
				}

				if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateSyncBit) {
					// in sync
					nextLagMemberState.Synchronization = 1
				} else {
					// out of sync
					nextLagMemberState.Synchronization = 0
				}
				// short 1, long 0
				if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateTimeoutBit) {
					// short
					nextLagMemberState.Timeout = 1
				} else {
					// long
					nextLagMemberState.Timeout = 0
				}

				if lacp.LacpStateIsSet(p.ActorOper.State, lacp.LacpStateActivityBit) {
					// active
					nextLagMemberState.Activity = 0
				} else {
					// passive
					nextLagMemberState.Activity = 1
				}

				nextLagMemberState.OperKey = int16(p.ActorOper.Key)
				nextLagMemberState.IntfRef = p.IntfNum
				nextLagMemberState.IfIndex = utils.GetIfIndexFromName(p.IntfNum)

				if p.AggAttached != nil {
					nextLagMemberState.LagIntfRef = p.AggAttached.AggName
					//		nextLagMemberState.Mode = ConvertLaAggModeToModelLacpMode(p.AggAttached.Config.Mode)
				}

				// partner info
				nextLagMemberState.PartnerId = p.PartnerOper.System.LacpSystemConvertSystemIdToString()
				nextLagMemberState.PartnerKey = int16(p.PartnerOper.Key)

				// System
				//nextLagMemberState.SystemIdMac = p.ActorOper.System.LacpSystemConvertSystemIdToString()[6:]
				//nextLagMemberState.LagType = ConvertLaAggTypeToModelLagType(p.AggAttached.AggType)
				nextLagMemberState.SystemId = p.ActorOper.System.LacpSystemConvertSystemIdToString()
				//nextLagMemberState.Interval = ConvertLaAggIntervalToLacpPeriod(p.AggAttached.Config.Interval)
				//nextLagMemberState.Enabled = p.PortEnabled

				// stats
				nextLagMemberState.LacpInPkts = int64(p.LacpCounter.AggPortStatsLACPDUsRx)
				nextLagMemberState.LacpOutPkts = int64(p.LacpCounter.AggPortStatsLACPDUsTx)
				nextLagMemberState.LacpRxErrors = int64(p.LacpCounter.AggPortStatsIllegalRx)
				nextLagMemberState.LacpTxErrors = 0
				nextLagMemberState.LacpUnknownErrors = int64(p.LacpCounter.AggPortStatsUnknownRx)
				nextLagMemberState.LacpErrors = int64(p.LacpCounter.AggPortStatsIllegalRx) + int64(p.LacpCounter.AggPortStatsUnknownRx)
				nextLagMemberState.LampInPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsRx)
				nextLagMemberState.LampInResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsRx)
				nextLagMemberState.LampOutPdu = int64(p.LacpCounter.AggPortStatsMarkerPDUsTx)
				nextLagMemberState.LampOutResponsePdu = int64(p.LacpCounter.AggPortStatsMarkerResponsePDUsTx)

				// debug
				nextLagMemberState.DebugId = int32(p.AggPortDebug.AggPortDebugInformationID)
				nextLagMemberState.RxMachine = ConvertRxMachineStateToYangState(p.AggPortDebug.AggPortDebugRxState)
				nextLagMemberState.RxTime = int32(p.AggPortDebug.AggPortDebugLastRxTime)
				nextLagMemberState.MuxMachine = ConvertMuxMachineStateToYangState(p.AggPortDebug.AggPortDebugMuxState)
				nextLagMemberState.MuxReason = string(p.AggPortDebug.AggPortDebugMuxReason)
				nextLagMemberState.ActorChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugActorChurnState)
				nextLagMemberState.PartnerChurnMachine = ConvertCdmMachineStateToYangState(p.AggPortDebug.AggPortDebugPartnerChurnState)
				nextLagMemberState.ActorChurnCount = int64(p.AggPortDebug.AggPortDebugActorChurnCount)
				nextLagMemberState.PartnerChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerChurnCount)
				nextLagMemberState.ActorSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugActorSyncTransitionCount)
				nextLagMemberState.PartnerSyncTransitionCount = int64(p.AggPortDebug.AggPortDebugPartnerSyncTransitionCount)
				nextLagMemberState.ActorChangeCount = int64(p.AggPortDebug.AggPortDebugActorChangeCount)
				nextLagMemberState.PartnerChangeCount = int64(p.AggPortDebug.AggPortDebugPartnerChangeCount)
				nextLagMemberState.ActorCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugActorCDSChurnState)
				nextLagMemberState.PartnerCdsChurnMachine = int32(p.AggPortDebug.AggPortDebugPartnerCDSChurnState)
				nextLagMemberState.ActorCdsChurnCount = int64(p.AggPortDebug.AggPortDebugActorCDSChurnCount)
				nextLagMemberState.PartnerCdsChurnCount = int64(p.AggPortDebug.AggPortDebugPartnerCDSChurnCount)

				if len(returnLagMemberStates) == 0 {
					returnLagMemberStates = make([]*lacpd.LaPortChannelIntfRefListState, 0)
				}
				returnLagMemberStates = append(returnLagMemberStates, nextLagMemberState)
				validCount++
				toIndex++
			}
		}

		// lets try and get the next agg if one exists then there are more routes
		if p != nil {
			moreRoutes = lacp.LaGetPortNext(&p)
		}
	}

	fmt.Printf("Returning %d list of lagMembers\n", validCount)
	obj.LaPortChannelIntfRefListStateList = returnLagMemberStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount

	return obj, nil
}

func (la *LACPDServiceHandler) convertDbObjDataToDRCPData(objData *objects.DistributedRelay, cfgData *drcp.DistributedRelayConfig) {
	cfgData.DrniName = objData.DrniName
	cfgData.DrniPortalAddress = objData.PortalAddress
	cfgData.DrniPortalPriority = uint16(objData.PortalPriority)
	cfgData.DrniThreePortalSystem = objData.ThreePortalSystem
	cfgData.DrniPortalSystemNumber = objData.PortalSystemNumber
	for i, val := range objData.IntfReflist {
		cfgData.DrniIntraPortalLinkList[i] = uint32(utils.GetIfIndexFromName(val))
	}
	cfgData.DrniAggregator = uint32(GetKeyByAggName(objData.IntfRef))

	cfgData.DrniGatewayAlgorithm = objData.GatewayAlgorithm
	cfgData.DrniNeighborAdminGatewayAlgorithm = objData.NeighborGatewayAlgorithm
	cfgData.DrniNeighborAdminPortAlgorithm = objData.NeighborPortAlgorithm
	cfgData.DrniNeighborAdminDRCPState = objData.NeighborAdminDRCPState
	cfgData.DrniEncapMethod = objData.EncapMethod
	cfgData.DrniIntraPortalPortProtocolDA = objData.IntraPortalPortProtocolDA
}

func (la *LACPDServiceHandler) CreateDistributedRelay(config *lacpd.DistributedRelay) (bool, error) {

	data := &objects.DistributedRelay{}
	objects.ConvertThriftTolacpdDistributedRelayObj(config, data)

	conf := &drcp.DistributedRelayConfig{}
	// convert to drcp module config data
	la.convertDbObjDataToDRCPData(data, conf)
	err1 := drcp.DistributedRelayConfigCreateCheck(conf.DrniName, conf.DrniAggregator)
	err2 := drcp.DistributedRelayConfigParamCheck(conf)
	if err1 != nil {
		return false, err1
	} else if err2 != nil {
		return false, err2
	} else {
		if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {

			cfg := server.LAConfig{
				Msgtype: server.LAConfigMsgCreateDistributedRelay,
				Msgdata: conf,
			}
			la.svr.ConfigCh <- cfg
		}
	}

	return true, nil
}
func (la *LACPDServiceHandler) UpdateDistributedRelay(origconfig *lacpd.DistributedRelay, updateconfig *lacpd.DistributedRelay, attrset []bool, op []*lacpd.PatchOpInfo) (bool, error) {

	objTyp := reflect.TypeOf(*origconfig)
	olddata := &objects.DistributedRelay{}
	objects.ConvertThriftTolacpdDistributedRelayObj(origconfig, olddata)

	newdata := &objects.DistributedRelay{}
	objects.ConvertThriftTolacpdDistributedRelayObj(updateconfig, newdata)

	oldconf := &drcp.DistributedRelayConfig{}
	// convert to drcp module config data
	la.convertDbObjDataToDRCPData(olddata, oldconf)
	newconf := &drcp.DistributedRelayConfig{}
	// convert to drcp module config data
	la.convertDbObjDataToDRCPData(newdata, newconf)
	err1 := drcp.DistributedRelayConfigCreateCheck(newconf.DrniName, newconf.DrniAggregator)
	err2 := drcp.DistributedRelayConfigParamCheck(newconf)
	if err1 != nil {
		return false, err1
	} else if err2 != nil {
		return false, err2
	} else {
		if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {
			// TODO need to set valid attribute types for update
			attrMap := map[string]server.LaConfigMsgType{}
			// lets deal with Members attribute first
			for i := 0; i < objTyp.NumField(); i++ {
				objName := objTyp.Field(i).Name
				//fmt.Println("UpdateDistributedRelay (server): (index, objName) ", i, objName)
				if attrset[i] {
					fmt.Println("UpdateDistributedRelay (server): changed ", objName)

					if msgtype, ok := attrMap[objName]; ok {
						// set message type
						cfg := server.LAConfig{
							Msgdata: newconf,
							Msgtype: msgtype,
						}
						la.svr.ConfigCh <- cfg
					}
				}
			}
		}
	}
	return true, nil
}
func (la *LACPDServiceHandler) DeleteDistributedRelay(config *lacpd.DistributedRelay) (bool, error) {

	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {
		data := &objects.DistributedRelay{}
		objects.ConvertThriftTolacpdDistributedRelayObj(config, data)

		conf := &drcp.DistributedRelayConfig{}
		// convert to drcp module config data
		la.convertDbObjDataToDRCPData(data, conf)
		err := drcp.DistributedRelayConfigDeleteCheck(config.DrniName)
		if err != nil {
			return false, err
		} else {
			cfg := server.LAConfig{
				Msgtype: server.LAConfigMsgDeleteDistributedRelay,
				Msgdata: conf,
			}
			la.svr.ConfigCh <- cfg
		}
	}
	return true, nil
}

func (la *LACPDServiceHandler) GetDistributedRelayState(drname string) (*lacpd.DistributedRelayState, error) {

	drs := &lacpd.DistributedRelayState{}
	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {
		var dr *drcp.DistributedRelay
		if drcp.DrFindByName(drname, &dr) {

			var a *lacp.LaAggregator
			aggName := ""
			if lacp.LaFindAggById(int(dr.DrniAggregator), &a) {
				aggName = a.AggName
			}

			drs.DrniName = dr.DrniName
			drs.IntfRef = aggName
			drs.PortalAddress = dr.DrniPortalAddr.String()
			drs.PortalPriority = int16(dr.DrniPortalPriority)
			drs.ThreePortalSystem = dr.DrniThreeSystemPortal
			drs.PortalSystemNumber = int8(dr.DrniPortalSystemNumber)
			for _, ifindex := range dr.DrniIntraPortalLinkList {
				if ifindex != 0 {
					drs.IntfRefList = append(drs.IntfRefList, utils.GetNameFromIfIndex(int32(ifindex)))
				}
			}
			drs.GatewayAlgorithm = dr.DrniGatewayAlgorithm.String()
			drs.NeighborGatewayAlgorithm = dr.DrniNeighborAdminGatewayAlgorithm.String()
			drs.NeighborPortAlgorithm = dr.DrniNeighborAdminPortAlgorithm.String()
			drs.NeighborAdminDRCPState = fmt.Sprintf("%s", strconv.FormatInt(int64(dr.DrniNeighborAdminDRCPState), 2))
			drs.EncapMethod = dr.DrniEncapMethod.String()
			drs.PSI = dr.DrniPSI
			drs.IntraPortalPortProtocolDA = dr.DrniPortalPortProtocolIDA.String()

		}
	}
	return drs, nil
}

func (la *LACPDServiceHandler) GetBulkDistributedRelayState(fromIndex lacpd.Int, count lacpd.Int) (obj *lacpd.DistributedRelayStateGetInfo, err error) {
	var drcpStateList []lacpd.DistributedRelayState = make([]lacpd.DistributedRelayState, count)
	var nextDrcpState *lacpd.DistributedRelayState
	var returnDrcpStates []*lacpd.DistributedRelayState
	var returnDrcpStateGetInfo lacpd.DistributedRelayStateGetInfo
	var dr *drcp.DistributedRelay
	validCount := lacpd.Int(0)
	toIndex := fromIndex
	obj = &returnDrcpStateGetInfo

	for currIndex := lacpd.Int(0); validCount != count && drcp.DrGetDrcpNext(&dr); currIndex++ {

		if currIndex < fromIndex {
			continue
		} else {

			var a *lacp.LaAggregator
			aggName := ""
			if lacp.LaFindAggById(int(dr.DrniAggregator), &a) {
				aggName = a.AggName
			}

			nextDrcpState = &drcpStateList[validCount]
			nextDrcpState.DrniName = dr.DrniName
			nextDrcpState.IntfRef = aggName
			nextDrcpState.PortalAddress = dr.DrniPortalAddr.String()
			nextDrcpState.PortalPriority = int16(dr.DrniPortalPriority)
			nextDrcpState.ThreePortalSystem = dr.DrniThreeSystemPortal
			nextDrcpState.PortalSystemNumber = int8(dr.DrniPortalSystemNumber)
			utils.GlobalLogger.Info(fmt.Sprintf("GetBulkDistributedRelay, IPP Port List %v", dr.DrniIntraPortalLinkList))
			for _, ifindex := range dr.DrniIntraPortalLinkList {
				if ifindex != 0 {
					nextDrcpState.IntfRefList = append(nextDrcpState.IntfRefList, utils.GetNameFromIfIndex(int32(ifindex)))
				}
			}
			nextDrcpState.GatewayAlgorithm = dr.DrniGatewayAlgorithm.String()
			nextDrcpState.NeighborGatewayAlgorithm = dr.DrniNeighborAdminGatewayAlgorithm.String()
			nextDrcpState.NeighborPortAlgorithm = dr.DrniNeighborAdminPortAlgorithm.String()
			nextDrcpState.NeighborAdminDRCPState = fmt.Sprintf("%s", strconv.FormatInt(int64(dr.DrniNeighborAdminDRCPState), 2))
			nextDrcpState.EncapMethod = dr.DrniEncapMethod.String()
			nextDrcpState.PSI = dr.DrniPSI
			nextDrcpState.IntraPortalPortProtocolDA = dr.DrniPortalPortProtocolIDA.String()

			if len(returnDrcpStates) == 0 {
				returnDrcpStates = make([]*lacpd.DistributedRelayState, 0)
			}
			returnDrcpStates = append(returnDrcpStates, nextDrcpState)
			validCount++
			toIndex++
		}
	}
	// lets try and get the next agg if one exists then there are more routes
	moreRoutes := false
	if dr != nil {
		moreRoutes = drcp.DrGetDrcpNext(&dr)
	}

	obj.DistributedRelayStateList = returnDrcpStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount

	return obj, err
}

func (la *LACPDServiceHandler) GetIppLinkState(intref, drnameref string) (obj *lacpd.IppLinkState, err error) {
	return obj, err
}

func (la *LACPDServiceHandler) GetBulkIppLinkState(fromIndex lacpd.Int, count lacpd.Int) (obj *lacpd.IppLinkStateGetInfo, err error) {
	// TODO
	return obj, err
}

func (la *LACPDServiceHandler) GetLacpGlobalState(vrf string) (*lacpd.LacpGlobalState, error) {
	obj := &lacpd.LacpGlobalState{}
	obj.Vrf = "default"
	obj.AdminState = "UP"
	if utils.LacpGlobalStateGet() != utils.LACP_GLOBAL_ENABLE {
		obj.AdminState = "DOWN"
	}
	var a *lacp.LaAggregator
	if utils.LacpGlobalStateGet() == utils.LACP_GLOBAL_ENABLE {
		for lacp.LaGetAggNext(&a) {
			obj.AggList = append(obj.AggList, a.AggName)
			if a.OperState {
				obj.AggOperStateUpList = append(obj.AggOperStateUpList, a.AggName)
				if a.DrniName != "" {
					obj.DistributedRelayUpList = append(obj.DistributedRelayUpList, a.DrniName)
				}
			}
			if a.DrniName != "" {
				obj.DistributedRelayAttachedList = append(obj.DistributedRelayAttachedList, fmt.Sprintf("%s-%s", a.DrniName, a.AggName))
			}
		}
		var dr *drcp.DistributedRelay
		for drcp.DrGetDrcpNext(&dr) {
			obj.DistributedRelayList = append(obj.DistributedRelayList, dr.DrniName)
		}

		var p *lacp.LaAggPort
		for lacp.LaGetPortNext(&p) {
			obj.LacpErrorsInPkts += int64(p.LacpCounter.AggPortStatsIllegalRx) + int64(p.LacpCounter.AggPortStatsUnknownRx)
			obj.LacpMissMatchPkts += int64(p.LacpCounter.AggPortStateMissMatchInfoRx)
			obj.LacpTotalRxPkts += int64(p.LacpCounter.AggPortStatsLACPDUsRx)
			obj.LacpTotalTxPkts += int64(p.LacpCounter.AggPortStatsLACPDUsTx)
		}
	} else {
		var currIndex lacpd.Int
		var ac *lacp.LaAggConfig
		for currIndex = 0; lacp.LaAggConfigGetByIndex(int(currIndex), &ac); currIndex++ {
			obj.AggList = append(obj.AggList, ac.Name)
			// TODO need to add DRCP info
		}
	}

	return obj, nil
}

func (la *LACPDServiceHandler) GetBulkLacpGlobalState(fromIndex lacpd.Int, count lacpd.Int) (obj *lacpd.LacpGlobalStateGetInfo, err error) {
	var returnLacpGlobalStates []*lacpd.LacpGlobalState
	var returnLacpGlobalStateGetInfo lacpd.LacpGlobalStateGetInfo
	toIndex := fromIndex
	obj = &returnLacpGlobalStateGetInfo

	nextLacpGlobalState, gserr := la.GetLacpGlobalState("default")
	if gserr == nil {
		if len(returnLacpGlobalStates) == 0 {
			returnLacpGlobalStates = make([]*lacpd.LacpGlobalState, 0)
		}
		returnLacpGlobalStates = append(returnLacpGlobalStates, nextLacpGlobalState)
	}
	obj.LacpGlobalStateList = returnLacpGlobalStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = false
	obj.Count = 1
	return obj, err
}
