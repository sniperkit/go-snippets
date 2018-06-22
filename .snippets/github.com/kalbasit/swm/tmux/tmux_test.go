package tmux

import (
	"io/ioutil"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// discard logs
	log.Logger = zerolog.New(ioutil.Discard)
}
