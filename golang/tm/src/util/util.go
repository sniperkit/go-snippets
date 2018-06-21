package util

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func GenerateIDByString(str string) string {
	hash := sha1.New()
	return hex.EncodeToString(hash.Sum([]byte(str)))
}

func GenerateID() string {
	var id = "0x"

	rand.Seed(time.Now().Unix())
	result := rand.Perm(128)
	for _, i := range result {
		id = id + strconv.Itoa(i)
	}

	return id
}
