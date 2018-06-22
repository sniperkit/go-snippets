package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ngageoint/seed-cli/constants"
	common_const "github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/util"
	"github.com/xeipuuv/gojsonschema"
)

//DockerRun Runs image described by Seed spec
func DockerRun(imageName, outputDir, metadataSchema string, inputs, json, settings, mounts []string, rmDir, quiet bool) (int, error) {
	util.InitPrinter(util.PrintErr)
	if quiet {
		util.InitPrinter(util.Quiet)
	}

	if imageName == "" {
		return 0, errors.New("ERROR: No input image specified.")
	}

	if exists, err := util.ImageExists(imageName); !exists {
		msg := fmt.Sprintf("Unable to find image: %s. Did you specify a valid tag?", imageName)
		util.PrintUtil("%s\n", msg)
		return 0, err
	}

	// Parse seed information off of the label
	seed := objects.SeedFromImageLabel(imageName)

	// build docker run command
	dockerArgs := []string{"run"}

	if rmDir {
		dockerArgs = append(dockerArgs, "--rm")
	}

	var mountsArgs []string
	var envArgs []string
	var resourceArgs []string
	var inputSize float64
	var outputSize float64

	// expand INPUT_FILEs to specified Inputs files
	if seed.Job.Interface.Inputs.Files != nil {
		inMounts, size, temp, err := DefineInputs(&seed, inputs)
		for _, v := range temp {
			defer util.RemoveAllFiles(v)
		}
		if err != nil {
			util.PrintUtil("ERROR: Error occurred processing inputs arguments.\n%s", err.Error())
			return -1, err
		} else if inMounts != nil {
			mountsArgs = append(mountsArgs, inMounts...)
			inputSize = size
		}
	}

	// add -e args for input json
	if seed.Job.Interface.Inputs.Json != nil {
		inJson, err := DefineInputJson(&seed, json)
		if err != nil {
			util.PrintUtil("ERROR: Error occurred processing json arguments.\n%s", err.Error())
			return -1, err
		} else if inJson != nil {
			envArgs = append(envArgs, inJson...)
		}
	}

	if len(seed.Job.Resources.Scalar) > 0 {
		inResources, diskSize, err := DefineResources(&seed, inputSize)
		if err != nil {
			util.PrintUtil("ERROR: Error occurred processing resources\n%s", err.Error())
			return -1, err
		} else if inResources != nil {
			resourceArgs = append(resourceArgs, inResources...)
			outputSize = diskSize
		}
	}

	// mount the JOB_OUTPUT_DIR (outDir flag)
	var outDir string
	outDir = SetOutputDir(imageName, &seed, outputDir)
	if outDir != "" {
		mountsArgs = append(mountsArgs, "-v")
		mountsArgs = append(mountsArgs, outDir+":"+outDir)
		mountsArgs = append(mountsArgs, "-e")
		mountsArgs = append(mountsArgs, "OUTPUT_DIR="+outDir)
	} else {
		util.PrintUtil("ERROR: Empty output directory string!\n")
	}

	// Settings
	if seed.Job.Interface.Settings != nil {
		inSettings, err := DefineSettings(&seed, settings)
		if err != nil {
			util.PrintUtil("ERROR: Error occurred processing settings arguments.\n%s", err.Error())
			return -1, err
		} else if inSettings != nil {
			envArgs = append(envArgs, inSettings...)
		}
	}

	// Additional Mounts defined in seed.json
	if seed.Job.Interface.Mounts != nil {
		inMounts, err := DefineMounts(&seed, mounts)
		if err != nil {
			util.PrintUtil("ERROR: Error occurred processing mount arguments.\n%s", err.Error())
			return -1, err
		} else if inMounts != nil {
			mountsArgs = append(mountsArgs, inMounts...)
		}
	}

	// Build Docker command arguments:
	// 		run
	//		-rm if specified
	//		env injection
	// 		all mounts
	//		image name
	//		Job.Interface.Command
	dockerArgs = append(dockerArgs, mountsArgs...)
	dockerArgs = append(dockerArgs, envArgs...)
	dockerArgs = append(dockerArgs, resourceArgs...)
	dockerArgs = append(dockerArgs, imageName)

	// Parse out command arguments from seed.Job.Interface.Command
	args := strings.Split(seed.Job.Interface.Command, " ")
	dockerArgs = append(dockerArgs, args...)

	// Run
	var cmd bytes.Buffer
	cmd.WriteString("docker ")
	for _, s := range dockerArgs {
		cmd.WriteString(s + " ")
	}
	util.PrintUtil("INFO: Running Docker command:\n%s\n", cmd.String())

	// Run Docker command and capture output
	dockerRun := exec.Command("docker", dockerArgs...)
	var errs bytes.Buffer
	if !quiet {
		dockerRun.Stderr = io.MultiWriter(&errs)
		dockerRun.Stdout = os.Stderr
	}

	// Run docker run
	runTime := time.Now()
	err := dockerRun.Run()
	util.TimeTrack(runTime, "INFO: "+imageName+" run")
	exitCode := 0
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
			util.PrintUtil("Exited with error code %v\n", exitCode)
			match := false
			for _, e := range seed.Job.Errors {
				if e.Code == exitCode {
					util.PrintUtil("Title: \t %s\n", e.Title)
					util.PrintUtil("Description: \t %s\n", e.Description)
					util.PrintUtil("Category: \t %s \n \n", e.Category)
					match = true
					util.PrintUtil("Exiting seed...\n")
					return exitCode, err
				}
			}
			if !match {
				util.PrintUtil("No matching error code found in Seed manifest\n")
			}
		} else {
			util.PrintUtil("ERROR: error executing docker run. %s\n",
				err.Error())
		}
	}

	if errs.String() != "" {
		util.PrintUtil("ERROR: Error running image '%s':\n%s\n",
			imageName, errs.String())
		util.PrintUtil("Exiting seed...\n")
		return exitCode, errors.New(errs.String())
	}

	// Validate output against pattern
	if seed.Job.Interface.Outputs.Files != nil ||
		seed.Job.Interface.Outputs.JSON != nil {
		CheckRunOutput(&seed, outDir, metadataSchema, outputSize)
	}

	return exitCode, err
}

