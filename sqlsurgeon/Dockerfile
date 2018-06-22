FROM golang:1.4
RUN apt-get update && \
    apt-get install -y graphviz

COPY src/github.com/hrharkins/sqlsurgeon /go/src/sqlsurgeon
WORKDIR /go/src/sqlsurgeon
RUN go-wrapper download
RUN go-wrapper install

