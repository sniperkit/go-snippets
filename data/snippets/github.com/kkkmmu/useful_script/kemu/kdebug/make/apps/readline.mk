SOURCE=$(APPS_DIR)/readline/readline-$(VERSION)
VERSION="6.1"

all: build install 

build: config 
	@cd $(SOURCE) && make 

config: 
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./configure --prefix=$(ROOT_DIR) --host=$(ARCH)-unknown-linux-gnu CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
	fi

install:
	@cd $(SOURCE) && make install

clean:
	@cd $(SOURCE) && make clean && make distclean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
