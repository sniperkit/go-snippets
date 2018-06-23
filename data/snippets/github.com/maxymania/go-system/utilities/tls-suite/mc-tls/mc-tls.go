/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

package main

import "strings"

import "flag"
import "crypto/rand"
import "crypto/rsa"
import "crypto/elliptic"
import "crypto/ecdsa"
import "crypto/x509"

import "crypto/x509/pkix"

import "math/big"
import "encoding/pem"
import "os"

import "syscall"
import "fmt"

import "net"
import "net/textproto"
import "bufio"
import "time"
import "io/ioutil"

var out = os.Stdout
var pkOut = os.Stdout
var bits = flag.Int("b", 1024, "bits (RSA)")
var primes = flag.Int("p", 2, "primes (RSA)")
var useRsa = flag.Bool("rsa", false, "RSA")
var useEcdsa = flag.Bool("ecdsa", false, "ECDSA")

var dumpSub = flag.Bool("dump-sub", false, "dump Cert Subject Parameters")
var loadSub = flag.Bool("load-sub", false, "load Cert Subject Parameters")

var dumpIp = flag.Bool("dump-dns", false, "dump Cert DNS/IP")
var loadIp = flag.Bool("load-dns", false, "load Cert DNS/IP")

//P256, P384 or P521
var ecdsaP256 = flag.Bool("p256", false, "ECDSA P256 (default)")
var ecdsaP384 = flag.Bool("p384", false, "ECDSA P384")
var ecdsaP521 = flag.Bool("p521", false, "ECDSA P521")

var subjCountry = flag.String("country","","Cert Subject's Country/s (,-seperated)")
var subjOrg = flag.String("org","","Cert Subject's Organization/s (,-seperated)")
var subjOu = flag.String("ou","","Cert Subject's OrganizationalUnit/s (,-seperated)")
var subjLocality = flag.String("loc","","Cert Subject's Locality/s (,-seperated)")
var subjProvince = flag.String("prov","","Cert Subject's Province/s (,-seperated)")
var subjStreetAddress = flag.String("addr","","Cert Subject's Street Address/es (,-seperated)")
var subjPostalCode = flag.String("postal","","Cert Subject's Postal Code/s (,-seperated)")
var subjSerialNumber = flag.String("serial","","Cert Subject's Serial Number")
var subjCommonName = flag.String("name","","Cert Subject's Common Name")

var certDNSNames = flag.String("dns-name","","Cert's DNS Name/s (,-seperated)")
var certEmailAddresses = flag.String("email","","Cert's Email Addresse/s (,-seperated)")
var certIPAddresses = flag.String("ip","","Cert's IP Addresse/s (,-seperated)")

var privateKeyName = flag.String("priv-key","","Private key destination file")
var certFileName = flag.String("cert-file","","Certificate destination file")

var duration   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")

var caTrue = flag.Bool("ca",false,"Certificate Autority?")

var signer = flag.String("signer","","Signer's Certificate")
var signerPriv = flag.String("signer-key","","Signer's Private Key")

type kdsak struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Priv    *big.Int
	Pub     *big.Int
}

var pkxName pkix.Name
var pkCert x509.Certificate
var DNSNames []string
var EmailAddresses []string
var IPAddresses []net.IP

var pkPriv interface{}

