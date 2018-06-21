CFLAGS += -O2 -nostdinc
CFLAGS += -I.
CFLAGS += -I..
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/usr/include
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/include
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/include
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/lib/gcc/$(toolchain)/5.1.0/plugin/include
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/lib/gcc/$(toolchain)/5.1.0/include
CFLAGS += -I$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/lib/gcc/$(toolchain)/5.1.0/install-tools/include

LDFLAGS += -L.
LDFLAGS += -L..
LDFLAGS += -L$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/lib
LDFLAGS += -L$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/lib32
LDFLAGS += -L$(PROJECT_DIR)/tools/toolchain/$(ARCH)/$(toolchain)/$(toolchain)/sysroot/lib64
