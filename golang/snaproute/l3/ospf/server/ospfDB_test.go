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

package server

import (
	"fmt"
	//   "l3/ospf/config"
	"testing"
)

func initDBTestParams() {
	fmt.Println("\n Get Server object")
	ospf = getServerObject()
	initAttr()
	go startDummyChannels(ospf)
	//fmt.Println(" Initialize db ")
	//ospf.InitializeDB()
	fmt.Println(" Init DB channels. ")
	ospf.InitDBChannels()
	fmt.Println("Start DB listener ")
	go ospf.StartDBListener()
}

func TestOspfDB(t *testing.T) {
	fmt.Println("\n**************** STATE DB ************\n")
	initDBTestParams()
	for index := 1; index < 21; index++ {
		err := dbTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test case for state db")
		}
	}
}

func dbTestLogic(tNum int) int {

	switch tNum {
	case 1:
		fmt.Println(tNum, ": Running DbReadConfig ")
		ospf.DbReadConfig <- true

	case 2:
		fmt.Println(tNum, ": Running DbRouteOp ")
		msg := DbRouteMsg{
			entry: rKey,
			op:    true,
		}

		ospf.DbRouteOp <- msg

	case 3:
		ospf.DbLsdbOp <- lsdbMsg

	case 4:
		ospf.DbEventOp <- eventMsg

	case 5:
		fmt.Println(tNum, ": Running ReadOspfCfgFromDB ")
		ospf.ReadOspfCfgFromDB()

	case 6:
		fmt.Println(tNum, "applyOspfGlobalConf ")
		//ospf.applyOspfGlobalConf(gConf)
	}

	return SUCCESS
}
