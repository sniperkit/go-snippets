package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	logger := log.New(file, "PREFIX: ", log.Ldate|log.Ltime|log.Lshortfile)
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//This writes 'INFO: 2017/02/15 19:01:16 exampleprogram.go:exampleline: example error' to log
	Info.Printf("example error")
	//This writes 'WARNING: 2017/02/15 19:01:16 exampleprogram.go:exampleline: example error' to log
	Warning.Printf("example error")
	//This writes 'ERROR: 2017/02/15 19:01:16 exampleprogram.go:exampleline: example error' to log
	Error.Printf("example error")

	fmt.Println("Hello, World!")

}
