all: compile-prepare

compile-prepare:
	@if [ ! -d $(BUILD_DIR) ]; then \
	    $(MKDIR) $(BUILD_DIR); \
	fi
	@if [ ! -d $(ROOT_DIR) ]; then \
	    $(MKDIR) $(ROOT_DIR); \
	fi
	$(MKDIR) -p $(ROOT_DIR)/sbin 
	$(MKDIR) -p $(ROOT_DIR)/bin
	$(MKDIR) -p $(ROOT_DIR)/usr/bin
	$(MKDIR) -p $(ROOT_DIR)/usr/sbin

clean:
	$(RM) -rf $(ROOT_DIR)
	$(RM) -rf $(BUILD_DIR)

.PHONY: all clean
