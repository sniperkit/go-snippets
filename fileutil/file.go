package fileutil

import (
	"os"
	"time"
)

// A File represents anything on the filesystem, including directories,
// symlinks, pipes, etc.  It's a simple, concrete encapsulation of FileInfo
// minus the Sys() data, since that's always of unknown type.
type File struct {
	Name    string      // base name of the file
	Size    int64       // length in bytes for regular files; system-dependent for others
	Mode    os.FileMode // file mode bits
	ModTime time.Time   // modification time
}

// InfoToFile converts a FileInfo interface into File data
func InfoToFile(fi os.FileInfo) *File {
	return &File{Name: fi.Name(), Size: fi.Size(), Mode: fi.Mode(), ModTime: fi.ModTime()}
}

// InfosToFiles converts a slice of FileInfos into a slice of Files
func InfosToFiles(fiList []os.FileInfo) []*File {
	var files = make([]*File, len(fiList))
	for i, fi := range fiList {
		files[i] = InfoToFile(fi)
	}

	return files
}

// IsRegular returns true if the file isn't dir/symlink/etc
func (f *File) IsRegular() bool {
	return f.Mode&os.ModeType == 0
}

// IsDir returns true if the file is a directory
func (f *File) IsDir() bool {
	return f.Mode&os.ModeDir != 0
}

// IsSymlink reports if the file is a symlink
func (f *File) IsSymlink() bool {
	return f.Mode&os.ModeSymlink != 0
}

// IsNamedPipe reports if the file is a named pipe
func (f *File) IsNamedPipe() bool {
	return f.Mode&os.ModeNamedPipe != 0
}

// IsSocket reports if the file is a socket
func (f *File) IsSocket() bool {
	return f.Mode&os.ModeSocket != 0
}

// IsDevice reports if the file is a device
func (f *File) IsDevice() bool {
	return f.Mode&os.ModeDevice != 0
}
