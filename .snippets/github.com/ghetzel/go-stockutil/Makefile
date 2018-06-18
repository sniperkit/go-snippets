.PHONY: test deps

PKGS=`go list ./... | grep -v /vendor/`
LOCALS=`find . -type f -name '*.go' -not -path "./vendor/*"`

all: fmt deps test

deps:
	@go list github.com/mjibson/esc || go get github.com/mjibson/esc/...
	go generate -x ./...
	go get ./...

clean-bundle:
	-rm -rf public

clean:
	-rm -rf bin

fmt:
	@go list golang.org/x/tools/cmd/goimports || go get golang.org/x/tools/cmd/goimports
	goimports -w $(LOCALS)
	go vet $(PKGS)

test:
	go test $(PKGS)

