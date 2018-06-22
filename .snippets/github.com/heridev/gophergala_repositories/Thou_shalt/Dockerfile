FROM ubuntu
MAINTAINER Mitch Ruebush <mruebush@gmail.com>
RUN apt-get update -y && apt-get install -y -q curl build-essential ca-certificates git mercurial

# (set an explicit GOARM of 5 for maximum compatibility)
ENV GOARM 5
ENV GOOS    linux
ENV GOARCH  amd64

# Install Go
# Save the SHA1 checksum from http://golang.org/dl
# RUN echo '9f9dfcbcb4fa126b2b66c0830dc733215f2f056e go1.4.src.tar.gz' > go1.4.src.tar.gz.sha1
RUN curl -O -s https://storage.googleapis.com/golang/go1.4.src.tar.gz
# RUN sha1sum —check go1.4.src.tar.gz.sha1
RUN tar -xzf go1.4.src.tar.gz -C /usr/local
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
RUN cd /usr/local/go/src && ./make.bash —no-clean 2>&1

# Get the code coverage tool and godep
RUN go get golang.org/x/tools/cmd/cover github.com/tools/godep

WORKDIR /go/src/github.com/mruebush/thoushalt/thoushalt
# Build the thoushalt application
RUN mkdir -p /go/src/github.com/mruebush/thoushalt
ADD . /go/src/github.com/mruebush/thoushalt
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
RUN cp thoushalt /thoushalt
# ENV PORT 80
# EXPOSE 80
# ENTRYPOINT ["/"]