#!/bin/bash

export BINDIR=/opt/driver/bin/default/rel/pre-alpha/erts-8.1/bin;
ERLEXEC=/opt/driver/bin/default/rel/pre-alpha/erts-8.1/bin/erlexec;
BOOT=/opt/driver/bin/default/rel/pre-alpha/releases/0.0.1/pre-alpha;
ERTSLIB=/opt/driver/bin/default/rel/pre-alpha/lib;
CONFIG=/opt/driver/bin/default/rel/pre-alpha/releases/0.0.1/sys.config;
ARGS=/opt/driver/bin/default/rel/pre-alpha/releases/0.0.1/vm.args;

exec $ERLEXEC +d -noshell -boot $BOOT -boot_var ERTS_LIB_DIR $ERTSLIB 2> /dev/null
