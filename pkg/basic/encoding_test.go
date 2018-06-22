package basic

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"testing"
)

func TestMd5(t *testing.T) {
	var data []byte
	md5sum := md5.Sum(data)
	md5str := hex.EncodeToString(md5sum[:])

	log.Println(md5str)
	log.Println(len(data))
}
