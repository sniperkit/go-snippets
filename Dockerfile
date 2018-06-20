FROM ubuntu:16.04

MAINTAINER Chris Tava <chris1tava@gmail.com>

#Begin: install prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        curl \
        git \
        mercurial \
        libcurl3-dev \
        libfreetype6-dev \
        libpng12-dev \
        libzmq3-dev \
        locate \
        pkg-config \
        rsync \
        software-properties-common \
        sudo \
        unzip \
        zip \
        zlib1g-dev \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
#End: install prerequisites

#Begin: install golang
ENV GOLANG_VERSION 1.10.1
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_SHA256_CHECKSUM 72d820dec546752e5a8303b33b009079c15c2390ce76d67cf514991646c6127b
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin:/usr/local/go/bin
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz && \
    echo "$GOLANG_SHA256_CHECKSUM golang.tar.gz" | sha256sum -c - && \
    sudo tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#End: install golang

#Begin: install tensorflow library
ENV TENSORFLOW_LIB_GZIP libtensorflow-cpu-linux-x86_64-1.7.0.tar.gz
ENV TARGET_DIRECTORY /usr/local
RUN  curl -fsSL "https://storage.googleapis.com/tensorflow/libtensorflow/$TENSORFLOW_LIB_GZIP" -o $TENSORFLOW_LIB_GZIP && \
     tar -C $TARGET_DIRECTORY -xzf $TENSORFLOW_LIB_GZIP && \
     rm -Rf $TENSORFLOW_LIB_GZIP
ENV LD_LIBRARY_PATH $TARGET_DIRECTORY/lib
ENV LIBRARY_PATH $TARGET_DIRECTORY/lib
RUN go get github.com/tensorflow/tensorflow/tensorflow/go
#End: install tensorflow library

#Begin: install protoc
ENV PROTOC_VERSION 3.5.1
ENV PROTOC_LIB_ZIP protoc-$PROTOC_VERSION-linux-x86_64.zip
ENV TARGET_DIRECTORY /usr/local
RUN  curl -fsSL "https://github.com/google/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_LIB_ZIP" -o $PROTOC_LIB_ZIP && \
     sudo unzip $PROTOC_LIB_ZIP -d $TARGET_DIRECTORY && \
     rm -Rf $TPROTOC_LIB_ZIP
#End: install protoc

#Begin: generate go files from protobufs
RUN go get github.com/golang/protobuf/proto
RUN go get github.com/golang/protobuf/protoc-gen-go
ENV TF_DIR /go/src/github.com/tensorflow/tensorflow
ENV TF_PB_DIR /go/src/github.com/tensorflow/tensorflow/tensorflow/go/pb
RUN mkdir -p $TF_PB_DIR
RUN /usr/local/bin/protoc -I $TF_DIR \
  --go_out=$TF_PB_DIR \
  $TF_DIR/tensorflow/core/framework/*.proto
#End: generate go files from protobufs

#Begin: add supplemental tensorflow go files to proper location
ADD ./tensorflow/graphio.go /go/src/github.com/tensorflow/tensorflow/tensorflow/go/graphio.go
#End: add supplemental tensorflow go files to proper location

#Begin: install gonum
RUN go get github.com/gonum/floats
#RUN go get github.com/gonum/plot //go: missing Mercurial command
#End: install gonum

#Begin: install tfcgo
RUN go get github.com/pkg/errors
RUN go get github.com/kniren/gota/dataframe
RUN go get github.com/ctava/tfcgo
#End: install tfcgo

#Begin: install delve
RUN go get github.com/derekparker/delve/cmd/dlv
#End: install delve

WORKDIR "/go/src/github.com/ctava/tfcgo/examples"
