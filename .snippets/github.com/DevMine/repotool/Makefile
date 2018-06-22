PKG  = github.com/DevMine/repotool
EXEC1 = repotool
EXEC2 = repotool-db
VERSION = 1.0.0
DIR1 = ${EXEC1}-${VERSION}
DIR2 = ${EXEC2}-${VERSION}

all: check test build

install:
	go install ${PKG}/cmd/repotool
	go install ${PKG}/cmd/repotool-db

build:
	go build ${PKG}/cmd/repotool
	go build ${PKG}/cmd/repotool-db

test:
	go test ${PKG}/...
	
package: clean deps build
	test -d ${DIR1} || mkdir ${DIR1}
	cp ${EXEC1} ${DIR1}/
	cp README.md ${DIR1}/
	tar czvf ${DIR1}.tar.gz ${DIR1}
	rm -rf ${DIR1}
	test -d ${DIR2} || mkdir ${DIR2}
	cp ${EXEC2} ${DIR2}/
	cp README.md ${DIR2}/
	cp repotool.conf.sample ${DIR2}/
	cp -r db ${DIR2}/
	tar czvf ${DIR2}.tar.gz ${DIR2}
	rm -rf ${DIR2}


# FIXME: we shall compile libgit2 statically with git2go to prevent libgit2
# from being a dependency to run repotool
deps:
	go get -u github.com/golang/glog
	go get -u github.com/libgit2/git2go
	go get -u github.com/lib/pq
	go get -u -f github.com/DevMine/srcanlzr/src

dev-deps:
	go get -u github.com/golang/lint/golint

check:
	go vet ${PKG}/...
	golint ${PKG}/...

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./repotool ./repotool-db
