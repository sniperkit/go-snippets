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

package dmnBase

import (
	"encoding/json"
	"flag"
	"fmt"
	//_nanomsg "github.com/op/go-nanomsg"
	"io/ioutil"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/dbutils"
	"utils/keepalive"
	"utils/logging"
)

const (
	CLIENTS_FILE_NAME = "clients.json"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type FSBaseDmn struct {
	DmnName     string
	ParamsDir   string
	LogPrefix   string
	Logger      *logging.Writer
	DbHdl       *dbutils.DBUtil
	ClientsList []commonDefs.ClientJson
}

func (dmn *FSBaseDmn) InitLogger() (err error) {
	fmt.Println(dmn.LogPrefix, " Starting ", dmn.DmnName, "logger")
	dmnLogger, err := logging.NewLogger(dmn.DmnName, dmn.LogPrefix, true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
		return err
	}
	dmn.Logger = dmnLogger
	return err
}

func (dmn *FSBaseDmn) InitDBHdl() (err error) {
	dbHdl := dbutils.NewDBUtil(dmn.Logger)
	err = dbHdl.Connect()
	if err != nil {
		dmn.Logger.Err("Failed to dial out to Redis server")
		return err
	}
	dmn.DbHdl = dbHdl
	return err
}

func (dmn *FSBaseDmn) Init() bool {
	err := dmn.InitLogger()
	if err != nil {
		return false
	}
	err = dmn.InitDBHdl()
	if err != nil {
		return false
	}
	configFile := dmn.ParamsDir + "clients.json"
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		dmn.Logger.Info("Error in reading configuration file ", configFile)
		return false
	}
	err = json.Unmarshal(bytes, &dmn.ClientsList)
	if err != nil {
		dmn.Logger.Info("Error in Unmarshalling Json")
		return false
	}
	dmn.Logger.Info("Base daemon init completed")
	return true
}

func (dmn *FSBaseDmn) GetParams() string {
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	dirName := *paramsDir
	if dirName[len(dirName)-1] != '/' {
		dirName = dirName + "/"
	}
	return dirName
}

func (dmn *FSBaseDmn) StartKeepAlive() {
	go keepalive.InitKeepAlive(dmn.DmnName, dmn.ParamsDir)
}

func NewBaseDmn(dmnName, logPrefix string) *FSBaseDmn {
	var dmn = new(FSBaseDmn)
	dmn.DmnName = dmnName
	dmn.LogPrefix = logPrefix
	dmn.ParamsDir = dmn.GetParams()
	return dmn
}

func (dmn *FSBaseDmn) InitSwitch(plugin, dmnName, logPrefix string, switchHdl commonDefs.AsicdClientStruct) asicdClient.AsicdClientIntf {
	// @TODO: need to change second argument
	return asicdClient.NewAsicdClientInit(plugin, dmn.ParamsDir+"clients.json", switchHdl)

}
