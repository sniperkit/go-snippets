FROM ubuntu:14.04
MAINTAINER Quinn Slack <sqs@sourcegraph.com>

RUN apt-get update -qq
RUN apt-get install -qq git curl

# Install Go
RUN curl -Ls https://golang.org/dl/go1.4.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV GOROOT /usr/local/go
ENV GOBIN /usr/local/bin
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /srclib

# Install godep
RUN apt-get install -qq build-essential
RUN go get github.com/tools/godep

# Add this toolchain
ADD . /srclib/src/sourcegraph.com/sourcegraph/srclib-docker/
WORKDIR /srclib/src/sourcegraph.com/sourcegraph/srclib-docker
RUN godep go install

RUN useradd -ms /bin/bash srclib
RUN mkdir /src
RUN chown -R srclib /src /srclib
USER srclib

# Now set the GOPATH for the project source code, which is mounted at /src.
ENV GOPATH /
WORKDIR /src

ENTRYPOINT ["srclib-docker"]
