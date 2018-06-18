package main

import (
	"fmt"
	"time"
)

func main() {
	// Use time.Now with Unix or UnixNano to get elapsed time since the Unix epoch in seconds or nanoseconds, respectively.
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println(now)
	// Note that there is no UnixMillis, so to get the milliseconds since epoch you’ll need to manually divide from nanoseconds.
	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)
	// You can also convert integer seconds or nanoseconds since the epoch into the corresponding time.
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}

// 2012-10-31 16:13:58.292387 +0000 UTC
// 1351700038
// 1351700038292
// 1351700038292387000
// 2012-10-31 16:13:58 +0000 UTC
// 2012-10-31 16:13:58.292387 +0000 UTC
