// DO NOT MODIFY!

package main

import (
	"flag"
	"fmt"
	"github.com/cmu440-F15/paxosapp/slave"
	"log"
)

var (
	hostport = flag.String("hostport", "", "hostport of this slave")
	id       = flag.Int("id", -1, "id of this slave")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()

	fmt.Println("Slave got hostport as ", *hostport, " and id as ", *id)

	// Create and start the Paxos Node.
	_, err := slave.NewSlaveNode(*hostport, *id)
	if err != nil {
		log.Fatalln("Failed to create slave node:", err)
	}

	// Run the paxos node forever.
	select {}
}
