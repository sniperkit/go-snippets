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

package rpc

import (
	"bfdd"
	"github.com/garyburd/redigo/redis"
	"l3/bfd/server"
	"models/objects"
	"utils/logging"
)

type BFDHandler struct {
	server *server.BFDServer
	logger *logging.Writer
}

func NewBFDHandler(logger *logging.Writer, server *server.BFDServer) *BFDHandler {
	h := new(BFDHandler)
	h.server = server
	h.logger = logger
	return h
}

func (h *BFDHandler) ReadGlobalConfigFromDB(dbHdl redis.Conn) error {
	h.logger.Info("Reading BfdGlobal")
	if dbHdl != nil {
		var dbObj objects.BfdGlobal
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			h.logger.Err("DB query failed for global config")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := bfdd.NewBfdGlobal()
			dbObject := objList[idx].(objects.BfdGlobal)
			objects.ConvertbfddBfdGlobalObjToThrift(&dbObject, obj)
			rv, _ := h.CreateBfdGlobal(obj)
			if rv == false {
				h.logger.Err("BfdGlobal create failed")
			}
		}
	}
	return nil
}

func (h *BFDHandler) ReadSessionParamConfigFromDB(dbHdl redis.Conn) error {
	h.logger.Info("Reading BfdSessionParam")
	if dbHdl != nil {
		var dbObj objects.BfdSessionParam
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			h.logger.Err("DB query failed for session param config")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := bfdd.NewBfdSessionParam()
			dbObject := objList[idx].(objects.BfdSessionParam)
			objects.ConvertbfddBfdSessionParamObjToThrift(&dbObject, obj)
			rv, _ := h.CreateBfdSessionParam(obj)
			if rv == false {
				h.logger.Err("BfdSessionParam create failed for ", dbObject.Name)
			}
		}
	}
	return nil
}

func (h *BFDHandler) ReadSessionConfigFromDB(dbHdl redis.Conn) error {
	h.logger.Info("Reading BfdSession")
	if dbHdl != nil {
		var dbObj objects.BfdSession
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			h.logger.Err("DB query failed for session config")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := bfdd.NewBfdSession()
			dbObject := objList[idx].(objects.BfdSession)
			objects.ConvertbfddBfdSessionObjToThrift(&dbObject, obj)
			rv, _ := h.CreateBfdSession(obj)
			if rv == false {
				h.logger.Err("BfdSession create failed for ", dbObject.IpAddr)
			}
		}
	}
	return nil
}

func (h *BFDHandler) ReadConfigFromDB(dbHdl redis.Conn) error {
	// BfdGlobalConfig
	h.ReadGlobalConfigFromDB(dbHdl)
	// BfdIntfConfig
	h.ReadSessionParamConfigFromDB(dbHdl)
	// BfdSessionConfig
	h.ReadSessionConfigFromDB(dbHdl)
	return nil
}
