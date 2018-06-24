package storageserver

import (
	"errors"
	"fmt"
	"github.com/cmu440/tribbler/rpc/storagerpc"
	"hash/fnv"
	"net"
	"net/http"
	"net/rpc"
	"sort"
	"strings"
	"sync"
	"time"
)

type storageServer struct {
	numNodes             int
	nodeID               uint32
	port                 int
	masterServerHostPort string
	tribbleMap           map[string]string   //need lock
	listMap              map[string][]string //need lock
	ackedSlaves          int
	serverList           []storagerpc.Node
	ackedSlavesMap       map[storagerpc.Node]bool //need lock
	//need lock
	leaseMap map[string][]string //For keeping a track of which key is cached by which libstore
	//TODO: I think, can combine revokeKeysmap and pendingPuts together
	revokeKeysMap map[string]bool     //need lock
	putBlockChans map[string]chan int //need lock
	pendingPuts   map[string]bool     //need lock

	putLock     *sync.Mutex
	putLocksMap map[string]*sync.Mutex

	appendLock     *sync.Mutex
	appendLocksMap map[string]*sync.Mutex

	allSlavesAcked chan int

	minhash uint32
	maxhash uint32

	first int
}

// NewStorageServer creates and starts a new StorageServer. masterServerHostPort
// is the master storage server's host:port address. If empty, then this server
// is the master; otherwise, this server is a slave. numNodes is the total number of
// servers in the ring. port is the port number that this server should listen on.
// nodeID is a random, unsigned 32-bit ID identifying this server.
//
// This function should return only once all storage servers have joined the ring,
// and should return a non-nil error if the storage server could not be started.

func NewStorageServer(masterServerHostPort string, numNodes, port int, nodeID uint32) (StorageServer, error) {
	defer fmt.Println("Leaving NewStorageServer")
	fmt.Println("Entered NewStorageServer")

	var a StorageServer

	server := storageServer{}

	server.numNodes = numNodes
	server.port = port
	server.masterServerHostPort = masterServerHostPort
	server.nodeID = nodeID
	server.tribbleMap = make(map[string]string)
	server.listMap = make(map[string][]string)
	server.ackedSlavesMap = make(map[storagerpc.Node]bool)
	server.leaseMap = make(map[string][]string)
	server.revokeKeysMap = make(map[string]bool)
	server.putBlockChans = make(map[string]chan int)
	server.pendingPuts = make(map[string]bool)
	server.putLocksMap = make(map[string]*sync.Mutex)
	server.appendLocksMap = make(map[string]*sync.Mutex)
	server.allSlavesAcked = make(chan int, 100)
	//server.serverList = make([]storagerpc.Node, 32)

	server.putLock = &sync.Mutex{}
	server.appendLock = &sync.Mutex{}

	a = &server

	if len(masterServerHostPort) == 0 {
		fmt.Println("This is the master speaking, with numNode = ", numNodes, " and port = ", port)
		/* Now register for RPCs */
		server.ackedSlaves = 1
		self := storagerpc.Node{fmt.Sprintf("localhost:%d", port), nodeID}

		server.serverList = append(server.serverList, self)

		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			return nil, err
		}

		err = rpc.RegisterName("StorageServer", storagerpc.Wrap(a))
		if err != nil {
			return nil, err
		}

		rpc.HandleHTTP()
		go http.Serve(listener, nil)

		/* Now wait until all nodes have joined */
		fmt.Println("Master waiting for all children")
		if numNodes != 1 {
			_ = <-server.allSlavesAcked
		}
		fmt.Println("Master heard from all children")
		sort.Sort(ServerSlice(server.serverList))

		i := 0

		for server.serverList[i].NodeID != nodeID {
			i = i + 1
		}
		if i == 0 {
			server.minhash = server.serverList[len(server.serverList)-1].NodeID
			server.maxhash = nodeID
			server.first = 1
		} else {
			server.maxhash = nodeID
			server.minhash = server.serverList[i-1].NodeID
			server.first = 0
		}

	} else {
		fmt.Println("I am just a lowly Slave.")
		fmt.Println("Number of nodes in the client is ", numNodes, " and the port is ", port, " and the nodeID is ", nodeID, " and the masterserverhostport is ", masterServerHostPort)
		/* Now try connecting to the ring by calling the RegisterServer RPC */

		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			fmt.Println("Returning from NewStorageServer because couldn't listen on given port")
			return nil, err
		}

		err = rpc.RegisterName("StorageServer", storagerpc.Wrap(a))
		if err != nil {
			fmt.Println("Returning from NewStorageServer because couldn't register rpc")
			return nil, err
		}

		rpc.HandleHTTP()
		go http.Serve(listener, nil)

		srvr, err := rpc.DialHTTP("tcp", masterServerHostPort)
		for {
			if err != nil {
				fmt.Println("Oops! Returning because couldn't dial master host port. Let's try after 1 sec")
				//return nil, errors.New("Couldn't Dial Master Host Port")
			} else {
				break
			}
			time.Sleep(1 * time.Second)
			srvr, err = rpc.DialHTTP("tcp", masterServerHostPort)
		}

		args := storagerpc.RegisterArgs{}
		args.ServerInfo.HostPort = fmt.Sprintf("localhost:%d", port)
		args.ServerInfo.NodeID = nodeID

		var reply storagerpc.RegisterReply

		for {
			err := srvr.Call("StorageServer.RegisterServer", args, &reply)
			if err != nil {
				fmt.Println("Returning from NewStorageServer because Call returned : ", err)
				return nil, err
			}
			if reply.Status == storagerpc.OK {
				server.serverList = reply.Servers
				break
			}
			time.Sleep(1 * time.Second)
		}
		sort.Sort(ServerSlice(server.serverList))
		i := 0

		for server.serverList[i].NodeID != nodeID {
			i = i + 1
		}
		if i == 0 {
			server.minhash = server.serverList[len(server.serverList)-1].NodeID
			server.maxhash = nodeID
			server.first = 1
		} else {
			server.maxhash = nodeID
			server.minhash = server.serverList[i-1].NodeID
			server.first = 0
		}

	}

	return a, nil
}

