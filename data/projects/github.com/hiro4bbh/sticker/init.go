package sticker

import (
	"encoding/gob"
)

func init() {
	gob.Register(map[string]int(nil))
}
