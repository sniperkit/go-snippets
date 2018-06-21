SOURCE=$(KERNEL_DIR)/linux

all: show config-kernel build-kernel

show:
	@echo "--------------------------------------------------------------" 
	@echo "" 
	@echo "			Building "linux-$(VERSION)"		"	
	@echo ""  
	@echo "--------------------------------------------------------------" 

config: config-kernel

config-kernel: $(BUILD_DIR)/.config

build-kernel:
	@cd $(SOURCE) && $(MAKE) O=$(BUILD_DIR)	ARCH=$(ARCH) CROSS_COMPILE=$(CROSS_COMPILE) bzImage

$(BUILD_DIR)/.config:
	@if [ ! -f "$(BUILD_DIR)/.config" ]; then \
		if [ $(ARCH)x = "x86_64"x ]; then \
			cd $(SOURCE) && $(MAKE) O=$(BUILD_DIR) $(ARCH)_defconfig; \
		elif [ $(ARCH)x = "arm"x ]; then \
			cd $(SOURCE) && $(MAKE) O=$(BUILD_DIR) versatile_defconfig; \
		else \
			echo "Unsupported ARCH: $(ARCH)" && exit -1; \
		fi \
	fi
	sed -i 's:CONFIG_INITRAMFS_SOURCE=.*:CONFIG_INITRAMFS_SOURCE="$(PROJECT_DIR)/rootfs":' $(BUILD_DIR)/.config

clean:
	@echo "--------------------------------------------------------------" 
	@echo "" 
	@echo "			Cleaning "linux-$(VERSION)"		"	
	@echo ""  
	@echo "--------------------------------------------------------------" 
	@cd $(SOURCE) && $(MAKE) O=$(BUILD_DIR)	clean

.PHONY: rootfs all apps modules kernel
