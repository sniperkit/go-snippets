package lobby

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"git.andrewcsellers.com/acsellers/card_sharp/config"
	"git.andrewcsellers.com/acsellers/card_sharp/store"
	"github.com/gorilla/securecookie"
)

var (
	activeLobbies map[string]*Lobby
	knownPlayers  map[string]*Player
)

func init() {
	activeLobbies = map[string]*Lobby{}
	knownPlayers = map[string]*Player{}
}

type Lobby struct {
	ID         string
	Name       string
	State      string
	Presenter  chan []byte
	Players    []*Player
	PlayerChan map[string]chan []byte
	Creation   time.Time
	Czar       *Player
	Game       *store.Deck
	Instance   *store.Game
}

func Create(d store.Deck) *Lobby {
	l := Lobby{
		ID:         fmt.Sprint(rand.Intn(100000)),
		Name:       "Lobby Name Here",
		State:      "rules",
		Creation:   time.Now(),
		PlayerChan: map[string]chan []byte{},
		Players:    []*Player{},
		Game:       &d,
	}

	// make sure we aren't reusing another lobby's name
	_, ok := activeLobbies[l.ID]
	for ok {
		l.ID = fmt.Sprint(rand.Intn(100000))
		_, ok = activeLobbies[l.ID]
	}

	activeLobbies[l.ID] = &l
	return &l
}

func Find(id string) *Lobby {
	return activeLobbies[id]
}

func (l Lobby) FindPlayer(pid string) *Player {
	for _, p := range l.Players {
		if p.ID == pid {
			return p
		}
	}
	return nil
}

func (g *Lobby) Add(p *Player) {
	for _, ap := range g.Players {
		if p == ap {
			return
		}
	}
	if len(g.Players) == 0 {
		g.Czar = p
		p.Czar = true
	}
	g.Players = append(g.Players, p)
	fmt.Println("Player:", p)
	go g.NotifyFor(p)
	if len(g.Players) >= g.Game.MinPlayer {
		go g.GameReady()
	}
}

func (g *Lobby) GameReady() {
	if ws, ok := g.PlayerChan[g.Czar.ID]; ok {
		g.Czar.Status = "czar-ready"
		b := &bytes.Buffer{}
		config.Tmpl.ExecuteTemplate(b, "player_show.html", map[string]interface{}{"Lobby": g, "Player": g.Czar})

		je, _ := json.Marshal(Event{
			Type: "game_ready",
			HTML: b.String(),
		})
		ws <- je
	} else {
		// handle missing czar
	}
}

func (g *Lobby) GameStart() {
	in := g.Game.BuildGame(len(g.Players))
	g.Instance = &in
	for i, p := range g.Players {
		g.Instance.Players[p.ID] = i
		if g.Czar != p {
			p.Status = "play"
		} else {
			p.Status = "judge-wait"
		}
		if ws, ok := g.PlayerChan[p.ID]; ok {
			b := &bytes.Buffer{}
			err := config.Tmpl.ExecuteTemplate(
				b,
				"player_show.html",
				map[string]interface{}{
					"Lobby":  g,
					"Player": p,
					"Hand":   g.Instance.Hands[i],
				})
			if err == nil {
				je, _ := json.Marshal(
					Event{
						Type: "game_start",
						HTML: b.String(),
					},
				)
				ws <- je
			} else {
				fmt.Println("GameStart", err)
				je, _ := json.Marshal(
					Event{
						Type: "sync",
					},
				)
				ws <- je
			}
		}
	}
	g.State = "play"
	je, _ := json.Marshal(Event{Type: "game_start"})
	g.Presenter <- je
}

func (g *Lobby) NotifyFor(p *Player) {
	b := &bytes.Buffer{}
	config.Tmpl.ExecuteTemplate(b, "player_list.html", p)
	je, _ := json.Marshal(Event{
		Type: "new_player",
		HTML: b.String(),
		Data: *p,
	})
	g.Presenter <- je
}

type Player struct {
	ID          string
	Name        string
	Czar        bool
	Status      string
	Cards       []*store.Card
	RoundPoints int
}

func (p Player) StatusMessage() string {
	switch p.Status {
	case "waiting", "czar-waiting":
		return "Waiting"
	case "czar-ready":
		return "Ready to Start"
	case "play":
		return "Choose a Card"
	case "play-wait":
		return "Waiting on Judge"
	case "judge-wait":
		return "Waiting on Players"
	case "judge":
		return "Judging"
	default:
		return "Missing: " + p.Status
	}
}

func (p Player) ShowCards() bool {
	switch p.Status {
	case "play", "play-wait":
		return true
	default:
		return false
	}
}

func (p Player) Actionable() bool {
	switch p.Status {
	case "play", "judge":
		return false
	case "waiting", "czar-waiting":
		return false
	case "play-wait", "judge-wait":
		return false
	case "czar-ready":
		return true
	default:
		return false
	}
}
func (p Player) Ready() bool {
	switch p.Status {
	case "play-wait":
		return true
	}
	return false
}
func (p Player) ButtonMessage() string {
	switch p.Status {
	case "czar-ready":
		return "Start Game"
	case "play", "judge":
		return "Pick a Card"
	case "waiting", "czar-waiting":
		return "Waiting"
	case "play-wait", "judge-wait":
		return "Waiting"
	default:
		return "Unknown"
	}
}
func (p Player) ButtonDisabled() string {
	switch p.Status {
	case "czar-ready":
		return ""
	default:
		return "disabled"
	}
}

