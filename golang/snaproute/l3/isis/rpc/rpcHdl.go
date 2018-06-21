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
        "isisd"
)

func (rpcHdl *rpcServiceHandler) CreateIsisGlobal(cfg *isisd.IsisGlobal) (bool, error) {
        rpcHdl.logger.Info("Calling CreateIsisGlobal", cfg)
        return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateIsisGlobal(oldCfg, newCfg *isisd.IsisGlobal, attrset []bool, op []*isisd.PatchOpInfo) (bool, error) {
        rpcHdl.logger.Info("Calling UpdateIsisGlobal", oldCfg, newCfg)
        return true, nil
}

func (rpcHdl rpcServiceHandler) DeleteIsisGlobal(cfg *isisd.IsisGlobal) (bool, error) {
        rpcHdl.logger.Info("Calling DeleteIsisGlobal", cfg)
        return true, nil
}

func (rpcHdl *rpcServiceHandler) GetIsisGlobalState(key string) (obj *isisd.IsisGlobalState, err error) {
        rpcHdl.logger.Info("Calling GetIsisGlobalState", key)
        return obj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkIsisGlobalState(fromIdx, count isisd.Int) (*isisd.IsisGlobalStateGetInfo, error) {
        var getBulkInfo isisd.IsisGlobalStateGetInfo
        var err error
        //info, err := api.GetBulkIsisGlobalState(int(fromIdx), int(count))
        //getBulkInfo.StartIdx = fromIdx
        //getBulkInfo.EndIdx = isisd.Int(info.EndIdx)
        //getBulkInfo.More = info.More
        //getBulkInfo.Count = isisd.Int(len(info.List))
        // Fill in data, remember to convert back to thrift format
        return &getBulkInfo, err
}

