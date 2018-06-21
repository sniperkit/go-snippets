package rpc

import (
	"testing"

	. "core/net/lan"
	"sync"
	"fmt"
)

//
var _ = testing.Coverage

//
var g_srvAddr = "localhost:8090"
var g_srvAddrs = []string{g_srvAddr}
var g_rpcSrv = NewServer(NewLanCfg("rpc server", g_srvAddr))
var g_clientPool = NewClientPool(&PoolConfig{"rpc client", 2, 5, 10})

//
func init() {
	g_rpcSrv.Init(4, handleMsgs)
	g_rpcSrv.Serve()

	g_clientPool.Update("server", g_srvAddrs)
}

//
func handleMsgs(msg []byte) []byte {
	ret := append([]byte("hello "), msg...)
	return ret
}

//
func getClient() *Client {
	return g_clientPool.GetClient()
}

func returnClient(c *Client) {
	g_clientPool.ReturnClient(c)
}

//
func TestClient(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i:=0; i<10; i++ {
		wg.Add(1)
		go testOneClient(t, wg, i)
	}
	wg.Wait()
}

//
func testOneClient(t *testing.T, wg *sync.WaitGroup, index int) {
	defer wg.Done()
	client := getClient()
	if nil == client {
		t.Fatal("client is nil")
	}
	defer returnClient(client)

	msg := fmt.Sprintf("test %v", index)

	r, e := client.Call([]byte(msg))
	t.Log(string(r), e)

	if e != nil {
		t.Fatal(e)
	}
}

//
func TestClientPool_Update(t *testing.T) {
	ds := [][]string{
		{"test", "127.0.0.1:8801"},
		{"test", "127.0.0.1:8811", "127.0.0.1:8812", "127.0.0.1:8813"},
		{"test", "127.0.0.1:8801", "127.0.0.1:8802"},
		{"test"},
		{"test"},
		{"test", "127.0.0.1:8811", "127.0.0.1:8812", "127.0.0.1:8813"},
	}

	for _, v := range ds {
		g_clientPool.Update(v[0], v[1:])
	}
}
