# APP Dockerfile
# For more control, you can copy and build manually
FROM golang:alpine

LABEL Name=userservice Version=0.0.1

RUN mkdir -p /go/src \
    && mkdir -p /go/bin \
    && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/app

ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app

RUN go build -o main .
RUN pwd && ls
CMD ["/go/src/app/main"]

EXPOSE 3000

