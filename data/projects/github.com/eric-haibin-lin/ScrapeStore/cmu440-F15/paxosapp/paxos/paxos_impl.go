package paxos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cmu440-F15/paxosapp/common"
	"github.com/cmu440-F15/paxosapp/rpc/monitorrpc"
	"github.com/cmu440-F15/paxosapp/rpc/paxosrpc"
	"github.com/cmu440-F15/paxosapp/rpc/slaverpc"
	"net"
	"net/http"
	"net/rpc"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/* NumCopies holds the number of copies that have to be made of each data element */
const (
	NumCopies = 5
)

/* The paxosNode struct
   Each paxos node acts as a decentralized metadata master in ScrapeStore
*/
type paxosNode struct {
	/* Hostport of this paxos node */
	myHostPort string

	/* Monitor connection to send heartbeat */
	monitor common.Conn

	/* The number of paxos nodes in the ring */
	numNodes int

	/* The number of slave (storage) nodes in ScrapeStore */
	numSlaves int

	/* The ID of this paxos node */
	srvId int

	/* This map holds a mapping of ID->hostport of all paxos nodes in the paxos ring */
	hostMap map[int]string

	/* This map holds a mapping of ID->hostport of all slave nodes in ScrapeStore */
	slaveMap map[int]string

	/* These maps cache the connections to all nodes in ScrapeStore */
	paxosDialerMap map[int]*rpc.Client
	slaveDialerMap map[int]*rpc.Client

	/* The main key-value store.
	   The key is either a string hodling a URL, or is a key of the form SLAVEID:size.
	   The value is an array of NumCopies ints.
	   When the key is a URL, the value is a list of all slave node IDs which have the
	   data for that key.
	   When the key is of the form SLAVEID:size, the first element of the array holds
	   the size (in Bytes) of the amount of data held by that slave node.
	*/
	valuesMap map[string][NumCopies]int

	/* Disclaimer - acceptedValuesMap and acceptedSeqNumMap should be used hand-in-hand */

	/* Temporary map which should be populated only when accept is called and
	should be cleared when commit is called */
	acceptedValuesMap map[string][NumCopies]int

	/* Temporary map which should be populated only when accept is called and
	should be cleared when commit is called */
	acceptedSeqNumMap map[string]int

	/* Again, a temporary map which should be populated when prepare is called */
	maxSeqNumSoFar map[string]int

	/* Next seqNum for particular key. Strictly increasing per key per node */
	nextSeqNumMap map[string]int

	/* Locks to protect all crucial data structures */
	maxSeqNumSoFarLock    *sync.Mutex
	valuesMapLock         *sync.Mutex
	acceptedValuesMapLock *sync.Mutex
	acceptedSeqNumMapLock *sync.Mutex
	nextSeqNumMapLock     *sync.Mutex
}

