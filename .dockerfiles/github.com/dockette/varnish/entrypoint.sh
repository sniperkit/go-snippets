#!/bin/bash

varnishd -a :80 -f ${VARNISH_CONFIG} -S ${VARNISH_SECRET} -s malloc,${VARNISH_CACHE}
varnishlog