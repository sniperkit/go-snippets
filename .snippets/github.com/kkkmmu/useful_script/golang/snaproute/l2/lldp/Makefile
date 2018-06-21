RM=rm -f
RMFORCE=rm -rf
GOLDFLAGS=-r /opt/flexswitch/sharedlib
DESTDIR=$(SR_CODE_BASE)/snaproute/src/out/bin
GENERATED_IPC=$(SR_CODE_BASE)/generated/src
IPC_GEN_CMD=thrift
IPC_SRCS=flexswitch/lldpd.thrift
SRCS=main.go
COMP_NAME=lldpd
all:ipc exe
ipc:
	 $(IPC_GEN_CMD) --gen go -out $(GENERATED_IPC) $(IPC_SRCS)

exe: $(SRCS)
	 go build -o $(DESTDIR)/$(COMP_NAME) -ldflags="$(GOLDFLAGS)" $(SRCS)

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

install:
	@echo "LLDP has no files to install"

clean:guard
	 $(RM) $(DESTDIR)/$(COMP_NAME)
	 $(RMFORCE) $(GENERATED_IPC)/$(COMP_NAME)
