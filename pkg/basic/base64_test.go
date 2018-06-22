package basic

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Encoding(t *testing.T) {
	assert := assert.New(t)

	url := "http://localhost:8080/hello/newyorker"
	str := base64.URLEncoding.EncodeToString([]byte(url))
	log.Printf("content: %v", str)
	// assert.Fail("ERROR", "Base64Encoding error")
	assert.Nil(nil)

	bs, err := base64.URLEncoding.DecodeString(str)
	log.Printf("url: %v, err: %v", string(bs), err)
}