func (ss *storageServer) RegisterServer(args *storagerpc.RegisterArgs, reply *storagerpc.RegisterReply) error {
	defer fmt.Println("Leaving RegisterServer")
	fmt.Println("RegisterServer invoked! By: ", args.ServerInfo.HostPort)

	if _, ok := ss.ackedSlavesMap[args.ServerInfo]; ok {
		if ss.ackedSlaves == ss.numNodes {
			reply.Status = storagerpc.OK
			//TODO: Sort this list!
			reply.Servers = ss.serverList
			ss.allSlavesAcked <- 1
			return nil
		}
		reply.Status = storagerpc.NotReady
		return nil
	}

	ss.serverList = append(ss.serverList, args.ServerInfo)
	ss.ackedSlaves += 1
	ss.ackedSlavesMap[args.ServerInfo] = true

	if ss.ackedSlaves == ss.numNodes {
		//TODO: Sort this list!
		reply.Status = storagerpc.OK
		reply.Servers = ss.serverList
		ss.allSlavesAcked <- 1
	} else {
		reply.Status = storagerpc.NotReady
	}
	return nil
}

func (ss *storageServer) GetServers(args *storagerpc.GetServersArgs, reply *storagerpc.GetServersReply) error {
	defer fmt.Println("Leaving GetServers")
	fmt.Println("GetServers invoked!")

	if ss.ackedSlaves == ss.numNodes {
		reply.Status = storagerpc.OK
		reply.Servers = ss.serverList
		fmt.Println("Returning OK")
		fmt.Println("The list of servers is ", reply.Servers)
		return nil
	}
	reply.Status = storagerpc.NotReady
	fmt.Println("Returning NotReady")
	return nil
}

