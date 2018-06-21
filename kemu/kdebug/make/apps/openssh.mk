SOURCE=$(APPS_DIR)/openssh
ZLIB_SOURCE=$(APPS_DIR)/zlib
OPENSSL_SOURCE=$(APPS_DIR)/openssl

all: build install 

build: config 
	@if  [ ! -f $(SOURCE)/.openssh.p ]; then \
		cd $(SOURCE) && make && echo 1 > $(SOURCE)/.openssh.p; \
	fi
config: zlib openssl
	@if [ ! -f $(SOURCE)/Makefile ]; then \
	    cd $(SOURCE) && ./configure --build=x86_64-unknown-linux-gnu --host=x86_64 --with-zlib=$(ZLIB_SOURCE) --with-ssl-dir=$(OPENSSL_SOURCE) CFLAGS="$(CFLAGS) -I$(OPENSSL_SOURCE)/include" LDFLAGS="$(LDFLAGS) -L$(OPENSSL_SOURCE)/"; \
	fi

# Do not run configure for openssl. It's disgusting.
openssl_config:
	cd $(OPENSSL_SOURCE) && ./Configure shared linux-x86_64 --prefix=$(ROOT_DIR); 
	cd $(OPENSSL_SOURCE) && make depend; 
	sed -i 's:CC= .*:CC=$(CROSS_COMPILE)gcc :' $(OPENSSL_SOURCE)/Makefile;
	#This is necessary for some version of openssl
	#sed -i 's:AR=.*:AR=$(CROSS_COMPILE)ar rcs:' $(OPENSSL_SOURCE)/Makefile;

openssl: openssl_config
	@if  [ ! -f $(OPENSSL_SOURCE)/.openssl.p ]; then \
		cd $(OPENSSL_SOURCE) && make && echo 1 > $(OPENSSL_SOURCE)/.openssl.p; \
	fi

zlib: zlib_config
	@if  [ ! -f $(ZLIB_SOURCE)/.zlib.p ]; then \
		cd $(ZLIB_SOURCE) && make; \
		echo 1 > $(ZLIB_SOURCE)/.zlib.p; \
	fi

zlib_config:
	cd $(ZLIB_SOURCE) && ./configure 

install:
	$(INSTALL) $(SOURCE)/sshd $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/scp $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/sftp $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh-agent $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh-add $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh-keygen $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh-keysign $(ROOT_DIR)/sbin/
	$(INSTALL) $(SOURCE)/ssh-keyscan $(ROOT_DIR)/sbin/

clean:
	@cd $(SOURCE) && make clean
	@rm -rf $(ZLIB_SOURCE)/.zlib.p
	@rm -rf $(OPENSSL_SOURCE)/.openssl.p
	@rm -rf $(OPENSSH_SOURCE)/.openssh.p

distclean:
	@cd $(SOURCE) && make distclean

.PHONY: config build clean install zlib

#	@find $(SOURCE)/ -perm 775 | xargs -i sudo $(INSTALL) {} $(ROOT_DIR)/sbin
