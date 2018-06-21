package busybody

import (
	// "github.com/gdamore/mangos/transport/ipc"
	// "github.com/gdamore/mangos/transport/tlstcp"

	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/mangos"
	"github.com/zerklabs/auburn/log"
	"github.com/zerklabs/busybody/protocol"
)

// tls+tcp://x.x.x.x:ZZZZ

type BusyMember struct {
	lock             sync.RWMutex
	bussock          mangos.Socket
	config           *BusyConfig
	id               string
	hostname         string
	peers            []Introduction
	terminate        bool
	handlers         []Handler
	incomingMessages chan *protocol.Message
	StopChan         chan int
	swimTicker       *time.Ticker
	swimTimeout      *time.Timer
	swimWaitGroup    sync.WaitGroup
	polling          bool
	listening        bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(config []byte) (*BusyMember, error) {
	bussock, err := newBusSocket(make(map[string]interface{}, 0))
	if err != nil {
		return nil, err
	}

	if len(config) == 0 {
		return nil, fmt.Errorf("configuration missing")
	}

	conf, err := ParseConfig(config)
	if err != nil {
		return nil, err
	}

	member := &BusyMember{
		hostname:         hostname,
		id:               crc32hash(hostname),
		bussock:          bussock,
		config:           conf,
		terminate:        false,
		listening:        false,
		swimTicker:       time.NewTicker(conf.SwimInterval),
		swimTimeout:      time.NewTimer(conf.SwimTimeout),
		peers:            make([]Introduction, 0),
		incomingMessages: make(chan *protocol.Message),
		StopChan:         make(chan int),
		handlers:         make([]Handler, 0),
		polling:          false,
	}

	for _, v := range member.config.Peers {
		if err := member.AddPeer(v); err != nil {
			return nil, err
		}
	}

	// stop the timeout ticker
	member.swimTimeout.Stop()

	return member, nil
}

func (m *BusyMember) Uri() string {
	return m.config.Uri
}

// Generates an introduction message for this node
func (m *BusyMember) Introduction() *Introduction {
	return &Introduction{
		Key: m.config.SharedKey,
		Id:  m.id,
		Uri: m.config.Uri,
	}
}

func (m *BusyMember) connectToPeers() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, peer := range m.peers {
		if !peer.connected {
			if err := m.DialBus(&peer); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *BusyMember) AddPeer(peer string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// we don't want to add empty peer addresses
	if peer == "" {
		return fmt.Errorf("peer cannot be empty")
	}

	// make sure the stored uri is normalized
	peer = strings.ToLower(peer)

	exists := false

	for _, v := range m.peers {
		if v.Uri == peer {
			exists = true
		}
	}

	if !exists {
		intro := Introduction{Key: m.config.SharedKey, Uri: peer, connected: false, state: HealthyState}

		if m.listening {
			if err := m.bussock.Dial(peer); err != nil {
				return err
			}

			// flag the connection state
			intro.connected = true
			// pause for join
			time.Sleep(time.Second)
			if m.config.LogLevel >= log.INFO {
				log.Infof("successfully connected to: %s", peer)
			}
		}

		m.peers = append(m.peers, intro)
	}

	return nil
}

func (m *BusyMember) Members() []Introduction {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.peers
}

func (m *BusyMember) DialBus(p *Introduction) error {
	if err := m.bussock.Dial(p.Uri); err != nil {
		return err
	}

	// pause for join
	time.Sleep(time.Second)
	if m.config.LogLevel >= log.INFO {
		log.Infof("added and successfully connected to: %s (%s)", p.Id, p.Uri)
	}

	p.connected = true

	return nil
}

// randomPeer selects a random peer from this nodes peer list
func (m *BusyMember) randomPeer() *Introduction {
	time.After(time.Second * 2)
	idx := rand.Intn(len(m.peers))

	return &m.peers[idx]
}

// selectPeerGroup returns k peers for this node, excluding the target
func (m *BusyMember) selectPeerGroup(target *Introduction) []*Introduction {
	var k int

	group := make([]*Introduction, 0)

	if len(m.peers) == 0 {
		return nil
	}

	for {
		time.After(time.Second * 2)
		k = rand.Intn(len(m.peers))

		if k != 0 {
			break
		}
	}

	// log.Infof("selectPeerGroup(k): %d", k)

	// if our peer list is smaller than than 3, add the other peers
	// excluding the target
	if k <= 0 {
		for i := range m.peers {
			if m.peers[i].Id != target.Id {
				group = append(group, &m.peers[i])
			}
		}
	} else {
		count := 0
		for i := range m.peers {
			if m.peers[i].Id != target.Id {
				group = append(group, &m.peers[i])
				count += 1
			}

			if count == k {
				break
			}
		}
	}

	return group
}

func (m *BusyMember) updatePeer(intro *Introduction) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	exists := false
	if intro.Id == m.id {
		return nil
	}

	for i, v := range m.peers {
		if v.Uri == intro.Uri {
			exists = true
			if v.Id == "" {

				if !m.peers[i].connected {
					if err := m.DialBus(&m.peers[i]); err != nil {
						return err
					}
				}
				m.peers[i] = *intro
				m.peers[i].connected = true
				m.peers[i].state = HealthyState

				if m.config.LogLevel >= log.INFO {
					log.Infof("updated peer: %#v", intro)
				}
			}
		}
	}

	if !exists {
		if err := m.DialBus(intro); err != nil {
			return err
		}

		intro.state = HealthyState
		intro.connected = true

		m.peers = append(m.peers, *intro)
	}

	return nil
}

func (m *BusyMember) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.swimTicker.Stop()
	m.terminate = true
	close(m.StopChan)
	close(m.incomingMessages)

	return nil
}

