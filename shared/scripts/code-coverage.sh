#!/bin/sh
#Execute all test files and echo the code coverage of each package (and global)

set -e

#test require git to work
which git

echo 'mode: set' > coverage.cov
go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.cov'
rm coverage.tmp

go tool cover -func=coverage.cov
go tool cover -html=coverage.cov -o=coverage.html