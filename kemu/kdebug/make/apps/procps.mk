SOURCE=$(APPS_DIR)/procps/procps-$(VERSION)
VERSION=3.2.8

all: build install 

CFLAGS+=-I$(APPS_DIR)/ncurses/ncurses-5.4 -I$(APPS_DIR)/ncurses/ncurses-5.4/include
build:  $(SOURCE)/curses.h $(SOURCE)/ncurses.h
	@cd $(SOURCE) && make 

$(SOURCE)/curses.h $(SOURCE)/ncurses.h:
	ln -sf $(APPS_DIR)/ncurses/ncurses-5.4/include/ncurses.h $(SOURCE)/curses.h 
	ln -sf $(APPS_DIR)/ncurses/ncurses-5.4/include/curses.h $(SOURCE)/ncurses.h 

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && chmod a+x ./configure && ./configure --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"  --with-ncurses-dir=$(APPS_DIR)/ncurses/ncurses-5.4;   \
	fi


install:
	find $(SOURCE)/ps/ -perm 775 -a ! -name ".deps" -a ! -type d | xargs -i $(INSTALL) {} $(ROOT_DIR)/bin/
	find $(SOURCE)/proc/ -perm 775 -a ! -name ".deps" -a ! -type d | xargs -i $(INSTALL) {} $(ROOT_DIR)/lib/

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
