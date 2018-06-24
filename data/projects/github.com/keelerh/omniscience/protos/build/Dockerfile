FROM golang:1.10

WORKDIR /tmp
RUN apt-get update && apt-get -y install unzip
RUN wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip -O protoc.zip && unzip protoc.zip && mv bin/* /usr/local/bin && mv include/* /usr/local/include/

WORKDIR /go
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/golang/protobuf/protoc-gen-go
COPY generate_protos.sh generate_protos.sh

