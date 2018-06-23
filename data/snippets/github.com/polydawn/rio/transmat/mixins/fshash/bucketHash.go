package fshash

import (
	"fmt"
	"hash"
	"sort"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/tok"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/lib/treewalk"
)

/*
	Walks the tree of files and metadata arrayed in a `Bucket` and
	constructs a tree hash over them.  The root of the tree hash is returned.
	The returned root has can be said to verify the integrity of the
	entire tree (much like a Merkle tree).

	The serial structure is expressed something like the following:

		{"node": $dir.metadata.hash,
		 "leaves": [
			{"node": $file1.metadata.hash, "content": $file1.contentHash},
			{"node": $subdir.metadata.hash,
			 "leaves": [ ... ]},
		 ]
		}

	This expression is made in cbor (rfc7049) format with indefinite-length
	arrays and a fixed order for all map fields.  Every structure starting
	with "node" is itself hashed and that value substituted in before
	hashing the parent.  Since the metadata hash contains the file/dir name,
	and the tree itself is traversed in sorted order, the entire structure
	is computed deterministically and unambiguously.
*/
func HashBucket(bucket Bucket, hasherFactory func() hash.Hash) []byte {
	// At every point in the visitation, children need to submit their hashes back up the tree.
	// Prime the pump with a special reaction for when the root returns; every directory preVisit attaches hoppers for children thereon.
	upsubs := make(upsubStack, 0)
	var finalAnswer []byte
	upsubs.Push(func(x []byte) {
		finalAnswer = x
	})
	// Also keep a stack of hashers in use because they jump across the pre/post visit gap.
	hashers := make(hasherStack, 0)
	// Keep a count of how many nodes visited in total.  Cheap sanity check.
	var visitCount int

	// Visitor definitions
	preVisit := func(node treewalk.Node) error {
		record := node.(RecordIterator).Record()
		visitCount++
		hasher := hasherFactory()
		enc := cbor.NewEncoder(hasher)

		// Begin map for the entry.
		//  Length two for dirs and files: it's metadata + one of either leaves list or contenthash.
		//  Length one for everything else: their attributes are all in the metadata.
		switch record.Metadata.Type {
		case fs.Type_Dir, fs.Type_File:
			enc.Step(&tok.Token{Type: tok.TMapOpen, Length: 2})
		default:
			enc.Step(&tok.Token{Type: tok.TMapOpen, Length: 1})
		}

		// Encode the metadata.
		enc.Step(&tok.Token{Type: tok.TString, Str: "m"})
		marshalMetadata(enc, record.Metadata)

		// Switch for leaves list for dirs, or pure content hash for files.
		//  (All the other kinds of file are fully described by their metadata alone.)
		switch record.Metadata.Type {
		case fs.Type_Dir:
			// Open the "leaves" array.
			//  This may end up being an empty dir, but we act the same regardless
			//  (and we don't have that information here since the iterator has tunnel vision).
			//  The array will eventually be closed in the postVisit hook.
			enc.Step(&tok.Token{Type: tok.TString, Str: "l"})
			enc.Step(&tok.Token{Type: tok.TArrOpen, Length: -1})
			upsubs.Push(func(x []byte) {
				enc.Step(&tok.Token{Type: tok.TBytes, Bytes: x})
			})
			hashers.Push(hasher)
		case fs.Type_File:
			// heap the object's content hash in
			enc.Step(&tok.Token{Type: tok.TString, Str: "h"})
			enc.Step(&tok.Token{Type: tok.TBytes, Bytes: record.ContentHash})
			// finalize our hash here and upsub to save us the work of hanging onto the hasher until the postvisit call
			upsubs.Peek()(hasher.Sum(nil))
		}
		return nil
	}
	postVisit := func(node treewalk.Node) error {
		record := node.(RecordIterator).Record()
		switch record.Metadata.Type {
		case fs.Type_Dir:
			hasher := hashers.Pop()
			// Close off the "leaves" array.
			//  No map-close necessary because we used a fixed length map.
			//  This is a weird reach-around using a hardcoded cbor constant because we didn't hang on to the encoder.
			hasher.Write([]byte{0xff})
			hash := hasher.Sum(nil)
			// pop out this dir's hoppers for children data
			upsubs.Pop()
			// hash and upsub
			upsubs.Peek()(hash)
		default:
		}
		return nil
	}
	// Traverse
	if err := treewalk.Walk(bucket.Iterator(), preVisit, postVisit); err != nil {
		panic(err) // none of our code has known believable error returns.
	}
	// Sanity check no node left behind
	_ = upsubs.Pop()
	if !upsubs.Empty() || !hashers.Empty() {
		panic(fmt.Errorf("invariant failed after bucket records walk: stacks not empty"))
	}
	if visitCount != bucket.Length() {
		panic(fmt.Errorf("invariant failed after bucket records walk: visited %d of %d nodes", visitCount, bucket.Length()))
	}
	// return the result upsubbed by the root
	return finalAnswer
}

