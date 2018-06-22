# TODO add api-server once the matrix library issue has been solved
GO_PROJECTS=crawld featscomp ght2dm repotool srcanlzr srctool
# TODO add srccat once it gets a makefile
PROJECTS=crawld featscomp ght2dm repotool srcanlzr srctool


build:
	for p in ${PROJECTS}; do GOPATH=${CURDIR}/build make -C $$p build; done

deps:
	test -d build || mkdir build
	for p in ${GO_PROJECTS}; do GOPATH=${CURDIR}/build go get -u github.com/DevMine/$$p; done
	for p in ${PROJECTS}; do GOPATH=${CURDIR}/build make -C $$p deps; done

install: build
	test -d bin || mkdir bin
	for p in ${GO_PROJECTS}; do GOPATH=${CURDIR}/build GOBIN=${CURDIR}/bin make -C $$p install; done

clean:
	rm -rf build

submodules-init:
	git submodule update --init

submodules-update:
	git submodule update --remote

.PHONY: deps build clean submodules-init submodules-update