//DefineInputs extracts the paths to any input data given by the 'run' command
// flags 'inputs' and sets the path in the json object. Returns:
// 	[]string: docker command args for input files in the format:
//	"-v /path/to/file1:/path/to/file1 -v /path/to/file2:/path/to/file2 etc"
func DefineInputs(seed *objects.Seed, inputs []string) ([]string, float64, map[string]string, error) {
	// Validate inputs given vs. inputs defined in manifest

	var mountArgs []string
	var sizeMiB float64

	inMap := inputMap(inputs)

	// Valid by default
	valid := true
	var keys []string
	var unrequired []string
	var tempDirectories map[string]string
	tempDirectories = make(map[string]string)
	for _, f := range seed.Job.Interface.Inputs.Files {
		if f.Multiple {
			tempDir := "temp-" + time.Now().Format(time.RFC3339)
			tempDir = strings.Replace(tempDir, ":", "_", -1)
			os.Mkdir(tempDir, os.ModePerm)
			tempDirectories[f.Name] = tempDir
			mountArgs = append(mountArgs, "-v")
			mountArgs = append(mountArgs, util.GetFullPath(tempDir, "")+":/"+tempDir)
			mountArgs = append(mountArgs, "-e")
			mountArgs = append(mountArgs, f.Name+"=/"+tempDir)
		}
		if f.Required == false {
			unrequired = append(unrequired, f.Name)
			continue
		}
		keys = append(keys, f.Name)
		if _, prs := inMap[f.Name]; !prs {
			valid = false
		}
	}

	if !valid {
		var buffer bytes.Buffer
		buffer.WriteString("ERROR: Incorrect input data files key/values provided. -i arguments should be in the form:\n")
		buffer.WriteString("  seed run -i KEY1=path/to/file1 -i KEY2=path/to/file2 ...\n")
		buffer.WriteString("The following input file keys are expected:\n")
		for _, n := range keys {
			buffer.WriteString("  " + n + "\n")
		}
		buffer.WriteString("\n")
		return nil, 0.0, tempDirectories, errors.New(buffer.String())
	}

	for key, val := range inMap {
		// Expand input VALUE
		val = util.GetFullPath(val, "")

		//get total size of input files in MiB
		info, err := os.Stat(val)
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("ERROR: Input file %s not found\n", val)
			return nil, 0.0, tempDirectories, errors.New(msg)
		}
		sizeMiB += (1.0 * float64(info.Size())) / (1024.0 * 1024.0) //fileinfo's Size() returns bytes, convert to MiB

		// Replace key if found in args strings
		// Handle replacing KEY or ${KEY} or $KEY
		value := val
		if directory, ok := tempDirectories[key]; ok {
			value = directory //replace with the temp directory if multiple files
		}
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
			"${"+key+"}", value, -1)
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command, "$"+key,
			value, -1)

		for _, k := range seed.Job.Interface.Inputs.Files {
			if k.Name == key {
				if k.Multiple {
					//directory has already been added to mount args, just link file into that directory
					os.Link(val, filepath.Join(tempDirectories[key], info.Name()))
				} else {
					mountArgs = append(mountArgs, "-v")
					mountArgs = append(mountArgs, val+":"+val)
					mountArgs = append(mountArgs, "-e")
					mountArgs = append(mountArgs, key+"="+val)
				}
			}
		}
	}

	//remove unspecified unrequired inputs from cmd string
	for _, k := range unrequired {
		key := k
		value := ""
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
			"${"+key+"}", value, -1)
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command, "$"+key,
			value, -1)
	}

	return mountArgs, sizeMiB, tempDirectories, nil
}

