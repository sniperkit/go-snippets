package daihinmin

import "sort"

type Player struct {
	Name    string
	Number  int
	Hand    Cards
	Miracle *Card

	Place int
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

func (p *Player) Give(c Card) {
	p.Hand = append(p.Hand, c)
}

func (p *Player) HasCard(c Card) bool {
	for _, card := range p.Hand {
		if c == card {
			return true
		}
	}
	return false
}

type Game struct {
	ID      string
	Players []*Player
	Current int
	Rules   Rules
	Deck    Cards
	Started bool

	Pile
	KakumeiChuu bool
	Shibari     bool
	LastPlay    *Player
}

func NewGame() *Game {
	return &Game{
		ID:    generateID("g:"),
		Rules: FuchuuRules,
	}
}

func (g *Game) Join(p *Player) {
	p.Number = len(g.Players)
	g.Players = append(g.Players, p)
}

func (g *Game) Start() {
	g.Deck = NewDeck(g.Rules.Jokers)
	playerCount := g.PlayerCount()

	dealUntil := g.Rules.Cards * playerCount
	if dealUntil > len(g.Deck) {
		dealUntil = len(g.Deck)
	}

	// deal everyone's basic hand
	for i := 0; i < dealUntil; i++ {
		pID := i % playerCount
		player := g.Players[pID]
		card, ok := g.draw()
		if ok {
			player.Give(card)
		}
	}

	// kiseki
	if g.Rules.Miracle && len(g.Deck) >= playerCount {
		for _, player := range g.Players {
			card, ok := g.draw()
			if ok {
				player.Miracle = &card
			}
		}
	}

	// sort everyone's hands
	for _, player := range g.Players {
		sort.Sort(player.Hand)
	}

	// find the first player
	first := g.firstPlayer()
	g.Current = first.Number
	g.Started = true
}

func (g *Game) Play(player *Player, cards Cards) (ok bool, events []Event) {
	//cards.TransformJokers()
	if g.Current != player.Number {
		// wrong player
		ok = false
		return
	}
	if len(cards) == 0 {
		// didn't play anything
		ok = false
		return
	}
	if !player.Hand.HasAll(cards) {
		// tried to play cards he doesn't have
		ok = false
		return
	}
	if len(g.Pile) != 0 {
		if !cards.CanPlayOn(g.Pile.Top(), g.KakumeiChuu, g.Shibari) {
			// can't play these cards
			ok = false
			return
		} else {
			if g.Rules.Shibari {
				if !g.Shibari && cards.CanShibari(g.Pile.Top()) {
					g.Shibari = true
					events = append(events, Shibatta)
				}
			}
		}
	}

	ok = true

	if g.Rules.Kakumei {
		if len(cards) >= g.Rules.KakumeiMin {
			if !(g.Rules.KaidanKakumei && cards.Trick() == Kaidan) {
				g.KakumeiChuu = !g.KakumeiChuu
				events = append(events, KakumeiOkoshita)
			}
		}
	}

	g.Pile = append(g.Pile, cards)
	player.Hand = player.Hand.Without(cards)
	var next *Player

	if g.Rules.HachiGiri {
		ok, rank := cards.SameRank()
		if ok && rank == Rank(8) {
			g.clearPile()
			events = append(events, HachiGiri)
			events = append(events, ClearPile)
		}
		next = player
	}

	if next == nil {
		next = g.nextPlayer()
	}
	g.Current = next.Number
	g.LastPlay = player
	// TODO: annonce next player
	return
}

func (g *Game) Pass(player *Player) (events []Event) {
	next := g.nextPlayer()
	if next == g.LastPlay {
		// play is returning to the person who last played a card
		// so we clear the pile
		g.clearPile()
		events = append(events, ClearPile)
	}
	g.Current = next.Number
	return
}

func (g *Game) PlayerCount() int {
	return len(g.Players)
}

func (g *Game) draw() (c Card, ok bool) {
	if len(g.Deck) > 1 {
		c, ok = g.Deck[0], true
		g.Deck = g.Deck[1:]
	} else if len(g.Deck) == 1 {
		c, ok = g.Deck[0], true
		g.Deck = []Card{}
	} else {
		ok = false
	}
	return
}

func (g *Game) firstPlayer() *Player {
	for rank := Rank(3); rank <= King; rank++ {
		for _, player := range g.Players {
			if player.HasCard(Card{Diamond, rank}) {
				return player
			}
		}
	}

	// no diamonds? fuck it
	return g.Players[0]
}

func (g *Game) nextPlayer() *Player {
	for cur, next := g.Current, g.Players[(g.Current+1)%g.PlayerCount()]; next.Place == 0; cur++ {
		return next
	}
	return nil
}

func (g *Game) clearPile() {
	g.Pile = nil
	g.Shibari = false
}

type Rules struct {
	Kakumei       bool
	Kaidan        bool
	KaidanKakumei bool
	HachiGiri     bool
	Spe3Gaeshi    bool
	Shibari       bool
	JackBack      bool
	Miracle       bool

	Cards      int
	Jokers     int
	KaidanMin  int
	KakumeiMin int
}

var FuchuuRules = Rules{
	Kakumei:       true,
	Kaidan:        true,
	KaidanKakumei: true,
	HachiGiri:     true,
	Spe3Gaeshi:    true,
	JackBack:      true,
	Shibari:       true,
	Miracle:       true,

	Cards:      9,
	Jokers:     2,
	KaidanMin:  3,
	KakumeiMin: 4,
}

type Event string

const (
	Shibatta        Event = "shibari"
	KakumeiOkoshita       = "kakumei"
	HachiGiri             = "8giri"
	ClearPile             = "clear"
	ErrorEvent            = "error"
)
