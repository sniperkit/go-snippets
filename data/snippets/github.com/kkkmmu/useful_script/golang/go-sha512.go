/*
SHA-2 (Secure Hash Algorithm 2) is a set of cryptographic hash functions designed by the United States National Security Agency (NSA).[3] Cryptographic hash functions are mathematical operations run on digital data; by comparing the computed "hash" (the output from execution of the algorithm) to a known and expected hash value, a person can determine the data's integrity. For example, computing the hash of a downloaded file and comparing the result to a previously published hash result can show whether the download has been modified or tampered with.[4] A key aspect of cryptographic hash functions is their collision resistance: nobody should be able to find two different input values that result in the same hash output.

SHA-2 includes significant changes from its predecessor, SHA-1. The SHA-2 family consists of six hash functions with digests (hash values) that are 224, 256, 384 or 512 bits: SHA-224, SHA-256, SHA-384, SHA-512, SHA-512/224, SHA-512/256.

SHA-256 and SHA-512 are novel hash functions computed with 32-bit and 64-bit words, respectively. They use different shift amounts and additive constants, but their structures are otherwise virtually identical, differing only in the number of rounds. SHA-224 and SHA-384 are simply truncated versions of the first two, computed with different initial values. SHA-512/224 and SHA-512/256 are also truncated versions of SHA-512, but the initial values are generated using the method described in Federal Information Processing Standards (FIPS) PUB 180-4. SHA-2 was published in 2001 by the National Institute of Standards and Technology (NIST) a U.S. federal standard (FIPS). The SHA-2 family of algorithms are patented in US patent 6829355.[5] The United States has released the patent under a royalty-free license.[6]

In 2005, an algorithm emerged for finding SHA-1 collisions in about 2,000 times fewer steps than was previously thought possible.[7] In 2017, an example of a SHA-1 collision was published.[8] The security margin left by SHA-1 is weaker than intended, and its use is therefore no longer recommended for applications that depend on collision resistance, such as digital signatures. Although SHA-2 bears some similarity to the SHA-1 algorithm, these attacks have not been successfully extended to SHA-2.

Currently, the best public attacks break preimage resistance for 52 rounds of SHA-256 or 57 rounds of SHA-512, and collision resistance for 46 rounds of SHA-256, as shown in the Cryptanalysis and validation section below.[1][2]
*/
package main

import (
	"crypto/sha512"
	"log"
)

func main() {
	hash := sha512.New()
	log.Println(hash.Sum([]byte("Hello World!")))
	log.Println(sha512.Sum384([]byte("Hello World!")))
	log.Println(sha512.Sum512([]byte("Hello World!")))
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
