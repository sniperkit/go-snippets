package pkix

import (
  "encoding/pem"
  "crypto"
  "io/ioutil"
  "errors"
  "crypto/x509"
  "crypto/rsa"
  "crypto/rand"
)

type Key struct {
  /*
  PublicKey *crypto.PublicKey
  PrivateKey *crypto.PrivateKey
  */
  PublicKey crypto.PublicKey
  PrivateKey *rsa.PrivateKey
  DerBytes []byte
}

func NewKey() (*Key, error) {
  privateKey, err := rsa.GenerateKey(rand.Reader, 2048) 
  if err != nil {
    return nil, err
  }

  derBytes := x509.MarshalPKCS1PrivateKey(privateKey)
  if derBytes == nil {
    return nil, errors.New("marshal rsa failed")
  }

  newKey := &Key{
    PrivateKey: privateKey,
    PublicKey: privateKey.Public(),
    DerBytes: derBytes,
  }

  return newKey, nil
}
func NewKeyFromPrivateKeyPEM(pemBytes []byte) (*Key, error) {
  // currently we only support rsa

  pemBlock, _ := pem.Decode(pemBytes)
  if pemBlock == nil {
    return nil, errors.New("decode pem failed")
  }

  privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  newKey := &Key{
    PrivateKey: privateKey,
    PublicKey: privateKey.Public(),
    DerBytes: pemBlock.Bytes,
  }

  return newKey, nil
}
func NewKeyFromPrivateKeyPEMFile(filename string) (*Key, error) {

  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  return NewKeyFromPrivateKeyPEM(data)

}
func (key *Key) ToPEM() ([]byte, error) {

  pemBlock := &pem.Block{
    Type: "RSA PRIVATE KEY",
    Bytes: key.DerBytes,
  }
  pemBytes := pem.EncodeToMemory(pemBlock)

  return pemBytes, nil
}
func (key *Key) ToPEMFile(filename string) (error) {
  pemBytes, err := key.ToPEM()
  if err != nil {
    return err
  }

  return ioutil.WriteFile(filename, pemBytes, 0400)
}
