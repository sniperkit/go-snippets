package fileutil

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

// CopyDirectory attempts to copy all files from srcPath to dstPath
// recursively.  dstPath must not exist.  Anything that isn't a file or a
// directory returns an error.  This includes symlinks for now.  The operation
// stops on the first error, and the partial copy is left in place.  Currently,
// permissions are not preserved.
func CopyDirectory(srcPath, dstPath string) error {
	var err error

	// Figure out absolute paths for clarity
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("source %q error: %s", srcPath, err)
	}
	dstPath, err = filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("destination %q error: %s", dstPath, err)
	}

	// Validate source exists and destination does not
	if !Exists(srcPath) {
		return fmt.Errorf("source %q does not exist", srcPath)
	}
	if !MustNotExist(dstPath) {
		return fmt.Errorf("destination %q already exists", dstPath)
	}

	// Destination parent must already exist
	if !IsDir(filepath.Dir(dstPath)) {
		return fmt.Errorf("destination's parent %q does not exist", dstPath)
	}

	// Get source path info and validate it's a directory
	var srcInfo os.FileInfo
	srcInfo, err = os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("source %q error: %s", srcPath, err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source %q is not a directory", srcPath)
	}

	return copyRecursive(srcPath, dstPath)
}

// copyRecursive is the actual file-copying function which CopyDirectory uses
func copyRecursive(srcPath, dstPath string) error {
	var err = os.MkdirAll(dstPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create directory %q: %s", dstPath, err)
	}

	var infos []os.FileInfo
	infos, err = Readdir(srcPath)
	if err != nil {
		return fmt.Errorf("unable to read source directory %q: %s", srcPath, err)
	}

	for _, info := range infos {
		var srcFull = filepath.Join(srcPath, info.Name())
		var dstFull = filepath.Join(dstPath, info.Name())

		var file = InfoToFile(info)
		switch {
		case file.IsDir():
			err = copyRecursive(srcFull, dstFull)
			if err != nil {
				return err
			}

		case file.IsRegular():
			err = CopyVerify(srcFull, dstFull)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unable to copy special file %q", srcFull)
		}
	}

	return nil
}

// CRC32 returns the checksum of the given file.  This is intended for
// verifying file copies immediately after the copy happens.  It should not be
// relied upon to detect malicious file changes.
func CRC32(file string) (string, error) {
	var f, err = os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var h = crc32.NewIEEE()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// CopyFile attempts to copy the bytes from src into dst, returning an error if
// applicable.  Does not use os.Link regardless of where the two files reside,
// as that can cause massive confusion when copying a file in order to back it
// up while writing out to the original.  The destination file permissions
// aren't set here, and must be managed externally.
func CopyFile(src, dst string) error {
	var err error
	var srcInfo os.FileInfo

	srcInfo, err = os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot stat %#v: %s", src, err)
	}
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("cannot copy non-regular file %#v: %s", src, err)
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot stat %#v: %s", dst, err)
	}

	return copyFileContents(src, dst)
}

// CopyVerify copies the bytes from src into dst using CopyFile, then verifies
// the two files have the same CRC32, giving a small measure of certainty that
// the copy succeeded.
func CopyVerify(src, dst string) error {
	var err = CopyFile(src, dst)
	if err != nil {
		return err
	}

	var srcChecksum, dstChecksum string
	srcChecksum, err = CRC32(src)
	if err != nil {
		return fmt.Errorf("unable to get source file's checksum: %s", err)
	}
	dstChecksum, err = CRC32(dst)
	if err != nil {
		return fmt.Errorf("unable to get destination file's checksum: %s", err)
	}
	if srcChecksum != dstChecksum {
		return fmt.Errorf("checksum failure")
	}

	return nil
}

// copyFileContents actually copies bytes from src to dst.  On any error, an
// attempt is made to clean up the state of the filesystem (though this is not
// guaranteed) and the first error encountered is returned.  i.e., if there's a
// failure in the io.Copy call, the caller will get that error, not the
// potentially meaningless error in the call to close the destination file.
func copyFileContents(src, dst string) error {
	var srcFile, dstFile *os.File
	var err error

	// Open source file or exit
	srcFile, err = os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to read %#v: %s", src, err)
	}
	defer srcFile.Close()

	// Create destination file or exit
	dstFile, err = os.Create(dst)
	if err != nil {
		return fmt.Errorf("unable to create %#v: %s", dst, err)
	}

	// Attempt to copy, and if the operation fails, attempt to clean up, then exit
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		err = fmt.Errorf("unable to copy data from %#v to %#v: %s", src, dst, err)
		dstFile.Close()
		os.Remove(dst)
		return err
	}

	// Attempt to sync the destination file
	err = dstFile.Sync()
	if err != nil {
		dstFile.Close()
		return fmt.Errorf("error syncing %#v: %s", dst, err)
	}

	// Attempt to close the destination file
	err = dstFile.Close()
	if err != nil {
		return fmt.Errorf("errro closing %#v: %s", dst, err)
	}

	return nil
}
