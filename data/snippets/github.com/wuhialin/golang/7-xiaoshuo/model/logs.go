package model

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"path/filepath"
)

var myLog *logs.BeeLogger

func init() {
	myLog = logs.GetBeeLogger()
	myLog.SetLogger(logs.AdapterConsole)

	path, err := filepath.Abs("./logs")
	if err != nil {
		panic(err)
	}
	fileConfig := map[string]interface{}{}
	fileConfig["filename"] = path + "/app.log"
	fileConfig["daily"] = true
	jsonFileConfig, err := json.Marshal(fileConfig)
	if err != nil {
		panic(err)
	}
	myLog.SetLogger(logs.AdapterMultiFile, string(jsonFileConfig))

	myLog.EnableFuncCallDepth(true)
	myLog.SetLogFuncCallDepth(3)
	myLog.Async(1e3)
	myLog.SetLevel(logs.LevelTrace)
}

func Log() *logs.BeeLogger {
	return myLog
}
