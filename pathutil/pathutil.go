// Additional helpers for working with file paths and filesystem information
package pathutil

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

func ExpandUser(path string) (string, error) {
	if u, err := user.Current(); err == nil {
		fullTilde := fmt.Sprintf("~%s", u.Name)

		if strings.HasPrefix(path, `~/`) || path == `~` {
			return strings.Replace(path, `~`, u.HomeDir, 1), nil
		}

		if strings.HasPrefix(path, fullTilde+`/`) || path == fullTilde {
			return strings.Replace(path, fullTilde, u.HomeDir, 1), nil
		}

		return path, nil
	} else {
		return path, err
	}
}

func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.Mode().IsRegular()
	}

	return false
}

func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}

	return false
}

func IsSymlink(mode os.FileMode) bool {
	return (mode&os.ModeSymlink != 0)
}

func IsDevice(mode os.FileMode) bool {
	return (mode&os.ModeDevice != 0)
}

func IsCharDevice(mode os.FileMode) bool {
	return (mode&os.ModeCharDevice != 0)
}

func IsNamedPipe(mode os.FileMode) bool {
	return (mode&os.ModeNamedPipe != 0)
}

func IsSocket(mode os.FileMode) bool {
	return (mode&os.ModeSocket != 0)
}

func IsSticky(mode os.FileMode) bool {
	return (mode&os.ModeSticky != 0)
}

func IsSetuid(mode os.FileMode) bool {
	return (mode&os.ModeSetuid != 0)
}

func IsSetgid(mode os.FileMode) bool {
	return (mode&os.ModeSetgid != 0)
}

func IsTemporary(mode os.FileMode) bool {
	return (mode&os.ModeTemporary != 0)
}

func IsExclusive(mode os.FileMode) bool {
	return (mode&os.ModeExclusive != 0)
}

func IsAppend(mode os.FileMode) bool {
	return (mode&os.ModeAppend != 0)
}

func IsReadable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

func IsWritable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}
