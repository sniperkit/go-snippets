#!/bin/sh
ln /dev/null /dev/raw1394 2>/dev/null
exec /usr/bin/alpr "$@"