/* NewPaxosNode creates a new PaxosNode. This function should return only when
   all nodes have joined the ring, and should return a non-nil error if the node
   could not be started in spite of dialing the other nodes numRetries times.

   hostMap is a map from node IDs to their hostports, numNodes is the number
   of nodes in the ring, replace is a flag which indicates whether this node
   is a replacement for a node which failed.
*/
func NewPaxosNode(myHostPort, monitorHostPort string, hostMap map[int]string, numNodes, srvId, numRetries int, replace bool, slaveMap map[int]string, numSlaves int) (PaxosNode, error) {
	fmt.Println("myhostport is ", myHostPort, "Numnodes is ", numNodes, "srvid is ", srvId)

	var a paxosrpc.RemotePaxosNode

	/* Init */
	node := paxosNode{}

	/* Create and populate the members of the paxosNode struct */
	node.srvId = srvId
	node.numNodes = numNodes
	node.myHostPort = myHostPort
	node.hostMap = make(map[int]string)
	node.slaveMap = make(map[int]string)

	node.valuesMap = make(map[string][NumCopies]int)

	node.acceptedValuesMap = make(map[string][NumCopies]int)
	node.acceptedSeqNumMap = make(map[string]int)

	node.maxSeqNumSoFar = make(map[string]int)
	node.nextSeqNumMap = make(map[string]int)

	node.paxosDialerMap = make(map[int]*rpc.Client)
	node.slaveDialerMap = make(map[int]*rpc.Client)

	node.numSlaves = numSlaves

	/* Create locks */
	node.valuesMapLock = &sync.Mutex{}
	node.acceptedValuesMapLock = &sync.Mutex{}
	node.acceptedSeqNumMapLock = &sync.Mutex{}
	node.maxSeqNumSoFarLock = &sync.Mutex{}
	node.nextSeqNumMapLock = &sync.Mutex{}

	/* Populate maps */
	for k, v := range hostMap {
		node.hostMap[k] = v
	}

	for k, v := range slaveMap {
		node.slaveMap[k] = v
	}

	/* Set up RPCs on this node */
	listener, err := net.Listen("tcp", myHostPort)
	if err != nil {
		return nil, err
	}

	a = &node

	/* Register this node to listen RPCs */
	err = rpc.RegisterName("PaxosNode", paxosrpc.Wrap(a))
	if err != nil {
		return nil, err
	}

	/* Create a Goroutine to accept RPCs */
	rpc.HandleHTTP()
	go http.Serve(listener, nil)

	/* Iterate over hostmap, and dial each paxos node. */
	for key, v := range hostMap {
		/* Dial */
		dialer, err := rpc.DialHTTP("tcp", v)

		/* number of times we've tried dialing this paxos node */
		cntr := 0

		if err != nil {
			/* Dial failed. Keep trying! */
			for {
				fmt.Println(myHostPort, " couldn't dial", v, ". Trying again.")

				cntr = cntr + 1

				/* If couldn't dial all paxos nodes, abort. */
				if cntr == numRetries {
					fmt.Println("Couldn't connect even after all retries.", myHostPort, " aborting.")
					return nil, errors.New("Couldn't dial a node")
				}

				/* Try after 1 second */
				time.Sleep(1 * time.Second)
				dialer, err := rpc.DialHTTP("tcp", v)
				if err == nil {
					/* Success! Cache this dialer and continue */
					node.paxosDialerMap[key] = dialer
					break
				}
			}
		} else {
			/* Success! Cache this dialer and continue */
			node.paxosDialerMap[key] = dialer
		}

		fmt.Println(myHostPort, " dialed fellow paxosnode", v, " successfully")
	}

	/* Get updated values map when it's a node for replacement */
	if replace {
		nextSrv := ""
		var k int
		var v string
		for k, v = range hostMap {
			/* Use the first entry in hostMap to retrieve the values map.
			   Do not send the request to the new node itself */
			if v != myHostPort {
				nextSrv = v
				break
			}
		}
		args := paxosrpc.ReplaceCatchupArgs{}
		reply := paxosrpc.ReplaceCatchupReply{}

		/* Call ReplaceCatchup on this paxos node */
		err = node.paxosDialerMap[k].Call("PaxosNode.RecvReplaceCatchup", &args, &reply)
		if err != nil {
			fmt.Println("ERROR: Couldn't Dial RecvReplaceCatchup on ", nextSrv)
			return nil, errors.New("ERROR: Couldn't Dial RecvReplaceCatchup")
		}

		var f interface{}
		json.Unmarshal(reply.Data, &f)
		node.valuesMapLock.Lock()

		/* Convert []interface{} to [numCopies]int by iterating the interface{} slice */
		for key, arr := range f.(map[string]interface{}) {
			slices := arr.([]interface{})
			var copies [NumCopies]int
			for index, value := range slices {
				fmt.Println(index, int(value.(float64)))
				copies[index] = int(value.(float64))
			}
			node.valuesMap[key] = copies
		}
		node.valuesMapLock.Unlock()

		fmt.Println("Received values from peers. The value map has the following entries:")
		for k, v := range node.valuesMap {
			fmt.Println(k, v)
		}

		/* Now call RecvReplaceServer on each of the other nodes to inform them that
		   I am now taking the place of the failed node
		*/
		for k, v := range hostMap {
			if v != myHostPort {
				args := paxosrpc.ReplaceServerArgs{}
				reply := paxosrpc.ReplaceServerReply{}

				args.SrvID = srvId
				args.Hostport = myHostPort
				/* Make the RPC call */
				err = node.paxosDialerMap[k].Call("PaxosNode.RecvReplaceServer", &args, &reply)
				if err != nil {
					fmt.Println("ERROR: Couldn't Dial RecvReplaceServer on ", nextSrv)
				}
			}
		}
	} else {
		/* This node has been generated for the first time
		   set the sizes of all nodes to 0
		*/
		for i := 0; i < numSlaves; i++ {
			sizeKey := strconv.Itoa(i)
			/* Key is like 0:size, 1:size and so on */
			sizeKey = sizeKey + ":size"
			var sizeSlice [NumCopies]int
			sizeSlice[0] = 0
			node.valuesMap[sizeKey] = sizeSlice
		}
	}

	/* Now try dialing all slave nodes. They should ideally already be up */
	for k, v := range node.slaveMap {
		dialer, err := rpc.DialHTTP("tcp", v)

		cntr := 0

		/* Didn't succeed. Try again! */
		if err != nil {
			for {
				fmt.Println(myHostPort, " couldn't dial slave", v, ". Trying again.")

				cntr = cntr + 1
				/* If couldn't dial slave even after all attempts, abort */
				if cntr == numRetries {
					fmt.Println("Couldn't connect even after all retries.", myHostPort, " aborting.")
					return nil, errors.New("Couldn't dial a node")
				}
				/* Try again after 1 second */
				time.Sleep(1 * time.Second)
				/* Make RPC call */
				dialer, err := rpc.DialHTTP("tcp", v)
				if err == nil {
					/* Success! Cache this connection and continue */
					node.slaveDialerMap[k] = dialer
					break
				}
			}
		} else {
			/* Success! Cache this connection and continue */
			node.slaveDialerMap[k] = dialer
		}

		fmt.Println(myHostPort, " dialed slave ", v, " successfully")
	}

	/* Create a Goroutine to send HeartBeat to the Monitor node */
	go (&node).HeartBeat(monitorHostPort)

	/* Return the struct */
	return a, nil
}

