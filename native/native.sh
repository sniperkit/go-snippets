#!/bin/sh
SCRIPT=index.js
BIN="`readlink -f $0`"
DIR="`dirname "$BIN"`"
exec node "$DIR/$SCRIPT"