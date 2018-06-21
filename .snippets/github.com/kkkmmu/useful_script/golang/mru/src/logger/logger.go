package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"util"
)

var loggerLock = sync.Mutex{}
var PrintToTerminal bool = false

func Push(ctx context.Context, message string) {
	sid := ctx.Value("SESSIONID")
	sessionid, _ := sid.(string)

	loggerLock.Lock()
	defer loggerLock.Unlock()

	logfile := sessionid
	if sessionid == "" {
		logfile = "default"
	}

	file, err := os.OpenFile("asset/log/"+logfile+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Println("cannot Open file: ", logfile+".log", " ", err.Error())
		return
	}
	file.WriteString(message)
	file.Close()

	//log to detail log file
	full, err := os.OpenFile("asset/log/"+logfile+"_full.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Println("cannot Open file: ", logfile+"_full.log", " ", err.Error())
		return
	}
	full.WriteString(message)
	full.Close()

	if PrintToTerminal {
		fmt.Printf("%s", message)
	}
}

func init() {
	exist, err := util.DirExists("asset/log")
	if !exist {
		err = os.MkdirAll("asset/log", 0777)
		if err != nil {
			fmt.Println("Log file will not created due to: ", err)
		}
	} else if err != nil {
		fmt.Println("Log file will not created due to: ", err)
	}
}
