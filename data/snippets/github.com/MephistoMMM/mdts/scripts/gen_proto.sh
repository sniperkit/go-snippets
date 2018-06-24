#!/bin/bash
# This script generate golang code according to proto file

PROJECT_PATH=~/Desktop/work/school/graduate_design/src

if which protoc 2>&1 > /dev/null ; then
    # generate brokermsg
    protoc -I $PROJECT_PATH --go_out=plugins=grpc:src mdts/protocols/brokermsg/brokermsg.proto

    # generate s2t/broker
    protoc -I $PROJECT_PATH  --go_out=plugins=grpc:src mdts/brokerSDK/s2t/broker/broker.proto

    # generate t2s/broker
    protoc -I $PROJECT_PATH  --go_out=plugins=grpc:src mdts/brokerSDK/t2s/broker/broker.proto
else
    echo "Not exist command protoc.Please install it first."
fi
