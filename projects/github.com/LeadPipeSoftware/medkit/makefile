VERSION=0.0.1-alpha
DATE=`date +%FT%T%z`
COMMIT=`git describe --always`
BIN_DIR=bin

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${DATE} -X main.commit=${COMMIT}"

build:
	go build -o ${BIN_DIR}/medkit ${LDFLAGS}

install:
	go install ${LDFLAGS}

clean:
	if [ -f medkit ] ; then rm medkit ; fi

.PHONY: clean install
