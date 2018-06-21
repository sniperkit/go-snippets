RM=rm -f
RMFORCE=rm -rf
DESTDIR=$(SR_CODE_BASE)/snaproute/src/out/bin
GENERATED_IPC=$(SR_CODE_BASE)/generated/src
IPC_GEN_CMD=thrift
SRCS=main.go
IPC_SRCS=rpc/bfdd.thrift
COMP_NAME=bfdd
GOLDFLAGS=-r /opt/flexswitch/sharedlib
all:ipc exe
ipc:
	$(IPC_GEN_CMD) -r --gen go -out $(GENERATED_IPC) $(IPC_SRCS)

exe: $(SRCS)
	go build -o $(DESTDIR)/$(COMP_NAME) -ldflags="$(GOLDFLAGS)" $(SRCS)

guard:
ifndef SR_CODE_BASE
	$(error SR_CODE_BASE is not set)
endif

install:
	@echo "BFD has no files to install"
clean:guard
	$(RM) $(DESTDIR)/$(COMP_NAME) 
	$(RMFORCE) $(GENERATED_IPC)/$(COMP_NAME)
