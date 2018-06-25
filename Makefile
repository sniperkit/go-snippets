#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = storage

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)


test:
	if [ "$(s3)" != "" ]; then go test -s3; else go test; fi

format:
	gofmt -w storage.go
	gofmt -w storage_test.go

lint:
	golint storage.go
	golint storage_test.go

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)


