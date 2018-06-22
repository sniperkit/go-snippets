package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type CertAuthority struct {
	CertFile, KeyFile string
	root              *tls.Certificate
}

func NewCertAuthority(certFile, keyFile string) *CertAuthority {
	return &CertAuthority{
		CertFile: certFile, KeyFile: keyFile,
	}
}

func (ca *CertAuthority) Init() error {
	//	if err = os.MkdirAll(certDirectory, 0755); err != nil {
	//		return
	//	}
	var err error
	if fileExists(ca.CertFile) && fileExists(ca.KeyFile) {
		root, err := tls.LoadX509KeyPair(ca.CertFile, ca.KeyFile)
		if err != nil {
			return err
		}
		ca.root = &root
		ca.root.Leaf, err = x509.ParseCertificate(ca.root.Certificate[0])
		return nil
	}
	now := time.Now()
	template := &x509.Certificate{
		IsCA: true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte("ket"),
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Earth"},
			Organization: []string{"Ket Certificate Authority"},
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(5, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	// create a self-signed certificate.
	ca.root, err = generate(template, nil, ca.CertFile, ca.KeyFile)
	ca.root.Leaf, err = x509.ParseCertificate(ca.root.Certificate[0])
	return err
}

func (ca *CertAuthority) Get(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certFile, keyFile := getCertFileName(clientHello.ServerName)

	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err == nil {
		return &tlsCert, nil
	}
	hostName := clientHello.ServerName
	//log.Printf("Creating SSL certificate for %s...", hostName)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   hostName,
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	if ip := net.ParseIP(hostName); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, hostName)
	}

	return generate(template, ca.root, certFile, keyFile)
}

func getCertFileName(host string) (string, string) {
	h := sha1.New()
	h.Write([]byte(host))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	certFile := "./data/tmp/" + hash + ".crt"
	keyFile := "./data/tmp/" + hash + ".key"
	return certFile, keyFile
}

func generate(template *x509.Certificate, root *tls.Certificate, certFile, keyFile string) (*tls.Certificate, error) {
	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	parent := template
	//otherKey := crypto.PrivateKey(privatekey)
	var cert []byte
	if root != nil {
		parent = root.Leaf
		//otherKey = root.PrivateKey

		cert, err = x509.CreateCertificate(rand.Reader, template, parent, &privatekey.PublicKey, root.PrivateKey)
	} else {
		cert, err = x509.CreateCertificate(rand.Reader, template, parent, &privatekey.PublicKey, privatekey)
	}
	if err != nil {
		return nil, err
	}
	err = write(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return nil, err
	}
	err = write(keyFile, &pem.Block{Type: "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey)})
	return nil, err

	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	return &tlsCert, err

}

func write(filename string, block *pem.Block) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	pem.Encode(file, block)
	file.Close()
	return nil
}
