package basic

import (
	"log"
	"testing"
	"time"
)

func init() {
	log.Println("init again")
}

var startTime = startup()

func init() {
	log.Println("init.....")
}

func startup() time.Time {
	log.Println("startup.....")

	return time.Now()
}

func TestPackageInit(t *testing.T) {
	log.Println("test....")
}
