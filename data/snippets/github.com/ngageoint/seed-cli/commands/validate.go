package commands

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ngageoint/seed-cli/assets"
	"github.com/ngageoint/seed-cli/constants"
	common_const "github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/util"
	"github.com/xeipuuv/gojsonschema"
)

//Validate seed validate: Validate seed.manifest.json. Does not require docker
func Validate(schemaFile, dir, version string) error {
	var err error = nil
	var seedFileName string

	seedFileName, err = util.SeedFileName(dir)
	if err != nil {
		return err
	}

	if schemaFile != "" {
		schemaFile = "file:///" + util.GetFullPath(schemaFile, dir)
	}

	err = ValidateSeedFile(schemaFile, version, seedFileName, common_const.SchemaManifest)

	return err
}

//PrintValidateUsage prints the seed validate usage, then exits the program
func PrintValidateUsage() {
	util.PrintUtil("\nUsage:\tseed validate [OPTIONS] \n")
	util.PrintUtil("\nValidates the given %s by verifying it is compliant with the Seed spec.\n",
		common_const.SeedFileName)
	util.PrintUtil("\nOptions:\n")
	util.PrintUtil("  -%s -%s\tSpecifies directory in which Seed is located (default is current directory)\n",
		constants.ShortJobDirectoryFlag, constants.JobDirectoryFlag)
	util.PrintUtil("  -%s -%s   \tExternal Seed schema file; Overrides built in schema to validate Seed spec against\n",
		constants.ShortSchemaFlag, constants.SchemaFlag)
	util.PrintUtil(
		"  -%s -%s\tVersion of built in seed manifest to validate against (default is 1.0.0).\n",
		constants.ShortVersionFlag, constants.VersionFlag)
	return
}

