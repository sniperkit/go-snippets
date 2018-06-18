package bagit

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// A Hasher represents a hash implementation for generating tag files
type Hasher struct {
	Hash func() hash.Hash
	Name string
}

// Convenience hasher variables for setting up tag-file hashing
var (
	HashMD5    = &Hasher{md5.New, "md5"}
	HashSHA1   = &Hasher{sha1.New, "sha1"}
	HashSHA256 = &Hasher{sha256.New, "sha256"}
	HashSHA512 = &Hasher{sha512.New, "sha512"}
)
