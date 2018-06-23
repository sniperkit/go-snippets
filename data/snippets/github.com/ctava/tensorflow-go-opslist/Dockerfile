FROM tensorflow/tensorflow

#Begin: install git
RUN apt-get update && apt-get install -y --no-install-recommends git
#End: install git

#Begin: install golang
ENV GOLANG_VERSION 1.8.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_SHA256_CHECKSUM 1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin:/usr/local/go/bin
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz && \
    echo "$GOLANG_SHA256_CHECKSUM golang.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#End: install golang

#Begin: install tensorflow library
ENV TENSORFLOW_LIB_GZIP libtensorflow-cpu-linux-x86_64-1.2.0.tar.gz
ENV TARGET_DIRECTORY /usr/local
RUN  curl -fsSL "https://storage.googleapis.com/tensorflow/libtensorflow/$TENSORFLOW_LIB_GZIP" -o $TENSORFLOW_LIB_GZIP && \
     sudo tar -C $TARGET_DIRECTORY -xzf $TENSORFLOW_LIB_GZIP && \
     rm -Rf $TENSORFLOW_LIB_GZIP
ENV LD_LIBRARY_PATH $TARGET_DIRECTORY/lib
ENV LIBRARY_PATH $TARGET_DIRECTORY/lib
RUN go get github.com/tensorflow/tensorflow/tensorflow/go
#End: install tensorflow library

#Begin: install protoc
ENV PROTOC_VERSION 3.3.0
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

#Begin: install tensorflow-go files
RUN go get github.com/ctava/tensorflow-go-version
RUN mkdir -p /go/src/github.com/ctava/tensorflow-go-opslist
ADD main.go /go/src/github.com/ctava/tensorflow-go-opslist/main.go
WORKDIR "/go"
#End: install tensorflow-go files