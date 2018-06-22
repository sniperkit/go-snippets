# gopathexec [![Build Status](https://travis-ci.org/sourcegraph/gopathexec.svg?branch=master)](https://travis-ci.org/sourcegraph/gopathexec) [![GoDoc](https://godoc.org/sourcegraph.com/sourcegraph/gopathexec?status.svg)](https://godoc.org/sourcegraph.com/sourcegraph/gopathexec)

Command gopathexec executes program with arguments, while expanding $GOPATH with multiple workspaces into multiple arguments.

Usage: gopathexec program [args]

It is helpful for writing go generate directives for programs that do not understand import paths, but you need to specify include
paths that are other Go package folders.

For example, consider a go generate directive like this:

	//go:generate protoc -I$GOPATH/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

That will only work if your GOPATH env var happens to contain one workspace. By prepending the above command with gopathexec:

	//go:generate gopathexec protoc -I$GOPATH/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

It will effectively execute:

	protoc -I/workspace1/src/github.com/gogo/protobuf/protobuf -I/workspace2/src/github.com/gogo/protobuf/protobuf --go_out=. sample.proto

If you have 1 or no GOPATH workspaces, gopathexec executes the given program and arguments without modification.

Installation
------------

```bash
go get -u sourcegraph.com/sourcegraph/gopathexec
```

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
