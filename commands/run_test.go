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

func TestDockerRun(t *testing.T) {
	cases := []struct {
		directory        string
		imageName        string
		inputs           []string
		json             []string
		settings         []string
		mounts           []string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/", "addition-job-0.0.1-seed:1.0.0",
			[]string{"INPUT_FILE=../examples/addition-job/inputs.txt"},
			[]string{},
			[]string{"SETTING_ONE=one", "SETTING_TWO=two"},
			[]string{"MOUNT_BIN=../testdata", "MOUNT_TMP=../testdata"},
			true, ""},
		{"../examples/extractor/", "extractor-0.1.0-seed:0.1.0",
			[]string{"ZIP=../testdata/seed-scale.zip", "MULTIPLE=../testdata/"},
			[]string{},
			[]string{"HELLO=Hello"}, []string{"MOUNTAIN=../examples/"},
			true, ""},
	}

	for _, c := range cases {
		//make sure the image exists
		outputDir := "output"
		metadataSchema := ""
		version := "1.0.0"
		DockerBuild(c.directory, version, "", "", ".", ".", "")
		_, err := DockerRun(c.imageName, outputDir, metadataSchema,
			c.inputs, c.json, c.settings, c.mounts, true, true)
		success := err == nil
		if success != c.expected {
			t.Errorf("DockerRun(%q, %q, %q, %q, %q, %q) == %v, expected %v", c.imageName, outputDir, metadataSchema, c.inputs, c.settings, c.mounts, err, nil)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErrorMsg) {
				t.Errorf("DockerRun(%q, %q, %q, %q, %q, %q) == %v, expected %v", c.imageName, outputDir, metadataSchema, c.inputs, c.settings, c.mounts, err.Error(), c.expectedErrorMsg)
			}
		}
	}
}

func TestDefineInputs(t *testing.T) {
	cases := []struct {
		seedFileName     string
		inputs           []string
		expectedVol      string
		expectedSize     string
		expectedTempDir  string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json",
			[]string{"INPUT_FILE=../examples/addition-job/inputs.txt"},
			"[-v $INPUT_FILE$:$INPUT_FILE$ -e INPUT_FILE=$INPUT_FILE$]", "0.0",
			"map[]", true, ""},
		{"../examples/extractor/seed.manifest.json",
			[]string{"ZIP=../testdata/seed-scale.zip", "MULTIPLE=../testdata/"},
			"[-v $MULTIPLE$:/$MULTIPLETEMP$ -e MULTIPLE=/$MULTIPLETEMP$ -v $ZIP$:$ZIP$ -e ZIP=$ZIP$]", "0.1",
			"map[MULTIPLE:$MULTIPLETEMP$]", true, ""},
	}

	for _, c := range cases {
		seedFileName := util.GetFullPath(c.seedFileName, "")
		seed := objects.SeedFromManifestFile(seedFileName)
		volumes, size, tempDir, err := DefineInputs(&seed, c.inputs)

		if c.expected != (err == nil) {
			t.Errorf("DefineInputs(%q, %q) == %v, expected %v", seedFileName, c.inputs, err, nil)
		}

		expectedVol := c.expectedVol
		expectedTempDir := c.expectedTempDir
		for _, f := range c.inputs {
			x := strings.Split(f, "=")
			tempDir, ok := tempDir[x[0]]
			if ok {
				defer util.RemoveAllFiles(tempDir)
				tempVarStr := fmt.Sprintf("$%sTEMP$", x[0])
				expectedVol = strings.Replace(expectedVol, tempVarStr, tempDir, -1)
				path := util.GetFullPath(tempDir, "")
				replaceNameStr := fmt.Sprintf("$%s$", x[0])
				expectedVol = strings.Replace(expectedVol, replaceNameStr, path, -1)
				expectedTempDir = strings.Replace(expectedTempDir, tempVarStr, tempDir, -1)
			} else {
				path := util.GetFullPath(x[1], "")
				replaceNameStr := fmt.Sprintf("$%s$", x[0])
				expectedVol = strings.Replace(expectedVol, replaceNameStr, path, -1)
			}
		}
		tempStr := fmt.Sprintf("%v", volumes)
		if expectedVol != tempStr {
			t.Errorf("DefineInputs(%q, %q) == \n%v, expected \n%v", seedFileName, c.inputs, tempStr, expectedVol)
		}

		sizeStr := fmt.Sprintf("%.1f", size)
		if c.expectedSize != sizeStr {
			t.Errorf("DefineInputs(%q, %q) == %v, expected %v", seedFileName, c.inputs, sizeStr, c.expectedSize)
		}

		tempStr = fmt.Sprintf("%v", tempDir)
		if expectedTempDir != tempStr {
			t.Errorf("DefineInputs(%q, %q) == \n%v, expected \n%v", seedFileName, c.inputs, tempStr, expectedTempDir)
		}

	}
}

