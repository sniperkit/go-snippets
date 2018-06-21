package main

import (
	"os"
	"sync"

	"github.com/zerklabs/auburn/log"
	"github.com/zerklabs/busybody"
	"github.com/zerklabs/busybody/protocol"
)

var testConfig2 = `
uri = "ipc:///tmp/ipc1.ipc"
peers = [ ]
shared_key = "default_shared_key"
swim_interval = "30s"
swim_timeout = "10s"
snappy_compression = false
zlib_compression = false
deflate_compression = true
deflate_compression_level = 6
log_level = 7
`

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	member2, err := busybody.New([]byte(testConfig2))
	if err != nil {
		log.Error(err)
	}

	member2.AddHandler(busybody.HandlerFunc(func(m *protocol.Message) error {
		body, err := m.Body()
		if err != nil {
			log.Error(err)

			return err
		}

		log.Infof("member2: message received: %s", string(body))
		wg.Done()

		return nil
	}))

	go func() {
		os.Remove(member2.Uri())
		if err := member2.Listen(); err != nil {
			log.Error(err)
		}
	}()

	member2.AddPeer("ipc:///tmp/ipc0.ipc")

	wg.Wait()
}
