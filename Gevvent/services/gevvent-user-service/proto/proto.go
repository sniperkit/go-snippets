//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. authorised/authorised.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. login/login.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. logout/logout.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. readuser/readuser.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. register/register.proto

package proto
