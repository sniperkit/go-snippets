package daihinmin

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

type reqResult struct {
	ok  bool
	err string
}

type doReq struct {
	sesh
	name   string
	result chan reqResult
}

func generateID(prefix string) string {
	const size = 30
	var b = make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error in generateID(%s): %v", prefix, err)
	}
	security := base64.URLEncoding.EncodeToString(b)
	id := prefix + security
	return id
}