func (m *BusyMember) notificationLoop() {
	for {
		select {
		case <-m.swimTimeout.C:
			if m.polling {
				m.polling = false
			}
		case <-m.swimTicker.C:
			m.swimTimeout.Reset(m.config.SwimTimeout)

			for {
				if err := m.hello(); err != nil {
					log.Error(err)
				}

				if err := m.share(); err != nil {
					log.Error(err)
				}

				break
			}
		}
	}
}

func (m *BusyMember) handlerLoop() {
	for {
		message, ok := <-m.incomingMessages
		if !ok {
			goto exit
		}

		if message.MessageType() == protocol.StandardMessage {
			for _, handler := range m.handlers {
				if err := handler.HandleMessage(message); err != nil {
					log.Errorf("error during HandleMessage: %v", err)
				}
			}
		}

		if message.MessageType() == protocol.HelloMessage {
			intro, err := UnmarshalIntroduction(message)
			if err != nil {
				log.Error(err)
				continue
			}

			if intro.Key != m.config.SharedKey {
				if m.config.LogLevel >= log.WARN {
					log.Warnf("received unauthorized introduction: %#v", intro)
				}
				continue
			}

			if err := m.updatePeer(intro); err != nil {
				log.Error(err)
			}
		}
	}

exit:
	log.Infof("stopping handler")
	m.Close()
}

func (m *BusyMember) AddHandler(handler Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.handlers = append(m.handlers, handler)
}

func (m *BusyMember) Listen() error {
	var msg []byte
	var err error

	defer m.Close()

	// start our listener
	if err := m.bussock.Listen(m.config.Uri); err != nil {
		return err
	}

	m.listening = true

	if err := m.connectToPeers(); err != nil {
		return err
	}

	time.After(time.Second * 2)

	// start dealing with incoming messages
	go m.handlerLoop()
	go m.notificationLoop()

	for {
		if m.terminate {
			return nil
		}

		if msg, err = m.bussock.Recv(); err != nil {
			return err
		}

		bmsg, err := protocol.Decode(msg)
		if err != nil {
			log.Errorf("error decoding message: %v", err)
			continue
		}

		// we ignore messages from ourselves
		if bmsg.Sender() != m.id {
			m.incomingMessages <- bmsg
		}
	}

	return nil
}
