.PHONY: all install build bench test cover lint

TIMESTAMP = $(shell date +"%Y-%m-%d_%H-%M-%S")
PKG = $(shell go list ./... | grep -v /vendor/)

all: install build test

install:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover

build:
	go build ${PKG}

bench:
	go test -bench=. -benchmem ${PKG}

test:
	go test -v ${PKG}

cover:
	echo "" > coverage.txt
	for d in ${PKG}; \
		do echo "" > profile.out; \
		go test -coverprofile=profile.out -covermode=set $$d; \
		cat profile.out >> coverage.txt; \
		rm profile.out; \
	done

lint:
	golint ${PKG}
	go vet ${PKG}