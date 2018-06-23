SOURCE=$(APPS_DIR)/ncurses/ncurses-$(VERSION)
VERSION=5.4

all: build install 

build: config 
	@cd $(SOURCE) && make 

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./configure --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
	    sed -i 's:MANS = .*:MANS = :' $(SOURCE)/Makefile; \
	fi


install:
	find $(SOURCE)/src/ -perm 775 -a ! -name ".deps" -a ! -type d | xargs -i $(INSTALL) {} $(ROOT_DIR)/bin/

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
