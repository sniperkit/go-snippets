PKG  = github.com/DevMine/srctool
EXEC = srctool
VERSION = 1.0.0
DIR = ${EXEC}-${VERSION}

all: check test build

install:
	go install ${PKG}

build:
	go build -o ${EXEC} ${PKG}

test:
	go test ${PKG}/...

check:
	go vet ${PKG}/...
	golint ${PKG}/...

package: deps build
	test -d ${DIR} || mkdir ${DIR}
	cp ${EXEC} ${DIR}/
	cp README.md ${DIR}/
	tar czvf ${DIR}.tar.gz ${DIR}
	rm -rf ${DIR}	

deps:
	go get -u -v github.com/codegangsta/cli
	go get -u -v -f github.com/DevMine/srcanlzr
	go get -u -v github.com/gilliek/go-xterm256/xterm256
	go get -u -v github.com/mitchellh/ioprogress
	go get -u -v golang.org/x/crypto/ssh/terminal
	go get -u -v -f github.com/DevMine/repotool/model

dev-deps:
	go get -u github.com/golang/lint/golint

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./${EXEC}