/* Returns the next proposal number for that key */
func (pn *paxosNode) GetNextProposalNumber(args *paxosrpc.ProposalNumberArgs, reply *paxosrpc.ProposalNumberReply) error {
	fmt.Println("GetNextProposalNumber invoked on ", pn.srvId)
	key := args.Key

	/* increase the nextNum for this key, append distinct srvId */
	pn.nextSeqNumMapLock.Lock()
	defer pn.nextSeqNumMapLock.Unlock()

	pn.nextSeqNumMap[key] += 1
	nextNum := pn.nextSeqNumMap[key]
	nextNum = nextNum*1000 + pn.srvId
	reply.N = nextNum
	return nil
}

/* Structs to send on the channel during prepare and accept phases */
type prepReplyAndTimeout struct {
	prepReply paxosrpc.PrepareReply
	timeout   int
}

type accReplyAndTimeout struct {
	accReply paxosrpc.AcceptReply
	timeout  int
}

/* Prepare calls the RecvPrepare RPC on the paxos node given by srvId for the sequence number
   seqnum, and the key given by key
*/
func prepare(pn *paxosNode, srvId int, key string, seqnum int, preparechan chan prepReplyAndTimeout) {
	args := paxosrpc.PrepareArgs{}
	args.Key = key
	args.N = seqnum

	reply := paxosrpc.PrepareReply{}

	/* Call the RecvPrepare RPC */
	err := pn.paxosDialerMap[srvId].Call("PaxosNode.RecvPrepare", &args, &reply)

	if err != nil {
		fmt.Println("RPC RecvPrepare failed!")
		return
	}

	var ret prepReplyAndTimeout
	/* If you get a reply, send it over the channel */
	ret.prepReply = reply
	/* Since this came in from the paxos node, timeout is set to 0 */
	ret.timeout = 0
	fmt.Println("Got Prepare reply from ", pn.hostMap[srvId], ". The N is ", reply.N_a, " and the value is ", reply.V_a)
	/* Send away... */
	preparechan <- ret
}

/* accept calls the RecvAccept RPC on the paxos node given by srvId for the sequence number
   seqnum, and the key given by key, and the value given by value
*/
func accept(pn *paxosNode, srvId int, value [NumCopies]int, key string, seqnum int, acceptchan chan accReplyAndTimeout) {
	args := paxosrpc.AcceptArgs{}
	args.Key = key
	args.N = seqnum
	args.V = value

	reply := paxosrpc.AcceptReply{}

	/* Call the RecvAccept RPC */
	err := pn.paxosDialerMap[srvId].Call("PaxosNode.RecvAccept", &args, &reply)

	if err != nil {
		fmt.Println("RPC RecvAccept failed!")
		return
	}

	var ret accReplyAndTimeout
	/* If you get a reply, send it over the channel */
	ret.accReply = reply
	/* Since this came in from the paxos node, timeout is set to 0 */
	ret.timeout = 0
	fmt.Println("Got Accept reply from ", pn.hostMap[srvId], ". The Status is ", reply.Status)
	/* Send away... */
	acceptchan <- ret
}

/* commit calls the RecvCommit RPC on the paxos node given by srvId, and the key given by key
   and the value given by value.
*/
func commit(pn *paxosNode, srvId int, value [NumCopies]int, key string, commitchan chan int) {
	args := paxosrpc.CommitArgs{}
	args.Key = key
	args.V = value

	reply := paxosrpc.CommitReply{}

	/* Call the RecvCommit RPC */
	err := pn.paxosDialerMap[srvId].Call("PaxosNode.RecvCommit", &args, &reply)

	if err != nil {
		fmt.Println("RPC RecvCommit failed!")
		return
	}

	/* Just write 0 onto the channel */
	fmt.Println("Got Commit reply from ", pn.hostMap[srvId])
	commitchan <- 0
}

/* This go routine writes onto the channels that the timeout has occured, and that the
   propose phase must be aborted
*/
func wakeMeUpAfter15Seconds(preparechan chan prepReplyAndTimeout, acceptchan chan accReplyAndTimeout,
	commitchan chan int) {
	/* sleep for 15 seconds */
	time.Sleep(15 * time.Second)

	var ret1 prepReplyAndTimeout
	/* set timeout field to 1 and send away */
	ret1.timeout = 1
	preparechan <- ret1

	var ret2 accReplyAndTimeout
	/* set timeout field to 1 and send away */
	ret2.timeout = 1
	acceptchan <- ret2

	/* 1 indicates timeout hit */
	commitchan <- 1
}

