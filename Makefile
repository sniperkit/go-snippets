-include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

NATIVE_BIN := driver.native

test-native-internal:
	cd native; \
	echo "not implemented"

build-native-internal:
	cd native; \
	eval `opam config env --root=/opt/driver/opam`; \
	ocamlbuild -tag debug -use-ocamlfind $(NATIVE_BIN); \
	cp $(NATIVE_BIN) $(BUILD_PATH)/bin/native
