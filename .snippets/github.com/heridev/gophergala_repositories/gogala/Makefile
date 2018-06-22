SHELL := /bin/bash

COV_LIB_FILE := cov.lib.out
COV_SRC_FILE := cov.src.out

default: test

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: lint
	go vet ./...

test: fmt
	go test github.com/julien/gogala/lib --coverprofile=$(COV_LIB_FILE)
	go test github.com/julien/gogala --coverprofile=$(COV_SRC_FILE)

cover_lib: test
	go tool cover --html=$(COV_LIB_FILE)

cover_src: test
	go tool cover --html=$(COV_SRC_FILE)

run: fmt
	go run main.go




