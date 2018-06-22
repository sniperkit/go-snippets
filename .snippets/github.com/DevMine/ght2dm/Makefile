PKG  = github.com/DevMine/ght2dm
EXEC = ght2dm

all: build

install:
	go install ${PKG}

build:
	go build -o ${EXEC} ${PKG}

deps:
	go get -u labix.org/v2/mgo/bson

check:
	go vet ${PKG}
	golint ${PKG}

clean:
	rm -f ./${EXEC}

