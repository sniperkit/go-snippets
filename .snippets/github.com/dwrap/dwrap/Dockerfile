FROM golang:1.6.2-alpine
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN set -x && apk add --no-cache git make zip bash gcc build-base
RUN go get -u github.com/kardianos/govendor

ADD . $GOPATH/src/github.com/dwrap/cli

WORKDIR $GOPATH/src/github.com/dwrap/cli
RUN make build

ENTRYPOINT $GOPATH/src/github.com/dwrap/cli/bin/dwrap