func (ss *storageServer) Get(args *storagerpc.GetArgs, reply *storagerpc.GetReply) error {
	/*defer fmt.Println("Leaving Get")
	fmt.Println("Get invoked!")*/
	/*fmt.Println("Key is ", args.Key)
	fmt.Println("The hash is ", StoreHash(args.Key), ", my nodeID is ", ss.nodeID, " and my port is ", ss.port)*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("Returning because hash doesn't match")
			fmt.Println("Server list is ", ss.serverList)
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("Returning because hash doesn't match")
			fmt.Println("Server list is ", ss.serverList)
			return nil
		}
	}

	//fmt.Println("Key is ", args.Key, ", WantLease is ", args.WantLease, " and HostPort is ", args.HostPort)

	val, ok := ss.tribbleMap[args.Key]

	if !ok {
		reply.Status = storagerpc.KeyNotFound
		return nil
	}

	reply.Status = storagerpc.OK
	reply.Value = val

	//If revoking in process, ok will be true, and don't give lease!
	_, ok = ss.revokeKeysMap[args.Key]

	if args.WantLease == true && !ok {
		if _, ok := ss.leaseMap[args.Key]; !ok {
			var templist []string
			templist = append(templist, args.HostPort)
			ss.leaseMap[args.Key] = templist
		} else {
			ss.leaseMap[args.Key] = append(ss.leaseMap[args.Key], args.HostPort)
		}
		reply.Status = storagerpc.OK
		reply.Lease.Granted = true
		reply.Lease.ValidSeconds = storagerpc.LeaseSeconds
		go handleLeaseExpire(ss, storagerpc.LeaseSeconds, args.Key, args.HostPort)
	} else {
		reply.Lease.Granted = false
	}

	return nil
}

func (ss *storageServer) GetList(args *storagerpc.GetArgs, reply *storagerpc.GetListReply) error {
	/*defer fmt.Println("Leaving GetList")
	fmt.Println("GetList invoked!")*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			return nil
		}
	}
	//fmt.Println("Key is ", args.Key, ", WantLease is ", args.WantLease, " and HostPort is ", args.HostPort)

	val, ok := ss.listMap[args.Key]

	if !ok {
		reply.Status = storagerpc.KeyNotFound
		return nil
	}
	reply.Status = storagerpc.OK
	reply.Value = val

	//If revoking in process, ok will be true, and don't give lease!
	_, ok = ss.revokeKeysMap[args.Key]

	if args.WantLease == true && !ok {
		if _, ok := ss.leaseMap[args.Key]; !ok {
			var templist []string
			templist = append(templist, args.HostPort)
			ss.leaseMap[args.Key] = templist
		} else {
			ss.leaseMap[args.Key] = append(ss.leaseMap[args.Key], args.HostPort)
		}
		reply.Status = storagerpc.OK
		reply.Lease.Granted = true
		reply.Lease.ValidSeconds = storagerpc.LeaseSeconds
		go handleLeaseExpire(ss, storagerpc.LeaseSeconds, args.Key, args.HostPort)
	} else {
		reply.Lease.Granted = false
	}
	return nil
}

func revoke(ss *storageServer, key string, libstore *rpc.Client, hostport string) {

	args2 := storagerpc.RevokeLeaseArgs{}
	args2.Key = key

	var reply storagerpc.RevokeLeaseReply

	err := libstore.Call("LeaseCallbacks.RevokeLease", args2, &reply)
	if err != nil {
		fmt.Println("RevokeLeaseFailed")
		return
	}

	//now remove entry from leaseMap
	i := 0
	for _, val := range ss.leaseMap[key] {
		if val == hostport {
			ss.leaseMap[key] = append((ss.leaseMap[key])[:i], (ss.leaseMap[key])[i+1:]...)

			_, ok := ss.pendingPuts[key]
			// if last one and pending put, write to channel
			if len(ss.leaseMap[key]) == 0 && ok {
				ss.putBlockChans[key] <- 1
			}
			return
		}
		i = i + 1
	}

}

func (ss *storageServer) Delete(args *storagerpc.DeleteArgs, reply *storagerpc.DeleteReply) error {
	/*defer fmt.Println("Leaving Delete")
	fmt.Println("Delete invoked!")*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			return nil
		}
	}
	//fmt.Println("Key is ", args.Key)

	ss.putLock.Lock() // lock map of locks
	_, ok := ss.putLocksMap[args.Key]
	if ok {
		ss.putLocksMap[args.Key].Lock()
	} else {
		ss.putLocksMap[args.Key] = &sync.Mutex{}
		ss.putLocksMap[args.Key].Lock()
	}
	ss.putLock.Unlock()

	_, ok = ss.tribbleMap[args.Key]

	if !ok {
		reply.Status = storagerpc.KeyNotFound
		return nil
	}

	// now we know that the key is present in the map
	// now check if we need to revoke any leases

	if list, ok := ss.leaseMap[args.Key]; ok {
		ss.pendingPuts[args.Key] = true
		fmt.Println("In Delete, will be revoking leases for key ", args.Key)
		ss.revokeKeysMap[args.Key] = true //for get

		channel := make(chan int, 100)
		ss.putBlockChans[args.Key] = channel

		for _, node := range list {
			fmt.Println("Trying to revoke lease on ", node)
			libstore, err := rpc.DialHTTP("tcp", node)

			if err != nil {
				fmt.Println("Oops! Returning because couldn't dial libstore")
				return errors.New("Couldn't Dial Master Host Port")
			}

			fmt.Println("Before RevokeLease")

			go revoke(ss, args.Key, libstore, node)

			fmt.Println("After RevokeLease")
		}

		fmt.Println("Waiting for info that all revokes are done for key ", args.Key)
		_ = <-ss.putBlockChans[args.Key]
		fmt.Println("Got flag that all revokes are done for key ", args.Key)

		delete(ss.pendingPuts, args.Key)
		delete(ss.revokeKeysMap, args.Key)
		delete(ss.leaseMap, args.Key) //should be cleared also by handleLeaseTimeouts
		delete(ss.putBlockChans, args.Key)
	}

	delete(ss.tribbleMap, args.Key)
	reply.Status = storagerpc.OK
	ss.putLocksMap[args.Key].Unlock()

	return nil
}

func (ss *storageServer) Put(args *storagerpc.PutArgs, reply *storagerpc.PutReply) error {
	/*defer fmt.Println("Leaving Put")
	fmt.Println("Put invoked!")
	fmt.Println("Key is ", args.Key, " and value is ", args.Value)
	fmt.Println("The hash is ", StoreHash(args.Key), ", my nodeID is ", ss.nodeID, " and my port is ", ss.port)*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("Returning because hashes don't match")
			fmt.Println("Server list is ", ss.serverList)
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("Returning because hashes don't match")
			fmt.Println("Server list is ", ss.serverList)
			return nil
		}
	}

	ss.putLock.Lock() // lock map of locks
	_, ok := ss.putLocksMap[args.Key]
	if ok {
		ss.putLocksMap[args.Key].Lock()
	} else {
		ss.putLocksMap[args.Key] = &sync.Mutex{}
		ss.putLocksMap[args.Key].Lock()
	}
	ss.putLock.Unlock()

	//fmt.Println("Key is ", args.Key, " and value is ", args.Value)
	/* TODO: If key exists, revoke leases! */

	//Key exists in lease map!

	if list, ok := ss.leaseMap[args.Key]; ok {
		ss.pendingPuts[args.Key] = true
		fmt.Println("In Put, will be revoking leases for key ", args.Key)
		ss.revokeKeysMap[args.Key] = true //for get

		channel := make(chan int, 100)
		ss.putBlockChans[args.Key] = channel

		for _, node := range list {
			fmt.Println("Trying to revoke lease on ", node)
			libstore, err := rpc.DialHTTP("tcp", node)

			if err != nil {
				fmt.Println("Oops! Returning because couldn't dial libstore")
				return errors.New("Couldn't Dial Master Host Port")
			}

			fmt.Println("Before RevokeLease")

			go revoke(ss, args.Key, libstore, node)

			fmt.Println("After RevokeLease")
		}

		fmt.Println("Waiting for info that all revokes are done for key ", args.Key)
		_ = <-ss.putBlockChans[args.Key]
		fmt.Println("Got flag that all revokes are done for key ", args.Key)

		delete(ss.pendingPuts, args.Key)
		delete(ss.revokeKeysMap, args.Key)
		delete(ss.leaseMap, args.Key) //should be cleared also by handleLeaseTimeouts
		delete(ss.putBlockChans, args.Key)
	}

	ss.tribbleMap[args.Key] = args.Value
	reply.Status = storagerpc.OK
	ss.putLocksMap[args.Key].Unlock()

	return nil
}

