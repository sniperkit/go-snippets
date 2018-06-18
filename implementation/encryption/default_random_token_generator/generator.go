package default_random_token_generator

import (
	"crypto/rand"

	"github.com/go-zero-boilerplate/ddd/interfaces/encryption"
)

type generator struct{}

func (g *generator) generateRandomString(length int, dictionary string) string {
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func (g *generator) AlphaNumericString(length int) string {
	return g.generateRandomString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

func (g *generator) AlphaString(length int) string {
	return g.generateRandomString(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

func (g *generator) NumericString(length int) string {
	return g.generateRandomString(length, "0123456789")
}

func (g *generator) CustomString(length int, charlist string) string {
	return g.generateRandomString(length, charlist)
}

func New() encryption.RandomTokenGenerator {
	return &generator{}
}
