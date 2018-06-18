package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	genDir := func(name string) string {

		dir := filepath.Join("data_dir/", "file_storage_path")

		for i := 0; i<3; i++ {
			dir = filepath.Join(dir, name[i*3:(i+1)*3])
		}

		return dir
	}

	str := "6b3ae5fcbeaf49c480baca60f88e7d40"
	fmt.Println(filepath.Join(genDir(str), str))
}