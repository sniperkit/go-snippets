#!/usr/bin/env bash

set -e

# Ensure script directory is CWD
pushd "${0%/*}" > /dev/null

VERSION=$1
if [[ "${VERSION}x" == "x" ]]
then
    echo Missing version parameter - setting to snapshot
    VERSION=snapshot
fi


UNAME=$(uname -s)

vendor/go-bindata-${UNAME} -pkg assets -o assets/assets.go ./schema/*
echo Building cross platform Seed CLI.
echo Building for Linux...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-X main.cliVersion=$VERSION -extldflags=\"-static\"" -o output/seed-linux-amd64
echo Building for OSX...
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags "-X main.cliVersion=$VERSION -extldflags=\"-static\"" -o output/seed-darwin-amd64
echo Building for Windows...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "-X main.cliVersion=$VERSION -extldflags=\"-static\"" -o output/seed-windows-amd64
echo CLI build complete

popd >/dev/null
