go get github.com/golang/protobuf/proto
go get github.com/golang/protobuf/protoc-gen-go
export TF_DIR=/go/src/github.com/tensorflow/tensorflow
export TF_PB_DIR=/go/src/github.com/tensorflow/tensorflow/tensorflow/go/pb
mkdir -p $TF_PB_DIR
/usr/local/bin/protoc -I $TF_DIR \
  --go_out=$TF_PB_DIR \
  $TF_DIR/tensorflow/core/framework/*.proto