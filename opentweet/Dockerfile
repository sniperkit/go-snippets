FROM golang:1.4
WORKDIR /go/src/github.com/gophergala/opentweet
ADD . /go/src/github.com/gophergala/opentweet
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"]
