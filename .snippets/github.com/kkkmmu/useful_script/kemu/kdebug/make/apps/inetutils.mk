SOURCE=$(APPS_DIR)/inetutils

all: build install 

build: config 
	@cd $(SOURCE) && make 
config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
		cd $(SOURCE) && autoreconf;\
		cd $(SOURCE) && ./configure --disable-ifconfig --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
		sed -i 's#all-am:.*#all-am: Makefile#' $(SOURCE)/man/Makefile; \
	fi


install:
	find $(SOURCE) -name "*d" -a ! -name "*.valgrind" -a ! -name "*.sed" -a ! -type d | xargs -i sudo $(INSTALL) {} $(ROOT_DIR)/usr/sbin/

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
