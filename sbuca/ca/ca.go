package ca

import (
  "strconv"
  "time" 
  //"errors"
  "math/big"
  "io/ioutil"
  "strings"
  "github.com/gophergala/sbuca/pkix"
  "crypto/x509"
  "crypto/rand"
  "os"
)


type CA struct {
  RootDir string
  CertStore *CertStore
  Certificate *pkix.Certificate
  Key *pkix.Key
}

func isPathNotExisted(path string) bool {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return true
  }
  return false
}

func NewCA(rootDir string) (*CA, error) {

  // mkdir if needed
  if isPathNotExisted(rootDir + "/ca") {
    if err := os.Mkdir(rootDir + "/ca", 0755); err != nil {
      return nil, err
    }
  }

  if isPathNotExisted(rootDir + "/certs") {
    if err := os.Mkdir(rootDir + "/certs", 0755); err != nil {
      return nil, err
    }
  }

  var key *pkix.Key
  var certificate *pkix.Certificate
  var err error
  if isPathNotExisted(rootDir + "/ca/ca.key") {
    // gen priv key
    key, err = pkix.NewKey()
    if err != nil {
      return nil, err
    }
    if err := key.ToPEMFile(rootDir + "/ca/ca.key"); err != nil {
      return nil, err
    }

    // gen self-signed cert
    // should refactor, move to cert.go
    notBefore := time.Now()
    notAfter  := notBefore.Add(time.Hour*365*24)
    keyUsage  := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
    extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
    template := &x509.Certificate{
      SerialNumber: big.NewInt(1),
      Subject: pkix.GenSubject("try.subca.com"),
      NotBefore: notBefore,
      NotAfter: notAfter,
      KeyUsage: keyUsage,
      ExtKeyUsage: extKeyUsage,
      BasicConstraintsValid: true,
    }

    derBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.PublicKey, key.PrivateKey)
    if err != nil {
      return nil, err
    }
    certificate, err = pkix.NewCertificateFromDER(derBytes) 
    if err != nil {
      return nil, err
    }
    if err := certificate.ToPEMFile(rootDir + "/ca/ca.crt"); err != nil {
      return nil, err
    }

  } else {

    certificate, err = pkix.NewCertificateFromPEMFile(rootDir + "/ca/ca.crt")
    if err != nil {
      return nil, err
    }
    key, err = pkix.NewKeyFromPrivateKeyPEMFile(rootDir + "/ca/ca.key")
    if err != nil {
      return nil, err
    }

  }

  if isPathNotExisted(rootDir + "/ca/ca.srl") {
    ioutil.WriteFile(rootDir + "/ca/ca.srl", []byte("2"), 0644)
  }

  certStore := NewCertStore(rootDir + "/certs")
  newCA := &CA{
    RootDir: rootDir,
    CertStore: certStore,
    Certificate: certificate,
    Key: key,
  }

  return newCA, nil
}
func (ca *CA) GetCertificate(id int64) (*pkix.Certificate, error){
  return ca.CertStore.Get(id)
}
func (ca *CA) PutCertificate(id int64, cert *pkix.Certificate) error {
  return ca.CertStore.Put(id, cert)
}
func (ca *CA) GetSerialNumber() (*big.Int, error) {
  snStr, err := ioutil.ReadFile(ca.RootDir + "/ca/ca.srl")
  if err != nil {
    panic(err)
  }
  snInt, err := strconv.Atoi(strings.Trim(string(snStr), "\n"))
  if err != nil {
    panic(err)
  }
  sn := big.NewInt(int64(snInt))

  return sn, nil
}
func (ca *CA) IncreaseSerialNumber() error {
  snStr, err := ioutil.ReadFile(ca.RootDir + "/ca/ca.srl")
  if err != nil {
    panic(err)
  }
  snInt, err := strconv.Atoi(strings.Trim(string(snStr), "\n"))
  if err != nil {
    panic(err)
  }
  nextSnInt := snInt + 1
  nextSnStr := strconv.Itoa(nextSnInt) + "\n"
  ioutil.WriteFile(ca.RootDir + "/ca/ca.srl", []byte(nextSnStr), 0600)

  return nil
}
func (ca *CA) IssueCertificate(csr *pkix.CertificateRequest) (*pkix.Certificate, error) {

  serialNumber, err := ca.GetSerialNumber()
  if err != nil {
    return nil, err
  }
  notBefore := time.Now()
  notAfter  := notBefore.Add(time.Hour*365*24)
  keyUsage  := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
  extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
  template := &x509.Certificate{
    SerialNumber: serialNumber,
    Subject: csr.Csr.Subject,
    NotBefore: notBefore,
    NotAfter: notAfter,
    KeyUsage: keyUsage,
    ExtKeyUsage: extKeyUsage,
    BasicConstraintsValid: true,
  }

  derBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate.Crt, ca.Key.PublicKey, ca.Key.PrivateKey)
  if err != nil {
    return nil, err
  }

  // increase sn
  if err = ca.IncreaseSerialNumber(); err != nil {
    return nil, err
  }

  // gen new cert
  cert, err := pkix.NewCertificateFromDER(derBytes)
  if err != nil {
    return nil, err
  }

  // put in certstore
  if err = ca.PutCertificate(serialNumber.Int64(), cert); err != nil {
    return nil, err
  }

  return cert, nil
}
