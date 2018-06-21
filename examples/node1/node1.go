package main

import (
	"os"
	"time"

	"github.com/zerklabs/auburn/log"
	"github.com/zerklabs/busybody"
	"github.com/zerklabs/busybody/protocol"
)

var testConfig = `
uri = "ipc:///tmp/ipc0.ipc"
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
	member1, err := busybody.New([]byte(testConfig))
	if err != nil {
		log.Error(err)
		return
	}

	member1.AddHandler(busybody.HandlerFunc(func(m *protocol.Message) error {
		body, err := m.Body()
		if err != nil {
			log.Error(err)

			return err
		}

		log.Infof("member1: message received: %s", string(body))

		return nil
	}))

	go func() {
		os.Remove(member1.Uri())

		if err := member1.Listen(); err != nil {
			log.Error(err)
		}
	}()

	time.After(time.Second * 10)

	member1.AddPeer("ipc:///tmp/ipc1.ipc")

	for i := 0; i < 5; i++ {
		if err := member1.Send([]byte("hello from member1")); err != nil {
			log.Error(err)
		}
	}

	member1.Close()
}
