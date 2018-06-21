SOURCE=$(APPS_DIR)/iproute2

all: build install 

build: config 
	@cd $(SOURCE) && $(MAKE)

config:
	cd $(SOURCE) && ./configure --host=$(ARCH) CC=$(CC) CFLAGS=$(CFLAGS) LDFLAGS=$(LDFLAGS)
	sed -i '/TARGETS/s/ arpd//g' $(SOURCE)/misc/Makefile 

install:
	@cd $(SOURCE) && $(MAKE) DESTDIR=$(ROOT_DIR) install

clean:
	@cd $(SOURCE) && $(MAKE) clean

distclean:
	@cd $(SOURCE) && $(MAKE) distclean


.PHONY: build config install all clean