/* a struct for finding the least loaded NumCopies slave nodes */
type idAndSize struct {
	id   int
	size int
}

/* PutRank puts a key-value pair onto the slave nodes, where the key is a URL, and value is pagerank */
func (pn *paxosNode) PutRank(args *paxosrpc.PutRankArgs, reply *paxosrpc.PutRankReply) error {
	fmt.Println("PutRank called on ", pn.myHostPort, " with key = ", args.Key)
	var putargs slaverpc.PutRankArgs
	putargs.Key = args.Key
	putargs.Value = args.Value

	var putreply slaverpc.PutRankReply

	/* see if there is already a mapping */
	slaveList, ok := pn.valuesMap[args.Key]

	if ok {
		/* There is a mapping. Directly call PutRank on all those node */
		fmt.Println("This key already has some slaves. Calling PutRank on them")

		/* Iterate over the list of slave nodes associated with this key */
		for _, slaveId := range slaveList {
			err := pn.slaveDialerMap[slaveId].Call("SlaveNode.PutRank", &putargs, &putreply)
			if err != nil {
				fmt.Println("PutRank RPC Failed on ", pn.slaveMap[slaveId])
			}
		}
		return nil
	}

	var propNumArgs paxosrpc.ProposalNumberArgs
	var propNumReply paxosrpc.ProposalNumberReply
	propNumArgs.Key = args.Key

	/* Get next proposal number for this key */
	pn.GetNextProposalNumber(&propNumArgs, &propNumReply)

	fmt.Println("Got proposal number as ", propNumReply.N)

	var idAndSizeList []idAndSize

	/* now select a group of slaves to store this value on */
	/* first create a list of slaveid-size pairs */
	for i := 0; i < pn.numSlaves; i++ {
		key := strconv.Itoa(i)
		key = key + ":size"
		size := (pn.valuesMap[key])[0]

		ias := idAndSize{id: i, size: size}
		idAndSizeList = append(idAndSizeList, ias)
	}

	/* sort the list based on size */
	sort.Sort(SlaveSlice(idAndSizeList))
	i := 0

	/* take the first NumCopies servers, which are the least loaded ones */
	var targetSlaves [NumCopies]int
	for i < NumCopies {
		targetSlaves[i] = idAndSizeList[i].id
		i += 1
	}

	fmt.Println("Generated the targetSlaves slice as ", targetSlaves)

	var proposeArgs paxosrpc.ProposeArgs
	var proposeReply paxosrpc.ProposeReply

	proposeArgs.N = propNumReply.N
	proposeArgs.Key = args.Key
	proposeArgs.V = targetSlaves

	fmt.Println(pn.myHostPort, " will now propose for key = ", proposeArgs.Key, " and value = ", proposeArgs.V)

	/* now propose this list of slaves */
	err := pn.Propose(&proposeArgs, &proposeReply)
	if err != nil {
		fmt.Println("Propose failed!")
		return errors.New("Propose failed!")
	}

	fmt.Println("Propose succeeded and the value committed was ", proposeReply.V)

	/* accept whatever was committed and call PutRank on all those slaves */
	for _, slaveId := range proposeReply.V {
		/* Make the RPC */
		err := pn.slaveDialerMap[slaveId].Call("SlaveNode.PutRank", &putargs, &putreply)
		if err != nil {
			fmt.Println("PutRank RPC Failed on ", pn.slaveMap[slaveId])
		}
	}

	return nil
}

/* This returns the page rank for a given URL */
func (pn *paxosNode) GetRank(args *paxosrpc.GetRankArgs, reply *paxosrpc.GetRankReply) error {
	fmt.Println("GetRank invoked on Node ", pn.srvId, " for key ", args.Key)

	/* See if it is present */
	_, ok := pn.valuesMap[args.Key]

	if !ok {
		/* Not there */
		fmt.Println("Key not found")
		return errors.New("Key not found")
	}

	/* Key found! So the page rank has to be present! Find and return it */
	for _, slaveId := range pn.valuesMap[args.Key] {
		var getRankArgs slaverpc.GetRankArgs
		getRankArgs.Key = args.Key

		var getRankReply slaverpc.GetRankReply
		/* Make the GetRank call */
		err := pn.slaveDialerMap[slaveId].Call("SlaveNode.GetRank", &getRankArgs, &getRankReply)
		if err == nil {
			/* Found on this slave node. Return the value */
			reply.Value = getRankReply.Value
			fmt.Println("Slave ", slaveId, " has the data for ", args.Key, "!")
			return nil
		}
	}
	/* None of the slaves replied */
	return errors.New("No slave replied.")
}

