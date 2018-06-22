MVN_CMD := mvn

test-native-internal:
	cd native; \
	$(MVN_CMD) test

build-native-internal:
	cd native; \
	$(MVN_CMD) package
	cp native/target/native-jar-with-dependencies.jar $(BUILD_PATH); \
	echo "#!/bin/sh" > $(BUILD_PATH)/native; \
	echo "java -jar native-jar-with-dependencies.jar" >> $(BUILD_PATH)/native; \
	chmod +x $(BUILD_PATH)/native


include .sdk/Makefile