type upsubStack []func([]byte)

func (s upsubStack) Empty() bool          { return len(s) == 0 }
func (s upsubStack) Peek() func([]byte)   { return s[len(s)-1] }
func (s *upsubStack) Push(x func([]byte)) { (*s) = append((*s), x) }
func (s *upsubStack) Pop() func([]byte) {
	x := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return x
}

// look me in the eye and tell me again how generics are a bad idea
type hasherStack []hash.Hash

func (s hasherStack) Empty() bool       { return len(s) == 0 }
func (s hasherStack) Peek() hash.Hash   { return s[len(s)-1] }
func (s *hasherStack) Push(x hash.Hash) { (*s) = append((*s), x) }
func (s *hasherStack) Pop() hash.Hash {
	x := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return x
}

// Marshal an `fs.Metadata` into the given encoder.
// This manual marshalling implementation has a stable order and the correctness
// of the HashBucket method over time relies on this.
func marshalMetadata(enc *cbor.Encoder, m fs.Metadata) {
	// Count up how many fields we're about to encode.
	fieldCount := 7
	if m.Linkname != "" {
		fieldCount++
	}
	xattrsLen := len(m.Xattrs)
	if xattrsLen > 0 {
		fieldCount++
	}
	if m.Type == fs.Type_Device || m.Type == fs.Type_CharDevice {
		fieldCount += 2 // devmajor and devminor will be included for these types
	}
	// Let us begin!
	enc.Step(&tok.Token{Type: tok.TMapOpen, Length: fieldCount})
	// Name
	enc.Step(&tok.Token{Type: tok.TString, Str: "n"})
	enc.Step(&tok.Token{Type: tok.TString, Str: m.Name.Last()}) // uses basename so hash subtrees are severable
	// Type
	enc.Step(&tok.Token{Type: tok.TString, Str: "t"})
	enc.Step(&tok.Token{Type: tok.TString, Str: string(m.Type)})
	// Permission mode bits (this is presumed to already be basic perms (0777) and setuid/setgid/sticky (07000) only, per fs.Metadata standard).
	enc.Step(&tok.Token{Type: tok.TString, Str: "p"})
	enc.Step(&tok.Token{Type: tok.TInt, Int: int64(m.Perms)})
	// UID (numeric)
	enc.Step(&tok.Token{Type: tok.TString, Str: "u"})
	enc.Step(&tok.Token{Type: tok.TInt, Int: int64(m.Uid)})
	// GID (numeric)
	enc.Step(&tok.Token{Type: tok.TString, Str: "g"})
	enc.Step(&tok.Token{Type: tok.TInt, Int: int64(m.Gid)})
	// Skipped: size -- because that's fairly redundant (and we never use hashes that are subject to length extension)
	// Linkname, if it's a symlink
	if m.Linkname != "" {
		enc.Step(&tok.Token{Type: tok.TString, Str: "l"})
		enc.Step(&tok.Token{Type: tok.TString, Str: m.Linkname})
	}
	// devMajor and devMinor numbers, if it's a device
	if m.Type == fs.Type_Device || m.Type == fs.Type_CharDevice {
		enc.Step(&tok.Token{Type: tok.TString, Str: "dM"})
		enc.Step(&tok.Token{Type: tok.TInt, Int: m.Devmajor})
		enc.Step(&tok.Token{Type: tok.TString, Str: "dm"})
		enc.Step(&tok.Token{Type: tok.TInt, Int: m.Devminor})
	}
	// Modtime
	enc.Step(&tok.Token{Type: tok.TString, Str: "m"})
	enc.Step(&tok.Token{Type: tok.TInt, Int: m.Mtime.Unix()})
	enc.Step(&tok.Token{Type: tok.TString, Str: "mn"})
	enc.Step(&tok.Token{Type: tok.TInt, Int: int64(m.Mtime.Nanosecond())})
	// Xattrs are a mite more complicated because we have to handle unknown keys:
	if xattrsLen > 0 {
		enc.Step(&tok.Token{Type: tok.TString, Str: "x"})
		sorted := make([]stringPair, 0, xattrsLen)
		for k, v := range m.Xattrs {
			sorted = append(sorted, stringPair{k, v})
		}
		sort.Sort(sortableStringPair(sorted))
		enc.Step(&tok.Token{Type: tok.TMapOpen, Length: xattrsLen})
		for _, line := range sorted {
			enc.Step(&tok.Token{Type: tok.TString, Str: line.a})
			enc.Step(&tok.Token{Type: tok.TString, Str: line.b})
		}
	}
	// There is no map-end to encode in cbor since we used the fixed-length map.  We're done.
}

type stringPair struct{ a, b string }
type sortableStringPair []stringPair

func (p sortableStringPair) Len() int           { return len(p) }
func (p sortableStringPair) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p sortableStringPair) Less(i, j int) bool { return p[i].a < p[j].a }
