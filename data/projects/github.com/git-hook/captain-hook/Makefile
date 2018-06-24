# Written by Eric Crosson
# 2017-03-12
#
# This makefile provides targets to test captain-hook.

tests := $(wildcard test/*.sh)

.PHONY: test

test: $(tests)
	./test/libs/bats/bin/bats $(tests)
	#./test/libs/bats/bin/bats @$(foreach test,$(tests),$(test);)
