all: etc lib

etc:
	$(call print_build,etc, dk)
	$(CP) -a $(PROJECT_DIR)/samples/rootfs/etc $(ROOT_DIR)
	$(CHMOD) a+x $(ROOT_DIR)/etc/init.d/rcS
	cd $(ROOT_DIR) && ln -sf sbin/init init
#	$(CHOWN) root -R $(ROOT_DIR)/etc
#	$(CHGRP) root -R $(ROOT_DIR)/etc

lib:
	$(call print_build,lib, kk)
	$(CP) -a $(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/lib $(ROOT_DIR)
	$(CP) -a $(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/lib64 $(ROOT_DIR)
#	$(CHOWN) root -R $(ROOT_DIR)/lib

.PHONY: all etc lib clean dir

