include .sdk/Makefile

test-native:
	cd native; \
	go test -v .

build-native:
	cd native; \
	go build -o $(BUILD_PATH)/native parse.go; \
	chmod +x $(BUILD_PATH)/native
