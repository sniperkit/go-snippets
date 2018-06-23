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

package ipcutils

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"models/objects"
	"sync"
	"utils/dbutils"
)

type IPCClientBase struct {
	Name               string
	Address            string
	TTransport         thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	Enabled            bool
	IsConnected        bool
	ApiHandlerMutex    sync.RWMutex
}

func (clnt *IPCClientBase) IsConnectedToServer() bool {
	return clnt.IsConnected
}

func (clnt *IPCClientBase) GetBulkObject(obj objects.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []objects.ConfigObj) {
	//logger.Println("### Get Bulk request called with", currMarker, count)
	return nil, 0, 0, false, make([]objects.ConfigObj, 0)
}

//
// This method gets Thrift related IPC handles.
//
func CreateIPCHandles(address string) (thrift.TTransport, *thrift.TBinaryProtocolFactory, error) {
	var transportFactory thrift.TTransportFactory
	var ttransport thrift.TTransport
	var protocolFactory *thrift.TBinaryProtocolFactory
	var err error

	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	ttransport, err = thrift.NewTSocket(address)
	if err != nil {
		return nil, nil, err
	}
	ttransport = transportFactory.GetTransport(ttransport)
	if err = ttransport.Open(); err != nil {
		//logger.Println("Failed to Open Transport", transport, protocolFactory)
		return nil, nil, err
	}
	return ttransport, protocolFactory, err
}

func (clnt *IPCClientBase) CloseIPCHandles() error {
	clnt.PtrProtocolFactory = nil
	if err := clnt.TTransport.Close(); err != nil {
		return err
	}
	clnt.IsConnected = false
	return nil
}

func (clnt *IPCClientBase) PreUpdateValidation(dbObj, obj objects.ConfigObj, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	return nil
}

func (clnt *IPCClientBase) PostUpdateProcessing(dbObj, obj objects.ConfigObj, attSet []bool, dbHdl *dbutils.DBUtil) error {
	return nil
}
