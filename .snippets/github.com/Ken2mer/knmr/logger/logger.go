package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

// Verbose sets log level. If false: INFO(default), if true: DEBUG.
var Verbose = false

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

// DumpJSON prints with double space
func DumpJSON(v interface{}) (n int, err error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return 0, err
	}
	return fmt.Printf("%s\n\n", data)
}

// Debug println
func Debug(v ...interface{}) {
	log.SetPrefix("[DEBUG] ")
	if Verbose == true {
		log.Println(v...)
	}
}

// Debugf printf
func Debugf(format string, v ...interface{}) {
	log.SetPrefix("[DEBUG] ")
	if Verbose == true {
		log.Printf(format, v...)
	}
}

// Info println
func Info(v ...interface{}) {
	log.SetPrefix("[INFO] ")
	log.Println(v...)
}

// Infof printf
func Infof(format string, v ...interface{}) {
	log.SetPrefix("[INFO] ")
	log.Printf(format, v...)
}

// Error fatal
func Error(v ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Fatal(v...)
}

// Errorf fatalf
func Errorf(format string, v ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Fatalf(format, v...)
}