func TestDefineInputJson(t *testing.T) {
	cases := []struct {
		seedFileName     string
		inputs           []string
		expectedSet      string
		expected         bool
		expectedErrorMsg string
	}{
		{"../testdata/complete/seed.manifest.json",
			[]string{"wrong=input"},
			"[]", false, ""},
		{"../examples/addition-job/seed.manifest.json",
			[]string{"a=2", "b=2"},
			"[-e a=2 -e b=2]", true, ""},
		{"../testdata/complete/seed.manifest.json",
			[]string{"INPUT_JSON={a: 1, b: 2}"},
			"[-e INPUT_JSON={a: 1, b: 2}]", true, ""},
	}

	for _, c := range cases {
		seedFileName := util.GetFullPath(c.seedFileName, "")
		seed := objects.SeedFromManifestFile(seedFileName)
		settings, err := DefineInputJson(&seed, c.inputs)

		if c.expected != (err == nil) {
			t.Errorf("DefineInputJson(%q, %q) == %v, expected %v", seedFileName, c.inputs, err, nil)
		}

		expectedSet := c.expectedSet
		tempStr := fmt.Sprintf("%v", settings)
		if expectedSet != tempStr {
			t.Errorf("DefineInputJson(%q, %q) == \n%v, expected \n%v", seedFileName, c.inputs, tempStr, expectedSet)
		}
	}
}

func TestDefineMounts(t *testing.T) {
	cases := []struct {
		seedFileName     string
		mounts           []string
		expectedVol      string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json",
			[]string{"MOUNT_BIN=../testdata", "MOUNT_TMP=../testdata"},
			"[-v MOUNT_BIN:/usr/bin/:ro -v MOUNT_TMP:/tmp/:rw]", true, ""},
		{"../examples/extractor/seed.manifest.json",
			[]string{"MOUNTAIN=../examples/"},
			"[-v MOUNTAIN:/the/mountain:ro]", true, ""},
	}

	for _, c := range cases {
		seedFileName := util.GetFullPath(c.seedFileName, "")
		seed := objects.SeedFromManifestFile(seedFileName)
		volumes, err := DefineMounts(&seed, c.mounts)

		if c.expected != (err == nil) {
			t.Errorf("DefineMounts(%q, %q) == %v, expected %v", seedFileName, c.mounts, err, nil)
		}

		expectedVol := c.expectedVol
		for _, f := range c.mounts {
			x := strings.Split(f, "=")
			path := util.GetFullPath(x[1], "")
			expectedVol = strings.Replace(expectedVol, x[0], path, -1)
		}
		tempStr := fmt.Sprintf("%v", volumes)
		if expectedVol != tempStr {
			t.Errorf("DefineMounts(%q, %q) == \n%v, expected \n%v", seedFileName, c.mounts, tempStr, expectedVol)
		}
	}
}

func TestDefineResources(t *testing.T) {
	cases := []struct {
		seedFileName     string
		inputSize        float64
		expectedResource string
		expectedOutSize  float64
		expectedResult   bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json",
			4.0, "[-e CPUS=0.100000 -m 16m -e MEM=16 -e DISK=5.000000 --shm-size=128m -e SHAREDMEM=128]", 5.0, true, ""},
		{"../examples/extractor/seed.manifest.json",
			1.0, "[-e CPUS=1.000000 -m 16m -e MEM=16 --shm-size=1m -e SHAREDMEM=1 -e DISK=1.010000]", 1.01, true, ""},
		{"../examples/extractor/seed.manifest.json",
			16.0, "[-e CPUS=1.000000 -m 16m -e MEM=16 --shm-size=1m -e SHAREDMEM=1 -e DISK=16.010000]", 16.01, true, ""},
	}

	for _, c := range cases {
		seedFileName := util.GetFullPath(c.seedFileName, "")
		seed := objects.SeedFromManifestFile(seedFileName)
		resources, outSize, err := DefineResources(&seed, c.inputSize)

		if c.expectedResult != (err == nil) {
			t.Errorf("DefineResources(%v, %v) returned unexpected error: %v", seedFileName, c.inputSize, err)
		}

		tempStr := fmt.Sprintf("%v", resources)
		if c.expectedResource != tempStr {
			t.Errorf("DefineResources(%v, %v) == \n%v, expected \n%v", seedFileName, c.inputSize, tempStr, c.expectedResource)
		}

		if c.expectedOutSize != outSize {
			t.Errorf("DefineResources(%v, %v) == \n%v, expected \n%v", seedFileName, c.inputSize, outSize, c.expectedOutSize)

		}
	}
}

func TestDefineSettings(t *testing.T) {
	cases := []struct {
		seedFileName     string
		settings         []string
		expectedSet      string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json",
			[]string{"SETTING_ONE=One", "SETTING_TWO=two"},
			"[-e SETTING_ONE=One -e SETTING_TWO=two]", true, ""},
		{"../examples/extractor/seed.manifest.json",
			[]string{"HELLO=Hello"}, "[-e HELLO=Hello]", true, ""},
	}

	for _, c := range cases {
		seedFileName := util.GetFullPath(c.seedFileName, "")
		seed := objects.SeedFromManifestFile(seedFileName)
		settings, err := DefineSettings(&seed, c.settings)

		if c.expected != (err == nil) {
			t.Errorf("DefineSettings(%q, %q) == %v, expected %v", seedFileName, c.settings, err, nil)
		}

		tempStr := fmt.Sprintf("%v", settings)
		if c.expectedSet != tempStr {
			t.Errorf("DefineSettings(%q, %q) == \n%v, expected \n%v", seedFileName, c.settings, tempStr, c.expectedSet)
		}
	}
}
