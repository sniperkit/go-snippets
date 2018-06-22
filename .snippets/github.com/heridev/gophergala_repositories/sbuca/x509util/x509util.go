package x509util

import (
  "errors"
  "io/ioutil"
  "fmt"
  "crypto/x509"
  "crypto/rsa"
  "encoding/pem"
)

func PemStringToPemBlock(pemString string) (*pem.Block, error) {

  pemBlock, _ := pem.Decode([]byte(pemString))
  if pemBlock == nil {
    return nil, errors.New("pem decode failed")
  }

  return pemBlock, nil
}
func PemFileToPemBlock(filename string) (*pem.Block, error) {
  certData, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  pemBlock, _ := pem.Decode(certData)
  if pemBlock == nil {
    return nil, errors.New("pem decode failed")
  }

  return pemBlock, nil
}

func PemStringToCertificate(pemString string) (*x509.Certificate, error) {

  pemBlock, err := PemStringToPemBlock(pemString)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParseCertificate(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}
func PemFileToCertificate(filename string) (*x509.Certificate, error) {

  pemBlock, err := PemFileToPemBlock(filename)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParseCertificate(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}
func PemStringToCertificateRequest(pemString string) (*x509.CertificateRequest, error) {

  pemBlock, err := PemStringToPemBlock(pemString)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParseCertificateRequest(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}
func PemFileToCertificateRequest(filename string) (*x509.CertificateRequest, error) {

  pemBlock, err := PemFileToPemBlock(filename)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParseCertificateRequest(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}
func PemStringToRsaPrivateKey(pemString string) (*rsa.PrivateKey, error) {

  pemBlock, err := PemStringToPemBlock(pemString)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}
func PemFileToRsaPrivateKey(filename string) (*rsa.PrivateKey, error) {

  pemBlock, err := PemFileToPemBlock(filename)
  if err != nil {
    return nil, err
  }

  cert, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  return cert, nil
}

func main() {
  cert, err := PemFileToCertificate("ca.crt")
  if err != nil {
    panic(err)
  }
  fmt.Println(cert.Version)

  certData, err := ioutil.ReadFile("ca.crt")
  if err != nil {
    panic(err)
  }
  cert2, err := PemStringToCertificate(string(certData))
  if err != nil {
    panic(err)
  }
  fmt.Println(cert2.Version)
}
