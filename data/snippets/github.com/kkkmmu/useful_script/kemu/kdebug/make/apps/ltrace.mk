SOURCE=$(APPS_DIR)/ltrace

all: build install 

build: config 
	@cd $(SOURCE) && make 

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./autogen.sh && ./configure --host=$(ARCH)-unknown-linux-gnu CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)" --disable-werror; \
	fi


install:
	$(INSTALL) $(SOURCE)/ltrace $(ROOT_DIR)/bin/ltrace
	$(STRIP) $(ROOT_DIR)/bin/ltrace

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
