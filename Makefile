TEST?=$$(go list ./... | grep -v '/vendor/')
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

docker-clean:
	docker rmi `docker images -aq dwrap-image/*`

build: clean vet
	govendor build -o $(CURDIR)/bin/dwrap $(CURDIR)/cmd/main.go

build-x: clean vet
	sh -c "'$(CURDIR)/scripts/build.sh'"

docker-build: clean vet
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"

test: vet
	govendor test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

vet: fmt
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)


.PHONY: default test vet testacc fmt fmtcheck
