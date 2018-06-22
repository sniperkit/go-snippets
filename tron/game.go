package tron

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/glog"
)

type Color string

var Colors = []Color{"blue", "red", "green", "orange", "black", "purple"}

const ColorWall = "wall"

type JoinCmd struct {
	ColorC chan Color
	ArenaC chan Arena
}

type LeaveCmd struct {
	Color Color
}

type Direction string

const (
	DirectionUp    = "u"
	DirectionDown  = "d"
	DirectionLeft  = "l"
	DirectionRight = "r"
)

type MoveCmd struct {
	Color     Color
	Direction Direction
}

type Point struct {
	X int
	Y int
}

type Loser struct {
	Color       Color
	CollideWith Color
}

type Arena struct {
	Snakes map[Color][]Point
	Points map[Color]map[Point]struct{}
	Losers []Loser

	Size  Point
	Ratio float64
}

const (
	DefaultSizeX     = 1000
	DefaultSizeY     = 600
	DefaultSizeRatio = 4
)

func NewArena(snakes map[Color][]Point, ratio float64) *Arena {
	a := Arena{
		Snakes: snakes,
		Points: make(map[Color]map[Point]struct{}),
		Losers: make([]Loser, 0),
		Ratio:  ratio,
	}
	a.Size = Point{X: int(DefaultSizeX / a.Ratio), Y: int(DefaultSizeY / a.Ratio)}
	for color, s := range snakes {
		a.Points[color] = make(map[Point]struct{})
		for _, point := range s {
			a.Points[color][point] = struct{}{}
		}
	}
	return &a
}

// ChangeInitDirt sets the initial direction by altering the second point of the color's snake.
func (a *Arena) ChangeInitDirt(cmd MoveCmd) bool {
	snake := a.Snakes[cmd.Color]
	prevDirt := computeDirection(snake)
	changed := false
	if cmd.Direction != prevDirt {
		changed = true

		delete(a.Points[cmd.Color], snake[1])
		first := snake[0]
		switch cmd.Direction {
		case DirectionUp:
			snake[1] = Point{X: first.X, Y: first.Y + 1}
		case DirectionDown:
			snake[1] = Point{X: first.X, Y: first.Y - 1}
		case DirectionLeft:
			snake[1] = Point{X: first.X - 1, Y: first.Y}
		case DirectionRight:
			snake[1] = Point{X: first.X + 1, Y: first.Y}
		}
		a.Points[cmd.Color][snake[1]] = struct{}{}
	}
	return changed
}

func computeDirection(snake []Point) Direction {
	var prevDirt Direction
	last := snake[len(snake)-1]
	pntm := snake[len(snake)-2]
	if last.X == pntm.X {
		if last.Y > pntm.Y {
			prevDirt = DirectionUp
		} else {
			prevDirt = DirectionDown
		}
	} else {
		if last.X > pntm.X {
			prevDirt = DirectionRight
		} else {
			prevDirt = DirectionLeft
		}
	}
	return prevDirt
}

func oppositeDirections(a, b Direction) bool {
	if (a == DirectionUp && b == DirectionDown) || (a == DirectionDown && b == DirectionUp) || (a == DirectionLeft && b == DirectionRight) || (a == DirectionRight && b == DirectionLeft) {
		return true
	}
	return false
}

// Update updates the state of an arena for a timestep.
func (a *Arena) Update(acts map[Color]Direction) {
	for color, snake := range a.Snakes {
		lost := false
		for _, l := range a.Losers {
			if l.Color == color {
				lost = true
				break
			}
		}
		if lost {
			continue
		}

		prevDirt := computeDirection(snake)
		dirt := prevDirt
		if act, ok := acts[color]; ok && !oppositeDirections(act, prevDirt) {
			dirt = act
		}

		var p Point
		last := snake[len(snake)-1]
		switch dirt {
		case DirectionUp:
			p = Point{X: last.X, Y: last.Y + 1}
		case DirectionDown:
			p = Point{X: last.X, Y: last.Y - 1}
		case DirectionLeft:
			p = Point{X: last.X - 1, Y: last.Y}
		case DirectionRight:
			p = Point{X: last.X + 1, Y: last.Y}
		}

		// Check for collision with the wall
		if p.X <= 0 || p.X >= a.Size.X || p.Y <= 0 || p.Y >= a.Size.Y {
			a.Losers = append(a.Losers, Loser{Color: color, CollideWith: ColorWall})
			continue
		}
		if dirt != prevDirt {
			a.Snakes[color] = append(snake, p)
		} else {
			a.Snakes[color][len(snake)-1] = p
		}

		// Check collisions
		shouldProfile := false
		if rand.Intn(60) == 0 {
			shouldProfile = true
		}
		var t int64
		if shouldProfile {
			t = time.Now().UnixNano()
		}
		for otherColor, points := range a.Points {
			_, ok := points[p]
			if ok {
				a.Losers = append(a.Losers, Loser{Color: color, CollideWith: otherColor})
				break
			}
		}
		if shouldProfile {
			glog.Infof("collision time spent: %d ns", time.Now().UnixNano()-t)
		}
		a.Points[color][p] = struct{}{}
	}
}

type Player struct {
	Arena     chan *Arena
	GameEnd   chan Color // this is the color of the winner
	Countdown chan int
}

type Room struct {
	sync.RWMutex
	Players    map[*Player]struct{}
	MaxPlayers int
	Game       *Game

	Watchers map[*Player]struct{}
}

