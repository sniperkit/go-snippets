package ca

import (
  "github.com/gophergala/sbuca/pkix"
  "strconv"
  "io/ioutil"
)

type CertStore struct {

  RootDir string

}

func NewCertStore(rootDir string) *CertStore{

  store := &CertStore {
    RootDir: rootDir,
  }

  return store
}

func (store *CertStore) Get(id int64) (*pkix.Certificate, error) {
  // FIXME
  // currently using serialnumber as id, should change to something which can be
  // mapped to (host, sn) pair
  filename := strconv.Itoa(int(id)) + ".crt"

  cert, err := pkix.NewCertificateFromPEMFile(store.RootDir + "/" + filename)
  if err != nil {
    return nil, err
  }

  return cert, nil
}

func (store *CertStore) Put(id int64, cert *pkix.Certificate) error {

  pemBytes, err := cert.ToPEM()
  if err != nil {
    return err
  }
  filename := strconv.Itoa(int(id)) + ".crt"
  err = ioutil.WriteFile(store.RootDir + "/" + filename, pemBytes, 0400)
  if err != nil {
    return err
  }

  return nil
}

func (store *CertStore) GetAllNames() ([]string, error) {

  files, err := ioutil.ReadDir(store.RootDir + "/")
  if err != nil {
    return nil, err
  }
  names := make([]string, len(files))
  for _, f := range files {
    names = append(names, f.Name())
  }

  return names, nil

}
  // should limit to 100 FIXME
