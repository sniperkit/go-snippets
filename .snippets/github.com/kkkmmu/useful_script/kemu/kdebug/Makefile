.EXPORT_ALL_VARIABLES:

PROJECT_DIR	=	$(PWD)
MAKE_DIR	=	$(PROJECT_DIR)/make
APPS_DIR	=	$(PROJECT_DIR)/apps
MODULES_DIR	=	$(PROJECT_DIR)/modules
KERNEL_DIR	=	$(PROJECT_DIR)/kernel
SCRIPT_DIR	=	$(PROJECT_DIR)/script

BUILD_DIR	=	$(PROJECT_DIR)/result/build
ROOT_DIR	=	$(PROJECT_DIR)/rootfs
MODULES_DIR	=	$(PROJECT_DIR)/modules

ifeq ($(ARCH),)
    $(error Please set the target platform!) 
endif

all: prepare config build 

prepare: 
	@case $(ARCH) in  \
		x86_64 | arm)	\
			ln -sf $(PROJECT_DIR)/make/rule/$(ARCH).mk $(PROJECT_DIR)/make/.rule;;  \
		*) 	\
			echo "Unsupported archetecture!";;  \
	esac

build:
	$(MAKE) -C $(PROJECT_DIR)/make 

modules:
	$(MAKE) -C $(PROJECT_DIR)/make $@

kernel:
	$(MAKE) -C $(PROJECT_DIR)/make $@

config:
	$(MAKE) -C $(PROJECT_DIR)/make $@

lib:
	$(MAKE) -C $(PROJECT_DIR)/make $@

reset:
	$(MAKE) -C $(PROJECT_DIR)/make $@

boot: host1

mboot: host1 host2 router

host1 host2 router: 
	$(MAKE) -C $(PROJECT_DIR)/make $@
	sh $(SCRIPT_DIR)/boot.sh $@

clean:
	$(MAKE) -C $(PROJECT_DIR)/make $@

distclean:
	$(MAKE) -C $(PROJECT_DIR)/make $@

.PHONY: all build clean prepare modules make kernel
