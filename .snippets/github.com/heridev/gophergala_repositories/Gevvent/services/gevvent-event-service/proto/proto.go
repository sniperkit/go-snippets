//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. attendees/attendees.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. create/create.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. delete/delete.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. invite/invite.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. list/list.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. newest/newest.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. read/read.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. readuser/readuser.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. rsvp/rsvp.proto
//go:generate protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. search/search.proto

package proto
