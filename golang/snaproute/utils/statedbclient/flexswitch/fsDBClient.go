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

package flexswitch

import (
	//	"fmt"
	"models/objects"
	"utils/dbutils"
	"utils/logging"
)

const (
	objAdd uint8 = iota
	objUpdate
	objDelete
)

var objOperation = map[uint8]string{
	objAdd:    "add",
	objUpdate: "update",
	objDelete: "delete",
}

type objInfo struct {
	operation uint8
	obj       objects.ConfigObj
}

type FSDBClient struct {
	logger     *logging.Writer
	dbUtil     *dbutils.DBUtil
	objStateCh chan objInfo
}

func NewFSDBClient(logger *logging.Writer) *FSDBClient {
	return &FSDBClient{
		logger:     logger,
		dbUtil:     dbutils.NewDBUtil(logger),
		objStateCh: make(chan objInfo, 30000),
	}
}

func (fs *FSDBClient) Init() error {
	err := fs.dbUtil.Connect()
	if err != nil {
		fs.logger.Err("FSDBClient - DB connect failed with error ", err)
		return err
	}

	go fs.StartStateObjectReceiver()
	return nil
}

func (fs *FSDBClient) AddObject(obj objects.ConfigObj) error {
	fs.logger.Info("AddObject object %s", obj.GetKey())
	fs.objStateCh <- objInfo{objAdd, obj}
	return nil
}

func (fs *FSDBClient) DeleteObject(obj objects.ConfigObj) error {
	fs.logger.Info("DeleteObject object %s", obj.GetKey())
	fs.objStateCh <- objInfo{objDelete, obj}
	return nil
}

func (fs *FSDBClient) UpdateObject(obj objects.ConfigObj) error {
	fs.logger.Info("UpdateObject object %s", obj.GetKey())
	fs.objStateCh <- objInfo{objUpdate, obj}
	return nil
}

/* This is done synchronously as we delete all the objects in the state DB when a process comes up */
func (fs *FSDBClient) DeleteAllObjects(obj objects.ConfigObj) error {
	fs.logger.Info("DeleteAllObjects object %s", obj.GetKey())
	objs, err := fs.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		fs.logger.Err("DeleteAllObjects - GetAllObjFromDb failed with error %s", err)
		return err
	}

	for idx, _ := range objs {
		fs.delObjToDB(objs[idx])
	}
	return nil
}

func (fs *FSDBClient) addObjToDB(obj objects.ConfigObj) (err error) {
	//fs.logger.Info("addObjToDB object %s", obj.GetKey())
	err = fs.dbUtil.StoreObjectInDb(obj)
	if err != nil {
		fs.logger.Err("Failed to add state object %s to DB with error %s", obj.GetKey(), err)
		return err
	}
	//fs.logger.Info("Added state object %s to DB", obj.GetKey())
	return nil
}

func (fs *FSDBClient) delObjToDB(obj objects.ConfigObj) (err error) {
	//fs.logger.Info("delObjToDB object %s", obj.GetKey())
	err = fs.dbUtil.DeleteObjectFromDb(obj)
	if err != nil {
		fs.logger.Err("Failed to delete state object %s from DB with error", obj.GetKey(), err)
		return err
	}
	fs.logger.Info("Deleted state object %s from DB", obj.GetKey())
	return nil
}

func (fs *FSDBClient) StartStateObjectReceiver() {
	fs.logger.Info("Starting the state object receiver")
	var err error

	for {
		err = nil
		select {
		case info := <-fs.objStateCh:
			if info.operation == objAdd {
				err = fs.addObjToDB(info.obj)
			} else if info.operation == objDelete {
				err = fs.delObjToDB(info.obj)
			} else if info.operation == objUpdate {
				//err = fs.updObjToDB(info.obj)
				err = fs.delObjToDB(info.obj)
				if err == nil {
					err = fs.addObjToDB(info.obj)
				}
			} else {
				fs.logger.Err("Recieved unknown operation %d for state object %s", info.operation,
					info.obj.GetKey())
			}

			if err != nil {
				fs.logger.Err("Failed to %s state object %s", objOperation[info.operation],
					info.obj.GetKey())
			}
		}
	}
}