func (ss *storageServer) AppendToList(args *storagerpc.PutArgs, reply *storagerpc.PutReply) error {
	/*defer fmt.Println("Leaving AppendToList")
	fmt.Println("AppendToList invoked!")*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("AppendToList returning WrongServer")
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("AppendToList returning WrongServer")
			return nil
		}
	}
	//fmt.Println("Key is ", args.Key, " and Value is ", args.Value)

	ss.appendLock.Lock() // lock map of locks
	_, ok := ss.appendLocksMap[args.Key]
	if ok {
		ss.appendLocksMap[args.Key].Lock()
	} else {
		ss.appendLocksMap[args.Key] = &sync.Mutex{}
		ss.appendLocksMap[args.Key].Lock()
	}
	ss.appendLock.Unlock()

	var templist []string

	_, ok = ss.listMap[args.Key]

	if !ok {
		templist = append(templist, args.Value)
		ss.listMap[args.Key] = templist
	} else {
		for _, val := range ss.listMap[args.Key] {
			if val == args.Value {
				reply.Status = storagerpc.ItemExists
				ss.appendLocksMap[args.Key].Unlock()
				//fmt.Println("AppendToList returning ItemExists")
				return nil
			}
		}

		// first revoke leases if any
		if list, ok := ss.leaseMap[args.Key]; ok {
			ss.pendingPuts[args.Key] = true
			fmt.Println("In AppendList, will be revoking leases for key ", args.Key)
			//TODO Change getlist to look like get
			ss.revokeKeysMap[args.Key] = true //for getlist

			channel := make(chan int, 100)
			ss.putBlockChans[args.Key] = channel

			for _, node := range list {
				fmt.Println("Trying to revoke lease on ", node)
				libstore, err := rpc.DialHTTP("tcp", node)

				if err != nil {
					fmt.Println("Oops! Returning because couldn't dial libstore")
					return errors.New("Couldn't Dial Master Host Port")
				}

				fmt.Println("Before RevokeLease")

				go revoke(ss, args.Key, libstore, node)

				fmt.Println("After RevokeLease")
			}

			fmt.Println("Waiting for info that all revokes are done for key ", args.Key)
			_ = <-ss.putBlockChans[args.Key]
			fmt.Println("Got flag that all revokes are done for key ", args.Key)

			delete(ss.pendingPuts, args.Key)
			delete(ss.revokeKeysMap, args.Key)
			delete(ss.leaseMap, args.Key) //should be cleared also by handleLeaseTimeouts
			delete(ss.putBlockChans, args.Key)
		}

		ss.listMap[args.Key] = append(ss.listMap[args.Key], args.Value)
	}

	reply.Status = storagerpc.OK
	ss.appendLocksMap[args.Key].Unlock()
	return nil
}

func (ss *storageServer) RemoveFromList(args *storagerpc.PutArgs, reply *storagerpc.PutReply) error {
	/*defer fmt.Println("Leaving RemoveFromList")
	fmt.Println("RemoveFromList invoked!")*/

	hash := StoreHash(args.Key)
	if ss.first == 1 {
		if !(hash > ss.minhash || hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("RemoveFromList returning WrongServer")
			return nil
		}
	} else {
		if !(hash > ss.minhash && hash <= ss.maxhash) {
			reply.Status = storagerpc.WrongServer
			fmt.Println("RemoveFromList returning WrongServer")
			return nil
		}
	}
	//fmt.Println("Key is ", args.Key, " and value is ", args.Value)

	ss.appendLock.Lock() // lock map of locks
	_, ok := ss.appendLocksMap[args.Key]
	if ok {
		ss.appendLocksMap[args.Key].Lock()
	} else {
		ss.appendLocksMap[args.Key] = &sync.Mutex{}
		ss.appendLocksMap[args.Key].Lock()
	}
	ss.appendLock.Unlock()

	i := 0
	for _, val := range ss.listMap[args.Key] {
		if val == args.Value {

			if list, ok := ss.leaseMap[args.Key]; ok {
				ss.pendingPuts[args.Key] = true
				fmt.Println("In RemoveFromList, will be revoking leases for key ", args.Key)
				//TODO Change getlist to look like get
				ss.revokeKeysMap[args.Key] = true //for getlist

				channel := make(chan int, 100)
				ss.putBlockChans[args.Key] = channel

				for _, node := range list {
					fmt.Println("Trying to revoke lease on ", node)
					libstore, err := rpc.DialHTTP("tcp", node)

					if err != nil {
						fmt.Println("Oops! Returning because couldn't dial libstore")
						return errors.New("Couldn't Dial Master Host Port")
					}

					fmt.Println("Before RevokeLease")

					go revoke(ss, args.Key, libstore, node)

					fmt.Println("After RevokeLease")
				}

				fmt.Println("Waiting for info that all revokes are done for key ", args.Key)
				_ = <-ss.putBlockChans[args.Key]
				fmt.Println("Got flag that all revokes are done for key ", args.Key)

				delete(ss.pendingPuts, args.Key)
				delete(ss.revokeKeysMap, args.Key)
				delete(ss.leaseMap, args.Key) //should be cleared also by handleLeaseTimeouts
				delete(ss.putBlockChans, args.Key)
			}

			ss.listMap[args.Key] = append((ss.listMap[args.Key])[:i], (ss.listMap[args.Key])[i+1:]...)
			reply.Status = storagerpc.OK
			ss.appendLocksMap[args.Key].Unlock()
			return nil
		}
		i = i + 1
	}
	reply.Status = storagerpc.ItemNotFound
	//fmt.Println("RemoveFromList returning ItemNotFound")
	ss.appendLocksMap[args.Key].Unlock()
	return nil
}

