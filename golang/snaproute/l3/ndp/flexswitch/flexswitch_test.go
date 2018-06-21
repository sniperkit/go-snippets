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
	"l3/ndp/api"
	"ndpd"
	"testing"
)

func baseNdpFSInit() {
	api.Init(nil)
}

func TestCreateNDPGlobal(t *testing.T) {
	baseNdpFSInit()
	h := NewConfigHandler()
	config := &ndpd.NDPGlobal{}
	h.CreateNDPGlobal(config)
}

func TestDeleteNDPGlobal(t *testing.T) {
	h := NewConfigHandler()
	config := &ndpd.NDPGlobal{}
	h.DeleteNDPGlobal(config)
}

func TestUpdateNDPGlobal(t *testing.T) {
	h := NewConfigHandler()
	config := &ndpd.NDPGlobal{}
	newConfig := &ndpd.NDPGlobal{}
	attrset := make([]bool, 0)
	op := make([]*ndpd.PatchOpInfo, 0)
	h.UpdateNDPGlobal(config, newConfig, attrset, op)
}

func TestGetBulkNDPEntry(t *testing.T) {
	/*
		h := NewConfigHandler()
		h.GetBulkNDPEntryState(0, 10)
	*/
}

func TestGetNDPEntry(t *testing.T) {
	/*
		h := NewConfigHandler()
		h.GetNDPEntryState("")
	*/
}
