package busybody

import (
	"sync"
	"testing"
	"time"

	"github.com/zerklabs/busybody/protocol"
)

var testConfig = `
uri = "ipc:///tmp/ipc0.ipc"
peers = [ "ipc:///tmp/ipc1.ipc", "ipc:///tmp/ipc2.ipc", "ipc:///tmp/ipc3.ipc", "ipc:///tmp/ipc4.ipc" ]
shared_key = "default_shared_key"
swim_interval = "1m0s"
swim_timeout = "30s"
snappy_compression = false
zlib_compression = false
deflate_compression = true
deflate_compression_level = 6
log_level = 6
`

var testConfig2 = `
uri = "ipc:///tmp/ipc1.ipc"
peers = [ "ipc:///tmp/ipc0.ipc", "ipc:///tmp/ipc2.ipc", "ipc:///tmp/ipc3.ipc", "ipc:///tmp/ipc4.ipc" ]
shared_key = "default_shared_key"
swim_interval = "1m0s"
swim_timeout = "30s"
snappy_compression = false
zlib_compression = false
deflate_compression = true
deflate_compression_level = 6
log_level = 6
`

func TestNew(t *testing.T) {
	member, err := New([]byte(testConfig))
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", member)
}

func TestRandomPeer(t *testing.T) {
	member, err := New([]byte(testConfig))
	if err != nil {
		t.Error(err)
	}

	peer := member.randomPeer()
	if peer == nil {
		t.Errorf("expected a peer to be returned")
	}

	t.Logf("%#v", peer)
}

// func TestRandomPeerGroup(t *testing.T) {
// 	member, err := New([]byte(testConfig))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	group := member.selectPeerGroup(&member.peers[0])
// 	if group == nil {
// 		t.Errorf("expected a non-nil slice to be returned")
// 	}

// 	if len(group) <= 0 {
// 		t.Errorf("expected a non-empty slice to be returned")
// 	}

// 	for i := range group {
// 		t.Logf("%#v", group[i])
// 	}
// }

func TestSend(t *testing.T) {
	wg := new(sync.WaitGroup)

	member1, err := New([]byte(testConfig))
	if err != nil {
		t.Error(err)
	}

	member2, err := New([]byte(testConfig2))
	if err != nil {
		t.Error(err)
	}

	member1.AddHandler(HandlerFunc(func(m *protocol.Message) error {
		body, _ := m.Body()
		t.Logf("member1: message received: %s", string(body))

		return nil
	}))

	member2.AddHandler(HandlerFunc(func(m *protocol.Message) error {
		body, _ := m.Body()
		t.Logf("member2: message received: %s", string(body))

		wg.Done()

		return nil
	}))

	go func(m *BusyMember) {
		if err := m.Listen(); err != nil {
			t.Errorf("error listening: %v", err)
		}
	}(member2)

	go func(m *BusyMember) {
		if err := m.Listen(); err != nil {
			t.Errorf("error listening: %v", err)
		}
	}(member1)

	time.After(time.Second * 5)

	member1.Send([]byte("hello"))

	wg.Wait()
}
