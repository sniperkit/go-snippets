package libstore

import (
	"errors"
	"fmt"
	"github.com/cmu440/tribbler/rpc/librpc"
	"github.com/cmu440/tribbler/rpc/storagerpc"
	//"net"
	"time"
	//"net/http"
	"net/rpc"
	"sort"
	"sync"
)

type libstore struct {
	masterServerHostPort string
	myHostPort           string
	mode                 LeaseMode
	listServers          []storagerpc.Node
	listDialers          []*rpc.Client
	countMtx 			 *sync.Mutex
	valueMtx 			 *sync.Mutex
	listMtx	 			 *sync.Mutex
	countMap			 map[string]int
	valueMap        	 map[string]CacheLease
}

// NewLibstore creates a new instance of a TribServer's libstore. masterServerHostPort
// is the master storage server's host:port. myHostPort is this Libstore's host:port
// (i.e. the callback address that the storage servers should use to send back
// notifications when leases are revoked).
//
// The mode argument is a debugging flag that determines how the Libstore should
// request/handle leases. If mode is Never, then the Libstore should never request
// leases from the storage server (i.e. the GetArgs.WantLease field should always
// be set to false). If mode is Always, then the Libstore should always request
// leases from the storage server (i.e. the GetArgs.WantLease field should always
// be set to true). If mode is Normal, then the Libstore should make its own
// decisions on whether or not a lease should be requested from the storage server,
// based on the requirements specified in the project PDF handout.  Note that the
// value of the mode flag may also determine whether or not the Libstore should
// register to receive RPCs from the storage servers.
//
// To register the Libstore to receive RPCs from the storage servers, the following
// line of code should suffice:
//
//     rpc.RegisterName("LeaseCallbacks", librpc.Wrap(libstore))
//
// Note that unlike in the NewTribServer and NewStorageServer functions, there is no
// need to create a brand new HTTP handler to serve the requests (the Libstore may
// simply reuse the TribServer's HTTP handler since the two run in the same process).
func NewLibstore(masterServerHostPort, myHostPort string, mode LeaseMode) (Libstore, error) {
	defer fmt.Println("Leaving NewLibStore")
	fmt.Println("Entered NewLibStore")

	var a Libstore
	var b LeaseCallbacks

	newlibstore := libstore{}
	
	newlibstore.masterServerHostPort = masterServerHostPort
	newlibstore.myHostPort = myHostPort
	newlibstore.mode = mode
	newlibstore.countMtx = &sync.Mutex{}
	newlibstore.valueMtx = &sync.Mutex{}
	newlibstore.listMtx = &sync.Mutex{}
	newlibstore.countMap = make(map[string]int)
	newlibstore.valueMap = make(map[string]CacheLease)
	
	a = &newlibstore
	b = &newlibstore

	if mode != Never {
		err := rpc.RegisterName("LeaseCallbacks", librpc.Wrap(b))
		if err != nil {
			fmt.Println("Oops! Couldn't register lease call backs")
			return nil, err
		}
	}

	srvr, err := rpc.DialHTTP("tcp", masterServerHostPort)
	if err != nil {
		fmt.Println("Oops! Returning because couldn't dial master host port")
		return nil, errors.New("Couldn't Dial Master Host Port")
	}

	//newlibstore.storageserver = srvr

	args := storagerpc.GetServersArgs{}

	var reply storagerpc.GetServersReply

	for {
		err = srvr.Call("StorageServer.GetServers", args, &reply)
		if err != nil {
			return nil, err
		}
		if reply.Status != storagerpc.OK {
			time.Sleep(1)
		} else {
			break
		}
	}

	sort.Sort(ServerSlice(reply.Servers))

	for _, node := range reply.Servers {
		srvr, _ := rpc.DialHTTP("tcp", node.HostPort)
		newlibstore.listDialers = append(newlibstore.listDialers, srvr)
	}
	//fmt.Println("Status of GetServers ", reply.Status, ", and the list of servers is ", reply.Servers)
	
	newlibstore.listServers = reply.Servers
	go newlibstore.CountChecker()
	go newlibstore.LeaseChecker()
	return a, nil
}

