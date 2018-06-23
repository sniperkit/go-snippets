package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestDockerBuild(t *testing.T) {
	cases := []struct {
		directory        string
		version          string
		manifest         string
		dockerfile       string
		expected         bool
		expectedErrorMsg string
	}{
		/*0*/ {"../examples/addition-job/", "1.0.0", ".", ".", true, ""},
		/*1*/ {"../examples/extractor/", "1.0.0", ".", ".", true, ""},
		/*2*/ {"../examples/extractor/", "1.0.0", "../examples/addition-job/seed.manifest.json", ".", true, ""},
		/*3*/ {"../examples/addition-job/", "1.0.0", ".", "../examples/extractor", true, ""},
		/*4*/ {"../examples/extractor/", "1.0.0", "../examples/addition-job/seed.manifest.json", "../examples/addition-job", true, ""},
		/*5*/ {"", "", ".", ".", false, "seed.manifest.json cannot be found"},
	}

	for _, c := range cases {
		_, err := DockerBuild(c.directory, c.version, "", "", c.manifest, c.dockerfile, "")
		success := err == nil
		if success != c.expected {
			t.Errorf("DockerBuild(%v, %v, %v, %v) == %v, expected %v", c.directory, c.version, "", "", success, c.expected)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErrorMsg) {
				t.Errorf("DockerBuild(%v, %v, %v, %v) == %v, expected %v", c.directory, c.version, "", "", err.Error(), c.expectedErrorMsg)
			}
		}
	}
}

// func TestDockerBuildPublish(t *testing.T) {
// 	cases := []struct {
// 		directory        string
// 		version          string
// 		expected         bool
// 		imageName        string
// 		registry         string
// 		org              string
// 		force            bool
// 		pkgpatch         bool
// 		pkgmin           bool
// 		pkgmaj           bool
// 		jobpatch         bool
// 		jobmin           bool
// 		jobmaj           bool
// 		expectedImgName  string
// 		expected         bool
// 		expectedErrorMsg string
// 	} {

// 		{"../examples/addition-job/", "1.0.0", true, ""},
// 		{"../examples/extractor/", "1.0.0", true, ""},
// 	}

// 	for _, c := range cases {
// 		_, err := DockerBuild(c.directory, c.version, "", "", ".", ".", "")
// 	}
// }

func TestSeedLabel(t *testing.T) {
	cases := []struct {
		directory        string
		version          string
		imageName        string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/", "1.0.0", "addition-job-0.0.1-seed:1.0.0", true, ""},
		{"../examples/extractor/", "1.0.0", "extractor-0.1.0-seed:0.1.0", true, ""},
	}

	for _, c := range cases {
		DockerBuild(c.directory, c.version, "", "", ".", ".", "")
		seedFileName, exist, _ := util.GetSeedFileName(c.directory)
		if !exist {
			t.Errorf("ERROR: %s cannot be found.\n",
				seedFileName)
			t.Errorf("Make sure you have specified the correct directory.\n")
		}

		// retrieve seed from seed manifest
		seed := objects.SeedFromManifestFile(seedFileName)

		seed2 := objects.SeedFromImageLabel(c.imageName)
		seedStr1 := fmt.Sprintf("%v", seed)
		seedStr2 := fmt.Sprintf("%v", seed2)

		success := seedStr1 != "" && seedStr1 == seedStr2
		if success != c.expected {
			t.Errorf("SeedFromImageLabel(%q) == %v, expected %v", seedFileName, seedStr1, seedStr2)
		}
	}
}

func TestImageName(t *testing.T) {
	cases := []struct {
		filename         string
		expected         string
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json", "addition-job-0.0.1-seed:1.0.0", ""},
		{"../examples/extractor/seed.manifest.json", "extractor-0.1.0-seed:0.1.0", ""},
	}

	for _, c := range cases {

		seedFileName := util.GetFullPath(c.filename, "")
		// retrieve seed from seed manifest
		seed := objects.SeedFromManifestFile(seedFileName)

		// Retrieve docker image name
		imageName := objects.BuildImageName(&seed)

		if imageName != c.expected {
			t.Errorf("BuildImageName(%q) == %v, expected %v", seedFileName, imageName, c.expected)
		}
	}
}