/* GetLinks returns the crawled data and the corresponding links for a given URL */
func (pn *paxosNode) GetLinks(args *paxosrpc.GetLinksArgs, reply *paxosrpc.GetLinksReply) error {
	fmt.Println("GetLinks invoked on Node ", pn.srvId, " for key ", args.Key)
	/* See if it is present in the map */
	list, ok := pn.valuesMap[args.Key]
	if !ok {
		/* Not there, return error */
		fmt.Println("Key", args.Key, "not found")
		return errors.New("Key " + args.Key + " not found")
	}

	/* iterate over the slavelist, and call the RPC on each one of them one after the other
	   until one responds */
	for _, slaveId := range list {
		var getArgs slaverpc.GetArgs
		getArgs.Key = args.Key

		var getReply slaverpc.GetReply
		/* Call RPC */
		err := pn.slaveDialerMap[slaveId].Call("SlaveNode.Get", &getArgs, &getReply)
		if err == nil && getReply.Status == 1 {
			reply.Value = getReply.Value
			fmt.Println("Slave ", slaveId, " has the data for ", args.Key, "!")
			return nil
		}
	}
	return errors.New("No slave replied.")
}

/* GetAllLinks returns ALL the crawled data and the corresponding links */
func (pn *paxosNode) GetAllLinks(args *paxosrpc.GetAllLinksArgs, reply *paxosrpc.GetAllLinksReply) error {
	fmt.Println("GetAllLinks invoked on Node ", pn.srvId)
	/* make a map
	   The key is the URL, and value is the crawled data for that URL
	*/
	reply.LinksMap = make(map[string][]string)
	/* iterate over the whole values map */
	for key, value := range pn.valuesMap {
		var getArgs slaverpc.GetArgs
		getArgs.Key = key

		var getReply slaverpc.GetReply
		/* Ask the slaves only if the key is not for the sizes! */
		if !strings.HasSuffix(key, ":size") {
			for _, slaveId := range value {
				/* Call RPC */
				err := pn.slaveDialerMap[slaveId].Call("SlaveNode.Get", &getArgs, &getReply)
				if err == nil && getReply.Status == 1 {
					reply.LinksMap[key] = getReply.Value
					break
				}
			}
		}
	}
	return nil
}

/* Do a paxos round to help update size metadata on all paxos nodes */
func (pn *paxosNode) updateSizes(url string, value []string) error {
	sliceSize := 0

	/* First get the size of the string slice */
	for _, str := range value {
		sliceSize += len(str)
	}

	var propNumArgs paxosrpc.ProposalNumberArgs
	var propNumReply paxosrpc.ProposalNumberReply

	/* Make one paxos round for updating the size of each slave node in the mapping */
	for _, i := range pn.valuesMap[url] {
		key := strconv.Itoa(i)
		key = key + ":size"

		propNumArgs.Key = key

		for {
			/* Get proposal number for proposing this update */
			pn.GetNextProposalNumber(&propNumArgs, &propNumReply)

			fmt.Println("For updating size, got next proposal number as ", propNumReply.N)

			var proposeArgs paxosrpc.ProposeArgs
			var proposeReply paxosrpc.ProposeReply

			var newSizeSlice [NumCopies]int
			/* the value is nothing but the current size added to the size of this string slice */
			newSizeSlice[0] = (pn.valuesMap[key])[0] + sliceSize

			proposeArgs.N = propNumReply.N
			proposeArgs.Key = key
			proposeArgs.V = newSizeSlice

			fmt.Println(pn.myHostPort, " will now propose for key = ", proposeArgs.Key, " and value = ", proposeArgs.V)

			/* propose this value */
			pn.Propose(&proposeArgs, &proposeReply)

			/* if the value was proposed successfully, break and continue paxos for other slaves */
			if proposeReply.Status == paxosrpc.OK {
				fmt.Println("Propose succeeded and the value committed was ", proposeReply.V)
				break
			} else if proposeReply.Status == paxosrpc.OtherValueCommitted { /* If some other value was
				committed, try again */
				fmt.Println("Some other value was committed! Retry")
			} else {
				/* If rejected, try again */
				fmt.Println("Paxos rejected! Retry")
			}
		}
	}

	return nil
}

