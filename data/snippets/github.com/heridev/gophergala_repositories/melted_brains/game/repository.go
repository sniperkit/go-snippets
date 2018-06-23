package game

import (
	"math/rand"
	"sync"
)

type GameRepository struct {
	Games
	sync.Mutex
}

func (gs *GameRepository) RandomJoinable() *Game {
	games := gs.Games.Joinable()
	if len(games) == 0 {
		return gs.Create()
	}
	return games[rand.Intn(len(games))]
}

func (gs *GameRepository) Find(id string) *Game {
	return gs.Games.Find(id)
}

func (gs *GameRepository) Delete(g *Game) {
	gs.Lock()
	newGames := Games{}
	for _, game := range gs.Games {
		if game != g {
			newGames = append(newGames, game)
		}
	}
	gs.Games = newGames
	gs.Unlock()
}

func (gs *GameRepository) Create() *Game {
	gs.Lock()
	game := NewGame()
	gs.Games = append(gs.Games, game)
	gs.Unlock()
	return game
}

var Repository = &GameRepository{Games: Games{}, Mutex: sync.Mutex{}}
