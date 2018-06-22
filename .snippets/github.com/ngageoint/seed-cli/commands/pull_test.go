package commands

import (
	"fmt"
	"os"
	"testing"
	"time"

	"strings"

	common_const "github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestDockerPull(t *testing.T) {
	util.RestartRegistry()

	registry := "localhost:5000"
	username := "testuser"
	password := "testpassword"

	//set config dir so we don't stomp on other users' logins with sudo
	configDir := common_const.DockerConfigDir + time.Now().Format(time.RFC3339)
	os.Setenv(common_const.DockerConfigKey, configDir)
	defer util.RemoveAllFiles(configDir)
	defer os.Unsetenv(common_const.DockerConfigKey)

	err := util.Login(registry, username, password)
	if err != nil {
		fmt.Println(err)
	}

	imgDirs := []string{"../testdata/complete/"}
	origImg := "my-job-0.1.0-seed:0.1.0"
	remoteImg := []string{"localhost:5000/my-job-0.1.0-seed:0.1.0", "localhost:5000/my-job-1.0.0-seed:1.0.0", "localhost:5000/not-a-valid-image"}
	version := "1.0.0"

	for _, dir := range imgDirs {
		_, err := DockerBuild(dir, version, "", "", ".", ".", "")
		if err != nil {
			t.Errorf("Error building image from %v for DockerPull test: %v", dir, err)
		}
	}

	for _, img := range remoteImg {
		err := util.Tag(origImg, img)
		if err != nil {
			t.Errorf("Error tagging image %v for DockerPull test: %v", img, err)
		}

		err = util.Push(img)
		if err != nil {
			t.Errorf("Error pushing image %v for DockerPull test: %v", img, err)
		}
	}

	cases := []struct {
		image            string
		registry         string
		org              string
		username         string
		password         string
		expectedResult   bool
		expectedErrorMsg string
	}{
		{"my-job-1.0.0-seed:1.0.0", "localhost:5000", "", "testuser", "wrongpassword",
			false, "401 Unauthorized"},
		{"not-a-valid-image", "localhost:5000", "", "testuser", "testpassword",
			true, ""},
	}

	for _, c := range cases {
		err := DockerPull(c.image, c.registry, c.org, c.username, c.password)

		success := err == nil
		if success != c.expectedResult {
			t.Errorf("DockerPull returned %v, expected %v\n", success, c.expectedResult)
		}

		if err != nil {
			errMsg := err.Error()
			if !strings.Contains(errMsg, c.expectedErrorMsg) {
				t.Errorf("DockerPull returned error %v, expected %v\n", errMsg, c.expectedErrorMsg)
			}
		}
	}
}
