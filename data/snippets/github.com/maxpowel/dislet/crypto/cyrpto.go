package crypto

import (

	"github.com/fatih/color"
	"github.com/maxpowel/dislet"
	"crypto/aes"
	"io"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type Config struct {
	Key string
}


type Crypto struct {
	key []byte
}

func (c *Crypto) EncryptString(plaintext string) (string, error) {
	data, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func (c *Crypto) DecryptString(inputData string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(inputData)
	if err != nil {
		return "", err
	}
	data, err := c.Decrypt(ciphertext)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Crypto) Encrypt(data []byte) ([]byte, error) {

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream2 := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream2.XORKeyStream(data, data)
	//fmt.Printf("%s", []byte(ciphertext2))
	return data, nil
}

func NewCrypto(key string) *Crypto {
	binKey, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	return &Crypto{
		key: binKey,
	}
}


func Bootstrap(k *dislet.Kernel) {
	//fmt.Println("DATABASE BOOT")
	mapping := k.Config.Mapping
	mapping["crypto"] = &Config{}

	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		color.Green("Booting Crypto")
		conf := k.Config.Mapping["crypto"].(*Config)
		k.Container.RegisterType("crypto", NewCrypto, conf.Key)
		/*
		fmt.Println("PROBANDO", conf.Key)
		c := NewCrypto(conf.Key)
		d,err := c.EncryptString("lolazo")
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(d)
		d2,err := c.DecryptString(d)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(d2)
		*/
	}
	k.Subscribe(baz)

}
