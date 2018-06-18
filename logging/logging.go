package logging

import "log"

func LogError(v ...interface{}) {
	log.Println("[ERROR] ", v)
}

func LogWarn(v ...interface{}) {
	log.Println("[WARN] ", v)
}

func LogInfo(v ...interface{}) {
	log.Println("[INFO] ", v)
}
