include .sdk/Makefile

test-native-internal:
	cd native; \
	lein test

build-native-internal:
	cd native; \
	lein uberjar; \
	mv ./target/uberjar/*-standalone.jar $(BUILD_PATH)/bin/; \
	echo -e "#!/bin/bash\njava -jar /opt/driver/bin/*-standalone.jar" > $(BUILD_PATH)/bin/native; \
	chmod +x $(BUILD_PATH)/bin/native
