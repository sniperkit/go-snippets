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

// ribDB.go
package server

import (
	defs "l3/rib/ribdCommonDefs"
	"models/objects"
	"ribd"
	"utils/dbutils"
)

func (ribdServiceHandler *RIBDServer) UpdateRoutesFromDB() (err error) {
	logger.Debug("UpdateRoutesFromDB")
	ribdServiceHandler.DBRouteCh <- RIBdServerConfig{Op: defs.DBFetch}
	/*	dbHdl := ribdServiceHandler.DbHdl
		if dbHdl != nil {
			var dbObjCfg objects.IPv4Route
			objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
			if err == nil {
				logger.Debug(fmt.Sprintln("Number of routes from DB: ", len((objList))))
				for idx := 0; idx < len(objList); idx++ {
					obj := ribd.NewIPv4Route()
					dbObj := objList[idx].(objects.IPv4Route)
					objects.ConvertribdIPv4RouteObjToThrift(&dbObj, obj)
					err = ribdServiceHandler.RouteConfigValidationCheck(obj, "add")
					if err != nil {
						logger.Err("Route validation failed when reading from db")
						continue
					}
					rv, _ := ribdServiceHandler.ProcessRouteCreateConfig(obj)
					if rv == false {
						logger.Err("IPv4Route create failed during init")
					}
				}
			} else {
				logger.Err("DB Query failed during IPv4Route query: RIBd init")
			}
		}*/
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyPrefixSetsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyPrefixSetsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyPrefixSet
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyPrefixSet()
				dbObj := objList[idx].(objects.PolicyPrefixSet)
				objects.ConvertribdPolicyPrefixSetObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyPrefixSet,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyPrefixSet query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyASPathSetsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyASPathSetsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyASPathSet
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyASPathSet()
				dbObj := objList[idx].(objects.PolicyASPathSet)
				objects.ConvertribdPolicyASPathSetObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyASPathSet,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyASPathSet query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyCommunitySetsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyCommunitySetsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyCommunitySet
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyCommunitySet()
				dbObj := objList[idx].(objects.PolicyCommunitySet)
				objects.ConvertribdPolicyCommunitySetObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyCommunitySet,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyCommunitySet query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyExtendedCommunitySetsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyExtendedCommunitySetsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyExtendedCommunitySet
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyExtendedCommunitySet()
				dbObj := objList[idx].(objects.PolicyExtendedCommunitySet)
				objects.ConvertribdPolicyExtendedCommunitySetObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyExtendedCommunitySet,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyExtendedCommunitySet query: RIBd init")
		}
	}
	return err
}

func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyConditionsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyConditionsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyCondition
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyCondition()
				dbObj := objList[idx].(objects.PolicyCondition)
				objects.ConvertribdPolicyConditionObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyCondition,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyCondition query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyStmtsFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyStmtsFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyStmt
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyStmt()
				dbObj := objList[idx].(objects.PolicyStmt)
				objects.ConvertribdPolicyStmtObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyStmt,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyStmt query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdateGlobalPolicyFromDB(dbHdl *dbutils.DBUtil) (err error) {
	logger.Debug("UpdateGlobalPolicyFromDB")
	if dbHdl != nil {
		var dbObjCfg objects.PolicyDefinition
		objList, err := dbHdl.GetAllObjFromDb(dbObjCfg)
		if err == nil {
			for idx := 0; idx < len(objList); idx++ {
				obj := ribd.NewPolicyDefinition()
				dbObj := objList[idx].(objects.PolicyDefinition)
				objects.ConvertribdPolicyDefinitionObjToThrift(&dbObj, obj)
				ribdServiceHandler.PolicyConfCh <- RIBdServerConfig{
					OrigConfigObject: obj,
					Op:               defs.AddPolicyDefinition,
				}
				err = <-ribdServiceHandler.PolicyConfDone
			}
		} else {
			logger.Err("DB Query failed during PolicyDefinition query: RIBd init")
		}
	}
	return err
}
func (ribdServiceHandler *RIBDServer) UpdatePolicyObjectsFromDB() { //(paramsDir string) (err error) {
	logger.Debug("UpdateFromDB")
	dbHdl := ribdServiceHandler.DbHdl
	ribdServiceHandler.UpdateGlobalPolicyPrefixSetsFromDB(dbHdl)            //paramsDir, dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyASPathSetsFromDB(dbHdl)            //paramsDir, dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyCommunitySetsFromDB(dbHdl)         //paramsDir, dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyExtendedCommunitySetsFromDB(dbHdl) //paramsDir, dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyConditionsFromDB(dbHdl)            //paramsDir, dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyStmtsFromDB(dbHdl)
	ribdServiceHandler.UpdateGlobalPolicyFromDB(dbHdl)
	return
}
