[![Go Report Card](https://goreportcard.com/badge/github.com/peeyushsrj/golang-snippets)](https://goreportcard.com/report/github.com/peeyushsrj/golang-snippets)
![Build Status](https://circleci.com/gh/peeyushsrj/golang-snippets.png?style=shield)

# golang-snippets

Some useful functions.

## Subpackages

### files
- [Read line by line from file pointer](https://github.com/peeyushsrj/golang-snippets/blob/master/files/read-lines-from-file-pointer.go)
- [Read line by line from file path](https://github.com/peeyushsrj/golang-snippets/blob/master/files/read-lines-from-file-path.go)
- [List files in a directory](https://github.com/peeyushsrj/golang-snippets/blob/master/files/browse-files.go)
- [Get filename from URL](https://github.com/peeyushsrj/golang-snippets/blob/master/files/filename-from-url.go)
- [Append text to file](https://github.com/peeyushsrj/golang-snippets/blob/master/files/append-text-to-file.go)


### http

- [Anonymous visit website via socks5 client](https://github.com/peeyushsrj/golang-snippets/blob/master/http/append-text-to-file.go) (Pre-requisite: tor must be installed)


## Install

```
go get -u github.com/peeyushsrj/golang-snippets/...
```

### Example

```golang
package main

import (
	"fmt"
	"log"
	"os"
	files "github.com/peeyushsrj/golang-snippets/files"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify the directory")
		fmt.Println("Usage: [path]\n")
		return
	}
	basePath := os.Args[1]
	if _, err := os.Stat(basePath); err == nil {
		fileNames, err := files.BrowseXFiles(".md", basePath)
		if err != nil {
			log.Fatal(1)
		}
    		fmt.Println(fileNames)
	} else {
		fmt.Println("The input path doesn't exist.")
	}
}
```
