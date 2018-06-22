default: ci

SCRIPTS_PATH:=$(CURDIR)/../scripts-build-go/script

.PHONY: build check ci dependencies test

dependencies:
	go get -v -d -u "github.com/crucibuild/scripts-build-go"
	go get -v -d ./...

build: dependencies
	$(SCRIPTS_PATH)/resource.sh "$(CURDIR)/resources.go" resources/* schema/*
	go build .

test: build
	go test .

check: dependencies
	$(SCRIPTS_PATH)/check.sh

ci: test check
	true
