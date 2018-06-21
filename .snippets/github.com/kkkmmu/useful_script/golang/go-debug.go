package main

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"time"
)

func debugPoint() string {
	pc, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("\033[31m%v %s %s %d\x1b[0m", time.Now(),
		runtime.FuncForPC(pc).Name(), path.Base(file), line)
}

func f() {
	fmt.Println(debugPoint())
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("log with flags")

	f()

}