func (l *Lobby) PlayCard(p *Player, card string) {
	hi, ok := l.Instance.Players[p.ID]
	if !ok {
		return
	}
	for i, c := range l.Instance.Hands[hi] {
		// draw a new card, add mutex later
		if c.Name == card {
			l.Instance.CurrentPlays[p.ID] = c
			p.Status = "play-wait"
			l.Instance.Hands[hi][i] = l.Instance.PC[0]
			l.Instance.PC = l.Instance.PC[1:]
			b := &bytes.Buffer{}
			config.Tmpl.ExecuteTemplate(
				b,
				"player_show.html",
				map[string]interface{}{
					"Lobby":  l,
					"Player": p,
					"Hand":   l.Instance.Hands[hi],
				},
			)
			je, _ := json.Marshal(
				Event{
					Type: "new_hand",
					HTML: b.String(),
				},
			)
			l.PlayerChan[p.ID] <- je

			if len(l.Instance.CurrentPlays)+1 < len(l.Instance.Players) {
				je, _ = json.Marshal(
					Event{
						Type: "player_commit",
						Data: p.ID,
					},
				)
				l.Presenter <- je
				return
			} else {
				je, _ = json.Marshal(
					Event{
						Type: "start_judge",
					},
				)
				l.Presenter <- je

				j := l.Players[l.Instance.CurrentJudge]
				j.Status = "judge"
				b := &bytes.Buffer{}
				config.Tmpl.ExecuteTemplate(
					b,
					"player_show.html",
					map[string]interface{}{
						"Lobby":   l,
						"Player":  p,
						"Judging": l.Instance.CurrentPlays,
					})

				je, _ = json.Marshal(
					Event{
						Type: "start_judge",
						HTML: b.String(),
					},
				)
				l.PlayerChan[j.ID] <- je
			}
		}
	}
	l.Sync(p)
}

func (l *Lobby) PickCard(pid string) {
	pi := l.FindPlayer(pid)
	pi.RoundPoints++
	if pi.RoundPoints >= 6 {
		// do win things
		return
	}
	l.Instance.RoundWinner = pid
	l.Instance.RoundCard = l.Instance.CurrentPlays[pid]

	wm := fmt.Sprintf(
		`<div class="page-header"><h2>%s<small>%s</small></h2></div>`,
		l.Instance.CurrentWith(l.Instance.RoundCard),
		l.FindPlayer(pid).Name,
	)
	je, _ := json.Marshal(
		Event{
			Type: "pick_winner",
			HTML: wm,
			Data: pid,
		},
	)
	l.Presenter <- je

	je, _ = json.Marshal(
		Event{
			Type: "round_win",
		},
	)
	l.PlayerChan[pid] <- je
	for ci, ws := range l.PlayerChan {
		if ci != pid {
			je, _ = json.Marshal(
				Event{Type: "round_queue"},
			)
			ws <- je
		}
	}

	// start the next round in 30 seconds
	go func() {
		time.Sleep(10 * time.Second)
		l.NextRound()
	}()
}

func (l *Lobby) NextRound() {
	l.Instance.CurrentPlays = map[string]store.Card{}
	l.Instance.CurrentJudge++
	if l.Instance.CurrentJudge >= len(l.Players) {
		l.Instance.CurrentJudge = 0
	}
	l.Instance.AdvanceCard()
	for i, p := range l.Players {
		if i == l.Instance.CurrentJudge {
			p.Status = "judge-wait"
		} else {
			p.Status = "play"
		}
	}

	for pid, ws := range l.PlayerChan {
		b := &bytes.Buffer{}
		fmt.Println(l.FindPlayer(pid).Name, config.Tmpl.ExecuteTemplate(
			b,
			"player_show.html",
			map[string]interface{}{
				"Lobby":  l,
				"Player": l.FindPlayer(pid),
				"Hand":   l.Instance.Hands[l.Instance.Players[pid]],
			},
		))
		je, _ := json.Marshal(
			Event{
				Type: "new_round",
				HTML: b.String(),
			},
		)
		ws <- je

	}
	je, _ := json.Marshal(
		Event{
			Type: "game_start",
		},
	)
	l.Presenter <- je
}

func (l *Lobby) Sync(p *Player) {
	fmt.Println("Syncing", p.Name)
	ws, ok := l.PlayerChan[p.ID]
	if !ok {
		return
	}
	je, _ := json.Marshal(
		Event{
			Type: "sync",
		},
	)
	ws <- je

}
func FindPlayer(id string) *Player {
	return knownPlayers[id]
}

func CreatePlayer(name string) *Player {
	p := &Player{Name: name, Status: "waiting"}
	p.ID = base64.URLEncoding.EncodeToString(securecookie.GenerateRandomKey(16))
	knownPlayers[p.ID] = p
	return p
}

func (p Player) String() string {
	return p.Name
}

func (p Player) Judge() bool {
	switch p.Status {
	case "judge-wait", "judge":
		return true
	default:
		return false
	}
}

type Event struct {
	Type string
	HTML string
	Data interface{}
}
