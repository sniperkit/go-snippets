#!/bin/bash
# This script install grpc and its dependences.

if [ "`uname`" == "Darwin" ];then
    brew install protobuf
fi

go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