/* Append the given string slice to the given key URL */
func (pn *paxosNode) Append(args *paxosrpc.AppendArgs, reply *paxosrpc.AppendReply) error {
	fmt.Println("Append called on ", pn.myHostPort, " with key = ", args.Key)
	var appendArgs slaverpc.AppendArgs
	appendArgs.Key = args.Key
	appendArgs.Value = args.Value

	var appendReply slaverpc.AppendReply
	/* see if there is already a list of slaves for this key */
	slaveList, ok := pn.valuesMap[args.Key]

	if ok {
		/* this key already has a list of dedicated slaves */
		fmt.Println("This key already has some slaves. First update the sizes.")

		/* first do paxos rounds to update sizes of all appropriate slaves */
		pn.updateSizes(args.Key, args.Value)

		fmt.Println("Size updation done! Now call Append on all")

		/* Now just call append on all */
		for _, slaveId := range slaveList {
			err := pn.slaveDialerMap[slaveId].Call("SlaveNode.Append", &appendArgs, &appendReply)
			if err != nil {
				fmt.Println("Append RPC Failed on ", pn.slaveMap[slaveId])
			}
		}
		return nil
	}

	var propNumArgs paxosrpc.ProposalNumberArgs
	var propNumReply paxosrpc.ProposalNumberReply
	propNumArgs.Key = args.Key

	/* No list of slaves found for this key.
	   Hence, we must initiate a paxos round to propose a list of dedicated slaves for this key */
	pn.GetNextProposalNumber(&propNumArgs, &propNumReply)

	fmt.Println("Got proposal number as ", propNumReply.N)

	var idAndSizeList []idAndSize
	/* now select a group of slaves to store this value on */
	for i := 0; i < pn.numSlaves; i++ {
		key := strconv.Itoa(i)
		key = key + ":size"
		size := (pn.valuesMap[key])[0]

		ias := idAndSize{id: i, size: size}
		idAndSizeList = append(idAndSizeList, ias)
	}

	/* sort the list of slaves based on size */
	sort.Sort(SlaveSlice(idAndSizeList))
	i := 0

	/* get the first NumCopies least loaded slaves */
	var targetSlaves [NumCopies]int
	for i < NumCopies {
		targetSlaves[i] = idAndSizeList[i].id
		i += 1
	}

	fmt.Println("Generated the targetSlaves slice as ", targetSlaves)

	var proposeArgs paxosrpc.ProposeArgs
	var proposeReply paxosrpc.ProposeReply

	proposeArgs.N = propNumReply.N
	proposeArgs.Key = args.Key
	proposeArgs.V = targetSlaves

	fmt.Println(pn.myHostPort, " will now propose for key = ", proposeArgs.Key, " and value = ", proposeArgs.V)

	/* now initiate a paxos round to propose this list of slaves */
	err := pn.Propose(&proposeArgs, &proposeReply)
	if err != nil {
		fmt.Println("Propose failed!")
		return errors.New("Propose failed!")
	}

	fmt.Println("Propose succeeded and the value committed was ", proposeReply.V)

	fmt.Println("Now will try to propose size values")

	/* now again do paxos to update the sizes of all associated slaves */
	pn.updateSizes(args.Key, args.Value)

	fmt.Println("Size updation done! Now call Append on all appropriate nodes")
	/* Now call append on all appropriate slave nodes */
	for _, slaveId := range proposeReply.V {
		err := pn.slaveDialerMap[slaveId].Call("SlaveNode.Append", &appendArgs, &appendReply)
		if err != nil {
			fmt.Println("Append RPC Failed on ", pn.slaveMap[slaveId])
		}
	}

	reply.Status = paxosrpc.OK
	return nil

}

