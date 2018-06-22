package daihinmin

import "crypto/rand"
import "math/big"
import "strconv"
import "fmt"
import "strings"

type Suit rune

const (
	Spade   Suit = '♠'
	Heart   Suit = '♥'
	Diamond Suit = '♦'
	Club    Suit = '♣'
	Joker   Suit = 'J'
)

type Rank int

const (
	Ace   Rank = 1
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13

	JokerRank Rank = 100
)

func (r Rank) Normal() int {
	switch {
	case r == Ace, r == Rank(2):
		return int(r + Rank(10))
	case r >= Rank(3) && r <= Rank(King):
		return int(r - Rank(3))
	case r >= JokerRank:
		return int(JokerRank)
	}
	return int(r)
}

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return strconv.Itoa(int(r))
	}
}

type Card struct {
	Suit
	Rank
}

var (
	Spe3 = Card{Spade, 3}
	Dia3 = Card{Diamond, 3}
)

func (c Card) String() string {
	return fmt.Sprintf("%s%s", string(c.Suit), c.Rank.String())
}

type Cards []Card

func (c Cards) Shuffle() {
	for i := range c {
		max := big.NewInt(int64(i + 1))
		jBig, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		j := jBig.Int64()
		c[i], c[j] = c[j], c[i]
	}
}

// TODO: JOKER
func (c Cards) Trick() Trick {
	if len(c) == 1 {
		return Single
	} else if len(c) > 1 {
		// is it a pair, triplet, etc?
		same := true
		for _, card := range c {
			if card.Rank != c[0].Rank {
				same = false
			}
		}
		if same {
			return Tuple
		}

		// test for kaidan
		same = true
		for i := 1; i < len(c); i++ {
			last, this := c[i-1], c[i]
			if this.Normal() != last.Normal()+1 {
				same = false
			}
			if this.Suit != last.Suit {
				same = false
			}
		}
		if same {
			return Kaidan
		}
	}
	return NoTrick
}

func (c Cards) CanPlayOn(top Cards, reverse bool, shibari bool) bool {
	trick := top.Trick()
	if c.Trick() != trick {
		return false
	}

	switch trick {
	case Kaidan:
		fallthrough
	case Tuple:
		if len(c) != len(top) {
			return false
		}
		fallthrough
	case Single:
		if !reverse {
			if c[0].Normal() < top[0].Normal() {
				return false
			}
		} else {
			if c[0].Normal() > top[0].Normal() {
				return false
			}
		}
	}

	if shibari {
		if !c.CanShibari(top) {
			return false
		}
	}

	return true
}

func (c Cards) CanShibari(top Cards) bool {
	return c.SuitCount().Equals(top.SuitCount())
}

func (c Cards) SameSuit() (same bool, suit Suit) {
	if len(c) == 0 {
		same = false
		return
	}

	var s *Suit
	for _, card := range c {
		if s == nil {
			s = &card.Suit
		}
		if card.Suit != *s {
			same = false
			return
		}
	}

	same, suit = true, *s
	return
}

func (c Cards) SameRank() (same bool, rank Rank) {
	if len(c) == 0 {
		same = false
		return
	}

	var r *Rank
	for _, card := range c {
		if r == nil {
			r = &card.Rank
		}
		if card.Rank != *r {
			same = false
			return
		}
	}

	same, rank = true, *r
	return
}

func (c Cards) SuitCount() SuitCount {
	ct := make(SuitCount)
	for _, card := range c {
		ct[card.Suit] = ct[card.Suit] + 1
	}
	return ct
}

func (c Cards) Map() map[Card]bool {
	m := make(map[Card]bool)
	for _, card := range c {
		m[card] = true
	}
	return m
}

func (c Cards) HasAll(test Cards) bool {
	m := c.Map()
	for _, card := range test {
		if !m[card] {
			return false
		}
	}
	return true
}

func (c Cards) Without(remove Cards) Cards {
	var fresh Cards
	deleteMap := remove.Map()
	for _, card := range c {
		if !deleteMap[card] {
			fresh = append(fresh, card)
		}
	}
	return fresh
}

// sorting stuff
func (c Cards) Len() int { return len(c) }
func (c Cards) Less(i, j int) bool {
	return c[i].Normal() < c[j].Normal()
}
func (c Cards) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c Cards) String() string {
	var s []string
	for _, card := range c {
		s = append(s, card.String())
	}
	return strings.Join(s, " ")
}

func NewDeck(jokers int) Cards {
	var deck Cards

	for _, s := range []Suit{Spade, Heart, Diamond, Club} {
		for i := Ace; i <= King; i++ {
			deck = append(deck, Card{s, i})
		}
	}

	for i := 0; i < jokers; i++ {
		deck = append(deck, Card{Joker, JokerRank + Rank(i)})
	}

	deck.Shuffle()
	return deck
}

type Trick int

const (
	NoTrick Trick = iota
	Single
	Tuple
	Kaidan
)

type SuitCount map[Suit]int

func (sc SuitCount) Equals(cmp SuitCount) bool {
	diff := 0
	diff += abs(sc[Spade] - cmp[Spade])
	diff += abs(sc[Heart] - cmp[Heart])
	diff += abs(sc[Diamond] - cmp[Diamond])
	diff += abs(sc[Club] - cmp[Club])
	diff += abs(sc[Joker] - cmp[Joker]) // for joker tanpin only
	return diff == 0
}

type Pile []Cards

func (p Pile) Top() Cards {
	return p[len(p)-1]
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
