package main
import (
	"sync"
	"fmt"
)

func main() {
	var data int
	var memAccess sync.Mutex
	go func() {
		memAccess.Lock()
		data++
		memAccess.Unlock()
	}()

	memAccess.Lock()
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	}
	memAccess.Unlock()
}