/* Initiate a paxos round to propose a given key-value pair with the specified proposal number */
func (pn *paxosNode) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error {
	preparechan := make(chan prepReplyAndTimeout, 100)
	acceptchan := make(chan accReplyAndTimeout, 100)
	commitchan := make(chan int, 100)

	fmt.Println("In Propose of ", pn.srvId)

	fmt.Println("Key is ", args.Key, ", V is ", args.V, " and N is ", args.N)

	/* create a go routine to tell this propose to quit after 15 seconds */
	go wakeMeUpAfter15Seconds(preparechan, acceptchan, commitchan)

	/* call prepare on all master nodes */
	for k, v := range pn.hostMap {
		fmt.Println("Will call Prepare on ", v)
		go prepare(pn, k, args.Key, args.N, preparechan)
	}

	okcount := 0

	max_n := 0
	var max_v [NumCopies]int
	max_v = args.V

	/* wait for reply from each paxos node */
	reply.Status = paxosrpc.OK
	for i := 0; i < pn.numNodes; i++ {
		ret, _ := <-preparechan
		if ret.timeout == 1 { /* timeout hit, abort */
			fmt.Println("Didn't finish prepare stage even after 15 seconds")
			reply.Status = paxosrpc.Reject
			return errors.New("Didn't finish prepare stage even after 15 seconds")
		}
		/* if prepare accepted, increment counter */
		if ret.prepReply.Status == paxosrpc.OK {
			okcount++
		}
		/* take into account whatever was returned by that node in the reply for prepare */
		if ret.prepReply.N_a != 0 && ret.prepReply.N_a > max_n {
			max_n = ret.prepReply.N_a
			max_v = ret.prepReply.V_a
			reply.Status = paxosrpc.OtherValueCommitted
		}
		if okcount >= ((pn.numNodes / 2) + 1) { /* break out early if a majority is reached */
			break
		}
	}

	/* if we didn't get a majority in prepare accept, abort */
	if !(okcount >= ((pn.numNodes / 2) + 1)) {
		reply.Status = paxosrpc.Reject
		return errors.New("Didn't get a majority in prepare phase")
	}

	/* now call accept on all other paxos masters */
	for k, v := range pn.hostMap {
		fmt.Println("Will call Accept on ", v)
		go accept(pn, k, max_v, args.Key, args.N, acceptchan)
	}

	okcount = 0

	/* wait for reply from all paxos nodes */
	for i := 0; i < pn.numNodes; i++ {
		ret, _ := <-acceptchan
		if ret.timeout == 1 { /* timeout hit, abort */
			fmt.Println("Didn't finish accept stage even after 15 seconds")
			reply.Status = paxosrpc.Reject
			return errors.New("Didn't finish accept stage even after 15 seconds")
		}
		/* if accept accepted, increment counter */
		if ret.accReply.Status == paxosrpc.OK {
			okcount++
		}
		if okcount >= ((pn.numNodes / 2) + 1) { /* break out earlier if majority have accepted */
			break
		}
	}

	/* If didn't get majority, abort */
	if !(okcount >= ((pn.numNodes / 2) + 1)) {
		reply.Status = paxosrpc.Reject
		return errors.New("Didn't get a majority in accept phase")
	}

	okcount = 0
	/* now call commit on all nodes */
	for k, v := range pn.hostMap {
		fmt.Println("Will call Commit on ", v)
		go commit(pn, k, max_v, args.Key, commitchan)
	}

	/* now we must wait for reply from all nodes, and we cannot return earlier like in previous cases */
	for i := 0; i < pn.numNodes; i++ {
		_, ok := <-commitchan
		if !ok {
			fmt.Println("Didn't finish commit stage even after 15 seconds")
			reply.Status = paxosrpc.Reject
			return errors.New("Didn't finish commit stage even after 15 seconds")
		}
	}

	reply.V = max_v

	fmt.Println(pn.myHostPort, " has successfully proposed the set ", max_v)
	return nil
}

/* Receive a prepare from another paxos node */
func (pn *paxosNode) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	defer fmt.Println("Leaving RecvPrepare of ", pn.myHostPort)
	fmt.Println("In RecvPrepare of ", pn.myHostPort)
	key := args.Key
	num := args.N

	/* get the maximum sequence number seen for this key */
	pn.maxSeqNumSoFarLock.Lock()
	maxNum := pn.maxSeqNumSoFar[key]
	pn.maxSeqNumSoFarLock.Unlock()

	/* see if there's any value that was already accepted, but hasn't yet been committed */
	pn.acceptedValuesMapLock.Lock()
	val, ok := pn.acceptedValuesMap[key]
	pn.acceptedValuesMapLock.Unlock()

	pn.acceptedSeqNumMapLock.Lock()
	seqNum := pn.acceptedSeqNumMap[key]
	pn.acceptedSeqNumMapLock.Unlock()

	if !ok {
		reply.N_a = -1 /* none accepted so far */
	} else {
		/* a value has been accepted previously, tell the proposer accordingly */
		reply.V_a = val
		reply.N_a = seqNum
	}

	/* reject proposal when its proposal number is not higher than the highest number it's ever seen */
	if maxNum > num {
		fmt.Println("In RecvPrepare of ", pn.myHostPort, "rejected proposal:", key, num, "maxNum:", maxNum)
		reply.Status = paxosrpc.Reject

		return nil
	}
	/* promise proposal when its higher. return with the number and value accepted */
	fmt.Println("In RecvPrepare of ", pn.myHostPort, "accepted proposal:", key, num, "maxNum:", maxNum)

	pn.maxSeqNumSoFarLock.Lock()
	pn.maxSeqNumSoFar[key] = num
	pn.maxSeqNumSoFarLock.Unlock()

	/* fill reply with accepted seqNum and value. default is -1 */
	reply.Status = paxosrpc.OK
	return nil
}

