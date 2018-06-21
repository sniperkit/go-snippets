go build -gcflags="-S -N"   // go tool compile -help
go build -x   //  see all the invocations
go test -race   //  race detector
go test -run=Func1   // run TestFunc1 
go get -u   // update
go get -d   // clone
go get -t   // get deps for test
go list -f  // list packages with a custom format
`
