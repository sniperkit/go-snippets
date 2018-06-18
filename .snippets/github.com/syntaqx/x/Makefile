# vim: ft=make

# Self-documenting Makefiles
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# resolve source directories and files
PACKAGES = $(shell go list ./... | grep -v /vendor/)
SOURCES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# If Go isn't configured correctly, this should circumvent it.
GOPATH ?= $(go env GOPATH)
GOBIN  ?= $(GOPATH)/bin

# Ensure our PATH contains installed go packages
export PATH := $(GOBIN):$(PATH)

.SILENT:
all: help

setup: ## install necessary development tooling
	go get -u github.com/golang/lint/golint
	go get -u github.com/stretchr/testify
	go get -u github.com/unrolled/render

clean: ## remove object files
	go clean -i ./...
	rm -f coverage.txt *.out

test: ## test all packages
	go test -cover -race ./...

cover:
	> coverage.txt
	for d in $(PACKAGES); do \
		go test -race -coverprofile=profile.out -covermode=atomic $$d; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done

cover-html: cover
	go tool cover -html=coverage.txt

fmt: ## run gofmt on all package sources
	bash -c "diff -u <(echo -n) <(gofmt -d $(SOURCES))"

lint: ## to run golint
	golint $(PACKAGES)

vet: ## run go tool vet all packages
	go vet ./...

check: test fmt lint vet ## Runs static code analysis checks

.PHONY: help
help: ## displays help information
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
