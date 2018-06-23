/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 Windows compatible version of mc-ssh.
 */
package main

import "flag"
import "crypto/rand"
import "crypto/rsa"
import "crypto/dsa"
import "crypto/elliptic"
import "crypto/ecdsa"
import "crypto/x509"
import "math/big"
import "encoding/pem"
import "encoding/asn1"
import "os"

import "fmt"

var out = os.Stdout
var bits = flag.Int("b", 1024, "bits (RSA)")
var primes = flag.Int("p", 2, "primes (RSA)")
var useRsa = flag.Bool("rsa", false, "RSA")
var useDsa = flag.Bool("dsa", false, "DSA")
var useEcdsa = flag.Bool("ecdsa", false, "ECDSA")

var dsaL1k = flag.Bool("l1k", false, "DSA L1024N160 (default)")
var dsaL2k = flag.Bool("l2k", false, "DSA L2048N224")
var dsaL2kb = flag.Bool("l2kb", false, "DSA L2048N256")
var dsaL3k = flag.Bool("l3k", false, "DSA L3072N256")

//P256, P384 or P521
var ecdsaP256 = flag.Bool("p256", false, "ECDSA P256 (default)")
var ecdsaP384 = flag.Bool("p384", false, "ECDSA P384")
var ecdsaP521 = flag.Bool("p521", false, "ECDSA P521")

type kdsak struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Priv    *big.Int
	Pub     *big.Int
}

func genRsa() (e error) {
	var pk *rsa.PrivateKey
	if *primes>2 {
		pk,e = rsa.GenerateMultiPrimeKey(rand.Reader,*primes,*bits)
	} else {
		pk,e = rsa.GenerateKey(rand.Reader,*bits)
	}
	if e!=nil { return }
	b := x509.MarshalPKCS1PrivateKey(pk)
	e = pem.Encode(out,&pem.Block{"RSA PRIVATE KEY",nil,b})
	return
}

func genDsa() (e error) {
	var b []byte
	pk := new(dsa.PrivateKey)
	sz := dsa.L1024N160
	switch {
	case *dsaL2k:sz = dsa.L2048N224
	case *dsaL2kb:sz = dsa.L2048N256
	case *dsaL3k:sz = dsa.L3072N256
	}
	e = dsa.GenerateParameters(&pk.Parameters,rand.Reader,sz)
	if e!=nil { return }
	e = dsa.GenerateKey(pk,rand.Reader)
	if e!=nil { return }
	/* OpenSSL format
		ASN.1 SEQUENCE consisting of the values
		of version (currently zero), p, q, g, the public and private
		key components respectively as ASN.1 INTEGERs.
	*/
	k := kdsak{
		Version : 0,
		P : pk.P,
		Q : pk.Q,
		G : pk.G,
		Priv : pk.Y,
		Pub : pk.X,
	}
	b,e = asn1.Marshal(k)
	if e!=nil { return }
	e = pem.Encode(out,&pem.Block{"DSA PRIVATE KEY",nil,b})
	return
}

func genEcdsa() (e error) {
	var pk *ecdsa.PrivateKey
	var b []byte
	curve := elliptic.P256()
	switch {
	case *ecdsaP384:curve = elliptic.P384()
	case *ecdsaP521:curve = elliptic.P521()
	}
	pk,e = ecdsa.GenerateKey(curve,rand.Reader)
	if e!=nil { return }
	b,e = x509.MarshalECPrivateKey(pk)
	if e!=nil { return }
	e = pem.Encode(out,&pem.Block{"EC PRIVATE KEY",nil,b})
	return
}

func do_off() {
	out = os.Stdout
	os.Stdout = os.Stderr
}

func main() {
	do_off()
	var err error
	flag.Parse()
	switch {
	case *useRsa:
		err = genRsa()
		if err!=nil { fmt.Fprintln(os.Stderr,"error gen RSA:",err) }
	case *useDsa:
		err = genDsa()
		if err!=nil { fmt.Fprintln(os.Stderr,"error gen DSA:",err) }
	case *useEcdsa:
		err = genEcdsa()
		if err!=nil { fmt.Fprintln(os.Stderr,"error gen ECDSA:",err) }
	default: flag.PrintDefaults()
	}
}


