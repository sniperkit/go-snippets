.PHONY: all test clean man glide fast release install help 

################################################################################################
## local - runtime

ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  RUNTIME_OS_SLUG			:= osx
else
  RUNTIME_OS_SLUG 			:= nix
endif
RUNTIME_OS_VERSION 			?= $(shell uname -r)
RUNTIME_OS_ARCH 			?= $(shell uname -m)
RUNTIME_OS_INFO 			?= $(shell uname -a)
RUNTIME_OS_NAME 			?= $(shell uname -s)

################################################################################################
## local - build

## program
PROG_NAME 					:= storage
PROG_NAME_SUFFIX 			:= colly-
PROG_SRCS 					:= $(shell git ls-files '*.go' | grep -v '^vendor/')

## local build
BIN_PREFIX_DIR 				:= ./bin
BIN_BASE_NAME 				:= $(PROG_NAME_SUFFIX)$(PROG_NAME)
BIN_FILE_PATH 				:= $(BIN_PREFIX_DIR)/$(BIN_BASE_NAME)

## local dist
DIST_PREFIX_DIR 			:= ./dist
DIST_BASE_NAME 				:= $(PROG_NAME_SUFFIX)$(PROG_NAME)
DIST_FILE_PATH 				:= $(DIST_PREFIX_DIR)/$(DIST_BASE_NAME)
DIST_ARCHS 					:= "linux darwin"
DIST_OSS 					:= "amd64"

################################################################################################
## docker

#### build
DOCKER_PREFIX_DIR 			:= ./docker
DOCKER_BIN_FILE_PATH 		:= $(DOCKER_PREFIX_DIR)/$(BIN_BASE_NAME)

#### image
DOCKER_IMAGE_OWNER 			:= sniperkit
DOCKER_IMAGE_BASENAME 		:= colly-storage
DOCKER_IMAGE_TAG 			:= 3.7-alpine
DOCKER_IMAGE 				:= $(DOCKER_IMAGE_OWNER)/$(DOCKER_IMAGE_BASENAME):$(DOCKER_IMAGE_TAG)
DOCKER_MULTI_STAGE_IMAGE 	:= $(DOCKER_IMAGE_OWNER)/$(DOCKER_IMAGE_BASENAME)-multi:$(DOCKER_IMAGE_TAG)

################################################################################################
## version

# vcs
REPO_VCS 		:= github.com
REPO_OWNER 		:= sniperkit
REPO_NAME 		:= colly-storage
REPO_URI 		:= $(REPO_VCS)/$(REPO_OWNER)/$(REPO_NAME)
REPO_BRANCH 	:= $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))

#### vcs - commit 
COMMIT_ID   	?= $(shell git describe --tags --always --dirty=-dev)
COMMIT_UNIX 	?= $(shell git show -s --format=%ct HEAD)
COMMIT_HASH 	?= $(shell git rev-parse HEAD)

#### semantic version 
BUILD_COUNT 	?= $(shell git rev-list --count HEAD)
BUILD_UNIX  	?= $(shell date +%s)
BUILD_VERSION 	:= $(shell cat $(CURDIR)/VERSION)
BUILD_TIME 		:= $(shell date)

################################################################################################
## golang

GO15VENDOREXPERIMENT=1
BUILD_LDFLAGS = \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).VERSION=$(BUILD_VERSION)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).branchName=$(REPO_BRANCH)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).commitHash=$(COMMIT_HASH)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).commitID=$(COMMIT_ID)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).commitUnix=$(COMMIT_UNIX)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).buildVersion=$(BUILD_VERSION)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).buildCount=$(BUILD_COUNT)' \
	-X '$(REPO_URI)/cmd/$(PROG_NAME).buildUnix=$(BUILD_UNIX)'

################################################################################################
## makefile
INFO_BREAKLINE := "\n"
INFO_HEADER := "$(INFO_BREAKLINE)------------------------------------------------------------------------------------------"

INFO_FOOTER := "$(INFO_BREAKLINE)------------------------------$(INFO_BREAKLINE)"

default: help

all: deps test build install version dist ## Trigger targets for generating a new release: deps, test, build, install, version and dist targets

