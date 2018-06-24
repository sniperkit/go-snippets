#! /bin/sh

linuxRelease=github.linux

GOOS=linux GOARCH=amd64 go build -o $linuxRelease