//ValidateSeedFile Validates the seed.manifest.json file based on the given schema
func ValidateSeedFile(schemaFile, version, seedFileName string, schemaType common_const.SchemaType) error {
	var result *gojsonschema.Result
	var err error

	seedFileName = strings.Replace(seedFileName, "\\", "/", -1)

	typeStr := "manifest"
	if schemaType == common_const.SchemaMetadata {
		typeStr = "metadata"
	}

	// Load supplied schema file
	if schemaFile != "" {
		util.PrintUtil("INFO: Validating seed %s file %s against schema file %s...\n",
			typeStr, seedFileName, schemaFile)
		schemaLoader := gojsonschema.NewReferenceLoader(schemaFile)
		docLoader := gojsonschema.NewReferenceLoader("file:///" + seedFileName)
		result, err = gojsonschema.Validate(schemaLoader, docLoader)

		// Load baked-in schema file
	} else {
		util.PrintUtil("INFO: Validating seed %s file %s against schema...\n",
			typeStr, seedFileName)
		if version == "" {
			version = "1.0.0"
		}
		assetName := fmt.Sprintf("schema/%s/seed.manifest.schema.json", version)
		schemaBytes, err := assets.Asset(assetName)
		if schemaType == common_const.SchemaMetadata {
			assetName = fmt.Sprintf("schema/%s/seed.metadata.schema.json", version)
			schemaBytes, err = assets.Asset(assetName)
		}

		if schemaBytes == nil || err != nil {
			return fmt.Errorf("This version of seed-cli does not support validating against version %s seed manifests", version)
		}

		schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
		docLoader := gojsonschema.NewReferenceLoader("file:///" + seedFileName)
		result, err = gojsonschema.Validate(schemaLoader, docLoader)
	}

	// Error occurred loading the schema or seed.manifest.json
	if err != nil {
		return errors.New("ERROR: Error validating seed file against schema. Error is:" + err.Error() + "\n")
	}

	// Validation failed. Print results
	var buffer bytes.Buffer
	if !result.Valid() {
		buffer.WriteString("ERROR:" + seedFileName + " is not valid. See errors:\n")
		for _, e := range result.Errors() {
			buffer.WriteString("-ERROR " + e.Description() + "\n")
			buffer.WriteString("\tField: " + e.Field() + "\n")
			buffer.WriteString("\tContext: " + e.Context().String() + "\n")
		}
	}

	//Identify any name collisions for the follwing reserved variables:
	//		OUTPUT_DIR, ALLOCATED_CPUS, ALLOCATED_MEM, ALLOCATED_SHARED_MEM, ALLOCATED_STORAGE
	util.PrintUtil("INFO: Checking for variable name collisions...\n")
	seed := objects.SeedFromManifestFile(seedFileName)

	//skip resource and name collision checking for metadata files
	if schemaType != common_const.SchemaManifest {
		if buffer.String() != "" {
			return errors.New(buffer.String())
		}
		return nil
	}

	recommendedResources := []string{"mem", "cpus", "disk"}
	if seed.Job.Resources.Scalar != nil {
		for _, s := range seed.Job.Resources.Scalar {
			recommendedResources = util.RemoveString(recommendedResources, s.Name)
		}
	}
	if len(recommendedResources) > 0 {
		util.PrintUtil("WARNING: %s does not specify some recommended resources\n", seedFileName)
		util.PrintUtil("Specifying cpu, memory and disk requirements are highly recommended\n")
		util.PrintUtil("The following resources are not defined: %s\n", recommendedResources)
	}

	// Grab all scalar resource names (verify none are set to OUTPUT_DIR)
	var allocated []string
	// var vars map[string]string
	vars := make(map[string][]string)
	if seed.Job.Resources.Scalar != nil {
		for _, s := range seed.Job.Resources.Scalar {
			name := util.GetNormalizedVariable(s.Name)
			allocated = append(allocated, "ALLOCATED_"+strings.ToUpper(name))
			if util.IsReserved(s.Name, nil) {
				buffer.WriteString("ERROR: job.resources.scalar Name " +
					s.Name + " is a reserved variable. Please choose a different name value.\n")
			}

			util.IsInUse(s.Name, "job.resources.scalar", vars)
		}
	}

	if seed.Job.Interface.Inputs.Files != nil {
		for _, f := range seed.Job.Interface.Inputs.Files {
			// check against the ALLOCATED_* and OUTPUT_DIR
			if util.IsReserved(f.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.inputs.files Name " +
					f.Name + " is a reserved variable. Please choose a different name value.\n")
			}

			util.IsInUse(f.Name, "job.interface.inputs.files", vars)
		}
	}

	if seed.Job.Interface.Inputs.Json != nil {
		for _, f := range seed.Job.Interface.Inputs.Json {
			if util.IsReserved(f.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.inputs.json Name " +
					f.Name + " is a reserved variable. Please choose a different name value.\n")
			}

			util.IsInUse(f.Name, "job.interface.inputs.json", vars)
		}
	}

	if seed.Job.Interface.Outputs.Files != nil {
		for _, f := range seed.Job.Interface.Outputs.Files {
			// check against the ALLOCATED_* and OUTPUT_DIR
			if util.IsReserved(f.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.outputs.files Name " +
					f.Name + " is a reserved variable. Please choose a different name value.\n")
			}
			util.IsInUse(f.Name, "job.interface.outputData.files", vars)
		}
	}

	if seed.Job.Interface.Outputs.JSON != nil {
		for _, f := range seed.Job.Interface.Outputs.JSON {
			// check against the ALLOCATED_* and OUTPUT_DIR
			if util.IsReserved(f.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.outputData.json Name " +
					f.Name + " is a reserved variable. Please choose a different name value.\n")
			}
			util.IsInUse(f.Name, "job.interface.outputData.json", vars)
		}
	}

	if seed.Job.Interface.Mounts != nil {
		for _, m := range seed.Job.Interface.Mounts {
			// check against the ALLOCATED_* and OUTPUT_DIR
			if util.IsReserved(m.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.mounts Name " + m.Name +
					" is a reserved variable. Please choose a different name value.\n")
			}
			util.IsInUse(m.Name, "job.interface.mounts", vars)
		}
	}

	if seed.Job.Interface.Settings != nil {
		for _, s := range seed.Job.Interface.Settings {
			// check against the ALLOCATED_* and OUTPUT_DIR
			if util.IsReserved(s.Name, allocated) {
				buffer.WriteString("ERROR: job.interface.settings Name " + s.Name +
					" is a reserved variable. Please choose a different name value.\n")
			}
			util.IsInUse(s.Name, "job.interface.settings", vars)
		}
	}

	// Find any name collisions
	for key, val := range vars {
		if len(val) > 1 {
			buffer.WriteString("ERROR: Multiple Name values are assigned the same " +
				key + " Name value. Each Name value must be unique.\n")
			for _, v := range val {
				buffer.WriteString("\t" + v + "\n")
			}
		}
	}

	// Return error if issues found
	if buffer.String() != "" {
		return errors.New(buffer.String())
	}

	// Validation succeeded
	util.PrintUtil("SUCCESS: No errors found. %s is valid.\n\n", seedFileName)
	return nil
}
