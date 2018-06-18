package bagit

import (
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uoregon-libraries/gopkg/fileutil"
)

// FileChecksum holds a path to a file and its checksum
type FileChecksum struct {
	Path     string
	Checksum string
}

// Bag holds state for the generation of bag manifest and other tag files
type Bag struct {
	root      string
	Hasher    *Hasher
	Checksums []*FileChecksum
}

// New returns Bag structure for processing the given root path, and sets the
// hasher to the built-in SHA256
func New(root string) *Bag {
	return &Bag{root: root, Hasher: HashSHA256}
}

// WriteTagFiles traverses all files under the bag's root/data, generates
// hashes for each, and writes out "manifest-[hashtype].txt".  Upon completion,
// bagit.txt and tagmanifest-[hashtype].txt are then written.
//
// This is not parallelized as it seems unlikely any advantage would be gained
// since file IO is likely to be the main cost, not CPU.
func (b *Bag) WriteTagFiles() (err error) {
	err = b.GenerateChecksums()
	if err == nil {
		err = b.writeManifest()
	}
	if err == nil {
		err = b.writeBagitFile()
	}
	if err == nil {
		err = b.writeTagManifest()
	}

	return
}

// GenerateChecksums iterates over all files in the data path and generates
// each file's checksum in turn, returning the resulting slice of
// FileChecksums.  The checksum path is always relative to the bag's root.
//
// If there are any errors, GenerateChecksums returns an empty slice along with
// relevant error information.
//
// This is typically used internally to generate the manifest file, but can be
// useful for testing, bag validation, or making use of the BagIt data
// structure in cases where checksums need to be stored externally to the data.
func (b *Bag) GenerateChecksums() error {
	var err error
	var realroot string

	realroot, err = filepath.Abs(b.root)
	if err != nil {
		return fmt.Errorf("unable to determine bag's absolute root path from %q: %s", b.root, err)
	}
	b.root = realroot

	var dataPath = filepath.Join(b.root, "data")
	if !fileutil.IsDir(dataPath) {
		return fmt.Errorf("%q is not a directory", dataPath)
	}

	b.Checksums = nil
	err = filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			var chksum, err = b.getsum(path)
			if err == nil {
				b.Checksums = append(b.Checksums, chksum)
			}
			return err
		}

		return nil
	})

	return err
}

func (b *Bag) getsum(path string) (*FileChecksum, error) {
	var err error
	var f *os.File
	var relPath, hexSum string
	var h hash.Hash

	f, err = os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %q: %s", path, err)
	}
	defer f.Close()

	h = b.Hasher.Hash()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q for hashing: %s", path, err)
	}

	relPath, err = filepath.Rel(b.root, path)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q's relative file path: %s", path, err)
	}

	hexSum = fmt.Sprintf("%x", h.Sum(nil))
	return &FileChecksum{Path: relPath, Checksum: hexSum}, nil
}

func (b *Bag) writeManifest() error {
	var manifestFile = filepath.Join(b.root, "manifest-"+b.Hasher.Name+".txt")
	if !fileutil.MustNotExist(manifestFile) {
		return fmt.Errorf("manifest file %q must not exist", manifestFile)
	}

	var f = fileutil.NewSafeFile(manifestFile)
	for _, ck := range b.Checksums {
		fmt.Fprintf(f, "%s  %s\n", ck.Checksum, ck.Path)
	}

	var err = f.Close()
	if err != nil {
		return fmt.Errorf("error writing manifest file: %s", err)
	}

	return nil
}

func (b *Bag) writeBagitFile() error {
	var f = fileutil.NewSafeFile(filepath.Join(b.root, "bagit.txt"))
	f.Write([]byte("BagIt-Version: 0.97\nTag-File-Character-Encoding: UTF-8\n"))
	return f.Close()
}

func (b *Bag) writeTagManifest() error {
	// Technically all files at the bag root are considered tag files, as a bag
	// can have custom tag files, so we iterate all files at that top level
	var infos, err = ioutil.ReadDir(b.root)
	if err != nil {
		return fmt.Errorf("error reading bag root: %s", err)
	}

	var f = fileutil.NewSafeFile(filepath.Join(b.root, "tagmanifest-"+b.Hasher.Name+".txt"))
	for _, info := range infos {
		if info.Mode().IsRegular() {
			var path = filepath.Join(b.root, info.Name())
			var chksum, err = b.getsum(path)
			if err != nil {
				f.Cancel()
				return fmt.Errorf("error getting %q's checksum: %s", path, err)
			}

			_, err = fmt.Fprintf(f, "%s  %s\n", chksum.Checksum, info.Name())
			if err != nil {
				f.Cancel()
				return fmt.Errorf("error writing checksum: %s", err)
			}
		}
	}

	f.Close()
	return nil
}
