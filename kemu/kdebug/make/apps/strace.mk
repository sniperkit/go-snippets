SOURCE=$(APPS_DIR)/strace

all: build install 

build: config 
	@cd $(SOURCE) && make 

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./configure --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
	fi

install:
	$(INSTALL) $(SOURCE)/strace $(ROOT_DIR)/bin/strace
	$(STRIP) $(ROOT_DIR)/bin/strace

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean


.PHONY: build config install all clean
