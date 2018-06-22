# GNUmakefile

REPOSITORY = dweomer/consul
VERSIONS   = $(sort $(notdir $(realpath $(wildcard versions/*))))
TAGS       = $(filter-out $(VERSIONS),$(sort $(notdir $(wildcard versions/*))))

export REPOSITORY VERSIONS TAGS

all: $(VERSIONS) $(TAGS)

build:
ifndef VERSION
	$(error VERSION is required)
endif
	docker build \
		--build-arg CONSUL_VERSION=$(VERSION) \
		--pull \
		--tag $(REPOSITORY):$(VERSION) \
		.

$(TAGS): $(VERSIONS)
	docker tag $(REPOSITORY):$(notdir $(realpath versions/$@)) $(REPOSITORY):$@

$(VERSIONS):
	@make build VERSION=$@

.PHONY: all build $(VERSIONS) $(TAGS) push
