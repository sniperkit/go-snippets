package util

import (
	"go.polydawn.net/go-timeless-api"
)

/*
	Return a first, second, and remaining chunk of a ware's hash as strings.

	These are the first three, second three, and remaining bytes of the string.
	For base58 encoded values, these first two chunks used as dir prefixes are a
	cozy density for storing many many thousands of objects:

	If the hash is too short, we return a bunch of dashes.  (The hash is probably
	invalid semantically anyway, but we're not going to error about that here.)
	A hash of empty string will result in a return of `"---", "---", "-"` (in other
	words, as if the hash had been padded to a min of 7 characts, all dashes).
*/
func ChunkifyHash(wareID api.WareID) (string, string, string) {
	hash := wareID.Hash
	if len(hash) < 7 {
		hash = hash + "-------"[:7-len(hash)]
	}
	return hash[0:3], hash[3:6], hash[6:]
}
