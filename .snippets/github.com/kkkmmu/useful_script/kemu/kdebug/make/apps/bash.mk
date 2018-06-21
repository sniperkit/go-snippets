SOURCE=$(APPS_DIR)/bash

all: build install 

build: config
	@cd $(SOURCE) && make

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./configure --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
	fi

install:
	find $(SOURCE)/src/ -perm 775 -a ! -name ".deps" -a ! -type d | xargs -i $(INSTALL) {} $(ROOT_DIR)/bin/
	#$(INSTALL) $(SOURCE)/bash $(ROOT_DIR)/bin/bash
	#$(STRIP) $(ROOT_DIR)/bin/bash
	@cd $(ROOT_DIR)/bin && ln -sf bash sh

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
