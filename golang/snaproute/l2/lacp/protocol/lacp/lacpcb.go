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
// callback signature for drcp in order to get event updates
//
// lacpcb.go
package lacp

// port events
type LacpPortEvtCb func(ifindex int32)
type LacpAggEvtCb func(ifindex int32)

type LacpCbDbEntry struct {
	PortCreateDbList  map[string]LacpPortEvtCb
	PortDeleteDbList  map[string]LacpPortEvtCb
	PortUpDbList      map[string]LacpPortEvtCb
	PortDownDbList    map[string]LacpPortEvtCb
	AggCreateDbList   map[string]LacpAggEvtCb
	AggDeleteDbList   map[string]LacpAggEvtCb
	AggOperUpDbList   map[string]LacpAggEvtCb
	AggOperDownDbList map[string]LacpAggEvtCb
}

var LacpCbDb LacpCbDbEntry

func RegisterLaPortCreateCb(owner string, cb LacpPortEvtCb) {
	LacpCbDb.PortCreateDbList[owner] = cb
}

func RegisterLaPortDeleteCb(owner string, cb LacpPortEvtCb) {
	LacpCbDb.PortDeleteDbList[owner] = cb
}

func RegisterLaPortUpCb(owner string, cb LacpPortEvtCb) {
	LacpCbDb.PortUpDbList[owner] = cb
}

func RegisterLaPortDownCb(owner string, cb LacpPortEvtCb) {
	LacpCbDb.PortDownDbList[owner] = cb
}

func RegisterLaAggCreateCb(owner string, cb LacpAggEvtCb) {
	LacpCbDb.AggCreateDbList[owner] = cb
}

func RegisterLaAggDeleteCb(owner string, cb LacpAggEvtCb) {
	LacpCbDb.AggDeleteDbList[owner] = cb
}

func RegisterLaAggOperStateUpCb(owner string, cb LacpAggEvtCb) {
	LacpCbDb.AggOperUpDbList[owner] = cb
}

func RegisterLaAggOperStateDownCb(owner string, cb LacpAggEvtCb) {
	LacpCbDb.AggOperDownDbList[owner] = cb
}

func DeRegisterLaAggCbAll(owner string) {
	delete(LacpCbDb.PortCreateDbList, owner)
	delete(LacpCbDb.PortDeleteDbList, owner)
	delete(LacpCbDb.PortUpDbList, owner)
	delete(LacpCbDb.PortDownDbList, owner)
	delete(LacpCbDb.AggCreateDbList, owner)
	delete(LacpCbDb.AggDeleteDbList, owner)
}
