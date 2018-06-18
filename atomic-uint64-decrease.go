package main

import (
	"log"
	"sync/atomic"
)

func main() {
	var count uint64

	// increase 2
	atomic.AddUint64(&count, 2)
	log.Printf("count increase one: %d", atomic.LoadUint64(&count))

	// Method A: (Nobody can't know this way)
	// decrease 1 no way, type of 'count' is unsigned
	//atomic.AddUint64(&count, -1)
	//log.Printf("count decrease one: %d", atomic.LoadUint64(&count))

	// Method B: (bit operation)
	atomic.AddUint64(&count, ^uint64(0))
	log.Printf("count decrease one: %d", atomic.LoadUint64(&count))
}
