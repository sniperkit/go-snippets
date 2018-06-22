package time

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	aa := time.Now()
	time.Sleep(1 * time.Second)
	log.Println(aa)
	bb := time.Since(aa)
	bb = bb - (bb % time.Second)
	log.Println(bb)
}

func TestUnixTimestamp(t *testing.T) {
	log.Println(time.Now().Unix())
	log.Println(time.Unix(0, 0))
}

func TestTimeString(t *testing.T) {
	log.Println(time.Now().String())
}

// Duration must come with unit

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	// loc, _ := time.LoadLocation("America/Denver")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return fmt.Print(time.Now().In(loc).Format("2006-01-02T15:04:05.999Z0700") + " [DEBUG] " + string(bytes))
}
func TestTimeZoneTransform(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("hello new time format")
}