info: clear info-runtime info-vcs info-docker info-footer ## Print all Makefile related variables

clear: ## Clear terminal screen 
	@clear

info-header:
	@echo ""
	@echo "------------------------------"

info-footer:
	@echo "$(INFO_FOOTER)"

info-runtime:  ## Print local runtime env variables
	@echo "$(INFO_HEADER)"
	@echo "Runtime:"
	@echo " - RUNTIME_OS_NAME: $(RUNTIME_OS_NAME)"
	@echo " - RUNTIME_ARCH: $(RUNTIME_OS_ARCH)"
	@echo " - RUNTIME_OS_VERSION: $(RUNTIME_OS_VERSION)"
	@echo " - RUNTIME_OS_SLUG: $(RUNTIME_OS_SLUG)"
	@echo " - RUNTIME_OS_INFO: $(RUNTIME_OS_INFO)"

info-vcs:  ## Print source-control related variables
	@echo "$(INFO_HEADER)"
	@echo "Source-Control:"
	@echo " - REPO_URI: $(REPO_URI)"
	@echo " - REPO_BRANCH: $(REPO_BRANCH)"
	@echo " - COMMIT_ID: $(COMMIT_ID)"
	@echo " - COMMIT_UNIX: $(COMMIT_UNIX)"
	@echo " - COMMIT_HASH: $(COMMIT_HASH)"
	@echo " - BUILD_COUNT: $(BUILD_COUNT)"
	@echo " - BUILD_UNIX: $(BUILD_UNIX)"
	@echo " - BUILD_VERSION: $(BUILD_VERSION)"
	@echo " - BUILD_TIME: $(BUILD_TIME)"

info-docker:
	@echo "$(INFO_HEADER)"
	@echo "Docker:"
	@echo " - DOCKER_PREFIX_DIR: $(DOCKER_PREFIX_DIR)"
	@echo " - DOCKER_BIN_FILE_PATH: $(DOCKER_BIN_FILE_PATH)"
	@echo " - DOCKER_IMAGE_OWNER: $(DOCKER_IMAGE_OWNER)"
	@echo " - DOCKER_IMAGE_TAG: $(DOCKER_IMAGE_TAG)"
	@echo " - DOCKER_IMAGE: $(DOCKER_IMAGE)"
	@echo " - DOCKER_MULTI_STAGE_IMAGE: $(DOCKER_MULTI_STAGE_IMAGE)"

docker: docker-build  # docker-tag docker-commit docker-push ## Generate, tag and push a new docker image for this program.

docker-quick: docker-build docker-run ## Build and run quickly a docker container for this program

#docker-multistage: ## Build docker multi-stage container
#	@cd $(DOCKER_PREFIX_DIR) 
#	@docker build --force-rm -t $(DOCKER_MULTI_STAGE_IMAGE) --no-cache -f $(CURDIR)/docker/dockerfile-multi-stage-alpine3.7 .

docker-build: ## Build docker container
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(BUILD_LDFLAGS)" -o $(DOCKER_BIN_FILE_PATH)-linux -v *.go
	@cd $(DOCKER_PREFIX_DIR) && docker build --force-rm -t $(DOCKER_IMAGE) --no-cache -f dockerfile-alpine3.7 .

docker-run: ## Run docker container locally
	@docker run -ti --rm $(DOCKER_IMAGE)

docker-info: ## Get docker client info and env variables
	@echo "'docker-info' is not implemented yet..."

docker-summary: ## Get docker image(s)/container(s) summary 
	@echo "'docker-summary' is not implemented yet..."

docker-commit: ## Commit latest docker image for this program
	@echo "'docker-commit' is not implemented yet..."

docker-tag: ## Tag latest docker image for this program
	@echo "'docker-push' is not implemented yet..."

docker-push: ## Push docker image to image registry
	@echo "'docker-push' is not implemented yet..."

build: ## Build binary for local operating system 
	@go build -ldflags "$(BUILD_LDFLAGS)" -o $(BIN_FILE_PATH) *.go

install: deps ## Install binary in your GOBIN path
	@go build -ldflags "$(BUILD_LDFLAGS)" -o $(BIN_FILE_PATH) *.go
	@$(BIN_BASE_NAME) --version

