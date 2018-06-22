PKG  = github.com/DevMine/srcanlzr
EXEC = srcanlzr
VERSION = 0.0b
DIR = ${EXEC}-${VERSION}

all: check test build

install:
	go install ${PKG}

build:
	go build -o ${EXEC} ${PKG}

test:
	go test ${PKG}/...

bench:
	go test -bench=. -benchmem ${PKG}/...

package: clean deps build
	test -d ${DIR} || mkdir ${DIR}
	cp ${EXEC} ${DIR}/
	cp README.md ${DIR}/
	tar czvf ${DIR}.tar.gz ${DIR}
	rm -rf ${DIR}

deps:
	go get -u golang.org/x/tools/imports
	go get -u -f github.com/DevMine/repotool/model

dev-deps:
	 go get -u github.com/golang/lint/golint
	 go get -u golang.org/x/tools/imports

check:
	go vet ${PKG}/...
	golint ${PKG}/...

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./${EXEC}
