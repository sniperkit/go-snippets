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

	"github.com/cmu440/tribbler/libstore"
	"github.com/cmu440/tribbler/rpc/storagerpc"
	"github.com/cmu440/tribbler/tests/proxycounter"
)

type testFunc struct {
	name string
	f    func()
}

var (
	portnum   = flag.Int("port", 9010, "port to listen on")
	testRegex = flag.String("t", "", "test to run")
)

var (
	pc         proxycounter.ProxyCounter
	ls         libstore.Libstore
	revokeConn *rpc.Client
	passCount  int
	failCount  int
)

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

// Initialize proxy and libstore
func initLibstore(storage, server, myhostport string, alwaysLease bool) (net.Listener, error) {
	l, err := net.Listen("tcp", server)
	if err != nil {
		LOGE.Println("Failed to listen:", err)
		return nil, err
	}

	// The ProxyServer acts like a "StorageServer" in the system, but also has some
	// additional functionalities that allow us to enforce the number of RPCs made
	// to the storage server, etc.
	proxyCounter, err := proxycounter.NewProxyCounter(storage, server)
	if err != nil {
		LOGE.Println("Failed to setup test:", err)
		return nil, err
	}
	pc = proxyCounter

	// Normally the StorageServer would register itself to receive RPCs,
	// but we don't call NewStorageServer here, do we need to do it here instead.
	rpc.RegisterName("StorageServer", storagerpc.Wrap(pc))

	// Normally the TribServer would call the two methods below when it is first
	// created, but these tests mock out the TribServer all together, so we do
	// it here instead.
	rpc.HandleHTTP()
	go http.Serve(l, nil)

	var leaseMode libstore.LeaseMode
	if alwaysLease {
		leaseMode = libstore.Always
	} else if myhostport == "" {
		leaseMode = libstore.Never
	} else {
		leaseMode = libstore.Normal
	}

	// Create and start the Libstore.
	libstore, err := libstore.NewLibstore(server, myhostport, leaseMode)
	if err != nil {
		LOGE.Println("Failed to create Libstore:", err)
		return nil, err
	}
	ls = libstore
	return l, nil
}

// Cleanup libstore and rpc hooks
func cleanupLibstore(l net.Listener) {
	// Close listener to stop http serve thread
	if l != nil {
		l.Close()
	}
	// Recreate default http serve mux
	http.DefaultServeMux = http.NewServeMux()
	// Recreate default rpc server
	rpc.DefaultServer = rpc.NewServer()
	// Unset libstore just in case
	ls = nil
}

// Check rpc and byte count limits
func checkLimits(rpcCountLimit, byteCountLimit uint32) bool {
	if pc.GetRpcCount() > rpcCountLimit {
		LOGE.Println("FAIL: using too many RPCs")
		failCount++
		return true
	}
	if pc.GetByteCount() > byteCountLimit {
		LOGE.Println("FAIL: transferring too much data")
		failCount++
		return true
	}
	return false
}

