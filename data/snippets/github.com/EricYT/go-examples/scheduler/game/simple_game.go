package game

import (
	"time"

	log "github.com/Sirupsen/logrus"

	tomb "gopkg.in/tomb.v1"
)

var slog = log.WithFields(log.Fields{
	"module": "simple_game",
})

type StartFunc func() chan error

type simpleGame struct {
	tomb *tomb.Tomb
	id   string
	f    StartFunc
	cfg  *Config `inject`
}

func NewSimpleGame(id string, f StartFunc) interface{} {
	// return a function that can be injectd
	return func(cfg *Config) *simpleGame {
		return &simpleGame{
			tomb: new(tomb.Tomb),
			id:   id,
			f:    f,
			cfg:  cfg,
		}
	}
}

func (s simpleGame) Id() string { return s.id }

func (s *simpleGame) Wait() error {
	return s.tomb.Wait()
}

func (s *simpleGame) Kill() {
	s.tomb.Kill(nil)
}

func (s *simpleGame) Run() error {
	slog.Debugf("simple game %s run. config: %s now: %s", s.id, *s.cfg.Name, time.Now())
	var signal chan error = s.f()
	select {
	case err := <-signal:
		slog.Debugf("simple game %s done. now: %s", s.id, time.Now())
		return err
	case <-s.tomb.Dying():
		slog.Debugf("simple game %s was killed. now: %s", s.id, time.Now())
		signal <- nil
		return nil
	}
}
