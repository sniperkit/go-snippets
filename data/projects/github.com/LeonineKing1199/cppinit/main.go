package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cleanFlag := flag.Bool("clean", false, "remove all generated files")
	projFlag := flag.String("name", "", "name of the CMake project to be generated")

	flag.Parse()

	if cleanFlag != nil {
		clean := *cleanFlag
		if clean {
			paths := []string{
				".vscode",
				"cmake.modules",
				"include",
				"test",
				"CMakeLists.txt",
				"build_Debug",
				"build_Release",
			}

			matches, err := filepath.Glob("*.cmake")
			if err != nil {
				log.Fatal(err)
			}

			if matches != nil {
				paths = append(paths, matches...)
			}

			for idx := range paths {
				err := os.RemoveAll(paths[idx])
				if err != nil {
					log.Fatal(err)
				}
			}

			return
		}
	}

	if projFlag == nil {
		log.Fatal("projFlag must be non-nil")
	}

	proj := *projFlag
	if proj == "" {
		log.Fatal("project name must be non-empty")
	}

	// first write settings for VSCode
	//
	settingsFile, err := createFileWithPath(".vscode/settings.json")
	if err != nil {
		log.Fatal(err)
	}

	err = writeVSCodeSettings(proj, settingsFile)
	if err != nil {
		log.Fatal(err)
	}

	// write project's main include dir
	//
	err = os.MkdirAll("include/"+proj, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	// create the Catch main file
	//
	testMainFile, err := createFileWithPath("test/main.cpp")
	if err != nil {
		log.Fatal(err)
	}

	err = writeTestMain(testMainFile)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("test/include/"+proj+"/test", os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	// add Catch CTest support
	//
	cmakeModulesFile, err := createFileWithPath("cmake.modules/ParseAndAddCatchTests.cmake")
	if err != nil {
		log.Fatal(err)
	}

	err = writeCMakeModules(cmakeModulesFile)
	if err != nil {
		log.Fatal(err)
	}

	// write the root CMakeLists.txt file
	//
	cmlFile, err := createFileWithPath("CMakeLists.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = writeCMakeLists(proj, cmlFile)
	if err != nil {
		log.Fatal(err)
	}

	// write an empty toolchain file for the user to fill in
	//
	_, err = createFileWithPath(proj + ".cmake")
	if err != nil {
		log.Fatal(err)
	}
}