//DefineInputJson passes input json values from the 'run' command
// to the image as environment variables.  simple int/string/bool/etc.
// types are passed as single value strings. complex objects can be read
// in from a given file which is read in and passed as an environment variable
// with the full json appropriately escaped
//Returns: []string: docker command args for input files in the format:
//	"-e JSON_NAME1=value -e JSON_NAME2="{json object}" etc"
func DefineInputJson(seed *objects.Seed, inputs []string) ([]string, error) {
	inMap := inputMap(inputs)

	var envArgs []string

	// Valid by default
	valid := true
	var keys []string
	var unrequired []string
	for _, f := range seed.Job.Interface.Inputs.Json {
		if f.Required == false {
			unrequired = append(unrequired, f.Name)
			continue
		}
		keys = append(keys, f.Name)
		if _, prs := inMap[f.Name]; !prs {
			valid = false
		}
	}

	if !valid {
		var buffer bytes.Buffer
		buffer.WriteString("ERROR: Missing json input key.")
		buffer.WriteString("The following input json keys are expected:\n")
		for _, n := range keys {
			buffer.WriteString("  " + n + "\n")
		}
		buffer.WriteString("\n")
		buffer.WriteString("JSON inputs should be provided in the following form: \n")
		buffer.WriteString("seed run -j KEY1=path/to/file1 -j KEY2=path/to/file2 ...\n")
		return nil, errors.New(buffer.String())
	}

	for _, f := range inputs {
		x := strings.SplitN(f, "=", 2)
		if len(x) != 2 {
			util.PrintUtil("ERROR: Input json should be specified in KEY=VALUE format.\n")
			util.PrintUtil("ERROR: Unknown key for input %v encountered.\n",
				inputs)
			continue
		}

		key := x[0]
		val := x[1]

		value, err := ReadJsonFile(val)
		if err != nil {
			value = val
		}

		// Replace key if found in args strings
		// Handle replacing KEY or ${KEY} or $KEY
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
			"${"+key+"}", value, -1)
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command, "$"+key,
			value, -1)

		for _, k := range seed.Job.Interface.Inputs.Json {
			if k.Name == key {
				envArgs = append(envArgs, "-e")
				envArgs = append(envArgs, key+"="+value)
			}
		}
	}

	//remove unspecified unrequired inputs from cmd string
	for _, k := range unrequired {
		key := k
		value := ""
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
			"${"+key+"}", value, -1)
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command, "$"+key,
			value, -1)
	}

	return envArgs, nil
}