func getServer(ls *libstore, key string) *rpc.Client {
	hash := StoreHash(key)

	i := 0
	for _, node := range ls.listServers {
		if hash <= node.NodeID {
			return ls.listDialers[i]
		}
		i = i + 1
	}
	return ls.listDialers[0]
}

func (ls *libstore) Get(key string) (string, error) {
	//fmt.Println("Get invoked with key ", key)
	val, ok := ls.GetCacheValue(key)
	if ok {
		return val.Value.(string), nil
	} 
	
	args := storagerpc.GetArgs{}
	args.Key = key
	args.HostPort = ls.myHostPort
	if ls.mode == Always {
		args.WantLease = true	
	} else if ls.mode == Never {
		args.WantLease = false
	} else {
		args.WantLease = ls.IncreCount(key)
	}

	var reply storagerpc.GetReply
	err := getServer(ls, args.Key).Call("StorageServer.Get", &args, &reply)
	if reply.Status != storagerpc.OK {
		return "", errors.New("RPC GET Didn't return OK")
	}
	if err != nil {
		fmt.Println("RPC GET Failed")
		return "", errors.New("RPC GET Failed")
	}
	if reply.Lease.Granted {
		lease := CacheLease{ Time: time.Now(), Value: reply.Value, ValidSeconds: reply.Lease.ValidSeconds}
		ls.PutCacheValue(key, lease)
	}
	//fmt.Println("Returning ", reply.Value)
	return reply.Value, nil		
}

func (ls *libstore) Put(key, value string) error {
	//fmt.Println("PUT Invoked")
	args := storagerpc.PutArgs{Key: key, Value: value}
	var reply storagerpc.PutReply

	err := getServer(ls, args.Key).Call("StorageServer.Put", &args, &reply)
	if err != nil {
		fmt.Println("RPC Get Failed")
		return errors.New("RPC Get Failed")
	}
	if reply.Status != storagerpc.OK {
		fmt.Println("RPC Put Didn't return OK. It returned ", reply.Status)
		return errors.New(fmt.Sprintf("RPC Put Didn't return OK. It returned %d", reply.Status))
	}
	return nil
}

func (ls *libstore) Delete(key string) error {
	//fmt.Println("DELETE Invoked")

	args := storagerpc.DeleteArgs{key}

	var reply storagerpc.DeleteReply

	err := getServer(ls, args.Key).Call("StorageServer.Delete", &args, &reply)

	if err != nil {
		fmt.Println("RPC Delete Failed")
		return errors.New("RPC Delete Failed")
	}

	if reply.Status != storagerpc.OK {
		fmt.Println("RPC Delete Didn't return OK")
		return errors.New("RPC Delete Didn't return OK")
	}

	return nil
}

func (ls *libstore) GetList(key string) ([]string, error) {
	//fmt.Println("GETLIST Invoked")
	val, ok := ls.GetCacheValue(key)
	if ok {
		return val.Value.([]string), nil
	} 

	args := storagerpc.GetArgs{}
	args.Key = key
	args.HostPort = ls.myHostPort
	if ls.mode == Always {
		args.WantLease = true	
	} else if ls.mode == Never {
		args.WantLease = false
	} else {
		args.WantLease = ls.IncreCount(key)
	}
	
	var reply storagerpc.GetListReply
	err := getServer(ls, args.Key).Call("StorageServer.GetList", &args, &reply)

	if err != nil {
		fmt.Println("RPC GetList Failed")
		return nil, errors.New("RPC GetList Failed")
	}

	// TODO: Now, intelligently return error messages such that the tribserver knows what exactly
	// happened. For example, if it returns KeyNotFound, return blank list
	if reply.Status == storagerpc.KeyNotFound {
		var list []string
		return list, errors.New("RPC GetList returns storagerpc.KeyNotFound")
	}
	if reply.Status != storagerpc.OK {
		fmt.Println("RPC GetList Didn't return OK")
		return nil, errors.New(fmt.Sprintf("RPC GetList Didn't return OK. It returned %d", reply.Status))
	}
	if reply.Lease.Granted {
		lease := CacheLease{ Time: time.Now(), Value: reply.Value, ValidSeconds: reply.Lease.ValidSeconds}
		ls.PutCacheValue(key, lease)
	}
	return reply.Value, nil
}

