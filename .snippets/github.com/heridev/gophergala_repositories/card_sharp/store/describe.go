package store

func SetupDescribe(c *Conn) {
	d := Deck{
		Name:        "Description Roulette",
		Description: "Players try to match their cards to the topic. Are Banks Dreary? You have to decide. But watch out, the topic can get modified so you may get a topic of 'The Furthest Thing From Dreary'.",
		FullGame:    true,
		GameType:    "adjective",
		AccountID:   1,
		MinPlayer:   3,
	}
	d.Save(c)
	cards := []Card{
		Card{
			Name:   "Dreary",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Avant-Garde",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Enlightening",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Jubilant",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Serious",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hopeless",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Lively",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Lazy",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Powerful",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Impressive",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Weak",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Wicked",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Slow",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Fun",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Enormous",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Exclusive",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Plentiful",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Essential",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Cheap",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Agile",
			Type:   "parent",
			DeckID: d.ID,
		},

		// child cards
		Card{
			Name:   "Leonardo da Vinci",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Cookies",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Pickup Trucks",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Chocolate Chip Cookies",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hope",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Smartphones",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Trains",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Couches",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Kittens",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Birthdays",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Giraffes",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Pencils",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Grandparents",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Forks",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hammer",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Televisions",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Children",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Winter Jackets",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Austrailia",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Paris",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Beagles",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hashtags",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Books",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Banks",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Milk",
			Type:   "child",
			DeckID: d.ID,
		},
	}
	c.Card.SaveAll(cards)
}
