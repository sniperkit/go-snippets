package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// SafeFile wraps the process of creating a temporary file, writing to it,
// copying it to the final file once all writing is successful, and then
// deleting the temp file.  This reduces the opportunities for errors
// (file-writing or program logic) to leave behind remnants of files.
//
// SafeFile implements io.WriteCloser for general-case usage, but also exposes
// a Cancel() method so callers can state that something failed and the file
// shouldn't be copied to its final location.
//
// Errors are returned from WriteCloser methods, but are also stored internally
// to allow the SafeFile to automatically skip certain methods and clean up
// after itself for simpler use-cases.  This means that when errors occur which
// aren't related to the caller's logic (e.g., an error is returned by Write),
// calling Cancel is unnecessary.
type SafeFile struct {
	temp      io.WriteCloser
	tempName  string
	finalPath string
	err       error
	closed    bool
}

// NewSafeFile returns a new SafeFile construct, wrapping the given path
func NewSafeFile(path string) *SafeFile {
	var f, err = ioutil.TempFile("", "")
	if err != nil {
		err = fmt.Errorf("couldn't generate temp file: %s", err)
	}
	return &SafeFile{temp: f, tempName: f.Name(), finalPath: path, err: err}
}

// Write delegates to the temporary file handle
func (f *SafeFile) Write(p []byte) (n int, err error) {
	if f.err != nil {
		return 0, fmt.Errorf("cannot write to SafeFile with errors")
	}

	n, f.err = f.temp.Write(p)
	return n, f.err
}

// Close closes the temporary file, copies its data to the final location, and
// removes the temp file from disk.  If an error occurs, an attempt is made to
// remove all files from disk to avoid broken files, and the error is returned.
func (f *SafeFile) Close() error {
	// Closing a file twice isn't great, but it should simply return the error
	// from f.temp.Close, not cancel the otherwise successful operation
	if f.closed {
		return f.temp.Close()
	}

	if f.err != nil {
		f.Cancel()
		return f.err
	}

	var err = f.temp.Close()
	if err != nil {
		f.Cancel()
		f.err = fmt.Errorf("unable to close temp file: %s", f.err)
		return f.err
	}

	err = CopyVerify(f.tempName, f.finalPath)
	if err != nil {
		f.Cancel()
		f.err = fmt.Errorf("unable to copy temp file: %s", f.err)
		return f.err
	}

	f.closed = true
	return os.Remove(f.tempName)
}

// Cancel attempts to close and delete all files
func (f *SafeFile) Cancel() {
	f.temp.Close()
	os.Remove(f.tempName)
	os.Remove(f.finalPath)
}
