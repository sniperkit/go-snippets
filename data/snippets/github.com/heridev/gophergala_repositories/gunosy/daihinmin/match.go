package daihinmin

import (
	"log"
	"sync"
	"time"
)

var matches = make(map[string]*match)
var matchesMutex = &sync.RWMutex{}

type match struct {
	name    string
	id      string
	size    int
	game    *Game
	users   map[sesh]*client
	players map[sesh]*Player
	readies map[sesh]bool

	join    chan joinReq
	part    chan sesh
	infoplz chan *client
	play    chan playReq
	timeout <-chan time.Time
	die     chan struct{}
}

type ExtInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewMatch(name string) *match {
	m := &match{
		name:    name,
		id:      generateID("m:"),
		size:    4,
		game:    NewGame(),
		users:   make(map[sesh]*client),
		players: make(map[sesh]*Player),
		readies: make(map[sesh]bool),

		join:    make(chan joinReq),
		part:    make(chan sesh),
		play:    make(chan playReq),
		infoplz: make(chan *client),
		timeout: make(<-chan time.Time),
		die:     make(chan struct{}),
	}
	m.register()
	return m
}

func ListMatches() []ExtInfo {
	ls := make([]ExtInfo, 0)
	for _, m := range matches {
		ls = append(ls, ExtInfo{m.id, m.name})
	}
	return ls
}

func (m *match) register() {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	if _, exists := matches[m.name]; exists {
		panic("Remaking match: " + m.name)
	}
	matches[m.id] = m

	go m.run()
}

func (m *match) unregister() {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	log.Printf("match dying: [%s] %s", m.id, m.name)
	if _, exists := matches[m.id]; !exists {
		panic("Deleting non-existent match: " + m.id)
	}
	delete(matches, m.id)
}

func (m match) String() string {
	return m.name + " (" + m.id + ")"
}

func (m *match) run() {
	log.Printf("Running match: %s", m.name)
	defer m.unregister()

	for {
		select {
		case req := <-m.join:
			// TODO reconnect voodoo
			var ok bool
			var err string
			var playerId int
			if ok = m.size > m.usercount(); ok {
				m.users[req.sesh] = req.from
				m.broadcast(UserJoinPartReply{
					X:    "user-join",
					Chan: m.id,
					User: req.from.username(),
				})
				m.broadcast(m.info())
				p := NewPlayer(req.from.username())
				m.game.Join(p)
				m.players[req.sesh] = p
				playerId = p.Number
				if m.size == m.usercount() {
					m.game.Start()
					m.broadcast(GameInfo{
						X:     "game-started",
						ID:    m.id,
						Name:  m.name,
						Users: m.usernames(),
					})
					m.notifyNextTurn()
				}
			} else {
				err = "Exceed the limit size."
			}
			if req.result != nil {
				req.result <- joinResult{ok: ok, err: err, playerId: playerId}
			}
		case s := <-m.part:
			m.goodbye(s)
			// if everyone leaves, die
			if m.usercount() == 0 {
				return
			}
		case p := <-m.play:
			ok, hand, events := m.doPlay(p.sesh, p.cards)
			if p.result != nil {
				p.result <- playResult{ok: ok, hand: hand, events: events}
			}
			m.notifyNextTurn()
		case c := <-m.infoplz:
			c.send(m.info())
		case <-m.die:
			return
		}
	}
}

func (m *match) broadcast(msg interface{}) {
	for _, c := range m.users {
		c.send(msg)
	}
}

func (m *match) usercount() int {
	return len(m.users)
}

func (m *match) usernames() []string {
	var names []string
	for _, u := range m.users {
		names = append(names, u.username())
	}
	return names
}

func (m *match) find(name string) (s sesh, ok bool) {
	for _, u := range m.users {
		if u.username() == name {
			return u.session, true
		}
	}
	return
}

func (m *match) goodbye(s sesh) bool {
	c, ok := m.users[s]
	if !ok {
		return false
	}

	delete(m.users, s)
	c.match = nil // TODO fix data race?
	// TODO readies players
	m.broadcast(UserJoinPartReply{
		X:    "user-part",
		Chan: m.id,
		User: c.username(),
	})
	m.broadcast(m.info())
	return true
}

func (m *match) doPlay(s sesh, cards Cards) (ok bool, hand Cards, events []Event) {
	p, ok := m.players[s]
	if !ok {
		return
	}
	ok, events = m.game.Play(p, cards)
	hand = p.Hand
	return
}

func (m *match) notifyNextTurn() {
	cnum := m.game.Current
	for sesh, p := range m.players {
		if p.Number == cnum {
			if c, ok := m.users[sesh]; ok {
				c.send(YourTurn{
					X:        "your-turn",
					PlayerID: p.Number,
					Sesh:     sesh,
					Hand:     p.Hand,
				})
			}
		}
	}
	stats := make([]int, m.size)
	for _, p := range m.players {
		stats[p.Number] = p.Hand.Len()
	}
	m.broadcast(GameStatus{
		X:             "current-status",
		CurrentPlayer: m.game.Current,
		Stats:         stats,
		Pile:          m.game.Pile,
	})
}

func (m *match) info() GameInfo {
	return GameInfo{
		X:     "game-info",
		ID:    m.id,
		Name:  m.name,
		Users: m.usernames(),
	}
}

func matchExists(id string) bool {
	matchesMutex.RLock()
	defer matchesMutex.RUnlock()

	_, exists := matches[id]
	return exists
}

func getMatch(id string) *match {
	matchesMutex.RLock()
	m, exists := matches[id]
	matchesMutex.RUnlock()

	if !exists {
		log.Printf("no match: %s", id)
	}
	return m
}

type joinReq struct {
	sesh
	from     *client
	password string
	result   chan joinResult
}

type joinResult struct {
	ok       bool
	err      string
	playerId int
}

type matchReq struct {
	from     sesh
	name     string
	password string
	result   chan *match
}

type playReq struct {
	sesh
	from   *client
	cards  Cards
	result chan playResult
}

type playResult struct {
	ok     bool
	err    string
	hand   Cards
	events []Event
}
