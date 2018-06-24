package monitor

import (
	"fmt"
	"github.com/cmu440-F15/paxosapp/rpc/monitorrpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

type monitorNode struct {
	//the host port information of all master nodes
	masterHostPortMap map[int]string
	//the heart beat information of all master nodes
	masterHeartBeatMap map[int]int
	//the host port of the monitor itself
	myHostPort string
}

// NewMonitorNode creates a new monitor node. The monitor node has the host port
// information of all master nodes. It also allows heartbeat rpc calls from master node
// to report that the master node is not down. It will not return until the monitor
// node can successfully handle rpc calls from masters
func NewMonitorNode(myHostPort string, masterHostPort []string) (MonitorNode, error) {
	log.SetOutput(ioutil.Discard)
	log.Println("myhostport is ", myHostPort, "masterHostPort is ", masterHostPort)
	defer log.Println("Leaving NewMonitorNode")
	var a monitorrpc.RemoteMonitorNode
	node := monitorNode{}
	node.masterHostPortMap = make(map[int]string)
	node.masterHeartBeatMap = make(map[int]int)
	node.myHostPort = myHostPort
	for id := 0; id < len(masterHostPort); id++ {
		node.masterHostPortMap[id] = masterHostPort[id]
	}
	//listen to monitor hostport and start to serve RPC's
	listener, err := net.Listen("tcp", myHostPort)
	if err != nil {
		return nil, err
	}
	a = &node
	err = rpc.RegisterName("MonitorNode", monitorrpc.Wrap(a))
	if err != nil {
		return nil, err
	}
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	go (&node).CheckHealth()
	return a, nil
}

// HeartBeat RPC is for master nodes to contact monitor server. The monitor server
// mark down the master's id for liveness check
func (sn *monitorNode) HeartBeat(args *monitorrpc.HeartBeatArgs, reply *monitorrpc.HeartBeatReply) error {
	if args.Type == monitorrpc.Master {
		sn.masterHeartBeatMap[args.Id] += 1
		fmt.Println("Received heartbeat from server", args.Id)
	}
	return nil
}

// CheckHealth count the number of heartbeats from every master node at some fixed interval.
// If some master node doesn't heartbeat monitor, their id's will be recorded and restarted by
// the monitor.
func (sn *monitorNode) CheckHealth() {
	time.Sleep(time.Second * 8)
	fmt.Println("master nodes health monitor is started")
	for {
		//Check all master nodes' heartbeat
		for index, _ := range sn.masterHostPortMap {
			_, ok := sn.masterHeartBeatMap[index]
			if !ok {
				//Not receiving heartbeat indicates that this master is down,
				//replace with new master node on the same port
				fmt.Println("Master ", index, "is down. Replacing this node..")
				filePath := os.Getenv("GOPATH") + "/src/github.com/cmu440-F15/scripts/server_id.txt"
				idToReplace := []byte(strconv.Itoa(index) + "\n")
				ioutil.WriteFile(filePath, idToReplace, 0666)
				f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660)
				if err != nil {
					log.Println(err)
				}
				f.Write(idToReplace)
				f.Sync()
				f.Close()
			} else {
				//Received heartbeat from this master before. Clear the entry for next round
				delete(sn.masterHeartBeatMap, index)
			}
		}
		time.Sleep(time.Second * 6)
	}
}
