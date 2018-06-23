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
	"ndpd"
	"ribd"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/logging"

	nanomsg "github.com/op/go-nanomsg"
)

/*  Router manager will handle all the communication with ribd
 */
type FSRouteMgr struct {
	plugin          string
	logger          *logging.Writer
	ribdClient      *ribd.RIBDServicesClient
	ribSubSocket    *nanomsg.SubSocket
	ribSubBGPSocket *nanomsg.SubSocket
}

/*  Interface manager will handle all the communication with asicd
 */
type FSIntfMgr struct {
	plugin                   string
	logger                   *logging.Writer
	AsicdClient              asicdClntIntfs.AsicdClntIntf
	NdpdClient               *ndpd.NDPDServicesClient
	asicdL3IntfSubSocket     *nanomsg.SubSocket
	ndpIntfSubSocket         *nanomsg.SubSocket
	enableAsicdNotifications bool
}

/*  @FUTURE: this will be using in future if FlexSwitch is planning to support
 *	     daemon which is handling policy statments
 */
type FSPolicyMgr struct {
	plugin          string
	logger          *logging.Writer
	policySubSocket *nanomsg.SubSocket
}

/*  BFD manager will handle all the communication with bfd daemon
 */
type FSBfdMgr struct {
	plugin       string
	logger       *logging.Writer
	bfddClient   *bfdd.BFDDServicesClient
	bfdSubSocket *nanomsg.SubSocket
}

func (mgr *FSIntfMgr) PortStateChange() {

}
