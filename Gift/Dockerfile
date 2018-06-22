# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.4.1-onbuild

RUN mkdir -p /go/src/github.com/gophergala/; ln -s /go/src/app/ /go/src/github.com/gophergala/Gift

CMD cd /go/src/app/cmd; go run main.go
# Document that the service listens on port 8080.
EXPOSE 8080
