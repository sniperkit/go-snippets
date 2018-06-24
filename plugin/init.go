// Package plugin provides plugin functions for sticker.
package plugin

import (
	"encoding/gob"
)

// InitializePlugin does nothing, because init functions in this package registers functions to sticker.
// Thus it is unnecessary for users to call any function in this package.
// Users can call this function for avoiding any import error.
func InitializePlugin() {
}

func init() {
	gob.Register(map[int]int(nil))
	gob.Register(map[uint32]int(nil))
	gob.Register(map[uint32]float32(nil))
	gob.Register(map[string]interface{}(nil))
}
