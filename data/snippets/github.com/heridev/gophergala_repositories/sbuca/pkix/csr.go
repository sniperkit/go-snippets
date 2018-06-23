package pkix

import (

  "errors"
  "crypto/x509"
  "io/ioutil"
  "encoding/pem"
  "crypto/rand"
)

type CertificateRequest struct {

  DerBytes []byte

  Csr *x509.CertificateRequest

}

func NewCertificateRequest(key *Key) (*CertificateRequest, error) {
  template := &x509.CertificateRequest{
    Subject: GenSubject(""), //FIXME
    //Attributes:
    //SignatureAlgorithm,
    //Extensions:
    DNSNames: []string{},
    //EmailAddress:
    //IPAddresses:
  }

  derBytes, err := x509.CreateCertificateRequest(rand.Reader, template, key.PrivateKey)
  if err != nil {
    return nil, err
  }
  csr, err := NewCertificateRequestFromDER(derBytes)
  if err != nil {
    return nil, err
  }

  return csr, nil
}

func NewCertificateRequestFromDER(derBytes []byte) (*CertificateRequest, error) {

  csr, err := x509.ParseCertificateRequest(derBytes)
  if err != nil {
    return nil, err
  }

  certificateRequest := &CertificateRequest{
    DerBytes: derBytes,
    Csr: csr,
  }

  return certificateRequest, nil
}
func NewCertificateRequestFromPEM(pemBytes []byte) (*CertificateRequest, error) {

  pemBlock, _ := pem.Decode(pemBytes)
  if pemBlock == nil {
    return nil, errors.New("PEM decode failed")
  }

  csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
  if err != nil {
    return nil, err
  }

  certificateRequest := &CertificateRequest{
    DerBytes: pemBlock.Bytes,
    Csr: csr,
  }

  return certificateRequest, nil
}
func NewCertificateRequestFromPEMFile(filename string) (*CertificateRequest, error) {

  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  return NewCertificateRequestFromPEM(data)
}

func (csr *CertificateRequest) ToPEM() ([]byte, error) {

  pemBlock := &pem.Block{
    Type: "CERTIFICATE REQUEST",
    Bytes: csr.DerBytes,
  }

  pemBytes := pem.EncodeToMemory(pemBlock)

  return pemBytes, nil
}