/* receive an accept RPC from a proposer */
func (pn *paxosNode) RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	defer fmt.Println("Leaving RecvAccept of ", pn.myHostPort)
	fmt.Println("In RecvAccept of ", pn.myHostPort)
	key := args.Key
	num := args.N
	value := args.V

	/* see what's the maximum sequence number seen so far this key */
	pn.maxSeqNumSoFarLock.Lock()
	maxNum := pn.maxSeqNumSoFar[key]
	pn.maxSeqNumSoFarLock.Unlock()

	/* reject proposal when its proposal number is not higher than the highest number it's ever seen */
	if maxNum > num {
		fmt.Println("In RecvAccept of ", pn.myHostPort, "rejected proposal:", key, num, value, "maxNum:", maxNum)
		reply.Status = paxosrpc.Reject
		return nil
	}
	/* accept proposal when its higher. update with the number and value accepted */
	fmt.Println("In RecvAccept of ", pn.myHostPort, "accepted proposal:", key, num, value)

	/* accept all maps */
	pn.maxSeqNumSoFarLock.Lock()
	pn.maxSeqNumSoFar[key] = num
	pn.maxSeqNumSoFarLock.Unlock()

	pn.acceptedValuesMapLock.Lock()
	pn.acceptedValuesMap[key] = value
	pn.acceptedValuesMapLock.Unlock()

	pn.acceptedSeqNumMapLock.Lock()
	pn.acceptedSeqNumMap[key] = num
	pn.acceptedSeqNumMapLock.Unlock()

	reply.Status = paxosrpc.OK
	return nil
}

/* receive a commit message from another paxos node */
func (pn *paxosNode) RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	defer fmt.Println("Leaving RecvCommit of ", pn.myHostPort)
	key := args.Key
	value := args.V

	// update the value and clear the map for accepted value and number
	fmt.Println("In RecvCommit of ", pn.myHostPort, "committing:", key, value)
	pn.valuesMapLock.Lock()
	pn.valuesMap[key] = value
	pn.valuesMapLock.Unlock()

	/* delete the entries from accepted maps */
	pn.acceptedValuesMapLock.Lock()
	delete(pn.acceptedValuesMap, key)
	pn.acceptedValuesMapLock.Unlock()

	pn.acceptedSeqNumMapLock.Lock()
	delete(pn.acceptedSeqNumMap, key)
	pn.acceptedSeqNumMapLock.Unlock()
	return nil
}

/* receive a replace server update from another paxos node which is trying to replace a failed node */
func (pn *paxosNode) RecvReplaceServer(args *paxosrpc.ReplaceServerArgs, reply *paxosrpc.ReplaceServerReply) error {
	fmt.Println("In ReplaceServer of ", pn.myHostPort)
	fmt.Println("Some node is replacing the node with ID : ", args.SrvID)
	fmt.Println("Its hostport is ", args.Hostport)

	dialer, _ := rpc.DialHTTP("tcp", args.Hostport)

	/* update the host map and cached dialer */
	pn.hostMap[args.SrvID] = args.Hostport
	pn.paxosDialerMap[args.SrvID] = dialer
	return nil
}

/* receive a request from a node trying to replace a failed paxos node to update it with the values map */
func (pn *paxosNode) RecvReplaceCatchup(args *paxosrpc.ReplaceCatchupArgs, reply *paxosrpc.ReplaceCatchupReply) error {
	fmt.Println("In RecvReplaceCatchup of ", pn.myHostPort)

	pn.valuesMapLock.Lock()
	fmt.Println("Marshalling", len(pn.valuesMap), "values for RecvReplaceCatchup")
	/* marshal the values map */
	marshaledMap, err := json.Marshal(pn.valuesMap)
	pn.valuesMapLock.Unlock()

	if err != nil {
		fmt.Println("Failed to marshall", err)
		return err
	}
	/* reply with the marshaled map */
	reply.Data = marshaledMap
	return nil
}

/* call heartbeat on the monitor node */
func (pn *paxosNode) HeartBeat(hostPort string) {
	fmt.Println("Heartbeat invoked on Master node", pn.srvId)
	for {
		/* first try dialing the monitor node */
		time.Sleep(time.Second * 2)
		fmt.Println("Try dialing Monitor node from Master", pn.srvId)
		var monitor common.Conn
		hostPorts := make([]string, 0)
		hostPorts = append(hostPorts, hostPort)
		monitorDialer, err := rpc.DialHTTP("tcp", hostPort)
		monitor.HostPort = hostPorts
		monitor.Dialer = monitorDialer
		if err != nil { /* couldn't dial monitor, try again */
			fmt.Println(err)
		} else { /* success, cache this connection and continue */
			pn.monitor = monitor
			break
		}
	}

	for {
		var args monitorrpc.HeartBeatArgs
		var reply monitorrpc.HeartBeatReply
		args.Id = pn.srvId
		args.Type = monitorrpc.Master
		//fmt.Println("Calling MonitorNode.HeartBeat to monitor node")
		/* call heartbeat on the monitor node to tell it that this node is alive */
		err := pn.monitor.Dialer.Call("MonitorNode.HeartBeat", &args, &reply)
		if err != nil {
			fmt.Println(err)
		}
		/* sleep for 3 seconds, and send again */
		time.Sleep(time.Second * 3)
	}
}

/* methods to sort a given list of slaves based on size */
type SlaveSlice []idAndSize

func (slice SlaveSlice) Len() int {
	return len(slice)
}

func (slice SlaveSlice) Less(i, j int) bool {
	return slice[i].size < slice[j].size
}

func (slice SlaveSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
