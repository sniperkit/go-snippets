# srctool: Tool to manage source code parsers
[![Build Status](https://travis-ci.org/DevMine/srctool.png?branch=master)](https://travis-ci.org/DevMine/srctool)
[![GoDoc](http://godoc.org/github.com/DevMine/srctool?status.svg)](http://godoc.org/github.com/DevMine/srctool)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/DevMine/srctool)

`srctool` is a command line tool to manage source code parsers. It is able to
download parsers from a web server, install them and run them. In short, it is a
manager for source code parsers.


## Installation

### Install from source

Assuming you have [Go](http://golang.org) installed and your `$GOPATH` correctly set, you can
simply issue the following command:

```
go get github.com/DevMine/srctool
```

Make sure that `$GOPATH/bin` is in the `$PATH`.

### Install a binary version

Binary for the latest version can be found [here](http://devmine.ch/downloads/).


## Usage

The command `srctool help` provides a short description of all available
options.

### Setup

Before doing anything, you must first create a configuration file. To do
so, issue:

```
srctool config
```

This will create a default configuration file located in
`$XDG_CONFIG_HOME/srctool/srctool.conf`

If the `$XDG_CONFIG_HOME` environment variable is not set, it will use the
`$HOME/.config` directory.

By default, the download server URL is `http://dl.devmine.ch/parsers`. You can
change this value with the `config` command:

```
srctool config --server-url "http://my-server.com"
```

### Install language parsers

The command `srctool list -r` lists all compatible parsers available on the
download server. For installing them, just issue the following command:

```
srctool install [language]
```

### Parse projects

The main purpose of `srctool` is to parse source code. After installing at least
one parser, issue the following command:

```
srctool parse [project path]
```

This will parse the whole project with each installed parser and merge all the
output to produce a final JSON.

## Running your own download server

Running your own download server requires nothing more than a HTTP server
(nginx, apache, lighttpd, ...). The only requirement is to keep the following
structure for the files organization:

```
.
├── MD5SUMS
├── darwin
│   └── amd64
│       ├── parser-go.zip
│       └── parser-java.zip
├── dragonfly
│   └── amd64
│       ├── parser-go.zip
│       └── parser-java.zip
├── freebsd
│   ├── 386
│   │   ├── parser-go.zip
│   │   └── parser-java.zip
│   └── amd64
│       ├── parser-go.zip
│       └── parser-java.zip
├── linux
│   ├── 386
│   │   ├── parser-go.zip
│   │   └── parser-java.zip
│   └── amd64
│       ├── parser-go.zip
│       └── parser-java.zip
├── netbsd
│   ├── 386
│   │   ├── parser-go.zip
│   │   └── parser-java.zip
│   └── amd64
│       ├── parser-go.zip
│       └── parser-java.zip
└── opendbsd
    ├── 386
    │   ├── parser-go.zip
    │   └── parser-java.zip
    └── amd64
        ├── parser-go.zip
        └── parser-java.zip
```

The folder `tools/` of the project contains several useful scripts for managing
the MD5SUMS file, cross compiling the Go parser, etc. See
[tools/README.md](https://github.com/DevMine/srctool/blob/master/tools/README.md)
for more information.
