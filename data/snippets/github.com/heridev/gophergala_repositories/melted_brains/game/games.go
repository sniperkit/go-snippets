package game

type Games []*Game

func (g Games) Find(id string) *Game {
	for _, game := range g {
		if game.Id == id {
			return game
		}
	}
	return nil
}

func (games Games) Joinable() Games {
	res := Games{}
	for _, g := range games {
		if g.Status == Created {
			res = append(res, g)
		}
	}
	return res
}
