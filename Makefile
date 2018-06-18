.PHONY: all test

all:
	./validate.sh
	go build github.com/uoregon-libraries/gopkg/...

test:
	go test ./...
