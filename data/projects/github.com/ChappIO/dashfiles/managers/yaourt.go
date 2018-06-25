package managers

import (
	"os"
	"os/exec"
	"strings"
)

type Yaourt struct {
	PackageManager
}

const YAOURT = "yaourt"

func init() {
	Register(&Yaourt{})
}

func (yaourt *Yaourt) GetName() string {
	return YAOURT
}

func (yaourt *Yaourt) GetFileFormat() string {
	return "packages"
}

func (yaourt *Yaourt) IsInstalled() bool {
	_, err := exec.Command(YAOURT, "--version").Output()
	return err == nil
}

func (yaourt *Yaourt) InstallFromFile(file string) error {
	packages := getToBeInstalledPackages(yaourt, file)
	if len(packages) != 0 {
		args := append([]string{"-S"}, packages...)
		return command(YAOURT, args...).Run()
	}
	return nil
}

func (yaourt *Yaourt) Update() error {
	return command(YAOURT, "-Syu").Run()
}

func (yaourt *Yaourt) GetPackages() []string {
	cmd := exec.Command(YAOURT, "-Qe")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if data, err := cmd.Output(); err != nil {
		return []string{}
	} else {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")

		for i, line := range lines {
			// Format Line: <repo>/<package-name> <version> (<tags>)
			// We only need <package-name>
			lines[i] = line[strings.Index(line, "/")+1 : strings.Index(line, " ")]
		}
		return lines
	}
}
