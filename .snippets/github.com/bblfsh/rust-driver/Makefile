include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

RUSTUP_CMD := rustup run $(RUNTIME_NATIVE_VERSION)
CARGO_CMD := $(RUSTUP_CMD) cargo

test-native-internal:
	cd native; \
	$(RUN_VERBOSE) $(CARGO_CMD) test

build-native-internal:
	$(header)
	cd native; \
	$(CARGO_CMD) install; \
	$(CARGO_CMD) build --release; \
	cp target/release/rust-parser $(BUILD_PATH)/bin/native
