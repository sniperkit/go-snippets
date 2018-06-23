APP_MAKE = $(MAKE_DIR)/apps

#APPS = $(basename $(notdir $(wildcard $(APP_MAKE)/*.mk)))
APPS = $(shell  sed -n 's/enable_\(.*\)=yes/\1/p' ./apps.rule)

all: apps

apps: config
	@for i in $(APPS); do  \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Building "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -f apps/$$i.mk || exit 1; \
	done	

config:
	@for i in $(APPS); do  \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Configuring "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -f apps/$$i.mk $@ || exit 1; \
	done	

clean: config
	for i in $(APPS); do \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Cleaning "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -f apps/$$i.mk $@ || exit 1; \
	done	

distclean:
	for i in $(APPS); do \
	    $(MAKE) -f apps/$$i.mk $@ || exit 1; \
	done	

.PHONY: rootfs all apps modules kernel prepare config
