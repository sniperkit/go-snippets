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

package ovsMgr

import (
	ovsdb "github.com/socketplane/libovsdb"

	_ "fmt"
	"l3/bgp/rpc"
	"reflect"
	"utils/logging"
)

const (
	// OVSDB Server Location
	OVSDB_HANDLER_HOST_IP   = "10.1.10.229"
	OVSDB_HANDLER_HOST_PORT = 6640

	// OVSDB Table
	OVSDB_HANDLER_DB_TABLE = "OpenSwitch"

	// OVSDB macro defines
	OVSDB_HANDLER_OPERATIONS_SIZE = 1024
	OVSDB_FS_INITIAL_SIZE         = 100
)

type BGPOvsdbNotifier struct {
	updateCh chan *ovsdb.TableUpdates
}

type BGPOvsOperations struct {
	operations []ovsdb.Operation
}

type BGPOvsRouterInfo struct {
	asn           uint32 // this is bgp asn number
	uuid          UUID   // this is key
	routerId      string // This is ip address
	HoldTime      uint32 // default is 180 seconds
	KeepaliveTime uint32 // default is 60 seconds
}

type BGPOvsdbHandler struct {
	logger         *logging.Writer
	bgpovs         *ovsdb.OvsdbClient
	ovsUpdateCh    chan *ovsdb.TableUpdates
	cache          map[string]map[string]ovsdb.Row
	operCh         chan *BGPOvsOperations
	bgpCachedOvsdb map[UUID]BGPFlexSwitch
	routerInfo     *BGPOvsRouterInfo
	rpcHdl         *rpc.BGPHandler
}

func NewBGPOvsdbNotifier(ch chan *ovsdb.TableUpdates) *BGPOvsdbNotifier {
	return &BGPOvsdbNotifier{
		updateCh: ch,
	}
}

func NewBGPOvsdbHandler(logger *logging.Writer, handler *rpc.BGPHandler) (*BGPOvsdbHandler, error) {
	ovs, err := ovsdb.Connect(OVSDB_HANDLER_HOST_IP, OVSDB_HANDLER_HOST_PORT)
	if err != nil {
		return nil, err
	}
	ovsUpdateCh := make(chan *ovsdb.TableUpdates)
	n := NewBGPOvsdbNotifier(ovsUpdateCh)
	ovs.Register(n)

	return &BGPOvsdbHandler{
		logger:         logger,
		bgpovs:         ovs,
		ovsUpdateCh:    ovsUpdateCh,
		operCh:         make(chan *BGPOvsOperations, OVSDB_HANDLER_OPERATIONS_SIZE),
		cache:          make(map[string]map[string]ovsdb.Row),
		bgpCachedOvsdb: make(map[UUID]BGPFlexSwitch, OVSDB_FS_INITIAL_SIZE),
		rpcHdl:         handler,
	}, nil
}

/*  BGP OVS DB populate cache with the latest update information from the
 *  notification channel
 */
func (ovsHdl *BGPOvsdbHandler) PopulateOvsdbCache(updates ovsdb.TableUpdates) {
	for table, tableUpdate := range updates.Updates {
		if _, ok := ovsHdl.cache[table]; !ok {
			ovsHdl.cache[table] = make(map[string]ovsdb.Row)
		}

		for uuid, row := range tableUpdate.Rows {
			empty := ovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				ovsHdl.cache[table][uuid] = row.New
			} else {
				delete(ovsHdl.cache[table], uuid)
			}
		}
	}
}

/* Stub interfaces for ovsdb library notifier
 */
func (ovsHdl BGPOvsdbNotifier) Update(context interface{}, tableUpdates ovsdb.TableUpdates) {
	ovsHdl.updateCh <- &tableUpdates
}

/* Stub interfaces for ovsdb library notifier
 */
func (ovsHdl BGPOvsdbNotifier) Locked([]interface{}) {
}

/* Stub interfaces for ovsdb library notifier
 */
func (ovsHdl BGPOvsdbNotifier) Stolen([]interface{}) {
}

/* Stub interfaces for ovsdb library notifier
 */
func (ovsHdl BGPOvsdbNotifier) Echo([]interface{}) {
}

/* Stub interfaces for ovsdb library notifier
 */
func (ovsHdl BGPOvsdbNotifier) Disconnected(client *ovsdb.OvsdbClient) {
}

/*  BGP OVS DB transaction api handler
 */
func (ovsHdl *BGPOvsdbHandler) Transact(operations []ovsdb.Operation) error {
	return nil
}

/*  BGP OVS DB handle update information
 */
func (ovsHdl *BGPOvsdbHandler) UpdateInfo(updates ovsdb.TableUpdates) {
	table, ok := updates.Updates[OVSDB_BGP_ROUTER_TABLE]
	if ok {
		ovsHdl.logger.Info("BGP_Router table Update")
		err := ovsHdl.HandleBGPRouteUpd(table)
		if err != nil {
			ovsHdl.logger.Err(err)
			return
		}
	}
	table, ok = updates.Updates[OVSDB_BGP_NEIGHBOR_TABLE]
	if ok {
		err := ovsHdl.HandleBGPNeighborUpd(table)
		if err != nil {
			ovsHdl.logger.Err(err)
			return
		}
		ovsHdl.logger.Info(ovsHdl.routerInfo)
	}
}

/*
 *  BGP OVS DB server.
 *	This API will handle reading operations from table... It can also do
 *	transactions.... In short its read/write bgp ovsdb handler
 */
func (ovsHdl *BGPOvsdbHandler) StartMonitoring() error {
	initial, err := ovsHdl.bgpovs.MonitorAll(OVSDB_HANDLER_DB_TABLE, "")
	if err != nil {
		return err
	}

	go func() {
		ovsHdl.ovsUpdateCh <- initial
	}()

	for {
		select {
		case updates := <-ovsHdl.ovsUpdateCh:
			ovsHdl.PopulateOvsdbCache(*updates)
			ovsHdl.UpdateInfo(*updates)
		case oper := <-ovsHdl.operCh:
			if err := ovsHdl.Transact(oper.operations); err != nil {
				//@FIXME: add some error message if needed
			}
		}
	}
	return nil
}
