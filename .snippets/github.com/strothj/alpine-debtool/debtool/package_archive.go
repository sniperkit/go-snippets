package debtool

import "net/url"

// ArchiveType specifies that the archive is for binary or source packages.
// It represents the "deb" or "deb-src" in a sources.list line.
type ArchiveType int

const (
	// BinaryArchive is a repository which provides binary packages.
	// It is the "deb" in a sources.list line.
	BinaryArchive ArchiveType = iota
	// SourceArchive is a repository which provides source packages.
	// It is the "deb-src" in a sources.list line.
	SourceArchive
)

// PackageArchive represents a Debian package archive. It represents a line in
// the typical "sources.list".
//
// deb http://ftp.debian.org/debian squeeze main contrib non-free
// could be represented as:
//	pa := debtool.PackageArchive{
//		Type:         debtool.BinaryArchive,
//		Root:         &url.URL{Path: "http://ftp.debian.org/debian"},
//		Distribution: "squeeze",
//		Components:   []string{"main", "contrib", "non-free"},
//	}
type PackageArchive struct {
	Type         ArchiveType
	Root         *url.URL
	Distribution string
	Components   []string
}
