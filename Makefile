GOIMPORTS_FILES?=$$(find . -name '*.go')

all: goimportscheck vet megacheck test

cover: goimportscheck vet
	@sh -c "'$(CURDIR)/scripts/cover.sh'"

megacheck:
	@megacheck ./...

test:
	@go test ./...

goimports:
	@goimports -w $(GOIMPORTS_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

vet:
	@go list -f '{{.Dir}}' ./... \
		| xargs go tool vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
		fi

bootstrap:
	@sh -c "'$(CURDIR)/scripts/bootstrap.sh'"

.PHONY: all cover test goimports goimportscheck vet bootstrap
