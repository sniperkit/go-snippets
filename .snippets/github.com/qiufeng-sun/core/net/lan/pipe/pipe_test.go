package pipe

import (
	"sync"
	"testing"
	"time"

	. "core/net/lan"
)

var _ = time.Now

//
func TestLan(t *testing.T) {
	url := "tcp://127.0.0.1:8801"
	loop := 10
	wg := &sync.WaitGroup{}
	t.Log("loop:", loop)

	wg.Add(1)
	go func() {
		t.Log("create server!")
		s := NewServer(NewLanCfg("test", url))

		for i := 0; i < loop; i++ {
			msg, e := s.Recv()
			if e != nil {
				t.Fatal("server recv", i, e)
			}
			t.Log("server recv", i, string(msg))
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		c := NewClient(NewLanCfg("test", url))

		t.Log("create client!")
		for i := 0; i < loop; i++ {
			t.Log("send:", i)
			e := c.Send([]byte("test"))
			if e != nil {
				t.Fatal("client send", i, e)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

//
func TestClients_Update(t *testing.T) {
	ds := [][]string{
		{"data", "127.0.0.1:8801"},
		{"match", "127.0.0.1:8811", "127.0.0.1:8812", "127.0.0.1:8813"},
		{"data", "127.0.0.1:8801", "127.0.0.1:8802"},
		{"match"},
		{"data"},
		{"match", "127.0.0.1:8811", "127.0.0.1:8812", "127.0.0.1:8813"},
	}

	cs := NewClients("gw")

	for _, v := range ds {
		cs.Update(v[0], v[1:])
	}
}