//SetOutputDir replaces the OUTPUT_DIR argument with the given output directory.
// Returns output directory string
func SetOutputDir(imageName string, seed *objects.Seed, outputDir string) string {
	// #37: if -o is not specified, auto create a time-stamped subdirectory with the name of the form:
	//		imagename-iso8601timestamp
	if outputDir == "" {
		outputDir = "output-" + imageName + "-" + time.Now().Format(time.RFC3339)
		outputDir = strings.Replace(outputDir, ":", "_", -1)
	}

	outdir := util.GetFullPath(outputDir, "")

	// Check if outputDir exists. Create if not
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		// Create the directory
		// Didn't find the specified directory
		util.PrintUtil("INFO: %s not found; creating directory...\n",
			outdir)
		os.Mkdir(outdir, os.ModePerm)
	}

	// Check if outdir is empty. Create time-stamped subdir if not
	f, err := os.Open(outdir)
	if err != nil {
		// complain
		util.PrintUtil("ERROR: Error with %s. %s\n", outdir, err.Error())
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err != io.EOF {
		// Directory is not empty
		t := time.Now().Format("20060102_150405")
		util.PrintUtil(
			"INFO: Output directory %s is not empty. Creating sub-directory %s for Job Output Directory.\n",
			outdir, t)
		outdir = filepath.Join(outdir, t)
		os.Mkdir(outdir, os.ModePerm)
	}

	seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
		"$OUTPUT_DIR", outdir, -1)
	seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
		"${OUTPUT_DIR}", outdir, -1)
	return outdir
}

//DefineMounts defines any seed specified mounts.
func DefineMounts(seed *objects.Seed, inputs []string) ([]string, error) {
	inMap := inputMap(inputs)

	// Valid by default
	valid := true
	var keys []string
	for _, f := range seed.Job.Interface.Mounts {
		keys = append(keys, f.Name)
		if _, prs := inMap[f.Name]; !prs {
			valid = false
		}
	}

	if !valid {
		var buffer bytes.Buffer
		buffer.WriteString("ERROR: Incorrect mount key/values provided. -m arguments should be in the form:\n")
		buffer.WriteString("  seed run -m MOUNT=path/to ...\n")
		buffer.WriteString("The following mount keys are expected:\n")
		for _, n := range keys {
			buffer.WriteString("  " + n + "\n")
		}
		buffer.WriteString("\n")
		return nil, errors.New(buffer.String())
	}

	var mounts []string
	if seed.Job.Interface.Mounts != nil {
		for _, mount := range seed.Job.Interface.Mounts {
			mounts = append(mounts, "-v")
			localPath := util.GetFullPath(inMap[mount.Name], "")
			mountPath := localPath + ":" + mount.Path

			if mount.Mode != "" {
				mountPath += ":" + mount.Mode
			} else {
				mountPath += ":ro"
			}
			mounts = append(mounts, mountPath)
		}
		return mounts, nil
	}

	return mounts, nil
}

