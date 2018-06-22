
.PHONY: all
all: runtime

.PHONY: clean
clean:
	docker rmi -f smizy/keras-tensorflow:${TAG} || :

.PHONY: runtime
runtime:
	docker build \
		--build-arg BUILD_DATE=${BUILD_DATE} \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg VERSION=${VERSION} \
		--build-arg "EXTRA_BAZEL_ARGS=${EXTRA_BAZEL_ARGS}" \
		--rm -t smizy/keras-tensorflow:${TAG} .
	docker images | grep keras-tensorflow

.PHONY: test
test:
	docker run smizy/keras-tensorflow:${TAG} python -c 'import keras'