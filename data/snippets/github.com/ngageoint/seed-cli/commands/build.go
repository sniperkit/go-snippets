package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/ngageoint/seed-cli/constants"
	common_const "github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/util"
)

//DockerBuild Builds the docker image with the given image tag.
func DockerBuild(jobDirectory, version, username, password, manifest, dockerfile, cacheFrom string) (string, error) {
	if username != "" {
		//set config dir so we don't stomp on other users' logins with sudo
		configDir := common_const.DockerConfigDir + time.Now().Format(time.RFC3339)
		os.Setenv(common_const.DockerConfigKey, configDir)
		defer util.RemoveAllFiles(configDir)
		defer os.Unsetenv(common_const.DockerConfigKey)

		registry, err := util.DockerfileBaseRegistry(jobDirectory)
		if err != nil {
			util.PrintUtil("Error getting registry from dockerfile: %s\n", err.Error())
		}
		err = util.Login(registry, username, password)
		if err != nil {
			util.PrintUtil("Error calling docker login: %s\n", err.Error())
		}
	}

	var seedFileName string
	var err error
	if manifest != "." {
		seedFileName = util.GetFullPath(manifest, "")
		if _, err = os.Stat(seedFileName); os.IsNotExist(err) {
			util.PrintUtil("ERROR: Seed manifest not found. %s\n", err.Error())
			return "", err
		}
	} else {
		seedFileName, err = util.SeedFileName(jobDirectory)
		if err != nil && !os.IsNotExist(err) {
			util.PrintUtil("ERROR: %s\n", err.Error())
			return "", err
		}
	}

	// Validate seed file
	err = ValidateSeedFile("", version, seedFileName, common_const.SchemaManifest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: seed file could not be validated. See errors for details.")
		util.PrintUtil("%s", err.Error())
		util.PrintUtil("Exiting seed...\n")
		return "", err
	}

	// retrieve seed from seed manifest
	seed := objects.SeedFromManifestFile(seedFileName)

	// Retrieve docker image name
	imageName := objects.BuildImageName(&seed)

	// Build Docker image
	util.PrintUtil("INFO: Building %s\n", imageName)
	buildArgs := []string{"build"}

	// docker doesn't care about validating the cache-from image
	if cacheFrom != "" {
		buildArgs = append(buildArgs, "--cache-from")
		buildArgs = append(buildArgs, cacheFrom)
	}

	buildArgs = append(buildArgs, "-t")
	buildArgs = append(buildArgs, imageName)

	util.PrintUtil("dockerfile: %s\n", dockerfile)
	if dockerfile != "." {
		dfile := util.GetFullPath(dockerfile, "")
		if _, err = os.Stat(dfile); os.IsNotExist(err) {
			util.PrintUtil("ERROR: Dockerfile not found. %s\n", err.Error())
			return imageName, err
		}
		buildArgs = append(buildArgs, dfile)
	} else {
		buildArgs = append(buildArgs, jobDirectory)
	}

	if util.DockerVersionHasLabel() {
		// Set the seed.manifest.json contents as an image label
		label := "com.ngageoint.seed.manifest=" + objects.GetManifestLabel(seedFileName)
		buildArgs = append(buildArgs, "--label", label)
	}
	cmd := exec.Command("docker", buildArgs...)
	var errs bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errs)
	cmd.Stdout = os.Stderr

	// Run docker build
	if err := cmd.Run(); err != nil {
		util.PrintUtil("ERROR: Error executing docker build. %s\n",
			err.Error())
		return imageName, err
	}

	// check for errors on stderr
	if errs.String() != "" {
		util.PrintUtil("ERROR: Error building image '%s':\n%s\n",
			imageName, errs.String())
		util.PrintUtil("Exiting seed...\n")
		return imageName, errors.New(errs.String())
	}

	return imageName, nil
}

//PrintBuildUsage prints the seed build usage arguments, then exits the program
func PrintBuildUsage() {
	util.PrintUtil("\nUsage:\tseed build [-d JOB_DIRECTORY]\n")
	util.PrintUtil("\nOptions:\n")
	util.PrintUtil("  -%s -%s  Utilizes the --cache-from option when building the docker image\n",
		constants.ShortCacheFromFlag, constants.CacheFromFlag)
	util.PrintUtil(
		"  -%s -%s\t  Directory containing Seed spec and Dockerfile (default is current directory)\n",
		constants.ShortJobDirectoryFlag, constants.JobDirectoryFlag)
	util.PrintUtil("  -%s -%s  Specifies the Dockerfile to use (default is Dockerfile within current directory)\n",
		constants.ShortDockerfileFlag, constants.DockerfileFlag)
	util.PrintUtil("  -%s -%s\t  Specifies the seed manifest file to use (default is seed.manifest.json within the current directory)\n",
		constants.ShortManifestFlag, constants.ManifestFlag)
	util.PrintUtil(
		"  -%s -%s\t  Version of built in seed manifest to validate against (default is 1.0.0).\n",
		constants.ShortVersionFlag, constants.VersionFlag)
	util.PrintUtil("  -%s -%s\t  Username to login if needed to pull images (default anonymous).\n",
		constants.ShortUserFlag, constants.UserFlag)
	util.PrintUtil("  -%s -%s\t  Password to login if needed to pull images (default anonymous).\n",
		constants.ShortPassFlag, constants.PassFlag)

	util.PrintUtil("\nBuild and Publish options:\n")
	util.PrintUtil("  -%s\t  Will publish image after a successful build.\n",
		constants.PublishCommand)
	util.PrintUtil("  -%s -%s\t  Specifies a specific registry to publish the image\n",
		constants.ShortRegistryFlag, constants.RegistryFlag)
	util.PrintUtil("  -%s -%s\t  Specifies a specific organization to publish the image\n",
		constants.ShortOrgFlag, constants.OrgFlag)
	util.PrintUtil("  -%s\t\t  Overwrite remote image if publish conflict found\n",
		constants.ForcePublishFlag)

	util.PrintUtil("\nPublish Conflict Options:\n")
	util.PrintUtil("If the force flag (-f) is not set, the following options specify how a publish conflict is handled:\n")
	util.PrintUtil("  -%s\t\t  Force Patch version bump of 'packageVersion' in manifest on disk if publish conflict found\n",
		constants.PkgVersionPatch)
	util.PrintUtil("  -%s\t\t  Force Minor version bump of 'packageVersion' in manifest on disk if publish conflict found\n",
		constants.PkgVersionMinor)
	util.PrintUtil("  -%s\t\t  Force Major version bump of 'packageVersion' in manifest on disk if publish conflict found\n",
		constants.PkgVersionMajor)
	util.PrintUtil("  -%s\t\t  Force Patch version bump of 'jobVersion' in manifest on disk if publish conflict found\n",
		constants.JobVersionPatch)
	util.PrintUtil("  -%s\t\t  Force Minor version bump of 'jobVersion' in manifest on disk if publish conflict found\n",
		constants.JobVersionMinor)
	util.PrintUtil("  -%s\t\t  Force Major version bump of 'jobVersion' in manifest on disk if publish conflict found\n",
		"JM")
	return
}
