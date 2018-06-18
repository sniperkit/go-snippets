#!/usr/bin/env python

import sys
from xml.etree import ElementTree

def extract_rss_urls_from_opml(filename):
    urls = []
    with open(filename, 'rt') as f:
        tree = ElementTree.parse(f)
    for node in tree.findall('.//outline'):
        url = node.attrib.get('xmlUrl')
        title = node.attrib.get('title')
        if url:
            urls.append( (title, url) )
    return urls

if __name__ == '__main__':
    if len(sys.argv) > 1:
        print("Parsing {}".format(sys.argv[1]))
        urls = extract_rss_urls_from_opml( sys.argv[1] )
        for url in urls:
            print(url)