//DefineSettings defines any seed specified docker settings.
// Return []string of docker command arguments in form of:
//	"-e setting1=val1 -e setting2=val2 etc"
func DefineSettings(seed *objects.Seed, inputs []string) ([]string, error) {
	inMap := inputMap(inputs)

	// Valid by default
	valid := true
	var keys []string
	for _, s := range seed.Job.Interface.Settings {
		keys = append(keys, s.Name)
		if _, prs := inMap[s.Name]; !prs {
			valid = false
		}

	}

	if !valid {
		var buffer bytes.Buffer
		buffer.WriteString("ERROR: Incorrect setting key/values provided. -e arguments should be in the form:\n")
		buffer.WriteString("  seed run -e SETTING=somevalue ...\n")
		buffer.WriteString("The following settings are expected:\n")
		for _, n := range keys {
			buffer.WriteString("  " + n + "\n")
		}
		buffer.WriteString("\n")
		return nil, errors.New(buffer.String())
	}

	var settings []string
	for _, key := range keys {
		// Replace key if found in args strings
		// Handle replacing KEY or ${KEY} or $KEY
		value := inMap[key]
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command,
			"${"+key+"}", value, -1)
		seed.Job.Interface.Command = strings.Replace(seed.Job.Interface.Command, "$"+key,
			value, -1)

		settings = append(settings, "-e")
		settings = append(settings, util.GetNormalizedVariable(key)+"="+value)
	}

	return settings, nil
}

//DefineResources defines any seed specified docker resource requirements
//based on the seed spec and the size of the input in MiB
// returns array of arguments to pass to docker to restrict/specify the resources required
// returns the total disk space requirement to be checked when validating output
func DefineResources(seed *objects.Seed, inputSizeMiB float64) ([]string, float64, error) {
	var resources []string
	var disk float64

	for _, s := range seed.Job.Resources.Scalar {
		value := fmt.Sprintf("%f", s.Value)
		if s.Name == "mem" {
			//resourceRequirement = inputVolume * inputMultiplier + constantValue
			mem := (s.InputMultiplier * inputSizeMiB) + s.Value
			mem = math.Max(mem, 4.0)        //docker memory requirement must be > 4MiB
			intMem := int64(math.Ceil(mem)) //docker expects integer, get the ceiling of the specified value and convert
			resources = append(resources, "-m")
			resources = append(resources, fmt.Sprintf("%dm", intMem))
			value = fmt.Sprintf("%d", intMem)
		}
		if s.Name == "disk" {
			//resourceRequirement = inputVolume * inputMultiplier + constantValue
			disk = (s.InputMultiplier * inputSizeMiB) + s.Value
			value = fmt.Sprintf("%f", disk)
		}
		if s.Name == "sharedMem" {
			//resourceRequirement = inputVolume * inputMultiplier + constantValue
			mem := (s.InputMultiplier * inputSizeMiB) + s.Value
			intMem := int64(math.Ceil(mem)) //docker expects integer, get the ceiling of the specified value and convert
			resources = append(resources, fmt.Sprintf("--shm-size=%dm", intMem))
			value = fmt.Sprintf("%d", intMem)
		}

		envVar := util.GetNormalizedVariable(s.Name)
		resources = append(resources, "-e")
		resources = append(resources, fmt.Sprintf("%s=%s", envVar, value))
	}

	return resources, disk, nil
}