func NewRoom(maxPlayers int) *Room {
	r := Room{
		Players:    make(map[*Player]struct{}),
		MaxPlayers: maxPlayers,
		Watchers:   make(map[*Player]struct{}),
	}
	return &r
}

func (r *Room) Ready(player *Player) (*Game, Color) {
	r.Lock()
	defer r.Unlock()
	if r.Game == nil {
		r.Game = NewGame(r.MaxPlayers)
	}
	game := r.Game

	var color Color
	for _, c := range Colors {
		if _, ok := game.Players[c]; !ok {
			color = c
			break
		}
	}
	game.Players[color] = player

	if len(game.Players) >= game.MinPlayers {
		game.Watchers = r.Watchers
		go game.Start()
		r.Game = nil
	}
	return game, color
}

type Hall struct {
	sync.RWMutex
	m map[string]*Room
}

func (h *Hall) EnterRoom(name string, player *Player) (*Room, error) {
	h.Lock()
	defer h.Unlock()
	room, ok := h.m[name]
	if !ok {
		room = NewRoom(4)
		h.m[name] = room
	}

	if len(room.Players) >= room.MaxPlayers {
		return nil, fmt.Errorf("max players reached")
	}
	room.Players[player] = struct{}{}

	return room, nil
}

func (h *Hall) LeaveRoom(name string, player *Player) {
	h.Lock()
	defer h.Unlock()
	room, ok := h.m[name]
	if !ok {
		return
	}

	delete(room.Players, player)
	if len(room.Players) == 0 {
		delete(h.m, name)
	}
}

func (h *Hall) WatchRoom(name string, player *Player) error {
	h.Lock()
	defer h.Unlock()
	room, ok := h.m[name]
	if !ok {
		return fmt.Errorf("no such room")
	}

	room.Watchers[player] = struct{}{}

	return nil
}

func (h *Hall) UnwatchRoom(name string, player *Player) error {
	h.Lock()
	defer h.Unlock()
	room, ok := h.m[name]
	if !ok {
		return fmt.Errorf("no such room")
	}

	delete(room.Watchers, player)

	return nil
}

type Game struct {
	Players    map[Color]*Player
	MinPlayers int
	Move       chan MoveCmd

	Watchers map[*Player]struct{}
}

func NewGame(minPlayers int) *Game {
	game := Game{
		Players:    make(map[Color]*Player),
		MinPlayers: minPlayers,
		Move:       make(chan MoveCmd),
	}
	return &game
}

func (g *Game) Ended(a *Arena) bool {
	if len(g.Players)-1 <= len(a.Losers) {
		return true
	}
	return false

	// TODO
	for _, snake := range a.Snakes {
		if len(snake) > 5 {
			return true
		}
	}
	return false
}

func (g *Game) broadcastCountdown(cnt int) {
	for _, p := range g.Players {
		select {
		case p.Countdown <- cnt:
		default:
		}
	}
	for p, _ := range g.Watchers {
		select {
		case p.Countdown <- cnt:
		default:
		}
	}
}

func (g *Game) broadcastArena(arena *Arena) {
	for _, p := range g.Players {
		select {
		case p.Arena <- arena:
		default:
		}
	}
	for p, _ := range g.Watchers {
		select {
		case p.Arena <- arena:
		default:
		}
	}
}

func (g *Game) broadcastGameEnd(arena *Arena) {
	var winner Color
	for color, _ := range g.Players {
		lost := false
		for _, loser := range arena.Losers {
			if loser.Color == color {
				lost = true
			}
		}
		if !lost {
			winner = color
			break
		}
	}

	for _, p := range g.Players {
		select {
		case p.GameEnd <- winner:
		default:
		}
	}
	for p, _ := range g.Watchers {
		select {
		case p.GameEnd <- winner:
		default:
		}
	}
}

var initColors = []Point{
	Point{X: 333, Y: 200},
	Point{X: 666, Y: 200},
	Point{X: 333, Y: 400},
	Point{X: 666, Y: 400},
}

func (g *Game) Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cnt := 3; cnt >= 1; cnt -= 1 {
			g.broadcastCountdown(cnt)
			glog.Infof("Counting down %d", cnt)
			<-time.After(time.Second)
		}
	}()

	// Select initial direction
	snakes := make(map[Color][]Point)
	var ratio float64 = DefaultSizeRatio
	i := 0
	for color, _ := range g.Players {
		s := make([]Point, 2)
		s[0] = Point{X: int(float64(initColors[i].X) / ratio), Y: int(float64(initColors[i].Y) / ratio)}
		s[1] = Point{s[0].X + 1, s[0].Y}
		snakes[color] = s
		i += 1
	}
	arena := NewArena(snakes, ratio)

	timer := time.After(3 * time.Second)
InitDirt:
	for {
		select {
		case cmd := <-g.Move:
			if ok := arena.ChangeInitDirt(cmd); ok {
				g.broadcastArena(arena)
			}
		case <-timer:
			break InitDirt
		}
	}
	wg.Wait()
	g.broadcastCountdown(0)

	// Game begins!
	for {
		acts := make(map[Color]Direction)
		timer = time.After(50 * time.Millisecond)
	CollectActs:
		for {
			select {
			case cmd := <-g.Move:
				acts[cmd.Color] = cmd.Direction
			case <-timer:
				break CollectActs
			}
		}
		arena.Update(acts)
		g.broadcastArena(arena)

		if g.Ended(arena) {
			g.broadcastGameEnd(arena)
			return
		}
	}
}