func (ls *libstore) RemoveFromList(key, removeItem string) error {
	//fmt.Println("REMOVEFROMLIST Invoked")
	//fmt.Println("Key is ", key, " and value is ", removeItem)

	args := storagerpc.PutArgs{key, removeItem}
	var reply storagerpc.PutReply

	err := getServer(ls, args.Key).Call("StorageServer.RemoveFromList", &args, &reply)

	if err != nil {
		fmt.Println("RPC RemoveFromList Failed")
		return errors.New("RPC RemoveFromList Failed")
	}

	if reply.Status != storagerpc.OK {
		//fmt.Println("RPC RemoveFromList Didn't return OK")
		return errors.New("RPC RemoveFromList Didn't return OK")
	}

	return nil
}

func (ls *libstore) AppendToList(key, newItem string) error {
	//fmt.Println("APPENDTOLIST Invoked")
	//fmt.Println("Key is ", key, " and value is ", newItem)

	args := storagerpc.PutArgs{key, newItem}
	var reply storagerpc.PutReply

	err := getServer(ls, args.Key).Call("StorageServer.AppendToList", &args, &reply)

	if err != nil {
		fmt.Println("RPC AppendToList Failed")
		return errors.New("RPC AppendToList Failed")
	}

	if reply.Status != storagerpc.OK {
		//fmt.Println("RPC ApppendToList Didn't return OK")
		return errors.New("RPC AppendToList Didn't return OK")
	}
	
	return nil
}

func (ls *libstore) RevokeLease(args *storagerpc.RevokeLeaseArgs, reply *storagerpc.RevokeLeaseReply) error {
	fmt.Println("REVOKELEASE Invoked")
	key := args.Key 
	ls.ClearCacheValue(key)
	reply.Status = storagerpc.OK
	return nil
}

type CacheLease struct {
	Value interface{}
	Time time.Time
	ValidSeconds int
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

func (ls *libstore) LeaseChecker() {
	for {
		time.Sleep(time.Second * storagerpc.LeaseSeconds)
		ls.valueMtx.Lock()
		for key, lease := range ls.valueMap{
			if int(time.Since(lease.Time).Seconds()) > lease.ValidSeconds {
				delete(ls.valueMap, key)
			}
		}
		ls.valueMtx.Unlock()
	}
}


func (ls *libstore) CountChecker() {
	for {
		time.Sleep(time.Second * storagerpc.QueryCacheSeconds)
		ls.countMtx.Lock()
		ls.countMap = make(map[string]int)
		ls.countMtx.Unlock()
	}
}

func (ls *libstore) GetCacheValue(key string) (CacheLease, bool) {
	ls.valueMtx.Lock()
	defer ls.valueMtx.Unlock()
	
	lease, ok := ls.valueMap[key]
	if !ok {
		return CacheLease{}, false
	}
	//just lazily check if any lease is out of date
	duration := int(time.Since(lease.Time).Seconds())
	fmt.Println(duration, "seconds has passed for key", key)
	if duration >= lease.ValidSeconds{
		fmt.Println("revoking", key, "since it has passed", storagerpc.LeaseSeconds, "seconds")
		delete(ls.valueMap, key)
		return CacheLease{}, false	
	} else {
		fmt.Println("key", key, "not expired yet.")
		return lease, true
	}
}

func (ls *libstore) PutCacheValue(key string, value CacheLease) {
	ls.valueMtx.Lock()
	defer ls.valueMtx.Unlock()
	ls.valueMap[key] = value
}

func (ls *libstore) IncreCount(key string) bool {
	ls.countMtx.Lock()
	defer ls.countMtx.Unlock()
	count := ls.countMap[key]
	count += 1
	if count >= storagerpc.QueryCacheThresh {
		return true
	}
	ls.countMap[key] = count
	return false
}

func (ls *libstore) ClearCacheValue(key string) {
	ls.countMtx.Lock()
	defer ls.countMtx.Unlock()
	ls.valueMtx.Lock()
	defer ls.valueMtx.Unlock()
	ls.countMap[key] = 0	
	delete(ls.valueMap, key)
}

/*package storagerpc

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
}*/
