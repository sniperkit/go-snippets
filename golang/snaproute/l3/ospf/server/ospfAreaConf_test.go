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

/*
ospfAreaConf_test.go
Area conf testing routines.
*/

package server

import (
	"fmt"
	// "l3/ospf/config"
	"testing"
)

func initHelloTestParams() {
	ospf = getServerObject()
	initAttr()
	go startDummyChannels(ospf)
}

func initAreaTestParams() {
	ospf = getServerObject()
	initAttr()
	go startDummyChannels(ospf)
}

func TestOspfArea(t *testing.T) {
	fmt.Println("\n**************** AREACONF ************\n")
	initHelloTestParams()
	for index := 1; index < 21; index++ {
		err := areaTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test case for Hello protocol ")
		}
	}

}

func areaTestLogic(tNum int) int {

	switch tNum {
	case 1:
		fmt.Println(tNum, ": Running initAreaConfDefault")
		ospf.initAreaConfDefault()

	case 2:
		fmt.Println(tNum, ": Running processAreaConfig")
		err := ospf.processAreaConfig(areaConf)
		if err == nil {
			fmt.Println("Successful to add areaconf to the map ")
		}

	case 3:
		fmt.Println(tNum, ": Running processAreaConfig")
		err := ospf.processAreaConfig(areaConf)
		if err == nil {
			ospf.initAreaStateSlice(areaConfKey)
		}

	case 4:
		fmt.Println(tNum, ": Running areaStateRefresh ")
		ospf.areaStateRefresh()

	case 5:
		fmt.Println(tNum, ": Running updateIntfToAreaMap ")
		oldAreaId := "10.0.0.0"
		newAreaId := "10.10.0.0"
		err := ospf.processAreaConfig(areaConf)
		if err == nil {
			ospf.updateIntfToAreaMap(key, oldAreaId, newAreaId, false)
		}

	case 6:
		fmt.Println(tNum, ": Running updateIfABR ")
		ospf.processAreaConfig(areaConf)
		ospf.updateIfABR()
		fmt.Println(" New status for ABR : ", ospf.ospfGlobalConf.AreaBdrRtrStatus)

	case 7:
		fmt.Println(tNum, ": Running isStubArea ")
		err := ospf.processAreaConfig(areaConf)
		if err == nil {
			stub := ospf.isStubArea(areaConfKey.AreaId)
			fmt.Println(" Is stub area : ", stub)
		}

	case 8:
		fmt.Println(tNum, ": Running initOspfGlobalConfDefault ")
		ospf.initOspfGlobalConfDefault()

	case 9:
		fmt.Println(tNum, ": Running processASBdrRtrStatus")
		ospf.processASBdrRtrStatus(true)

	case 10:
		fmt.Println(tNum, ": Running processGlobalConfig ")
		ospf.processGlobalConfig(gConf)
		conf := ospf.GetOspfGlobalState()
		fmt.Println("Global conf ", conf)
		checkAsicdAPIs()
	}

	return SUCCESS
}

func checkAsicdAPIs() {
	//ospf.listenForASICdUpdates("ribd")
	ospf.processAsicdNotification(hello)
	err := ospf.initAsicdForRxMulticastPkt()
	fmt.Println("Asicd initialised with err ", err)
}