func handleLeaseExpire(ss *storageServer, seconds int, key string, hostport string) {
	defer fmt.Println("Leaving handleLeaseExpire for key ", key, " and hostport ", hostport)
	fmt.Println("HandleLeaseExpire invoked!")

	fmt.Println("Seconds is ", seconds, ", key is ", key, " and hostport is ", hostport)
	time.Sleep((storagerpc.LeaseSeconds + storagerpc.LeaseGuardSeconds) * time.Second)

	//TODO: Lock the maps before this!
	i := 0
	for _, val := range ss.leaseMap[key] {
		if val == hostport {
			ss.leaseMap[key] = append((ss.leaseMap[key])[:i], (ss.leaseMap[key])[i+1:]...)

			//write to channel if last one and pending put, otherwise simply remove entry from map
			_, ok := ss.pendingPuts[key]
			if len(ss.leaseMap[key]) == 0 {
				if ok {
					ss.putBlockChans[key] <- 1
				} else {
					delete(ss.leaseMap, key)
				}
			}
			return
		}
		i = i + 1
	}
}

type ServerSlice []storagerpc.Node

func (slice ServerSlice) Len() int {
	return len(slice)
}

func (slice ServerSlice) Less(i, j int) bool {
	return slice[i].NodeID < slice[j].NodeID
}

