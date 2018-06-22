PKG  = github.com/DevMine/fluxio
EXEC = fluxio
VERSION = 1.0.0
DIR = ${EXEC}-${VERSION}

all: check test build

install:
	go install ${PKG}

build:
	go build -o ${EXEC} ${PKG}

test:
	go test -v ${PKG}/...

package: clean deps build
	test -d ${DIR} || mkdir ${DIR}
	cp ${EXEC} ${DIR}/
	cp README.md ${DIR}/
	cp fluxio.conf.sample ${DIR}/
	tar czvf ${DIR}.tar.gz ${DIR}
	rm -rf ${DIR}

deps:
	go get -u github.com/lib/pq

dev-deps:
	go get -u github.com/golang/lint/golint

check:
	go vet ${PKG}/...
	golint ${PKG}/...

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./${EXEC}
