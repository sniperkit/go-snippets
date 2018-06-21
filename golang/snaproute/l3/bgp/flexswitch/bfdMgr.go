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

package FSMgr

import (
	"bfdd"
	"encoding/json"
	"errors"
	_ "fmt"
	"l3/bfd/bfddCommonDefs"
	"l3/bgp/api"
	"l3/bgp/config"
	"l3/bgp/rpc"
	"utils/logging"

	nanomsg "github.com/op/go-nanomsg"
)

/*  Init bfd manager with bfd client as its core
 */
func NewFSBfdMgr(logger *logging.Writer, fileName string) (*FSBfdMgr, error) {
	var bfddClient *bfdd.BFDDServicesClient = nil
	bfddClientChan := make(chan *bfdd.BFDDServicesClient)

	logger.Info("Connecting to BFDd")
	go rpc.StartBfddClient(logger, fileName, bfddClientChan)
	bfddClient = <-bfddClientChan
	if bfddClient == nil {
		logger.Err("Failed to connect to BFDd")
		return nil, errors.New("Failed to connect to BFDd")
	} else {
		logger.Info("Connected to BFDd")
	}
	mgr := &FSBfdMgr{
		plugin:     "ovsdb",
		logger:     logger,
		bfddClient: bfddClient,
	}

	return mgr, nil
}

/*  Do any necessary init. Called from server..
 */
func (mgr *FSBfdMgr) Start() {
	// create bfd sub socket listener
	mgr.bfdSubSocket, _ = mgr.SetupSubSocket(bfddCommonDefs.PUB_SOCKET_ADDR)
	go mgr.listenForBFDNotifications()
}

/*  Listen for any BFD notifications
 */
func (mgr *FSBfdMgr) listenForBFDNotifications() {
	for {
		mgr.logger.Info("Read on BFD subscriber socket...")
		rxBuf, err := mgr.bfdSubSocket.Recv(0)
		if err != nil {
			mgr.logger.Err("Recv on BFD subscriber socket failed with error:", err)
			continue
		}
		mgr.logger.Info("BFD subscriber recv returned:", rxBuf)
		mgr.handleBfdNotifications(rxBuf)
	}
}

func (mgr *FSBfdMgr) handleBfdNotifications(rxBuf []byte) {
	bfd := bfddCommonDefs.BfddNotifyMsg{}
	err := json.Unmarshal(rxBuf, &bfd)
	if err != nil {
		mgr.logger.Errf("Unmarshal BFD notification failed with err %s", err)
		return
	}

	if bfd.State {
		api.SendBfdNotification(bfd.DestIp, bfd.State,
			config.BFD_STATE_VALID)
	} else {
		api.SendBfdNotification(bfd.DestIp, bfd.State,
			config.BFD_STATE_INVALID)
	}
}

func (mgr *FSBfdMgr) CreateBfdSession(ipAddr string, iface string, sessionParam string) (bool, error) {
	bfdSession := bfdd.NewBfdSession()
	bfdSession.IpAddr = ipAddr
	bfdSession.ParamName = sessionParam
	bfdSession.Interface = iface
	bfdSession.Owner = "bgp"
	mgr.logger.Info("Creating BFD Session: ", bfdSession)
	ret, err := mgr.bfddClient.CreateBfdSession(bfdSession)
	return ret, err
}

func (mgr *FSBfdMgr) DeleteBfdSession(ipAddr string, iface string) (bool, error) {
	bfdSession := bfdd.NewBfdSession()
	bfdSession.IpAddr = ipAddr
	bfdSession.Interface = iface
	bfdSession.Owner = "bgp"
	mgr.logger.Info("Deleting BFD Session: ", bfdSession)
	ret, err := mgr.bfddClient.DeleteBfdSession(bfdSession)
	return ret, err
}

func (mgr *FSBfdMgr) SetupSubSocket(address string) (*nanomsg.SubSocket, error) {
	var err error
	var socket *nanomsg.SubSocket
	if socket, err = nanomsg.NewSubSocket(); err != nil {
		mgr.logger.Errf("Failed to create subscribe socket %s, error:%s", address, err)
		return nil, err
	}

	if err = socket.Subscribe(""); err != nil {
		mgr.logger.Errf("Failed to subscribe to \"\" on subscribe socket %s, error:%s", address, err)
		return nil, err
	}

	if _, err = socket.Connect(address); err != nil {
		mgr.logger.Errf("Failed to connect to publisher socket %s, error:%s", address, err)
		return nil, err
	}

	mgr.logger.Infof("Connected to publisher socket %s", address)
	if err = socket.SetRecvBuffer(1024 * 1024); err != nil {
		mgr.logger.Err("Failed to set the buffer size for subsriber socket %s, error:", address, err)
		return nil, err
	}
	return socket, nil
}
