-include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

DIR="_build/default/rel"
RELEASE="pre-alpha"

test-native-internal:
	cd native; \
	rebar3 eunit;
	if [ -d .cache ]; then rm -r .cache; fi;

build-native-internal:

	rm -rf build; \
	cd native; \
	rm -rf _build; \
	rebar3 release; \
	awk '{gsub("-noinput","");print}' $(DIR)/$(RELEASE)/bin/$(RELEASE) > $(DIR)/$(RELEASE)/bin/tmp;\
	cat $(DIR)/$(RELEASE)/bin/tmp > $(DIR)/$(RELEASE)/bin/$(RELEASE); \
  rm $(DIR)/$(RELEASE)/bin/tmp; \
	sed -i '/echo "Exec:/ d' $(DIR)/$(RELEASE)/bin/$(RELEASE) && \
	sed -i '/echo "Root:/ d' $(DIR)/$(RELEASE)/bin/$(RELEASE) && \
	sed -i '547 d' $(DIR)/$(RELEASE)/bin/$(RELEASE) && \
	cp -r _build $(BUILD_PATH); \
  cp native.sh $(BUILD_PATH)/native; \
	chmod +x $(BUILD_PATH)/native;
	if [ -d .cache ]; then rm -r .cache; fi;
