/*
In cryptography, SHA-1 (Secure Hash Algorithm 1) is a cryptographic hash function designed by the United States National Security Agency and is a U.S. Federal Information Processing Standard published by the United States NIST.[3] SHA-1 produces a 160-bit (20-byte) hash value known as a message digest. A SHA-1 hash value is typically rendered as a hexadecimal number, 40 digits long.

SHA-1 is no longer considered secure against well-funded opponents. In 2005, cryptanalysts found attacks on SHA-1 suggesting that the algorithm might not be secure enough for ongoing use,[4] and since 2010 many organizations have recommended its replacement by SHA-2 or SHA-3.[5][6][7] Microsoft,[8] Google,[9] Apple[10] and Mozilla[11][12][13] have all announced that their respective browsers will stop accepting SHA-1 SSL certificates by 2017.

On February 23, 2017 CWI Amsterdam and Google announced they had performed a collision attack against SHA-1,[14][15] publishing two dissimilar PDF files which produce the same SHA-1 hash as proof of concept.[16]
*/

package main

import (
	"crypto/sha1"
	"log"
)

func main() {
	hash := sha1.New()
	log.Printf("%x\n", hash.Sum([]byte("Hello World!")))
	log.Printf("%x\n", hash.Sum([]byte("How are you!")))
	log.Printf("%x\n", hash.Sum([]byte("Hello World!")))
	log.Printf("%x\n", hash.Sum([]byte("How are you!")))
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
