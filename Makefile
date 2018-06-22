.PHONY: all assert-go-installed validate-go-version assert-gopath-available list-packages assert-project-under-gopath error-project build install test ci

# project and version
PROJECT_PATH=github.com/kalbasit/swm
VERSION=$(shell git describe --always HEAD)
REQUIRED_GO_VERSION=8
SOURCES := $(shell find . -name "*.go")

all: install

swm: $(SOURCES)
	go build -ldflags "-X main.version=$(VERSION)" -o swm ./cmd/swm/*.go

build: prerequisites vendor swm

vendor: Gopkg.toml Gopkg.lock
	dep ensure -v

install: prerequisites vendor
	go install -v -ldflags "-X main.version=$(VERSION)" ./cmd/swm

test: prerequisites vendor
	go test -v -race -cover -bench=. $(shell go list ./... | grep -v /vendor/)

ci: test

prerequisites:
	@make validate-go-version
	@make assert-gopath-available
	@make assert-project-under-gopath

assert-go-installed:
ifeq (, $(shell which go))
	$(error "No go in $(PATH), please install go from https://golang.org/dl")
endif
	@:

validate-go-version: assert-go-installed
ifneq "$(TRAVIS_GO_VERSION)" "tip"
ifneq "$(shell expr `go version 2>/dev/null | sed -e 's:.*go\([^ ]*\) .*:\1:g' | cut -d. -f2` \>= $(REQUIRED_GO_VERSION) 2>/dev/null)" "1"
	$(error "Minimum required version of Go is 1.$(REQUIRED_GO_VERSION)")
endif
	@:
endif
	@:

assert-gopath-available:
ifndef GOPATH
	$(error GOPATH is not defined, please define it and make sure this project is accessible at $$GOPATH/src/$(PROJECT_PATH))
endif
	@:

list-packages:
	@go list $(PROJECT_PATH)/... | grep -v /vendor/ &>/dev/null

assert-project-under-gopath: assert-gopath-available
	@make list-packages || make error-project

error-project:
	$(error this project is not accessible at $$GOPATH/src/$(PROJECT_PATH))