//CheckRunOutput validates the output of the docker run command. Output data is
// validated as defined in the seed.Job.Interface.Outputs.
func CheckRunOutput(seed *objects.Seed, outDir, metadataSchema string, diskLimit float64) {
	// Validate any Outputs.Files
	if seed.Job.Interface.Outputs.Files != nil {
		util.PrintUtil("INFO: Validating output files found under %s...\n",
			outDir)

		var dirSize int64
		readSize := func(path string, file os.FileInfo, err error) error {
			if !file.IsDir() {
				dirSize += file.Size()
			}

			return nil
		}
		filepath.Walk(outDir, readSize)
		sizeMB := float64(dirSize) / (1024.0 * 1024.0)
		if diskLimit > 0 && sizeMB > diskLimit {
			util.PrintUtil("ERROR: Output directory exceeds disk space limit (%f MiB vs. %f MiB)\n", sizeMB, diskLimit)
		}

		// For each defined Outputs file:
		//	#1 Check file media type
		// 	#2 Check file names match output pattern
		//  #3 Check number of files (if defined)
		for _, f := range seed.Job.Interface.Outputs.Files {
			// find all pattern matches in OUTPUT_DIR
			matches, _ := filepath.Glob(path.Join(outDir, f.Pattern))

			// Check media type of matches
			count := 0
			var matchList []string
			for _, match := range matches {
				ext := filepath.Ext(match)
				mType := mime.TypeByExtension(ext)
				if strings.Contains(mType, f.MediaType) ||
					strings.Contains(f.MediaType, mType) {
					count++
					matchList = append(matchList, "\t"+match+"\n")
					metadata := match + ".metadata.json"
					if _, err := os.Stat(metadata); err == nil {
						schema := metadataSchema
						if schema != "" {
							schema = util.GetFullPath(schema, "")
						}
						err := ValidateSeedFile(schema, seed.SeedVersion, metadata, common_const.SchemaMetadata)
						if err != nil {
							util.PrintUtil("ERROR: Side-car metadata file %s validation error: %s", metadata, err.Error())
						}
					}
				}
			}

			expected := 1
			errStr := "ERROR: Required file expected for output %v, %v found.\n"
			if f.Multiple == true {
				expected = 2
				errStr = "ERROR: Multiple required files expected for output %v, %v found.\n"
			}

			// Validate that any required fields are present
			if f.Required && len(matchList) < expected {
				util.PrintUtil(errStr, f.Name, strconv.Itoa(len(matchList)))
			} else if !f.Multiple && len(matchList) > 1 {
				util.PrintUtil("WARNING: Multiple files found for single output %v, %v found.\n",
					f.Name, strconv.Itoa(len(matchList)))
				for _, s := range matchList {
					util.PrintUtil(s)
				}
			} else {

				util.PrintUtil("SUCCESS: Files found for output %v:\n",
					f.Name, strconv.Itoa(len(matchList)))
				for _, s := range matchList {
					util.PrintUtil(s)
				}
			}
		}
	}

	// Validate any defined Outputs.Json
	// Look for ResultsFileManifestName.json in the root of the OUTPUT_DIR
	// and then validate any keys identified in Outputs exist
	if seed.Job.Interface.Outputs.JSON != nil {
		util.PrintUtil("INFO: Validating %s...\n",
			filepath.Join(outDir, constants.ResultsFileManifestName))
		// look for results manifest
		manfile := filepath.Join(outDir, constants.ResultsFileManifestName)
		if _, err := os.Stat(manfile); os.IsNotExist(err) {
			util.PrintUtil("ERROR: %s specified but cannot be found. %s\n Exiting testrunner.\n",
				constants.ResultsFileManifestName, err.Error())
			return
		}

		bites, err := ioutil.ReadFile(filepath.Join(outDir,
			constants.ResultsFileManifestName))
		if err != nil {
			util.PrintUtil("ERROR: Error reading %s.%s\n",
				constants.ResultsFileManifestName, err.Error())
			return
		}

		documentLoader := gojsonschema.NewStringLoader(string(bites))
		_, err = documentLoader.LoadJSON()
		if err != nil {
			util.PrintUtil("ERROR: Error loading results manifest file: %s. %s\n Exiting testrunner.\n",
				constants.ResultsFileManifestName, err.Error())
			return
		}

		schemaFmt := "{ \"type\": \"object\", \"properties\": { %s }, \"required\": [ %s ] }"
		schema := ""
		required := ""

		// Loop through defined name/key values to extract from seed.outputs.json
		for _, jsonStr := range seed.Job.Interface.Outputs.JSON {
			key := jsonStr.Name
			if jsonStr.Key != "" {
				key = jsonStr.Key
			}

			schema += fmt.Sprintf("\"%s\": { \"type\": \"%s\" },", key, jsonStr.Type)

			if jsonStr.Required {
				required += fmt.Sprintf("\"%s\",", key)
			}
		}
		//remove trailing commas
		if len(schema) > 0 {
			schema = schema[:len(schema)-1]
		}
		if len(required) > 0 {
			required = required[:len(required)-1]
		}

		schema = fmt.Sprintf(schemaFmt, schema, required)

		schemaLoader := gojsonschema.NewStringLoader(schema)
		schemaResult, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			util.PrintUtil("ERROR: Error running validator: %s\n Exiting testrunner.\n",
				err.Error())
			return
		}

		if len(schemaResult.Errors()) == 0 {
			util.PrintUtil("SUCCESS: Results manifest file is valid.\n")
		}

		for _, desc := range schemaResult.Errors() {
			util.PrintUtil("ERROR: %s is invalid: - %s\n", constants.ResultsFileManifestName, desc)
		}
	}
}

