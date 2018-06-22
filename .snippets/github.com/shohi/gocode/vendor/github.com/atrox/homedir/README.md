# HomeDir

[![Build Status](https://img.shields.io/travis/Atrox/homedir.svg?style=flat-square)](https://travis-ci.org/Atrox/homedir)
[![Coverage Status](https://img.shields.io/coveralls/Atrox/homedir.svg?style=flat-square)](https://coveralls.io/r/Atrox/homedir)
[![Go Report Card](https://goreportcard.com/badge/github.com/atrox/homedir?style=flat-square)](https://goreportcard.com/report/github.com/atrox/homedir)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/Atrox/homedir)

> This library is based on and fully compatible with [mitchellh/go-homedir](https://github.com/mitchellh/go-homedir) but uses `os/user` because since go 1.9 there is no longer cgo compilation required.

## Installation

```sh
go get -u github.com/atrox/homedir
# or with dep
dep ensure -add github.com/atrox/homedir
```

## Usage

Usage is incredibly simple, just call `homedir.Dir()` to get the home directory
for a user, and `homedir.Expand(path string)` to expand the `~` in a path to the home
directory.

## Example

```go
package main

import (
    "fmt"

    "github.com/atrox/homedir"
)

func main() {
    dir, err := homedir.Dir()
    if err != nil {
        panic(err)
    }
    fmt.Printf("'%s' is your users home directory\n", dir)

    path, err := homedir.Expand("~/.config")
    if err != nil {
        panic(err)
    }
    fmt.Printf("'%s' is the expanded path to the .config directory\n", path)
}
```

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/atrox/homedir/issues)
- Fix bugs and [submit pull requests](https://github.com/atrox/homedir/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features
