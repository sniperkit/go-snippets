package game

import (
	"fmt"
	"sync"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestGameDirector(t *testing.T) {
	t.Logf("game main go ...")
	name := "test_game_director"
	cfg := Config{Name: &name}
	gd := NewGameDirector(2, 100, &cfg)
	gd.Start()
	defer gd.Stop()

	var wg sync.WaitGroup
	var games []interface{}
	for index := 1; index <= 5; index++ {
		gfunc := NewSimpleGame(fmt.Sprintf("simpleGame#%d", index), func() chan error {
			var errc chan error = make(chan error)
			go func() {
				select {
				case <-time.After(time.Millisecond * 500):
					wg.Done()
					select {
					case errc <- nil:
					default:
					}
				case <-errc:
					return
				}
			}()
			return errc
		})
		games = append(games, gfunc)
	}

	var gs []Game
	for _, g := range games {
		wg.Add(1)
		game, err := gd.Schedule(g)
		if err != nil {
			t.Errorf("game test schedule game error: %s", err)
		}
		gs = append(gs, game)
	}

	for _, g := range gs {
		if g.Id() == "simpleGame#2" || g.Id() == "simpleGame#7" {
			err := gd.Cancel(g.Id())
			if err != nil {
				continue
			}
			wg.Done()
		}
	}

	wg.Wait()
	time.Sleep(time.Second)
	pending := gd.Pending()
	if pending != 0 {
		t.Fatalf("game director has not complete games %d", pending)
	}
	running := gd.Running()
	if running != 0 {
		t.Fatalf("game director has not finished games %d", running)
	}
}
