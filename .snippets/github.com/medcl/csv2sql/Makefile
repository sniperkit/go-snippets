SHELL=/bin/bash

# Default version
APP_VERSION := 0.1.0_SNAPSHOT
APP_NAME := csv2sql

# Get release version from environment
ifneq "$(VERSION)" ""
   APP_VERSION := $(VERSION)
endif

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

# Go environment
CURDIR := $(shell pwd)
OLDGOPATH:= $(GOPATH)
NEWGOPATH:= $(CURDIR):$(CURDIR)/vendor:$(GOPATH)

GO        := GO15VENDOREXPERIMENT="1" go
GOBUILD  := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) build -ldflags -s
GOTEST   := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) test -ldflags -s

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"

.PHONY: all build update test clean

default: build

build: config
	@#echo $(GOPATH)
	@echo $(NEWGOPATH)
	$(GOBUILD) -o bin/$(APP_NAME)
	@$(MAKE) restore-generated-file

build-linux:
	GOOS=linux  GOARCH=amd64  $(GOBUILD) -o bin/$(APP_NAME)-linux64
	GOOS=linux  GOARCH=386    $(GOBUILD) -o bin/$(APP_NAME)-linux32

build-darwin:
	GOOS=darwin  GOARCH=amd64     $(GOBUILD) -o bin/$(APP_NAME)-darwin64
	GOOS=darwin  GOARCH=386       $(GOBUILD) -o bin/$(APP_NAME)-darwin32

cross-build: clean config
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-darwin64
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-linux64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-windows64.exe
	@$(MAKE) restore-generated-file

update-generated-file:
	@echo "update generated info"
	@echo -e "package config\n\nconst LastCommitLog = \""`git log -1 --pretty=format:"%h, %ad, %an, %s"` "\"\nconst BuildDate = \"`date`\"" > config/generated.go
	@echo -e "\nconst Version  = \"$(APP_VERSION)\"" >> config/generated.go


restore-generated-file:
	@echo "restore generated info"
	@echo -e "package config\n\nconst LastCommitLog = \"N/A\"\nconst BuildDate = \"N/A\"" > config/generated.go
	@echo -e "\nconst Version = \"0.0.1-SNAPSHOT\"" >> config/generated.go


format:
	gofmt -l -s -w .

clean_data:
	rm -rif data
	rm -rif log

clean: clean_data
	rm -rif bin
	mkdir bin

init-version:
	@echo building $(APP_NAME) $(APP_VERSION)


update-ui:
	@echo "generate static files"
	@$(GO) get github.com/infinitbyte/esc
	@(cd static&& esc -ignore="static.go|build_static.sh|.DS_Store" -o static.go -pkg static ../static )

update-template-ui:
	@echo "generate UI pages"
	@$(GO) get github.com/infinitbyte/ego/cmd/ego
	@cd ui/ && ego

config: init-version update-ui update-template-ui update-generated-file
	@echo "update configs"
	@# $(GO) env
	@mkdir -p bin
	@cp $(APP_NAME).yml bin/$(APP_NAME).yml

package-darwin-platform:
	@echo "Packaging Darwin"
	cd bin && tar cfz ../bin/darwin64.tar.gz      $(APP_NAME)-darwin64 $(APP_NAME).yml
	cd bin && tar cfz ../bin/darwin32.tar.gz      $(APP_NAME)-darwin32 $(APP_NAME).yml

package-linux-platform:
	@echo "Packaging Linux"
	cd bin && tar cfz ../bin/linux64.tar.gz     $(APP_NAME)-linux64 $(APP_NAME).yml
	cd bin && tar cfz ../bin/linux32.tar.gz     $(APP_NAME)-linux32 $(APP_NAME).yml