// Check error
func checkError(err error, expectError bool) bool {
	if expectError {
		if err == nil {
			LOGE.Println("FAIL: error should be returned")
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

// Test libstore returns nil when it cannot connect to the server
func testNonexistentServer() {
	if l, err := libstore.NewLibstore(fmt.Sprintf("localhost:%d", *portnum), fmt.Sprintf("localhost:%d", *portnum), libstore.Normal); l == nil || err != nil {
		fmt.Println("PASS")
		passCount++
	} else {
		LOGE.Println("FAIL: libstore does not return a non-nil error when it cannot connect to nonexistent storage server")
		failCount++
	}
	cleanupLibstore(nil)
}

// Handle get error
func testGetError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	_, err := ls.Get("key:1")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle get error reply status
func testGetErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.KeyNotFound)
	defer pc.OverrideOff()
	_, err := ls.Get("key:2")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid get
func testGetValid() {
	ls.Put("key:3", "value")
	pc.Reset()
	v, err := ls.Get("key:3")
	if checkError(err, false) {
		return
	}
	if v != "value" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle put error
func testPutError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	err := ls.Put("key:4", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle put error reply status
func testPutErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.WrongServer /* use arbitrary status */)
	defer pc.OverrideOff()
	err := ls.Put("key:5", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid put
func testPutValid() {
	pc.Reset()
	err := ls.Put("key:6", "value")
	if checkError(err, false) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	v, err := ls.Get("key:6")
	if checkError(err, false) {
		return
	}
	if v != "value" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle delete error
func testDeleteError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	err := ls.Delete("key:4")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle delete error reply status
func testDeleteErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.WrongServer /* use arbitrary status */)
	defer pc.OverrideOff()
	err := ls.Delete("key:5")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid delete
func testDeleteValid() {
	pc.Reset()
	err := ls.Put("key:6", "value")
	if checkError(err, false) {
		return
	}
	if checkLimits(5, 50) {
		return
	}

	err = ls.Delete("key:6")
	if checkError(err, false) {
		return
	}

	_, err = ls.Get("key:6")
	if checkError(err, true) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle get list error
func testGetListError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	_, err := ls.GetList("keylist:1")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle get list error reply status
func testGetListErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.KeyNotFound)
	defer pc.OverrideOff()
	_, err := ls.GetList("keylist:2")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid get list
func testGetListValid() {
	ls.AppendToList("keylist:3", "value")
	pc.Reset()
	v, err := ls.GetList("keylist:3")
	if checkError(err, false) {
		return
	}
	if len(v) != 1 || v[0] != "value" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle append to list error
func testAppendToListError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	err := ls.AppendToList("keylist:4", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle append to list error reply status
func testAppendToListErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.ItemExists)
	defer pc.OverrideOff()
	err := ls.AppendToList("keylist:5", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid append to list
func testAppendToListValid() {
	pc.Reset()
	err := ls.AppendToList("keylist:6", "value")
	if checkError(err, false) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	v, err := ls.GetList("keylist:6")
	if checkError(err, false) {
		return
	}
	if len(v) != 1 || v[0] != "value" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle remove from list error
func testRemoveFromListError() {
	pc.Reset()
	pc.OverrideErr()
	defer pc.OverrideOff()
	err := ls.RemoveFromList("keylist:7", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle remove from list error reply status
func testRemoveFromListErrorStatus() {
	pc.Reset()
	pc.OverrideStatus(storagerpc.ItemNotFound)
	defer pc.OverrideOff()
	err := ls.RemoveFromList("keylist:8", "value")
	if checkError(err, true) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	fmt.Println("PASS")
	passCount++
}

// Handle valid remove from list
func testRemoveFromListValid() {
	err := ls.AppendToList("keylist:9", "value1")
	if checkError(err, false) {
		return
	}
	err = ls.AppendToList("keylist:9", "value2")
	if checkError(err, false) {
		return
	}
	pc.Reset()
	err = ls.RemoveFromList("keylist:9", "value1")
	if checkError(err, false) {
		return
	}
	if checkLimits(5, 50) {
		return
	}
	v, err := ls.GetList("keylist:9")
	if checkError(err, false) {
		return
	}
	if len(v) != 1 || v[0] != "value2" {
		LOGE.Println("FAIL: got wrong value")
		failCount++
		return
	}
	fmt.Println("PASS")
	passCount++
}

func main() {
	initTests := []testFunc{
		{"testNonexistentServer", testNonexistentServer},
	}
	tests := []testFunc{
		{"testGetError", testGetError},
		{"testGetErrorStatus", testGetErrorStatus},
		{"testGetValid", testGetValid},
		{"testPutError", testPutError},
		{"testPutErrorStatus", testPutErrorStatus},
		{"testPutValid", testPutValid},
		{"testDeleteError", testDeleteError},
		{"testDeleteErrorStatus", testDeleteErrorStatus},
		{"testDeleteValid", testDeleteValid},
		{"testGetListError", testGetListError},
		{"testGetListErrorStatus", testGetListErrorStatus},
		{"testGetListValid", testGetListValid},
		{"testAppendToListError", testAppendToListError},
		{"testAppendToListErrorStatus", testAppendToListErrorStatus},
		{"testAppendToListValid", testAppendToListValid},
		{"testRemoveFromListError", testRemoveFromListError},
		{"testRemoveFromListErrorStatus", testRemoveFromListErrorStatus},
		{"testRemoveFromListValid", testRemoveFromListValid},
	}

	flag.Parse()
	if flag.NArg() < 1 {
		LOGE.Fatalln("Usage: libtest <storage master host:port>")
	}

	var err error

	// Run init tests
	for _, t := range initTests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			fmt.Printf("Running %s:\n", t.name)
			t.f()
		}
		// Give the current Listener some time to close before creating
		// a new Libstore.
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	_, err = initLibstore(flag.Arg(0), fmt.Sprintf("localhost:%d", *portnum), fmt.Sprintf("localhost:%d", *portnum), false)
	if err != nil {
		return
	}
	revokeConn, err = rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", *portnum))
	if err != nil {
		LOGE.Println("Failed to connect to Libstore RPC:", err)
		return
	}

	// Run tests
	for _, t := range tests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			fmt.Printf("Running %s:\n", t.name)
			t.f()
		}
	}

	fmt.Printf("Passed (%d/%d) tests\n", passCount, passCount+failCount)
}
