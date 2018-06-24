package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"regexp"
	"time"

	"github.com/cmu440/tribbler/rpc/storagerpc"
)

type storageTester struct {
	srv        *rpc.Client
	myhostport string
}

type testFunc struct {
	name string
	f    func()
}

var (
	portnum   = flag.Int("port", 9019, "port # to listen on")
	testType  = flag.Int("type", 1, "type of test, 1: jtest, 2: btest")
	numServer = flag.Int("N", 1, "(jtest only) total # of storage servers")
	myID      = flag.Int("id", 1, "(jtest only) my id")
	testRegex = flag.String("t", "", "test to run")
	passCount int
	failCount int
	st        *storageTester
)

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

var statusMap = map[storagerpc.Status]string{
	storagerpc.OK:           "OK",
	storagerpc.KeyNotFound:  "KeyNotFound",
	storagerpc.ItemNotFound: "ItemNotFound",
	storagerpc.WrongServer:  "WrongServer",
	storagerpc.ItemExists:   "ItemExists",
	storagerpc.NotReady:     "NotReady",
	0:                       "Unknown",
}

func initStorageTester(server, myhostport string) (*storageTester, error) {
	tester := new(storageTester)
	tester.myhostport = myhostport

	// Create RPC connection to storage server.
	srv, err := rpc.DialHTTP("tcp", server)
	if err != nil {
		return nil, fmt.Errorf("could not connect to server %s", server)
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if err != nil {
		LOGE.Fatalln("Failed to listen:", err)
	}
	go http.Serve(l, nil)
	tester.srv = srv
	return tester, nil
}

func (st *storageTester) RegisterServer() (*storagerpc.RegisterReply, error) {
	node := storagerpc.Node{HostPort: st.myhostport, NodeID: uint32(*myID)}
	args := &storagerpc.RegisterArgs{ServerInfo: node}
	var reply storagerpc.RegisterReply
	err := st.srv.Call("StorageServer.RegisterServer", args, &reply)
	return &reply, err
}

func (st *storageTester) GetServers() (*storagerpc.GetServersReply, error) {
	args := &storagerpc.GetServersArgs{}
	var reply storagerpc.GetServersReply
	err := st.srv.Call("StorageServer.GetServers", args, &reply)
	return &reply, err
}

func (st *storageTester) Put(key, value string) (*storagerpc.PutReply, error) {
	args := &storagerpc.PutArgs{Key: key, Value: value}
	var reply storagerpc.PutReply
	err := st.srv.Call("StorageServer.Put", args, &reply)
	return &reply, err
}

func (st *storageTester) Delete(key string) (*storagerpc.DeleteReply, error) {
	args := &storagerpc.DeleteArgs{Key: key}
	var reply storagerpc.DeleteReply
	err := st.srv.Call("StorageServer.Delete", args, &reply)
	return &reply, err
}

func (st *storageTester) Get(key string, wantlease bool) (*storagerpc.GetReply, error) {
	args := &storagerpc.GetArgs{Key: key, WantLease: wantlease, HostPort: st.myhostport}
	var reply storagerpc.GetReply
	err := st.srv.Call("StorageServer.Get", args, &reply)
	return &reply, err
}

func (st *storageTester) GetList(key string, wantlease bool) (*storagerpc.GetListReply, error) {
	args := &storagerpc.GetArgs{Key: key, WantLease: wantlease, HostPort: st.myhostport}
	var reply storagerpc.GetListReply
	err := st.srv.Call("StorageServer.GetList", args, &reply)
	return &reply, err
}

func (st *storageTester) RemoveFromList(key, removeitem string) (*storagerpc.PutReply, error) {
	args := &storagerpc.PutArgs{Key: key, Value: removeitem}
	var reply storagerpc.PutReply
	err := st.srv.Call("StorageServer.RemoveFromList", args, &reply)
	return &reply, err
}

func (st *storageTester) AppendToList(key, newitem string) (*storagerpc.PutReply, error) {
	args := &storagerpc.PutArgs{Key: key, Value: newitem}
	var reply storagerpc.PutReply
	err := st.srv.Call("StorageServer.AppendToList", args, &reply)
	return &reply, err
}

// Check error and status
func checkErrorStatus(err error, status, expectedStatus storagerpc.Status) bool {
	if err != nil {
		LOGE.Println("FAIL: unexpected error returned:", err)
		failCount++
		return true
	}
	if status != expectedStatus {
		LOGE.Printf("FAIL: incorrect status %s, expected status %s\n", statusMap[status], statusMap[expectedStatus])
		failCount++
		return true
	}
	return false
}

// Check error
func checkError(err error, expectError bool) bool {
	if expectError {
		if err == nil {
			LOGE.Println("FAIL: non-nil error should be returned")
			failCount++
			return true
		}
	} else {
		if err != nil {
			LOGE.Println("FAIL: unexpected error returned:", err)
			failCount++
			return true
		}
	}
	return false
}

// Check list
func checkList(list []string, expectedList []string) bool {
	if len(list) != len(expectedList) {
		LOGE.Printf("FAIL: incorrect list %v, expected list %v\n", list, expectedList)
		failCount++
		return true
	}
	m := make(map[string]bool)
	for _, s := range list {
		m[s] = true
	}
	for _, s := range expectedList {
		if m[s] == false {
			LOGE.Printf("FAIL: incorrect list %v, expected list %v\n", list, expectedList)
			failCount++
			return true
		}
	}
	return false
}

// We treat a RPC call finihsed in 0.5 seconds as OK
func isTimeOK(d time.Duration) bool {
	return d < 500*time.Millisecond
}

// Cache a key
func cacheKey(key string) bool {
	replyP, err := st.Put(key, "old-value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return true
	}

	// get and cache key
	replyG, err := st.Get(key, true)
	if checkErrorStatus(err, replyG.Status, storagerpc.OK) {
		return true
	}
	if !replyG.Lease.Granted {
		LOGE.Println("FAIL: Failed to get lease")
		failCount++
		return true
	}
	return false
}

// Cache a list key
func cacheKeyList(key string) bool {
	replyP, err := st.AppendToList(key, "old-value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return true
	}

	// get and cache key
	replyL, err := st.GetList(key, true)
	if checkErrorStatus(err, replyL.Status, storagerpc.OK) {
		return true
	}
	if !replyL.Lease.Granted {
		LOGE.Println("FAIL: Failed to get lease")
		failCount++
		return true
	}
	return false
}

/////////////////////////////////////////////
//  test storage server initialization
/////////////////////////////////////////////

// make sure to run N-1 servers in shell before entering this function
func testInitStorageServers() {
	// test get server
	replyGS, err := st.GetServers()
	if checkError(err, false) {
		return
	}
	if replyGS.Status == storagerpc.OK {
		LOGE.Println("FAIL: storage system should not be ready:", err)
		failCount++
		return
	}

	// test register
	replyR, err := st.RegisterServer()
	if checkError(err, false) {
		return
	}
	if replyR.Status != storagerpc.OK || replyR.Servers == nil {
		LOGE.Println("FAIL: storage system should be ready and Servers field should be non-nil:", err)
		failCount++
		return
	}
	if len(replyR.Servers) != (*numServer) {
		LOGE.Println("FAIL: storage system returned wrong server list:", err)
		failCount++
		return
	}

	fmt.Println("PASS")
	passCount++
}

/////////////////////////////////////////////
//  test basic storage operations
/////////////////////////////////////////////

// Get keys without wantlease
func testPutGetWithoutLease() {
	// get an invalid key
	replyG, err := st.Get("nullkey:1", false)
	if checkErrorStatus(err, replyG.Status, storagerpc.KeyNotFound) {
		return
	}

	replyP, err := st.Put("keyputget:1", "value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// without asking for a lease
	replyG, err = st.Get("keyputget:1", false)
	if checkErrorStatus(err, replyG.Status, storagerpc.OK) {
		return
	}
	if replyG.Value != "value" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	if replyG.Lease.Granted {
		LOGE.Println("FAIL: did not apply for lease")
		failCount++
		return
	}

	fmt.Println("PASS")
	passCount++
}

// Without leasing, we should not expect revoke
func testDeleteWithoutLease() {
	key := "revokekey:0"

	replyP, err := st.Put(key, "value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// delete this key
	replyD, err := st.Delete(key)
	if checkErrorStatus(err, replyD.Status, storagerpc.OK) {
		return
	}

	// read it back
	replyG, err := st.Get(key, false)
	if checkErrorStatus(err, replyG.Status, storagerpc.KeyNotFound) {
		return
	}

	fmt.Println("PASS")
	passCount++
}

// list related operations
func testAppendGetRemoveList() {
	// test AppendToList
	replyP, err := st.AppendToList("keylist:1", "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// test GetList
	replyL, err := st.GetList("keylist:1", false)
	if checkErrorStatus(err, replyL.Status, storagerpc.OK) {
		return
	}
	if len(replyL.Value) != 1 || replyL.Value[0] != "value1" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}

	// test AppendToList for a duplicated item
	replyP, err = st.AppendToList("keylist:1", "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.ItemExists) {
		return
	}

	// test AppendToList for a different item
	replyP, err = st.AppendToList("keylist:1", "value2")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// test RemoveFromList for the first item
	replyP, err = st.RemoveFromList("keylist:1", "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// test RemoveFromList for removed item
	replyP, err = st.RemoveFromList("keylist:1", "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.ItemNotFound) {
		return
	}

	// test GetList after RemoveFromList
	replyL, err = st.GetList("keylist:1", false)
	if checkErrorStatus(err, replyL.Status, storagerpc.OK) {
		return
	}
	if len(replyL.Value) != 1 || replyL.Value[0] != "value2" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}

	fmt.Println("PASS")
	passCount++
}

// Without leasing, we should not expect revoke
func testUpdateWithoutLease() {
	key := "revokekey:0"

	replyP, err := st.Put(key, "value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// get without caching this item
	replyG, err := st.Get(key, false)
	if checkErrorStatus(err, replyG.Status, storagerpc.OK) {
		return
	}

	// update this key
	replyP, err = st.Put(key, "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// get without caching this item
	replyG, err = st.Get(key, false)
	if checkErrorStatus(err, replyG.Status, storagerpc.OK) {
		return
	}

	fmt.Println("PASS")
	passCount++
}

// Without leasing, we should not expect revoke
func testUpdateListWithoutLease() {
	key := "revokelistkey:0"

	replyP, err := st.AppendToList(key, "value")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// get without caching this item
	replyL, err := st.GetList(key, false)
	if checkErrorStatus(err, replyL.Status, storagerpc.OK) {
		return
	}

	// update this key
	replyP, err = st.AppendToList(key, "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}
	replyP, err = st.RemoveFromList(key, "value1")
	if checkErrorStatus(err, replyP.Status, storagerpc.OK) {
		return
	}

	// get without caching this item
	replyL, err = st.GetList(key, false)
	if checkErrorStatus(err, replyL.Status, storagerpc.OK) {
		return
	}

	fmt.Println("PASS")
	passCount++
}

func main() {
	jtests := []testFunc{{"testInitStorageServers", testInitStorageServers}}
	btests := []testFunc{
		{"testPutGetWithoutLease", testPutGetWithoutLease},
		{"testDeleteWithoutLease", testDeleteWithoutLease},
		{"testAppendGetRemoveList", testAppendGetRemoveList},
		{"testUpdateWithoutLease", testUpdateWithoutLease},
		{"testUpdateListWithoutLease", testUpdateListWithoutLease},
	}

	flag.Parse()
	if flag.NArg() < 1 {
		LOGE.Fatalln("Usage: storagetest <storage master>")
	}

	// Run the tests with a single tester
	storageTester, err := initStorageTester(flag.Arg(0), fmt.Sprintf("localhost:%d", *portnum))
	if err != nil {
		LOGE.Fatalln("Failed to initialize test:", err)
	}
	st = storageTester

	switch *testType {
	case 1:
		for _, t := range jtests {
			if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
				fmt.Printf("Running %s:\n", t.name)
				t.f()
			}
		}
	case 2:
		for _, t := range btests {
			if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
				fmt.Printf("Running %s:\n", t.name)
				t.f()
			}
		}
	}

	fmt.Printf("Passed (%d/%d) tests\n", passCount, passCount+failCount)
}
