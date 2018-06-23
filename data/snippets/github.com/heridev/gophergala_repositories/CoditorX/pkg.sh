#!/bin/bash

# CoditorX package tool.
# 
# Command: 
#  ./pkg.sh ${version} ${target}
# Example:
#  ./pkg.sh 1.0.0 /home/daniel/CoditorX

ver=$1
target=$2
list="conf i18n static views workspaces README.md LICENSE"

mkdir -p ${target}

echo version=${ver}
echo target=${target}

## darwin
os=darwin

export GOOS=${os}
export GOARCH=386
go build
tar zcvf ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.tar.gz ${list} CoditorX --exclude-vcs --exclude conf/*.go --exclude i18n/*.go
rm -f CoditorX

export GOOS=${os}
export GOARCH=amd64
go build
tar zcvf ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.tar.gz ${list} CoditorX --exclude-vcs --exclude conf/*.go --exclude i18n/*.go
rm -f CoditorX

## linux
os=linux

export GOOS=${os}
export GOARCH=386
go build
tar zcvf ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.tar.gz ${list} CoditorX --exclude-vcs --exclude conf/*.go --exclude i18n/*.go
rm -f CoditorX

export GOOS=${os}
export GOARCH=amd64
go build
tar zcvf ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.tar.gz ${list} CoditorX --exclude-vcs --exclude conf/*.go --exclude i18n/*.go
rm -f CoditorX

## windows
os=windows

export GOOS=${os}
export GOARCH=386
go build
zip -r ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.zip ${list} CoditorX.exe --exclude=conf/*.go --exclude=i18n/*.go
rm -f CoditorX.exe

export GOOS=${os}
export GOARCH=amd64
go build
zip -r ${target}/CoditorX-${ver}-${GOOS}-${GOARCH}.zip ${list} CoditorX.exe --exclude=conf/*.go --exclude=i18n/*.go
rm -f CoditorX.exe
