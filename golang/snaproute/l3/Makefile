COMPS=arp\
      rib\
      bgp\
      ndp\
      dhcp\
      ospf\
      ospfv2\
      dhcp_relay\
      bfd \
      vrrp\
      tunnel/vxlan\
      isis


IPCS=arp\
     rib\
     bgp\
     ndp\
     dhcp\
     ospf\
     ospfv2\
     dhcp_relay\
     bfd\
     vrrp\
     tunnel/vxlan\
     isis

define timedMake
@echo -n "Building component $(1) started at :`date`\n"
make -C $(1) exe 
@echo -n "Done building component $(1) at :`date`\n\n"
endef
all: ipc exe install

exe: $(COMPS)
	@$(foreach f,$^, $(call timedMake, $(f)))

ipc: $(IPCS)
	 $(foreach f,$^, make -C $(f) ipc;)

clean: $(COMPS)
	 $(foreach f,$^, make -C $(f) clean;)

install: $(COMPS)
	 $(foreach f,$^, make -C $(f) install;)

