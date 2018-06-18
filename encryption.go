package toolkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func MD5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(baseChars string, n int) string {
	if baseChars == "" {
		baseChars = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz-_"
	}
	baseCharsLen := len(baseChars)

	rnd := ""
	for i := 0; i < n; i++ {
		x := RandInt(baseCharsLen)
		rnd += string(baseChars[x])
	}
	return rnd
}

func FileChecksum(fileLocation string) string {
	f, err := os.Open(fileLocation)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return ""
	}

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return ""
	}

	hashed := hash.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}

func EncryptAES(text, key string) (string, error) {
	plaintext := []byte(text)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	result := fmt.Sprintf("%x", ciphertext)
	return result, nil
}

func DecryptAES(text, key string) (string, error) {
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	res, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%s", res)
	return result, nil
}
