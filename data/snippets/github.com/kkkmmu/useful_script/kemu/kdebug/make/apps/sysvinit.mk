SOURCE=$(APPS_DIR)/sysvinit

all: build install 

CFLAGS+=-D _PATH_UTMP='"/var/run/utmp"' -D _PATH_WTMP='"/var/run/wtmp"'
LDFLAGS+= -lcrypt
build: config
	@cd $(SOURCE)/src && make

config:
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE)/src && ./configure --host=$(ARCH) CFLAGS="$(CFLAGS)" LDFLAGS="$(LDFLAGS)"; \
	fi

install:
	find $(SOURCE)/src/ -perm 775 -a ! -name ".deps" -a ! -type d | xargs -i $(INSTALL) {} $(ROOT_DIR)/sbin/
	@cd $(ROOT_DIR)/sbin && ln -sf init ../init

clean:
	@cd $(SOURCE) && make clean

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
