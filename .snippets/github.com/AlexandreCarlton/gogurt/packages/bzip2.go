// Should really learn what 'package' does
package packages

import (
	"fmt"

	"github.com/alexandrecarlton/gogurt"
)

type Bzip2 struct {}

func (bzip2 Bzip2) Name() string {
	return "bzip2"
}

func (bzip2 Bzip2) URL(version string) string {
	return fmt.Sprintf("http://www.bzip.org/%s/bzip2-%s.tar.gz", version, version)
}


func (bzip2 Bzip2) Build(config gogurt.Config) error {
	cmd := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	fmt.Println(cmd)
	return cmd.Run()
}

func (bzip2 Bzip2) Install(config gogurt.Config) error {
	fmt.Println("Running bzip install...")
	cmd := gogurt.MakeCmd{
		Args: []string{
			"install",
			"PREFIX=" + config.InstallDir(bzip2),
		},
	}.Cmd()
	fmt.Println(cmd)
	return cmd.Run()
}

func (bzip2 Bzip2) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
