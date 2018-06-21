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

package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"ospfv2d"
	"utils/dbutils"
	"utils/logging"
)

type rpcServiceHandler struct {
	logger logging.LoggerIntf
	dbHdl  dbutils.DBIntf
}

type RPCServer struct {
	*thrift.TSimpleServer
}

func newRPCServiceHandler(logger logging.LoggerIntf, dbHdl dbutils.DBIntf) *rpcServiceHandler {
	hdl := &rpcServiceHandler{
		logger: logger,
		dbHdl:  dbHdl,
	}
	ok, err := hdl.restoreConfigFromDB()
	if !ok {
		logger.Err("Failed to restore configuration from DB-", err)
	}
	return hdl
}

func NewRPCServer(rpcAddr string, logger logging.LoggerIntf, dbHdl dbutils.DBIntf) *RPCServer {
	transport, err := thrift.NewTServerSocket(rpcAddr)
	if err != nil {
		panic(err)
	}
	handler := newRPCServiceHandler(logger, dbHdl)
	processor := ospfv2d.NewOSPFV2DServicesProcessor(handler)
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	return &RPCServer{
		TSimpleServer: server,
	}
}

func (rpcHdl *rpcServiceHandler) restoreConfigFromDB() (bool, error) {
	ok, err := rpcHdl.restoreOspfv2GlobalConfFromDB()
	if !ok {
		return ok, err
	}
	ok, err = rpcHdl.restoreOspfv2AreaConfFromDB()
	if !ok {
		return ok, err
	}
	ok, err = rpcHdl.restoreOspfv2IntfConfFromDB()
	return ok, err
}
