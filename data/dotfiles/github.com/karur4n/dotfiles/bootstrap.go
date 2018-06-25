package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

func main() {
	err := installFish()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePaths, err := getDotFilePaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	currentPath, err := os.Getwd()
	user, err := user.Current()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, path := range filePaths {
		sourcePath := currentPath + "/home/" + path
		targetPath := user.HomeDir + "/" + path

		isExist := fileExist(targetPath)

		if isExist {
			continue
		}

		err = os.Symlink(sourcePath, targetPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getDotFilePaths() ([]string, error) {
	files, err := ioutil.ReadDir("./home")
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if !(file.IsDir()) {
			paths = append(paths, file.Name())
		}
	}

	return paths, nil
}

func installFish() error {
	isExist := fileExist("/usr/local/bin/fish")
	if isExist {
		return nil
	}
	err := exec.Command("brew", "install", "fish").Run()
	if err != nil {
		return err
	}
	return nil
}
