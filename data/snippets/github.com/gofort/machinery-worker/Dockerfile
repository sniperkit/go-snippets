FROM golang

ADD . /go/src/github.com/gofort/machinery-worker

RUN cd /go/src/github.com/gofort/machinery-worker && go get && go install github.com/gofort/machinery-worker

ENTRYPOINT /go/bin/machinery-worker