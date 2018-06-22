GOPATH := $(shell pwd)
GOROOT := $(HOME)/App/go
GO := $(GOROOT)/bin/go
SRC := $(GOPATH)/src
GOBIN := $(GOPATH)/bin

.PHONY: all
all: build

.PHONY: setup build clean


build : clean setup
	echo "building";
	export GOPATH=$(GOPATH) && \
	mkdir bin && \
	cd bin && \
	$(GO) build $(SRC)/main.go
	echo "success"

run : clean setup
	echo "running";
	export GOPATH=$(GOPATH) && \
	mkdir bin && \
	cd bin && \
	$(GO) build $(SRC)/main.go
	bin/main
	echo "success"
clean :
	rm -rf bin
test :
	export GOPATH=$(GOPATH) && \
	$(GO) test -ttimeout=30s  ./...

.PHONY: setup
setup:
	export GOROOT=$(GOROOT)
	export GOBIN=$GOPATH/bin
	export PATH=$PATH:$GOBIN