func getPK(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func s2ip(s []string) []net.IP {
	o := make([]net.IP,len(s))
	p := 0
	for i,z := range s {
		ipa := net.ParseIP(z)
		if ipa==nil {
			p++
		} else {
			o[i-p] = ipa
		}
	}
	return o[:len(s)-p]
}

func ip2s(s []net.IP) []string {
	o := make([]string,len(s))
	for i,z := range s {
		o[i] = z.String()
	}
	return o
}

func toSubArray(s string) []string{
	if s=="" { return nil }
	return strings.Split(s, ",")
}
func dumps(k string,v []string){
	for _,vv := range v {
		out.Write([]byte(k+": "+vv+"\n"))
	}
}

func parseSubject() error{
	var h textproto.MIMEHeader = nil
	if *loadSub || *loadIp {
		r := textproto.NewReader(bufio.NewReader(os.Stdin))
		hh,e := r.ReadMIMEHeader()
		if e!=nil { return e }
		h = hh
	}
	if *loadSub {
		pkxName.Country = h["Country"]
		pkxName.Organization = h["Org"]
		pkxName.OrganizationalUnit = h["Ou"]
		pkxName.Locality = h["Locality"]
		pkxName.Province = h["Province"]
		pkxName.StreetAddress = h["Street-Address"]
		pkxName.PostalCode = h["Postal-Code"]
		pkxName.SerialNumber = h.Get("Serial-Number")
		pkxName.CommonName = h.Get("Common-Name")
		EmailAddresses = h["Email"]
	}else{
		pkxName.Country = toSubArray(*subjCountry)
		pkxName.Organization = toSubArray(*subjOrg)
		pkxName.OrganizationalUnit = toSubArray(*subjOu)
		pkxName.Locality = toSubArray(*subjLocality)
		pkxName.Province = toSubArray(*subjProvince)
		pkxName.StreetAddress = toSubArray(*subjStreetAddress)
		pkxName.PostalCode = toSubArray(*subjPostalCode)
		pkxName.SerialNumber = *subjSerialNumber
		pkxName.CommonName = *subjCommonName
		EmailAddresses = toSubArray(*certEmailAddresses)
	}
	if *loadIp {
		DNSNames = h["Dns"]
		IPAddresses = s2ip(h["Ip"])
	}else{
		DNSNames = toSubArray(*certDNSNames)
		IPAddresses = s2ip(toSubArray(*certIPAddresses))
	}
	return nil
}

func makeCertObject() error{
	ctsBegin := time.Now()
	ctsEnd := ctsBegin.Add(*duration)
	sn128Bit := new(big.Int).Lsh(big.NewInt(1), 128)
	snRnd, e := rand.Int(rand.Reader, sn128Bit)
	if e!=nil { return e }
	pkCert = x509.Certificate{
		SerialNumber: snRnd,
		Subject: pkxName,
		NotBefore: ctsBegin,
		NotAfter: ctsEnd,
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	if *caTrue {
		pkCert.IsCA = true
		pkCert.KeyUsage |= x509.KeyUsageCertSign
	}
	return nil
}

func dumpSubject(){
	if *dumpSub {
		dumps("Country",pkxName.Country)
		dumps("Org",pkxName.Organization)
		dumps("Ou",pkxName.OrganizationalUnit)
		dumps("Locality",pkxName.Locality)
		dumps("Province",pkxName.Province)
		dumps("Street-Address",pkxName.StreetAddress)
		dumps("Postal-Code",pkxName.PostalCode)
		if pkxName.SerialNumber!="" {
			out.Write([]byte("Serial-Number: "+pkxName.SerialNumber+"\n"))
		}
		if pkxName.CommonName!="" {
			out.Write([]byte("Common-Name: "+pkxName.CommonName+"\n"))
		}
		dumps("Email",EmailAddresses)
	}
	if *dumpIp {
		dumps("Dns",DNSNames)
		dumps("Ip",ip2s(IPAddresses))
	}
	out.Write([]byte("\n"))
}

func genRsa() (e error) {
	var pk *rsa.PrivateKey
	if *primes>2 {
		pk,e = rsa.GenerateMultiPrimeKey(rand.Reader,*primes,*bits)
	} else {
		pk,e = rsa.GenerateKey(rand.Reader,*bits)
	}
	if e!=nil { return }
	pkPriv = pk
	b := x509.MarshalPKCS1PrivateKey(pk)
	e = pem.Encode(pkOut,&pem.Block{"RSA PRIVATE KEY",nil,b})
	return
}



func genEcdsa() (e error) {
	var pk *ecdsa.PrivateKey
	var b []byte
	curve := elliptic.P256()
	switch {
	case *ecdsaP384:curve = elliptic.P256()
	case *ecdsaP521:curve = elliptic.P256()
	}
	pk,e = ecdsa.GenerateKey(curve,rand.Reader)
	if e!=nil { return }
	pkPriv = pk
	b,e = x509.MarshalECPrivateKey(pk)
	if e!=nil { return }
	e = pem.Encode(pkOut,&pem.Block{"EC PRIVATE KEY",nil,b})
	return
}

func loadPem(s string) (*pem.Block,error) {
	f,e := ioutil.ReadFile(s)
	if e!=nil { return nil,e }
	b,_ := pem.Decode(f)
	return b,nil
}

func eat(i interface{}) {}

func generate() error{
	if *signer!="" && *signerPriv!="" {
		var lPriv interface{} = nil
		pb,e := loadPem(*signer)
		if e!=nil { return e }
		pk,e := loadPem(*signerPriv)
		if e!=nil { return e }
		cert,e := x509.ParseCertificate(pb.Bytes)
		if e!=nil { return e }
		switch pubk := cert.PublicKey.(type) {
		case *rsa.PublicKey:
			lPriv,e = x509.ParsePKCS1PrivateKey(pk.Bytes)
			if e!=nil { return e }
			eat(pubk)
		case *ecdsa.PublicKey:
			lPriv,e = x509.ParseECPrivateKey(pk.Bytes)
			if e!=nil { return e }
			eat(pubk)
		}
		cbin, err := x509.CreateCertificate(rand.Reader, &pkCert, cert, getPK(pkPriv), lPriv)
		if err!=nil { return err }
		return pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: cbin})
	}

	cbin, err := x509.CreateCertificate(rand.Reader, &pkCert, &pkCert, getPK(pkPriv), pkPriv)
	if err!=nil { return err }
	return pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: cbin})
}

func do_off() {
	syscall.Dup2(1,3)
	syscall.Dup2(2,1)
	out = os.NewFile(3,"dest")
}

func do_key() {
	var err error
	if *privateKeyName=="" { fmt.Fprintln(os.Stderr,"error: -priv-key is required"); os.Exit(0) }
	pkOut,err = os.Create(*privateKeyName)
	if err!=nil { fmt.Fprintln(os.Stderr,"error: bad file (-priv-key)"); os.Exit(0) }
	if *certFileName!="" {
		out,err = os.Create(*certFileName)
		if err!=nil { fmt.Fprintln(os.Stderr,"error: bad file (-cert-file)"); os.Exit(0) }
	}
}

func main() {
	do_off()
	var err error
	flag.Parse()
	err = parseSubject()
	if err!=nil { fmt.Fprintln(os.Stderr,"error parsing Subject data:",err); return }
	err = makeCertObject()
	if err!=nil { fmt.Fprintln(os.Stderr,"error making certificate:",err); return }
	switch {
	case *dumpSub || *dumpIp:
		dumpSubject()
	case *useRsa:
		do_key()
		err = genRsa()
		if err!=nil { fmt.Fprintln(os.Stderr,"error gen RSA:",err) }
		err = generate()
		if err!=nil { fmt.Fprintln(os.Stderr,"error creating Cert:",err) }
	case *useEcdsa:
		do_key()
		err = genEcdsa()
		if err!=nil { fmt.Fprintln(os.Stderr,"error gen ECDSA:",err) }
		err = generate()
		if err!=nil { fmt.Fprintln(os.Stderr,"error creating Cert:",err) }
	default: flag.PrintDefaults()
	}
}


