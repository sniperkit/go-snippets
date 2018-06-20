all:	build

clean:
	rm -f mlfmt coverage.txt
	
format:
	gofmt -w .
	gofmt -w json/
	gofmt -w yaml/
	gofmt -w toml/

test:
	golint -set_exit_status .
	golint -set_exit_status json
	golint -set_exit_status yaml
	golint -set_exit_status toml
	go vet ./...
	errcheck ./...
	go test ./... -v -covermode=atomic

build: clean format test
	go build

.PHONY: clean format test build
