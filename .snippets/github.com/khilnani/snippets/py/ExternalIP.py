#!/usr/bin/env python

import urllib
import re

def get_external_ip():
    site = urllib.urlopen("http://checkip.dyndns.org/").read()
    grab = re.findall('([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)', site)
    address = grab[0]
    return address

if __name__ == '__main__':
	print( get_external_ip() )