xbuild: ## Build binaries for linux, darwin in amd64 arch.
	@gox build -ldflags "$(BUILD_LDFLAGS)" -os="darwin linux" -arch="amd64" -output="$(DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" *.go

version-current: ## Check current version of command build
	@which $(BIN_BASE_NAME)
	@$(BIN_BASE_NAME) --version

clean: ## Clean previous build outputs 
	@go clean
	@rm -fr ./bin/$(BIN_FILE_PATH)
	@rm -fr ./dist/$(BIN_FILE_PATH)*

release: $(PROG_NAME) ## Push a new release version to the remote repository
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

cover: ## Execute coverage tests
	@rm -rf *.coverprofile
	@go test -coverprofile=$(PROG_NAME).coverprofile ./pkg/...
	@gover
	@go tool cover -html=$(PROG_NAME).coverprofile ./pkg/...

deps: deps-ensure deps-dev deps-test ## Ensure all required dependencies and helpers

deps-all: deps-create deps-ensure deps-dev deps-test ## Re-create all dependencies list and ensure all locally

deps-create: ## Create program's dependencies list
	@rm -f glide.*
	@rm -f *Gopkg*
	@yes no | glide create

deps-ensure: ## Ensure locally all external dependencies required (package manager: glide)
	@glide install --strip-vendor

deps-dev: ## Install required build helpers in GOBIN 
	@go get -v -u github.com/sniperkit/crane/cmd/crane
	@go get -v -u github.com/sniperkit/gox/cmd/gox

deps-test:  ## Install required program testing an ci helpers in GOBIN
	@go get -v -u github.com/go-playground/overalls
	@go get -v -u github.com/mattn/goveralls
	@go get -v -u golang.org/x/tools/cmd/cover
	@go get -v -u github.com/alexkohler/prealloc
	@go get -v -u github.com/FiloSottile/vendorcheck
	@go get -v -u github.com/golang/dep/cmd/dep
	@go get -v -u github.com/golang/lint/golint
	@go get -v -u github.com/kisielk/errcheck
	@go get -v -u github.com/mdempsky/unconvert
	@go get -v -u github.com/opennota/check/...
	@go get -v -u honnef.co/go/tools/...
	@go get -v -u mvdan.cc/interfacer
	@go get -v -u github.com/dominikh/go-tools/...

lint: ## Lint program's source code
	@errors=$$(gofmt -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi
	@errors=$$(glide novendor | xargs -n 1 golint -min_confidence=0.3); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

vet: ## Vet program's source code
	@go vet $$(glide novendor)

errcheck: ## Check for errors
	@errcheck $(PACKAGES)

interfacer: ## Suggest interface types
	@interfacer $(PACKAGES)

aligncheck: ## Find inefficiently packed structs
	@aligncheck $(PACKAGES)

structcheck: ## Find unused struct fields
	@structcheck $(PACKAGES)

varcheck: ## Find unused global variables and constants
	@varcheck $(PACKAGES)

unconvert: ## Remove unnecessary type conversions from Go source
	@unconvert -v $(PACKAGES)

gosimple: ## Suggest code simplifications
	@gosimple $(PACKAGES)

staticcheck: ## Execute a ton of static analysis checks
	@staticcheck $(PACKAGES)

unused: ## Find for unused constants, variables, functions and types. 
	@unused $(PACKAGES)

vendorcheck: ## Check that all Go dependencies are properly vendored
	@vendorcheck $(PACKAGES)
	@vendorcheck -u $(PACKAGES)

prealloc: ## Find slice declarations that could potentially be preallocated.
	@prealloc $(PACKAGES)

test: ## Execute cover tests on program's sources
	@go test -cover $(PACKAGES)

coverage: ## Execute all coverage tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out

help: ## Display the list of available targets.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

#generate-webui: build-webui ## Build the web front-end
#	if [ ! -d "static" ]; then \
#		mkdir -p static; \
#		docker run --rm -v "$$PWD/static":'/src/static' dtk-webui npm run build; \
#		echo 'For more informations show `webui/readme.md`' > $$PWD/static/DONT-EDIT-FILES-IN-THIS-DIRECTORY.md; \
#	fi