# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/Project-V
WORKDIR /go/src/Project-V

RUN go get

EXPOSE 8100

CMD ["go", "run", "server.go", "routes.go"] 
