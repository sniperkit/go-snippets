COMPS=lacp\
	stp\
	lldp
BUILD_DIR=out/bin
DESTDIR=$(SR_CODE_BASE)/snaproute/src/$(BUILD_DIR)

IPCS=lacp\
	stp\
	lldp

all: ipc exe install 

exe: $(COMPS)
	 $(foreach f,$^, make -C $(f) exe DESTDIR=$(DESTDIR)/$(EXE_DIR) GOLDFLAGS="-r /opt/flexswitch/sharedlib";)

ipc: $(IPCS)
	 $(foreach f,$^, make -C $(f) ipc;)

clean: $(COMPS)
	$(foreach f,$^, make -C $(f) clean;)

install:
	$(foreach f,$^, make -C $(f) install;)
