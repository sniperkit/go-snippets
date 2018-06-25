package managers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var packageManagers []PackageManager
var commentRegex, _ = regexp.Compile("#.*")

func Register(manager PackageManager) {
	packageManagers = append(packageManagers, manager)
}

func GetManagers() []PackageManager {
	return packageManagers
}

type PackageManager interface {
	GetName() string
	IsInstalled() bool
	InstallFromFile(filePath string) error
	GetPackages() []string
	Update() error
	GetFileFormat() string
}

func getToBeInstalledPackages(manager PackageManager, filePath string) []string {
	if file, err := os.Open(filePath); err != nil {
		fmt.Println("Could not open " + filePath + ": " + err.Error())
		return nil
	} else {
		scanner := bufio.NewScanner(file)
		var packageData []string

		// Collect all packages that should be installed
		for scanner.Scan() {
			packageLine := scanner.Bytes()
			packageLine = commentRegex.ReplaceAll(packageLine, nil)

			packageName := strings.TrimSpace(string(packageLine))
			if packageName != "" {
				packageData = append(packageData, packageName)
			}
		}

		// Remove already installed packages
		removedCount := 0
		for _, installedPackage := range manager.GetPackages() {
			for i, matchPackage := range packageData {
				if installedPackage == matchPackage {
					packageData[i] = ""
					removedCount++
				}
			}
		}

		packages := make([]string, len(packageData)-removedCount)

		index := 0
		for _, element := range packageData {
			if element != "" {
				packages[index] = element
				index++
			}
		}

		return packages
	}
}