func (slice ServerSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func StoreHash(key string) uint32 {
	prefix := strings.Split(key, ":")[0]
	hasher := fnv.New32()
	hasher.Write([]byte(prefix))
	return hasher.Sum32()
}

/*

// This file contains constants and arguments used to perform RPCs between
// a TribServer's local Libstore and the storage servers. DO NOT MODIFY!

package storagerpc

// Status represents the status of a RPC's reply.
type Status int

const (
	OK           Status = iota + 1 // The RPC was a success.
	KeyNotFound                    // The specified key does not exist.
	ItemNotFound                   // The specified item does not exist.
	WrongServer                    // The specified key does not fall in the server's hash range.
	ItemExists                     // The item already exists in the list.
	NotReady                       // The storage servers are still getting ready.
)

// Lease constants.
const (
	QueryCacheSeconds = 10 // Time period used for tracking queries/determining whether to request leases.
	QueryCacheThresh  = 3  // If QueryCacheThresh queries in last QueryCacheSeconds, then request a lease.
	LeaseSeconds      = 10 // Number of seconds a lease should remain valid.
	LeaseGuardSeconds = 2  // Additional seconds a server should wait before invalidating a lease.
)

// Lease stores information about a lease sent from the storage servers.
type Lease struct {
	Granted      bool
	ValidSeconds int
}

type Node struct {
	HostPort string // The host:port address of the storage server node.
	NodeID   uint32 // The ID identifying this storage server node.
}

type RegisterArgs struct {
	ServerInfo Node
}

type RegisterReply struct {
	Status  Status
	Servers []Node
}

type GetServersArgs struct {
	// Intentionally left empty.
}

type GetServersReply struct {
	Status  Status
	Servers []Node
}

type GetArgs struct {
	Key       string
	WantLease bool
	HostPort  string // The Libstore's callback host:port.
}

type GetReply struct {
	Status Status
	Value  string
	Lease  Lease
}

type GetListReply struct {
	Status Status
	Value  []string
	Lease  Lease
}

type PutArgs struct {
	Key   string
	Value string
}

type PutReply struct {
	Status Status
}

type DeleteArgs struct {
	Key string
}

type DeleteReply struct {
	Status Status
}

type RevokeLeaseArgs struct {
	Key string
}

type RevokeLeaseReply struct {
	Status Status
}
*/
