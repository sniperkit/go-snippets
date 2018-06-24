#!/bin/bash

set -e

OS="darwin linux windows"
ARCH="amd64 386"


echo "Ensuring code quality"
#go vet ./...
gofmt -w .

rm -Rf bin/
mkdir bin/

for GOOS in $OS; do
    for GOARCH in $ARCH; do
        arch="$GOOS-$GOARCH"
        binary="dwrap"
        if [ "$GOOS" = "windows" ]; then
          binary="${binary}.exe"
        fi
        echo "Building $binary $arch"
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 govendor build -o $binary cmd/main.go
        zip -r "bin/dwrap_$arch" $binary
        rm -f $binary
    done
done
