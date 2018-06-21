MODULES = $(shell sed -n 's/enable_\(.*\)=yes/\1/p' ./modules.rule)

all: modules

modules: modules_build modules_install

modules_build: kernel
	@for i in $(MODULES); do  \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Building Kernel moudles "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -C $(BUILD_DIR) M=$(MODULES_DIR)/$$i modules || exit 1; \
	done

modules_install:
	@for i in $(MODULES); do  \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Installing Kernel moudles "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -C $(BUILD_DIR) M=$(MODULES_DIR)/$$i INSTALL_MOD_PATH=$(ROOT_DIR) modules_install || exit 1; \
	    done

clean:
	@for i in $(MODULES); do  \
	    echo "--------------------------------------------------------------" ; \
	    echo "" ; \
	    echo "			Cleaning Kernel moudles "$$i"		";	\
	    echo "" ; \
	    echo "--------------------------------------------------------------"; \
	    $(MAKE) -C $(BUILD_DIR) M=$(MODULES_DIR)/$$i clean || exit 1; \
	done

.PHONY: rootfs all apps modules kernel clean
