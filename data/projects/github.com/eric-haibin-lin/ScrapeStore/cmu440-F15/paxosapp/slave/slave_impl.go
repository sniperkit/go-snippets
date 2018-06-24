package slave

import (
	"fmt"
	"github.com/cmu440-F15/paxosapp/rpc/slaverpc"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

/* slaveNode is a storage node, which acts like the dumb chunkservers in GFS.
   They are basically a storage medium for storing key-value maps */
type slaveNode struct {
	/* This map holds the scraped data, where the key is a URL, and the value is a slice of slices of strings */
	valuesMap map[string][][]string

	/* this map holds the page ranks of URLs, where the key is a URL, and the value is a float representing
	   the pagerank of that URL */
	ranksMap map[string]float64

	/* locks for the maps */
	valuesMapLock *sync.Mutex
	ranksMapLock  *sync.Mutex

	/* ID of this slave */
	srvId int

	/* hostport of this slave */
	myHostPort string
}

/* this method instantiates a slave node, and exports several RPCs that the ScrapeStore masters can call for interacting
   for their storage needs */
func NewSlaveNode(myHostPort string, srvId int) (SlaveNode, error) {
	fmt.Println("NewSlaveNode invoked on", srvId)
	var a slaverpc.RemoteSlaveNode

	/* instantiate and init a new slave */
	node := slaveNode{}

	/* set all values */
	node.valuesMap = make(map[string][][]string)
	node.ranksMap = make(map[string]float64)
	node.valuesMapLock = &sync.Mutex{}
	node.ranksMapLock = &sync.Mutex{}
	node.srvId = srvId
	node.myHostPort = myHostPort

	/* setup RPCs */
	listener, err := net.Listen("tcp", myHostPort)
	if err != nil {
		return nil, err
	}

	a = &node

	/* wrap and register the set of RPCs that this slave will support */
	err = rpc.RegisterName("SlaveNode", slaverpc.Wrap(a))
	if err != nil {
		return nil, err
	}

	/* create a goroutine to accept and listen for RPCs */
	rpc.HandleHTTP()
	go http.Serve(listener, nil)

	return a, nil
}

/* This method accepts a slice of strings, and appends it to the slice of slice of strings for that key */
func (sn *slaveNode) Append(args *slaverpc.AppendArgs, reply *slaverpc.AppendReply) error {
	fmt.Println("Append invoked on ", sn.srvId, " for Key: ", args.Key)
	defer fmt.Println("Leaving Append on ", sn.srvId)
	sn.valuesMapLock.Lock()
	defer sn.valuesMapLock.Unlock()
	key := args.Key
	value := args.Value

	/* see if the key is present */
	_, ok := sn.valuesMap[key]

	if !ok { /* not present currently, so first create a slice of slice of strings */
		fmt.Println("Didn't find anything on slave ", sn.srvId, ", so creating a new slice")
		sn.valuesMap[key] = make([][]string, 0)
	}

	/* Append this slice to the slice of slices */
	sn.valuesMap[key] = append(sn.valuesMap[key], value)
	return nil
}

/* get returns the last appended slice of strings for the given key */
func (sn *slaveNode) Get(args *slaverpc.GetArgs, reply *slaverpc.GetReply) error {
	fmt.Println("Get invoked on ", sn.srvId, " for key : ", args.Key)
	defer fmt.Println("Leaving Get on ", sn.srvId)
	sn.valuesMapLock.Lock()
	defer sn.valuesMapLock.Unlock()
	key := args.Key
	/* first see if this node has the value for this key */
	value, ok := sn.valuesMap[key]
	if ok { /* If the key was present, return the last appended slice */
		reply.Status = 1
		reply.Value = value[len(value)-1]
	} else { /* If not present, return Status as 0 */
		reply.Status = 0
	}
	return nil
}

/* GetRank returns the rank of the given key */
func (sn *slaveNode) GetRank(args *slaverpc.GetRankArgs, reply *slaverpc.GetRankReply) error {
	fmt.Println("GetRank invoked on ", sn.srvId)
	defer fmt.Println("Leaving GetRank on ", sn.srvId)
	sn.ranksMapLock.Lock()
	defer sn.ranksMapLock.Unlock()
	key := args.Key
	/* lookup the rank for the given key. This has to be present here, because if it weren't, the
	   master wouldn't make this RPC on this node in the first place */
	value := sn.ranksMap[key]
	fmt.Println("Rank found for key", key, ": ", value)
	reply.Value = value
	return nil
}

/* PutRank sets the page rank for a given key */
func (sn *slaveNode) PutRank(args *slaverpc.PutRankArgs, reply *slaverpc.PutRankReply) error {
	fmt.Println("PutRank invoked on ", sn.srvId)
	defer fmt.Println("Leaving PutRank on ", sn.srvId)
	sn.ranksMapLock.Lock()
	defer sn.ranksMapLock.Unlock()
	key := args.Key
	/* Just set the rank for this key */
	sn.ranksMap[key] = args.Value
	fmt.Println("Rank Put for key", key, ": ", sn.ranksMap[key])
	return nil
}
