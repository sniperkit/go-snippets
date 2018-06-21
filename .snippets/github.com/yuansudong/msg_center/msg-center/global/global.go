package global

import (
	"fmt"
	"jkt/jktgo/log"
	"jkt/jktgo/redis"
)

type (
	// Global 是用于存储全局变量的一些东西,并且初始化全局
	Global struct {
		storeServer string
		logServer   string
		logName     string
		runMode     string
	}
)

// createGlobal 用于实例化全局变量
func createGlobal() *Global {
	pGlobal := &Global{}
	pGlobal.initStoreServer()
	pGlobal.initLog()
	return nil
}

func (g *Global) initStoreServer() {
	redis.InitRedis(g.storeServer)
}
func (g *Global) initLog() {
	var err error
	err = log.InitLog(g.logServer, g.logName)
	if err != nil {
		panic(fmt.Sprintf("init log failed, because of %s", err.Error()))
	}
	if g.runMode == "dev" {
		log.SetLogLevel(log.LEVEL_DEBUG)
		log.SetLogModel(log.MODEL_DEV)
	} else if g.runMode == "info" {
		log.SetLogLevel(log.LEVEL_INFO)
		log.SetLogModel(log.MODEL_INFO)
	} else {
		log.SetLogLevel(log.LEVEL_WARNING)
		log.SetLogModel(log.MODEL_PRO)
	}
}


