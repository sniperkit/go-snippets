VERSION?=$(shell git describe --tags --always --dirty)
PROTOC=/usr/local/bin/protoc

.PHONY: proto

gen:
	@${PROTOC} -I=/usr/include -I=. --go_out=${GOPATH}/src *.proto
	@python -m grpc_tools.protoc -I=/usr/include -I=. --python_out=./qtypes qtypes.proto
	@ls -al | grep "pb.go"
	@ls -al ./qtypes | grep "_pb2"

version:
	@echo ${VERSION} > VERSION.txt