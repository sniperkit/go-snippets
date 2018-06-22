.PHONY: build run test clean dist

BINARY := entity-extractor
BUILDFILES := config.go extractor.go extractor_api.go main.go
IMPORT_BASE := github.com/alphagov
IMPORT_PATH := $(IMPORT_BASE)/entity-extractor
DISTFILES := $(BINARY) scripts/*

all: test build

dist/entity-extractor.zip: test build
	mkdir -p dist
	rm -f $@
	zip $@ $(DISTFILES)

dist: dist/entity-extractor.zip

build: _vendor
	gom build -o $(BINARY) $(BUILDFILES)

run: _vendor
	gom run $(BUILDFILES)

test_database:
	(psql -lqt | cut -d \| -f 1 | grep -q -w entity-extractor_test >/dev/null) || \
	( \
		(createdb -T template0 -E UTF8 entity-extractor_test) && \
		(psql entity-extractor_test < db/schema.sql) && \
		(psql entity-extractor_test < db/test-fixtures.sql) \
	)

drop_test_database:
	dropdb entity-extractor_test

test: _vendor test_database
	gom test -v

clean:
	rm -f $(BINARY)

_vendor: Gomfile _vendor/src/.symlink_created
	gom install
	touch $@

_vendor/src/.symlink_created:
	rm -f _vendor/src/$(IMPORT_PATH)
	mkdir -p _vendor/src/$(IMPORT_BASE)
	ln -s $(CURDIR) _vendor/src/$(IMPORT_PATH)
	touch $@
