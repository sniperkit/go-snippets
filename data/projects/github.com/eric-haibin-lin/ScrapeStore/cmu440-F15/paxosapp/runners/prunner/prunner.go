// DO NOT MODIFY!

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/cmu440-F15/paxosapp/paxos"
)

var (
	ports       = flag.String("ports", "", "ports for all paxos nodes")
	numNodes    = flag.Int("N", 1, "the number of nodes in the ring")
	nodeID      = flag.Int("id", 0, "node ID must match index of this node's port in the ports list")
	numRetries  = flag.Int("retries", 5, "number of times a node should retry dialing another node")
	slaveports  = flag.String("slaveports", "", "ports for all slave nodes")
	numslaves   = flag.Int("numslaves", 0, "the number of slaves in the ring")
	replace     = flag.Bool("replace", false, "whether it's a replacement of existing Paxos Nodes")
	monitorPort = flag.String("monitorport", "", "the port of the cluster monitor node")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()

	fmt.Println("Got slaveports as ", *slaveports)
	fmt.Println("Got numslaves as ", *numslaves)

	portStrings := strings.Split(*ports, ",")
	slaveStrings := strings.Split(*slaveports, ",")
	monitorString := "localhost:" + *monitorPort

	hostMap := make(map[int]string)
	slaveMap := make(map[int]string)

	for i, port := range portStrings {
		hostMap[i] = "localhost:" + port
	}

	for i, port := range slaveStrings {
		slaveMap[i] = "localhost:" + port
	}
	// Create and start the Paxos Node.

	var err error
	if *replace {
		_, err = paxos.NewPaxosNode(hostMap[*nodeID], monitorString, hostMap, *numNodes, *nodeID, *numRetries, true, slaveMap, *numslaves)
	} else {
		_, err = paxos.NewPaxosNode(hostMap[*nodeID], monitorString, hostMap, *numNodes, *nodeID, *numRetries, false, slaveMap, *numslaves)
	}

	if err != nil {
		log.Fatalln("Failed to create paxos node:", err)
	}

	// Run the paxos node forever.
	select {}
}
