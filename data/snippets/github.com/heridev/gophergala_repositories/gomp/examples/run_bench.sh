#!/bin/bash
gomp < bench.go > bench1.go
go run bench1.go
#./bench1
rm -rf bench1.go bench1
