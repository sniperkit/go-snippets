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
//   This is a auto-generated file, please do not edit!
// _______   __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----  \   \/    \/   /  |  |  ---|  |----    ,---- |  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |        |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |        `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package rpc

import (
	"dhcprelayd"
	"l3/dhcp_relay/server"
)

func (rpcHdl *rpcServiceHandler) CreateDHCPRelayGlobal(
	cfg *dhcprelayd.DHCPRelayGlobal) (bool, error) {

	rpcHdl.logger.Info("Calling CreateDHCPRelayGlobal", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_DHCPRELAY_GLOBAL,
		Data: interface{}(&server.CreateDHCPRelayGlobalInArgs{
			DHCPRelayGlobal: cfg,
		}),
	}
	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) UpdateDHCPRelayGlobal(
	oldCfg, newCfg *dhcprelayd.DHCPRelayGlobal, attrset []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	rpcHdl.logger.Info("Calling UpdateDHCPRelayGlobal", oldCfg, newCfg,
		attrset, op)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_DHCPRELAY_GLOBAL,
		Data: interface{}(&server.UpdateDHCPRelayGlobalInArgs{
			DHCPRelayGlobalOld: oldCfg,
			DHCPRelayGlobalNew: newCfg,
			AttrSet:              attrset,
		}),
	}
	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl rpcServiceHandler) DeleteDHCPRelayGlobal(
	cfg *dhcprelayd.DHCPRelayGlobal) (bool, error) {

	rpcHdl.logger.Info("Calling DeleteDHCPRelayGlobal", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_DHCPRELAY_GLOBAL,
		Data: interface{}(&server.DeleteDHCPRelayGlobalInArgs{
			Vrf: cfg.Vrf,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) CreateDHCPRelayIntf(
	cfg *dhcprelayd.DHCPRelayIntf) (bool, error) {

	rpcHdl.logger.Info("Calling CreateDHCPRelayIntf", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_DHCPRELAY_INTF,
		Data: interface{}(&server.CreateDHCPRelayIntfInArgs{
			DHCPRelayIntf: cfg,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) UpdateDHCPRelayIntf(
	oldCfg, newCfg *dhcprelayd.DHCPRelayIntf, attrset []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	rpcHdl.logger.Info("Calling UpdateDHCPRelayIntf", oldCfg, newCfg,
		attrset, op)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_DHCPRELAY_INTF,
		Data: interface{}(&server.UpdateDHCPRelayIntfInArgs{
			DHCPRelayIntfOld: oldCfg,
			DHCPRelayIntfNew: newCfg,
			AttrSet:            attrset,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl rpcServiceHandler) DeleteDHCPRelayIntf(cfg *dhcprelayd.DHCPRelayIntf) (bool, error) {
	rpcHdl.logger.Info("Calling DeleteDHCPRelayIntf", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_DHCPRELAY_INTF,
		Data: interface{}(&server.DeleteDHCPRelayIntfInArgs{
			IntfRef: cfg.IntfRef,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPRelayClientState(key string) (obj *dhcprelayd.DHCPRelayClientState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPRelayClientState", key)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPRELAY_CLIENT_STATE,
		Data: interface{}(&server.GetDHCPRelayClientStateInArgs{
			MacAddr: key,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPRelayClientStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPRelayClientState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPRelayClientStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayClientStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayClientState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPRELAY_CLIENT_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPRelayClientStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPRelayIntfState(key string) (obj *dhcprelayd.DHCPRelayIntfState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPRelayIntfState", key)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPRELAY_INTF_STATE,
		Data: interface{}(&server.GetDHCPRelayIntfStateInArgs{
			IntfRef: key,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPRelayIntfStateOutArgs).Obj, result.Obj.(*server.GetDHCPRelayIntfStateOutArgs).Err
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPRelayIntfState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPRelayIntfStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayIntfStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayIntfState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format
	//return &getBulkInfo, err

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPRELAY_INTF_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPRelayIntfStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPRelayIntfServerState(key1 string, key2 string) (obj *dhcprelayd.DHCPRelayIntfServerState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPRelayIntfServerState", key1, key2)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPRELAY_INTFSERVER_STATE,
		Data: interface{}(&server.GetDHCPRelayIntfServerStateInArgs{
			IntfRef:  key1,
			ServerIp: key2,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPRelayIntfServerStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPRelayIntfServerState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPRelayIntfServerStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayIntfServerStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayIntfServerState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format
	//return &getBulkInfo, err

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPRELAY_INTFSERVER_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPRelayIntfServerStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) CreateDHCPv6RelayGlobal(
	cfg *dhcprelayd.DHCPv6RelayGlobal) (bool, error) {

	rpcHdl.logger.Info("Calling CreateDHCPv6RelayGlobal", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_DHCPV6RELAY_GLOBAL,
		Data: interface{}(&server.CreateDHCPv6RelayGlobalInArgs{
			DHCPv6RelayGlobal: cfg,
		}),
	}
	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) UpdateDHCPv6RelayGlobal(
	oldCfg, newCfg *dhcprelayd.DHCPv6RelayGlobal, attrset []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	rpcHdl.logger.Info("Calling UpdateDHCPv6RelayGlobal", oldCfg, newCfg,
		attrset, op)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_DHCPV6RELAY_GLOBAL,
		Data: interface{}(&server.UpdateDHCPv6RelayGlobalInArgs{
			DHCPv6RelayGlobalOld: oldCfg,
			DHCPv6RelayGlobalNew: newCfg,
			AttrSet:              attrset,
		}),
	}
	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl rpcServiceHandler) DeleteDHCPv6RelayGlobal(
	cfg *dhcprelayd.DHCPv6RelayGlobal) (bool, error) {

	rpcHdl.logger.Info("Calling DeleteDHCPv6RelayGlobal", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_DHCPV6RELAY_GLOBAL,
		Data: interface{}(&server.DeleteDHCPv6RelayGlobalInArgs{
			Vrf: cfg.Vrf,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) CreateDHCPv6RelayIntf(
	cfg *dhcprelayd.DHCPv6RelayIntf) (bool, error) {

	rpcHdl.logger.Info("Calling CreateDHCPv6RelayIntf", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_DHCPV6RELAY_INTF,
		Data: interface{}(&server.CreateDHCPv6RelayIntfInArgs{
			DHCPv6RelayIntf: cfg,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) UpdateDHCPv6RelayIntf(
	oldCfg, newCfg *dhcprelayd.DHCPv6RelayIntf, attrset []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	rpcHdl.logger.Info("Calling UpdateDHCPv6RelayIntf", oldCfg, newCfg,
		attrset, op)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_DHCPV6RELAY_INTF,
		Data: interface{}(&server.UpdateDHCPv6RelayIntfInArgs{
			DHCPv6RelayIntfOld: oldCfg,
			DHCPv6RelayIntfNew: newCfg,
			AttrSet:            attrset,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl rpcServiceHandler) DeleteDHCPv6RelayIntf(cfg *dhcprelayd.DHCPv6RelayIntf) (bool, error) {
	rpcHdl.logger.Info("Calling DeleteDHCPv6RelayIntf", cfg)

	//Send message to server
	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_DHCPV6RELAY_INTF,
		Data: interface{}(&server.DeleteDHCPv6RelayIntfInArgs{
			IntfRef: cfg.IntfRef,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(bool), result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPv6RelayClientState(key string) (obj *dhcprelayd.DHCPv6RelayClientState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPv6RelayClientState", key)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPV6RELAY_CLIENT_STATE,
		Data: interface{}(&server.GetDHCPv6RelayClientStateInArgs{
			MacAddr: key,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPv6RelayClientStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPv6RelayClientState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPv6RelayClientStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayClientStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayClientState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPV6RELAY_CLIENT_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPv6RelayClientStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPv6RelayIntfState(key string) (obj *dhcprelayd.DHCPv6RelayIntfState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPv6RelayIntfState", key)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPV6RELAY_INTF_STATE,
		Data: interface{}(&server.GetDHCPv6RelayIntfStateInArgs{
			IntfRef: key,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPv6RelayIntfStateOutArgs).Obj, result.Obj.(*server.GetDHCPv6RelayIntfStateOutArgs).Err
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPv6RelayIntfState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPv6RelayIntfStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayIntfStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayIntfState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format
	//return &getBulkInfo, err

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPV6RELAY_INTF_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPv6RelayIntfStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetDHCPv6RelayIntfServerState(key1 string, key2 string) (obj *dhcprelayd.DHCPv6RelayIntfServerState, err error) {
	rpcHdl.logger.Info("Calling GetDHCPv6RelayIntfServerState", key1, key2)

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GET_DHCPV6RELAY_INTFSERVER_STATE,
		Data: interface{}(&server.GetDHCPv6RelayIntfServerStateInArgs{
			IntfRef:  key1,
			ServerIp: key2,
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetDHCPv6RelayIntfServerStateOutArgs).Obj, result.Error
}

func (rpcHdl *rpcServiceHandler) GetBulkDHCPv6RelayIntfServerState(fromIdx, count dhcprelayd.Int) (*dhcprelayd.DHCPv6RelayIntfServerStateGetInfo, error) {
	//var getBulkInfo dhcprelayd.DHCPv6RelayIntfServerStateGetInfo
	//var err error
	//info, err := api.GetBulkDHCPv6RelayIntfServerState(int(fromIdx), int(count))
	//getBulkInfo.StartIdx = fromIdx
	//getBulkInfo.EndIdx = dhcprelayd.Int(info.EndIdx)
	//getBulkInfo.More = info.More
	//getBulkInfo.Count = dhcprelayd.Int(len(info.List))
	// Fill in data, remember to convert back to thrift format
	//return &getBulkInfo, err

	rpcHdl.dmnServer.ReqChan <- &server.ServerRequest{
		Op: server.GETBLK_DHCPV6RELAY_INTFSERVER_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: int(fromIdx),
			Count:   int(count),
		}),
	}

	result := <-rpcHdl.dmnServer.ReplyChan
	return result.Obj.(*server.GetBulkDHCPv6RelayIntfServerStateOutArgs).Obj, result.Error
}
