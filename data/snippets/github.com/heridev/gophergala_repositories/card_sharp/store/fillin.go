package store

func SetupFillIn(c *Conn) {
	d := Deck{
		Name:        "Fill in the Blanks",
		Description: "Players try to complete a phrase in the funniest manner, using the cards in their hand. Players rotate through judging duties where they must choose the best card from the other players that completes the phrase.",
		FullGame:    true,
		GameType:    "blanks",
		AccountID:   1,
		MinPlayer:   3,
	}
	d.Save(c)
	cards := []Card{
		Card{
			Name:   "The Door was the Way to %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Oh God, not %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "I'm not dead, it's just %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "%[1]s Considered Harmful",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Friends come from %[1]s, they aren't made",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "What is life without the presence of %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "%[1]s, the root of all joy and sorrow",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Some day, the world will discover a lack of %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "Don't disturb my %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},
		Card{
			Name:   "It's about the experience, we'll bond over %[1]s",
			Type:   "parent",
			DeckID: d.ID,
		},

		// child cards
		Card{
			Name:   "Broken Dreams",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Sleeping in the Park",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Mystery Meat Monday",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Sandals over Socks",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Love",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "The Grim Reaper",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Reanimated Bigfoot",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Unsafe Carnival Rides",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Sleepy Kittens",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "My Birthday Celebrations",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Tiger Selfies",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Tax Day",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Growing Older",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Philosophy Books",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Macaroni Paintings",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Time Machine Shenanigans",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Leprechauns",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Financial Responsibility",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Behaving Rationally",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Plagues of Frogs",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "The Pumpkin Emperor",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hashtags",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Hippy Caravans",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Airborne Narcolepsy",
			Type:   "child",
			DeckID: d.ID,
		},
		Card{
			Name:   "Urban Goatherding",
			Type:   "child",
			DeckID: d.ID,
		},
	}
	c.Card.SaveAll(cards)
}
