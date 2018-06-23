#!/bin/sh -x

export GOPATH=/data/git/snippets/go/martini
go get github.com/go-martini/martini
go get gopkg.in/mgo.v2
go get gopkg.in/mgo.v2/bson
go get gopkg.in/mgo.v2/txn
go get github.com/martini-contrib/binding
go get github.com/martini-contrib/auth
go get github.com/martini-contrib/gzip
go get github.com/martini-contrib/render
go get github.com/martini-contrib/sessions
go get github.com/martini-contrib/strip
go get github.com/martini-contrib/method
go get github.com/martini-contrib/secure
go get github.com/martini-contrib/encoder
go get github.com/martini-contrib/cors
go get github.com/martini-contrib/oauth2
go get github.com/lib/pq
