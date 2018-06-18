package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	file, err := os.Open("for.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//read file state
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	//read file size
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println(err)
		return
	}

	str := string(bs)
	fmt.Println(str)

	err = os.Remove("/tmp/test/*.conf")
	fmt.Println(err)

	err = os.RemoveAll("/tmp/test/*.conf")
	fmt.Println(err)

	err = filepath.Walk("/tmp/test/", func(path string, info os.FileInfo, err error) error {
		isRegular := info.Mode().IsRegular()
		if isRegular {
			if filepath.Ext(path) == ".conf" {
				fmt.Println(path)
			}
		}
		return nil
	})
	fmt.Println(err)
}
