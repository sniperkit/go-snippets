ifeq (${OS},Windows_NT)
	EXE := .bin/srclib-quickstart.exe
else
	EXE := .bin/srclib-quickstart
endif

.PHONY: install clean

default: govendor install

install: ${EXE}

clean:
	rm -f ${EXE}

govendor:
	go get github.com/kardianos/govendor
	govendor sync

${EXE}: $(shell /usr/bin/find . -type f -and -name '*.go' -not -path './vendor/*')
	go build -o ${EXE}
