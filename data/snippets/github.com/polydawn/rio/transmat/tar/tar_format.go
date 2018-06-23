package tartrans

import (
	"archive/tar"
	"fmt"

	. "github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
)

// Mutate tar.Header fields to match the given fmeta.
func MetadataToTarHdr(fmeta *fs.Metadata, hdr *tar.Header) {
	hdr.Name = fmeta.Name.String()
	if fmeta.Type == fs.Type_Dir {
		hdr.Name += "/"
	}
	hdr.Typeflag = fsTypeToTarType(fmeta.Type)
	hdr.Mode = int64(fmeta.Perms)
	hdr.Uid = int(fmeta.Uid)
	hdr.Gid = int(fmeta.Gid)
	hdr.Size = fmeta.Size
	hdr.Linkname = fmeta.Linkname
	hdr.Devmajor = fmeta.Devmajor
	hdr.Devminor = fmeta.Devminor
	hdr.ModTime = fmeta.Mtime
	hdr.Xattrs = fmeta.Xattrs
}

func fsTypeToTarType(fsType fs.Type) byte {
	switch fsType {
	case fs.Type_File:
		return tar.TypeReg
	case fs.Type_Hardlink:
		return tar.TypeLink
	case fs.Type_Symlink:
		return tar.TypeSymlink
	case fs.Type_CharDevice:
		return tar.TypeChar
	case fs.Type_Device:
		return tar.TypeBlock
	case fs.Type_Dir:
		return tar.TypeDir
	case fs.Type_NamedPipe:
		return tar.TypeFifo
	case fs.Type_Socket:
		panic(fmt.Errorf("can't pack sockets into tar"))
	default:
		panic(fmt.Errorf("invalid fs.Type %q", fsType))

	}
}

// Mutate fs.Metadata fields to match the given tar header.
// Does not check for names that go above '.'; caller may want to do that.
func TarHdrToMetadata(hdr *tar.Header, fmeta *fs.Metadata) error {
	fmeta.Name = fs.MustRelPath(hdr.Name) // FIXME should not use the 'must' path
	fmeta.Type = tarTypeToFsType(hdr.Typeflag)
	if fmeta.Type == fs.Type_Invalid {
		return Errorf(rio.ErrWareCorrupt, "corrupt tar: %q is not a known file type", hdr.Typeflag)
	}
	fmeta.Perms = fs.Perms(hdr.Mode & 07777)
	fmeta.Uid = uint32(hdr.Uid)
	fmeta.Gid = uint32(hdr.Gid)
	fmeta.Size = hdr.Size
	fmeta.Linkname = hdr.Linkname
	fmeta.Devmajor = hdr.Devmajor
	fmeta.Devminor = hdr.Devminor
	fmeta.Mtime = hdr.ModTime
	fmeta.Xattrs = hdr.Xattrs
	return nil
}

func tarTypeToFsType(tarType byte) fs.Type {
	switch tarType {
	case tar.TypeReg, tar.TypeRegA:
		return fs.Type_File
	case tar.TypeLink:
		return fs.Type_Hardlink
	case tar.TypeSymlink:
		return fs.Type_Symlink
	case tar.TypeChar:
		return fs.Type_CharDevice
	case tar.TypeBlock:
		return fs.Type_Device
	case tar.TypeDir:
		return fs.Type_Dir
	case tar.TypeFifo:
		return fs.Type_NamedPipe
	// Notice that tar does not have a type for socket files
	default:
		return fs.Type_Invalid
	}
}
