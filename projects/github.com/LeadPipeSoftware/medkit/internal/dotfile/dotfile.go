package dotfile

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const bundlesDir = "bundles"

var AlwaysSkip bool
var AlwaysOverwrite bool
var HomeDirectory string

// InstallDotfiles will install dotfiles.
func InstallDotfiles(dotfilesDirectory string, homeDirectory string, alwaysSkip bool, alwaysOverwrite bool, backupExtension string) {
	AlwaysSkip = alwaysSkip
	AlwaysOverwrite = alwaysOverwrite

	if AlwaysSkip && AlwaysOverwrite {
		fmt.Print("Sorry, the always-skip and always-overwrite flags cannot be used together.")
		os.Exit(1)
	}

	fmt.Printf("Dotfile directory: %s", dotfilesDirectory)
	fmt.Printf("\n   Home directory: %s", homeDirectory)
	fmt.Println()

	if err := filepath.Walk(dotfilesDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == bundlesDir {
			return filepath.SkipDir
		}

		if info.IsDir() {
			if err := symlinkFilesInDirectory(path, homeDirectory, backupExtension); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println()
}

// ShowDotfiles will display all dotfiles (recursively) in the dotfiles directory.
func ShowDotfiles(dotfilesDirectory string) {
	if dotfiles, err := getAllDotfiles(dotfilesDirectory); err == nil {
		for _, dotfile := range dotfiles {
			fmt.Printf("\n%s", dotfile)
		}
	} else {
		fmt.Printf("\n%s", err)
	}

	fmt.Println()
}

// getAllDotfiles returns the paths of every dotfile in a directory (recursively).
func getAllDotfiles(rootpath string) ([]string, error) {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".symlink" {
			list = append(list, path)
		}

		return nil
	})

	return list, err
}

// backupThenRemoveFile creates a backup of a file, but removes existing backups first.
func backupThenRemoveFile(filename string, backupExtension string) error {
	backupFile := filename + backupExtension

	if fileExists(backupFile) {
		if err := os.Remove(backupFile); err != nil {
			fmt.Printf("\nERROR removing backup %s: %s", backupFile, err)
			return err
		}
	}

	if err := os.Rename(filename, backupFile); err != nil {
		fmt.Printf("\nERROR creating backup %s: %s", backupFile, err)
		return err
	}

	return nil
}

// createSymlink creates a symbolic link.
func createSymlink(source string, target string) {
	if err := os.Symlink(source, target); err == nil {
		fmt.Printf("\nSymlinked %s => %s", source, target)
	} else {
		fmt.Printf("\nERROR symlinking %s: %s", target, err)
	}
}

// fileDoesNotExist returns true if the specified file does not exist.
func fileDoesNotExist(targetFile string) bool {
	_, err := os.Stat(targetFile)

	return os.IsNotExist(err)
}

// fileExists returns true if the specified file exists.
func fileExists(targetFile string) bool {
	return !fileDoesNotExist(targetFile)
}

// getSymlinkTargetName builds and returns the symbolic link target name (the name without the .symlink extension).
func getSymlinkTargetName(fileName string) string {
	name := path.Base(fileName)
	re := regexp.MustCompile("\\.symlink")

	return re.ReplaceAllString(name, "")
}

// symlinkFilesInDirectory creates symbolic links for each .symlink file in a directory.
func symlinkFilesInDirectory(path string, home string, backupExtension string) error {
	matches, err := filepath.Glob(path + "/*.symlink")

	if err == nil {
		for _, match := range matches {
			targetFile := home + "/" + getSymlinkTargetName(match)

			if fileDoesNotExist(targetFile) {
				createSymlink(match, targetFile)
			} else {
				input := bufio.NewScanner(os.Stdin)

				if AlwaysSkip {
					fmt.Printf("\nSkipping %s", targetFile)
				} else if AlwaysOverwrite {
					fmt.Printf("\nOverwriting %s", targetFile)
					if err := backupThenRemoveFile(targetFile, backupExtension); err == nil {
						createSymlink(match, targetFile)
					}
				} else {
					fmt.Printf("\n\n%s already exists. What do you want to do?\n[s]kip, [S]kip All, [o]verwrite, [O]verwrite All: ", targetFile)

				InputLoop:
					for input.Scan() {

						answer := input.Text()

						switch answer {
						case "O", "o":
							fmt.Printf("\nOkay, overwriting %s", targetFile)
							if err := backupThenRemoveFile(targetFile, backupExtension); err == nil {
								createSymlink(match, targetFile)
							}
							if answer == "O" {
								AlwaysOverwrite = true
							}
							break InputLoop
						case "S", "s":
							fmt.Printf("\nOkay, skipping %s", targetFile)
							if answer == "S" {
								AlwaysSkip = true
							}
							break InputLoop
						default:
							fmt.Printf("\n%s already exists. What do you want to do?\n[s]kip, [S]kip All, [o]verwrite, [O]verwrite All: ", targetFile)
						}
					}
				}
			}
		}
	} else {
		return err
	}

	return nil
}
