package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"io/ioutil"

	"github.com/spf13/viper"
)

var parsedFiles []string

// Sample snippet:
// `CONFIG_DIR=/path-to-config-dir CONFIG_SUFFIX=sample go run viperparse.go`
//
// `CONFIG_SUFFIX` is optional. Above snippet will parse
// configurations from all files which endings contain string "sample".
func main() {
	v := viper.New()

	appConfigDir := os.Getenv("CONFIG_DIR")
	v.AddConfigPath(appConfigDir)
	err := v.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	files, err := ioutil.ReadDir(appConfigDir)
	if err != nil {
		panic(err)
	}

	parseFiles(v, files, os.Getenv("CONFIG_SUFFIX"))

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	fmt.Println("→ Parsed files:\n", parsedFiles)
	fmt.Println("\n→ Parsed value:")

	v.SetDefault("config-key-here", "Please adjust your config key (line 45)")
	spew.Dump(v.Get("config-key-here"))
}

func parseFiles(v *viper.Viper, files []os.FileInfo, fileSuffix string) {
	var suffix string
	if fileSuffix != "" {
		suffix = fmt.Sprintf("\"%s\"", fileSuffix)
	} else {
		suffix = "no"
	}

	fmt.Printf("Parsing %d files.. (using suffix: %s)\n", len(files), suffix)
	for _, file := range files {
		nameSplit := strings.Split(file.Name(), ".")
		name := strings.Join(nameSplit[:len(nameSplit)-1], ".")
		if fileSuffix != "" {
			if len(nameSplit) < 2 {
				continue
			}

			if nameSplit[len(nameSplit)-2] == fileSuffix {
				readConfigName(name, v)
			}
		} else if len(nameSplit) == 2 {
			readConfigName(name, v)
		}
	}
}

func readConfigName(filename string, v *viper.Viper) {
	parsedFiles = append(parsedFiles, filename)
	v.SetConfigName(filename)
	err := v.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
