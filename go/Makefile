NAME = parser-go
PKG  = github.com/DevMine/parsers/go
EXEC = parser

all: build

build:
	go build -o ${EXEC} ${PKG}

package: build
	mkdir ${NAME}
	cp ${EXEC} ${NAME}
	zip -r ${NAME}.zip ${NAME}
	rm -rf ./${NAME}

clean:
	rm -f ${EXEC} ${NAME}.zip

