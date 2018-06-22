#!/bin/sh

for GOOS in linux freebsd openbsd netbsd dragonfly darwin
do
	for GOARCH in amd64 386
	do
		gvm cross $GOOS $GOARCH
	done
done
