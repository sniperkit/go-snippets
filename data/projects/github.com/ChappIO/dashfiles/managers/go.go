package managers

import (
	"os/exec"
)

type Go struct {
	PackageManager
}

func init() {
	Register(&Go{})
}

func (gopm *Go) GetName() string {
	return "go"
}

func (gopm *Go) GetFileFormat() string {
	return "packages"
}

func (gopm *Go) IsInstalled() bool {
	_, err := exec.Command("go", "version").Output()
	return err == nil
}

func (gopm *Go) InstallFromFile(file string) error {
	packages := getToBeInstalledPackages(gopm, file)
	if len(packages) != 0 {
		return command("go", append([]string{"get", "-u"}, packages...)...).Run()
	}
	return nil
}

func (gopm *Go) Update() error {
	// We do not automatically update go packages due to performance
	return nil
}

func (gopm *Go) GetPackages() []string {
	// Go packages are reinstalled every time
	return nil
}
