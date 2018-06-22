package commands

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestDockerList(t *testing.T) {
	cases := []struct {
		directory        string
		imageName        string
		expectedErrorMsg string
	}{
		{"../testdata/dummy-scratch/", "test-seed", ""},
	}

	for _, c := range cases {
		buildArgs := []string{"build", "-t", c.imageName, c.directory}
		cmd := exec.Command("docker", buildArgs...)
		cmd.Run()
		output, err := DockerList()
		if err != nil {
			t.Errorf("DockerList returned an error: %v", err)
		}
		if !strings.Contains(output, c.imageName) {
			t.Errorf("DockerList() did not return expected image %v", c.imageName)
		}
	}
}
