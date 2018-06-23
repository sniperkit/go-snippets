package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"math/big"

	"github.com/gitchander/go-examples/crypto/hex"
)

func main() {
	if err := Example2(); err != nil {
		fmt.Println(err)
	}
}

func Example1() error {

	// Generate private key

	ellipticCurve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(ellipticCurve, rand.Reader)
	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey

	//--------------------------------

	data := []byte("This is a message to be signed and verified by ECDSA!")

	var h hash.Hash
	h = sha256.New()
	h.Write(data)
	dataHash := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, dataHash)
	if err != nil {
		return err
	}

	format := "%s: [ %s ]\n"

	fmt.Printf(format, "r", hex.HexQuad(r.Bytes()))
	fmt.Printf(format, "s", hex.HexQuad(s.Bytes()))

	var signature []byte
	signature = append(signature, r.Bytes()...)
	signature = append(signature, s.Bytes()...)

	fmt.Printf(format, "signature", hex.HexQuad(signature))

	//r.SetBytes()

	ok := ecdsa.Verify(&publicKey, dataHash, r, s)
	fmt.Println(ok)

	return nil
}

func vSign(privateKey *ecdsa.PrivateKey, h hash.Hash, data []byte) ([]byte, error) {

	h.Reset()
	h.Write(data)
	dataHash := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, dataHash)
	if err != nil {
		return nil, err
	}

	var signature []byte
	signature = append(signature, r.Bytes()...)
	signature = append(signature, s.Bytes()...)

	result := append(data, signature...)

	return result, nil
}

func vVerify(publicKey *ecdsa.PublicKey, h hash.Hash, data []byte) (bool, error) {

	const digLen = 32

	n := len(data)

	if n < digLen*2 {
		err := errors.New("wrong data len")
		return false, err
	}

	r := new(big.Int)
	s := new(big.Int)

	dataLen := n - digLen*2

	k := dataLen
	r.SetBytes(data[k : k+digLen])

	k += digLen
	s.SetBytes(data[k : k+digLen])

	h.Reset()
	h.Write(data[:dataLen])
	dataHash := h.Sum(nil)

	ok := ecdsa.Verify(publicKey, dataHash, r, s)
	return ok, nil
}

func Example2() error {

	ellipticCurve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(ellipticCurve, rand.Reader)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	data := []byte("This is a message to be signed and verified by ECDSA!")

	h := sha256.New()

	dataSign, err := vSign(privateKey, h, data)
	if err != nil {
		return err
	}

	ok, err := vVerify(publicKey, h, dataSign)
	if err != nil {
		return err
	}

	fmt.Println(ok)

	return nil
}
