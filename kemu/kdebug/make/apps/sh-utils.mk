SOURCE=$(APPS_DIR)/sh-utils

#all: build install 
all:  
	echo "This has problem"

build: config 
	@cd $(SOURCE) && make 

config:
	echo "This has problem"

install:
	sudo $(INSTALL) $(SOURCE)/env $(ROOT_DIR)/bin/env

clean:  
	echo "This has problem"

#	@cd $(SOURCE) && make clean

distclean:
	echo "This has problem"

#	@cd $(SOURCE) && make distclean

.PHONY: config build clean install