//PrintRunUsage prints the seed run usage arguments, then exits the program
func PrintRunUsage() {
	util.PrintUtil("\nUsage:\tseed run -in IMAGE_NAME [OPTIONS] \n")

	util.PrintUtil("\nRuns Docker image defined by seed spec.\n")

	util.PrintUtil("\nOptions:\n")
	util.PrintUtil("  -%s  -%s \tDocker image name to run\n",
		constants.ShortImgNameFlag, constants.ImgNameFlag)
	util.PrintUtil("  -%s   -%s \t\tSpecifies the key/value input data values of the seed spec in the format INPUT_FILE_KEY=INPUT_FILE_VALUE\n",
		constants.ShortInputsFlag, constants.InputsFlag)
	util.PrintUtil("  -%s   -%s \tSpecifies the key/value setting values of the seed spec in the format SETTING_KEY=VALUE\n",
		constants.ShortSettingFlag, constants.SettingFlag)
	util.PrintUtil("  -%s   -%s \t\tSpecifies the key/value mount values of the seed spec in the format MOUNT_KEY=HOST_PATH\n",
		constants.ShortMountFlag, constants.MountFlag)
	util.PrintUtil("  -%s   -%s \t\tJob Output Directory Location\n",
		constants.ShortJobOutputDirFlag, constants.JobOutputDirFlag)
	util.PrintUtil("  -%s  \t\t\tAutomatically remove the container when it exits (docker run --rm)\n",
		constants.RmFlag)
	util.PrintUtil("  -%s   -%s \t\tSuppress stdout when running docker image\n",
		constants.ShortQuietFlag, constants.QuietFlag)
	util.PrintUtil("  -%s -%s \tRun docker image multiple times (i.e. -rep 5 runs the image 5 times)\n",
		constants.ShortRepeatFlag, constants.RepeatFlag)
	util.PrintUtil("  -%s   -%s \t\tExternal Seed metadata schema file; Overrides built in schema to validate side-car metadata files\n",
		constants.ShortSchemaFlag, constants.SchemaFlag)
	return
}

func inputMap(inputs []string) map[string]string {
	// Ingest inputs into a map key = inputkey, value=inputpath
	inMap := make(map[string]string)
	for _, f := range inputs {
		if f == "" {
			//skip empty strings
			continue
		}
		x := strings.SplitN(f, "=", 2)
		if len(x) != 2 {
			util.PrintUtil("ERROR: Input should be specified in KEY=VALUE format.\n")
			util.PrintUtil("ERROR: Unknown key for input %v encountered.\n",
				x)
			continue
		}
		inMap[x[0]] = x[1]
	}
	return inMap
}

func ReadJsonFile(filename string) (string, error) {
	filebytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	valid := json.Valid(filebytes)
	if !valid {
		return "", errors.New("Invalid JSON")
	}

	dst := new(bytes.Buffer)
	err = json.Compact(dst, filebytes)
	if err != nil {
		return "", err
	}
	json := dst.String()
	return json, err
}