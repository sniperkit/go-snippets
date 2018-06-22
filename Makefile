NAME=ci-release-notes
ARCH=$(shell uname -m)
VERSION=1.0.1

build:
	mkdir -p build/Linux && GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o build/Linux/$(NAME)
	cp ./get-prs ./codeship build/Linux

release: build
	rm -rf release && mkdir release
	cd build/Linux && tar -zcf ../../release/$(NAME)_$(VERSION)_linux_$(ARCH).tar.bz2 *
	gh-release create lbn/$(NAME) $(VERSION)
